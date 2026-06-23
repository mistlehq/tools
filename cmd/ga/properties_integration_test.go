package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAccountSummariesList(t *testing.T) {
	env := setupCommandEnvironment(t)
	accountID := testAccountID(t)
	propertyID := testPropertyID(t)
	commandResult, err := runCommandWithInput(t, env, "", "ga", "account-summaries", "list", "--json")
	if err != nil {
		t.Fatal(err)
	}

	var list GAAccountSummariesList
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &list); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if len(list.AccountSummaries) == 0 {
		t.Fatal("expected at least one account summary")
	}

	accountSeen := false
	propertySeen := false
	for _, summary := range list.AccountSummaries {
		if summary.Account == accountID {
			accountSeen = true
		}
		for _, property := range summary.PropertySummaries {
			if property.Property == propertyID {
				propertySeen = true
			}
		}
	}
	if !accountSeen && !propertySeen {
		t.Fatalf("expected account summaries to include account %q or property %q", accountID, propertyID)
	}
}

func TestPropertiesGet(t *testing.T) {
	env := setupCommandEnvironment(t)
	propertyID := testPropertyID(t)
	commandResult, err := runCommandWithInput(t, env, "", "ga", "properties", "get", "--property", propertyID)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "Name: "+propertyID) {
		t.Fatalf("expected property output to mention %q, got %q", propertyID, output)
	}
}

func TestMetadataGet(t *testing.T) {
	env := setupCommandEnvironment(t)
	propertyID := testPropertyID(t)
	commandResult, err := runCommandWithInput(t, env, "", "ga", "metadata", "get", "--property", propertyID, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var metadata GAMetadata
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &metadata); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if len(metadata.Dimensions) == 0 {
		t.Fatal("expected metadata dimensions")
	}
	if len(metadata.Metrics) == 0 {
		t.Fatal("expected metadata metrics")
	}
}

func TestGoogleAdsLinksList(t *testing.T) {
	env := setupCommandEnvironment(t)
	propertyID := testPropertyID(t)
	_, err := runCommandWithInput(t, env, "", "ga", "google-ads-links", "list", "--property", propertyID, "--json")
	if err != nil {
		t.Fatal(err)
	}
}
