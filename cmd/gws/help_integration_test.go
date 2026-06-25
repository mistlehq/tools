package main

import (
	"strings"
	"testing"
)

func TestHelpCommands(t *testing.T) {
	testCases := [][]string{
		{"gws", "help"},
		{"gws", "auth", "help"},
		{"gws", "request", "help"},
		{"gws", "drive", "help"},
		{"gws", "drive", "files", "help"},
		{"gws", "drive", "permissions", "help"},
		{"gws", "sheets", "help"},
		{"gws", "sheets", "spreadsheets", "help"},
		{"gws", "sheets", "values", "help"},
		{"gws", "docs", "help"},
		{"gws", "docs", "documents", "help"},
		{"gws", "slides", "help"},
		{"gws", "slides", "presentations", "help"},
		{"gws", "mcp", "help"},
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
