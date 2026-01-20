package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/output"
)

var (
	messagesLimit int
)

var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Manage messages",
	Long:  `Retrieve messages from Beeper chats.`,
}

var messagesListCmd = &cobra.Command{
	Use:   "list --chat-id <chat-id>",
	Short: "List messages from a chat",
	RunE: func(cmd *cobra.Command, args []string) error {
		chatID, _ := cmd.Flags().GetString("chat-id")
		if chatID == "" {
			return fmt.Errorf("--chat-id is required")
		}

		client := getAPIClient()

		messages, err := client.ListMessages(chatID, messagesLimit)
		if err != nil {
			return fmt.Errorf("failed to list messages: %w", err)
		}

		formatted := output.FormatMessages(messages, getOutputFormat())
		fmt.Print(formatted)
		return nil
	},
}

func init() {
	messagesListCmd.Flags().String("chat-id", "", "Chat ID to retrieve messages from")
	messagesListCmd.Flags().IntVar(&messagesLimit, "limit", 50, "Maximum number of messages to retrieve")

	messagesCmd.AddCommand(messagesListCmd)
	rootCmd.AddCommand(messagesCmd)
}
