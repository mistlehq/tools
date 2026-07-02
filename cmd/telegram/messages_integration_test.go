package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func TestTelegramMessageLifecycle(t *testing.T) {
	env := setupCommandEnvironment(t)
	chatID := getRequiredEnv(t, "TELEGRAM_TEST_CHAT_ID")
	messageText := uniqueTestMessage("telegram integration")
	editedText := messageText + " edited"

	result, err := runCommand(t, env, "telegram", "messages", "send", "--chat", chatID, "--text", messageText)
	if err != nil {
		t.Fatal(err)
	}

	output := result.stdout.String()
	for _, expected := range []string{"id\t", "chat\t", messageText} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected send output to include %q, got:\n%s", expected, output)
		}
	}
	messageID := parseLineValue(t, output, "id\t")

	edited, err := runCommand(t, env, "telegram", "messages", "edit", "--chat", chatID, "--message", messageID, "--text", editedText)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(edited.stdout.String(), editedText) {
		t.Fatalf("expected edit output to include %q, got:\n%s", editedText, edited.stdout.String())
	}

	reaction, err := runCommand(t, env, "telegram", "reactions", "set", "--chat", chatID, "--message", messageID, "--emoji", "👍")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reaction.stdout.String(), "ok\ttrue") {
		t.Fatalf("expected reaction set output to report success, got:\n%s", reaction.stdout.String())
	}

	clearReaction, err := runCommand(t, env, "telegram", "reactions", "clear", "--chat", chatID, "--message", messageID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(clearReaction.stdout.String(), "ok\ttrue") {
		t.Fatalf("expected reaction clear output to report success, got:\n%s", clearReaction.stdout.String())
	}

	deleted, err := runCommand(t, env, "telegram", "messages", "delete", "--chat", chatID, "--message", messageID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(deleted.stdout.String(), "ok\ttrue") {
		t.Fatalf("expected delete output to report success, got:\n%s", deleted.stdout.String())
	}
}

func TestTelegramMessageLifecycleJSON(t *testing.T) {
	env := setupCommandEnvironment(t)
	chatID := getRequiredEnv(t, "TELEGRAM_TEST_CHAT_ID")
	messageText := uniqueTestMessage("telegram json")
	editedText := messageText + " edited"

	result, err := runCommand(t, env, "telegram", "messages", "send", "--chat", chatID, "--text", messageText, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var message TelegramMessage
	if err := json.Unmarshal(result.stdout.Bytes(), &message); err != nil {
		t.Fatal(err)
	}
	if message.MessageID == 0 || message.Chat.ID == 0 || message.Text != messageText {
		t.Fatalf("expected sent message JSON for %q, got %#v", messageText, message)
	}

	edited, err := runCommand(t, env, "telegram", "messages", "edit", "--chat", chatID, "--message", strconv.Itoa(message.MessageID), "--text", editedText, "--json")
	if err != nil {
		t.Fatal(err)
	}
	var editedMessage TelegramMessage
	if err := json.Unmarshal(edited.stdout.Bytes(), &editedMessage); err != nil {
		t.Fatal(err)
	}
	if editedMessage.MessageID != message.MessageID || editedMessage.Text != editedText {
		t.Fatalf("expected edited message JSON for %q, got %#v", editedText, editedMessage)
	}

	reaction, err := runCommand(t, env, "telegram", "reactions", "set", "--chat", chatID, "--message", strconv.Itoa(message.MessageID), "--emoji", "👍", "--json")
	if err != nil {
		t.Fatal(err)
	}
	var reactionResponse TelegramBoolResponse
	if err := json.Unmarshal(reaction.stdout.Bytes(), &reactionResponse); err != nil {
		t.Fatal(err)
	}
	if !reactionResponse.OK {
		t.Fatalf("expected reaction set JSON to report success, got %#v", reactionResponse)
	}

	clearReaction, err := runCommand(t, env, "telegram", "reactions", "clear", "--chat", chatID, "--message", strconv.Itoa(message.MessageID), "--json")
	if err != nil {
		t.Fatal(err)
	}
	var clearReactionResponse TelegramBoolResponse
	if err := json.Unmarshal(clearReaction.stdout.Bytes(), &clearReactionResponse); err != nil {
		t.Fatal(err)
	}
	if !clearReactionResponse.OK {
		t.Fatalf("expected reaction clear JSON to report success, got %#v", clearReactionResponse)
	}

	deleted, err := runCommand(t, env, "telegram", "messages", "delete", "--chat", chatID, "--message", strconv.Itoa(message.MessageID), "--json")
	if err != nil {
		t.Fatal(err)
	}
	var deleteResponse TelegramBoolResponse
	if err := json.Unmarshal(deleted.stdout.Bytes(), &deleteResponse); err != nil {
		t.Fatal(err)
	}
	if !deleteResponse.OK {
		t.Fatalf("expected delete JSON to report success, got %#v", deleteResponse)
	}
}

func TestTelegramMessageDeleteBatch(t *testing.T) {
	env := setupCommandEnvironment(t)
	chatID := getRequiredEnv(t, "TELEGRAM_TEST_CHAT_ID")
	firstID := sendTestMessage(t, env, chatID, uniqueTestMessage("telegram batch first"))
	secondID := sendTestMessage(t, env, chatID, uniqueTestMessage("telegram batch second"))

	result, err := runCommand(t, env, "telegram", "messages", "delete-batch", "--chat", chatID, "--messages", firstID+","+secondID, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var out TelegramBoolResponse
	if err := json.Unmarshal(result.stdout.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if !out.OK {
		t.Fatalf("expected delete batch JSON to report success, got %#v", out)
	}
}

func TestTelegramRequest(t *testing.T) {
	env := setupCommandEnvironment(t)
	chatID := getRequiredEnv(t, "TELEGRAM_TEST_CHAT_ID")
	body, err := json.Marshal(map[string]string{"chat_id": chatID})
	if err != nil {
		t.Fatal(err)
	}

	result, err := runCommand(t, env, "telegram", "request", "--method", "getChat", "--body", string(body), "--json")
	if err != nil {
		t.Fatal(err)
	}

	var chat TelegramChat
	if err := json.Unmarshal(result.stdout.Bytes(), &chat); err != nil {
		t.Fatal(err)
	}
	if chat.ID == 0 || chat.Type == "" {
		t.Fatalf("expected request output to contain chat JSON, got %#v", chat)
	}

	messageID := sendTestMessage(t, env, chatID, uniqueTestMessage("telegram request reaction"))
	numericMessageID, err := strconv.Atoi(messageID)
	if err != nil {
		t.Fatal(err)
	}
	reactionBody, err := json.Marshal(map[string]any{
		"chat_id":    chatID,
		"message_id": numericMessageID,
		"reaction":   []map[string]string{{"type": "emoji", "emoji": "👍"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	reactionResult, err := runCommand(t, env, "telegram", "request", "--method", "setMessageReaction", "--body", string(reactionBody), "--json")
	if err != nil {
		t.Fatal(err)
	}
	var reactionOK bool
	if err := json.Unmarshal(reactionResult.stdout.Bytes(), &reactionOK); err != nil {
		t.Fatal(err)
	}
	if !reactionOK {
		t.Fatal("expected request reaction to report success")
	}
	if _, err := runCommand(t, env, "telegram", "messages", "delete", "--chat", chatID, "--message", messageID); err != nil {
		t.Fatal(err)
	}
}

func TestTelegramForumTopicLifecycle(t *testing.T) {
	env := setupCommandEnvironment(t)
	forumChatID := getRequiredEnv(t, "TELEGRAM_TEST_FORUM_CHAT_ID")
	skipUnlessBotCanManageForumTopics(t, env, forumChatID)
	topicName := uniqueTestMessage("Mistle test topic")

	created, err := runCommand(t, env, "telegram", "topics", "create", "--chat", forumChatID, "--name", topicName, "--json")
	if err != nil {
		t.Fatal(err)
	}
	var topic TelegramForumTopic
	if err := json.Unmarshal(created.stdout.Bytes(), &topic); err != nil {
		t.Fatal(err)
	}
	if topic.MessageThreadID == 0 || topic.Name != topicName {
		t.Fatalf("expected created forum topic %q, got %#v", topicName, topic)
	}

	sent, err := runCommand(t, env, "telegram", "messages", "send", "--chat", forumChatID, "--thread", strconv.Itoa(topic.MessageThreadID), "--text", uniqueTestMessage("telegram threaded message"), "--json")
	if err != nil {
		t.Fatal(err)
	}
	var message TelegramMessage
	if err := json.Unmarshal(sent.stdout.Bytes(), &message); err != nil {
		t.Fatal(err)
	}
	if message.MessageID == 0 {
		t.Fatalf("expected threaded message JSON, got %#v", message)
	}

	deleted, err := runCommand(t, env, "telegram", "topics", "delete", "--chat", forumChatID, "--thread", strconv.Itoa(topic.MessageThreadID), "--json")
	if err != nil {
		t.Fatal(err)
	}
	var deletedTopic TelegramBoolResponse
	if err := json.Unmarshal(deleted.stdout.Bytes(), &deletedTopic); err != nil {
		t.Fatal(err)
	}
	if !deletedTopic.OK {
		t.Fatalf("expected topic delete JSON to report success, got %#v", deletedTopic)
	}
}

func skipUnlessBotCanManageForumTopics(t *testing.T, env Environment, chatID string) {
	t.Helper()

	body, err := json.Marshal(map[string]string{"chat_id": chatID})
	if err != nil {
		t.Fatal(err)
	}
	result, err := runCommand(t, env, "telegram", "request", "--method", "getChat", "--body", string(body), "--json")
	if err != nil {
		t.Fatal(err)
	}
	var chat struct {
		IsForum bool `json:"is_forum"`
	}
	if err := json.Unmarshal(result.stdout.Bytes(), &chat); err != nil {
		t.Fatal(err)
	}
	if !chat.IsForum {
		t.Skipf("skipping forum topic lifecycle: TELEGRAM_TEST_FORUM_CHAT_ID %s is not a forum supergroup with Topics enabled", chatID)
	}

	me, err := runCommand(t, env, "telegram", "request", "--method", "getMe", "--json")
	if err != nil {
		t.Fatal(err)
	}
	var bot TelegramUser
	if err := json.Unmarshal(me.stdout.Bytes(), &bot); err != nil {
		t.Fatal(err)
	}
	memberBody, err := json.Marshal(map[string]any{"chat_id": chatID, "user_id": bot.ID})
	if err != nil {
		t.Fatal(err)
	}
	memberResult, err := runCommand(t, env, "telegram", "request", "--method", "getChatMember", "--body", string(memberBody), "--json")
	if err != nil {
		t.Fatal(err)
	}
	var member struct {
		Status          string `json:"status"`
		CanManageTopics bool   `json:"can_manage_topics"`
	}
	if err := json.Unmarshal(memberResult.stdout.Bytes(), &member); err != nil {
		t.Fatal(err)
	}
	if member.Status != "administrator" || !member.CanManageTopics {
		t.Skipf("skipping forum topic lifecycle: bot status=%s can_manage_topics=%t", member.Status, member.CanManageTopics)
	}
}

func sendTestMessage(t *testing.T, env Environment, chatID string, text string) string {
	t.Helper()

	result, err := runCommand(t, env, "telegram", "messages", "send", "--chat", chatID, "--text", text)
	if err != nil {
		t.Fatal(err)
	}
	return parseLineValue(t, result.stdout.String(), "id\t")
}
