package main

import "testing"

func TestParseGraphRequestArgs(t *testing.T) {
	request, err := parseGraphRequestArgs([]string{
		"--method", "POST",
		"--path", "/act_123/campaigns",
		"--params", `{"fields":"id,name"}`,
		"--body", `{"name":"Example"}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if request.Method != "POST" || request.Path != "/act_123/campaigns" {
		t.Fatalf("unexpected request: %#v", request)
	}
	if request.Params["fields"] != "id,name" {
		t.Fatalf("unexpected params: %#v", request.Params)
	}
	if request.Body["name"] != "Example" {
		t.Fatalf("unexpected body: %#v", request.Body)
	}
}

func TestParseEdgeArgsRequiresAccountWhenRequested(t *testing.T) {
	_, err := parseEdgeArgs([]string{}, "campaigns search", true)
	if err == nil {
		t.Fatal("expected missing account id error")
	}
}

func TestParseInsightsArgs(t *testing.T) {
	input, err := parseInsightsArgs([]string{
		"--id", "act_123",
		"--fields", "impressions,spend",
		"--level", "campaign",
		"--time-range", `{"since":"2026-06-01","until":"2026-06-23"}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if input.ID != "act_123" || input.Level != "campaign" {
		t.Fatalf("unexpected input: %#v", input)
	}
}

func TestRequestRequiresAbsolutePath(t *testing.T) {
	client := NewMetaAdsClient(Config{GraphBaseURL: "https://graph.facebook.com/v25.0"})
	_, err := client.Request(MetaAdsRequest{Method: "GET", Path: "me"})
	if err == nil {
		t.Fatal("expected path validation error")
	}
}
