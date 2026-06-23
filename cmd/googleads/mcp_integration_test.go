package main

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHelp(t *testing.T) {
	result, err := runCommandWithInput(t, Environment{}, "", "googleads", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}
	if !stringsContainsAll(result.stdout.String(), []string{"googleads mcp", "googleads mcp serve", "Streamable HTTP"}) {
		t.Fatalf("unexpected mcp help: %s", result.stdout.String())
	}
}

func TestMCPServerListsGoogleAdsTools(t *testing.T) {
	session := newLocalGoogleAdsMCPTestSession(t)
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
		"googleads_request":                   googleAdsRequestDoc.Description,
		"googleads_auth_test":                 googleAdsAuthTestDoc.Description,
		"googleads_customers_list_accessible": googleAdsCustomersListAccessibleDoc.Description,
		"googleads_gaql_search":               googleAdsGAQLSearchDoc.Description,
		"googleads_gaql_search_stream":        googleAdsGAQLSearchStreamDoc.Description,
		"googleads_fields_search":             googleAdsFieldsSearchDoc.Description,
		"googleads_field_get":                 googleAdsFieldGetDoc.Description,
	}
	for name, description := range expected {
		tool, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected MCP tool %q to be listed", name)
		}
		if tool.Description != description {
			t.Fatalf("expected MCP tool %q description %q, got %q", name, description, tool.Description)
		}
	}
}

func TestMCPGoogleAdsReadTools(t *testing.T) {
	_, gc := setupGoogleAdsClient(t)
	customerID := testCustomerID(t)
	session := newGoogleAdsMCPTestSession(t, gc)
	defer session.Close()

	rawResult := callGoogleAdsMCPTool(t, session, "googleads_request", map[string]any{
		"method": "GET",
		"path":   "/customers:listAccessibleCustomers",
	})
	var raw map[string]any
	decodeMCPStructuredContent(t, rawResult, &raw)
	if raw["resourceNames"] == nil {
		t.Fatalf("expected raw accessible customers response, got %#v", raw)
	}

	searchResult := callGoogleAdsMCPTool(t, session, "googleads_gaql_search", map[string]any{
		"customer_id": customerID,
		"query":       "SELECT customer.id FROM customer LIMIT 1",
	})
	var search map[string]any
	decodeMCPStructuredContent(t, searchResult, &search)
	if search["results"] == nil {
		t.Fatalf("expected GAQL search response, got %#v", search)
	}
}

func TestMCPGoogleAdsToolValidation(t *testing.T) {
	session := newLocalGoogleAdsMCPTestSession(t)
	defer session.Close()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: "googleads_gaql_search", Arguments: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected tool validation to return a tool error")
	}
}
