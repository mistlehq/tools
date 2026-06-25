package main

import "testing"

func TestLoadConfigRejectsMissingBaseURLs(t *testing.T) {
	_, err := loadConfig(Environment{})
	if err == nil {
		t.Fatal("expected missing base URLs to fail")
	}
}

func TestLoadConfigRejectsTrailingSlash(t *testing.T) {
	env := validUnitEnv()
	env["GWS_GMAIL_BASE_URL"] = "https://gmail.googleapis.com/gmail/v1/"
	_, err := loadConfig(env)
	if err == nil {
		t.Fatal("expected trailing slash to fail")
	}
}

func TestLoadConfigAcceptsBaseURLs(t *testing.T) {
	config, err := loadConfig(validUnitEnv())
	if err != nil {
		t.Fatal(err)
	}
	if config.DriveBaseURL != "https://www.googleapis.com/drive/v3" {
		t.Fatalf("unexpected drive base URL: %s", config.DriveBaseURL)
	}
	if config.PeopleBaseURL != "https://people.googleapis.com/v1" {
		t.Fatalf("unexpected people base URL: %s", config.PeopleBaseURL)
	}
}

func validUnitEnv() Environment {
	return Environment{
		"GWS_DRIVE_BASE_URL":    "https://www.googleapis.com/drive/v3",
		"GWS_SHEETS_BASE_URL":   "https://sheets.googleapis.com/v4",
		"GWS_DOCS_BASE_URL":     "https://docs.googleapis.com/v1",
		"GWS_SLIDES_BASE_URL":   "https://slides.googleapis.com/v1",
		"GWS_GMAIL_BASE_URL":    "https://gmail.googleapis.com/gmail/v1",
		"GWS_CALENDAR_BASE_URL": "https://www.googleapis.com/calendar/v3",
		"GWS_CHAT_BASE_URL":     "https://chat.googleapis.com/v1",
		"GWS_PEOPLE_BASE_URL":   "https://people.googleapis.com/v1",
	}
}
