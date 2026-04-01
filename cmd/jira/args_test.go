package main

import "testing"

func TestParseArgsAcceptsPositionalsAndFlags(t *testing.T) {
	parsed, err := parseArgs([]string{
		"KAN-1",
		"--body", "hello",
		"extra",
	}, map[string]argSpec{
		"body": {takesValue: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(parsed.positionals) != 2 {
		t.Fatalf("expected 2 positionals, got %d", len(parsed.positionals))
	}

	if parsed.positionals[0] != "KAN-1" || parsed.positionals[1] != "extra" {
		t.Fatalf("unexpected positionals: %#v", parsed.positionals)
	}

	if parsed.first("body") != "hello" {
		t.Fatalf("expected body flag to equal hello, got %q", parsed.first("body"))
	}
}

func TestParseArgsAcceptsEqualsSyntax(t *testing.T) {
	parsed, err := parseArgs([]string{
		"--body=hello",
	}, map[string]argSpec{
		"body": {takesValue: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	if parsed.first("body") != "hello" {
		t.Fatalf("expected body flag to equal hello, got %q", parsed.first("body"))
	}
}

func TestParseArgsAcceptsBooleanFlags(t *testing.T) {
	parsed, err := parseArgs([]string{
		"--verbose",
	}, map[string]argSpec{
		"verbose": {},
	})
	if err != nil {
		t.Fatal(err)
	}

	if parsed.first("verbose") != "true" {
		t.Fatalf("expected boolean flag to be recorded as true, got %q", parsed.first("verbose"))
	}
}

func TestParseArgsRejectsUnsupportedFlag(t *testing.T) {
	_, err := parseArgs([]string{"--nope"}, map[string]argSpec{})
	if err == nil {
		t.Fatal("expected unsupported flag to fail")
	}

	if err.Error() != "unsupported flag: --nope" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseArgsRejectsMissingValue(t *testing.T) {
	_, err := parseArgs([]string{"--body"}, map[string]argSpec{
		"body": {takesValue: true},
	})
	if err == nil {
		t.Fatal("expected missing flag value to fail")
	}

	if err.Error() != "flag --body requires a value" {
		t.Fatalf("unexpected error: %v", err)
	}
}
