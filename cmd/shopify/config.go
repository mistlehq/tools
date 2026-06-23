package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AdminBaseURL string `json:"adminBaseUrl"`
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
	adminBaseURL := env["SHOPIFY_ADMIN_BASE_URL"]
	if adminBaseURL == "" {
		return Config{}, fmt.Errorf("missing SHOPIFY_ADMIN_BASE_URL")
	}

	if strings.HasSuffix(adminBaseURL, "/") {
		return Config{}, fmt.Errorf("SHOPIFY_ADMIN_BASE_URL must not end with '/'")
	}

	return Config{
		AdminBaseURL: adminBaseURL,
	}, nil
}
