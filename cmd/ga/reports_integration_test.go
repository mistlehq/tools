package main

import (
	"encoding/json"
	"testing"
)

func TestCompatibilityCheck(t *testing.T) {
	env := setupCommandEnvironment(t)
	propertyID := testPropertyID(t)
	requestFile := writeTempJSONRequest(t, minimalCompatibilityRequest())
	commandResult, err := runCommandWithInput(t, env, "", "ga", "compatibility", "check", "--property", propertyID, "--request-file", requestFile)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected compatibility response fields")
	}
}

func TestReportsRun(t *testing.T) {
	env := setupCommandEnvironment(t)
	propertyID := testPropertyID(t)
	requestFile := writeTempJSONRequest(t, minimalRunReportRequest())
	commandResult, err := runCommandWithInput(t, env, "", "ga", "reports", "run", "--property", propertyID, "--request-file", requestFile)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if result["kind"] != "analyticsData#runReport" {
		t.Fatalf("expected run report response kind analyticsData#runReport, got %#v", result)
	}
}

func TestReportsRealtime(t *testing.T) {
	env := setupCommandEnvironment(t)
	propertyID := testPropertyID(t)
	requestFile := writeTempJSONRequest(t, minimalRealtimeReportRequest())
	commandResult, err := runCommandWithInput(t, env, "", "ga", "reports", "realtime", "--property", propertyID, "--request-file", requestFile)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if result["kind"] != "analyticsData#runRealtimeReport" {
		t.Fatalf("expected realtime report response kind analyticsData#runRealtimeReport, got %#v", result)
	}
}

func TestReportsFunnel(t *testing.T) {
	env := setupCommandEnvironment(t)
	propertyID := testPropertyID(t)
	requestFile := writeTempJSONRequest(t, minimalFunnelReportRequest())
	commandResult, err := runCommandWithInput(t, env, "", "ga", "reports", "funnel", "--property", propertyID, "--request-file", requestFile)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if _, ok := result["funnelTable"]; !ok {
		t.Fatalf("expected funnel report response to include funnelTable, got %#v", result)
	}
}
