package output

import (
	"strings"
	"testing"

	"github.com/nerveband/beeper-api-cli/internal/api"
	"github.com/stretchr/testify/assert"
)

// Test data
var (
	testChats = []api.Chat{
		{
			ID:           "chat1",
			Name:         "Test Chat 1",
			Participants: []string{"Alice", "Bob"},
			UnreadCount:  5,
		},
		{
			ID:           "chat2",
			Name:         "Test Chat 2",
			Participants: []string{"Charlie"},
			UnreadCount:  0,
		},
	}

	testMessages = []api.Message{
		{
			ID:        "msg1",
			Text:      "Hello, world!",
			Sender:    "Alice",
			Timestamp: 1640000000,
		},
		{
			ID:        "msg2",
			Text:      "How are you?",
			Sender:    "Bob",
			Timestamp: 1640000100,
		},
	}
)

// TestFormatChatsJSON tests JSON formatting for chats
func TestFormatChatsJSON(t *testing.T) {
	result := FormatChats(testChats, FormatJSON)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "chat1")
	assert.Contains(t, result, "Test Chat 1")
	assert.Contains(t, result, "Alice")
	assert.Contains(t, result, "Bob")

	// Should be valid JSON
	assert.True(t, strings.HasPrefix(result, "["))
	assert.True(t, strings.HasSuffix(strings.TrimSpace(result), "]"))
}

// TestFormatChatsText tests plain text formatting for chats
func TestFormatChatsText(t *testing.T) {
	result := FormatChats(testChats, FormatText)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "chat1")
	assert.Contains(t, result, "Test Chat 1")
	assert.Contains(t, result, "Unread: 5")
	assert.Contains(t, result, "Participants: Alice, Bob")
}

// TestFormatChatsMarkdown tests Markdown formatting for chats
func TestFormatChatsMarkdown(t *testing.T) {
	result := FormatChats(testChats, FormatMarkdown)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "##")
	assert.Contains(t, result, "chat1")
	assert.Contains(t, result, "Test Chat 1")
	assert.Contains(t, result, "**Participants**")
}

// TestFormatMessagesJSON tests JSON formatting for messages
func TestFormatMessagesJSON(t *testing.T) {
	result := FormatMessages(testMessages, FormatJSON)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "msg1")
	assert.Contains(t, result, "Hello, world!")
	assert.Contains(t, result, "Alice")

	// Should be valid JSON array
	assert.True(t, strings.HasPrefix(result, "["))
	assert.True(t, strings.HasSuffix(strings.TrimSpace(result), "]"))
}

// TestFormatMessagesText tests plain text formatting for messages
func TestFormatMessagesText(t *testing.T) {
	result := FormatMessages(testMessages, FormatText)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Alice:")
	assert.Contains(t, result, "Hello, world!")
	assert.Contains(t, result, "Bob:")
	assert.Contains(t, result, "How are you?")
}

// TestFormatMessagesMarkdown tests Markdown formatting for messages
func TestFormatMessagesMarkdown(t *testing.T) {
	result := FormatMessages(testMessages, FormatMarkdown)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "**Alice**")
	assert.Contains(t, result, "Hello, world!")
	assert.Contains(t, result, ">")
}

// TestFormatEmptyChats tests formatting with empty chat list
func TestFormatEmptyChats(t *testing.T) {
	emptyChats := []api.Chat{}

	jsonResult := FormatChats(emptyChats, FormatJSON)
	assert.Equal(t, "[]\n", jsonResult)

	textResult := FormatChats(emptyChats, FormatText)
	assert.Contains(t, textResult, "No chats found")

	mdResult := FormatChats(emptyChats, FormatMarkdown)
	assert.Contains(t, mdResult, "No chats found")
}

// TestFormatEmptyMessages tests formatting with empty message list
func TestFormatEmptyMessages(t *testing.T) {
	emptyMessages := []api.Message{}

	jsonResult := FormatMessages(emptyMessages, FormatJSON)
	assert.Equal(t, "[]\n", jsonResult)

	textResult := FormatMessages(emptyMessages, FormatText)
	assert.Contains(t, textResult, "No messages found")

	mdResult := FormatMessages(emptyMessages, FormatMarkdown)
	assert.Contains(t, mdResult, "No messages found")
}

// TestFormatInvalidFormat tests handling of invalid format
func TestFormatInvalidFormat(t *testing.T) {
	result := FormatChats(testChats, Format("invalid"))
	// Should default to JSON
	assert.Contains(t, result, "[")
}

// TestFormatChatName tests chat name formatting edge cases
func TestFormatChatName(t *testing.T) {
	testCases := []struct {
		name     string
		chat     api.Chat
		expected string
	}{
		{
			name: "Chat with name",
			chat: api.Chat{
				ID:   "chat1",
				Name: "My Chat",
			},
			expected: "My Chat",
		},
		{
			name: "Chat without name",
			chat: api.Chat{
				ID:           "chat2",
				Name:         "",
				Participants: []string{"Alice", "Bob"},
			},
			expected: "Alice, Bob",
		},
		{
			name: "Chat with no name or participants",
			chat: api.Chat{
				ID:           "chat3",
				Name:         "",
				Participants: []string{},
			},
			expected: "chat3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatChats([]api.Chat{tc.chat}, FormatText)
			assert.Contains(t, result, tc.expected)
		})
	}
}

// TestFormatMessageTimestamp tests timestamp formatting
func TestFormatMessageTimestamp(t *testing.T) {
	msg := api.Message{
		ID:        "msg1",
		Text:      "Test",
		Sender:    "Alice",
		Timestamp: 1640000000,
	}

	result := FormatMessages([]api.Message{msg}, FormatText)
	// Should contain a formatted timestamp
	assert.NotEmpty(t, result)
}

// TestFormatLongMessage tests formatting of long messages
func TestFormatLongMessage(t *testing.T) {
	longText := strings.Repeat("This is a long message. ", 50)
	msg := api.Message{
		ID:        "msg1",
		Text:      longText,
		Sender:    "Alice",
		Timestamp: 1640000000,
	}

	result := FormatMessages([]api.Message{msg}, FormatText)
	assert.Contains(t, result, longText)
}

// TestFormatSpecialCharacters tests handling of special characters
func TestFormatSpecialCharacters(t *testing.T) {
	msg := api.Message{
		ID:        "msg1",
		Text:      "Special chars: < > & \" ' \n\t",
		Sender:    "Alice",
		Timestamp: 1640000000,
	}

	jsonResult := FormatMessages([]api.Message{msg}, FormatJSON)
	// JSON should escape special characters
	assert.Contains(t, jsonResult, "\\n")

	textResult := FormatMessages([]api.Message{msg}, FormatText)
	assert.Contains(t, textResult, "Special chars")
}
