package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "slack", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"slack mcp", "slack mcp serve", "Streamable HTTP"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp help to mention %q", want)
		}
	}
}

func TestMCPServeHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "slack", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"--addr <addr>", "--endpoint <path>"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp serve help to mention %q", want)
		}
	}
}

func TestMCPServerListsSlackTools(t *testing.T) {
	session := newLocalSlackMCPTestSession(t)
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
		"slack_auth_test":             slackAuthTestDoc.Description,
		"slack_conversations_list":    slackConversationsListDoc.Description,
		"slack_conversations_info":    slackConversationsInfoDoc.Description,
		"slack_conversations_history": slackConversationsHistoryDoc.Description,
		"slack_conversations_replies": slackConversationsRepliesDoc.Description,
		"slack_chat_post_message":     slackChatPostMessageDoc.Description,
		"slack_chat_update":           slackChatUpdateDoc.Description,
		"slack_chat_delete":           slackChatDeleteDoc.Description,
		"slack_chat_get_permalink":    slackChatGetPermalinkDoc.Description,
		"slack_reactions_add":         slackReactionsAddDoc.Description,
		"slack_reactions_remove":      slackReactionsRemoveDoc.Description,
		"slack_files_info":            slackFilesInfoDoc.Description,
		"slack_files_download":        slackFilesDownloadDoc.Description,
		"slack_files_upload":          slackFilesUploadDoc.Description,
		"slack_emoji_list":            slackEmojiListDoc.Description,
	}

	for name, description := range expected {
		tool, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected MCP tool %q to be listed", name)
		}
		if tool.Description != description {
			t.Fatalf("expected MCP tool %q description %q, got %q", name, description, tool.Description)
		}

		if name == "slack_chat_delete" || name == "slack_reactions_remove" {
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || !*tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be annotated as destructive", name)
			}
			continue
		}
		if strings.Contains(name, "post") || strings.Contains(name, "update") || name == "slack_reactions_add" || strings.Contains(name, "files_upload") || strings.Contains(name, "files_download") {
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || *tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be annotated as non-destructive mutation", name)
			}
			continue
		}
		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Fatalf("expected MCP tool %q to be annotated as read-only", name)
		}
	}
}

