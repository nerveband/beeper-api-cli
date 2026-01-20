package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify Beeper CLI configuration settings.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("API URL: %s\n", cfg.APIURL)
		fmt.Printf("Output Format: %s\n", cfg.OutputFormat)
		return nil
	},
}

var configSetURLCmd = &cobra.Command{
	Use:   "set-url <url>",
	Short: "Set the Beeper Desktop API URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.APIURL = args[0]
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Printf("API URL set to: %s\n", cfg.APIURL)
		return nil
	},
}

var configSetFormatCmd = &cobra.Command{
	Use:   "set-format <format>",
	Short: "Set the default output format (json, text, markdown)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format := args[0]
		if format != "json" && format != "text" && format != "markdown" {
			return fmt.Errorf("invalid format: %s (must be json, text, or markdown)", format)
		}

		cfg.OutputFormat = format
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Printf("Output format set to: %s\n", cfg.OutputFormat)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetURLCmd)
	configCmd.AddCommand(configSetFormatCmd)
	rootCmd.AddCommand(configCmd)
}
