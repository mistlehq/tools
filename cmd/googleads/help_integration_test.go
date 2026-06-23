package main

import (
	"strings"
	"testing"
)

func TestHelpCommands(t *testing.T) {
	testCases := [][]string{
		{"googleads", "help"},
		{"googleads", "auth", "help"},
		{"googleads", "request", "help"},
		{"googleads", "customers", "help"},
		{"googleads", "gaql", "help"},
		{"googleads", "fields", "help"},
		{"googleads", "mcp", "help"},
	}
	for _, args := range testCases {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			result, err := runCommandWithInput(t, Environment{}, "", args...)
			if err != nil {
				t.Fatal(err)
			}
			if result.stdout.Len() == 0 {
				t.Fatal("expected help output")
			}
		})
	}
}
