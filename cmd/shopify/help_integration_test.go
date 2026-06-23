package main

import (
	"strings"
	"testing"
)

func TestTopLevelHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "shopify", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expected := []string{
		"shopify auth help",
		"shopify graphql help",
		"shopify shop help",
		"shopify products help",
		"shopify orders help",
		"shopify customers help",
		"shopify inventory help",
		"shopify locations help",
		"shopify mcp help",
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
		{name: "auth", args: []string{"shopify", "auth", "help"}, want: "shopify auth test"},
		{name: "graphql", args: []string{"shopify", "graphql", "help"}, want: "shopify graphql request --help"},
		{name: "shop", args: []string{"shopify", "shop", "help"}, want: "shopify shop get"},
		{name: "products", args: []string{"shopify", "products", "help"}, want: "shopify products create --product-json <json>"},
		{name: "orders", args: []string{"shopify", "orders", "help"}, want: "shopify orders get --id <order-gid>"},
		{name: "customers", args: []string{"shopify", "customers", "help"}, want: "shopify customers get --id <customer-gid>"},
		{name: "inventory", args: []string{"shopify", "inventory", "help"}, want: "shopify inventory levels search --first <count>"},
		{name: "locations", args: []string{"shopify", "locations", "help"}, want: "shopify locations list --first <count>"},
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
