package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AnalyticsDataBaseURL  string `json:"analyticsDataBaseUrl"`
	AnalyticsAdminBaseURL string `json:"analyticsAdminBaseUrl"`
}

type Environment map[string]string

func loadEnvironment() Environment {
	env := make(Environment)

	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	return env
}

func loadConfig(env Environment) (Config, error) {
	dataBaseURL := env["GA_ANALYTICS_DATA_BASE_URL"]
	if dataBaseURL == "" {
		return Config{}, fmt.Errorf("missing GA_ANALYTICS_DATA_BASE_URL")
	}
	if strings.HasSuffix(dataBaseURL, "/") {
		return Config{}, fmt.Errorf("GA_ANALYTICS_DATA_BASE_URL must not end with '/'")
	}

	adminBaseURL := env["GA_ANALYTICS_ADMIN_BASE_URL"]
	if adminBaseURL == "" {
		return Config{}, fmt.Errorf("missing GA_ANALYTICS_ADMIN_BASE_URL")
	}
	if strings.HasSuffix(adminBaseURL, "/") {
		return Config{}, fmt.Errorf("GA_ANALYTICS_ADMIN_BASE_URL must not end with '/'")
	}

	return Config{
		AnalyticsDataBaseURL:  dataBaseURL,
		AnalyticsAdminBaseURL: adminBaseURL,
	}, nil
}
