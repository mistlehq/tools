package main

import (
	"strings"
	"testing"
)

func TestIssueAssignToMe(t *testing.T) {
	commandResult := setupAndRunCommand(t, "jira", "issue", "assign", "KAN-1", "--me")
	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")

	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d: %q", len(lines), output)
	}

	if !strings.HasPrefix(lines[0], "Issue: ") {
		t.Fatal("expected Issue line")
	}

	if !strings.HasPrefix(lines[1], "Assignee: ") {
		t.Fatal("expected Assignee line")
	}
}

func TestIssueAssignByAccountID(t *testing.T) {
	myselfResult := setupAndRunCommand(t, "jira", "auth", "whoami")
	output := strings.TrimSpace(myselfResult.stdout.String())
	lines := strings.Split(output, "\n")
	accountID := strings.TrimPrefix(lines[0], "Account ID: ")

	commandResult := setupAndRunCommand(t, "jira", "issue", "assign", "KAN-1", "--account-id", accountID)
	assignOutput := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(assignOutput, "Assignee: ") {
		t.Fatal("expected assignee output")
	}
}

func TestIssueAssignUnassigned(t *testing.T) {
	commandResult := setupAndRunCommand(t, "jira", "issue", "assign", "KAN-1", "--unassigned")
	output := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(output, "Assignee: Unassigned") {
		t.Fatal("expected issue to be unassigned")
	}
}

func TestIssueAssignRequiresExactlyOneTargetFlag(t *testing.T) {
	_, err := runCommand(t, Environment{}, "jira", "issue", "assign", "KAN-1")
	if err == nil {
		t.Fatal("expected issue assign without target flag to fail")
	}

	if err.Error() != "exactly one of --me, --account-id, or --unassigned is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueAssignRejectsConflictingTargetFlags(t *testing.T) {
	_, err := runCommand(t, Environment{}, "jira", "issue", "assign", "KAN-1", "--me", "--unassigned")
	if err == nil {
		t.Fatal("expected conflicting issue assign flags to fail")
	}

	if err.Error() != "exactly one of --me, --account-id, or --unassigned is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}
