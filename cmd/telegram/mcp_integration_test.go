package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPServerListsTelegramTools(t *testing.T) {
	session := newLocalTelegramMCPTestSession(t)
	defer session.Close()

	toolsResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	toolsByName := make(map[string]*mcp.Tool)
	for _, tool := range toolsResult.Tools {
		toolsByName[tool.Name] = tool
	}

	expected := map[string]string{
		"telegram_auth_test":             telegramAuthTestDoc.Description,
		"telegram_chats_get":             telegramChatsGetDoc.Description,
		"telegram_messages_send":         telegramMessagesSendDoc.Description,
		"telegram_messages_edit":         telegramMessagesEditDoc.Description,
		"telegram_messages_delete":       telegramMessagesDeleteDoc.Description,
		"telegram_messages_delete_batch": telegramMessagesDeleteBatchDoc.Description,
		"telegram_reactions_set":         telegramReactionsSetDoc.Description,
		"telegram_reactions_clear":       telegramReactionsClearDoc.Description,
		"telegram_reactions_delete":      telegramReactionsDeleteDoc.Description,
		"telegram_reactions_delete_all":  telegramReactionsDeleteAllDoc.Description,
		"telegram_topics_create":         telegramTopicsCreateDoc.Description,
		"telegram_topics_delete":         telegramTopicsDeleteDoc.Description,
		"telegram_request":               telegramRequestDoc.Description,
	}

	if len(toolsByName) != len(expected) {
		t.Fatalf("expected %d MCP tools, got %d", len(expected), len(toolsByName))
	}

	for name, description := range expected {
		tool, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected MCP tool %q to be listed", name)
		}
		if tool.Description != description {
			t.Fatalf("expected MCP tool %q description %q, got %q", name, description, tool.Description)
		}
		if name == "telegram_messages_send" && !strings.Contains(mustMarshalToolSchema(t, tool), "message_thread_id") {
			t.Fatalf("expected MCP tool %q schema to expose message_thread_id", name)
		}

		switch name {
		case "telegram_messages_delete", "telegram_messages_delete_batch", "telegram_reactions_delete", "telegram_reactions_delete_all", "telegram_topics_delete":
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || !*tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be annotated as destructive", name)
			}
		case "telegram_messages_send", "telegram_messages_edit", "telegram_reactions_set", "telegram_reactions_clear", "telegram_topics_create", "telegram_request":
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || *tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be annotated as non-destructive mutation", name)
			}
		default:
			if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
				t.Fatalf("expected MCP tool %q to be annotated as read-only", name)
			}
		}
	}
}

func TestMCPTelegramToolValidation(t *testing.T) {
	session := newLocalTelegramMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name          string
		tool          string
		arguments     map[string]any
		expectedError string
	}{
		{name: "chat get missing chat", tool: "telegram_chats_get", arguments: map[string]any{"chat": ""}, expectedError: "chat is required"},
		{name: "messages send missing text", tool: "telegram_messages_send", arguments: map[string]any{"chat": "123", "text": ""}, expectedError: "text is required"},
		{name: "messages send reports first required field", tool: "telegram_messages_send", arguments: map[string]any{"chat": "", "text": ""}, expectedError: "chat is required"},
		{name: "messages edit missing id", tool: "telegram_messages_edit", arguments: map[string]any{"chat": "123", "message_id": 0, "text": "edited"}, expectedError: "message_id is required"},
		{name: "messages delete batch missing ids", tool: "telegram_messages_delete_batch", arguments: map[string]any{"chat": "123", "message_ids": []int{}}, expectedError: "message_ids is required"},
		{name: "reactions set missing reaction", tool: "telegram_reactions_set", arguments: map[string]any{"chat": "123", "message_id": 1}, expectedError: "--emoji or --custom-emoji-id is required"},
		{name: "topics create missing name", tool: "telegram_topics_create", arguments: map[string]any{"chat": "123", "name": ""}, expectedError: "name is required"},
		{name: "topics delete missing thread", tool: "telegram_topics_delete", arguments: map[string]any{"chat": "123", "message_thread_id": 0}, expectedError: "message_thread_id is required"},
		{name: "request missing method", tool: "telegram_request", arguments: map[string]any{"method": ""}, expectedError: "method is required"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
				Name:      tc.tool,
				Arguments: tc.arguments,
			})
			if err != nil {
				t.Fatal(err)
			}
			if !result.IsError {
				t.Fatal("expected tool validation to return a tool error")
			}
			if !toolErrorContains(result, tc.expectedError) {
				t.Fatalf("expected tool error to contain %q, got %#v", tc.expectedError, result.Content)
			}
		})
	}
}

