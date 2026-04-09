package main

import (
	"fmt"
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

func TestIssueDelete(t *testing.T) {
	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	jc := NewJiraClient(config)
	template, err := getJiraTestIssueTemplate(jc, jiraTestTemplateIssueKey)
	if err != nil {
		t.Fatal(err)
	}

	issue, err := createJiraTestIssue(jc, template, fmt.Sprintf("delete integration test %s", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	deleted := false
	t.Cleanup(func() {
		if deleted {
			return
		}

		if err := deleteJiraTestIssue(jc, issue.Key); err != nil {
			t.Errorf("failed to delete issue %s: %v", issue.Key, err)
		}
	})

	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "delete", issue.Key)
	if err != nil {
		t.Fatal(err)
	}

	deleted = true

	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d: %q", len(lines), output)
	}

	if lines[0] != "Issue: "+issue.Key {
		t.Fatalf("unexpected issue line: %q", lines[0])
	}

	if lines[1] != "Deleted: true" {
		t.Fatalf("unexpected deleted line: %q", lines[1])
	}

	if _, err := jc.GetIssue(issue.Key); err == nil {
		t.Fatalf("expected deleted issue %s to be unavailable", issue.Key)
	}
}

func TestIssueDeleteRequiresIssueKey(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "delete")
	if err == nil {
		t.Fatal("expected missing issue key to fail")
	}

	if err.Error() != "issue delete expects exactly 1 positional argument" {
		t.Fatalf("unexpected error: %v", err)
	}
}
