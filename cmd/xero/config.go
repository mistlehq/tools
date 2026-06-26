package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	APIBaseURL string
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
	apiBaseURL, err := requiredBaseURL(env, "XERO_API_BASE_URL")
	if err != nil {
		return Config{}, err
	}

	return Config{APIBaseURL: apiBaseURL}, nil
}

func requiredBaseURL(env Environment, name string) (string, error) {
	value := env[name]
	if value == "" {
		return "", fmt.Errorf("missing %s", name)
	}
	if strings.HasSuffix(value, "/") {
		return "", fmt.Errorf("%s must not end with '/'", name)
	}
	return value, nil
}
