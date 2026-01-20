# Beeper API CLI

A cross-platform command-line interface for the Beeper Desktop API. Built for programmatic access to Beeper conversations with LLM-friendly output formats.

🚀 **[Quick Start Guide](QUICKSTART.md)** | 📚 **[API Documentation](API.md)** | 🔧 **[Examples](examples/)**

## Purpose

While existing tools read the Beeper SQLite database directly, this CLI interfaces with the Beeper Desktop HTTP API to provide both read and write capabilities. This enables sending messages, managing conversations, and full bidirectional communication through Beeper's unified chat platform.

## Features

- **Full API Coverage**: Read and write operations across all Beeper-supported chat networks
- **LLM-Friendly Output**: JSON, plain text, and markdown formats optimized for AI agent consumption
- **Auto-Discovery**: Automatically detect Beeper Desktop API endpoint when available
- **Cross-Platform**: Single binary for Linux, macOS, and Windows
- **Pipeable**: Unix-friendly design for scripting and tool composition

## Installation

```bash
# Download latest release for your platform
# macOS (arm64)
curl -L https://github.com/nerveband/beeper-api-cli/releases/latest/download/beeper-api-darwin-arm64 -o beeper
chmod +x beeper
sudo mv beeper /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/nerveband/beeper-api-cli/releases/latest/download/beeper-api-linux-amd64 -o beeper
chmod +x beeper
sudo mv beeper /usr/local/bin/

# Or build from source
go install github.com/nerveband/beeper-api-cli@latest
```

## Quick Start

```bash
# Auto-discover Beeper Desktop API
beeper discover

# Or manually configure
beeper config set-url http://localhost:39867

# List all chats
beeper chats list

# Get messages from a chat
beeper messages list --chat-id CHAT_ID

# Send a message
beeper send --chat-id CHAT_ID --message "Hello from CLI"

# Search messages
beeper search --query "important meeting"
```

## Configuration

The CLI stores configuration in `~/.beeper-api-cli/config.yaml`:

```yaml
api_url: http://localhost:39867
output_format: json  # json, text, markdown
```

## API Coverage

### Read Operations
- `chats list` - List all conversations
- `chats get` - Get chat details
- `messages list` - Retrieve messages from a chat
- `messages get` - Get specific message details
- `search` - Search across all messages
- `users get` - Get user information

### Write Operations
- `send` - Send new message
- `react` - Add reaction to message
- `edit` - Edit existing message
- `delete` - Delete message
- `read` - Mark messages as read

## Output Formats

All commands support multiple output formats via `--output` flag:

```bash
# JSON (default, ideal for LLM parsing)
beeper chats list --output json

# Plain text (human-readable)
beeper chats list --output text

# Markdown (formatted for documentation)
beeper chats list --output markdown
```

## Examples

### List chats with participants
```bash
beeper chats list | jq '.[] | {name, participants, last_message}'
```

### Send message and capture response
```bash
MESSAGE_ID=$(beeper send --chat-id CHAT --message "Status update" --output json | jq -r '.message_id')
```

### Search and export
```bash
beeper search --query "invoice" --output markdown > invoices.md
```

## Architecture

Built with Go for:
- Fast execution
- Easy cross-compilation
- Minimal dependencies
- Single binary distribution

Uses Cobra for CLI framework and Viper for configuration management.

## Philosophy & Design Goals

**beeper-api-cli** is designed for simplicity and minimal dependencies:

- **Pure API approach** - No SQLite dependencies, no keyring requirements
- **Simple authentication** - Manual token configuration (no OAuth flow)
- **Focused feature set** - Core read/write operations for programmatic access
- **LLM-friendly** - JSON/text/markdown output optimized for AI agents
- **Single binary** - Easy cross-compilation and distribution

This makes it ideal for:
- CI/CD pipelines and automation scripts
- LLM/AI agent integrations
- Environments where you want minimal dependencies
- Quick setup without browser auth flows

## Choosing the Right Tool

Different Beeper CLI tools have different strengths. **Use the best tool for your use case!**

### [salmonumbrella/beeper-cli](https://github.com/salmonumbrella/beeper-cli) 🌟
**Full-featured terminal client with desktop integration**
- OAuth browser authentication with keyring storage
- Real-time message streaming (tail/follow mode)
- Desktop control (focus window, navigate chats, pre-fill drafts)
- Advanced features: reminders, bulk archive, inbox view
- DB backend for faster local search
- Template output and JQ integration

**Best for:** Interactive terminal use, power users, desktop integration

### [KrauseFx/beeper-cli](https://github.com/KrauseFx/beeper-cli) 📖
**Read-only SQLite access** (original inspiration for this project)
- Direct database queries
- No API dependencies

**Best for:** Quick local reads, offline access

### beeper-api-cli (this tool) 🚀
**Lightweight API client for automation**
- Pure HTTP API (no SQLite, no keyring)
- Write operations (send, edit, react, delete)
- Simple manual token auth
- LLM-optimized JSON output
- Minimal dependencies

**Best for:** Automation, CI/CD, LLM integration, simple deployments

## Development

```bash
# Clone repository
git clone https://github.com/nerveband/beeper-api-cli
cd beeper-api-cli

# Install dependencies
go mod download

# Build
go build -o beeper

# Run tests
go test ./...

# Cross-compile
./build.sh
```

## Requirements

- Beeper Desktop application installed and running
- API server enabled (default: http://localhost:39867)

## License

MIT

## Contributing

Contributions welcome. Please open an issue before starting major work.
