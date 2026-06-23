package main

import (
	"strings"
	"testing"
)

func TestTopLevelHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "ga", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expected := []string{
		"ga auth help",
		"ga account-summaries help",
		"ga properties help",
		"ga metadata help",
		"ga compatibility help",
		"ga reports help",
		"ga google-ads-links help",
		"ga mcp help",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected top-level help to mention %q", want)
		}
	}
}

func TestNamespaceHelp(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want string
	}{
		{name: "auth", args: []string{"ga", "auth", "help"}, want: "ga auth test --property properties/<id>"},
		{name: "account summaries", args: []string{"ga", "account-summaries", "help"}, want: "ga account-summaries list"},
		{name: "properties", args: []string{"ga", "properties", "help"}, want: "ga properties get --property properties/<id>"},
		{name: "metadata", args: []string{"ga", "metadata", "help"}, want: "ga metadata get --property properties/<id>"},
		{name: "compatibility", args: []string{"ga", "compatibility", "help"}, want: "ga compatibility check --property properties/<id> --request-file <json>"},
		{name: "reports", args: []string{"ga", "reports", "help"}, want: "ga reports run --property properties/<id> --request-file <json>"},
		{name: "google ads links", args: []string{"ga", "google-ads-links", "help"}, want: "ga google-ads-links list --property properties/<id>"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			commandResult, err := runCommandWithInput(t, Environment{}, "", tc.args...)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(commandResult.stdout.String(), tc.want) {
				t.Fatalf("expected help to mention %q, got %q", tc.want, commandResult.stdout.String())
			}
		})
	}
}
