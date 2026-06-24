package main

import "testing"

func TestRawRequestValidation(t *testing.T) {
	_, err := parseRawRequestArgs([]string{"--api", "drive", "--method", "TRACE", "--path", "/files"})
	if err != nil {
		t.Fatal(err)
	}
	client := NewGWSClient(Config{DriveBaseURL: "http://127.0.0.1"})
	if _, err := client.Request(GWSRequest{API: "drive", Method: "TRACE", Path: "/files"}); err == nil {
		t.Fatal("expected unsupported method to fail")
	}
	if _, err := client.Request(GWSRequest{API: "drive", Method: "GET"}); err == nil {
		t.Fatal("expected missing path to fail")
	}
	if _, err := client.Request(GWSRequest{API: "drive", Method: "GET", Path: "files"}); err == nil {
		t.Fatal("expected relative path to fail")
	}
	if _, err := client.Request(GWSRequest{API: "calendar", Method: "GET", Path: "/events"}); err == nil {
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
}
