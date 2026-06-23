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
	commandResult, err := runCommandWithInput(t, Environment{}, "", "gsc", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"gsc mcp", "gsc mcp serve", "Streamable HTTP"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp help to mention %q", want)
		}
	}
}

func TestMCPServeHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "gsc", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"--addr <addr>", "--endpoint <path>", "gsc_searchanalytics_query"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp serve help to mention %q", want)
		}
	}
}

func TestMCPServerListsGSCTools(t *testing.T) {
	session := newLocalGSCMCPTestSession(t)
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
		"gsc_auth_test":              gscAuthTestDoc.Description,
		"gsc_sites_list":             gscSitesListDoc.Description,
		"gsc_site_get":               gscSiteGetDoc.Description,
		"gsc_searchanalytics_query":  gscSearchAnalyticsQueryDoc.Description,
		"gsc_sitemaps_list":          gscSitemapsListDoc.Description,
		"gsc_sitemap_get":            gscSitemapGetDoc.Description,
		"gsc_url_inspection_inspect": gscURLInspectionInspectDoc.Description,
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

func TestMCPGSCTools(t *testing.T) {
	_, gc := setupGSCClient(t)
	siteURL := testSiteURL(t)
	inspectionURL := testInspectionURL(t)
	sitemapURL := testSitemapURL(t)
	session := newGSCMCPTestSession(t, gc)
	defer session.Close()

	authResult := callGSCMCPTool(t, session, "gsc_auth_test", map[string]any{
		"siteUrl": siteURL,
	})
	var site GSCSite
	decodeMCPStructuredContent(t, authResult, &site)
	if site.SiteURL != siteURL {
		t.Fatalf("expected site %q, got %#v", siteURL, site)
	}

	sitesResult := callGSCMCPTool(t, session, "gsc_sites_list", map[string]any{})
	var sites GSCSitesList
	decodeMCPStructuredContent(t, sitesResult, &sites)
	if !containsSite(sites.SiteEntry, siteURL) {
		t.Fatalf("expected sites list to include %q, got %#v", siteURL, sites.SiteEntry)
	}

	siteResult := callGSCMCPTool(t, session, "gsc_site_get", map[string]any{
		"siteUrl": siteURL,
	})
	var fetchedSite GSCSite
	decodeMCPStructuredContent(t, siteResult, &fetchedSite)
	if fetchedSite.SiteURL != siteURL {
		t.Fatalf("expected fetched site %q, got %#v", siteURL, fetchedSite)
	}

	searchAnalyticsResult := callGSCMCPTool(t, session, "gsc_searchanalytics_query", map[string]any{
		"siteUrl": siteURL,
		"request": map[string]any{
			"startDate": recentSearchAnalyticsStartDate(),
			"endDate":   recentSearchAnalyticsEndDate(),
			"dimensions": []string{
				"query",
			},
			"rowLimit": 10,
		},
	})
	var searchAnalytics map[string]any
	decodeMCPStructuredContent(t, searchAnalyticsResult, &searchAnalytics)

	sitemapsResult := callGSCMCPTool(t, session, "gsc_sitemaps_list", map[string]any{
		"siteUrl": siteURL,
	})
	var sitemaps GSCSitemapsList
	decodeMCPStructuredContent(t, sitemapsResult, &sitemaps)

	sitemapResult := callGSCMCPTool(t, session, "gsc_sitemap_get", map[string]any{
		"siteUrl":  siteURL,
		"feedPath": sitemapURL,
	})
	var sitemap GSCSitemap
	decodeMCPStructuredContent(t, sitemapResult, &sitemap)
	if sitemap.Path != sitemapURL {
		t.Fatalf("expected sitemap path %q, got %#v", sitemapURL, sitemap)
	}

	inspectionResult := callGSCMCPTool(t, session, "gsc_url_inspection_inspect", map[string]any{
		"request": map[string]any{
			"siteUrl":       siteURL,
			"inspectionUrl": inspectionURL,
		},
	})
	var inspection map[string]any
	decodeMCPStructuredContent(t, inspectionResult, &inspection)
	if _, ok := inspection["inspectionResult"]; !ok {
		t.Fatalf("expected URL inspection response to include inspectionResult, got %#v", inspection)
	}
}

func TestMCPGSCToolValidation(t *testing.T) {
	session := newLocalGSCMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{name: "site get missing site URL", tool: "gsc_site_get", arguments: map[string]any{}},
		{name: "search analytics missing request", tool: "gsc_searchanalytics_query", arguments: map[string]any{"siteUrl": "https://example.com/"}},
		{name: "sitemap get missing feed path", tool: "gsc_sitemap_get", arguments: map[string]any{"siteUrl": "https://example.com/"}},
		{name: "url inspection missing request", tool: "gsc_url_inspection_inspect", arguments: map[string]any{}},
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

func newLocalGSCMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	return newGSCMCPTestSession(t, NewGSCClient(Config{
		SearchConsoleBaseURL: "http://127.0.0.1",
	}))
}

func newGSCMCPTestSession(t *testing.T, gc GSCClient) *mcp.ClientSession {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newGSCMCPHTTPHandler(gc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return connectGSCMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func connectGSCMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()

	client := mcp.NewClient(&mcp.Implementation{Name: "gsc-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func callGSCMCPTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
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
