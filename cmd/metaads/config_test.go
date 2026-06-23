package main

import "testing"

func TestLoadConfig(t *testing.T) {
	config, err := loadConfig(Environment{"METAADS_GRAPH_BASE_URL": "https://graph.facebook.com/v25.0"})
	if err != nil {
		t.Fatal(err)
	}
	if config.GraphBaseURL != "https://graph.facebook.com/v25.0" {
		t.Fatalf("unexpected graph base url: %q", config.GraphBaseURL)
	}
}

func TestLoadConfigRequiresBaseURL(t *testing.T) {
	_, err := loadConfig(Environment{})
	if err == nil {
		t.Fatal("expected missing base URL error")
	}
}

func TestLoadConfigRejectsTrailingSlash(t *testing.T) {
	_, err := loadConfig(Environment{"METAADS_GRAPH_BASE_URL": "https://graph.facebook.com/v25.0/"})
	if err == nil {
		t.Fatal("expected trailing slash error")
	}
}
