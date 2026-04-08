package textinput

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestReadReturnsInlineValue(t *testing.T) {
	text, err := Read(bytes.NewBuffer(nil), "body", "hello", "body-file", "")
	if err != nil {
		t.Fatal(err)
	}

	if text != "hello" {
		t.Fatalf("expected hello, got %q", text)
	}
}

func TestReadRejectsMissingSources(t *testing.T) {
	_, err := Read(bytes.NewBuffer(nil), "body", "", "body-file", "")
	if err == nil {
		t.Fatal("expected missing sources to fail")
	}

	if err.Error() != "exactly one of --body or --body-file is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReadRejectsConflictingSources(t *testing.T) {
	_, err := Read(bytes.NewBuffer(nil), "body", "hello", "body-file", "comment.txt")
	if err == nil {
		t.Fatal("expected conflicting sources to fail")
	}

	if err.Error() != "--body and --body-file are mutually exclusive" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReadReadsFile(t *testing.T) {
	tempDir := t.TempDir()
	commentFile := filepath.Join(tempDir, "comment.txt")
	if err := os.WriteFile(commentFile, []byte("comment from file"), 0o600); err != nil {
		t.Fatal(err)
	}

	text, err := Read(bytes.NewBuffer(nil), "body", "", "body-file", commentFile)
	if err != nil {
		t.Fatal(err)
	}

	if text != "comment from file" {
		t.Fatalf("expected file contents, got %q", text)
	}
}

func TestReadReadsStdin(t *testing.T) {
	text, err := Read(bytes.NewBufferString("comment from stdin"), "body", "", "body-file", "-")
	if err != nil {
		t.Fatal(err)
	}

	if text != "comment from stdin" {
		t.Fatalf("expected stdin contents, got %q", text)
	}
}
