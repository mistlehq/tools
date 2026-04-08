package argparse

import "testing"

func TestParseAcceptsPositionalsAndFlags(t *testing.T) {
	parsed, err := Parse([]string{
		"KAN-1",
		"--body", "hello",
		"extra",
	}, map[string]Spec{
		"body": {TakesValue: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(parsed.Positionals) != 2 {
		t.Fatalf("expected 2 positionals, got %d", len(parsed.Positionals))
	}

	if parsed.Positionals[0] != "KAN-1" || parsed.Positionals[1] != "extra" {
		t.Fatalf("unexpected positionals: %#v", parsed.Positionals)
	}

	if parsed.First("body") != "hello" {
		t.Fatalf("expected body flag to equal hello, got %q", parsed.First("body"))
	}
}

func TestParseAcceptsEqualsSyntax(t *testing.T) {
	parsed, err := Parse([]string{
		"--body=hello",
	}, map[string]Spec{
		"body": {TakesValue: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	if parsed.First("body") != "hello" {
		t.Fatalf("expected body flag to equal hello, got %q", parsed.First("body"))
	}
}

func TestParseAcceptsBooleanFlags(t *testing.T) {
	parsed, err := Parse([]string{
		"--verbose",
	}, map[string]Spec{
		"verbose": {},
	})
	if err != nil {
		t.Fatal(err)
	}

	if parsed.First("verbose") != "true" {
		t.Fatalf("expected boolean flag to be recorded as true, got %q", parsed.First("verbose"))
	}
}

func TestParseRejectsUnsupportedFlag(t *testing.T) {
	_, err := Parse([]string{"--nope"}, map[string]Spec{})
	if err == nil {
		t.Fatal("expected unsupported flag to fail")
	}

	if err.Error() != "unsupported flag: --nope" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsMissingValue(t *testing.T) {
	_, err := Parse([]string{"--body"}, map[string]Spec{
		"body": {TakesValue: true},
	})
	if err == nil {
		t.Fatal("expected missing flag value to fail")
	}

	if err.Error() != "flag --body requires a value" {
		t.Fatalf("unexpected error: %v", err)
	}
}
