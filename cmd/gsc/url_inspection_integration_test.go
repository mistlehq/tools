package main

import (
	"encoding/json"
	"testing"
)

func TestURLInspectionInspect(t *testing.T) {
	env := setupCommandEnvironment(t)
	siteURL := testSiteURL(t)
	inspectionURL := testInspectionURL(t)
	requestFile := writeTempJSONRequest(t, minimalURLInspectionRequest(siteURL, inspectionURL))
	commandResult, err := runCommandWithInput(t, env, "", "gsc", "url-inspection", "inspect", "--request-file", requestFile)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if _, ok := result["inspectionResult"]; !ok {
		t.Fatalf("expected URL inspection response to include inspectionResult, got %#v", result)
	}
}
