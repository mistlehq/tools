package main

import (
	"strings"
	"testing"
)

func TestAccountsGetRequiresAccount(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gbp", "accounts", "get")
	if err == nil || !strings.Contains(err.Error(), "accounts get requires --account") {
		t.Fatalf("expected missing account error, got %v", err)
	}
}

func TestLocationsListRequiresReadMask(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gbp", "locations", "list", "--account", testAccount)
	if err == nil || !strings.Contains(err.Error(), "locations list requires --read-mask") {
		t.Fatalf("expected missing read mask error, got %v", err)
	}
}

func TestReviewsGetRequiresReview(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gbp", "reviews", "get", "--account", testAccount, "--location", testLocation)
	if err == nil || !strings.Contains(err.Error(), "reviews get requires --review") {
		t.Fatalf("expected missing review error, got %v", err)
	}
}

func TestPerformanceDailyMetricsRequiresRequestFile(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gbp", "performance", "daily-metrics", "--location", testLocation)
	if err == nil || !strings.Contains(err.Error(), "performance daily-metrics requires --request-file") {
		t.Fatalf("expected missing request file error, got %v", err)
	}
}

func TestRequestFileRejectsInvalidJSON(t *testing.T) {
	path := writeTempJSONRequest(t, "{")
	_, err := runCommandWithInput(t, validUnitEnv(), "", "gbp", "performance", "daily-metrics", "--location", testLocation, "--request-file", path)
	if err == nil || !strings.Contains(err.Error(), "request file must contain valid JSON") {
		t.Fatalf("expected invalid JSON error, got %v", err)
	}
}

func TestEscapeResourceNamePreservesResourceSlashes(t *testing.T) {
	if got := escapeResourceName("accounts/123/locations/456"); got != "accounts/123/locations/456" {
		t.Fatalf("expected resource slashes to be preserved, got %q", got)
	}
}

func TestQueryFromJSONFlattensNestedGoogleQueryShape(t *testing.T) {
	query, err := queryFromJSON([]byte(`{"dailyMetric":"WEBSITE_CLICKS","dailyRange":{"start_date":{"year":2026,"month":6,"day":1}}}`))
	if err != nil {
		t.Fatal(err)
	}
	if got := query.Get("dailyRange.start_date.year"); got != "2026" {
		t.Fatalf("expected flattened year query value, got %q", got)
	}
}
