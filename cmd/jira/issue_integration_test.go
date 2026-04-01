package main

import (
	"strings"
	"testing"
)

func TestIssueGet(t *testing.T) {
	commandResult := setupAndRunCommandWithInput(t, "", "jira", "issue", "get", "KAN-1")
	stdout := commandResult.stdout

	if stdout.Len() == 0 {
		t.Fatal("expected stdout output")
	}

	output := strings.TrimSpace(stdout.String())
	lines := strings.Split(output, "\n")

	if len(lines) != 4 {
		t.Error("expected output to have 4 lines")
	}

	for i, line := range lines {
		if i == 0 && !strings.HasPrefix(line, "ID: ") {
			t.Error("expected line 1 to begin with 'ID: '")
		}

		if i == 1 && !strings.HasPrefix(line, "Key: ") {
			t.Error("expected line 2 to begin with 'Key: '")
		}

		if i == 2 && !strings.HasPrefix(line, "Summary: ") {
			t.Error("expected line 3 to begin with 'Summary: '")
		}

		if i == 3 && !strings.HasPrefix(line, "Status: ") {
			t.Error("expected line 4 to begin with 'Status: '")
		}
	}
}

func TestIssueSearch(t *testing.T) {
	commandResult := setupAndRunCommandWithInput(t, "", "jira", "issue", "search", "issuekey = KAN-1")
	stdout := commandResult.stdout

	if stdout.Len() == 0 {
		t.Fatal("expected stdout output")
	}

	output := strings.TrimSpace(stdout.String())
	lines := strings.Split(output, "\n")

	if strings.TrimSpace(lines[0]) == "" {
		t.Fatal("expected headers KEY, STATUS, and SUMMARY")
	}

	headerColumns := strings.Split(lines[0], "\t")

	if len(headerColumns) != 3 {
		t.Fatal("expected header row to have 3 columns")
	}

	if headerColumns[0] != "KEY" {
		t.Fatal("expected first header column to be KEY")
	}

	if headerColumns[1] != "STATUS" {
		t.Fatal("expected second header column to be STATUS")
	}

	if headerColumns[2] != "SUMMARY" {
		t.Fatal("expected third header column to be SUMMARY")
	}

	if len(lines[1:]) != 1 {
		t.Fatal("expected exactly one issue")
	}

	for _, line := range lines[1:] {
		if len(strings.Split(line, "\t")) != 3 {
			t.Fatal("expected row to have 3 columns")
		}
	}
}
