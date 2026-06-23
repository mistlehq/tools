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
	commandResult, err := runCommandWithInput(t, Environment{}, "", "ga", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"ga mcp", "ga mcp serve", "Streamable HTTP"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp help to mention %q", want)
		}
	}
}

func TestMCPServeHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "ga", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"--addr <addr>", "--endpoint <path>", "ga_report_run"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp serve help to mention %q", want)
		}
	}
}

func TestMCPServerListsGATools(t *testing.T) {
	session := newLocalGAMCPTestSession(t)
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
		"ga_auth_test":              gaAuthTestDoc.Description,
		"ga_account_summaries_list": gaAccountSummariesListDoc.Description,
		"ga_property_get":           gaPropertyGetDoc.Description,
		"ga_metadata_get":           gaMetadataGetDoc.Description,
		"ga_compatibility_check":    gaCompatibilityCheckDoc.Description,
		"ga_report_run":             gaReportRunDoc.Description,
		"ga_report_realtime":        gaReportRealtimeDoc.Description,
		"ga_report_funnel":          gaReportFunnelDoc.Description,
		"ga_google_ads_links_list":  gaGoogleAdsLinksListDoc.Description,
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

func TestMCPGATools(t *testing.T) {
	_, gc := setupGAClient(t)
	propertyID := testPropertyID(t)
	session := newGAMCPTestSession(t, gc)
	defer session.Close()

	authResult := callGAMCPTool(t, session, "ga_auth_test", map[string]any{
		"property": propertyID,
	})
	var property GAProperty
	decodeMCPStructuredContent(t, authResult, &property)
	if property.Name != propertyID {
		t.Fatalf("expected property %q, got %#v", propertyID, property)
	}

	summariesResult := callGAMCPTool(t, session, "ga_account_summaries_list", map[string]any{})
	var summaries GAAccountSummariesList
	decodeMCPStructuredContent(t, summariesResult, &summaries)
	if len(summaries.AccountSummaries) == 0 {
		t.Fatal("expected account summaries")
	}

	propertyResult := callGAMCPTool(t, session, "ga_property_get", map[string]any{
		"property": propertyID,
	})
	var fetchedProperty GAProperty
	decodeMCPStructuredContent(t, propertyResult, &fetchedProperty)
	if fetchedProperty.Name != propertyID {
		t.Fatalf("expected fetched property %q, got %#v", propertyID, fetchedProperty)
	}

	metadataResult := callGAMCPTool(t, session, "ga_metadata_get", map[string]any{
		"property": propertyID,
	})
	var metadata GAMetadata
	decodeMCPStructuredContent(t, metadataResult, &metadata)
	if len(metadata.Dimensions) == 0 || len(metadata.Metrics) == 0 {
		t.Fatalf("expected metadata dimensions and metrics, got %#v", metadata)
	}

	compatibilityResult := callGAMCPTool(t, session, "ga_compatibility_check", map[string]any{
		"property": propertyID,
		"request": map[string]any{
			"dimensions": []map[string]string{{"name": "country"}},
			"metrics":    []map[string]string{{"name": "activeUsers"}},
		},
	})
	var compatibility map[string]any
	decodeMCPStructuredContent(t, compatibilityResult, &compatibility)
	if len(compatibility) == 0 {
		t.Fatal("expected compatibility response")
	}

	reportResult := callGAMCPTool(t, session, "ga_report_run", map[string]any{
		"property": propertyID,
		"request": map[string]any{
			"dateRanges": []map[string]string{{"startDate": "7daysAgo", "endDate": "today"}},
			"metrics":    []map[string]string{{"name": "activeUsers"}},
			"limit":      "10",
		},
	})
	var report map[string]any
	decodeMCPStructuredContent(t, reportResult, &report)
	if report["kind"] != "analyticsData#runReport" {
		t.Fatalf("expected report kind analyticsData#runReport, got %#v", report)
	}

	realtimeResult := callGAMCPTool(t, session, "ga_report_realtime", map[string]any{
		"property": propertyID,
		"request": map[string]any{
			"metrics": []map[string]string{{"name": "activeUsers"}},
			"limit":   "10",
		},
	})
	var realtime map[string]any
	decodeMCPStructuredContent(t, realtimeResult, &realtime)
	if realtime["kind"] != "analyticsData#runRealtimeReport" {
		t.Fatalf("expected realtime report kind analyticsData#runRealtimeReport, got %#v", realtime)
	}

	funnelResult := callGAMCPTool(t, session, "ga_report_funnel", map[string]any{
		"property": propertyID,
		"request": map[string]any{
			"dateRanges": []map[string]string{{"startDate": "30daysAgo", "endDate": "today"}},
			"funnel": map[string]any{
				"steps": []map[string]any{
					{
						"name": "First visit",
						"filterExpression": map[string]any{
							"funnelEventFilter": map[string]any{"eventName": "first_visit"},
						},
					},
					{
						"name": "Page view",
						"filterExpression": map[string]any{
							"funnelEventFilter": map[string]any{"eventName": "page_view"},
						},
					},
				},
			},
		},
	})
	var funnel map[string]any
	decodeMCPStructuredContent(t, funnelResult, &funnel)
	if _, ok := funnel["funnelTable"]; !ok {
		t.Fatalf("expected funnel report response to include funnelTable, got %#v", funnel)
	}

	googleAdsLinksResult := callGAMCPTool(t, session, "ga_google_ads_links_list", map[string]any{
		"property": propertyID,
	})
	var googleAdsLinks GAGoogleAdsLinksList
	decodeMCPStructuredContent(t, googleAdsLinksResult, &googleAdsLinks)
}

func TestMCPGAToolValidation(t *testing.T) {
	session := newLocalGAMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{name: "property missing property", tool: "ga_property_get", arguments: map[string]any{}},
		{name: "metadata missing property", tool: "ga_metadata_get", arguments: map[string]any{}},
		{name: "report missing request", tool: "ga_report_run", arguments: map[string]any{"property": "properties/123"}},
		{name: "google ads links missing property", tool: "ga_google_ads_links_list", arguments: map[string]any{}},
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

func newLocalGAMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	return newGAMCPTestSession(t, NewGAClient(Config{
		AnalyticsDataBaseURL:  "http://127.0.0.1",
		AnalyticsAdminBaseURL: "http://127.0.0.1",
	}))
}

func newGAMCPTestSession(t *testing.T, gc GAClient) *mcp.ClientSession {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newGAMCPHTTPHandler(gc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return connectGAMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func connectGAMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()

	client := mcp.NewClient(&mcp.Implementation{Name: "ga-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func callGAMCPTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
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
