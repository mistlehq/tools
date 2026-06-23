package main

import (
	"strings"
	"testing"
)

func TestLoadConfigRequiresAnalyticsDataBaseURL(t *testing.T) {
	_, err := loadConfig(Environment{
		"GA_ANALYTICS_ADMIN_BASE_URL": "https://analyticsadmin.googleapis.com",
	})
	if err == nil || !strings.Contains(err.Error(), "missing GA_ANALYTICS_DATA_BASE_URL") {
		t.Fatalf("expected missing data base URL error, got %v", err)
	}
}

func TestLoadConfigRequiresAnalyticsAdminBaseURL(t *testing.T) {
	_, err := loadConfig(Environment{
		"GA_ANALYTICS_DATA_BASE_URL": "https://analyticsdata.googleapis.com",
	})
	if err == nil || !strings.Contains(err.Error(), "missing GA_ANALYTICS_ADMIN_BASE_URL") {
		t.Fatalf("expected missing admin base URL error, got %v", err)
	}
}

func TestLoadConfigRejectsTrailingSlash(t *testing.T) {
	testCases := []struct {
		name string
		env  Environment
		want string
	}{
		{
			name: "data URL",
			env: Environment{
				"GA_ANALYTICS_DATA_BASE_URL":  "https://analyticsdata.googleapis.com/",
				"GA_ANALYTICS_ADMIN_BASE_URL": "https://analyticsadmin.googleapis.com",
			},
			want: "GA_ANALYTICS_DATA_BASE_URL must not end with '/'",
		},
		{
			name: "admin URL",
			env: Environment{
				"GA_ANALYTICS_DATA_BASE_URL":  "https://analyticsdata.googleapis.com",
				"GA_ANALYTICS_ADMIN_BASE_URL": "https://analyticsadmin.googleapis.com/",
			},
			want: "GA_ANALYTICS_ADMIN_BASE_URL must not end with '/'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := loadConfig(tc.env)
			if err == nil || err.Error() != tc.want {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
}

func TestLoadConfigReturnsBaseURLs(t *testing.T) {
	config, err := loadConfig(Environment{
		"GA_ANALYTICS_DATA_BASE_URL":  "https://analyticsdata.googleapis.com",
		"GA_ANALYTICS_ADMIN_BASE_URL": "https://analyticsadmin.googleapis.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if config.AnalyticsDataBaseURL != "https://analyticsdata.googleapis.com" {
		t.Fatalf("unexpected data base URL: %q", config.AnalyticsDataBaseURL)
	}
	if config.AnalyticsAdminBaseURL != "https://analyticsadmin.googleapis.com" {
		t.Fatalf("unexpected admin base URL: %q", config.AnalyticsAdminBaseURL)
	}
}
