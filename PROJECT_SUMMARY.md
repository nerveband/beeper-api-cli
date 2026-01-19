# Beeper CLI - Project Summary

## Overview

A cross-platform command-line interface for the Beeper Desktop API, built in Go. Provides LLM-friendly interfaces for reading and sending messages across all Beeper-supported chat networks.

## Repository

**GitHub**: https://github.com/nerveband/beeper-cli

**Latest Release**: v0.1.0 - https://github.com/nerveband/beeper-cli/releases/tag/v0.1.0

## Architecture

### Technology Stack
- **Language**: Go 1.22+
- **CLI Framework**: Cobra
- **Configuration**: Viper
- **Build**: Native Go cross-compilation

### Project Structure
```
beeper-cli/
├── main.go                    # Entry point
├── cmd/                       # CLI commands
│   ├── root.go               # Root command
│   ├── discover.go           # Auto-discovery
│   ├── config.go             # Configuration management
│   ├── chats.go              # Chat operations
│   ├── messages.go           # Message retrieval
│   ├── send.go               # Send messages
│   └── search.go             # Search messages
├── internal/
│   ├── api/                  # API client
│   │   └── client.go        # HTTP client + models
│   ├── config/               # Configuration
│   │   └── config.go        # Config management
│   └── output/               # Output formatting
│       └── formatter.go     # JSON/text/markdown formatters
├── examples/                 # Usage examples
│   ├── basic-usage.sh
│   └── llm-integration.sh
├── .github/workflows/        # CI/CD
│   ├── build.yml            # Build on push
│   └── release.yml          # Release on tag
├── build.sh                  # Cross-compilation script
├── README.md                 # Main documentation
├── API.md                    # API documentation
├── CONTRIBUTING.md           # Contribution guide
└── LICENSE                   # MIT License
```

## Features

### Core Functionality
✅ Auto-discover Beeper Desktop API endpoint
✅ Configuration management (URL, output format)
✅ List all chats with metadata
✅ Get specific chat details
✅ List messages from chats (with limit)
✅ Send messages to chats
✅ Search messages across all chats

### Output Formats
- **JSON**: Structured data for LLM parsing
- **Plain Text**: Human-readable output
- **Markdown**: Formatted for documentation

### Cross-Platform Support
- macOS (amd64, arm64)
- Linux (amd64, arm64)
- Windows (amd64)

## API Coverage

The CLI wraps the Beeper Desktop HTTP API:

| Endpoint | CLI Command | Status |
|----------|-------------|--------|
| `GET /health` | `discover` | ✅ Implemented |
| `GET /chats` | `chats list` | ✅ Implemented |
| `GET /chats/{id}` | `chats get` | ✅ Implemented |
| `GET /chats/{id}/messages` | `messages list` | ✅ Implemented |
| `POST /messages/send` | `send` | ✅ Implemented |
| `GET /search` | `search` | ✅ Implemented |

## Usage Examples

### Basic Operations
```bash
# Auto-discover API
beeper discover

# List chats
beeper chats list

# Send message
beeper send --chat-id CHAT_ID --message "Hello!"

# Search messages
beeper search --query "important"
```

### LLM Integration
```bash
# Extract chat IDs for processing
beeper chats list | jq -r '.[] | .id'

# Find chats with unread messages
beeper chats list | jq -r '.[] | select(.unread_count > 0)'

# Send and capture message ID
MESSAGE_ID=$(beeper send --chat-id CHAT --message "Hi" | jq -r '.message_id')

# Export to markdown
beeper messages list --chat-id CHAT --output markdown > messages.md
```

## Development

### Build from Source
```bash
git clone https://github.com/nerveband/beeper-cli
cd beeper-cli
go build -o beeper .
```

### Cross-Compile
```bash
./build.sh
# Outputs to dist/ directory
```

### Run Tests
```bash
go test ./...
```

## CI/CD

### GitHub Actions Workflows

**Build Workflow** (`.github/workflows/build.yml`)
- Triggers on push/PR to main
- Runs on Ubuntu
- Builds for all platforms
- Runs tests

**Release Workflow** (`.github/workflows/release.yml`)
- Triggers on version tags (v*)
- Builds all platform binaries
- Creates GitHub release
- Uploads binaries as assets

### Creating a Release
```bash
git tag v0.2.0
git push origin v0.2.0
# GitHub Actions automatically builds and releases
```

## Installation

### Pre-built Binaries
Download from: https://github.com/nerveband/beeper-cli/releases/latest

```bash
# macOS (arm64)
curl -L https://github.com/nerveband/beeper-cli/releases/latest/download/beeper-darwin-arm64 -o beeper
chmod +x beeper
sudo mv beeper /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/nerveband/beeper-cli/releases/latest/download/beeper-linux-amd64 -o beeper
chmod +x beeper
sudo mv beeper /usr/local/bin/
```

### From Source
```bash
go install github.com/nerveband/beeper-cli@latest
```

## Configuration

Default config location: `~/.beeper-cli/config.yaml`

```yaml
api_url: http://localhost:39867
output_format: json  # json, text, markdown
```

## Requirements

- Beeper Desktop application running
- API server enabled (default port: 39867)

## Comparison with Existing Tools

### vs. beeper-tools
The existing [beeper-tools](https://github.com/beeper/beeper-tools):
- ✅ Read-only (direct SQLite access)
- ❌ Cannot send messages
- ❌ No real-time operations

**This CLI**:
- ✅ Read and write operations
- ✅ Uses HTTP API (no database locks)
- ✅ LLM-optimized output formats
- ✅ Cross-platform binaries

## Future Enhancements

Potential additions:
- [ ] React to messages
- [ ] Edit messages
- [ ] Delete messages
- [ ] Mark messages as read
- [ ] Get user information
- [ ] Media upload/download
- [ ] Typing indicators
- [ ] Real-time message streaming (WebSocket)

## License

MIT License - See [LICENSE](LICENSE)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## Links

- **Repository**: https://github.com/nerveband/beeper-cli
- **Issues**: https://github.com/nerveband/beeper-cli/issues
- **Releases**: https://github.com/nerveband/beeper-cli/releases
- **Beeper Desktop**: https://www.beeper.com/

## Credits

Created by Ashraf Ali (@nerveband)

Built for LLM agents and power users who need programmatic access to Beeper.
