package main

import (
	"strings"
	"testing"
)

func TestAuthTestRequiresSiteURL(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gsc", "auth", "test")
	if err == nil || !strings.Contains(err.Error(), "auth test requires --site-url") {
		t.Fatalf("expected missing site URL error, got %v", err)
	}
}

func TestSearchAnalyticsQueryRequiresRequestFile(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gsc", "searchanalytics", "query", "--site-url", "https://example.com/")
	if err == nil || !strings.Contains(err.Error(), "searchanalytics query requires --request-file") {
		t.Fatalf("expected missing request file error, got %v", err)
	}
}

func TestSitemapsGetRequiresFeedPath(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gsc", "sitemaps", "get", "--site-url", "https://example.com/")
	if err == nil || !strings.Contains(err.Error(), "sitemaps get requires --feed-path") {
		t.Fatalf("expected missing feed path error, got %v", err)
	}
}

func TestURLInspectionRequiresRequestFile(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gsc", "url-inspection", "inspect")
	if err == nil || !strings.Contains(err.Error(), "url-inspection inspect requires --request-file") {
		t.Fatalf("expected missing request file error, got %v", err)
	}
}

func TestRequestFileRejectsInvalidJSON(t *testing.T) {
	path := writeTempJSONRequest(t, "{")
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gsc", "searchanalytics", "query", "--site-url", "https://example.com/", "--request-file", path)
	if err == nil || !strings.Contains(err.Error(), "request file must contain valid JSON") {
		t.Fatalf("expected invalid JSON error, got %v", err)
	}
}

func TestEscapePathPartEscapesURLAndDomainProperties(t *testing.T) {
	testCases := map[string]string{
		"https://example.com/":        "https:%2F%2Fexample.com%2F",
		"sc-domain:example.com":       "sc-domain:example.com",
		"https://example.com/map.xml": "https:%2F%2Fexample.com%2Fmap.xml",
	}

	for input, want := range testCases {
		if got := escapePathPart(input); got != want {
			t.Fatalf("expected %q to escape to %q, got %q", input, want, got)
		}
	}
}
