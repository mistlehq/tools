package main

import (
	"strings"
	"testing"
)

func TestIssueEditMeta(t *testing.T) {
	commandResult := setupAndRunCommand(t, "jira", "issue", "editmeta", "KAN-1")
	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")

	if len(lines) < 2 {
		t.Fatalf("expected editmeta output to include at least one field, got %q", output)
	}

	headerColumns := strings.Split(lines[0], "\t")
	if len(headerColumns) != 4 {
		t.Fatalf("expected editmeta header row to have 4 columns, got %q", lines[0])
	}

	if headerColumns[0] != "FIELD ID" || headerColumns[1] != "NAME" || headerColumns[2] != "REQUIRED" || headerColumns[3] != "TYPE" {
		t.Fatalf("unexpected editmeta headers: %q", lines[0])
	}

	foundSummary := false
	for _, line := range lines[1:] {
		columns := strings.Split(line, "\t")
		if len(columns) != 4 {
			t.Fatalf("expected editmeta row to have 4 columns, got %q", line)
		}

		if columns[0] == "summary" {
			foundSummary = true
		}
	}

	if !foundSummary {
		t.Fatal("expected editmeta output to include summary field")
	}
}
