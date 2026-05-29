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

func TestMCPHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{
		"jira mcp",
		"jira mcp serve",
		"Streamable HTTP",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp help to mention %q", want)
		}
	}
}

func TestMCPServeHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{
		"--addr <addr>",
		"--endpoint <path>",
		"jira_issue_get",
		"jira_issue_search",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp serve help to mention %q", want)
		}
	}
}

func TestMCPServerListsJiraTools(t *testing.T) {
	jc := NewJiraClient(Config{BaseURL: "http://127.0.0.1"})
	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newJiraMCPHTTPHandler(jc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	session := connectJiraMCPTestClient(t, server.URL+defaultMCPEndpoint)
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
		"jira_auth_whoami":  jiraAuthWhoAmIDoc.Description,
		"jira_project_list": jiraProjectListDoc.Description,
		"jira_issue_get":    jiraIssueGetDoc.Description,
		"jira_issue_search": jiraIssueSearchDoc.Description,
	}

	for name, description := range expected {
		tool, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected MCP tool %q to be listed", name)
		}

		if tool.Description != description {
			t.Fatalf("expected MCP tool %q description %q, got %q", name, description, tool.Description)
		}

		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Fatalf("expected MCP tool %q to be annotated as read-only", name)
		}
	}
}

func TestMCPAuthWhoAmI(t *testing.T) {
	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newJiraMCPHTTPHandler(NewJiraClient(config)))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	session := connectJiraMCPTestClient(t, server.URL+defaultMCPEndpoint)
	defer session.Close()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "jira_auth_whoami",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.IsError {
		t.Fatalf("expected jira_auth_whoami to succeed, got tool error: %#v", result.Content)
	}

	var myself JiraMyself
	rawMyself, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(rawMyself, &myself); err != nil {
		t.Fatal(err)
	}

	if myself.AccountID == "" {
		t.Fatal("expected jira_auth_whoami structured output to include accountId")
	}
}

func connectJiraMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "jira-test-client",
		Version: "dev",
	}, nil)

	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{
		Endpoint: endpoint,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	return session
}
