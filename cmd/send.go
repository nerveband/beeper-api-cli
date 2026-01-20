package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nerveband/beeper-api-cli/internal/output"
)

var sendCmd = &cobra.Command{
	Use:   "send --chat-id <chat-id> --message <text>",
	Short: "Send a message to a chat",
	Long:  `Send a new message to the specified Beeper chat.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chatID, _ := cmd.Flags().GetString("chat-id")
		message, _ := cmd.Flags().GetString("message")

		if chatID == "" {
			return fmt.Errorf("--chat-id is required")
		}
		if message == "" {
			return fmt.Errorf("--message is required")
		}

		client := getAPIClient()

		messageID, err := client.SendMessage(chatID, message)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		// Format output based on format preference
		format := getOutputFormat()
		switch format {
		case output.FormatJSON:
			result := map[string]interface{}{
				"success":    true,
				"message_id": messageID,
				"chat_id":    chatID,
			}
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		case output.FormatMarkdown:
			fmt.Printf("**Message sent successfully**\n\nID: `%s`\n", messageID)
		default: // text
			fmt.Printf("Message sent successfully. ID: %s\n", messageID)
		}

		return nil
	},
}

func init() {
	sendCmd.Flags().String("chat-id", "", "Chat ID to send message to")
	sendCmd.Flags().String("message", "", "Message text to send")
	rootCmd.AddCommand(sendCmd)
}
