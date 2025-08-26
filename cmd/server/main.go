package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"youtube-mcp/pkg/server"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Command line flags
	useJSON := flag.Bool("json-config", false, "Use JSON config file instead of .env")
	configFile := flag.String("config", "config.json", "Configuration file path (only used with -json-config)")
	flag.Parse()

	// Load configuration
	var cfg *server.Config
	var err error
	
	if *useJSON {
		cfg, err = server.LoadConfigFromJSON(*configFile)
	} else {
		cfg, err = server.LoadConfig()
	}
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if cfg.YouTubeAPIKey == "" && !fileExists(cfg.OAuth2CredentialsFile) {
		log.Printf("Warning: No YouTube API key provided and no OAuth2 credentials file found")
		
		if *useJSON {
			log.Printf("Please set YOUTUBE_API_KEY in your config file or provide %s", cfg.OAuth2CredentialsFile)
			// Create a sample config file if it doesn't exist
			if !fileExists(*configFile) {
				if err := cfg.SaveConfig(*configFile); err != nil {
					log.Printf("Failed to create sample config file: %v", err)
				} else {
					log.Printf("Created sample configuration file: %s", *configFile)
				}
			}
		} else {
			log.Printf("Please create a .env file with YOUTUBE_API_KEY or provide %s", cfg.OAuth2CredentialsFile)
			log.Printf("Example .env file:")
			log.Printf("YOUTUBE_API_KEY=your_api_key_here")
		}
		
		log.Printf("You can create API credentials at: https://console.cloud.google.com/")
		
		os.Exit(1)
	}

	// Create YouTube client
	youtubeClient, err := server.NewYouTubeClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create YouTube client: %v", err)
	}

	// Create MCP server using official SDK
	implementation := &mcp.Implementation{
		Name:    cfg.ServerName,
		Version: cfg.ServerVersion,
	}
	
	mcpServer := mcp.NewServer(implementation, nil)
	
	// Register MCP tools using the official SDK
	if err := server.SetupOfficialMCPTools(mcpServer, youtubeClient); err != nil {
		log.Fatalf("Failed to setup MCP tools: %v", err)
	}

	// Start the server
	log.Printf("Starting %s v%s", cfg.ServerName, cfg.ServerVersion)
	log.Printf("Server description: %s", cfg.ServerDescription)
	
	// Run the server over stdin/stdout using official SDK
	ctx := context.Background()
	transport := &mcp.StdioTransport{}
	if err := mcpServer.Run(ctx, transport); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Example usage function (for documentation)
func printUsageExamples() {
	fmt.Println("YouTube MCP Server Usage Examples:")
	fmt.Println("")
	fmt.Println("1. Search for videos:")
	fmt.Println(`   {"method": "tools/call", "params": {"name": "search_videos", "arguments": {"query": "golang tutorial", "max_results": 5}}}`)
	fmt.Println("")
	fmt.Println("2. Get channel information:")
	fmt.Println(`   {"method": "tools/call", "params": {"name": "get_channel_info", "arguments": {"channel_id": "UCBUKHRdFNRo2gPXRXCKnN7w"}}}`)
	fmt.Println("")
	fmt.Println("3. Get video details:")
	fmt.Println(`   {"method": "tools/call", "params": {"name": "get_video_details", "arguments": {"video_id": "dQw4w9WgXcQ"}}}`)
	fmt.Println("")
	fmt.Println("4. Get playlist items:")
	fmt.Println(`   {"method": "tools/call", "params": {"name": "get_playlist_items", "arguments": {"playlist_id": "PLrAXtmRdnEQy4jrMVwm9wRbWqz4_0dmSh", "max_results": 10}}}`)
	fmt.Println("")
	fmt.Println("5. Search for channels:")
	fmt.Println(`   {"method": "tools/call", "params": {"name": "search_channels", "arguments": {"query": "programming", "max_results": 5}}}`)
}
