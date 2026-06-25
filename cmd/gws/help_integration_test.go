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
		{"gws", "gmail", "help"},
		{"gws", "gmail", "messages", "help"},
		{"gws", "gmail", "drafts", "help"},
		{"gws", "calendar", "help"},
		{"gws", "calendar", "calendar-list", "help"},
		{"gws", "calendar", "events", "help"},
		{"gws", "calendar", "freebusy", "help"},
		{"gws", "chat", "help"},
		{"gws", "chat", "spaces", "help"},
		{"gws", "chat", "messages", "help"},
		{"gws", "chat", "members", "help"},
		{"gws", "people", "help"},
		{"gws", "people", "people", "help"},
		{"gws", "people", "connections", "help"},
		{"gws", "people", "search-contacts", "help"},
		{"gws", "people", "search-directory", "help"},
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
