package main

import (
	"strings"
	"testing"
)

func TestTopLevelHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "slack", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expected := []string{
		"slack auth help",
		"slack conversations help",
		"slack chat help",
		"slack reactions help",
		"slack files help",
		"slack emoji help",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected top-level help to mention %q", want)
		}
	}
}

func TestAuthHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "slack", "auth", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "slack auth test") {
		t.Fatal("expected auth help to mention slack auth test")
	}
}

func TestConversationsHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "slack", "conversations", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expected := []string{
		"slack conversations list",
		"slack conversations info --channel <conversation-id>",
		"slack conversations history --channel <conversation-id>",
		"slack conversations replies --channel <conversation-id> --ts <thread-root-ts>",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected conversations help to mention %q", want)
		}
	}
}

func TestChatHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "slack", "chat", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expected := []string{
		"slack chat post-message --channel <conversation-id> --text <text>",
		"slack chat update --channel <conversation-id> --ts <ts> --text <text>",
		"slack chat delete --channel <conversation-id> --ts <ts>",
		"slack chat get-permalink --channel <conversation-id> --message-ts <ts>",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected chat help to mention %q", want)
		}
	}
}

func TestReactionsHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "slack", "reactions", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expected := []string{
		"slack reactions add --channel <conversation-id> --timestamp <ts> --name <emoji-name>",
		"slack reactions remove --channel <conversation-id> --timestamp <ts> --name <emoji-name>",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected reactions help to mention %q", want)
		}
	}
}

func TestFilesHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "slack", "files", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expected := []string{
		"slack files upload --path <path> --channel <conversation-id>",
		"slack files upload --path <path> --channel <conversation-id> --initial-comment <text>",
		"slack files upload --path <path> --channel <conversation-id> --thread-ts <ts>",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected files help to mention %q", want)
		}
	}
}

func TestEmojiHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "slack", "emoji", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expected := []string{
		"slack emoji list",
		"slack emoji list --include-categories",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected emoji help to mention %q", want)
		}
	}
}
