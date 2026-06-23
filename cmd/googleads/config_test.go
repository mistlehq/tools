package main

import "testing"

func TestLoadConfigRejectsMissingBaseURL(t *testing.T) {
	_, err := loadConfig(Environment{})
	if err == nil {
		t.Fatal("expected missing GOOGLEADS_BASE_URL to fail")
	}
}

func TestLoadConfigRejectsTrailingSlash(t *testing.T) {
	_, err := loadConfig(Environment{"GOOGLEADS_BASE_URL": "https://googleads.googleapis.com/v24/"})
	if err == nil {
		t.Fatal("expected trailing slash to fail")
	}
}

func TestLoadConfigAcceptsBaseURL(t *testing.T) {
	config, err := loadConfig(Environment{"GOOGLEADS_BASE_URL": "https://googleads.googleapis.com/v24"})
	if err != nil {
		t.Fatal(err)
	}
	if config.BaseURL != "https://googleads.googleapis.com/v24" {
		t.Fatalf("unexpected base URL: %s", config.BaseURL)
	}
}
