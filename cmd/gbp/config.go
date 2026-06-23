package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AccountManagementBaseURL   string `json:"accountManagementBaseUrl"`
	BusinessInformationBaseURL string `json:"businessInformationBaseUrl"`
	PerformanceBaseURL         string `json:"performanceBaseUrl"`
	MyBusinessBaseURL          string `json:"myBusinessBaseUrl"`
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
	accountManagementBaseURL, err := requiredBaseURL(env, "GBP_ACCOUNT_MANAGEMENT_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	businessInformationBaseURL, err := requiredBaseURL(env, "GBP_BUSINESS_INFORMATION_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	performanceBaseURL, err := requiredBaseURL(env, "GBP_PERFORMANCE_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	myBusinessBaseURL, err := requiredBaseURL(env, "GBP_MYBUSINESS_BASE_URL")
	if err != nil {
		return Config{}, err
	}

	return Config{
		AccountManagementBaseURL:   accountManagementBaseURL,
		BusinessInformationBaseURL: businessInformationBaseURL,
		PerformanceBaseURL:         performanceBaseURL,
		MyBusinessBaseURL:          myBusinessBaseURL,
	}, nil
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
