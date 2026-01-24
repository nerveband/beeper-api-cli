package cmd

import (
	"context"
	"fmt"
	"runtime"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

// Configure these for your repo
const repoOwner = "nerveband"
const repoName = "beeper-api-cli"

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to the latest version",
	Long:  "Check for and install the latest version from GitHub releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpgrade()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func runUpgrade() error {
	fmt.Printf("Current version: %s\n", Version)
	fmt.Println("Checking for updates...")

	// Create GitHub source (no auth needed for public repos)
	source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})
	if err != nil {
		return fmt.Errorf("failed to create update source: %w", err)
	}

	// Create updater (checksum validation optional - only if checksums.txt exists)
	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		Source: source,
	})
	if err != nil {
		return fmt.Errorf("failed to create updater: %w", err)
	}

	// Check for latest release
	latest, found, err := updater.DetectLatest(
		context.Background(),
		selfupdate.NewRepositorySlug(repoOwner, repoName),
	)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !found {
		fmt.Println("No releases found")
		return nil
	}

	// Compare versions (skip comparison for dev builds - always offer upgrade)
	if Version != "dev" {
		if latest.LessOrEqual(Version) {
			fmt.Printf("Already up to date (latest: %s)\n", latest.Version())
			return nil
		}
	} else {
		fmt.Println("Running development build")
	}

	// Download and install
	fmt.Printf("New version available: %s\n", latest.Version())
	fmt.Printf("Downloading for %s/%s...\n", runtime.GOOS, runtime.GOARCH)

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	if err := updater.UpdateTo(context.Background(), latest, exe); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	fmt.Printf("Successfully upgraded to %s\n", latest.Version())
	return nil
}
