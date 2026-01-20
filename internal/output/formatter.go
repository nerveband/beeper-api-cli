package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nerveband/beeper-api-cli/internal/api"
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
		sb.WriteString(fmt.Sprintf("Title: %s\n", chat.Title))
		sb.WriteString(fmt.Sprintf("Type: %s\n", chat.Type))
		sb.WriteString(fmt.Sprintf("Network: %s\n", chat.Network))
		sb.WriteString(fmt.Sprintf("Unread: %d\n", chat.UnreadCount))
		if chat.IsMuted {
			sb.WriteString("Muted: Yes\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatChatsMarkdown(chats []api.Chat) string {
	var sb strings.Builder
	sb.WriteString("# Chats\n\n")
	for _, chat := range chats {
		title := chat.Title
		if title == "" {
			title = chat.ID
		}
		sb.WriteString(fmt.Sprintf("## %s\n\n", title))
		sb.WriteString(fmt.Sprintf("- **ID**: %s\n", chat.ID))
		sb.WriteString(fmt.Sprintf("- **Type**: %s\n", chat.Type))
		sb.WriteString(fmt.Sprintf("- **Network**: %s\n", chat.Network))
		sb.WriteString(fmt.Sprintf("- **Unread**: %d\n\n", chat.UnreadCount))
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
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			msg.Timestamp,
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
		sb.WriteString(fmt.Sprintf("**%s** - %s\n\n",
			msg.Sender,
			msg.Timestamp,
		))
		sb.WriteString(fmt.Sprintf("> %s\n\n", msg.Text))
		sb.WriteString("---\n\n")
	}
	return sb.String()
}

