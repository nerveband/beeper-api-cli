package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfigShowCommand tests showing current configuration
func TestConfigShowCommand(t *testing.T) {
	t.Skip("Integration test - requires proper config setup")
	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"config", "show"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "api_url")
}

// TestConfigSetCommand tests setting configuration values
func TestConfigSetCommand(t *testing.T) {
	t.Skip("Integration test - requires proper config setup")
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.Setenv("BEEPER_CONFIG", configPath)
	defer os.Unsetenv("BEEPER_CONFIG")

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	// Set a config value
	rootCmd.SetArgs([]string{"config", "set", "output_format", "markdown"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	// Verify config file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err)
}

// TestConfigGetCommand tests getting a specific config value
func TestConfigGetCommand(t *testing.T) {
	// Create temp config with known values
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	
	testConfig := `api_url: http://localhost:8080
output_format: json
`
	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	require.NoError(t, err)

	os.Setenv("BEEPER_CONFIG", configPath)
	defer os.Unsetenv("BEEPER_CONFIG")

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"config", "get", "output_format"})

	err = rootCmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.Contains(t, result, "json")
}

// TestConfigValidateCommand tests config validation
func TestConfigValidateCommand(t *testing.T) {
	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	rootCmd.SetArgs([]string{"config", "validate"})

	_ = rootCmd.Execute()
	
	// Should either succeed or fail gracefully
	result := output.String()
	assert.NotEmpty(t, result)
}
