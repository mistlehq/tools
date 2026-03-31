package main

import (
	"strings"
	"testing"
)

func TestAuthWhoAmI(t *testing.T) {
	commandResult := setupAndRunCommand(t, "jira", "auth", "whoami")
	stdout := commandResult.stdout

	if stdout.Len() == 0 {
		t.Fatal("expected stdout output")
	}

	output := strings.TrimSpace(stdout.String())
	lines := strings.Split(output, "\n")

	if len(lines) != 3 {
		t.Fatalf("expected 3 output lines, got %d: %q", len(lines), output)
	}

	if strings.TrimSpace(lines[0]) == "" {
		t.Fatal("expected Account ID line to be non-empty")
	}

	if strings.TrimSpace(lines[1]) == "" {
		t.Fatal("expected Display Name line to be non-empty")
	}

	if strings.TrimSpace(lines[2]) == "" {
		t.Fatal("expected Email line to be non-empty")
	}
}
