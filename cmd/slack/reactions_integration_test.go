package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestReactionsAddAndRemove(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	posted, err := sc.PostMessage(SlackChatPostMessageInput{
		Channel: channelID,
		Text:    uniqueTestMessage("reactions"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, _ = sc.RemoveReaction(SlackReactionInput{
			Channel:   channelID,
			Timestamp: posted.TS,
			Name:      "eyes",
		})
		_, _ = sc.DeleteMessage(SlackChatDeleteInput{
			Channel: channelID,
			TS:      posted.TS,
		})
	})

	addResult, err := runCommandWithInput(t, env, "", "slack", "reactions", "add", "--channel", channelID, "--timestamp", posted.TS, "--name", "eyes")
	if err != nil {
		t.Fatal(err)
	}

	addOutput := strings.TrimSpace(addResult.stdout.String())
	if !strings.Contains(addOutput, "Action: added") {
		t.Fatalf("expected add output to include Action: added, got %q", addOutput)
	}

	removeResult, err := runCommandWithInput(t, env, "", "slack", "reactions", "remove", "--channel", channelID, "--timestamp", posted.TS, "--name", "eyes")
	if err != nil {
		t.Fatal(err)
	}

	removeOutput := strings.TrimSpace(removeResult.stdout.String())
	if !strings.Contains(removeOutput, "Action: removed") {
		t.Fatalf("expected remove output to include Action: removed, got %q", removeOutput)
	}
}

func TestReactionsAddJSON(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	posted, err := sc.PostMessage(SlackChatPostMessageInput{
		Channel: channelID,
		Text:    uniqueTestMessage("reactions json"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, _ = sc.RemoveReaction(SlackReactionInput{
			Channel:   channelID,
			Timestamp: posted.TS,
			Name:      "eyes",
		})
		_, _ = sc.DeleteMessage(SlackChatDeleteInput{
			Channel: channelID,
			TS:      posted.TS,
		})
	})

	commandResult, err := runCommandWithInput(t, env, "", "slack", "reactions", "add", "--channel", channelID, "--timestamp", posted.TS, "--name", "eyes", "--json")
	if err != nil {
		t.Fatal(err)
	}

	var response SlackReactionResponse
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &response); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !response.OK {
		t.Fatal("expected ok=true in JSON output")
	}
}
