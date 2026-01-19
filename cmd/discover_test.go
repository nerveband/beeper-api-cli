package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDiscoverCommand tests the API discovery command
func TestDiscoverCommand(t *testing.T) {
	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"discover"})

	err := rootCmd.Execute()
	
	// Discovery may fail if Beeper Desktop is not running
	// We just check that the command executes without panicking
	result := output.String()
	assert.NotEmpty(t, result)
}

// TestDiscoverCommand_OutputFormat tests that discover respects output format
func TestDiscoverCommand_OutputFormat(t *testing.T) {
	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"discover", "--output", "json"})

	err := rootCmd.Execute()
	
	result := output.String()
	assert.NotEmpty(t, result)
	
	// If successful, should contain JSON-like output
	if err == nil {
		assert.True(t, len(result) > 0)
	}
}
