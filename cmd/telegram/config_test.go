package main

import (
	"strings"
	"testing"
)

func TestLoadConfigRequiresTelegramBaseURL(t *testing.T) {
	_, err := loadConfig(Environment{})
	if err == nil {
		t.Fatal("expected missing TELEGRAM_BASE_URL to fail")
	}
	if !strings.Contains(err.Error(), "missing TELEGRAM_BASE_URL") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadConfigRejectsTrailingSlash(t *testing.T) {
	_, err := loadConfig(Environment{"TELEGRAM_BASE_URL": "https://api.telegram.org/"})
	if err == nil {
		t.Fatal("expected trailing slash to fail")
	}
	if !strings.Contains(err.Error(), "must not end with '/'") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadConfigUsesTelegramBaseURL(t *testing.T) {
	config, err := loadConfig(Environment{"TELEGRAM_BASE_URL": "https://api.telegram.org"})
	if err != nil {
		t.Fatal(err)
	}
	if config.BaseURL != "https://api.telegram.org" {
		t.Fatalf("expected base URL to be preserved, got %q", config.BaseURL)
	}
}

func TestParseOptionalPositiveIntFlag(t *testing.T) {
	empty, err := parseOptionalPositiveIntFlag("thread", "")
	if err != nil {
		t.Fatal(err)
	}
	if empty != 0 {
		t.Fatalf("expected empty optional integer to be zero, got %d", empty)
	}

	parsed, err := parseOptionalPositiveIntFlag("thread", "42")
	if err != nil {
		t.Fatal(err)
	}
	if parsed != 42 {
		t.Fatalf("expected parsed thread ID 42, got %d", parsed)
	}

	_, err = parseOptionalPositiveIntFlag("thread", "0")
	if err == nil || !strings.Contains(err.Error(), "--thread must be a positive integer") {
		t.Fatalf("expected invalid optional integer error, got %v", err)
	}
}
