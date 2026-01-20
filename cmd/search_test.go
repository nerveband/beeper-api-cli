package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSearchCommand tests the search command
func TestSearchCommand(t *testing.T) {
	t.Skip("Integration test - requires live Beeper API")
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping search test")
	}

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"search", "--query", "test", "--limit", "5", "--output", "json"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
}

// TestSearchCommand_Text tests text output for search
func TestSearchCommand_Text(t *testing.T) {
	t.Skip("Integration test - requires live Beeper API")
	apiURL := os.Getenv("BEEPER_API_URL")
	token := os.Getenv("BEEPER_TOKEN")
	
	if apiURL == "" || token == "" {
		t.Skip("BEEPER_API_URL or BEEPER_TOKEN not set - skipping search test")
	}

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"search", "--query", "hello", "--limit", "10", "--output", "text"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
}

// TestSearchCommand_MissingQuery tests error handling for missing query
func TestSearchCommand_MissingQuery(t *testing.T) {
	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"search"})

	err := rootCmd.Execute()
	assert.Error(t, err)
}

// TestSearchCommand_EmptyQuery tests handling of empty query
func TestSearchCommand_EmptyQuery(t *testing.T) {
	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"search", "--query", ""})

	err := rootCmd.Execute()
	assert.Error(t, err)
}
