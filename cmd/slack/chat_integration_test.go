package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestChatPostMessage(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	commandResult, err := runCommandWithInput(t, env, "", "slack", "chat", "post-message", "--channel", channelID, "--text", uniqueTestMessage("chat post-message"))
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	ts := parseLineValue(t, output, "TS: ")
	t.Cleanup(func() {
		_, _ = sc.DeleteMessage(SlackChatDeleteInput{
			Channel: channelID,
			TS:      ts,
		})
	})

	expectedPrefixes := []string{
		"Channel: ",
		"TS: ",
		"Thread TS: ",
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

func TestChatPostMessageJSON(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	commandResult, err := runCommandWithInput(t, env, "", "slack", "chat", "post-message", "--channel", channelID, "--text", uniqueTestMessage("chat post-message json"), "--json")
	if err != nil {
		t.Fatal(err)
	}

	var posted SlackChatPostMessage
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &posted); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	t.Cleanup(func() {
		if posted.TS != "" {
			_, _ = sc.DeleteMessage(SlackChatDeleteInput{
				Channel: channelID,
				TS:      posted.TS,
			})
		}
	})

	if !posted.OK {
		t.Fatal("expected ok=true in JSON output")
	}
}

func TestChatUpdate(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	posted, err := sc.PostMessage(SlackChatPostMessageInput{
		Channel: channelID,
		Text:    uniqueTestMessage("chat update original"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, _ = sc.DeleteMessage(SlackChatDeleteInput{
			Channel: channelID,
			TS:      posted.TS,
		})
	})

	updatedText := uniqueTestMessage("chat update replacement")
	commandResult, err := runCommandWithInput(t, env, "", "slack", "chat", "update", "--channel", channelID, "--ts", posted.TS, "--text", updatedText)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if parseLineValue(t, output, "TS: ") != posted.TS {
		t.Fatalf("expected update output to reference ts %q, got %q", posted.TS, output)
	}

	if !strings.Contains(output, updatedText) {
		t.Fatalf("expected updated text in output, got %q", output)
	}
}

func TestChatDelete(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	posted, err := sc.PostMessage(SlackChatPostMessageInput{
		Channel: channelID,
		Text:    uniqueTestMessage("chat delete"),
	})
	if err != nil {
		t.Fatal(err)
	}

	commandResult, err := runCommandWithInput(t, env, "", "slack", "chat", "delete", "--channel", channelID, "--ts", posted.TS)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if parseLineValue(t, output, "TS: ") != posted.TS {
		t.Fatalf("expected delete output to reference ts %q, got %q", posted.TS, output)
	}

	if !strings.Contains(output, "Deleted: true") {
		t.Fatalf("expected delete output to include Deleted: true, got %q", output)
	}
}

func TestChatGetPermalink(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	posted, err := sc.PostMessage(SlackChatPostMessageInput{
		Channel: channelID,
		Text:    uniqueTestMessage("chat permalink"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, _ = sc.DeleteMessage(SlackChatDeleteInput{
			Channel: channelID,
			TS:      posted.TS,
		})
	})

	commandResult, err := runCommandWithInput(t, env, "", "slack", "chat", "get-permalink", "--channel", channelID, "--message-ts", posted.TS)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if parseLineValue(t, output, "Message TS: ") != posted.TS {
		t.Fatalf("expected permalink output to reference ts %q, got %q", posted.TS, output)
	}

	if !strings.HasPrefix(parseLineValue(t, output, "Permalink: "), "https://") {
		t.Fatalf("expected permalink output to include https permalink, got %q", output)
	}
}
