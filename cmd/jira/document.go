package main

import (
	"fmt"
	"strings"
)

type JiraDocument struct {
	Type    string        `json:"type"`
	Version int           `json:"version"`
	Content []JiraDocNode `json:"content"`
}

type JiraDocNode struct {
	Type    string         `json:"type"`
	Text    string         `json:"text,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Marks   []JiraDocMark  `json:"marks,omitempty"`
	Content []JiraDocNode  `json:"content,omitempty"`
}

type JiraDocMark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

// NewJiraTextDocument converts plain text into the minimal Atlassian document
// structure needed for comment bodies.
func NewJiraTextDocument(text string) (JiraDocument, error) {
	if strings.TrimSpace(text) == "" {
		return JiraDocument{}, fmt.Errorf("text document must not be empty")
	}

	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	content := make([]JiraDocNode, 0, len(lines))

	for _, line := range lines {
		paragraph := JiraDocNode{
			Type: "paragraph",
		}

		if line != "" {
			paragraph.Content = []JiraDocNode{{
				Type: "text",
				Text: line,
			}}
		}

		content = append(content, paragraph)
	}

	return JiraDocument{
		Type:    "doc",
		Version: 1,
		Content: content,
	}, nil
}

func (document JiraDocument) PlainText() string {
	var parts []string
	for _, node := range document.Content {
		collectJiraDocText(node, &parts)
	}

	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func collectJiraDocText(node JiraDocNode, parts *[]string) {
	if node.Text != "" {
		*parts = append(*parts, node.Text)
	}

	if href, ok := node.Attrs["href"].(string); ok && href != "" {
		*parts = append(*parts, href)
	}
	for _, mark := range node.Marks {
		if href, ok := mark.Attrs["href"].(string); ok && href != "" {
			*parts = append(*parts, href)
		}
	}

	var childParts []string
	for _, child := range node.Content {
		collectJiraDocText(child, &childParts)
	}

	if len(childParts) == 0 {
		return
	}

	switch node.Type {
	case "paragraph", "heading", "blockquote", "listItem":
		*parts = append(*parts, strings.Join(childParts, " "))
	default:
		*parts = append(*parts, childParts...)
	}
}
