package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	DefaultAPIURL      = "http://localhost:39867"
	DefaultOutputFormat = "json"
)

type Config struct {
	APIURL       string `mapstructure:"api_url"`
	OutputFormat string `mapstructure:"output_format"`
}

// Load reads configuration from ~/.beeper-api-cli/config.yaml
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".beeper-api-cli")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("api_url", DefaultAPIURL)
	viper.SetDefault("output_format", DefaultOutputFormat)

	// Read config file if it exists (ignore if not found, use defaults)
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Save writes the current configuration to disk
func Save(cfg *Config) error {
	viper.Set("api_url", cfg.APIURL)
	viper.Set("output_format", cfg.OutputFormat)

	if err := viper.WriteConfig(); err != nil {
		// If config file doesn't exist yet, create it
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.SafeWriteConfig()
		}
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// LoadConfig loads configuration from a specific file path
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set defaults
	v.SetDefault("api_url", DefaultAPIURL)
	v.SetDefault("output_format", DefaultOutputFormat)

	// Read config file if it exists (ignore file not found errors)
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves configuration to a specific file path
func SaveConfig(configPath string, cfg *Config) error {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.Set("api_url", cfg.APIURL)
	v.Set("output_format", cfg.OutputFormat)

	// Create parent directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		APIURL:       DefaultAPIURL,
		OutputFormat: DefaultOutputFormat,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.APIURL == "" {
		return fmt.Errorf("api_url cannot be empty")
	}

	validFormats := map[string]bool{
		"json":     true,
		"text":     true,
		"markdown": true,
	}

	if !validFormats[c.OutputFormat] {
		return fmt.Errorf("invalid output format: %s (must be json, text, or markdown)", c.OutputFormat)
	}

	return nil
}

// Merge merges override config into base config (non-empty values from override take precedence)
func (c *Config) Merge(override *Config) *Config {
	merged := &Config{
		APIURL:       c.APIURL,
		OutputFormat: c.OutputFormat,
	}

	if override.APIURL != "" {
		merged.APIURL = override.APIURL
	}
	if override.OutputFormat != "" {
		merged.OutputFormat = override.OutputFormat
	}

	return merged
}

// GetConfigPath returns the default config file path
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".beeper-api-cli", "config.yaml")
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	cfg := DefaultConfig()

	if apiURL := os.Getenv("BEEPER_API_URL"); apiURL != "" {
		cfg.APIURL = apiURL
	}
	if format := os.Getenv("BEEPER_OUTPUT_FORMAT"); format != "" {
		cfg.OutputFormat = format
	}

	return cfg
}

// UpdateConfig updates specific fields in an existing config file
func UpdateConfig(configPath string, update *Config) error {
	// Load existing config
	existing, err := LoadConfig(configPath)
	if err != nil {
		existing = DefaultConfig()
	}

	// Merge updates
	merged := existing.Merge(update)

	// Save merged config
	return SaveConfig(configPath, merged)
}
