package main

import "testing"

func TestAuthTest(t *testing.T) {
	env := setupCommandEnvironment(t)
	result, err := runCommandWithInput(t, env, "", "metaads", "auth", "test")
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	decodeCommandJSON(t, result, &out)
	if out["id"] == nil {
		t.Fatalf("expected /me id, got %#v", out)
	}
}

func TestRawGraphRequest(t *testing.T) {
	env := setupCommandEnvironment(t)
	result, err := runCommandWithInput(t, env, "", "metaads", "graph", "request", "--method", "GET", "--path", "/me", "--params", `{"fields":"id,name"}`)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	decodeCommandJSON(t, result, &out)
	if out["id"] == nil {
		t.Fatalf("expected /me id, got %#v", out)
	}
}

func TestNamedReadCommands(t *testing.T) {
	env := setupCommandEnvironment(t)
	accountID := testAdAccountID(t)

	accountsResult, err := runCommandWithInput(t, env, "", "metaads", "ad-accounts", "list", "--limit", "5")
	if err != nil {
		t.Fatal(err)
	}
	var accounts map[string]any
	decodeCommandJSON(t, accountsResult, &accounts)
	if accounts["data"] == nil {
		t.Fatalf("expected ad accounts data, got %#v", accounts)
	}

	accountResult, err := runCommandWithInput(t, env, "", "metaads", "ad-accounts", "get", "--id", accountID)
	if err != nil {
		t.Fatal(err)
	}
	var account map[string]any
	decodeCommandJSON(t, accountResult, &account)
	if account["id"] != accountID {
		t.Fatalf("expected ad account %q, got %#v", accountID, account)
	}

	campaignsResult, err := runCommandWithInput(t, env, "", "metaads", "campaigns", "search", "--account-id", accountID, "--limit", "5")
	if err != nil {
		t.Fatal(err)
	}
	var campaigns map[string]any
	decodeCommandJSON(t, campaignsResult, &campaigns)
	if campaigns["data"] == nil {
		t.Fatalf("expected campaigns data, got %#v", campaigns)
	}

	insightsResult, err := runCommandWithInput(t, env, "", "metaads", "insights", "get", "--id", accountID, "--fields", "impressions,spend", "--params", `{"date_preset":"last_7d"}`)
	if err != nil {
		t.Fatal(err)
	}
	var insights map[string]any
	decodeCommandJSON(t, insightsResult, &insights)
	if insights["data"] == nil {
		t.Fatalf("expected insights data, got %#v", insights)
	}
}

func TestCampaignGetWhenFixturePresent(t *testing.T) {
	env := setupCommandEnvironment(t)
	campaignID := testCampaignID(t)
	result, err := runCommandWithInput(t, env, "", "metaads", "campaigns", "get", "--id", campaignID)
	if err != nil {
		t.Fatal(err)
	}
	var campaign map[string]any
	decodeCommandJSON(t, result, &campaign)
	if campaign["id"] != campaignID {
		t.Fatalf("expected campaign %q, got %#v", campaignID, campaign)
	}
}