func TestMCPServerCallsTelegramTools(t *testing.T) {
	session := newLiveTelegramMCPTestSession(t)
	defer session.Close()
	chatID := getRequiredEnv(t, "TELEGRAM_TEST_CHAT_ID")
	messageText := uniqueTestMessage("telegram mcp")

	authResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: "telegram_auth_test"})
	if err != nil {
		t.Fatal(err)
	}
	var user TelegramUser
	decodeStructuredContent(t, authResult, &user)
	if user.ID == 0 || !user.IsBot {
		t.Fatalf("expected bot user structured content, got %#v", user)
	}

	chatResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "telegram_chats_get",
		Arguments: map[string]any{"chat": chatID},
	})
	if err != nil {
		t.Fatal(err)
	}
	var chat TelegramChat
	decodeStructuredContent(t, chatResult, &chat)
	if chat.ID == 0 || chat.Type == "" {
		t.Fatalf("expected chat structured content, got %#v", chat)
	}

	messageResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "telegram_messages_send",
		Arguments: map[string]any{"chat": chatID, "text": messageText},
	})
	if err != nil {
		t.Fatal(err)
	}
	var message TelegramMessage
	decodeStructuredContent(t, messageResult, &message)
	if message.MessageID == 0 || message.Chat.ID == 0 || message.Text != messageText {
		t.Fatalf("expected sent message structured content for %q, got %#v", messageText, message)
	}

	editedText := messageText + " edited"
	editResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "telegram_messages_edit",
		Arguments: map[string]any{"chat": chatID, "message_id": message.MessageID, "text": editedText},
	})
	if err != nil {
		t.Fatal(err)
	}
	var edited TelegramMessage
	decodeStructuredContent(t, editResult, &edited)
	if edited.MessageID != message.MessageID || edited.Text != editedText {
		t.Fatalf("expected edited message structured content for %q, got %#v", editedText, edited)
	}

	reactionResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "telegram_reactions_set",
		Arguments: map[string]any{"chat": chatID, "message_id": message.MessageID, "emoji": "👍"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var reaction TelegramBoolResponse
	decodeStructuredContent(t, reactionResult, &reaction)
	if !reaction.OK {
		t.Fatalf("expected reaction set success, got %#v", reaction)
	}

	clearReactionResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "telegram_reactions_clear",
		Arguments: map[string]any{"chat": chatID, "message_id": message.MessageID},
	})
	if err != nil {
		t.Fatal(err)
	}
	var clearReaction TelegramBoolResponse
	decodeStructuredContent(t, clearReactionResult, &clearReaction)
	if !clearReaction.OK {
		t.Fatalf("expected reaction clear success, got %#v", clearReaction)
	}

	requestResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "telegram_request",
		Arguments: map[string]any{"method": "getChat", "body": map[string]any{"chat_id": chatID}},
	})
	if err != nil {
		t.Fatal(err)
	}
	var requestOutput struct {
		Result TelegramChat `json:"result"`
	}
	decodeStructuredContent(t, requestResult, &requestOutput)
	if requestOutput.Result.ID == 0 || requestOutput.Result.Type == "" {
		t.Fatalf("expected request structured content to contain chat, got %#v", requestOutput)
	}

	deleteResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "telegram_messages_delete",
		Arguments: map[string]any{"chat": chatID, "message_id": message.MessageID},
	})
	if err != nil {
		t.Fatal(err)
	}
	var deleted TelegramBoolResponse
	decodeStructuredContent(t, deleteResult, &deleted)
	if !deleted.OK {
		t.Fatalf("expected delete success, got %#v", deleted)
	}

	firstBatchID := sendMCPTestMessage(t, session, chatID, uniqueTestMessage("telegram mcp batch first"))
	secondBatchID := sendMCPTestMessage(t, session, chatID, uniqueTestMessage("telegram mcp batch second"))
	deleteBatchResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "telegram_messages_delete_batch",
		Arguments: map[string]any{"chat": chatID, "message_ids": []int{firstBatchID, secondBatchID}},
	})
	if err != nil {
		t.Fatal(err)
	}
	var deletedBatch TelegramBoolResponse
	decodeStructuredContent(t, deleteBatchResult, &deletedBatch)
	if !deletedBatch.OK {
		t.Fatalf("expected delete batch success, got %#v", deletedBatch)
	}
}

