package main

import "testing"

func TestRawRequestValidation(t *testing.T) {
	_, err := parseRawRequestArgs([]string{"--api", "drive", "--method", "TRACE", "--path", "/files"})
	if err != nil {
		t.Fatal(err)
	}
	client := NewGWSClient(validUnitConfig())
	if _, err := client.Request(GWSRequest{API: "drive", Method: "TRACE", Path: "/files"}); err == nil {
		t.Fatal("expected unsupported method to fail")
	}
	if _, err := client.Request(GWSRequest{API: "drive", Method: "GET"}); err == nil {
		t.Fatal("expected missing path to fail")
	}
	if _, err := client.Request(GWSRequest{API: "drive", Method: "GET", Path: "files"}); err == nil {
		t.Fatal("expected relative path to fail")
	}
	if _, err := client.Request(GWSRequest{API: "unsupported", Method: "GET", Path: "/events"}); err == nil {
		t.Fatal("expected unsupported api to fail")
	}
}

func TestRequestFileValidation(t *testing.T) {
	path := writeTempJSONRequest(t, `[]`)
	_, _, err := parseRequestFileAndParamsArgs([]string{"--request-file", path}, "files create", nil)
	if err == nil {
		t.Fatal("expected array request file to fail")
	}

	invalid := writeTempJSONRequest(t, `{`)
	_, _, err = parseRequestFileAndParamsArgs([]string{"--request-file", invalid}, "files create", nil)
	if err == nil {
		t.Fatal("expected invalid request file JSON to fail")
	}
}

func TestRequiredCommandArguments(t *testing.T) {
	if _, _, err := parseIDAndParamsArgs(nil, "files get", "file-id", nil); err == nil {
		t.Fatal("expected missing file-id to fail")
	}
	if _, _, _, err := parseIDRequestFileAndParamsArgs(nil, "files update", "file-id", nil); err == nil {
		t.Fatal("expected missing file-id/request-file to fail")
	}
	if _, _, _, _, err := parseSpreadsheetValuesUpdateArgs(nil, "values update"); err == nil {
		t.Fatal("expected missing spreadsheet values args to fail")
	}
	if _, _, _, err := parseUserIDChildIDAndParamsArgs(nil, "messages get", "message-id", nil); err == nil {
		t.Fatal("expected missing gmail ids to fail")
	}
	if _, _, _, err := parseTwoIDAndParamsArgs(nil, "events get", "calendar-id", "event-id", nil); err == nil {
		t.Fatal("expected missing calendar ids to fail")
	}
	if _, err := parseRequiredQueryMaskArgs(nil, "search-contacts", "read-mask", nil); err == nil {
		t.Fatal("expected missing people search args to fail")
	}
}

func TestGoogleResourcePathValidation(t *testing.T) {
	path, err := googleResourcePath("space-name", "spaces/AAA/messages/BBB", "spaces/")
	if err != nil {
		t.Fatal(err)
	}
	if path != "spaces/AAA/messages/BBB" {
		t.Fatalf("expected slashes to be preserved, got %q", path)
	}
	if _, err := googleResourcePath("resource-name", "people/me?personFields=names", "people/"); err == nil {
		t.Fatal("expected query characters to fail")
	}
	if _, err := googleResourcePath("resource-name", "contactGroups/all", "people/"); err == nil {
		t.Fatal("expected wrong prefix to fail")
	}
}

func validUnitConfig() Config {
	env := validUnitEnv()
	config, err := loadConfig(env)
	if err != nil {
		panic(err)
	}
	return config
}
