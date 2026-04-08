package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEmojiList(t *testing.T) {
	env := setupCommandEnvironment(t)
	commandResult, err := runCommandWithInput(t, env, "", "slack", "emoji", "list")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least header and one row, got %q", output)
	}

	if lines[0] != "NAME\tVALUE" {
		t.Fatalf("unexpected header row: %q", lines[0])
	}

	if len(strings.Split(lines[1], "\t")) != 2 {
		t.Fatalf("expected data row to have 2 columns, got %q", lines[1])
	}
}

func TestEmojiListJSON(t *testing.T) {
	env := setupCommandEnvironment(t)
	commandResult, err := runCommandWithInput(t, env, "", "slack", "emoji", "list", "--json")
	if err != nil {
		t.Fatal(err)
	}

	var list SlackEmojiList
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &list); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !list.OK {
		t.Fatal("expected ok=true in JSON output")
	}

	if len(list.Emoji) == 0 {
		t.Fatal("expected emoji map to be non-empty")
	}
}
