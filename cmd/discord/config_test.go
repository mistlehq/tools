package main

import (
	"testing"
)

func TestLoadConfigRequiresDiscordBaseURL(t *testing.T) {
	_, err := loadConfig(Environment{})
	if err == nil {
		t.Fatal("expected missing DISCORD_BASE_URL to fail")
	}
}

func TestLoadConfigRejectsTrailingSlash(t *testing.T) {
	_, err := loadConfig(Environment{"DISCORD_BASE_URL": "https://discord.com/api/v10/"})
	if err == nil {
		t.Fatal("expected trailing slash to fail")
	}
}

func TestLoadConfigReturnsDiscordBaseURL(t *testing.T) {
	config, err := loadConfig(Environment{"DISCORD_BASE_URL": "https://discord.com/api/v10"})
	if err != nil {
		t.Fatal(err)
	}
	if config.BaseURL != "https://discord.com/api/v10" {
		t.Fatalf("unexpected base URL: %s", config.BaseURL)
	}
}

func TestNoArgumentCommandsRejectUnexpectedArguments(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{name: "auth test positional", args: []string{"discord", "auth", "test", "typo"}, expectedError: "unexpected positional argument: typo"},
		{name: "guilds list unsupported flag", args: []string{"discord", "guilds", "list", "--guild", "G123"}, expectedError: "unsupported flag: --guild"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := runCommand(t, Environment{}, tc.args...)
			if err == nil {
				t.Fatal("expected unexpected argument to fail")
			}
			if err.Error() != tc.expectedError {
				t.Fatalf("expected error %q, got %q", tc.expectedError, err.Error())
			}
		})
	}
}
