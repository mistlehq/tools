package main

import (
	"strings"
	"testing"
)

func TestLoadConfigRequiresAllBaseURLs(t *testing.T) {
	names := []string{
		"GBP_ACCOUNT_MANAGEMENT_BASE_URL",
		"GBP_BUSINESS_INFORMATION_BASE_URL",
		"GBP_PERFORMANCE_BASE_URL",
		"GBP_MYBUSINESS_BASE_URL",
	}

	for _, missingName := range names {
		t.Run(missingName, func(t *testing.T) {
			env := validUnitEnv()
			delete(env, missingName)
			_, err := loadConfig(env)
			if err == nil || !strings.Contains(err.Error(), "missing "+missingName) {
				t.Fatalf("expected missing %s error, got %v", missingName, err)
			}
		})
	}
}

func TestLoadConfigRejectsTrailingSlash(t *testing.T) {
	env := validUnitEnv()
	env["GBP_MYBUSINESS_BASE_URL"] = "http://127.0.0.1/"

	_, err := loadConfig(env)
	if err == nil || !strings.Contains(err.Error(), "GBP_MYBUSINESS_BASE_URL must not end with '/'") {
		t.Fatalf("expected trailing slash error, got %v", err)
	}
}
