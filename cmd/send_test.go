package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSendCommand tests sending a message
func TestSendCommand(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	testChatID := os.Getenv("BEEPER_TEST_CHAT_ID")
	
	if apiURL == "" || token == "" || testChatID == "" {
		t.Skip("BEEPER_API_URL, BEEPER_TOKEN, or BEEPER_TEST_CHAT_ID not set - skipping send test")
	}

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	testMessage := "Test message from Beeper CLI test suite"
	rootCmd.SetArgs([]string{"send", "--chat-id", testChatID, "--message", testMessage})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "sent") // Success message
}

// TestSendCommand_MissingChatID tests error handling for missing chat ID
func TestSendCommand_MissingChatID(t *testing.T) {
	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"send", "--message", "test"})

	err := rootCmd.Execute()
	assert.Error(t, err)
}

// TestSendCommand_MissingMessage tests error handling for missing message
func TestSendCommand_MissingMessage(t *testing.T) {
	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"send", "--chat-id", "test-chat"})

	err := rootCmd.Execute()
	assert.Error(t, err)
}
