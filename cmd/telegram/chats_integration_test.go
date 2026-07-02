package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTelegramChatGet(t *testing.T) {
	env := setupCommandEnvironment(t)
	chatID := getRequiredEnv(t, "TELEGRAM_TEST_CHAT_ID")

	result, err := runCommand(t, env, "telegram", "chats", "get", "--chat", chatID)
	if err != nil {
		t.Fatal(err)
	}

	output := result.stdout.String()
	for _, expected := range []string{"id\t", "type\t"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected chat output to include %q, got:\n%s", expected, output)
		}
	}
}

func TestTelegramChatGetJSON(t *testing.T) {
	env := setupCommandEnvironment(t)
	chatID := getRequiredEnv(t, "TELEGRAM_TEST_CHAT_ID")

	result, err := runCommand(t, env, "telegram", "chats", "get", "--chat", chatID, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var chat TelegramChat
	if err := json.Unmarshal(result.stdout.Bytes(), &chat); err != nil {
		t.Fatal(err)
	}
	if chat.ID == 0 || chat.Type == "" {
		t.Fatalf("expected chat JSON, got %#v", chat)
	}
}
