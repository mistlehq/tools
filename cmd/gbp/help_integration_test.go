package main

import (
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "gbp", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{"gbp", "accounts", "locations", "reviews", "performance", "mcp"}
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
		{name: "auth", args: []string{"gbp", "auth", "help"}, want: "gbp auth test"},
		{name: "accounts", args: []string{"gbp", "accounts", "help"}, want: "gbp accounts get"},
		{name: "locations", args: []string{"gbp", "locations", "help"}, want: "gbp locations get"},
		{name: "reviews", args: []string{"gbp", "reviews", "help"}, want: "gbp reviews get"},
		{name: "media", args: []string{"gbp", "media", "help"}, want: "gbp media list"},
		{name: "local posts", args: []string{"gbp", "local-posts", "help"}, want: "gbp local-posts list"},
		{name: "performance", args: []string{"gbp", "performance", "help"}, want: "gbp performance search-keywords"},
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
