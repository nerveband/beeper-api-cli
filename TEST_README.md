# Beeper CLI Test Suite

Comprehensive TDD test coverage for the Beeper CLI tool.

## Test Coverage

### ✅ Unit Tests (Offline - No API Required)

#### 1. Output Formatter Tests (`internal/output`)
- ✅ JSON formatting for chats and messages
- ✅ Plain text formatting
- ✅ Markdown formatting  
- ✅ Empty list handling
- ✅ Invalid format fallback
- ✅ Chat name edge cases
- ✅ Timestamp formatting
- ✅ Long messages
- ✅ Special characters escaping

**Run:**
```bash
go test ./internal/output -v
```

#### 2. Configuration Tests (`internal/config`)
- ✅ Load/save configuration from file
- ✅ Default configuration values
- ✅ Configuration validation
- ✅ Environment variable override
- ✅ Configuration merging
- ✅ Partial configuration updates
- ✅ File permissions
- ✅ Invalid YAML handling

**Run:**
```bash
go test ./internal/config -v
```

### 🔌 Integration Tests (Require Live Beeper Desktop API)

#### 3. API Client Tests (`internal/api`)
Tests real HTTP communication with Beeper Desktop API.

**Prerequisites:**
- Beeper Desktop running
- `BEEPER_API_URL` set (default: http://[::1]:23373)
- `BEEPER_TOKEN` set (Bearer token from Beeper settings)

**Tests:**
- ListChats - fetch all conversations
- GetChat - fetch specific chat details
- ListMessages - fetch messages from a chat
- SendMessage - send a message (requires `BEEPER_TEST_CHAT_ID`)
- SearchMessages - search across all messages
- Ping - health check
- Error handling for invalid URLs

**Run:**
```bash
export BEEPER_API_URL="http://[::1]:23373"
export BEEPER_TOKEN="your-token-here"
export BEEPER_TEST_CHAT_ID="safe-test-chat-id"  # Optional, for send tests

go test ./internal/api -v
```

**Skip behavior:**
- Tests automatically skip if `BEEPER_API_URL` or `BEEPER_TOKEN` not set
- Send tests skip if `BEEPER_TEST_CHAT_ID` not set

#### 4. Command Tests (`cmd/`)
Tests CLI command execution with real API.

**Prerequisites:** Same as API tests above

**Tests:**
- `chats list` - JSON, text, markdown formats
- `chats get` - specific chat retrieval
- `messages list` - with limit parameter
- `send` - message sending + error cases
- `search` - keyword search with limits
- `discover` - API discovery
- `config` - show/set/get/validate

**Run:**
```bash
go test ./cmd -v
```

#### 5. End-to-End Integration Tests
Full workflow tests using the compiled binary.

**Prerequisites:** 
- Built binary at `./beeper`
- Live Beeper Desktop API
- Environment variables set

**Tests:**
- Complete workflow (discover → list → get → messages → send → search)
- All output formats (JSON/text/markdown)
- Error handling scenarios
- Configuration management workflow
- Unix pipeline compatibility (jq, grep)

**Run:**
```bash
# Build first
./build.sh

# Run integration tests
go test -tags=integration -v
```

## Test Environment Setup

### 1. Get Beeper Desktop API Token

1. Open Beeper Desktop
2. Go to Settings → Advanced → API
3. Enable API access
4. Copy the Bearer token

### 2. Set Environment Variables

```bash
export BEEPER_API_URL="http://[::1]:23373"
export BEEPER_TOKEN="your-bearer-token-here"

# Optional: for safe send testing
export BEEPER_TEST_CHAT_ID="your-test-chat-id"
```

### 3. Find a Test Chat ID

```bash
./beeper chats list --output json | jq -r '.[0].id'
```

## Running All Tests

### Unit Tests Only (No API Required)
```bash
go test ./internal/output ./internal/config -v
```

### Integration Tests (Requires Live API)
```bash
# Ensure Beeper Desktop is running and env vars are set
go test ./internal/api ./cmd -v
```

### Full Test Suite
```bash
# Unit + Integration
go test ./... -v

# Including E2E integration tests
./build.sh && go test -tags=integration ./... -v
```

## Test Coverage Report

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Structure

```
beeper-api-cli/
├── internal/
│   ├── api/
│   │   ├── client.go
│   │   └── client_test.go          # Real API tests
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go          # Config management tests
│   └── output/
│       ├── formatter.go
│       └── formatter_test.go       # Formatter unit tests
├── cmd/
│   ├── chats.go
│   ├── chats_test.go               # Command tests
│   ├── messages_test.go
│   ├── send_test.go
│   ├── search_test.go
│   ├── discover_test.go
│   └── config_test.go
└── integration_test.go             # E2E workflow tests
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: go test ./internal/output ./internal/config -v

  integration-tests:
    runs-on: ubuntu-latest
    services:
      beeper:
        # Mock Beeper API service (if available)
        # Or skip integration tests in CI
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run integration tests
        run: |
          export BEEPER_API_URL="${{ secrets.BEEPER_API_URL }}"
          export BEEPER_TOKEN="${{ secrets.BEEPER_TOKEN }}"
          go test ./internal/api ./cmd -v
        continue-on-error: true  # Allow failure if API not available
```

## Test Principles

1. **TDD Approach**: Tests written first, following Test-Driven Development
2. **Real API**: Integration tests use actual Beeper Desktop API (no mocks)
3. **Graceful Skipping**: Tests skip automatically if environment not configured
4. **Isolated Unit Tests**: Formatter and config tests don't require external services
5. **Table-Driven**: Common patterns use table-driven tests
6. **Clean Assertions**: Using testify/assert for readable test code

## Troubleshooting

### "connection refused" errors
- Ensure Beeper Desktop is running
- Verify API is enabled in Beeper settings
- Check `BEEPER_API_URL` matches your Beeper API port

### Tests skip automatically
- Set `BEEPER_API_URL` and `BEEPER_TOKEN` environment variables
- Tests requiring API will skip if these are not set

### Send tests fail
- Set `BEEPER_TEST_CHAT_ID` to a safe test chat
- Ensure you have permission to send to that chat

## Next Steps

- [ ] Increase coverage to >80% (current: comprehensive for core modules)
- [ ] Add benchmark tests for performance
- [ ] Create mock Beeper API server for CI/CD
- [ ] Add mutation testing
- [ ] Property-based testing for complex scenarios
