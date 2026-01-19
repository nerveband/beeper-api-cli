# Implementation Notes

Internal documentation for developers working on Beeper CLI.

## Architecture Decisions

### Why Go?
- Fast compilation and execution
- Excellent cross-platform support (single binary)
- Strong standard library (HTTP, JSON)
- Easy to distribute (no runtime dependencies)
- Good CLI frameworks (Cobra, Viper)

### Why Cobra + Viper?
- Industry-standard CLI framework
- Automatic help generation
- Subcommand structure
- Flag parsing
- Configuration management (Viper)
- POSIX-compliant

### API Client Design
- Thin wrapper around HTTP endpoints
- Strongly typed structs for requests/responses
- Error handling at client level
- No global state

### Output Formatting
- Separate formatter package
- Format-agnostic commands
- Easy to add new formats
- Consistent interface

## Project Structure Rationale

```
beeper-cli/
├── main.go              # Minimal entry point
├── cmd/                 # Command implementations
├── internal/            # Private packages
│   ├── api/            # API client (reusable)
│   ├── config/         # Configuration management
│   └── output/         # Output formatting
```

**Why `internal/`?**
- Prevents external packages from importing
- Keeps implementation details private
- Clean public API surface

**Why separate `api` package?**
- Reusable in other projects
- Testable independently
- Clear separation of concerns

## API Client Implementation

### Current Approach
The API client in `internal/api/client.go` makes assumptions about Beeper Desktop API structure based on common REST patterns:

```go
// Example endpoint structure
GET  /chats           -> List all chats
GET  /chats/{id}      -> Get specific chat
GET  /chats/{id}/messages -> List messages
POST /messages/send   -> Send message
GET  /search          -> Search messages
```

### Real API Integration
When connecting to actual Beeper Desktop API:

1. **Discover Endpoints**: Use `/api/spec` or similar to get OpenAPI spec
2. **Update Structs**: Adjust `Chat` and `Message` structs based on actual responses
3. **Error Handling**: Parse actual error response format
4. **Authentication**: Add auth headers if required
5. **Rate Limiting**: Implement backoff/retry if needed

### Testing Strategy
```bash
# Manual testing against real API
beeper --api-url http://localhost:39867 chats list

# Mock server for automated tests (future)
go test -v ./internal/api
```

## Configuration Management

### Current Config
```yaml
api_url: http://localhost:39867
output_format: json
```

### Future Enhancements
```yaml
api_url: http://localhost:39867
output_format: json
auth:
  token: optional_api_token
  username: optional_username
retry:
  max_attempts: 3
  backoff_ms: 1000
cache:
  enabled: true
  ttl_seconds: 300
```

## Output Formatting

### Current Formats
1. **JSON**: Machine-readable, ideal for LLMs
2. **Text**: Human-readable, simple
3. **Markdown**: Documentation-friendly

### Adding New Format
1. Add to `output.Format` enum
2. Implement format functions in `formatter.go`
3. Update command output handling

Example:
```go
// Add to formatter.go
func formatChatsCSV(chats []api.Chat) string {
    var sb strings.Builder
    sb.WriteString("ID,Name,Participants,Unread\n")
    for _, chat := range chats {
        sb.WriteString(fmt.Sprintf("%s,%s,%s,%d\n",
            chat.ID,
            chat.Name,
            strings.Join(chat.Participants, ";"),
            chat.UnreadCount))
    }
    return sb.String()
}
```

## Error Handling Strategy

### Current Approach
- Client-level: HTTP errors, JSON parsing
- Command-level: Missing parameters, validation
- User-facing: Clear error messages with hints

### Best Practices
```go
// Good
return fmt.Errorf("failed to list chats: %w", err)

// Bad
return fmt.Errorf("error: %v", err)
```

## Testing Strategy

