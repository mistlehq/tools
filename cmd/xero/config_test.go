package main

import "testing"

func TestLoadConfigRequiresAPIBaseURL(t *testing.T) {
	_, err := loadConfig(Environment{})
	if err == nil {
		t.Fatal("expected missing API base URL to fail")
	}
	if err.Error() != "missing XERO_API_BASE_URL" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadConfigRejectsTrailingSlash(t *testing.T) {
	_, err := loadConfig(Environment{"XERO_API_BASE_URL": "https://api.xero.com/"})
	if err == nil {
		t.Fatal("expected trailing slash to fail")
	}
	if err.Error() != "XERO_API_BASE_URL must not end with '/'" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadConfigAcceptsAPIBaseURL(t *testing.T) {
	config, err := loadConfig(Environment{"XERO_API_BASE_URL": "https://api.xero.com"})
	if err != nil {
		t.Fatal(err)
	}
	if config.APIBaseURL != "https://api.xero.com" {
		t.Fatalf("unexpected API base URL: %s", config.APIBaseURL)
	}
}
