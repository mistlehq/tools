package main

import (
	"encoding/json"
	"testing"
)

func TestSearchAnalyticsQuery(t *testing.T) {
	env := setupCommandEnvironment(t)
	siteURL := testSiteURL(t)
	requestFile := writeTempJSONRequest(t, minimalSearchAnalyticsRequest())
	commandResult, err := runCommandWithInput(t, env, "", "gsc", "searchanalytics", "query", "--site-url", siteURL, "--request-file", requestFile)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
}
