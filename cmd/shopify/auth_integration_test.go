package main

import (
	"strings"
	"testing"
)

func TestAuthTest(t *testing.T) {
	env := setupCommandEnvironment(t)
	commandResult, err := runCommandWithInput(t, env, "", "shopify", "auth", "test")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	for _, want := range []string{"Name:", "Myshopify domain:"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected auth output to mention %q, got %q", want, output)
		}
	}
}

func TestAuthTestJSON(t *testing.T) {
	env := setupCommandEnvironment(t)
	commandResult, err := runCommandWithInput(t, env, "", "shopify", "auth", "test", "--json")
	if err != nil {
		t.Fatal(err)
	}

	var shop ShopifyShop
	decodeCommandJSON(t, commandResult, &shop)
	if shop.ID == "" || shop.MyshopifyDomain == "" {
		t.Fatalf("expected shop identity fields, got %#v", shop)
	}
}
