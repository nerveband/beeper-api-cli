package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/api"
	"github.com/nerveband/beeper-api-cli/internal/config"
	"github.com/nerveband/beeper-api-cli/internal/output"
)

var (
	cfg          *config.Config
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:   "beeper",
	Short: "Beeper API CLI - Command-line interface for Beeper Desktop API",
	Long: `A cross-platform CLI for the Beeper Desktop API.
Provides LLM-friendly interfaces for reading and sending messages
across all Beeper-supported chat networks.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Override with environment variables if set
		envCfg := config.LoadFromEnv()
		cfg = cfg.Merge(envCfg)

		// Override output format if flag is set
		if outputFormat != "" {
			cfg.OutputFormat = outputFormat
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (json, text, markdown)")
}

// getOutputFormat returns the configured output format
func getOutputFormat() output.Format {
	return output.Format(cfg.OutputFormat)
}

// getAPIClient returns an API client with auth token
func getAPIClient() *api.Client {
	client := api.NewClient(cfg.APIURL)
	// Set auth token from environment variable
	if token := os.Getenv("BEEPER_TOKEN"); token != "" {
		client.SetAuthToken(token)
	}
	return client
}
