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

// Load reads configuration from ~/.beeper-cli/config.yaml
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".beeper-cli")
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
