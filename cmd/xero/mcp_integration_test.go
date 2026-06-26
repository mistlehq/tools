package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHelp(t *testing.T) {
	commandResult, err := runCommand(t, Environment{}, "xero", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"xero mcp", "xero mcp serve", "Streamable HTTP"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp help to mention %q", want)
		}
	}
}

func TestMCPServeHelp(t *testing.T) {
	commandResult, err := runCommand(t, Environment{}, "xero", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"--addr <addr>", "--endpoint <path>", "xero_api_get"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp serve help to mention %q", want)
		}
	}
}

func TestMCPServerListsXeroTools(t *testing.T) {
	session := newLocalXeroMCPTestSession(t)
	defer session.Close()

	toolsResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	toolsByName := make(map[string]*mcp.Tool)
	for _, tool := range toolsResult.Tools {
		toolsByName[tool.Name] = tool
	}
	if len(toolsResult.Tools) != 5 || len(toolsByName) != 5 {
		t.Fatalf("expected exactly 5 MCP tools, got %d tools and %d unique names", len(toolsResult.Tools), len(toolsByName))
	}

	readOnlyTools := map[string]string{
		"xero_tenants_list": xeroTenantsListDoc,
		"xero_api_get":      xeroAPIGetDoc,
	}
	for name, description := range readOnlyTools {
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

	writeTools := map[string]string{
		"xero_api_post":   xeroAPIPostDoc,
		"xero_api_put":    xeroAPIPutDoc,
		"xero_api_delete": xeroAPIDeleteDoc,
	}
	for name, description := range writeTools {
		tool, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected MCP tool %q to be listed", name)
		}
		if tool.Description != description {
			t.Fatalf("expected MCP tool %q description %q, got %q", name, description, tool.Description)
		}
		if tool.Annotations == nil || tool.Annotations.ReadOnlyHint {
			t.Fatalf("expected MCP tool %q to be annotated as write-capable", name)
		}
	}
	assertBoolPtr(t, "xero_api_post destructive", toolsByName["xero_api_post"].Annotations.DestructiveHint, false)
	assertBoolPtr(t, "xero_api_put destructive", toolsByName["xero_api_put"].Annotations.DestructiveHint, false)
	if !toolsByName["xero_api_put"].Annotations.IdempotentHint {
		t.Fatal("expected MCP tool xero_api_put to be annotated as idempotent")
	}
	assertBoolPtr(t, "xero_api_delete destructive", toolsByName["xero_api_delete"].Annotations.DestructiveHint, true)
	for name, tool := range toolsByName {
		assertBoolPtr(t, name+" open world", tool.Annotations.OpenWorldHint, true)
	}
}

func assertBoolPtr(t *testing.T, label string, got *bool, want bool) {
	t.Helper()
	if got == nil || *got != want {
		t.Fatalf("expected %s hint %v, got %v", label, want, got)
	}
}

func TestMCPXeroToolValidation(t *testing.T) {
	session := newLocalXeroMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{name: "get missing tenant", tool: "xero_api_get", arguments: map[string]any{"family": "accounting", "endpoint": "/Invoices"}},
		{name: "get missing endpoint", tool: "xero_api_get", arguments: map[string]any{"family": "accounting", "tenantId": "tenant-123"}},
		{name: "post missing request", tool: "xero_api_post", arguments: map[string]any{"family": "accounting", "tenantId": "tenant-123", "endpoint": "/Invoices"}},
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

func newLocalXeroMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newXeroMCPHTTPHandler(NewXeroClient(Config{APIBaseURL: "http://127.0.0.1"})))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return connectXeroMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func connectXeroMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()

	client := mcp.NewClient(&mcp.Implementation{Name: "xero-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}
