package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAuthTest(t *testing.T) {
	env := setupCommandEnvironment(t)
	commandResult, err := runCommandWithInput(t, env, "", "slack", "auth", "test")
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	lines := strings.Split(output, "\n")
	if len(lines) != 6 {
		t.Fatalf("expected 6 output lines, got %d: %q", len(lines), output)
	}

	expectedPrefixes := []string{
		"URL: ",
		"Team: ",
		"Team ID: ",
		"User: ",
		"User ID: ",
		"Bot ID: ",
	}

	for index, want := range expectedPrefixes {
		if !strings.HasPrefix(lines[index], want) {
			t.Fatalf("expected line %d to start with %q, got %q", index+1, want, lines[index])
		}
	}
}

func TestAuthTestJSON(t *testing.T) {
	env := setupCommandEnvironment(t)
	commandResult, err := runCommandWithInput(t, env, "", "slack", "auth", "test", "--json")
	if err != nil {
		t.Fatal(err)
	}

	var authTest SlackAuthTest
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &authTest); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !authTest.OK {
		t.Fatal("expected ok=true in JSON output")
	}
}
