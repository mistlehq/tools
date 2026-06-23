package main

import (
	"encoding/json"
	"testing"
)

func TestGraphQLRequest(t *testing.T) {
	env := setupCommandEnvironment(t)
	queryFile := writeTempTextFile(t, `query ProductByID($id: ID!) { product(id: $id) { id handle } }`)
	variablesFile := writeTempTextFile(t, `{"id": "`+testProductID(t)+`"}`)

	commandResult, err := runCommandWithInput(t, env, "", "shopify", "graphql", "request", "--query-file", queryFile, "--variables-file", variablesFile)
	if err != nil {
		t.Fatal(err)
	}

	var result struct {
		Data struct {
			Product ShopifyProduct `json:"product"`
		} `json:"data"`
	}
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if result.Data.Product.ID != testProductID(t) {
		t.Fatalf("expected product %q, got %#v", testProductID(t), result.Data.Product)
	}
}
