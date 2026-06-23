package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAuthTest(t *testing.T) {
	env := setupCommandEnvironment(t)
	propertyID := testPropertyID(t)
	commandResult, err := runCommandWithInput(t, env, "", "ga", "auth", "test", "--property", propertyID)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	if !strings.Contains(output, "Property: "+propertyID) {
		t.Fatalf("expected auth output to mention property %q, got %q", propertyID, output)
	}
}

func TestAuthTestJSON(t *testing.T) {
	env := setupCommandEnvironment(t)
	propertyID := testPropertyID(t)
	commandResult, err := runCommandWithInput(t, env, "", "ga", "auth", "test", "--property", propertyID, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var property GAProperty
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &property); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if property.Name != propertyID {
		t.Fatalf("expected property %q, got %#v", propertyID, property)
	}
}
