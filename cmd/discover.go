package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/api"
	"github.com/nerveband/beeper-api-cli/internal/config"
)

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Auto-discover Beeper Desktop API endpoint",
	Long: `Attempts to automatically discover the Beeper Desktop API URL
by trying common localhost ports. If successful, saves the URL to config.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Discovering Beeper Desktop API...")

		apiURL, err := api.DiscoverAPI()
		if err != nil {
			return fmt.Errorf("discovery failed: %w\n\nPlease ensure Beeper Desktop is running and try manually:\n  beeper config set-url <url>", err)
		}

		fmt.Printf("Found Beeper Desktop API at: %s\n", apiURL)

		// Update config
		cfg.APIURL = apiURL
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println("Configuration saved successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}
