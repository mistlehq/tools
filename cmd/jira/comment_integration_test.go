package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIssueCommentAddWithBody(t *testing.T) {
	commandResult := setupAndRunCommand(t, "jira", "issue", "comment", "add", "KAN-1", "--body", "comment from integration test")
	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")

	if len(lines) != 4 {
		t.Fatalf("expected 4 output lines, got %d: %q", len(lines), output)
	}

	if !strings.HasPrefix(lines[0], "Issue: ") {
		t.Fatal("expected Issue line")
	}

	if !strings.HasPrefix(lines[1], "Comment ID: ") {
		t.Fatal("expected Comment ID line")
	}

	if !strings.HasPrefix(lines[2], "Author: ") {
		t.Fatal("expected Author line")
	}

	if !strings.HasPrefix(lines[3], "Created: ") {
		t.Fatal("expected Created line")
	}
}

func TestIssueCommentAddWithBodyFile(t *testing.T) {
	tempDir := t.TempDir()
	commentFile := filepath.Join(tempDir, "comment.txt")
	if err := os.WriteFile(commentFile, []byte("comment from file"), 0o600); err != nil {
		t.Fatal(err)
	}

	commandResult := setupAndRunCommand(t, "jira", "issue", "comment", "add", "KAN-1", "--body-file", commentFile)
	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Comment ID: ") {
		t.Fatal("expected comment output to include Comment ID")
	}
}

func TestIssueCommentAddWithStdin(t *testing.T) {
	commandResult := setupAndRunCommandWithInput(t, "comment from stdin", "jira", "issue", "comment", "add", "KAN-1", "--body-file", "-")
	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Comment ID: ") {
		t.Fatal("expected comment output to include Comment ID")
	}
}

func TestIssueCommentAddRequiresSingleBodySource(t *testing.T) {
	_, err := runCommand(t, Environment{}, "jira", "issue", "comment", "add", "KAN-1")
	if err == nil {
		t.Fatal("expected missing body source to fail")
	}

	if err.Error() != "exactly one of --body or --body-file is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueCommentAddRejectsConflictingBodyFlags(t *testing.T) {
	_, err := runCommand(t, Environment{}, "jira", "issue", "comment", "add", "KAN-1", "--body", "a", "--body-file", "comment.txt")
	if err == nil {
		t.Fatal("expected conflicting body flags to fail")
	}

	if err.Error() != "--body and --body-file are mutually exclusive" {
		t.Fatalf("unexpected error: %v", err)
	}
}
