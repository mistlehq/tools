package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mistlehq/tools/internal/testproxy"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

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
	return commandResult{
		stdout: stdout,
		stderr: stderr,
	}, err
}

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

	shopDomain := getRequiredEnv(t, "SHOPIFY_TEST_SHOP_DOMAIN")
	apiVersion := getRequiredEnv(t, "SHOPIFY_TEST_ADMIN_API_VERSION")
	accessToken := mintShopifyTestAccessToken(t, shopDomain)

	proxy, err := testproxy.Start(testproxy.Config{
		UpstreamBaseURL: "https://" + shopDomain + "/admin/api/" + apiVersion,
		AuthMode:        testproxy.AuthModeHeader,
		HeaderName:      "X-Shopify-Access-Token",
		HeaderValue:     accessToken,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := proxy.Close(); err != nil {
			t.Errorf("failed to close proxy: %v", err)
		}
	})

	return Environment{
		"SHOPIFY_ADMIN_BASE_URL": proxy.BaseURL,
	}
}

func setupShopifyClient(t *testing.T) (Environment, ShopifyClient) {
	t.Helper()

	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	return env, NewShopifyClient(config)
}

func mintShopifyTestAccessToken(t *testing.T, shopDomain string) string {
	t.Helper()

	clientID := getRequiredEnv(t, "SHOPIFY_TEST_CLIENT_ID")
	clientSecret := getRequiredEnv(t, "SHOPIFY_TEST_CLIENT_SECRET")

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	response, err := http.PostForm("https://"+shopDomain+"/admin/oauth/access_token", form)
	if err != nil {
		t.Fatalf("failed to mint Shopify test access token: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read Shopify token response: %v", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		t.Fatalf("failed to mint Shopify test access token with status %d: %s", response.StatusCode, string(body))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		t.Fatalf("failed to decode Shopify token response: %v", err)
	}
	if tokenResponse.AccessToken == "" {
		t.Fatal("Shopify test access token was empty")
	}
	return tokenResponse.AccessToken
}

func testProductID(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "SHOPIFY_TEST_PRODUCT_ID")
}

func testProductHandle(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "SHOPIFY_TEST_PRODUCT_HANDLE")
}

func testOrderID(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "SHOPIFY_TEST_ORDER_ID")
}

func testOrderName(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "SHOPIFY_TEST_ORDER_NAME")
}

func testCustomerID(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "SHOPIFY_TEST_CUSTOMER_ID")
}

func testCustomerEmail(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "SHOPIFY_TEST_CUSTOMER_EMAIL")
}

func uniqueProductTitle(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("mistle-tools test %s %d", strings.ReplaceAll(t.Name(), "/", "-"), time.Now().UnixNano())
}

func writeTempTextFile(t *testing.T, body string) string {
	t.Helper()
	path := fmt.Sprintf("%s/request.txt", t.TempDir())
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func decodeCommandJSON(t *testing.T, result commandResult, out any) {
	t.Helper()
	if err := json.Unmarshal(result.stdout.Bytes(), out); err != nil {
		t.Fatalf("expected valid JSON output: %v\nstdout: %s", err, result.stdout.String())
	}
}

func newLocalShopifyMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	return newShopifyMCPTestSession(t, NewShopifyClient(Config{AdminBaseURL: "http://127.0.0.1"}))
}

func newShopifyMCPTestSession(t *testing.T, sc ShopifyClient) *mcp.ClientSession {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newShopifyMCPHTTPHandler(sc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return connectShopifyMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func connectShopifyMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()

	client := mcp.NewClient(&mcp.Implementation{Name: "shopify-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func callShopifyMCPTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
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

func toolErrorContains(result *mcp.CallToolResult, text string) bool {
	for _, content := range result.Content {
		textContent, ok := content.(*mcp.TextContent)
		if ok && strings.Contains(textContent.Text, text) {
			return true
		}
	}
	return false
}

func isProtectedCustomerDataError(err error) bool {
	return err != nil &&
		strings.Contains(err.Error(), "not approved to access") &&
		strings.Contains(err.Error(), "protected-customer-data")
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
