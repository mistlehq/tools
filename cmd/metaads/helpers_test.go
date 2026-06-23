package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mistlehq/tools/internal/testproxy"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func getRequiredEnv(t *testing.T, name string) string {
	t.Helper()
	value := os.Getenv(name)
	if value == "" {
		t.Skipf("skipping: %s is not set", name)
	}
	return value
}

func setupCommandEnvironment(t *testing.T) Environment {
	t.Helper()
	apiVersion := getOptionalEnv("METAADS_TEST_GRAPH_API_VERSION", "v25.0")
	accessToken := getRequiredEnv(t, "METAADS_TEST_ACCESS_TOKEN")

	proxy, err := testproxy.Start(testproxy.Config{
		UpstreamBaseURL: "https://graph.facebook.com/" + apiVersion,
		AuthMode:        testproxy.AuthModeBearer,
		Token:           accessToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := proxy.Close(); err != nil {
			t.Errorf("failed to close proxy: %v", err)
		}
	})
	return Environment{"METAADS_GRAPH_BASE_URL": proxy.BaseURL}
}

func setupMetaAdsClient(t *testing.T) (Environment, MetaAdsClient) {
	t.Helper()
	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	return env, NewMetaAdsClient(config)
}

func getOptionalEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func testAdAccountID(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "METAADS_TEST_AD_ACCOUNT_ID")
}

func testCampaignID(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "METAADS_TEST_CAMPAIGN_ID")
}

func decodeCommandJSON(t *testing.T, result commandResult, out any) {
	t.Helper()
	if err := json.Unmarshal(result.stdout.Bytes(), out); err != nil {
		t.Fatalf("expected valid JSON output: %v\nstdout: %s", err, result.stdout.String())
	}
}

func newLocalMetaAdsMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	return newMetaAdsMCPTestSession(t, NewMetaAdsClient(Config{GraphBaseURL: "http://127.0.0.1"}))
}

func newMetaAdsMCPTestSession(t *testing.T, mc MetaAdsClient) *mcp.ClientSession {
	t.Helper()
	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newMetaAdsMCPHTTPHandler(mc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return connectMetaAdsMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func connectMetaAdsMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()
	client := mcp.NewClient(&mcp.Implementation{Name: "metaads-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func callMetaAdsMCPTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
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
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, output); err != nil {
		t.Fatalf("expected structured content to decode: %v\ncontent: %s", err, string(raw))
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
