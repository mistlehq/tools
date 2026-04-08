package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BaseURL string `json:"baseUrl"`
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
	baseURL := env["SLACK_BASE_URL"]
	if baseURL == "" {
		return Config{}, fmt.Errorf("missing SLACK_BASE_URL")
	}

	if strings.HasSuffix(baseURL, "/") {
		return Config{}, fmt.Errorf("SLACK_BASE_URL must not end with '/'")
	}

	return Config{
		BaseURL: baseURL,
	}, nil
}
