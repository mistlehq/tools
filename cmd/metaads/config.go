package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	GraphBaseURL string `json:"graphBaseUrl"`
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
	graphBaseURL := env["METAADS_GRAPH_BASE_URL"]
	if graphBaseURL == "" {
		return Config{}, fmt.Errorf("missing METAADS_GRAPH_BASE_URL")
	}
	if strings.HasSuffix(graphBaseURL, "/") {
		return Config{}, fmt.Errorf("METAADS_GRAPH_BASE_URL must not end with '/'")
	}
	return Config{GraphBaseURL: graphBaseURL}, nil
}
