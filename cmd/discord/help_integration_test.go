package main

import (
	"strings"
	"testing"
)

func TestRootHelpDescribesDiscordCommandFamilies(t *testing.T) {
	result, err := runCommand(t, Environment{}, "discord", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := result.stdout.String()
	for _, expected := range []string{
		"CLI for Discord.",
		"discord messages help",
		"discord roles help",
		"discord mcp help",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected root help to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestMCPServeHelpListsDiscordTools(t *testing.T) {
	result, err := runCommand(t, Environment{}, "discord", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := result.stdout.String()
	for _, expected := range []string{
		"discord_auth_test",
		"discord_messages_send",
		"discord_roles_create",
		"discord_members_remove_role",
		"discord_members_ban",
		"discord_members_unban",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected MCP help to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestReactionsHelpDocumentsJSONOutput(t *testing.T) {
	result, err := runCommand(t, Environment{}, "discord", "reactions", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := result.stdout.String()
	for _, expected := range []string{
		"discord reactions add --channel <channel-id> --message <message-id> --emoji <emoji> [--json]",
		"discord reactions remove --channel <channel-id> --message <message-id> --emoji <emoji> [--json]",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected reactions help to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestLeafCommandsAcceptHelpWithoutConfig(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{name: "auth test", args: []string{"discord", "auth", "test", "--help"}, expected: "discord auth test [--json]"},
		{name: "guilds list", args: []string{"discord", "guilds", "list", "--help"}, expected: "discord guilds list [--json]"},
		{name: "messages send", args: []string{"discord", "messages", "send", "--help"}, expected: "discord messages send --channel <channel-id> --content <text> [--json]"},
		{name: "members ban", args: []string{"discord", "members", "ban", "--help"}, expected: "discord members ban --guild <guild-id> --user <user-id> [--json]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := runCommand(t, Environment{}, tc.args...)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(result.stdout.String(), tc.expected) {
				t.Fatalf("expected leaf help to contain %q, got:\n%s", tc.expected, result.stdout.String())
			}
		})
	}
}
