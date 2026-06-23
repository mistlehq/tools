package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mistlehq/tools/internal/testproxy"
	"golang.org/x/oauth2/google"
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

	upstreamBaseURL := getRequiredEnv(t, "GSC_TEST_SEARCH_CONSOLE_BASE_URL")
	token := mintGoogleSearchConsoleTestAccessToken(t)

	proxy, err := testproxy.Start(testproxy.Config{
		UpstreamBaseURL: upstreamBaseURL,
		AuthMode:        testproxy.AuthModeBearer,
		Token:           token,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := proxy.Close(); err != nil {
			t.Errorf("failed to close search console proxy: %v", err)
		}
	})

	return Environment{
		"GSC_SEARCH_CONSOLE_BASE_URL": proxy.BaseURL,
	}
}

func setupGSCClient(t *testing.T) (Environment, GSCClient) {
	t.Helper()

	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	return env, NewGSCClient(config)
}

func mintGoogleSearchConsoleTestAccessToken(t *testing.T) string {
	t.Helper()

	keyJSONBase64 := getRequiredEnv(t, "GSC_TEST_SERVICE_ACCOUNT_KEY_JSON_BASE64")
	keyJSON, err := base64.StdEncoding.DecodeString(keyJSONBase64)
	if err != nil {
		t.Fatalf("failed to decode GSC_TEST_SERVICE_ACCOUNT_KEY_JSON_BASE64: %v", err)
	}

	config, err := google.JWTConfigFromJSON(
		keyJSON,
		"https://www.googleapis.com/auth/webmasters.readonly",
	)
	if err != nil {
		t.Fatalf("failed to parse service account JSON key: %v", err)
	}

	token, err := config.TokenSource(context.Background()).Token()
	if err != nil {
		t.Fatalf("failed to mint Google Search Console test access token: %v", err)
	}
	if token.AccessToken == "" {
		t.Fatal("Google Search Console test access token was empty")
	}
	return token.AccessToken
}

func testSiteURL(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "GSC_TEST_SITE_URL")
}

func testInspectionURL(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "GSC_TEST_INSPECTION_URL")
}

func testSitemapURL(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "GSC_TEST_SITEMAP_URL")
}

func minimalSearchAnalyticsRequest() string {
	return fmt.Sprintf(`{"startDate":%q,"endDate":%q,"dimensions":["query"],"rowLimit":10}`, recentSearchAnalyticsStartDate(), recentSearchAnalyticsEndDate())
}

func recentSearchAnalyticsStartDate() string {
	return time.Now().AddDate(0, 0, -9).Format("2006-01-02")
}

func recentSearchAnalyticsEndDate() string {
	return time.Now().AddDate(0, 0, -2).Format("2006-01-02")
}

func minimalURLInspectionRequest(siteURL string, inspectionURL string) string {
	return fmt.Sprintf(`{"inspectionUrl":%q,"siteUrl":%q}`, inspectionURL, siteURL)
}

func writeTempJSONRequest(t *testing.T, body string) string {
	t.Helper()
	path := fmt.Sprintf("%s/request.json", t.TempDir())
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func validUnitEnv() Environment {
	return Environment{
		"GSC_SEARCH_CONSOLE_BASE_URL": "http://127.0.0.1",
	}
}
