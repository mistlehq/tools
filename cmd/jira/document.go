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
	Type    string        `json:"type"`
	Text    string        `json:"text,omitempty"`
	Content []JiraDocNode `json:"content,omitempty"`
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
