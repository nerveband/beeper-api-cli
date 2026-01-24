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
	baseURL    string
	authToken  string
	httpClient *http.Client
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
		return fmt.Errorf("failed to ping API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

// doRequest performs an HTTP request and returns the response body
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
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
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
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
	data, err := c.doRequest("GET", "/v1/chats", nil)
	if err != nil {
		return nil, err
	}

	var resp ChatsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chats: %w", err)
	}

	return resp.Items, nil
}

// GetChat retrieves a specific chat by ID
func (c *Client) GetChat(chatID string) (*Chat, error) {
	data, err := c.doRequest("GET", "/v1/chats/"+chatID, nil)
	if err != nil {
		return nil, err
	}

	var chat Chat
	if err := json.Unmarshal(data, &chat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chat: %w", err)
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
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp MessagesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
	}

	return resp.Items, nil
}

// SendMessage sends a message to a chat and returns the message ID
func (c *Client) SendMessage(chatID, message string) (string, error) {
	req := SendMessageRequest{
		Text: message,
	}

	data, err := c.doRequest("POST", "/v1/chats/"+chatID+"/messages", req)
	if err != nil {
		return "", err
	}

	var resp SendMessageResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
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
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp SearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
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
