package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAuthTest(t *testing.T) {
	env := setupCommandEnvironment(t)
	siteURL := testSiteURL(t)
	commandResult, err := runCommandWithInput(t, env, "", "gsc", "auth", "test", "--site-url", siteURL)
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	if !strings.Contains(output, "Site URL: "+siteURL) {
		t.Fatalf("expected auth output to include site URL %q, got %q", siteURL, output)
	}
}

func TestSitesList(t *testing.T) {
	env := setupCommandEnvironment(t)
	siteURL := testSiteURL(t)
	commandResult, err := runCommandWithInput(t, env, "", "gsc", "sites", "list", "--json")
	if err != nil {
		t.Fatal(err)
	}

	var result GSCSitesList
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if !containsSite(result.SiteEntry, siteURL) {
		t.Fatalf("expected sites list to include %q, got %#v", siteURL, result.SiteEntry)
	}
}

func TestSitesGet(t *testing.T) {
	env := setupCommandEnvironment(t)
	siteURL := testSiteURL(t)
	commandResult, err := runCommandWithInput(t, env, "", "gsc", "sites", "get", "--site-url", siteURL, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var result GSCSite
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if result.SiteURL != siteURL {
		t.Fatalf("expected site URL %q, got %#v", siteURL, result)
	}
}
