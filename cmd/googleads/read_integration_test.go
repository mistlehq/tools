package main

import "testing"

func TestAuthTest(t *testing.T) {
	env := setupCommandEnvironment(t)
	result, err := runCommandWithInput(t, env, "", "googleads", "auth", "test")
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	decodeCommandJSON(t, result, &out)
	if out["resourceNames"] == nil {
		t.Fatalf("expected accessible customer resourceNames, got %#v", out)
	}
}

func TestRawRequest(t *testing.T) {
	env := setupCommandEnvironment(t)
	result, err := runCommandWithInput(t, env, "", "googleads", "request", "--method", "GET", "--path", "/customers:listAccessibleCustomers")
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	decodeCommandJSON(t, result, &out)
	if out["resourceNames"] == nil {
		t.Fatalf("expected accessible customer resourceNames, got %#v", out)
	}
}

func TestNamedReadCommands(t *testing.T) {
	env := setupCommandEnvironment(t)
	customerID := testCustomerID(t)
	query := "SELECT customer.id, customer.descriptive_name FROM customer LIMIT 1"

	searchResult, err := runCommandWithInput(t, env, "", "googleads", "gaql", "search", "--customer-id", customerID, "--query", query, "--page-size", "10")
	if err != nil {
		t.Fatal(err)
	}
	var search map[string]any
	decodeCommandJSON(t, searchResult, &search)
	if search["results"] == nil {
		t.Fatalf("expected GAQL search results, got %#v", search)
	}

	fieldsResult, err := runCommandWithInput(t, env, "", "googleads", "fields", "search", "--query", `SELECT name, category, data_type WHERE name = "campaign.id"`)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	decodeCommandJSON(t, fieldsResult, &fields)
	if fields["results"] == nil {
		t.Fatalf("expected field search results, got %#v", fields)
	}
}
