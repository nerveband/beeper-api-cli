package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/api"
	"github.com/nerveband/beeper-api-cli/internal/config"
)

var testPermissions bool

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display CLI and environment information",
	Long: `Display comprehensive information about the CLI, configuration,
and API connectivity status. Useful for troubleshooting and verification.

The --test-permissions flag tests actual API access to verify your token works.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInfo()
	},
}

func init() {
	infoCmd.Flags().BoolVar(&testPermissions, "test-permissions", false, "Test API token permissions by making test requests")
	rootCmd.AddCommand(infoCmd)
}

func runInfo() error {
	fmt.Println("Beeper API CLI Information")
	fmt.Println("==========================")
	fmt.Println()

	// CLI Version
	fmt.Printf("Version:        %s\n", Version)
	fmt.Printf("Go Version:     %s\n", runtime.Version())
	fmt.Printf("Platform:       %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()

	// Configuration
	fmt.Println("Configuration")
	fmt.Println("-------------")
	configPath := config.GetConfigPath()
	fmt.Printf("Config File:    %s\n", configPath)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("                (not created yet, using defaults)\n")
	}

	fmt.Printf("API URL:        %s\n", cfg.APIURL)
	fmt.Printf("Output Format:  %s\n", cfg.OutputFormat)
	fmt.Println()

	// Authentication
	fmt.Println("Authentication")
	fmt.Println("--------------")
	token := os.Getenv("BEEPER_TOKEN")
	if token != "" {
		// Mask token for display
		maskedToken := token[:4] + "..." + token[len(token)-4:]
		if len(token) < 12 {
			maskedToken = "****"
		}
		fmt.Printf("BEEPER_TOKEN:   Set (%s)\n", maskedToken)
	} else {
		fmt.Printf("BEEPER_TOKEN:   Not set\n")
		fmt.Printf("                (Set this environment variable to authenticate API requests)\n")
	}
	fmt.Println()

	// API Connectivity
	fmt.Println("API Connectivity")
	fmt.Println("----------------")
	client := getAPIClient()

	if err := client.Ping(); err != nil {
		fmt.Printf("Status:         Unreachable\n")
		fmt.Printf("Error:          %v\n", err)
		if !quietMode {
			fmt.Println()
			fmt.Println("Hint: Make sure Beeper Desktop is running and the API is enabled.")
			fmt.Println("      Try 'beeper discover' to find the API endpoint.")
		}
	} else {
		fmt.Printf("Status:         Connected\n")
		if version := client.GetDesktopVersion(); version != "" {
			fmt.Printf("Desktop Ver:    %s\n", version)
		}
	}

	// Permission Testing
	if testPermissions {
		fmt.Println()
		fmt.Println("Permission Test")
		fmt.Println("---------------")
		testAPIPermissions(client)
	}

	return nil
}

func testAPIPermissions(client *api.Client) {
	// Test read permission by listing chats
	fmt.Print("Read (list chats):  ")
	_, err := client.ListChats()
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			switch apiErr.Category {
			case api.CategoryAuth:
				fmt.Println("FAILED - Authentication required")
			case api.CategoryPermission:
				fmt.Println("FAILED - Insufficient permissions")
			case api.CategoryNetwork:
				fmt.Println("FAILED - Network error")
			default:
				fmt.Printf("FAILED - %s\n", apiErr.Message)
			}
		} else {
			fmt.Printf("FAILED - %v\n", err)
		}
	} else {
		fmt.Println("OK")
	}

	// Test search capability
	fmt.Print("Search messages:    ")
	_, err = client.SearchMessages("test", 1)
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			switch apiErr.Category {
			case api.CategoryAuth:
				fmt.Println("FAILED - Authentication required")
			case api.CategoryPermission:
				fmt.Println("FAILED - Insufficient permissions")
			case api.CategoryNetwork:
				fmt.Println("FAILED - Network error")
			default:
				fmt.Printf("FAILED - %s\n", apiErr.Message)
			}
		} else {
			fmt.Printf("FAILED - %v\n", err)
		}
	} else {
		fmt.Println("OK")
	}

	fmt.Println()
	fmt.Println("Note: Write permissions cannot be tested without making actual changes.")
}
