package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// YouTubeClient wraps the YouTube Data API client
type YouTubeClient struct {
	service *youtube.Service
	config  *Config
}

// NewYouTubeClient creates a new YouTube client
func NewYouTubeClient(cfg *Config) (*YouTubeClient, error) {
	ctx := context.Background()
	
	var service *youtube.Service
	var err error
	
	// Try to create service with API key first (for public data)
	if cfg.YouTubeAPIKey != "" {
		service, err = youtube.NewService(ctx, option.WithAPIKey(cfg.YouTubeAPIKey))
		if err != nil {
			log.Printf("Failed to create service with API key: %v", err)
		}
	}
	
	// If API key fails or doesn't exist, try OAuth2
	if service == nil {
		client, err := getOAuth2Client(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to get OAuth2 client: %v", err)
		}
		
		service, err = youtube.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return nil, fmt.Errorf("failed to create YouTube service: %v", err)
		}
	}
	
	return &YouTubeClient{
		service: service,
		config:  cfg,
	}, nil
}

// getOAuth2Client gets an OAuth2 client for authenticated requests
func getOAuth2Client(cfg *Config) (*http.Client, error) {
	// Read OAuth2 credentials
	b, err := os.ReadFile(cfg.OAuth2CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}
	
	oauthConfig, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file: %v", err)
	}
	
	return getClient(oauthConfig, cfg.TokenFile), nil
}

// getClient gets an HTTP client with OAuth2 token
func getClient(config *oauth2.Config, tokenFile string) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// tokenFromFile retrieves a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// getTokenFromWeb uses config to request a Token
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// saveToken saves a token to a file path
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// SearchVideos searches for videos based on query
func (yc *YouTubeClient) SearchVideos(query string, maxResults int64, channelID string) ([]*youtube.SearchResult, error) {
	call := yc.service.Search.List([]string{"snippet"}).
		Q(query).
		Type("video").
		MaxResults(maxResults).
		Order("relevance")
	
	if channelID != "" {
		call = call.ChannelId(channelID)
	}
	
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error searching videos: %v", err)
	}
	
	return response.Items, nil
}

// GetChannelInfo gets information about a channel
func (yc *YouTubeClient) GetChannelInfo(channelID string) (*youtube.Channel, error) {
	call := yc.service.Channels.List([]string{"snippet", "statistics", "contentDetails"})
	
	if channelID != "" {
		call = call.Id(channelID)
	} else {
		call = call.Mine(true)
	}
	
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error getting channel info: %v", err)
	}
	
	if len(response.Items) == 0 {
		return nil, fmt.Errorf("channel not found")
	}
	
	return response.Items[0], nil
}

// GetVideoDetails gets detailed information about a video
func (yc *YouTubeClient) GetVideoDetails(videoID string) (*youtube.Video, error) {
	call := yc.service.Videos.List([]string{"snippet", "statistics", "contentDetails"}).
		Id(videoID)
	
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error getting video details: %v", err)
	}
	
	if len(response.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	}
	
	return response.Items[0], nil
}

// GetPlaylistItems gets items from a playlist
func (yc *YouTubeClient) GetPlaylistItems(playlistID string, maxResults int64) ([]*youtube.PlaylistItem, error) {
	call := yc.service.PlaylistItems.List([]string{"snippet", "contentDetails"}).
		PlaylistId(playlistID).
		MaxResults(maxResults)
	
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error getting playlist items: %v", err)
	}
	
	return response.Items, nil
}

// SearchChannels searches for channels based on query
func (yc *YouTubeClient) SearchChannels(query string, maxResults int64) ([]*youtube.SearchResult, error) {
	call := yc.service.Search.List([]string{"snippet"}).
		Q(query).
		Type("channel").
		MaxResults(maxResults).
		Order("relevance")
	
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error searching channels: %v", err)
	}
	
	return response.Items, nil
}
