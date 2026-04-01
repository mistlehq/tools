package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIssueUpdateSummary(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	summary := fmt.Sprintf("summary updated at %d", time.Now().UnixNano())
	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "update", issueKey, "--summary", summary)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Updated: summary") {
		t.Fatal("expected update output to mention summary")
	}

	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	issue, err := NewJiraClient(config).GetIssue(issueKey)
	if err != nil {
		t.Fatal(err)
	}

	if issue.Fields.Summary != summary {
		t.Fatalf("expected updated summary %q, got %q", summary, issue.Fields.Summary)
	}
}

func TestIssueUpdateDescription(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "update", issueKey, "--description", "description updated from integration test")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Updated: description") {
		t.Fatal("expected update output to mention description")
	}
}

func TestIssueUpdateDescriptionFile(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	tempDir := t.TempDir()
	descriptionFile := filepath.Join(tempDir, "description.txt")
	if err := os.WriteFile(descriptionFile, []byte("description from file"), 0o600); err != nil {
		t.Fatal(err)
	}

	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "update", issueKey, "--description-file", descriptionFile)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Updated: description") {
		t.Fatal("expected update output to mention description")
	}
}

func TestIssueUpdateSummaryAndDescription(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	summary := fmt.Sprintf("combined update at %d", time.Now().UnixNano())
	commandResult, err := runCommandWithInput(t, env, "combined description from stdin", "jira", "issue", "update", issueKey, "--summary", summary, "--description-file", "-")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Updated: summary, description") {
		t.Fatal("expected update output to mention summary and description")
	}
}

func TestIssueUpdateRequiresAtLeastOneField(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "update", "KAN-1")
	if err == nil {
		t.Fatal("expected update without fields to fail")
	}

	if err.Error() != "issue update requires at least one of --summary, --description, or --description-file" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueUpdateRejectsConflictingDescriptionFlags(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "update", "KAN-1", "--description", "a", "--description-file", "description.txt")
	if err == nil {
		t.Fatal("expected conflicting description flags to fail")
	}

	if err.Error() != "--description and --description-file are mutually exclusive" {
		t.Fatalf("unexpected error: %v", err)
	}
}
