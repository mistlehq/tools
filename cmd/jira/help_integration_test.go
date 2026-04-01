package main

import (
	"strings"
	"testing"
)

func TestTopLevelHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira auth help") {
		t.Fatal("expected top-level help to mention jira auth help")
	}

	if !strings.Contains(output, "jira issue help") {
		t.Fatal("expected top-level help to mention jira issue help")
	}
}

func TestAuthHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "auth", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira auth whoami") {
		t.Fatal("expected auth help to mention jira auth whoami")
	}
}

func TestProjectHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "project", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira project list") {
		t.Fatal("expected project help to mention jira project list")
	}
}

func TestIssueHelpListsNestedFamilies(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	expected := []string{
		"jira issue comment help",
		"jira issue assign help",
		"jira issue transition help",
		"jira issue update help",
		"jira issue editmeta help",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected issue help to mention %q", want)
		}
	}
}

func TestIssueCommentHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "comment", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira issue comment add <issue-key> --body <text>") {
		t.Fatal("expected issue comment help to mention --body usage")
	}

	if !strings.Contains(output, "jira issue comment add <issue-key> --body-file <path>") {
		t.Fatal("expected issue comment help to mention --body-file usage")
	}
}

func TestIssueAssignHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "assign", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira issue assign <issue-key> --me") {
		t.Fatal("expected issue assign help to mention --me")
	}

	if !strings.Contains(output, "jira issue assign <issue-key> --unassigned") {
		t.Fatal("expected issue assign help to mention --unassigned")
	}
}

func TestIssueTransitionHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "transition", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira issue transition list <issue-key>") {
		t.Fatal("expected issue transition help to mention list")
	}

	if !strings.Contains(output, "jira issue transition <issue-key> --to-id <transition-id>") {
		t.Fatal("expected issue transition help to mention --to-id")
	}
}

func TestIssueUpdateHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "update", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira issue update <issue-key> --summary <text>") {
		t.Fatal("expected issue update help to mention --summary")
	}

	if !strings.Contains(output, "jira issue update <issue-key> --description-file <path>") {
		t.Fatal("expected issue update help to mention --description-file")
	}
}

func TestIssueEditMetaHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "issue", "editmeta", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira issue editmeta <issue-key>") {
		t.Fatal("expected issue editmeta help to mention issue key usage")
	}
}
