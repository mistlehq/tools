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

func TestJiraDocumentPlainText(t *testing.T) {
	document, err := NewJiraTextDocument("first\nsecond")
	if err != nil {
		t.Fatal(err)
	}

	if got := document.PlainText(); got != "first\nsecond" {
		t.Fatalf("unexpected plain text: %q", got)
	}
}

func TestExtractJiraCommentAttachmentRefsFromPlainTextAttachmentURL(t *testing.T) {
	document, err := NewJiraTextDocument("see https://mistle.atlassian.net/rest/api/3/attachment/content/10001")
	if err != nil {
		t.Fatal(err)
	}

	refs := extractJiraCommentAttachmentRefs(document)
	if len(refs) != 1 {
		t.Fatalf("expected 1 attachment reference, got %#v", refs)
	}

	if refs[0].Type != "link" || refs[0].URL == "" {
		t.Fatalf("unexpected attachment reference: %#v", refs[0])
	}
}

func TestExtractJiraCommentAttachmentRefsFromLinkMark(t *testing.T) {
	attachmentURL := "https://mistle.atlassian.net/rest/api/3/attachment/content/10001"
	document := JiraDocument{
		Type:    "doc",
		Version: 1,
		Content: []JiraDocNode{{
			Type: "paragraph",
			Content: []JiraDocNode{{
				Type: "text",
				Text: "linked attachment",
				Marks: []JiraDocMark{{
					Type: "link",
					Attrs: map[string]any{
						"href": attachmentURL,
					},
				}},
			}},
		}},
	}

	refs := extractJiraCommentAttachmentRefs(document)
	if len(refs) != 1 {
		t.Fatalf("expected 1 attachment reference, got %#v", refs)
	}

	if refs[0].Type != "link" || refs[0].URL != attachmentURL {
		t.Fatalf("unexpected attachment reference: %#v", refs[0])
	}
}
