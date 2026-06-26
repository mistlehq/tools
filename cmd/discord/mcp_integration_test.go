package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPServerListsDiscordTools(t *testing.T) {
	session := newLocalDiscordMCPTestSession(t)
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
		"discord_auth_test":           discordAuthTestDoc.Description,
		"discord_guilds_list":         discordGuildsListDoc.Description,
		"discord_guilds_get":          discordGuildsGetDoc.Description,
		"discord_channels_list":       discordChannelsListDoc.Description,
		"discord_channels_get":        discordChannelsGetDoc.Description,
		"discord_messages_list":       discordMessagesListDoc.Description,
		"discord_messages_send":       discordMessagesSendDoc.Description,
		"discord_messages_edit":       discordMessagesEditDoc.Description,
		"discord_messages_delete":     discordMessagesDeleteDoc.Description,
		"discord_reactions_add":       discordReactionsAddDoc.Description,
		"discord_reactions_remove":    discordReactionsRemoveDoc.Description,
		"discord_roles_list":          discordRolesListDoc.Description,
		"discord_roles_create":        discordRolesCreateDoc.Description,
		"discord_roles_delete":        discordRolesDeleteDoc.Description,
		"discord_members_list":        discordMembersListDoc.Description,
		"discord_members_get":         discordMembersGetDoc.Description,
		"discord_members_add_role":    discordMembersAddRoleDoc.Description,
		"discord_members_remove_role": discordMembersRemoveRoleDoc.Description,
		"discord_members_ban":         discordMembersBanDoc.Description,
		"discord_members_unban":       discordMembersUnbanDoc.Description,
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

		switch {
		case isDiscordDestructiveTool(name):
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || !*tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be annotated as destructive", name)
			}
		case isDiscordMutatingTool(name):
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

func TestMCPDiscordToolValidation(t *testing.T) {
	session := newLocalDiscordMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name          string
		tool          string
		arguments     map[string]any
		expectedError string
	}{
		{name: "guild get missing guild", tool: "discord_guilds_get", arguments: map[string]any{"guild": ""}, expectedError: "guild is required"},
		{name: "channel get missing channel", tool: "discord_channels_get", arguments: map[string]any{"channel": ""}, expectedError: "channel is required"},
		{name: "messages edit reports first required field", tool: "discord_messages_edit", arguments: map[string]any{"channel": "", "message": "", "content": ""}, expectedError: "channel is required"},
		{name: "message delete missing message", tool: "discord_messages_delete", arguments: map[string]any{"channel": "C123", "message": ""}, expectedError: "message is required"},
		{name: "reaction add missing emoji", tool: "discord_reactions_add", arguments: map[string]any{"channel": "C123", "message": "M123", "emoji": ""}, expectedError: "emoji is required"},
		{name: "role create missing name", tool: "discord_roles_create", arguments: map[string]any{"guild": "G123", "name": ""}, expectedError: "name is required"},
		{name: "member add role missing role", tool: "discord_members_add_role", arguments: map[string]any{"guild": "G123", "user": "U123", "role": ""}, expectedError: "role is required"},
		{name: "member ban missing user", tool: "discord_members_ban", arguments: map[string]any{"guild": "G123", "user": ""}, expectedError: "user is required"},
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

func newLocalDiscordMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newDiscordMCPHTTPHandler(NewDiscordClient(Config{BaseURL: "http://127.0.0.1"})))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "discord-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: server.URL + defaultMCPEndpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func isDiscordDestructiveTool(name string) bool {
	return strings.Contains(name, "delete") ||
		strings.Contains(name, "remove") ||
		strings.HasSuffix(name, "_ban") ||
		strings.HasSuffix(name, "_unban")
}

func isDiscordMutatingTool(name string) bool {
	return strings.Contains(name, "send") ||
		strings.Contains(name, "edit") ||
		strings.Contains(name, "add") ||
		strings.Contains(name, "create")
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
