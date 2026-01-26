package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles communication with Beeper Desktop API
type Client struct {
	baseURL        string
	authToken      string
	httpClient     *http.Client
	desktopVersion string // Cached from X-Beeper-Desktop-Version header
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetAuthToken sets the Bearer authentication token
func (c *Client) SetAuthToken(token string) {
	c.authToken = token
}

// Ping checks if the API is reachable
func (c *Client) Ping() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return WrapNetworkError(err, "ping")
	}
	defer resp.Body.Close()

	// Capture desktop version from header
	if version := resp.Header.Get("X-Beeper-Desktop-Version"); version != "" {
		c.desktopVersion = version
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return NewAPIErrorFromStatus(resp.StatusCode, body, "ping")
	}

	return nil
}

// GetDesktopVersion returns the cached Beeper Desktop version from the last API response
func (c *Client) GetDesktopVersion() string {
	return c.desktopVersion
}

// GetBaseURL returns the base URL of the API
func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// HasAuthToken returns true if an auth token is set
func (c *Client) HasAuthToken() bool {
	return c.authToken != ""
}

// doRequest performs an HTTP request and returns the response body
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	return c.doRequestWithOp(method, path, body, "")
}

// doRequestWithOp performs an HTTP request with operation context for error messages
func (c *Client) doRequestWithOp(method, path string, body interface{}, operation string) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, &APIError{
				Message:    fmt.Sprintf("failed to marshal request body: %v", err),
				Category:   CategoryValidation,
				Operation:  operation,
				Underlying: err,
			}
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to create request: %v", err),
			Category:   CategoryConfig,
			Operation:  operation,
			Underlying: err,
		}
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add auth token if set
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, WrapNetworkError(err, operation)
	}
	defer resp.Body.Close()

	// Capture desktop version from header
	if version := resp.Header.Get("X-Beeper-Desktop-Version"); version != "" {
		c.desktopVersion = version
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to read response body: %v", err),
			Category:   CategoryNetwork,
			Operation:  operation,
			Underlying: err,
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, NewAPIErrorFromStatus(resp.StatusCode, respBody, operation)
	}

	return respBody, nil
}

// Chat represents a Beeper chat/conversation
// ChatsResponse represents the API response for listing chats
type ChatsResponse struct {
	Items  []Chat `json:"items"`
	HasMore bool `json:"hasMore"`
}

// ListChats retrieves all chats
func (c *Client) ListChats() ([]Chat, error) {
	data, err := c.doRequestWithOp("GET", "/v1/chats", nil, "list_chats")
	if err != nil {
		return nil, err
	}

	var resp ChatsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to unmarshal chats: %v", err),
			Category:   CategoryServer,
			Operation:  "list_chats",
			Underlying: err,
		}
	}

	return resp.Items, nil
}

// GetChat retrieves a specific chat by ID
func (c *Client) GetChat(chatID string) (*Chat, error) {
	data, err := c.doRequestWithOp("GET", "/v1/chats/"+chatID, nil, "get_chat")
	if err != nil {
		return nil, err
	}

	var chat Chat
	if err := json.Unmarshal(data, &chat); err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to unmarshal chat: %v", err),
			Category:   CategoryServer,
			Operation:  "get_chat",
			Underlying: err,
		}
	}

	return &chat, nil
}

// MessagesResponse represents the API response for listing messages
type MessagesResponse struct {
	Items []Message `json:"items"`
	HasMore bool `json:"hasMore"`
}

// ListMessages retrieves messages from a chat
func (c *Client) ListMessages(chatID string, limit int) ([]Message, error) {
	path := fmt.Sprintf("/v1/chats/%s/messages?limit=%d", chatID, limit)
	data, err := c.doRequestWithOp("GET", path, nil, "list_messages")
	if err != nil {
		return nil, err
	}

	var resp MessagesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to unmarshal messages: %v", err),
			Category:   CategoryServer,
			Operation:  "list_messages",
			Underlying: err,
		}
	}

	return resp.Items, nil
}

// SendMessage sends a message to a chat and returns the message ID
func (c *Client) SendMessage(chatID, message string) (string, error) {
	req := SendMessageRequest{
		Text: message,
	}

	data, err := c.doRequestWithOp("POST", "/v1/chats/"+chatID+"/messages", req, "send_message")
	if err != nil {
		return "", err
	}

	var resp SendMessageResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", &APIError{
			Message:    fmt.Sprintf("failed to unmarshal response: %v", err),
			Category:   CategoryServer,
			Operation:  "send_message",
			Underlying: err,
		}
	}

	return resp.ID, nil
}

// SearchResponse represents the API response for searching messages
type SearchResponse struct {
	Items []Message `json:"items"`
}

// SearchMessages searches for messages across all chats
func (c *Client) SearchMessages(query string, limit int) ([]Message, error) {
	path := fmt.Sprintf("/v1/messages/search?q=%s&limit=%d", query, limit)
	data, err := c.doRequestWithOp("GET", path, nil, "search_messages")
	if err != nil {
		return nil, err
	}

	var resp SearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to unmarshal messages: %v", err),
			Category:   CategoryServer,
			Operation:  "search_messages",
			Underlying: err,
		}
	}

	return resp.Items, nil
}

// DiscoverAPI attempts to auto-discover the Beeper Desktop API URL
func DiscoverAPI() (string, error) {
	// Try common ports
	ports := []int{39867, 39868, 39869}
	for _, port := range ports {
		url := fmt.Sprintf("http://localhost:%d", port)
		client := NewClient(url)
		if err := client.Ping(); err == nil {
			return url, nil
		}
	}

	return "", fmt.Errorf("could not auto-discover Beeper Desktop API")
}
