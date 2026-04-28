package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestIssueCreate(t *testing.T) {
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

	summary := fmt.Sprintf("create integration test %d", time.Now().UnixNano())
	commandResult, err := runCommandWithInput(
		t,
		env,
		"created from stdin",
		"jira",
		"issue",
		"create",
		"--project-id",
		template.Fields.Project.ID,
		"--issue-type-id",
		template.Fields.IssueType.ID,
		"--summary",
		summary,
		"--description-file",
		"-",
	)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 output lines, got %d: %q", len(lines), output)
	}

	if !strings.HasPrefix(lines[0], "ID: ") {
		t.Fatalf("unexpected id line: %q", lines[0])
	}

	if !strings.HasPrefix(lines[1], "Key: ") {
		t.Fatalf("unexpected key line: %q", lines[1])
	}

	if lines[2] != "Summary: "+summary {
		t.Fatalf("unexpected summary line: %q", lines[2])
	}

	issueKey := strings.TrimPrefix(lines[1], "Key: ")
	t.Cleanup(func() {
		if err := deleteJiraTestIssue(jc, issueKey); err != nil {
			t.Errorf("failed to delete issue %s: %v", issueKey, err)
		}
	})

	issue, err := jc.GetIssue(issueKey)
	if err != nil {
		t.Fatal(err)
	}

	if issue.Fields.Summary != summary {
		t.Fatalf("expected summary %q, got %q", summary, issue.Fields.Summary)
	}
}

func TestIssueCreateRequiresProjectSelector(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "create", "--issue-type", "Task", "--summary", "summary")
	if err == nil {
		t.Fatal("expected create without project selector to fail")
	}

	if err.Error() != "exactly one of --project-key or --project-id is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueCreateRequiresIssueTypeSelector(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "create", "--project-key", "KAN", "--summary", "summary")
	if err == nil {
		t.Fatal("expected create without issue type selector to fail")
	}

	if err.Error() != "exactly one of --issue-type or --issue-type-id is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueCreateRequiresSummary(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "create", "--project-key", "KAN", "--issue-type", "Task")
	if err == nil {
		t.Fatal("expected create without summary to fail")
	}

	if err.Error() != "--summary is required and must not be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueCreateRejectsConflictingDescriptionFlags(t *testing.T) {
	_, err := runCommandWithInput(
		t,
		Environment{},
		"",
		"jira",
		"issue",
		"create",
		"--project-key",
		"KAN",
		"--issue-type",
		"Task",
		"--summary",
		"summary",
		"--description",
		"a",
		"--description-file",
		"description.txt",
	)
	if err == nil {
		t.Fatal("expected conflicting description flags to fail")
	}

	if err.Error() != "--description and --description-file are mutually exclusive" {
		t.Fatalf("unexpected error: %v", err)
	}
}
