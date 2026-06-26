package main

import (
	"os"
	"strings"
	"testing"
)

func TestDiscordMessageLifecycle(t *testing.T) {
	env := setupCommandEnvironment(t)
	channelID := getRequiredEnv(t, "DISCORD_TEST_CHANNEL_ID")
	messageText := uniqueTestMessage("discord lifecycle")
	updatedText := messageText + " updated"

	created, err := runCommand(t, env, "discord", "messages", "send", "--channel", channelID, "--content", messageText)
	if err != nil {
		t.Fatal(err)
	}
	messageID := parseLineValue(t, created.stdout.String(), "id\t")

	listed, err := runCommand(t, env, "discord", "messages", "list", "--channel", channelID, "--limit", "10")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(listed.stdout.String(), messageText) {
		t.Fatalf("expected listed messages to include %q, got:\n%s", messageText, listed.stdout.String())
	}

	edited, err := runCommand(t, env, "discord", "messages", "edit", "--channel", channelID, "--message", messageID, "--content", updatedText)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(edited.stdout.String(), updatedText) {
		t.Fatalf("expected edited message output to include %q, got:\n%s", updatedText, edited.stdout.String())
	}

	if _, err := runCommand(t, env, "discord", "messages", "delete", "--channel", channelID, "--message", messageID); err != nil {
		t.Fatal(err)
	}
}

func TestDiscordChannelsList(t *testing.T) {
	env := setupCommandEnvironment(t)
	guildID := getRequiredEnv(t, "DISCORD_TEST_GUILD_ID")
	channelID := os.Getenv("DISCORD_TEST_CHANNEL_ID")

	result, err := runCommand(t, env, "discord", "channels", "list", "--guild", guildID)
	if err != nil {
		t.Fatal(err)
	}

	if channelID != "" && !strings.Contains(result.stdout.String(), channelID) {
		t.Fatalf("expected channels output to contain configured channel %q, got:\n%s", channelID, result.stdout.String())
	}
}
