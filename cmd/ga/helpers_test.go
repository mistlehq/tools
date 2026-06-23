package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/mistlehq/tools/internal/testproxy"
	"golang.org/x/oauth2/google"
	"os"
	"testing"
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

	dataBaseURL := getRequiredEnv(t, "GA_TEST_ANALYTICS_DATA_BASE_URL")
	adminBaseURL := getRequiredEnv(t, "GA_TEST_ANALYTICS_ADMIN_BASE_URL")
	token := mintGoogleAnalyticsTestAccessToken(t)

	dataProxy, err := testproxy.Start(testproxy.Config{
		UpstreamBaseURL: dataBaseURL,
		AuthMode:        testproxy.AuthModeBearer,
		Token:           token,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := dataProxy.Close(); err != nil {
			t.Errorf("failed to close analytics data proxy: %v", err)
		}
	})

	adminProxy, err := testproxy.Start(testproxy.Config{
		UpstreamBaseURL: adminBaseURL,
		AuthMode:        testproxy.AuthModeBearer,
		Token:           token,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := adminProxy.Close(); err != nil {
			t.Errorf("failed to close analytics admin proxy: %v", err)
		}
	})

	return Environment{
		"GA_ANALYTICS_DATA_BASE_URL":  dataProxy.BaseURL,
		"GA_ANALYTICS_ADMIN_BASE_URL": adminProxy.BaseURL,
	}
}

func setupGAClient(t *testing.T) (Environment, GAClient) {
	t.Helper()

	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	return env, NewGAClient(config)
}

func mintGoogleAnalyticsTestAccessToken(t *testing.T) string {
	t.Helper()

	keyJSONBase64 := getRequiredEnv(t, "GA_TEST_SERVICE_ACCOUNT_KEY_JSON_BASE64")
	keyJSON, err := base64.StdEncoding.DecodeString(keyJSONBase64)
	if err != nil {
		t.Fatalf("failed to decode GA_TEST_SERVICE_ACCOUNT_KEY_JSON_BASE64: %v", err)
	}

	config, err := google.JWTConfigFromJSON(
		keyJSON,
		"https://www.googleapis.com/auth/analytics.readonly",
	)
	if err != nil {
		t.Fatalf("failed to parse service account JSON key: %v", err)
	}

	token, err := config.TokenSource(context.Background()).Token()
	if err != nil {
		t.Fatalf("failed to mint Google Analytics test access token: %v", err)
	}
	if token.AccessToken == "" {
		t.Fatal("Google Analytics test access token was empty")
	}
	return token.AccessToken
}

func testPropertyID(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "GA_TEST_PROPERTY_ID")
}

func testAccountID(t *testing.T) string {
	t.Helper()
	return getRequiredEnv(t, "GA_TEST_ACCOUNT_ID")
}

func minimalRunReportRequest() string {
	return `{"dateRanges":[{"startDate":"7daysAgo","endDate":"today"}],"metrics":[{"name":"activeUsers"}],"limit":"10"}`
}

func minimalRealtimeReportRequest() string {
	return `{"metrics":[{"name":"activeUsers"}],"limit":"10"}`
}

func minimalFunnelReportRequest() string {
	return `{"dateRanges":[{"startDate":"30daysAgo","endDate":"today"}],"funnel":{"steps":[{"name":"First visit","filterExpression":{"funnelEventFilter":{"eventName":"first_visit"}}},{"name":"Page view","filterExpression":{"funnelEventFilter":{"eventName":"page_view"}}}]}}`
}

func minimalCompatibilityRequest() string {
	return `{"dimensions":[{"name":"country"}],"metrics":[{"name":"activeUsers"}]}`
}

func writeTempJSONRequest(t *testing.T, body string) string {
	t.Helper()
	path := fmt.Sprintf("%s/request.json", t.TempDir())
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
