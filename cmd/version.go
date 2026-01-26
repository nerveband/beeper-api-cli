package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/config"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Display the current version, build information, and platform details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("beeper-api-cli version %s\n", Version)
		fmt.Printf("  Go version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  Config:     %s\n", config.GetConfigPath())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
