package main

import (
	"strings"
	"testing"
)

func TestDiscordAuthTest(t *testing.T) {
	env := setupCommandEnvironment(t)

	result, err := runCommand(t, env, "discord", "auth", "test")
	if err != nil {
		t.Fatal(err)
	}

	output := result.stdout.String()
	if !strings.Contains(output, "id\t") {
		t.Fatalf("expected auth output to include an id, got:\n%s", output)
	}
}
