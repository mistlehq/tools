package main

import (
	"strings"
	"testing"
)

func TestLoadConfigRequiresSearchConsoleBaseURL(t *testing.T) {
	_, err := loadConfig(Environment{})
	if err == nil || !strings.Contains(err.Error(), "missing GSC_SEARCH_CONSOLE_BASE_URL") {
		t.Fatalf("expected missing Search Console base URL error, got %v", err)
	}
}

func TestLoadConfigRejectsTrailingSlash(t *testing.T) {
	_, err := loadConfig(Environment{
		"GSC_SEARCH_CONSOLE_BASE_URL": "https://searchconsole.googleapis.com/",
	})
	if err == nil || err.Error() != "GSC_SEARCH_CONSOLE_BASE_URL must not end with '/'" {
		t.Fatalf("expected trailing slash error, got %v", err)
	}
}

func TestLoadConfigReturnsBaseURL(t *testing.T) {
	config, err := loadConfig(Environment{
		"GSC_SEARCH_CONSOLE_BASE_URL": "https://searchconsole.googleapis.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if config.SearchConsoleBaseURL != "https://searchconsole.googleapis.com" {
		t.Fatalf("unexpected Search Console base URL: %q", config.SearchConsoleBaseURL)
	}
}
