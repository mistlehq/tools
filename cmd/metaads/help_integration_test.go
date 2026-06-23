package main

import (
	"bytes"
	"strings"
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
	return commandResult{stdout: stdout, stderr: stderr}, err
}

func TestHelp(t *testing.T) {
	result, err := runCommandWithInput(t, Environment{}, "", "metaads", "help")
	if err != nil {
		t.Fatal(err)
	}
	output := result.stdout.String()
	for _, want := range []string{"metaads", "graph", "campaigns", "mcp"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected help to mention %q", want)
		}
	}
}

func TestLeafHelp(t *testing.T) {
	result, err := runCommandWithInput(t, Environment{}, "", "metaads", "graph", "request", "--help")
	if err != nil {
		t.Fatal(err)
	}
	output := result.stdout.String()
	for _, want := range []string{"--method <method>", "--path <path>", "complete Meta Ads API coverage"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected graph request help to mention %q", want)
		}
	}
}
