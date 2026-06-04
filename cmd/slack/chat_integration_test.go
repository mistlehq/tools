package main

import (
	"encoding/json"
	"github.com/mistlehq/tools/internal/argparse"
	"strings"
	"testing"
)

func TestSlackChatMessageInputMarshalBlocksAndAttachments(t *testing.T) {
	blocks := json.RawMessage(`[{"type":"section","text":{"type":"mrkdwn","text":"*Deploy:* complete"},"accessory":{"type":"button","text":{"type":"plain_text","text":"Open"},"value":"deploy-1"}}]`)
	attachments := json.RawMessage(`[{"color":"#36a64f","blocks":[{"type":"section","text":{"type":"mrkdwn","text":"Attachment block"}}]}]`)

	body, err := json.Marshal(SlackChatPostMessageInput{
		Channel:     "C123",
		Text:        "deploy status",
		Blocks:      &blocks,
		Attachments: &attachments,
	})
	if err != nil {
		t.Fatal(err)
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatal(err)
	}

	expectedKeys := []string{"channel", "text", "blocks", "attachments"}
	for _, key := range expectedKeys {
		if _, ok := payload[key]; !ok {
			t.Fatalf("expected marshaled payload to include %q, got %s", key, string(body))
		}
	}
	if string(payload["blocks"]) != string(blocks) {
		t.Fatalf("expected blocks to be preserved, got %s", string(payload["blocks"]))
	}
	if string(payload["attachments"]) != string(attachments) {
		t.Fatalf("expected attachments to be preserved, got %s", string(payload["attachments"]))
	}
}

func TestSlackChatUpdateInputMarshalEmptyBlocksAndAttachments(t *testing.T) {
	blocks := json.RawMessage(`[]`)
	attachments := json.RawMessage(`[]`)

	body, err := json.Marshal(SlackChatUpdateInput{
		Channel:     "C123",
		TS:          "123.456",
		Blocks:      &blocks,
		Attachments: &attachments,
	})
	if err != nil {
		t.Fatal(err)
	}

	output := string(body)
	expected := []string{`"blocks":[]`, `"attachments":[]`}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected marshaled payload to include %s, got %s", want, output)
		}
	}
	if strings.Contains(output, `"text"`) {
		t.Fatalf("expected empty text to be omitted, got %s", output)
	}
}

func TestReadOptionalJSONArray(t *testing.T) {
	parsedArgs, err := argparse.Parse([]string{"--blocks", `[{"type":"section","elements":[{"type":"button","text":{"type":"plain_text","text":"Open"}}]}]`}, map[string]argparse.Spec{
		"blocks":      {TakesValue: true},
		"blocks-file": {TakesValue: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	blocks, err := readOptionalJSONArray(strings.NewReader(""), parsedArgs, "blocks", "blocks-file")
	if err != nil {
		t.Fatal(err)
	}
	if blocks == nil {
		t.Fatal("expected blocks JSON to be read")
	}
	if !strings.Contains(string(*blocks), `"elements"`) {
		t.Fatalf("expected block elements to be preserved, got %s", string(*blocks))
	}
}

func TestReadOptionalJSONArrayRejectsObject(t *testing.T) {
	parsedArgs, err := argparse.Parse([]string{"--blocks", `{"type":"section"}`}, map[string]argparse.Spec{
		"blocks":      {TakesValue: true},
		"blocks-file": {TakesValue: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = readOptionalJSONArray(strings.NewReader(""), parsedArgs, "blocks", "blocks-file")
	if err == nil {
		t.Fatal("expected object payload to be rejected")
	}
	if !strings.Contains(err.Error(), "must be a JSON array") {
		t.Fatalf("expected JSON array error, got %v", err)
	}
}

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

func TestChatPostMessageBlocksAndAttachmentsJSON(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	blocks := `[{"type":"section","text":{"type":"mrkdwn","text":"*Block Kit integration test*"},"accessory":{"type":"button","text":{"type":"plain_text","text":"Open"},"value":"open"}}]`
	attachments := `[{"fallback":"attachment fallback","color":"#36a64f","text":"attachment text"}]`
	commandResult, err := runCommandWithInput(t, env, "", "slack", "chat", "post-message", "--channel", channelID, "--text", uniqueTestMessage("chat post-message blocks attachments"), "--blocks", blocks, "--attachments", attachments, "--json")
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
	if len(posted.Message.Blocks) == 0 {
		t.Fatalf("expected posted message to include blocks, got %#v", posted.Message)
	}
	if len(posted.Message.Attachments) == 0 {
		t.Fatalf("expected posted message to include attachments, got %#v", posted.Message)
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

func TestChatUpdateBlocksAndAttachmentsJSON(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	posted, err := sc.PostMessage(SlackChatPostMessageInput{
		Channel: channelID,
		Text:    uniqueTestMessage("chat update blocks attachments original"),
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

	blocks := `[{"type":"section","text":{"type":"mrkdwn","text":"*Updated Block Kit integration test*"},"accessory":{"type":"button","text":{"type":"plain_text","text":"Review"},"value":"review"}}]`
	attachments := `[{"fallback":"updated attachment fallback","color":"#439FE0","text":"updated attachment text"}]`
	commandResult, err := runCommandWithInput(t, env, "", "slack", "chat", "update", "--channel", channelID, "--ts", posted.TS, "--text", uniqueTestMessage("chat update blocks attachments fallback"), "--blocks", blocks, "--attachments", attachments, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var updated SlackChatUpdate
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &updated); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !updated.OK {
		t.Fatal("expected ok=true in JSON output")
	}
	if updated.TS != posted.TS {
		t.Fatalf("expected update ts %q, got %q", posted.TS, updated.TS)
	}
	if len(updated.Message.Blocks) == 0 {
		t.Fatalf("expected updated message to include blocks, got %#v", updated.Message)
	}
	if len(updated.Message.Attachments) == 0 {
		t.Fatalf("expected updated message to include attachments, got %#v", updated.Message)
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
