package main

import "testing"

func TestNewJiraTextDocumentRejectsEmptyText(t *testing.T) {
	_, err := NewJiraTextDocument("   ")
	if err == nil {
		t.Fatal("expected empty text to fail")
	}

	if err.Error() != "text document must not be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewJiraTextDocumentCreatesSingleParagraph(t *testing.T) {
	document, err := NewJiraTextDocument("hello")
	if err != nil {
		t.Fatal(err)
	}

	if document.Type != "doc" || document.Version != 1 {
		t.Fatalf("unexpected document header: %#v", document)
	}

	if len(document.Content) != 1 {
		t.Fatalf("expected 1 paragraph, got %d", len(document.Content))
	}

	paragraph := document.Content[0]
	if paragraph.Type != "paragraph" {
		t.Fatalf("expected paragraph node, got %#v", paragraph)
	}

	if len(paragraph.Content) != 1 || paragraph.Content[0].Text != "hello" {
		t.Fatalf("unexpected paragraph content: %#v", paragraph.Content)
	}
}

func TestNewJiraTextDocumentPreservesBlankLines(t *testing.T) {
	document, err := NewJiraTextDocument("first\n\nthird")
	if err != nil {
		t.Fatal(err)
	}

	if len(document.Content) != 3 {
		t.Fatalf("expected 3 paragraphs, got %d", len(document.Content))
	}

	if len(document.Content[1].Content) != 0 {
		t.Fatalf("expected blank line to become empty paragraph, got %#v", document.Content[1])
	}
}

func TestNewJiraTextDocumentNormalizesWindowsNewlines(t *testing.T) {
	document, err := NewJiraTextDocument("first\r\nsecond")
	if err != nil {
		t.Fatal(err)
	}

	if len(document.Content) != 2 {
		t.Fatalf("expected 2 paragraphs, got %d", len(document.Content))
	}

	if document.Content[0].Content[0].Text != "first" || document.Content[1].Content[0].Text != "second" {
		t.Fatalf("unexpected normalized content: %#v", document.Content)
	}
}
