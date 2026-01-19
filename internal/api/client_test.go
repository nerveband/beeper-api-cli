package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_NewClient tests client initialization
func TestClient_NewClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewClient(baseURL)

	assert.NotNil(t, client)
	assert.Equal(t, baseURL, client.baseURL)
	assert.NotNil(t, client.httpClient)
}

// TestClient_ListChats tests listing chats against real Beeper API
func TestClient_ListChats(t *testing.T) {
	// Skip if no Beeper Desktop API available
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping live API test")
	}

	client := NewClient(apiURL)
	client.SetAuthToken(token)

	chats, err := client.ListChats()
	
	require.NoError(t, err, "ListChats should not return an error")
	assert.NotNil(t, chats, "Chats list should not be nil")
	
	// If we have chats, validate structure
	if len(chats) > 0 {
		firstChat := chats[0]
		assert.NotEmpty(t, firstChat.ID, "Chat ID should not be empty")
		// Name can be empty for some chats
		assert.NotNil(t, firstChat.Participants, "Participants should not be nil")
	}
}

// TestClient_GetChat tests getting a specific chat
func TestClient_GetChat(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping live API test")
	}

	client := NewClient(apiURL)
	client.SetAuthToken(token)

	// First get a chat ID
	chats, err := client.ListChats()
	require.NoError(t, err)
	
	if len(chats) == 0 {
		t.Skip("No chats available to test GetChat")
	}

	chatID := chats[0].ID
	chat, err := client.GetChat(chatID)
	
	require.NoError(t, err, "GetChat should not return an error")
	assert.NotNil(t, chat, "Chat should not be nil")
	assert.Equal(t, chatID, chat.ID, "Chat ID should match requested ID")
}

// TestClient_ListMessages tests listing messages from a chat
func TestClient_ListMessages(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping live API test")
	}

	client := NewClient(apiURL)
	client.SetAuthToken(token)

	// Get first available chat
	chats, err := client.ListChats()
	require.NoError(t, err)
	
	if len(chats) == 0 {
		t.Skip("No chats available to test ListMessages")
	}

	chatID := chats[0].ID
	messages, err := client.ListMessages(chatID, 10)
	
	require.NoError(t, err, "ListMessages should not return an error")
	assert.NotNil(t, messages, "Messages list should not be nil")
	
	// Validate message structure if we have messages
	if len(messages) > 0 {
		firstMsg := messages[0]
		assert.NotEmpty(t, firstMsg.ID, "Message ID should not be empty")
		assert.NotEmpty(t, firstMsg.Text, "Message text should not be empty")
		assert.NotEmpty(t, firstMsg.Sender, "Message sender should not be empty")
	}
}

// TestClient_SendMessage tests sending a message to a chat
func TestClient_SendMessage(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	testChatID := os.Getenv("BEEPER_TEST_CHAT_ID") // Set this to a safe test chat
	
	if apiURL == "" || token == "" || testChatID == "" {
		t.Skip("BEEPER_API_URL, BEEPER_TOKEN, or BEEPER_TEST_CHAT_ID not set - skipping send test")
	}

	client := NewClient(apiURL)
	client.SetAuthToken(token)

	testMessage := "Test message from Beeper CLI test suite"
	messageID, err := client.SendMessage(testChatID, testMessage)
	
	require.NoError(t, err, "SendMessage should not return an error")
	assert.NotEmpty(t, messageID, "Message ID should not be empty")
}

// TestClient_SearchMessages tests message search functionality
func TestClient_SearchMessages(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping live API test")
	}

	client := NewClient(apiURL)
	client.SetAuthToken(token)

	// Search for a common word
	results, err := client.SearchMessages("test", 5)
	
	require.NoError(t, err, "SearchMessages should not return an error")
	assert.NotNil(t, results, "Search results should not be nil")
	
	// Validate search result structure if we have results
	if len(results) > 0 {
		firstResult := results[0]
		assert.NotEmpty(t, firstResult.ID, "Message ID should not be empty")
	}
}

// TestClient_InvalidURL tests error handling for invalid API URL
func TestClient_InvalidURL(t *testing.T) {
	client := NewClient("http://invalid-url-that-does-not-exist:99999")
	
	_, err := client.ListChats()
	assert.Error(t, err, "Should return error for invalid URL")
}

// TestClient_Ping tests API health check
func TestClient_Ping(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping live API test")
	}

	client := NewClient(apiURL)
	client.SetAuthToken(token)

	err := client.Ping()
	assert.NoError(t, err, "Ping should succeed with valid API URL")
}
