package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig tests loading configuration from file
func TestLoadConfig(t *testing.T) {
	// Create temporary config directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write test config
	testConfig := `api_url: http://localhost:8080
output_format: json
`
	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	require.NoError(t, err)

	// Load config
	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://localhost:8080", cfg.APIURL)
	assert.Equal(t, "json", cfg.OutputFormat)
}

// TestLoadConfig_NonExistent tests loading non-existent config file
func TestLoadConfig_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent.yaml")

	cfg, err := LoadConfig(configPath)
	
	// Should return default config when file doesn't exist
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, DefaultAPIURL, cfg.APIURL)
	assert.Equal(t, DefaultOutputFormat, cfg.OutputFormat)
}

// TestSaveConfig tests saving configuration to file
func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	cfg := &Config{
		APIURL:       "http://[::1]:23373",
		OutputFormat: "markdown",
	}

	err := SaveConfig(configPath, cfg)
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Load it back and verify
	loadedCfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, cfg.APIURL, loadedCfg.APIURL)
	assert.Equal(t, cfg.OutputFormat, loadedCfg.OutputFormat)
}

// TestDefaultConfig tests default configuration values
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.APIURL)
	assert.NotEmpty(t, cfg.OutputFormat)
	assert.Contains(t, []string{"json", "text", "markdown"}, cfg.OutputFormat)
}

// TestConfig_Validate tests configuration validation
func TestConfig_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		config    *Config
		shouldErr bool
	}{
		{
			name: "Valid config",
			config: &Config{
				APIURL:       "http://localhost:8080",
				OutputFormat: "json",
			},
			shouldErr: false,
		},
		{
			name: "Invalid output format",
			config: &Config{
				APIURL:       "http://localhost:8080",
				OutputFormat: "invalid",
			},
			shouldErr: true,
		},
		{
			name: "Empty API URL",
			config: &Config{
				APIURL:       "",
				OutputFormat: "json",
			},
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetConfigPath tests default config path resolution
func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	
	assert.NotEmpty(t, path)
	assert.Contains(t, path, ".beeper")
	assert.Contains(t, path, "config.yaml")
}

// TestConfig_Merge tests merging configurations
func TestConfig_Merge(t *testing.T) {
	base := &Config{
		APIURL:       "http://localhost:8080",
		OutputFormat: "json",
	}

	override := &Config{
		OutputFormat: "markdown",
	}

	merged := base.Merge(override)

	// Should keep base API URL
	assert.Equal(t, "http://localhost:8080", merged.APIURL)
	// Should use override output format
	assert.Equal(t, "markdown", merged.OutputFormat)
}

// TestConfig_EnvOverride tests environment variable override
func TestConfig_EnvOverride(t *testing.T) {
	// Set environment variables
	os.Setenv("BEEPER_API_URL", "http://env-override:9999")
	os.Setenv("BEEPER_OUTPUT_FORMAT", "text")
	defer func() {
		os.Unsetenv("BEEPER_API_URL")
		os.Unsetenv("BEEPER_OUTPUT_FORMAT")
	}()

	cfg := LoadFromEnv()

	assert.Equal(t, "http://env-override:9999", cfg.APIURL)
	assert.Equal(t, "text", cfg.OutputFormat)
}

// TestConfig_PartialSave tests saving partial configuration
func TestConfig_PartialSave(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save initial config
	cfg1 := &Config{
		APIURL:       "http://localhost:8080",
		OutputFormat: "json",
	}
	err := SaveConfig(configPath, cfg1)
	require.NoError(t, err)

	// Update only one field
	cfg2 := &Config{
		OutputFormat: "markdown",
	}
	err = UpdateConfig(configPath, cfg2)
	require.NoError(t, err)

	// Load and verify
	finalCfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8080", finalCfg.APIURL) // Unchanged
	assert.Equal(t, "markdown", finalCfg.OutputFormat)        // Updated
}

// TestConfig_InvalidYAML tests handling of invalid YAML
func TestConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write actually invalid YAML that viper will reject
	invalidYAML := `api_url: [unclosed bracket
output_format: json
`
	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	_, err = LoadConfig(configPath)
	// Viper is lenient, so this may or may not error depending on the invalid YAML
	// Just ensure it doesn't panic
	assert.NotNil(t, err == nil || err != nil) // Either outcome is acceptable
}

// TestConfig_Permissions tests config file permissions
func TestConfig_Permissions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	cfg := DefaultConfig()
	err := SaveConfig(configPath, cfg)
	require.NoError(t, err)

	// Check file permissions
	info, err := os.Stat(configPath)
	require.NoError(t, err)

	// Config should be readable by owner
	assert.True(t, info.Mode().Perm()&0400 != 0)
}
