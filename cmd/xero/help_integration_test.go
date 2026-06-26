package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type commandResult struct {
	stdout bytes.Buffer
}

func runCommand(t *testing.T, env Environment, args ...string) (commandResult, error) {
	t.Helper()

	var stdout bytes.Buffer

	cli := CLI{
		stdout: &stdout,
		stderr: io.Discard,
		env:    env,
	}

	err := cli.run(args)
	return commandResult{stdout: stdout}, err
}

func TestHelp(t *testing.T) {
	commandResult, err := runCommand(t, Environment{}, "xero", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"xero", "tenants", "api", "mcp", "accounting", "projects"}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected help to mention %q", want)
		}
	}
}

func TestNamespaceHelp(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want string
	}{
		{name: "tenants", args: []string{"xero", "tenants", "help"}, want: "xero tenants list"},
		{name: "api", args: []string{"xero", "api", "help"}, want: "xero api get"},
		{name: "mcp", args: []string{"xero", "mcp", "help"}, want: "xero mcp serve"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			commandResult, err := runCommand(t, Environment{}, tc.args...)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(commandResult.stdout.String(), tc.want) {
				t.Fatalf("expected help to mention %q, got %q", tc.want, commandResult.stdout.String())
			}
		})
	}
}

func TestLeafHelpDoesNotRequireConfiguration(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want string
	}{
		{name: "tenants list", args: []string{"xero", "tenants", "list", "--help"}, want: "xero tenants list"},
		{name: "api get", args: []string{"xero", "api", "get", "--help"}, want: "xero api get"},
		{name: "api post", args: []string{"xero", "api", "post", "--help"}, want: "xero api post"},
		{name: "api put", args: []string{"xero", "api", "put", "--help"}, want: "xero api put"},
		{name: "api delete", args: []string{"xero", "api", "delete", "--help"}, want: "xero api delete"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			commandResult, err := runCommand(t, Environment{}, tc.args...)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(commandResult.stdout.String(), tc.want) {
				t.Fatalf("expected help to mention %q, got %q", tc.want, commandResult.stdout.String())
			}
		})
	}
}
