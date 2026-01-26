package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/api"
	"github.com/nerveband/beeper-api-cli/internal/config"
	"github.com/nerveband/beeper-api-cli/internal/output"
	"github.com/nerveband/beeper-api-cli/internal/update"
)

var (
	cfg           *config.Config
	outputFormat  string
	quietMode     bool
	jsonErrors    bool
	updateCheckCh <-chan *update.UpdateInfo
	// Version is set at build time via ldflags
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "beeper",
	Short: "Beeper API CLI - Command-line interface for Beeper Desktop API",
	Long:  BannerWithVersion() + "\nA cross-platform CLI for the Beeper Desktop API.\nProvides LLM-friendly interfaces for reading and sending messages\nacross all Beeper-supported chat networks.",
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

		// Start async update check (skip for version, upgrade, and help commands)
		cmdName := cmd.Name()
		if !quietMode && cmdName != "version" && cmdName != "upgrade" && cmdName != "help" {
			updateCheckCh = update.CheckAsync(Version)
		}

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Display update notification if available
		if updateCheckCh != nil && !quietMode {
			select {
			case info := <-updateCheckCh:
				if notice := update.FormatUpdateNotice(info); notice != "" {
					fmt.Fprint(os.Stderr, notice)
				}
			default:
				// Don't block if check hasn't completed
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}

// helpFooter contains documentation links shown at the bottom of help output
const helpFooter = `
Documentation & Support:
  GitHub:  https://github.com/nerveband/beeper-api-cli
  Issues:  https://github.com/nerveband/beeper-api-cli/issues
  API:     https://github.com/nerveband/beeper-api-cli/blob/main/API.md
`

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (json, text, markdown)")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Suppress non-essential output (hints, update notifications)")
	rootCmd.PersistentFlags().BoolVar(&jsonErrors, "json-errors", false, "Output errors as JSON to stderr")

	// Add help footer with documentation links
	defaultUsageTemplate := rootCmd.UsageTemplate()
	rootCmd.SetUsageTemplate(defaultUsageTemplate + helpFooter)
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
