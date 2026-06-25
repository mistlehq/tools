package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DriveBaseURL    string `json:"driveBaseUrl"`
	SheetsBaseURL   string `json:"sheetsBaseUrl"`
	DocsBaseURL     string `json:"docsBaseUrl"`
	SlidesBaseURL   string `json:"slidesBaseUrl"`
	GmailBaseURL    string `json:"gmailBaseUrl"`
	CalendarBaseURL string `json:"calendarBaseUrl"`
	ChatBaseURL     string `json:"chatBaseUrl"`
	PeopleBaseURL   string `json:"peopleBaseUrl"`
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
	driveBaseURL, err := requiredBaseURL(env, "GWS_DRIVE_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	sheetsBaseURL, err := requiredBaseURL(env, "GWS_SHEETS_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	docsBaseURL, err := requiredBaseURL(env, "GWS_DOCS_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	slidesBaseURL, err := requiredBaseURL(env, "GWS_SLIDES_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	gmailBaseURL, err := requiredBaseURL(env, "GWS_GMAIL_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	calendarBaseURL, err := requiredBaseURL(env, "GWS_CALENDAR_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	chatBaseURL, err := requiredBaseURL(env, "GWS_CHAT_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	peopleBaseURL, err := requiredBaseURL(env, "GWS_PEOPLE_BASE_URL")
	if err != nil {
		return Config{}, err
	}
	return Config{
		DriveBaseURL:    driveBaseURL,
		SheetsBaseURL:   sheetsBaseURL,
		DocsBaseURL:     docsBaseURL,
		SlidesBaseURL:   slidesBaseURL,
		GmailBaseURL:    gmailBaseURL,
		CalendarBaseURL: calendarBaseURL,
		ChatBaseURL:     chatBaseURL,
		PeopleBaseURL:   peopleBaseURL,
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
