package main

import (
	"strings"
	"testing"
)

func TestLocalCommandsCoverTelegramCommandSurface(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "no args prints root help",
			args:     []string{"telegram"},
			expected: []string{"CLI for Telegram Bot API.", "telegram auth help", "telegram chats help", "telegram messages help", "telegram reactions help", "telegram topics help", "telegram request help", "telegram mcp help"},
		},
		{
			name:     "help",
			args:     []string{"telegram", "help"},
			expected: []string{"CLI for Telegram Bot API.", "telegram auth help", "telegram chats help", "telegram messages help", "telegram reactions help", "telegram topics help", "telegram request help", "telegram mcp help"},
		},
		{
			name:     "--help",
			args:     []string{"telegram", "--help"},
			expected: []string{"CLI for Telegram Bot API.", "telegram auth help", "telegram chats help", "telegram messages help", "telegram reactions help", "telegram topics help", "telegram request help", "telegram mcp help"},
		},
		{
			name:     "version",
			args:     []string{"telegram", "version"},
			expected: []string{Version},
		},
		{
			name:     "--version",
			args:     []string{"telegram", "--version"},
			expected: []string{Version},
		},
		{
			name:     "auth help",
			args:     []string{"telegram", "auth", "help"},
			expected: []string{"telegram auth test [--json]"},
		},
		{
			name:     "chats help",
			args:     []string{"telegram", "chats", "help"},
			expected: []string{"telegram chats get --chat <chat-id-or-username> [--json]"},
		},
		{
			name:     "messages help",
			args:     []string{"telegram", "messages", "help"},
			expected: []string{"telegram messages send --chat <chat-id-or-username> --text <text> [--thread <message-thread-id>] [--parse-mode <mode>] [--json]", "telegram messages edit", "telegram messages delete", "telegram messages delete-batch"},
		},
		{
			name:     "reactions help",
			args:     []string{"telegram", "reactions", "help"},
			expected: []string{"telegram reactions set", "telegram reactions clear", "telegram reactions delete", "telegram reactions delete-all"},
		},
		{
			name:     "request help",
			args:     []string{"telegram", "request", "help"},
			expected: []string{"telegram request --method <telegram-method> [--body <json>] [--json]"},
		},
		{
			name:     "topics help",
			args:     []string{"telegram", "topics", "help"},
			expected: []string{"telegram topics create", "telegram topics delete"},
		},
		{
			name:     "mcp help",
			args:     []string{"telegram", "mcp", "help"},
			expected: []string{"telegram mcp serve", "Streamable HTTP"},
		},
		{
			name:     "mcp serve help",
			args:     []string{"telegram", "mcp", "serve", "--help"},
			expected: []string{"telegram_auth_test", "telegram_chats_get", "telegram_messages_send", "telegram_messages_edit", "telegram_messages_delete", "telegram_messages_delete_batch", "telegram_reactions_set", "telegram_reactions_clear", "telegram_reactions_delete", "telegram_reactions_delete_all", "telegram_topics_create", "telegram_topics_delete", "telegram_request"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := runCommand(t, Environment{}, tc.args...)
			if err != nil {
				t.Fatal(err)
			}
			output := result.stdout.String()
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
				}
			}
		})
	}
}

func TestLeafCommandsAcceptHelpWithoutConfig(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{name: "auth test", args: []string{"telegram", "auth", "test", "--help"}, expected: "telegram auth test [--json]"},
		{name: "chats get", args: []string{"telegram", "chats", "get", "--help"}, expected: "telegram chats get --chat <chat-id-or-username> [--json]"},
		{name: "messages send", args: []string{"telegram", "messages", "send", "--help"}, expected: "telegram messages send --chat <chat-id-or-username> --text <text> [--thread <message-thread-id>] [--parse-mode <mode>] [--json]"},
		{name: "messages edit", args: []string{"telegram", "messages", "edit", "--help"}, expected: "telegram messages edit --chat <chat-id-or-username> --message <message-id> --text <text> [--parse-mode <mode>] [--json]"},
		{name: "messages delete", args: []string{"telegram", "messages", "delete", "--help"}, expected: "telegram messages delete --chat <chat-id-or-username> --message <message-id> [--json]"},
		{name: "messages delete batch", args: []string{"telegram", "messages", "delete-batch", "--help"}, expected: "telegram messages delete-batch --chat <chat-id-or-username> --messages <message-id-csv> [--json]"},
		{name: "reactions set", args: []string{"telegram", "reactions", "set", "--help"}, expected: "telegram reactions set --chat <chat-id-or-username> --message <message-id> --emoji <emoji-csv> [--custom-emoji-id <id-csv>] [--big] [--json]"},
		{name: "reactions clear", args: []string{"telegram", "reactions", "clear", "--help"}, expected: "telegram reactions clear --chat <chat-id-or-username> --message <message-id> [--json]"},
		{name: "reactions delete", args: []string{"telegram", "reactions", "delete", "--help"}, expected: "telegram reactions delete --chat <chat-id-or-username> --message <message-id> [--user <user-id>] [--actor-chat-id <chat-id>] [--json]"},
		{name: "reactions delete all", args: []string{"telegram", "reactions", "delete-all", "--help"}, expected: "telegram reactions delete-all --chat <chat-id-or-username> [--user <user-id>] [--actor-chat-id <chat-id>] [--json]"},
		{name: "topics create", args: []string{"telegram", "topics", "create", "--help"}, expected: "telegram topics create --chat <chat-id-or-username> --name <name> [--icon-color <rgb-int>] [--icon-custom-emoji-id <emoji-id>] [--json]"},
		{name: "topics delete", args: []string{"telegram", "topics", "delete", "--help"}, expected: "telegram topics delete --chat <chat-id-or-username> --thread <message-thread-id> [--json]"},
		{name: "request", args: []string{"telegram", "request", "--help"}, expected: "telegram request --method <telegram-method> [--body <json>] [--json]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := runCommand(t, Environment{}, tc.args...)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(result.stdout.String(), tc.expected) {
				t.Fatalf("expected leaf help to contain %q, got:\n%s", tc.expected, result.stdout.String())
			}
		})
	}
}

func TestUnsupportedCommandsFailClearly(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{name: "root", args: []string{"telegram", "unknown"}, expected: "unsupported command: unknown"},
		{name: "auth", args: []string{"telegram", "auth", "unknown"}, expected: "unsupported auth command: unknown"},
		{name: "chats", args: []string{"telegram", "chats", "unknown"}, expected: "unsupported chats command: unknown"},
		{name: "messages", args: []string{"telegram", "messages", "unknown"}, expected: "unsupported messages command: unknown"},
		{name: "reactions", args: []string{"telegram", "reactions", "unknown"}, expected: "unsupported reactions command: unknown"},
		{name: "topics", args: []string{"telegram", "topics", "unknown"}, expected: "unsupported topics command: unknown"},
		{name: "mcp", args: []string{"telegram", "mcp", "unknown"}, expected: "unsupported mcp command: unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := runCommand(t, Environment{}, tc.args...)
			if err == nil {
				t.Fatal("expected command to fail")
			}
			if err.Error() != tc.expected {
				t.Fatalf("expected error %q, got %q", tc.expected, err.Error())
			}
		})
	}
}
