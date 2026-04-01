package main

import (
	"strings"
	"testing"
)

func TestIssueTransitionList(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "transition", "list", issueKey)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")

	if strings.TrimSpace(lines[0]) == "" {
		t.Fatal("expected transition headers")
	}

	headerColumns := strings.Split(lines[0], "\t")
	if len(headerColumns) != 3 {
		t.Fatal("expected transition list header to have 3 columns")
	}

	if headerColumns[0] != "ID" || headerColumns[1] != "NAME" || headerColumns[2] != "TO STATUS" {
		t.Fatalf("unexpected transition list headers: %q", lines[0])
	}
}

func TestIssueTransitionByID(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	listResult, err := runCommandWithInput(t, env, "", "jira", "issue", "transition", "list", issueKey)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(listResult.stdout.String())
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		t.Skip("skipping: issue has no available transitions")
	}

	transitionColumns := strings.Split(lines[1], "\t")
	if len(transitionColumns) != 3 {
		t.Fatalf("expected transition row to have 3 columns, got %q", lines[1])
	}

	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "transition", issueKey, "--to-id", transitionColumns[0])
	if err != nil {
		t.Fatal(err)
	}

	transitionOutput := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(transitionOutput, "Transition: "+transitionColumns[1]) {
		t.Fatal("expected transition output to include transition name")
	}
}

func TestIssueTransitionByName(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	listResult, err := runCommandWithInput(t, env, "", "jira", "issue", "transition", "list", issueKey)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(listResult.stdout.String())
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		t.Skip("skipping: issue has no available transitions")
	}

	transitionColumns := strings.Split(lines[1], "\t")
	if len(transitionColumns) != 3 {
		t.Fatalf("expected transition row to have 3 columns, got %q", lines[1])
	}

	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "transition", issueKey, "--to", transitionColumns[1])
	if err != nil {
		t.Fatal(err)
	}

	transitionOutput := strings.TrimSpace(commandResult.stdout.String())

	if !strings.Contains(transitionOutput, "Transition: "+transitionColumns[1]) {
		t.Fatal("expected transition output to include transition name")
	}
}

func TestIssueTransitionRejectsUnknownName(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "transition", "KAN-1")
	if err == nil {
		t.Fatal("expected missing transition target to fail")
	}

	if err.Error() != "exactly one of --to or --to-id is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIssueTransitionRejectsUnknownTransitionName(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	_, err := runCommandWithInput(t, env, "", "jira", "issue", "transition", issueKey, "--to", "definitely not a transition")
	if err == nil {
		t.Fatal("expected unknown transition name to fail")
	}

	if !strings.Contains(err.Error(), `no transition named "definitely not a transition"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
