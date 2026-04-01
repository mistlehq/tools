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
	expected := []string{
		"Thin Jira Cloud CLI for shells, scripts, and agent-driven workflows.",
		"Common Starting Points:",
		"jira issue get PROJ-123",
		"Leaf commands also accept --help",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected top-level help to mention %q", want)
		}
	}
}

func TestTopLevelDashDashHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira --help") {
		t.Fatal("expected top-level --help to mention jira --help")
	}
}

func TestLeafCommandsAcceptDashDashHelp(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "auth whoami",
			args: []string{"jira", "auth", "whoami", "--help"},
			want: "Show the Jira account behind the current auth context.",
		},
		{
			name: "project list",
			args: []string{"jira", "project", "list", "--help"},
			want: "List Jira projects visible to the current caller.",
		},
		{
			name: "issue get",
			args: []string{"jira", "issue", "get", "--help"},
			want: "Fetch a single Jira issue by key.",
		},
		{
			name: "issue search",
			args: []string{"jira", "issue", "search", "--help"},
			want: "Search Jira issues with a JQL query.",
		},
		{
			name: "issue comment add",
			args: []string{"jira", "issue", "comment", "add", "--help"},
			want: "Exactly one of --body or --body-file is required.",
		},
		{
			name: "issue assign",
			args: []string{"jira", "issue", "assign", "--help"},
			want: "Exactly one of --me, --account-id, or --unassigned is required.",
		},
		{
			name: "issue transition list",
			args: []string{"jira", "issue", "transition", "list", "--help"},
			want: "List the workflow transitions currently available for an issue.",
		},
		{
			name: "issue update",
			args: []string{"jira", "issue", "update", "--help"},
			want: "Provide at least one of --summary, --description, or --description-file.",
		},
		{
			name: "issue editmeta",
			args: []string{"jira", "issue", "editmeta", "--help"},
			want: "Show edit metadata for a Jira issue.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			commandResult, err := runCommandWithInput(t, Environment{}, "", tc.args...)
			if err != nil {
				t.Fatal(err)
			}

			output := strings.TrimSpace(commandResult.stdout.String())
			if !strings.Contains(output, tc.want) {
				t.Fatalf("expected %q help to mention %q", tc.name, tc.want)
			}
		})
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

func TestAuthHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "auth", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira auth whoami --help") {
		t.Fatal("expected auth help to mention jira auth whoami --help")
	}
}

func TestProjectHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "project", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "jira project list --help") {
		t.Fatal("expected project help to mention jira project list --help")
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
