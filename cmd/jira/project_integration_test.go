package main

import (
	"strings"
	"testing"
)

func TestProjectList(t *testing.T) {
	commandResult := setupAndRunCommand(t, "jira", "project", "list")
	stdout := commandResult.stdout

	if stdout.Len() == 0 {
		t.Fatal("expected stdout output")
	}

	output := strings.TrimSpace(stdout.String())
	lines := strings.Split(output, "\n")

	if strings.TrimSpace(lines[0]) == "" {
		t.Fatal("expected headers ID, KEY, and NAME")
	}

	headerColumns := strings.Split(lines[0], "\t")

	if len(headerColumns) != 3 {
		t.Fatal("expected header row to have 3 columns")
	}

	if headerColumns[0] != "ID" {
		t.Fatal("expected first header column to be ID")
	}

	if headerColumns[1] != "KEY" {
		t.Fatal("expected second header column to be KEY")
	}

	if headerColumns[2] != "NAME" {
		t.Fatal("expected third header column to be NAME")
	}

	if len(lines[1:]) == 0 {
		t.Fatal("expected at least one project")
	}

	for _, line := range lines[1:] {
		if len(strings.Split(line, "\t")) != 3 {
			t.Fatal("expected row to have 3 columns")
		}
	}
}
