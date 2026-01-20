# Beeper CLI Quick Start

Get up and running with Beeper CLI in under 5 minutes.

## Prerequisites

- Beeper Desktop installed and running
- Terminal access

## Step 1: Install

### Option A: Download Pre-built Binary (Recommended)

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/nerveband/beeper-api-cli/releases/latest/download/beeper-api-darwin-arm64 -o beeper
chmod +x beeper
sudo mv beeper /usr/local/bin/
```

**macOS (Intel)**
```bash
curl -L https://github.com/nerveband/beeper-api-cli/releases/latest/download/beeper-api-darwin-amd64 -o beeper
chmod +x beeper
sudo mv beeper /usr/local/bin/
```

**Linux (amd64)**
```bash
curl -L https://github.com/nerveband/beeper-api-cli/releases/latest/download/beeper-api-linux-amd64 -o beeper
chmod +x beeper
sudo mv beeper /usr/local/bin/
```

**Windows**
Download from: https://github.com/nerveband/beeper-api-cli/releases/latest/download/beeper-api-windows-amd64.exe

### Option B: Build from Source
```bash
git clone https://github.com/nerveband/beeper-api-cli
cd beeper-api-cli
go build -o beeper .
sudo mv beeper /usr/local/bin/
```

## Step 2: Configure

Auto-discover your Beeper Desktop API:
```bash
beeper discover
```

Expected output:
```
Discovering Beeper Desktop API...
Found Beeper Desktop API at: http://localhost:39867
Configuration saved successfully!
```

**Manual Configuration** (if auto-discovery fails):
```bash
beeper config set-url http://localhost:39867
```

## Step 3: Verify Installation

Check your configuration:
```bash
beeper config show
```

Expected output:
```
API URL: http://localhost:39867
Output Format: json
```

## Step 4: List Your Chats

```bash
beeper chats list
```

You should see JSON output with all your Beeper conversations.

### Human-Readable Format
```bash
beeper chats list --output text
```

## Step 5: Send Your First Message

1. **Get a chat ID** from the previous command
2. **Send a message**:

```bash
beeper send --chat-id YOUR_CHAT_ID --message "Hello from Beeper CLI!"
```

## Common Commands

### List messages from a chat
```bash
beeper messages list --chat-id CHAT_ID --limit 20
```

### Search for messages
```bash
beeper search --query "important meeting"
```

### Change output format
```bash
beeper chats list --output markdown
beeper chats list --output text
```

## Next Steps

- **LLM Integration**: See `examples/llm-integration.sh` for advanced usage
- **Full Documentation**: Read `README.md` for complete command reference
- **API Details**: Check `API.md` for endpoint documentation

## Troubleshooting

### "Could not auto-discover Beeper Desktop API"
- Ensure Beeper Desktop is running
- Check if port 39867 is accessible: `curl http://localhost:39867/health`
- Try manual configuration: `beeper config set-url http://localhost:PORT`

### "API error (status 404)"
- Verify Beeper Desktop API server is enabled in settings
- Check the API URL is correct: `beeper config show`

### Commands not working
- Verify installation: `beeper --help`
- Check you're using the latest version
- Open an issue: https://github.com/nerveband/beeper-api-cli/issues

## Getting Help

- **Command Help**: `beeper [command] --help`
- **GitHub Issues**: https://github.com/nerveband/beeper-api-cli/issues
- **Documentation**: https://github.com/nerveband/beeper-api-cli

## Example Workflow

```bash
# 1. Setup
beeper discover

# 2. Explore your chats
beeper chats list --output text

# 3. Get chat ID you want to message
CHAT_ID="your_chat_id_here"

# 4. Read recent messages
beeper messages list --chat-id $CHAT_ID --limit 10

# 5. Send a message
beeper send --chat-id $CHAT_ID --message "Testing Beeper CLI!"

# 6. Search across all messages
beeper search --query "project deadline" --output markdown
```

Now you're ready to automate your Beeper workflows!