### Unit Tests (Future)
```go
// internal/api/client_test.go
func TestListChats(t *testing.T) {
    // Mock HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`[{"id":"chat1","name":"Test"}]`))
    }))
    defer server.Close()

    client := api.NewClient(server.URL)
    chats, err := client.ListChats()
    
    assert.NoError(t, err)
    assert.Len(t, chats, 1)
    assert.Equal(t, "chat1", chats[0].ID)
}
```

### Integration Tests
```bash
# Test against real Beeper Desktop
./test-integration.sh
```

## Build Process

### Cross-Compilation
`build.sh` compiles for multiple platforms:
- darwin/amd64 (macOS Intel)
- darwin/arm64 (macOS Apple Silicon)
- linux/amd64
- linux/arm64
- windows/amd64

### Release Process
1. Update version in code (if versioning added)
2. Tag release: `git tag v0.2.0`
3. Push tag: `git push origin v0.2.0`
4. GitHub Actions builds and publishes

### Binary Size Optimization
Current binaries are ~10-11MB. To reduce:

```bash
# Strip debug symbols
go build -ldflags="-s -w" -o beeper .

# UPX compression (optional)
upx --best beeper
```

## Future Enhancements

### High Priority
- [ ] Add tests (unit + integration)
- [ ] Implement actual API discovery
- [ ] Add authentication support
- [ ] Better error messages with suggestions

### Medium Priority
- [ ] React to messages
- [ ] Edit/delete messages
- [ ] Mark as read
- [ ] User information commands
- [ ] Media upload/download

### Low Priority
- [ ] Shell completion scripts
- [ ] Config file templates
- [ ] Verbose/debug mode
- [ ] Profile management (multiple accounts)

## Code Style

### Follow Go Conventions
```go
// Good
func (c *Client) ListChats() ([]Chat, error)

// Bad
func (c *Client) list_chats() ([]Chat, error)
```

### Comment Exported Functions
```go
// ListChats retrieves all conversations from the Beeper API.
// Returns an empty slice if no chats are found.
func (c *Client) ListChats() ([]Chat, error)
```

### Error Wrapping
```go
if err != nil {
    return nil, fmt.Errorf("failed to list chats: %w", err)
}
```

## Performance Considerations

### Current Performance
- No caching
- Synchronous HTTP requests
- JSON parsing per request

### Potential Optimizations
1. **Caching**: Cache chat list for N seconds
2. **Concurrent Requests**: Batch operations with goroutines
3. **Streaming**: For large message lists
4. **Connection Pooling**: Reuse HTTP connections

Example caching:
```go
type CachedClient struct {
    client *Client
    cache  map[string]cacheEntry
    mu     sync.RWMutex
}

func (c *CachedClient) ListChats() ([]Chat, error) {
    c.mu.RLock()
    if entry, ok := c.cache["chats"]; ok && !entry.Expired() {
        c.mu.RUnlock()
        return entry.Data.([]Chat), nil
    }
    c.mu.RUnlock()
    
    // Fetch and cache...
}
```

## Security Considerations

### Current Security
- Local-only API access
- No network exposure
- No credential storage

### If Adding Authentication
- Use system keychain (keyring library)
- Never log tokens/passwords
- Environment variable override
- Secure config file permissions (0600)

## Debugging

### Enable Verbose Mode (Future)
```bash
beeper --debug chats list
```

### Manual API Testing
```bash
# Test endpoints directly
curl http://localhost:39867/chats

# Check API response format
curl -s http://localhost:39867/chats | jq .
```

## Contributing Workflow

1. Fork repository
2. Create feature branch: `git checkout -b feature/name`
3. Make changes
4. Run tests: `go test ./...`
5. Format code: `go fmt ./...`
6. Commit: Clear, descriptive messages
7. Push and open PR

## Resources

- [Cobra Documentation](https://github.com/spf13/cobra)
- [Viper Documentation](https://github.com/spf13/viper)
- [Go HTTP Client Best Practices](https://pkg.go.dev/net/http)
- [Effective Go](https://go.dev/doc/effective_go)
