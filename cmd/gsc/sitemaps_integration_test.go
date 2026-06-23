package main

import (
	"encoding/json"
	"testing"
)

func TestSitemapsList(t *testing.T) {
	env := setupCommandEnvironment(t)
	siteURL := testSiteURL(t)
	commandResult, err := runCommandWithInput(t, env, "", "gsc", "sitemaps", "list", "--site-url", siteURL, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var result GSCSitemapsList
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
}

func TestSitemapsGet(t *testing.T) {
	env := setupCommandEnvironment(t)
	siteURL := testSiteURL(t)
	sitemapURL := testSitemapURL(t)
	commandResult, err := runCommandWithInput(t, env, "", "gsc", "sitemaps", "get", "--site-url", siteURL, "--feed-path", sitemapURL, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var result GSCSitemap
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if result.Path != sitemapURL {
		t.Fatalf("expected sitemap path %q, got %#v", sitemapURL, result)
	}
}