func TestMCPSlackTools(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")
	session := newSlackMCPTestSession(t, sc)
	defer session.Close()

	authResult := callSlackMCPTool(t, session, "slack_auth_test", map[string]any{})
	var auth SlackAuthTest
	decodeMCPStructuredContent(t, authResult, &auth)
	if !auth.OK {
		t.Fatal("expected slack_auth_test ok=true")
	}

	listResult := callSlackMCPTool(t, session, "slack_conversations_list", map[string]any{
		"limit": "5",
		"types": "public_channel,private_channel",
	})
	var list SlackConversationsList
	decodeMCPStructuredContent(t, listResult, &list)
	if !list.OK {
		t.Fatal("expected slack_conversations_list ok=true")
	}

	infoResult := callSlackMCPTool(t, session, "slack_conversations_info", map[string]any{
		"channel":       channelID,
		"includeLocale": true,
	})
	var info SlackConversationsInfo
	decodeMCPStructuredContent(t, infoResult, &info)
	if !info.OK || info.Channel.ID != channelID {
		t.Fatalf("expected info for channel %q, got %#v", channelID, info.Channel)
	}

	text := uniqueTestMessage("mcp slack post")
	postResult := callSlackMCPTool(t, session, "slack_chat_post_message", map[string]any{
		"channel": channelID,
		"text":    text,
	})
	var posted SlackChatPostMessage
	decodeMCPStructuredContent(t, postResult, &posted)
	if !posted.OK || posted.TS == "" {
		t.Fatalf("expected posted message, got %#v", posted)
	}
	messageDeleted := false
	t.Cleanup(func() {
		if !messageDeleted && posted.TS != "" {
			_, _ = sc.DeleteMessage(SlackChatDeleteInput{Channel: channelID, TS: posted.TS})
		}
	})

	historyResult := callSlackMCPTool(t, session, "slack_conversations_history", map[string]any{
		"channel": channelID,
		"limit":   "5",
	})
	var history SlackConversationsHistory
	decodeMCPStructuredContent(t, historyResult, &history)
	if !history.OK {
		t.Fatal("expected slack_conversations_history ok=true")
	}

	repliesResult := callSlackMCPTool(t, session, "slack_conversations_replies", map[string]any{
		"channel": channelID,
		"ts":      posted.TS,
		"limit":   "5",
	})
	var replies SlackConversationsReplies
	decodeMCPStructuredContent(t, repliesResult, &replies)
	if !replies.OK {
		t.Fatal("expected slack_conversations_replies ok=true")
	}

	updatedText := uniqueTestMessage("mcp slack update")
	updateResult := callSlackMCPTool(t, session, "slack_chat_update", map[string]any{
		"channel": channelID,
		"ts":      posted.TS,
		"text":    updatedText,
	})
	var updated SlackChatUpdate
	decodeMCPStructuredContent(t, updateResult, &updated)
	if !updated.OK || updated.TS != posted.TS {
		t.Fatalf("expected updated message ts %q, got %#v", posted.TS, updated)
	}

	permalinkResult := callSlackMCPTool(t, session, "slack_chat_get_permalink", map[string]any{
		"channel":   channelID,
		"messageTs": posted.TS,
	})
	var permalink SlackChatGetPermalink
	decodeMCPStructuredContent(t, permalinkResult, &permalink)
	if !permalink.OK || !strings.HasPrefix(permalink.Permalink, "https://") {
		t.Fatalf("expected https permalink, got %#v", permalink)
	}

	callSlackMCPTool(t, session, "slack_reactions_add", map[string]any{
		"channel":   channelID,
		"timestamp": posted.TS,
		"name":      "eyes",
	})
	callSlackMCPTool(t, session, "slack_reactions_remove", map[string]any{
		"channel":   channelID,
		"timestamp": posted.TS,
		"name":      "eyes",
	})

	emojiResult := callSlackMCPTool(t, session, "slack_emoji_list", map[string]any{
		"includeCategories": true,
	})
	var emoji SlackEmojiList
	decodeMCPStructuredContent(t, emojiResult, &emoji)
	if !emoji.OK {
		t.Fatal("expected slack_emoji_list ok=true")
	}

	tempDir := t.TempDir()
	uploadPath := filepath.Join(tempDir, "mcp-upload.txt")
	if err := os.WriteFile(uploadPath, []byte("slack mcp upload integration test"), 0o600); err != nil {
		t.Fatal(err)
	}
	uploadResult := callSlackMCPTool(t, session, "slack_files_upload", map[string]any{
		"path":           uploadPath,
		"channel":        channelID,
		"initialComment": "uploaded from MCP test",
		"threadTs":       posted.TS,
	})
	var uploaded SlackFilesCompleteUploadExternal
	decodeMCPStructuredContent(t, uploadResult, &uploaded)
	file := firstUploadedFile(uploaded)
	if file.ID != "" {
		t.Cleanup(func() {
			_ = sc.DeleteFile(file.ID)
		})
	}
	if !uploaded.OK || file.ID == "" {
		t.Fatalf("expected uploaded file, got %#v", uploaded)
	}

	waitForSlackFileMessages(t, file.ID, func() ([]SlackMessage, string) {
		repliesResult := callSlackMCPTool(t, session, "slack_conversations_replies", map[string]any{
			"channel": channelID,
			"ts":      posted.TS,
			"limit":   "10",
		})
		var replies SlackConversationsReplies
		decodeMCPStructuredContent(t, repliesResult, &replies)
		return replies.Messages, fmt.Sprintf("%#v", replies.Messages)
	})

	historyUploadPath := filepath.Join(tempDir, "mcp-history-upload.txt")
	if err := os.WriteFile(historyUploadPath, []byte("slack mcp history upload integration test"), 0o600); err != nil {
		t.Fatal(err)
	}
	historyUploadResult := callSlackMCPTool(t, session, "slack_files_upload", map[string]any{
		"path":    historyUploadPath,
		"channel": channelID,
	})
	var historyUploaded SlackFilesCompleteUploadExternal
	decodeMCPStructuredContent(t, historyUploadResult, &historyUploaded)
	historyFile := firstUploadedFile(historyUploaded)
	if historyFile.ID != "" {
		t.Cleanup(func() {
			_ = sc.DeleteFile(historyFile.ID)
		})
	}
	if !historyUploaded.OK || historyFile.ID == "" {
		t.Fatalf("expected history uploaded file, got %#v", historyUploaded)
	}

	waitForSlackFileMessages(t, historyFile.ID, func() ([]SlackMessage, string) {
		historyResult := callSlackMCPTool(t, session, "slack_conversations_history", map[string]any{
			"channel": channelID,
			"limit":   "10",
		})
		var history SlackConversationsHistory
		decodeMCPStructuredContent(t, historyResult, &history)
		return history.Messages, fmt.Sprintf("%#v", history.Messages)
	})

	infoFileResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "slack_files_info",
		Arguments: map[string]any{
			"fileId": file.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if infoFileResult.IsError {
		if toolErrorContains(infoFileResult, "missing_scope") {
			t.Skip("skipping file info/download checks: Slack test bot token does not have files:read")
		}
		t.Fatalf("expected slack_files_info to succeed, got tool error: %#v", infoFileResult.Content)
	}

	var fileInfo SlackFilesInfo
	decodeMCPStructuredContent(t, infoFileResult, &fileInfo)
	if !fileInfo.OK || fileInfo.File.ID != file.ID {
		t.Fatalf("expected file info for %q, got %#v", file.ID, fileInfo)
	}

	downloadPath := filepath.Join(tempDir, "downloaded.txt")
	downloadResult := callSlackMCPTool(t, session, "slack_files_download", map[string]any{
		"fileId": file.ID,
		"output": downloadPath,
	})
	var downloaded SlackFilesDownload
	decodeMCPStructuredContent(t, downloadResult, &downloaded)
	if downloaded.Bytes == 0 {
		t.Fatal("expected downloaded bytes")
	}

	deleteResult := callSlackMCPTool(t, session, "slack_chat_delete", map[string]any{
		"channel": channelID,
		"ts":      posted.TS,
	})
	var deleted SlackChatDelete
	decodeMCPStructuredContent(t, deleteResult, &deleted)
	if !deleted.OK {
		t.Fatal("expected delete ok=true")
	}
	messageDeleted = true

	_ = env
}

func TestMCPSlackToolValidation(t *testing.T) {
	session := newLocalSlackMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{name: "conversation info missing channel", tool: "slack_conversations_info", arguments: map[string]any{}},
		{name: "history missing channel", tool: "slack_conversations_history", arguments: map[string]any{}},
		{name: "replies missing ts", tool: "slack_conversations_replies", arguments: map[string]any{"channel": "C123"}},
		{name: "post missing text", tool: "slack_chat_post_message", arguments: map[string]any{"channel": "C123"}},
		{name: "update missing ts", tool: "slack_chat_update", arguments: map[string]any{"channel": "C123", "text": "hello"}},
		{name: "delete missing channel", tool: "slack_chat_delete", arguments: map[string]any{"ts": "123.456"}},
		{name: "permalink missing message ts", tool: "slack_chat_get_permalink", arguments: map[string]any{"channel": "C123"}},
		{name: "reaction missing name", tool: "slack_reactions_add", arguments: map[string]any{"channel": "C123", "timestamp": "123.456"}},
		{name: "file info missing file", tool: "slack_files_info", arguments: map[string]any{}},
		{name: "file download missing output", tool: "slack_files_download", arguments: map[string]any{"fileId": "F123"}},
		{name: "file upload missing channel", tool: "slack_files_upload", arguments: map[string]any{"path": "file.txt"}},
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
		})
	}
}

func newLocalSlackMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	return newSlackMCPTestSession(t, NewSlackClient(Config{BaseURL: "http://127.0.0.1"}))
}

func newSlackMCPTestSession(t *testing.T, sc SlackClient) *mcp.ClientSession {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newSlackMCPHTTPHandler(sc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return connectSlackMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func connectSlackMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()

	client := mcp.NewClient(&mcp.Implementation{Name: "slack-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func callSlackMCPTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
	t.Helper()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("expected %s to succeed, got tool error: %#v", name, result.Content)
	}
	return result
}

func decodeMCPStructuredContent(t *testing.T, result *mcp.CallToolResult, output any) {
	t.Helper()

	rawOutput, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(rawOutput, output); err != nil {
		t.Fatal(err)
	}
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
