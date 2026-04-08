package main

import (
	"bytes"
	"fmt"
	"github.com/mistlehq/tools/internal/testproxy"
	"os"
	"strings"
	"testing"
	"time"
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

	upstreamBaseURL := getRequiredEnv(t, "SLACK_TEST_UPSTREAM_BASE_URL")
	token := getRequiredEnv(t, "SLACK_TEST_BOT_TOKEN")

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
			t.Errorf("failed to close proxy: %v", err)
		}
	})

	return Environment{
		"SLACK_BASE_URL": proxy.BaseURL,
	}
}

func setupSlackClient(t *testing.T) (Environment, SlackClient) {
	t.Helper()

	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	return env, NewSlackClient(config)
}

func uniqueTestMessage(prefix string) string {
	return fmt.Sprintf("%s %d", prefix, time.Now().UnixNano())
}

func parseLineValue(t *testing.T, output string, prefix string) string {
	t.Helper()

	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimPrefix(line, prefix)
		}
	}

	t.Fatalf("expected output to contain line with prefix %q, got %q", prefix, output)
	return ""
}
