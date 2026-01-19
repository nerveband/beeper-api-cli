package output

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nerveband/beeper-cli/internal/api"
)

// Format defines the output format type
type Format string

const (
	FormatJSON     Format = "json"
	FormatText     Format = "text"
	FormatMarkdown Format = "markdown"
)

// FormatChats formats a list of chats according to the specified format
func FormatChats(chats []api.Chat, format Format) string {
	if len(chats) == 0 {
		switch format {
		case FormatJSON:
			return "[]\n"
		case FormatText, FormatMarkdown:
			return "No chats found.\n"
		}
	}

	switch format {
	case FormatJSON:
		data, err := formatChatsJSON(chats)
		if err != nil {
			return fmt.Sprintf("Error formatting JSON: %v\n", err)
		}
		return data
	case FormatText:
		return formatChatsText(chats)
	case FormatMarkdown:
		return formatChatsMarkdown(chats)
	default:
		// Default to JSON for unknown formats
		data, _ := formatChatsJSON(chats)
		return data
	}
}

func formatChatsJSON(chats []api.Chat) (string, error) {
	data, err := json.MarshalIndent(chats, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

func formatChatsText(chats []api.Chat) string {
	var sb strings.Builder
	for _, chat := range chats {
		sb.WriteString(fmt.Sprintf("ID: %s\n", chat.ID))
		sb.WriteString(fmt.Sprintf("Name: %s\n", chat.Name))
		sb.WriteString(fmt.Sprintf("Participants: %s\n", strings.Join(chat.Participants, ", ")))
		sb.WriteString(fmt.Sprintf("Unread: %d\n", chat.UnreadCount))
		sb.WriteString(fmt.Sprintf("Last Message: %s\n", chat.LastMessage))
		sb.WriteString(fmt.Sprintf("Updated: %s\n", chat.UpdatedAt.Format(time.RFC3339)))
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatChatsMarkdown(chats []api.Chat) string {
	var sb strings.Builder
	sb.WriteString("# Chats\n\n")
	for _, chat := range chats {
		sb.WriteString(fmt.Sprintf("## %s\n\n", chat.Name))
		sb.WriteString(fmt.Sprintf("- **ID**: %s\n", chat.ID))
		sb.WriteString(fmt.Sprintf("- **Participants**: %s\n", strings.Join(chat.Participants, ", ")))
		sb.WriteString(fmt.Sprintf("- **Unread**: %d\n", chat.UnreadCount))
		sb.WriteString(fmt.Sprintf("- **Last Message**: %s\n", chat.LastMessage))
		sb.WriteString(fmt.Sprintf("- **Updated**: %s\n\n", chat.UpdatedAt.Format(time.RFC3339)))
	}
	return sb.String()
}

// FormatMessages formats a list of messages according to the specified format
func FormatMessages(messages []api.Message, format Format) string {
	if len(messages) == 0 {
		switch format {
		case FormatJSON:
			return "[]\n"
		case FormatText, FormatMarkdown:
			return "No messages found.\n"
		}
	}

	switch format {
	case FormatJSON:
		data, err := formatMessagesJSON(messages)
		if err != nil {
			return fmt.Sprintf("Error formatting JSON: %v\n", err)
		}
		return data
	case FormatText:
		return formatMessagesText(messages)
	case FormatMarkdown:
		return formatMessagesMarkdown(messages)
	default:
		// Default to JSON
		data, _ := formatMessagesJSON(messages)
		return data
	}
}

func formatMessagesJSON(messages []api.Message) (string, error) {
	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

func formatMessagesText(messages []api.Message) string {
	var sb strings.Builder
	for _, msg := range messages {
		timestamp := time.Unix(msg.Timestamp, 0)
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			timestamp.Format("2006-01-02 15:04:05"),
			msg.Sender,
			msg.Text,
		))
	}
	return sb.String()
}

func formatMessagesMarkdown(messages []api.Message) string {
	var sb strings.Builder
	sb.WriteString("# Messages\n\n")
	for _, msg := range messages {
		timestamp := time.Unix(msg.Timestamp, 0)
		sb.WriteString(fmt.Sprintf("**%s** - %s\n\n",
			msg.Sender,
			timestamp.Format("2006-01-02 15:04:05"),
		))
		sb.WriteString(fmt.Sprintf("> %s\n\n", msg.Text))
		sb.WriteString("---\n\n")
	}
	return sb.String()
}

// FormatSendResponse formats a send message response
func FormatSendResponse(resp *api.SendMessageResponse, format Format) (string, error) {
	switch format {
	case FormatJSON:
		data, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return string(data), nil
	case FormatText:
		if resp.Success {
			return fmt.Sprintf("Message sent successfully. ID: %s\n", resp.MessageID), nil
		}
		return "Failed to send message\n", nil
	case FormatMarkdown:
		if resp.Success {
			return fmt.Sprintf("**Message sent successfully**\n\nID: `%s`\n", resp.MessageID), nil
		}
		return "**Failed to send message**\n", nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}
