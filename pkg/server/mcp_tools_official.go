package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SearchVideosArgs represents arguments for video search
type SearchVideosArgs struct {
	Query      string `json:"query"`
	MaxResults int64  `json:"max_results,omitempty"`
	ChannelID  string `json:"channel_id,omitempty"`
}

// GetChannelInfoArgs represents arguments for getting channel information
type GetChannelInfoArgs struct {
	ChannelID string `json:"channel_id,omitempty"`
}

// GetVideoDetailsArgs represents arguments for getting video details
type GetVideoDetailsArgs struct {
	VideoID string `json:"video_id"`
}

// GetPlaylistItemsArgs represents arguments for getting playlist items
type GetPlaylistItemsArgs struct {
	PlaylistID string `json:"playlist_id"`
	MaxResults int64  `json:"max_results,omitempty"`
}

// SearchChannelsArgs represents arguments for channel search
type SearchChannelsArgs struct {
	Query      string `json:"query"`
	MaxResults int64  `json:"max_results,omitempty"`
}

// SetupOfficialMCPTools registers all MCP tools with the official SDK server
func SetupOfficialMCPTools(server *mcp.Server, youtubeClient *YouTubeClient) error {
	// Search videos tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_videos",
		Description: "Search for YouTube videos based on a query. Accepts query string, optional max_results (default 10), and optional channel_id to limit search to specific channel.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args SearchVideosArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 {
			args.MaxResults = 10
		}
		
		results, err := youtubeClient.SearchVideos(args.Query, args.MaxResults, args.ChannelID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to search videos: %v", err)
		}
		
		var videos []map[string]interface{}
		for _, item := range results {
			video := map[string]interface{}{
				"video_id":      item.Id.VideoId,
				"title":         item.Snippet.Title,
				"description":   item.Snippet.Description,
				"channel_id":    item.Snippet.ChannelId,
				"channel_title": item.Snippet.ChannelTitle,
				"published_at":  item.Snippet.PublishedAt,
				"thumbnail_url": item.Snippet.Thumbnails.Medium.Url,
			}
			videos = append(videos, video)
		}
		
		response, err := json.MarshalIndent(videos, "", "  ")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal response: %v", err)
		}
		
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(response)}},
		}, nil, nil
	})

	// Get channel info tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_channel_info",
		Description: "Get information about a YouTube channel",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GetChannelInfoArgs) (*mcp.CallToolResult, any, error) {
		channel, err := youtubeClient.GetChannelInfo(args.ChannelID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get channel info: %v", err)
		}
		
		channelInfo := map[string]interface{}{
			"channel_id":       channel.Id,
			"title":            channel.Snippet.Title,
			"description":      channel.Snippet.Description,
			"custom_url":       channel.Snippet.CustomUrl,
			"published_at":     channel.Snippet.PublishedAt,
			"country":          channel.Snippet.Country,
			"thumbnail_url":    channel.Snippet.Thumbnails.Medium.Url,
			"subscriber_count": channel.Statistics.SubscriberCount,
			"video_count":      channel.Statistics.VideoCount,
			"view_count":       channel.Statistics.ViewCount,
		}
		
		response, err := json.MarshalIndent(channelInfo, "", "  ")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal response: %v", err)
		}
		
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(response)}},
		}, nil, nil
	})

	// Get video details tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_video_details",
		Description: "Get detailed information about a YouTube video",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GetVideoDetailsArgs) (*mcp.CallToolResult, any, error) {
		video, err := youtubeClient.GetVideoDetails(args.VideoID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get video details: %v", err)
		}
		
		viewCount := video.Statistics.ViewCount
		likeCount := video.Statistics.LikeCount
		commentCount := video.Statistics.CommentCount
		
		videoInfo := map[string]interface{}{
			"video_id":       video.Id,
			"title":          video.Snippet.Title,
			"description":    video.Snippet.Description,
			"channel_id":     video.Snippet.ChannelId,
			"channel_title":  video.Snippet.ChannelTitle,
			"published_at":   video.Snippet.PublishedAt,
			"duration":       video.ContentDetails.Duration,
			"thumbnail_url":  video.Snippet.Thumbnails.Medium.Url,
			"view_count":     viewCount,
			"like_count":     likeCount,
			"comment_count":  commentCount,
			"tags":           video.Snippet.Tags,
			"category_id":    video.Snippet.CategoryId,
		}
		
		response, err := json.MarshalIndent(videoInfo, "", "  ")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal response: %v", err)
		}
		
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(response)}},
		}, nil, nil
	})

	// Get playlist items tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_playlist_items",
		Description: "Get items from a YouTube playlist",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GetPlaylistItemsArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 {
			args.MaxResults = 10
		}
		
		items, err := youtubeClient.GetPlaylistItems(args.PlaylistID, args.MaxResults)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get playlist items: %v", err)
		}
		
		var playlistItems []map[string]interface{}
		for _, item := range items {
			playlistItem := map[string]interface{}{
				"video_id":        item.ContentDetails.VideoId,
				"title":           item.Snippet.Title,
				"description":     item.Snippet.Description,
				"channel_id":      item.Snippet.ChannelId,
				"channel_title":   item.Snippet.ChannelTitle,
				"published_at":    item.Snippet.PublishedAt,
				"position":        item.Snippet.Position,
				"thumbnail_url":   item.Snippet.Thumbnails.Medium.Url,
			}
			playlistItems = append(playlistItems, playlistItem)
		}
		
		response, err := json.MarshalIndent(playlistItems, "", "  ")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal response: %v", err)
		}
		
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(response)}},
		}, nil, nil
	})

	// Search channels tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_channels",
		Description: "Search for YouTube channels based on a query",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args SearchChannelsArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 {
			args.MaxResults = 10
		}
		
		results, err := youtubeClient.SearchChannels(args.Query, args.MaxResults)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to search channels: %v", err)
		}
		
		var channels []map[string]interface{}
		for _, item := range results {
			channel := map[string]interface{}{
				"channel_id":     item.Id.ChannelId,
				"title":          item.Snippet.Title,
				"description":    item.Snippet.Description,
				"published_at":   item.Snippet.PublishedAt,
				"thumbnail_url":  item.Snippet.Thumbnails.Medium.Url,
			}
			channels = append(channels, channel)
		}
		
		response, err := json.MarshalIndent(channels, "", "  ")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal response: %v", err)
		}
		
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(response)}},
		}, nil, nil
	})

	return nil
}
