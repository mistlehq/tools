package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIssueCommentAddWithBody(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "comment", "add", issueKey, "--body", "comment from integration test")
	if err != nil {
		t.Fatal(err)
	}

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
	env, issueKey := setupIsolatedIssue(t)
	tempDir := t.TempDir()
	commentFile := filepath.Join(tempDir, "comment.txt")
	if err := os.WriteFile(commentFile, []byte("comment from file"), 0o600); err != nil {
		t.Fatal(err)
	}

	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "comment", "add", issueKey, "--body-file", commentFile)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Comment ID: ") {
		t.Fatal("expected comment output to include Comment ID")
	}
}

func TestIssueCommentAddWithStdin(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	commandResult, err := runCommandWithInput(t, env, "comment from stdin", "jira", "issue", "comment", "add", issueKey, "--body-file", "-")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Comment ID: ") {
		t.Fatal("expected comment output to include Comment ID")
	}
}

func TestIssueCommentAddRequiresSingleBodySource(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "comment", "add", "KAN-1")
	if err == nil {
		t.Fatal("expected missing body source to fail")
	}

	if err.Error() != "exactly one of --body or --body-file is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueCommentAddRejectsConflictingBodyFlags(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "comment", "add", "KAN-1", "--body", "a", "--body-file", "comment.txt")
	if err == nil {
		t.Fatal("expected conflicting body flags to fail")
	}

	if err.Error() != "--body and --body-file are mutually exclusive" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueCommentDelete(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	jc := NewJiraClient(config)
	comment, err := jc.AddIssueComment(issueKey, AddCommentInput{
		Body: "comment slated for deletion",
	})
	if err != nil {
		t.Fatal(err)
	}

	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "comment", "delete", issueKey, comment.ID)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 output lines, got %d: %q", len(lines), output)
	}

	if lines[0] != "Issue: "+issueKey {
		t.Fatalf("unexpected issue line: %q", lines[0])
	}

	if lines[1] != "Comment ID: "+comment.ID {
		t.Fatalf("unexpected comment line: %q", lines[1])
	}

	if lines[2] != "Deleted: true" {
		t.Fatalf("unexpected deleted line: %q", lines[2])
	}

	if _, err := jc.get(fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueKey, comment.ID)); err == nil {
		t.Fatalf("expected deleted comment %s on %s to be unavailable", comment.ID, issueKey)
	}
}

func TestIssueCommentDeleteRequiresIssueKeyAndCommentID(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "comment", "delete", "KAN-1")
	if err == nil {
		t.Fatal("expected missing comment id to fail")
	}

	if err.Error() != "issue comment delete expects exactly 2 positional arguments" {
		t.Fatalf("unexpected error: %v", err)
	}
}
