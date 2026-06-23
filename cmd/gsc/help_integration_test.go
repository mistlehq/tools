package main

import (
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "gsc", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"gsc", "sites", "searchanalytics", "url-inspection", "mcp"}
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
		{name: "auth", args: []string{"gsc", "auth", "help"}, want: "gsc auth test"},
		{name: "sites", args: []string{"gsc", "sites", "help"}, want: "gsc sites get"},
		{name: "searchanalytics", args: []string{"gsc", "searchanalytics", "help"}, want: "gsc searchanalytics query"},
		{name: "sitemaps", args: []string{"gsc", "sitemaps", "help"}, want: "gsc sitemaps get"},
		{name: "url inspection", args: []string{"gsc", "url-inspection", "help"}, want: "gsc url-inspection inspect"},
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
