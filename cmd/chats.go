package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/api"
	"github.com/nerveband/beeper-api-cli/internal/output"
)

var chatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Manage chats",
	Long:  `List and retrieve information about Beeper chats/conversations.`,
}

var chatsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all chats",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getAPIClient()

		chats, err := client.ListChats()
		if err != nil {
			return fmt.Errorf("failed to list chats: %w", err)
		}

		formatted := output.FormatChats(chats, getOutputFormat())
		fmt.Print(formatted)
		return nil
	},
}

var chatsGetCmd = &cobra.Command{
	Use:   "get <chat-id>",
	Short: "Get details of a specific chat",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getAPIClient()
		chatID := args[0]

		chat, err := client.GetChat(chatID)
		if err != nil {
			return fmt.Errorf("failed to get chat: %w", err)
		}

		// Format as single-item array for consistent output
		chats := []api.Chat{*chat}
		formatted := output.FormatChats(chats, getOutputFormat())
		fmt.Print(formatted)
		return nil
	},
}

func init() {
	chatsCmd.AddCommand(chatsListCmd)
	chatsCmd.AddCommand(chatsGetCmd)
	rootCmd.AddCommand(chatsCmd)
}
