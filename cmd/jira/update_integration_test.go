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
	summary := fmt.Sprintf("summary updated at %d", time.Now().UnixNano())
	commandResult := setupAndRunCommand(t, "jira", "issue", "update", "KAN-1", "--summary", summary)
	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Updated: summary") {
		t.Fatal("expected update output to mention summary")
	}

	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	issue, err := NewJiraClient(config).GetIssue("KAN-1")
	if err != nil {
		t.Fatal(err)
	}

	if issue.Fields.Summary != summary {
		t.Fatalf("expected updated summary %q, got %q", summary, issue.Fields.Summary)
	}
}

func TestIssueUpdateDescription(t *testing.T) {
	commandResult := setupAndRunCommand(t, "jira", "issue", "update", "KAN-1", "--description", "description updated from integration test")
	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Updated: description") {
		t.Fatal("expected update output to mention description")
	}
}

func TestIssueUpdateDescriptionFile(t *testing.T) {
	tempDir := t.TempDir()
	descriptionFile := filepath.Join(tempDir, "description.txt")
	if err := os.WriteFile(descriptionFile, []byte("description from file"), 0o600); err != nil {
		t.Fatal(err)
	}

	commandResult := setupAndRunCommand(t, "jira", "issue", "update", "KAN-1", "--description-file", descriptionFile)
	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Updated: description") {
		t.Fatal("expected update output to mention description")
	}
}

func TestIssueUpdateSummaryAndDescription(t *testing.T) {
	summary := fmt.Sprintf("combined update at %d", time.Now().UnixNano())
	commandResult := setupAndRunCommandWithInput(t, "combined description from stdin", "jira", "issue", "update", "KAN-1", "--summary", summary, "--description-file", "-")
	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Updated: summary, description") {
		t.Fatal("expected update output to mention summary and description")
	}
}

func TestIssueUpdateRequiresAtLeastOneField(t *testing.T) {
	_, err := runCommand(t, Environment{}, "jira", "issue", "update", "KAN-1")
	if err == nil {
		t.Fatal("expected update without fields to fail")
	}

	if err.Error() != "issue update requires at least one of --summary, --description, or --description-file" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueUpdateRejectsConflictingDescriptionFlags(t *testing.T) {
	_, err := runCommand(t, Environment{}, "jira", "issue", "update", "KAN-1", "--description", "a", "--description-file", "description.txt")
	if err == nil {
		t.Fatal("expected conflicting description flags to fail")
	}

	if err.Error() != "--description and --description-file are mutually exclusive" {
		t.Fatalf("unexpected error: %v", err)
	}
}
