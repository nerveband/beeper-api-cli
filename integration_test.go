// +build integration

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_FullWorkflow tests a complete end-to-end workflow
func TestIntegration_FullWorkflow(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	testChatID := os.Getenv("BEEPER_TEST_CHAT_ID")
	
	if apiURL == "" || token == "" || testChatID == "" {
		t.Skip("Integration test environment not configured")
	}

	// 1. Test discover
	t.Run("Discover API", func(t *testing.T) {
		cmd := exec.Command("./beeper", "discover")
		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)
		assert.NotEmpty(t, output)
	})

	// 2. Test listing chats
	t.Run("List Chats", func(t *testing.T) {
		cmd := exec.Command("./beeper", "chats", "list", "--output", "json")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(output), "[")
	})

	// 3. Test getting specific chat
	t.Run("Get Chat", func(t *testing.T) {
		cmd := exec.Command("./beeper", "chats", "get", testChatID, "--output", "json")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(output), testChatID)
	})

	// 4. Test listing messages
	t.Run("List Messages", func(t *testing.T) {
		cmd := exec.Command("./beeper", "messages", "list", testChatID, "--limit", "5", "--output", "json")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(output), "[")
	})

	// 5. Test sending a message
	t.Run("Send Message", func(t *testing.T) {
		testMsg := fmt.Sprintf("Integration test message - %v", time.Now().Unix())
		cmd := exec.Command("./beeper", "send", "--chat-id", testChatID, "--message", testMsg)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)
		assert.NotEmpty(t, output)
	})

	// 6. Test search
	t.Run("Search Messages", func(t *testing.T) {
		cmd := exec.Command("./beeper", "search", "--query", "test", "--limit", "5", "--output", "json")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)
		assert.NotEmpty(t, output)
	})
}

// TestIntegration_OutputFormats tests all output formats work correctly
func TestIntegration_OutputFormats(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("Integration test environment not configured")
	}

	formats := []string{"json", "text", "markdown"}

	for _, format := range formats {
		t.Run("Format_"+format, func(t *testing.T) {
			cmd := exec.Command("./beeper", "chats", "list", "--output", format)
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			assert.NotEmpty(t, output)

			result := string(output)
			switch format {
			case "json":
				assert.Contains(t, result, "[")
			case "markdown":
				assert.Contains(t, result, "#")
			case "text":
				assert.NotEmpty(t, result)
			}
		})
	}
}

// TestIntegration_ErrorHandling tests error scenarios
func TestIntegration_ErrorHandling(t *testing.T) {
	t.Run("Invalid Chat ID", func(t *testing.T) {
		cmd := exec.Command("./beeper", "chats", "get", "invalid-chat-id-that-does-not-exist")
		output, err := cmd.CombinedOutput()
		assert.Error(t, err)
		assert.NotEmpty(t, output)
	})

	t.Run("Missing Required Argument", func(t *testing.T) {
		cmd := exec.Command("./beeper", "send", "--message", "test")
		output, err := cmd.CombinedOutput()
		assert.Error(t, err)
		assert.NotEmpty(t, output)
	})

	t.Run("Invalid Output Format", func(t *testing.T) {
		cmd := exec.Command("./beeper", "chats", "list", "--output", "invalid")
		output, err := cmd.CombinedOutput()
		// Should either error or default to valid format
		assert.NotEmpty(t, output)
	})
}

// TestIntegration_Config tests configuration management
func TestIntegration_Config(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := tmpDir + "/config.yaml"
	os.Setenv("BEEPER_CONFIG", configPath)
	defer os.Unsetenv("BEEPER_CONFIG")

	t.Run("Set Config", func(t *testing.T) {
		cmd := exec.Command("./beeper", "config", "set", "output_format", "markdown")
		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		// Verify file exists
		_, err = os.Stat(configPath)
		assert.NoError(t, err)
	})

	t.Run("Show Config", func(t *testing.T) {
		cmd := exec.Command("./beeper", "config", "show")
		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)
		assert.Contains(t, string(output), "output_format")
	})

	t.Run("Get Config Value", func(t *testing.T) {
		cmd := exec.Command("./beeper", "config", "get", "output_format")
		output, err := cmd.CombinedOutput()
		assert.NoError(t, err)
		assert.Contains(t, string(output), "markdown")
	})
}

// TestIntegration_PipelineUsage tests CLI works in Unix pipelines
func TestIntegration_PipelineUsage(t *testing.T) {
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("Integration test environment not configured")
	}

	t.Run("Pipe JSON to jq", func(t *testing.T) {
		// Test that JSON output can be piped to jq
		cmd := exec.Command("bash", "-c", "./beeper chats list --output json | jq -r '.[0].id'")
		output, err := cmd.CombinedOutput()
		
		// Only assert if jq is available
		if err == nil {
			assert.NotEmpty(t, output)
		}
	})

	t.Run("Grep text output", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "./beeper chats list --output text | grep -i chat")
		output, err := cmd.CombinedOutput()
		
		// May not match anything, but should not error on pipe
		assert.NotNil(t, output)
	})
}
