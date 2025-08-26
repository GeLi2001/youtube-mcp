package server

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration for the YouTube MCP server
type Config struct {
	// YouTube API Key - can be used for public data access
	YouTubeAPIKey string `json:"youtube_api_key"`
	
	// OAuth2 credentials file path for user-specific operations
	OAuth2CredentialsFile string `json:"oauth2_credentials_file"`
	
	// Token file path to store OAuth2 tokens
	TokenFile string `json:"token_file"`
	
	// Server configuration
	ServerName        string `json:"server_name"`
	ServerVersion     string `json:"server_version"`
	ServerDescription string `json:"server_description"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		YouTubeAPIKey:         os.Getenv("YOUTUBE_API_KEY"),
		OAuth2CredentialsFile: "client_secret.json",
		TokenFile:             "token.json",
		ServerName:            "youtube-mcp-server",
		ServerVersion:         "1.0.0",
		ServerDescription:     "YouTube Data API v3 MCP Server for video search, channel info, and more",
	}
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() (*Config, error) {
	// Try to load .env file (optional)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}
	
	config := DefaultConfig()
	
	// Load from environment variables (which now include .env values)
	if apiKey := os.Getenv("YOUTUBE_API_KEY"); apiKey != "" {
		config.YouTubeAPIKey = apiKey
	}
	if credsFile := os.Getenv("OAUTH2_CREDENTIALS_FILE"); credsFile != "" {
		config.OAuth2CredentialsFile = credsFile
	}
	if tokenFile := os.Getenv("TOKEN_FILE"); tokenFile != "" {
		config.TokenFile = tokenFile
	}
	if serverName := os.Getenv("SERVER_NAME"); serverName != "" {
		config.ServerName = serverName
	}
	if serverVersion := os.Getenv("SERVER_VERSION"); serverVersion != "" {
		config.ServerVersion = serverVersion
	}
	if serverDesc := os.Getenv("SERVER_DESCRIPTION"); serverDesc != "" {
		config.ServerDescription = serverDesc
	}
	
	return config, nil
}

// LoadConfigFromJSON loads configuration from a JSON file (legacy support)
func LoadConfigFromJSON(filename string) (*Config, error) {
	config := DefaultConfig()
	
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return config, nil
	}
	
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	
	// Override with environment variables if they exist
	if apiKey := os.Getenv("YOUTUBE_API_KEY"); apiKey != "" {
		config.YouTubeAPIKey = apiKey
	}
	
	return config, nil
}

// SaveConfig saves configuration to a JSON file
func (c *Config) SaveConfig(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}