func TestMCPServeArgParsingCoversServeOptions(t *testing.T) {
	defaultConfig, err := parseTelegramMCPServeArgs(nil)
	if err != nil {
		t.Fatal(err)
	}
	if defaultConfig.Addr != defaultMCPAddr || defaultConfig.Endpoint != defaultMCPEndpoint {
		t.Fatalf("expected default MCP config, got %#v", defaultConfig)
	}

	customConfig, err := parseTelegramMCPServeArgs([]string{"--addr", "127.0.0.1:8080", "--endpoint", "/telegram-mcp"})
	if err != nil {
		t.Fatal(err)
	}
	if customConfig.Addr != "127.0.0.1:8080" || customConfig.Endpoint != "/telegram-mcp" {
		t.Fatalf("expected custom MCP config, got %#v", customConfig)
	}

	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{name: "positional", args: []string{"extra"}, expected: "mcp serve does not accept positional arguments"},
		{name: "empty addr", args: []string{"--addr", " "}, expected: "--addr must not be empty"},
		{name: "relative endpoint", args: []string{"--endpoint", "mcp"}, expected: "--endpoint must start with '/'"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseTelegramMCPServeArgs(tc.args)
			if err == nil {
				t.Fatal("expected parse to fail")
			}
			if err.Error() != tc.expected {
				t.Fatalf("expected error %q, got %q", tc.expected, err.Error())
			}
		})
	}
}

func newLocalTelegramMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newTelegramMCPHTTPHandler(NewTelegramClient(Config{BaseURL: "http://127.0.0.1"})))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "telegram-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: server.URL + defaultMCPEndpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func newLiveTelegramMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()

	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newTelegramMCPHTTPHandler(NewTelegramClient(config)))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "telegram-live-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: server.URL + defaultMCPEndpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func decodeStructuredContent(t *testing.T, result *mcp.CallToolResult, out any) {
	t.Helper()

	if result.IsError {
		t.Fatalf("expected successful tool result, got %#v", result.Content)
	}
	encoded, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(encoded, out); err != nil {
		t.Fatal(err)
	}
}

func sendMCPTestMessage(t *testing.T, session *mcp.ClientSession, chatID string, text string) int {
	t.Helper()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "telegram_messages_send",
		Arguments: map[string]any{"chat": chatID, "text": text},
	})
	if err != nil {
		t.Fatal(err)
	}
	var message TelegramMessage
	decodeStructuredContent(t, result, &message)
	if message.MessageID == 0 {
		t.Fatalf("expected sent message ID, got %#v", message)
	}
	return message.MessageID
}

func toolErrorContains(result *mcp.CallToolResult, text string) bool {
	for _, content := range result.Content {
		textContent, ok := content.(*mcp.TextContent)
		if ok && strings.Contains(textContent.Text, text) {
			return true
		}
	}
	return false
}

func mustMarshalToolSchema(t *testing.T, tool *mcp.Tool) string {
	t.Helper()

	encoded, err := json.Marshal(tool)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}
