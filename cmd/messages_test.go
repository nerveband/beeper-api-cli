package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMessagesListCommand tests the messages list command
func TestMessagesListCommand(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	testChatID := os.Getenv("BEEPER_TEST_CHAT_ID")
	
	if apiURL == "" || token == "" || testChatID == "" {
		t.Skip("BEEPER_API_URL, BEEPER_TOKEN, or BEEPER_TEST_CHAT_ID not set - skipping test")
	}

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"messages", "list", testChatID, "--limit", "10", "--output", "json"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "[") // JSON array
}

// TestMessagesListCommand_Text tests text output for messages
func TestMessagesListCommand_Text(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	testChatID := os.Getenv("BEEPER_TEST_CHAT_ID")
	
	if apiURL == "" || token == "" || testChatID == "" {
		t.Skip("BEEPER_API_URL, BEEPER_TOKEN, or BEEPER_TEST_CHAT_ID not set - skipping test")
	}

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"messages", "list", testChatID, "--limit", "5", "--output", "text"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
}

// TestMessagesListCommand_Limit tests limit parameter
func TestMessagesListCommand_Limit(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	testChatID := os.Getenv("BEEPER_TEST_CHAT_ID")
	
	if apiURL == "" || token == "" || testChatID == "" {
		t.Skip("BEEPER_API_URL, BEEPER_TOKEN, or BEEPER_TEST_CHAT_ID not set - skipping test")
	}

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"messages", "list", testChatID, "--limit", "3", "--output", "json"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
}
