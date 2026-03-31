package main

import (
	"bytes"
	"os"
	"testing"
)

func getRequiredEnv(t *testing.T, name string) string {
	t.Helper()

	value := os.Getenv(name)
	if value == "" {
		t.Skipf("skipping: %s is not set", name)
	}

	return value
}

type commandResult struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func setupAndRunCommand(t *testing.T, args ...string) commandResult {
	t.Helper()

	upstreamBaseURL := getRequiredEnv(t, "JIRA_TEST_UPSTREAM_BASE_URL")
	username := getRequiredEnv(t, "JIRA_TEST_USERNAME")
	password := getRequiredEnv(t, "JIRA_TEST_PASSWORD")

	proxy, err := startProxyServer(proxyConfig{
		UpstreamBaseURL: upstreamBaseURL,
		AuthMode:        ProxyAuthModeBasic,
		Username:        username,
		Password:        password,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := proxy.Close(); err != nil {
			t.Errorf("failed to close proxy: %v", err)
		}
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cli := CLI{
		stdout: &stdout,
		stderr: &stderr,
		env: Environment{
			"JIRA_BASE_URL": proxy.BaseURL,
		},
	}

	err = cli.run(args)
	if err != nil {
		t.Fatal(err)
	}

	return commandResult{
		stdout: stdout,
		stderr: stderr,
	}
}
