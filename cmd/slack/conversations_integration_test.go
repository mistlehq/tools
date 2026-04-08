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

func TestConversationsReplies(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	root, replies := createTestThread(t, sc, channelID)
	t.Cleanup(func() {
		cleanupTestThread(t, sc, channelID, root.TS, replies)
	})

	commandResult, err := runCommandWithInput(t, env, "", "slack", "conversations", "replies", "--channel", channelID, "--ts", root.TS)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	for _, want := range []string{root.Message.Text, replies[0].Message.Text, replies[1].Message.Text} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected replies output to contain %q, got %q", want, output)
		}
	}

	if strings.Count(output, "TS: ") < 3 {
		t.Fatalf("expected replies output to include at least 3 messages, got %q", output)
	}
}

func TestConversationsRepliesJSON(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	root, replies := createTestThread(t, sc, channelID)
	t.Cleanup(func() {
		cleanupTestThread(t, sc, channelID, root.TS, replies)
	})

	commandResult, err := runCommandWithInput(t, env, "", "slack", "conversations", "replies", "--channel", channelID, "--ts", root.TS, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var thread SlackConversationsReplies
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &thread); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !thread.OK {
		t.Fatal("expected ok=true in JSON output")
	}

	if len(thread.Messages) < 3 {
		t.Fatalf("expected at least 3 messages in thread JSON, got %d", len(thread.Messages))
	}
}

func createTestThread(t *testing.T, sc SlackClient, channelID string) (SlackChatPostMessage, []SlackChatPostMessage) {
	t.Helper()

	root, err := sc.PostMessage(SlackChatPostMessageInput{
		Channel: channelID,
		Text:    uniqueTestMessage("conversations replies root"),
	})
	if err != nil {
		t.Fatal(err)
	}

	replies := make([]SlackChatPostMessage, 0, 2)
	for _, text := range []string{
		uniqueTestMessage("conversations replies reply one"),
		uniqueTestMessage("conversations replies reply two"),
	} {
		threadTS := root.TS
		reply, err := sc.PostMessage(SlackChatPostMessageInput{
			Channel:  channelID,
			Text:     text,
			ThreadTS: &threadTS,
		})
		if err != nil {
			t.Fatal(err)
		}

		replies = append(replies, reply)
	}

	return root, replies
}

func cleanupTestThread(t *testing.T, sc SlackClient, channelID string, rootTS string, replies []SlackChatPostMessage) {
	t.Helper()

	for index := len(replies) - 1; index >= 0; index-- {
		if _, err := sc.DeleteMessage(SlackChatDeleteInput{
			Channel: channelID,
			TS:      replies[index].TS,
		}); err != nil {
			t.Errorf("failed to delete reply %s: %v", replies[index].TS, err)
		}
	}

	if _, err := sc.DeleteMessage(SlackChatDeleteInput{
		Channel: channelID,
		TS:      rootTS,
	}); err != nil {
		t.Errorf("failed to delete root %s: %v", rootTS, err)
	}
}
