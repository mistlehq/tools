package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestConversationsList(t *testing.T) {
	env := setupCommandEnvironment(t)
	commandResult, err := runCommandWithInput(t, env, "", "slack", "conversations", "list", "--limit", "1")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least header and one row, got %q", output)
	}

	if lines[0] != "ID\tNAME\tIS_PRIVATE\tIS_ARCHIVED\tIS_MEMBER" {
		t.Fatalf("unexpected header row: %q", lines[0])
	}

	if len(strings.Split(lines[1], "\t")) != 5 {
		t.Fatalf("expected data row to have 5 columns, got %q", lines[1])
	}
}

func TestConversationsListJSON(t *testing.T) {
	env := setupCommandEnvironment(t)
	commandResult, err := runCommandWithInput(t, env, "", "slack", "conversations", "list", "--limit", "1", "--json")
	if err != nil {
		t.Fatal(err)
	}

	var list SlackConversationsList
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &list); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !list.OK {
		t.Fatal("expected ok=true in JSON output")
	}
}

func TestConversationsInfo(t *testing.T) {
	env := setupCommandEnvironment(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")
	commandResult, err := runCommandWithInput(t, env, "", "slack", "conversations", "info", "--channel", channelID)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")
	expectedPrefixes := []string{
		"ID: ",
		"Name: ",
		"Is Private: ",
		"Is Archived: ",
		"Is Member: ",
	}

	if len(lines) != len(expectedPrefixes) {
		t.Fatalf("expected %d lines, got %d: %q", len(expectedPrefixes), len(lines), output)
	}

	for index, want := range expectedPrefixes {
		if !strings.HasPrefix(lines[index], want) {
			t.Fatalf("expected line %d to start with %q, got %q", index+1, want, lines[index])
		}
	}
}

func TestConversationsInfoJSON(t *testing.T) {
	env := setupCommandEnvironment(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")
	commandResult, err := runCommandWithInput(t, env, "", "slack", "conversations", "info", "--channel", channelID, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var info SlackConversationsInfo
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &info); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !info.OK {
		t.Fatal("expected ok=true in JSON output")
	}
}

func TestConversationsHistoryPublicChannel(t *testing.T) {
	env := setupCommandEnvironment(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")
	commandResult, err := runCommandWithInput(t, env, "", "slack", "conversations", "history", "--channel", channelID, "--limit", "1")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expectedPrefixes := []string{
		"TS: ",
		"Thread TS: ",
		"User: ",
		"Type: ",
		"Text:",
	}

	lines := strings.Split(output, "\n")
	if len(lines) < len(expectedPrefixes) {
		t.Fatalf("expected at least %d lines, got %d: %q", len(expectedPrefixes), len(lines), output)
	}

	for index, want := range expectedPrefixes {
		if !strings.HasPrefix(lines[index], want) {
			t.Fatalf("expected line %d to start with %q, got %q", index+1, want, lines[index])
		}
	}
}

func TestConversationsHistoryPrivateChannel(t *testing.T) {
	env := setupCommandEnvironment(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_PRIVATE_CHANNEL_ID")
	commandResult, err := runCommandWithInput(t, env, "", "slack", "conversations", "history", "--channel", channelID, "--limit", "1")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expectedPrefixes := []string{
		"TS: ",
		"Thread TS: ",
		"User: ",
		"Type: ",
		"Text:",
	}

	lines := strings.Split(output, "\n")
	if len(lines) < len(expectedPrefixes) {
		t.Fatalf("expected at least %d lines, got %d: %q", len(expectedPrefixes), len(lines), output)
	}

	for index, want := range expectedPrefixes {
		if !strings.HasPrefix(lines[index], want) {
			t.Fatalf("expected line %d to start with %q, got %q", index+1, want, lines[index])
		}
	}
}

func TestConversationsHistoryJSON(t *testing.T) {
	env := setupCommandEnvironment(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")
	commandResult, err := runCommandWithInput(t, env, "", "slack", "conversations", "history", "--channel", channelID, "--limit", "1", "--json")
	if err != nil {
		t.Fatal(err)
	}

	var history SlackConversationsHistory
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &history); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !history.OK {
		t.Fatal("expected ok=true in JSON output")
	}
}
