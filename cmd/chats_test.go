package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestChatsListCommand tests the chats list command
func TestChatsListCommand(t *testing.T) {
	t.Skip("Integration test - requires live Beeper API")
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping live test")
	}

	// Capture output
	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	// Set args
	rootCmd.SetArgs([]string{"chats", "list", "--output", "json"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "[") // JSON array
}

// TestChatsListCommand_Text tests text output format
func TestChatsListCommand_Text(t *testing.T) {
	t.Skip("Integration test - requires live Beeper API")
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping live test")
	}

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"chats", "list", "--output", "text"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
}

// TestChatsListCommand_Markdown tests markdown output format
func TestChatsListCommand_Markdown(t *testing.T) {
	t.Skip("Integration test - requires live Beeper API")
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping live test")
	}

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"chats", "list", "--output", "markdown"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "#") // Markdown headers
}

// TestChatsGetCommand tests getting a specific chat
func TestChatsGetCommand(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	testChatID := os.Getenv("BEEPER_TEST_CHAT_ID")
	
	if apiURL == "" || token == "" || testChatID == "" {
		t.Skip("BEEPER_API_URL, BEEPER_TOKEN, or BEEPER_TEST_CHAT_ID not set - skipping test")
	}

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"chats", "get", testChatID, "--output", "json"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, testChatID)
}
