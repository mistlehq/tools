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
	commandResult, err := runCommandWithInput(t, Environment{}, "", "gbp", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"gbp mcp", "gbp mcp serve", "Streamable HTTP"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp help to mention %q", want)
		}
	}
}

func TestMCPServeHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "gbp", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"--addr <addr>", "--endpoint <path>", "gbp_performance_search_keywords"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp serve help to mention %q", want)
		}
	}
}

func TestMCPServerListsGBPTools(t *testing.T) {
	session := newLocalGBPMCPTestSession(t)
	defer session.Close()

	toolsResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	toolsByName := make(map[string]*mcp.Tool)
	for _, tool := range toolsResult.Tools {
		toolsByName[tool.Name] = tool
	}

	readOnlyTools := map[string]string{
		"gbp_auth_test":                   gbpAuthTestDoc.Description,
		"gbp_accounts_list":               gbpAccountsListDoc.Description,
		"gbp_account_get":                 gbpAccountGetDoc.Description,
		"gbp_locations_list":              gbpLocationsListDoc.Description,
		"gbp_location_get":                gbpLocationGetDoc.Description,
		"gbp_reviews_list":                gbpReviewsListDoc.Description,
		"gbp_review_get":                  gbpReviewGetDoc.Description,
		"gbp_media_list":                  gbpMediaListDoc.Description,
		"gbp_media_get":                   gbpMediaGetDoc.Description,
		"gbp_local_posts_list":            gbpLocalPostsListDoc.Description,
		"gbp_local_post_get":              gbpLocalPostGetDoc.Description,
		"gbp_performance_daily_metrics":   gbpPerformanceDailyMetricsDoc.Description,
		"gbp_performance_search_keywords": gbpPerformanceSearchKeywordsDoc.Description,
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
		"gbp_location_create":            gbpLocationCreateDoc.Description,
		"gbp_location_patch":             gbpLocationPatchDoc.Description,
		"gbp_location_delete":            gbpLocationDeleteDoc.Description,
		"gbp_review_update_reply":        gbpReviewUpdateReplyDoc.Description,
		"gbp_review_delete_reply":        gbpReviewDeleteReplyDoc.Description,
		"gbp_media_create":               gbpMediaCreateDoc.Description,
		"gbp_media_patch":                gbpMediaPatchDoc.Description,
		"gbp_media_delete":               gbpMediaDeleteDoc.Description,
		"gbp_media_start_upload":         gbpMediaStartUploadDoc.Description,
		"gbp_local_post_create":          gbpLocalPostCreateDoc.Description,
		"gbp_local_post_patch":           gbpLocalPostPatchDoc.Description,
		"gbp_local_post_delete":          gbpLocalPostDeleteDoc.Description,
		"gbp_local_post_report_insights": gbpLocalPostReportInsightsDoc.Description,
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
}

func TestMCPGBPTools(t *testing.T) {
	_, gc, _ := setupGBPClient(t)
	session := newGBPMCPTestSession(t, gc)
	defer session.Close()

	authResult := callGBPMCPTool(t, session, "gbp_auth_test", map[string]any{})
	var accounts GBPAccountsList
	decodeMCPStructuredContent(t, authResult, &accounts)
	if !containsAccount(accounts.Accounts, testAccount) {
		t.Fatalf("expected accounts to include %q, got %#v", testAccount, accounts)
	}

	locationResult := callGBPMCPTool(t, session, "gbp_location_get", map[string]any{
		"location": testLocation,
		"readMask": "name,title",
	})
	var location GBPLocation
	decodeMCPStructuredContent(t, locationResult, &location)
	if stringField(location, "name") != testLocation {
		t.Fatalf("expected location %q, got %#v", testLocation, location)
	}

	reviewsResult := callGBPMCPTool(t, session, "gbp_reviews_list", map[string]any{
		"account":  testAccount,
		"location": testLocation,
	})
	var reviews GBPReviewsList
	decodeMCPStructuredContent(t, reviewsResult, &reviews)
	if _, ok := reviews["reviews"]; !ok {
		t.Fatalf("expected reviews response to include reviews, got %#v", reviews)
	}

	replyResult := callGBPMCPTool(t, session, "gbp_review_update_reply", map[string]any{
		"account":  testAccount,
		"location": testLocation,
		"review":   "abc",
		"request": map[string]any{
			"comment": "Thanks!",
		},
	})
	var reply GBPWriteResult
	decodeMCPStructuredContent(t, replyResult, &reply)
	if reply["comment"] != "Thanks!" {
		t.Fatalf("expected reply result to include comment, got %#v", reply)
	}

	mediaCreateResult := callGBPMCPTool(t, session, "gbp_media_create", map[string]any{
		"account":  testAccount,
		"location": testLocation,
		"request": map[string]any{
			"mediaFormat": "PHOTO",
		},
	})
	var media GBPMediaItem
	decodeMCPStructuredContent(t, mediaCreateResult, &media)
	if stringField(media, "name") != "accounts/123/locations/456/media/789" {
		t.Fatalf("expected media name, got %#v", media)
	}

	localPostCreateResult := callGBPMCPTool(t, session, "gbp_local_post_create", map[string]any{
		"account":  testAccount,
		"location": testLocation,
		"request": map[string]any{
			"summary": "Hello",
		},
	})
	var localPost GBPLocalPost
	decodeMCPStructuredContent(t, localPostCreateResult, &localPost)
	if stringField(localPost, "name") != "accounts/123/locations/456/localPosts/post-1" {
		t.Fatalf("expected local post name, got %#v", localPost)
	}

	performanceResult := callGBPMCPTool(t, session, "gbp_performance_daily_metrics", map[string]any{
		"location": testLocation,
		"request": map[string]any{
			"dailyMetric": "WEBSITE_CLICKS",
			"dailyRange": map[string]any{
				"start_date": map[string]any{"year": 2026, "month": 6, "day": 1},
				"end_date":   map[string]any{"year": 2026, "month": 6, "day": 2},
			},
		},
	})
	var performance GBPPerformanceResult
	decodeMCPStructuredContent(t, performanceResult, &performance)
	if _, ok := performance["timeSeries"]; !ok {
		t.Fatalf("expected performance response to include timeSeries, got %#v", performance)
	}
}

func TestMCPGBPToolValidation(t *testing.T) {
	session := newLocalGBPMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{name: "account get missing account", tool: "gbp_account_get", arguments: map[string]any{}},
		{name: "location get missing read mask", tool: "gbp_location_get", arguments: map[string]any{"location": testLocation}},
		{name: "reviews list missing location", tool: "gbp_reviews_list", arguments: map[string]any{"account": testAccount}},
		{name: "performance missing request", tool: "gbp_performance_daily_metrics", arguments: map[string]any{"location": testLocation}},
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

func newLocalGBPMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	return newGBPMCPTestSession(t, NewGBPClient(Config{
		AccountManagementBaseURL:   "http://127.0.0.1",
		BusinessInformationBaseURL: "http://127.0.0.1",
		PerformanceBaseURL:         "http://127.0.0.1",
		MyBusinessBaseURL:          "http://127.0.0.1",
	}))
}

func newGBPMCPTestSession(t *testing.T, gc GBPClient) *mcp.ClientSession {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newGBPMCPHTTPHandler(gc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return connectGBPMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func connectGBPMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()

	client := mcp.NewClient(&mcp.Implementation{Name: "gbp-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func callGBPMCPTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
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
