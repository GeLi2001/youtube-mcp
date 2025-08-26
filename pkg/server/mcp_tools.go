package server

import (
	"encoding/json"
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
)

// SearchVideosArgs represents arguments for video search
type SearchVideosArgs struct {
	Query      string `json:"query" jsonschema:"description=Search query for videos"`
	MaxResults int64  `json:"max_results,omitempty" jsonschema:"description=Maximum number of results to return (default: 10)"`
	ChannelID  string `json:"channel_id,omitempty" jsonschema:"description=Optional channel ID to search within"`
}

// GetChannelInfoArgs represents arguments for getting channel information
type GetChannelInfoArgs struct {
	ChannelID string `json:"channel_id,omitempty" jsonschema:"description=Channel ID to get info for (if empty, uses authenticated user's channel)"`
}

// GetVideoDetailsArgs represents arguments for getting video details
type GetVideoDetailsArgs struct {
	VideoID string `json:"video_id" jsonschema:"description=YouTube video ID"`
}

// GetPlaylistItemsArgs represents arguments for getting playlist items
type GetPlaylistItemsArgs struct {
	PlaylistID string `json:"playlist_id" jsonschema:"description=YouTube playlist ID"`
	MaxResults int64  `json:"max_results,omitempty" jsonschema:"description=Maximum number of results to return (default: 10)"`
}

// SearchChannelsArgs represents arguments for channel search
type SearchChannelsArgs struct {
	Query      string `json:"query" jsonschema:"description=Search query for channels"`
	MaxResults int64  `json:"max_results,omitempty" jsonschema:"description=Maximum number of results to return (default: 10)"`
}

// SetupMCPTools registers all MCP tools with the server
func SetupMCPTools(server *mcp_golang.Server, youtubeClient *YouTubeClient) error {
	// Search videos tool
	err := server.RegisterTool(
		"search_videos",
		"Search for YouTube videos based on a query",
		func(args SearchVideosArgs) (*mcp_golang.ToolResponse, error) {
			if args.MaxResults == 0 {
				args.MaxResults = 10
			}
			
			results, err := youtubeClient.SearchVideos(args.Query, args.MaxResults, args.ChannelID)
			if err != nil {
				return nil, fmt.Errorf("failed to search videos: %v", err)
			}
			
			var videos []map[string]interface{}
			for _, item := range results {
				video := map[string]interface{}{
					"video_id":     item.Id.VideoId,
					"title":        item.Snippet.Title,
					"description":  item.Snippet.Description,
					"channel_id":   item.Snippet.ChannelId,
					"channel_title": item.Snippet.ChannelTitle,
					"published_at": item.Snippet.PublishedAt,
					"thumbnail_url": item.Snippet.Thumbnails.Medium.Url,
				}
				videos = append(videos, video)
			}
			
			response, err := json.MarshalIndent(videos, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %v", err)
			}
			
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(response))), nil
		})
	if err != nil {
		return fmt.Errorf("failed to register search_videos tool: %v", err)
	}
	
	// Get channel info tool
	err = server.RegisterTool(
		"get_channel_info",
		"Get information about a YouTube channel",
		func(args GetChannelInfoArgs) (*mcp_golang.ToolResponse, error) {
			channel, err := youtubeClient.GetChannelInfo(args.ChannelID)
			if err != nil {
				return nil, fmt.Errorf("failed to get channel info: %v", err)
			}
			
			channelInfo := map[string]interface{}{
				"channel_id":      channel.Id,
				"title":           channel.Snippet.Title,
				"description":     channel.Snippet.Description,
				"custom_url":      channel.Snippet.CustomUrl,
				"published_at":    channel.Snippet.PublishedAt,
				"country":         channel.Snippet.Country,
				"thumbnail_url":   channel.Snippet.Thumbnails.Medium.Url,
				"subscriber_count": channel.Statistics.SubscriberCount,
				"video_count":     channel.Statistics.VideoCount,
				"view_count":      channel.Statistics.ViewCount,
			}
			
			response, err := json.MarshalIndent(channelInfo, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %v", err)
			}
			
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(response))), nil
		})
	if err != nil {
		return fmt.Errorf("failed to register get_channel_info tool: %v", err)
	}
	
	// Get video details tool
	err = server.RegisterTool(
		"get_video_details",
		"Get detailed information about a YouTube video",
		func(args GetVideoDetailsArgs) (*mcp_golang.ToolResponse, error) {
			video, err := youtubeClient.GetVideoDetails(args.VideoID)
			if err != nil {
				return nil, fmt.Errorf("failed to get video details: %v", err)
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
				return nil, fmt.Errorf("failed to marshal response: %v", err)
			}
			
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(response))), nil
		})
	if err != nil {
		return fmt.Errorf("failed to register get_video_details tool: %v", err)
	}
	
	// Get playlist items tool
	err = server.RegisterTool(
		"get_playlist_items",
		"Get items from a YouTube playlist",
		func(args GetPlaylistItemsArgs) (*mcp_golang.ToolResponse, error) {
			if args.MaxResults == 0 {
				args.MaxResults = 10
			}
			
			items, err := youtubeClient.GetPlaylistItems(args.PlaylistID, args.MaxResults)
			if err != nil {
				return nil, fmt.Errorf("failed to get playlist items: %v", err)
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
				return nil, fmt.Errorf("failed to marshal response: %v", err)
			}
			
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(response))), nil
		})
	if err != nil {
		return fmt.Errorf("failed to register get_playlist_items tool: %v", err)
	}
	
	// Search channels tool
	err = server.RegisterTool(
		"search_channels",
		"Search for YouTube channels based on a query",
		func(args SearchChannelsArgs) (*mcp_golang.ToolResponse, error) {
			if args.MaxResults == 0 {
				args.MaxResults = 10
			}
			
			results, err := youtubeClient.SearchChannels(args.Query, args.MaxResults)
			if err != nil {
				return nil, fmt.Errorf("failed to search channels: %v", err)
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
				return nil, fmt.Errorf("failed to marshal response: %v", err)
			}
			
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(response))), nil
		})
	if err != nil {
		return fmt.Errorf("failed to register search_channels tool: %v", err)
	}
	
	return nil
}
