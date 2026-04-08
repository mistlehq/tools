package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilesUpload(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "upload.txt")
	if err := os.WriteFile(filePath, []byte("slack file upload integration test"), 0o600); err != nil {
		t.Fatal(err)
	}

	commandResult, err := runCommandWithInput(t, env, "", "slack", "files", "upload", "--path", filePath, "--channel", channelID)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(commandResult.stdout.String())
	fileID := parseLineValue(t, output, "File ID: ")
	t.Cleanup(func() {
		if fileID != "" {
			_ = sc.DeleteFile(fileID)
		}
	})

	expectedPrefixes := []string{
		"Channel: ",
		"Thread TS: ",
		"File ID: ",
		"Name: ",
		"Title: ",
	}

	lines := strings.Split(output, "\n")
	if len(lines) != len(expectedPrefixes) {
		t.Fatalf("expected %d lines, got %d: %q", len(expectedPrefixes), len(lines), output)
	}

	for index, want := range expectedPrefixes {
		if !strings.HasPrefix(lines[index], want) {
			t.Fatalf("expected line %d to start with %q, got %q", index+1, want, lines[index])
		}
	}
}

func TestFilesUploadJSON(t *testing.T) {
	env, sc := setupSlackClient(t)
	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "upload-json.txt")
	if err := os.WriteFile(filePath, []byte("slack file upload json integration test"), 0o600); err != nil {
		t.Fatal(err)
	}

	commandResult, err := runCommandWithInput(t, env, "", "slack", "files", "upload", "--path", filePath, "--channel", channelID, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var uploaded SlackFilesCompleteUploadExternal
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &uploaded); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !uploaded.OK {
		t.Fatal("expected ok=true in JSON output")
	}

	file := firstUploadedFile(uploaded)
	if file.ID != "" {
		t.Cleanup(func() {
			_ = sc.DeleteFile(file.ID)
		})
	}
}
