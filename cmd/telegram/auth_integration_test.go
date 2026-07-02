package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTelegramAuthTest(t *testing.T) {
	env := setupCommandEnvironment(t)

	result, err := runCommand(t, env, "telegram", "auth", "test")
	if err != nil {
		t.Fatal(err)
	}

	output := result.stdout.String()
	for _, expected := range []string{"id\t", "bot\ttrue"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected auth output to include %q, got:\n%s", expected, output)
		}
	}
}

func TestTelegramAuthTestJSON(t *testing.T) {
	env := setupCommandEnvironment(t)

	result, err := runCommand(t, env, "telegram", "auth", "test", "--json")
	if err != nil {
		t.Fatal(err)
	}

	var user TelegramUser
	if err := json.Unmarshal(result.stdout.Bytes(), &user); err != nil {
		t.Fatal(err)
	}
	if user.ID == 0 || !user.IsBot {
		t.Fatalf("expected bot user JSON, got %#v", user)
	}
}
