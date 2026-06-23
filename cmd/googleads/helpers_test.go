package main

import (
	"bytes"
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

type commandResult struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func runCommandWithInput(t *testing.T, env Environment, input string, args ...string) (commandResult, error) {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cli := CLI{
		stdin:  bytes.NewBufferString(input),
		stdout: &stdout,
		stderr: &stderr,
		env:    env,
	}
	err := cli.run(args)
	return commandResult{stdout: stdout, stderr: stderr}, err
}

func getRequiredEnv(t *testing.T, name string) string {
	t.Helper()
	value := os.Getenv(name)
	if value == "" {
		t.Skipf("skipping: %s is not set", name)
	}
	return value
}

func getOptionalEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func setupCommandEnvironment(t *testing.T) Environment {
	t.Helper()
	apiVersion := getOptionalEnv("GOOGLEADS_TEST_API_VERSION", "v24")
	accessToken := getRequiredEnv(t, "GOOGLEADS_TEST_ACCESS_TOKEN")
	developerToken := getRequiredEnv(t, "GOOGLEADS_TEST_DEVELOPER_TOKEN")
	headers := map[string]string{"developer-token": developerToken}
	proxy, err := testproxy.Start(testproxy.Config{
		UpstreamBaseURL: "https://googleads.googleapis.com/" + apiVersion,
		AuthMode:        testproxy.AuthModeBearer,
		Token:           accessToken,
		Headers:         headers,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := proxy.Close(); err != nil {
			t.Errorf("failed to close proxy: %v", err)
		}
	})
	return Environment{"GOOGLEADS_BASE_URL": proxy.BaseURL}
}

func setupGoogleAdsClient(t *testing.T) (Environment, GoogleAdsClient) {
	t.Helper()
	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	return env, NewGoogleAdsClient(config)
}

func testCustomerID(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "GOOGLEADS_TEST_CUSTOMER_ID")
}

func optionalLoginCustomerArgs() []string {
	if loginCustomerID := os.Getenv("GOOGLEADS_TEST_LOGIN_CUSTOMER_ID"); loginCustomerID != "" {
		return []string{"--login-customer-id", loginCustomerID}
	}
	return nil
}

func optionalLoginCustomerID() string {
	return os.Getenv("GOOGLEADS_TEST_LOGIN_CUSTOMER_ID")
}

func decodeCommandJSON(t *testing.T, result commandResult, out any) {
	t.Helper()
	if err := json.Unmarshal(result.stdout.Bytes(), out); err != nil {
		t.Fatalf("expected valid JSON output: %v\nstdout: %s", err, result.stdout.String())
	}
}

func newLocalGoogleAdsMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	return newGoogleAdsMCPTestSession(t, NewGoogleAdsClient(Config{BaseURL: "http://127.0.0.1"}))
}

func newGoogleAdsMCPTestSession(t *testing.T, gc GoogleAdsClient) *mcp.ClientSession {
	t.Helper()
	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newGoogleAdsMCPHTTPHandler(gc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return connectGoogleAdsMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func connectGoogleAdsMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()
	client := mcp.NewClient(&mcp.Implementation{Name: "googleads-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func callGoogleAdsMCPTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
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

func stringsContainsAll(text string, values []string) bool {
	for _, value := range values {
		if !strings.Contains(text, value) {
			return false
		}
	}
	return true
}
