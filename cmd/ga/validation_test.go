package main

import (
	"strings"
	"testing"
)

func TestAuthTestRequiresProperty(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "ga", "auth", "test")
	if err == nil || !strings.Contains(err.Error(), "auth test requires --property") {
		t.Fatalf("expected missing property error, got %v", err)
	}
}

func TestReportRequiresRequestFile(t *testing.T) {
	_, err := runCommandWithInput(t, validUnitEnv(), "", "ga", "reports", "run", "--property", "properties/123")
	if err == nil || !strings.Contains(err.Error(), "reports run requires --request-file") {
		t.Fatalf("expected missing request file error, got %v", err)
	}
}

func TestReportRejectsInvalidRequestFileJSON(t *testing.T) {
	path := writeTempJSONRequest(t, "{")
	_, err := runCommandWithInput(t, validUnitEnv(), "", "ga", "reports", "run", "--property", "properties/123", "--request-file", path)
	if err == nil || !strings.Contains(err.Error(), "request file must contain valid JSON") {
		t.Fatalf("expected invalid JSON error, got %v", err)
	}
}

func TestPropertyResourceNameRequiresPropertiesPrefix(t *testing.T) {
	err := validateGAResourceName("property", "123", "properties/")
	if err == nil || !strings.Contains(err.Error(), "property must start with properties/") {
		t.Fatalf("expected property prefix error, got %v", err)
	}
}

func validUnitEnv() Environment {
	return Environment{
		"GA_ANALYTICS_DATA_BASE_URL":  "http://127.0.0.1",
		"GA_ANALYTICS_ADMIN_BASE_URL": "http://127.0.0.1",
	}
}
