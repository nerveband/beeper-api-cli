package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/output"
)

var (
	searchLimit int
)

var searchCmd = &cobra.Command{
	Use:   "search --query <text>",
	Short: "Search messages across all chats",
	Long:  `Search for messages containing the specified query text across all Beeper chats.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		if query == "" {
			return fmt.Errorf("--query is required")
		}

		client := getAPIClient()

		messages, err := client.SearchMessages(query, searchLimit)
		if err != nil {
			return fmt.Errorf("failed to search messages: %w", err)
		}

		formatted := output.FormatMessages(messages, getOutputFormat())
		fmt.Print(formatted)
		return nil
	},
}

func init() {
	searchCmd.Flags().String("query", "", "Search query text")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "Maximum number of results (max: 20)")
	rootCmd.AddCommand(searchCmd)
}
