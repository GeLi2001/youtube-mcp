# YouTube MCP Server

A Go-based Model Context Protocol (MCP) server that provides access to the YouTube Data API v3. This server allows AI assistants and other MCP clients to search for videos, get channel information, retrieve video details, and more.

## Features

- **Video Search**: Search for YouTube videos with filtering options
- **Channel Information**: Get detailed information about YouTube channels
- **Video Details**: Retrieve comprehensive details about specific videos
- **Playlist Items**: Get items from YouTube playlists
- **Channel Search**: Search for YouTube channels

## Prerequisites

1. **Go 1.19 or later** installed on your system
2. **YouTube Data API v3 credentials** from Google Cloud Console

## Setup

### 1. Get YouTube API Credentials

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the YouTube Data API v3:
   - Navigate to "APIs & Services" → "Library"
   - Search for "YouTube Data API v3"
   - Click "Enable"
4. Create credentials:
   - Go to "APIs & Services" → "Credentials"
   - Click "Create Credentials" → "API Key" (for public data access)
   - Or create "OAuth client ID" (for user-specific data)

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configuration

#### Option A: Using .env File (Recommended)

1. Create a `.env` file from the example:

```bash
cp .env.example .env
```

2. Edit `.env` with your API key:

```bash
# YouTube Data API v3 Configuration
YOUTUBE_API_KEY=your_actual_api_key_here

# OAuth2 Configuration (optional, for user-specific data)
OAUTH2_CREDENTIALS_FILE=client_secret.json
TOKEN_FILE=token.json
```

#### Option B: Using Environment Variables

```bash
export YOUTUBE_API_KEY="your_api_key_here"
```

#### Option C: Using JSON Config File (Legacy)

```bash
cp config.example.json config.json
# Edit config.json with your settings
# Run with: ./youtube-mcp-server -json-config
```

## Usage

### Build and Run

```bash
# Build the server
go build -o youtube-mcp-server ./cmd/server

# Run the server (uses .env file)
./youtube-mcp-server

# Or run with JSON config file
./youtube-mcp-server -json-config -config config.json
```

### Running with Docker (Coming Soon)

```bash
# Build Docker image
docker build -t youtube-mcp-server .

# Run with environment variables
docker run -e YOUTUBE_API_KEY="your_api_key" youtube-mcp-server
```

## Available MCP Tools

### 1. search_videos

Search for YouTube videos based on a query.

**Parameters:**

- `query` (string, required): Search query for videos
- `max_results` (integer, optional): Maximum number of results (default: 10)
- `channel_id` (string, optional): Limit search to specific channel

**Example:**

```json
{
  "method": "tools/call",
  "params": {
    "name": "search_videos",
    "arguments": {
      "query": "golang tutorial",
      "max_results": 5
    }
  }
}
```

### 2. get_channel_info

Get information about a YouTube channel.

**Parameters:**

- `channel_id` (string, optional): Channel ID to get info for (if empty, uses authenticated user's channel)

**Example:**

```json
{
  "method": "tools/call",
  "params": {
    "name": "get_channel_info",
    "arguments": {
      "channel_id": "UCBUKHRdFNRo2gPXRXCKnN7w"
    }
  }
}
```

### 3. get_video_details

Get detailed information about a YouTube video.

**Parameters:**

- `video_id` (string, required): YouTube video ID

**Example:**

```json
{
  "method": "tools/call",
  "params": {
    "name": "get_video_details",
    "arguments": {
      "video_id": "dQw4w9WgXcQ"
    }
  }
}
```

### 4. get_playlist_items

Get items from a YouTube playlist.

**Parameters:**

- `playlist_id` (string, required): YouTube playlist ID
- `max_results` (integer, optional): Maximum number of results (default: 10)

**Example:**

```json
{
  "method": "tools/call",
  "params": {
    "name": "get_playlist_items",
    "arguments": {
      "playlist_id": "PLrAXtmRdnEQy4jrMVwm9wRbWqz4_0dmSh",
      "max_results": 10
    }
  }
}
```

### 5. search_channels

Search for YouTube channels based on a query.

**Parameters:**

- `query` (string, required): Search query for channels
- `max_results` (integer, optional): Maximum number of results (default: 10)

**Example:**

```json
{
  "method": "tools/call",
  "params": {
    "name": "search_channels",
    "arguments": {
      "query": "programming",
      "max_results": 5
    }
  }
}
```

## Configuration Options

The server can be configured via a JSON file or environment variables:

| Setting                   | Environment Variable | Description                     |
| ------------------------- | -------------------- | ------------------------------- |
| `youtube_api_key`         | `YOUTUBE_API_KEY`    | YouTube Data API v3 key         |
| `oauth2_credentials_file` | -                    | Path to OAuth2 credentials file |
| `token_file`              | -                    | Path to store OAuth2 tokens     |
| `server_name`             | -                    | MCP server name                 |
| `server_version`          | -                    | MCP server version              |
| `server_description`      | -                    | MCP server description          |

## Development

### Project Structure

```
youtube-mcp/
├── cmd/
│   └── server/
│       └── main.go              # Main server entry point
├── pkg/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── youtube/
│   │   └── youtube_client.go    # YouTube API client wrapper
│   └── mcp/
│       └── mcp_tools.go         # MCP tool definitions
├── .env.example                 # Example environment file
├── config.example.json          # Example configuration file (legacy)
├── go.mod                       # Go module definition
└── README.md                   # This file
```

### Adding New Tools

1. Add new argument structs in `pkg/mcp/mcp_tools.go`
2. Implement the YouTube API logic in `pkg/youtube/youtube_client.go`
3. Register the new tool in the `SetupMCPTools` function

### Testing

```bash
# Run tests
go test ./...

# Run with race detection
go test -race ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Troubleshooting

### Common Issues

1. **"quota exceeded" error**: You've exceeded your YouTube API quota. Wait for it to reset or upgrade your quota.

2. **"invalid API key" error**: Check that your API key is correct and the YouTube Data API v3 is enabled for your project.

3. **OAuth2 authentication issues**: Ensure your `client_secret.json` file is in the correct location and properly formatted.

### Getting Help

- Check the [YouTube Data API documentation](https://developers.google.com/youtube/v3)
- Review the [MCP specification](https://modelcontextprotocol.io/)
- Open an issue on GitHub for bugs or feature requests

## Resources

- [YouTube Data API v3 Documentation](https://developers.google.com/youtube/v3)
- [Google Cloud Console](https://console.cloud.google.com/)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [metoro-io/mcp-golang](https://github.com/metoro-io/mcp-golang)
