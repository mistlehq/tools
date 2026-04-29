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

func TestFilesInfoJSON(t *testing.T) {
	env, sc := setupSlackClient(t)
	_, preflightErr := sc.GetFileInfo("F0000000000")
	if preflightErr != nil && strings.Contains(preflightErr.Error(), "missing_scope") {
		t.Skip("skipping: Slack test bot token does not have files:read")
	}

	channelID := getRequiredEnv(t, "SLACK_TEST_CHANNEL_ID")

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "info-json.txt")
	if err := os.WriteFile(filePath, []byte("slack file info json integration test"), 0o600); err != nil {
		t.Fatal(err)
	}

	uploaded, err := sc.UploadFile(SlackFilesUploadInput{
		Path:    filePath,
		Channel: channelID,
	})
	if err != nil {
		t.Fatal(err)
	}

	file := firstUploadedFile(uploaded)
	if file.ID != "" {
		t.Cleanup(func() {
			_ = sc.DeleteFile(file.ID)
		})
	}

	commandResult, err := runCommandWithInput(t, env, "", "slack", "files", "info", "--file", file.ID, "--json")
	if err != nil {
		t.Fatal(err)
	}

	var info SlackFilesInfo
	if err := json.Unmarshal(commandResult.stdout.Bytes(), &info); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}

	if !info.OK {
		t.Fatal("expected ok=true in JSON output")
	}

	if info.File.ID != file.ID {
		t.Fatalf("expected file id %q, got %q", file.ID, info.File.ID)
	}
}

func TestFilesInfoRequiresFile(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{"SLACK_BASE_URL": "http://127.0.0.1"}, "", "slack", "files", "info")
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "files info requires --file" {
		t.Fatalf("expected missing file error, got %q", err.Error())
	}
}

func TestFilesDownloadRequiresOutput(t *testing.T) {
	_, err := runCommandWithInput(t, Environment{"SLACK_BASE_URL": "http://127.0.0.1"}, "", "slack", "files", "download", "--file", "F123")
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "files download requires --output" {
		t.Fatalf("expected missing output error, got %q", err.Error())
	}
}

func TestResolveSlackFileDownloadURL(t *testing.T) {
	downloadURL := resolveSlackFileDownloadURL(SlackFile{
		URLPrivate:         "https://files.slack.com/files-pri/T123-F123/file.txt",
		URLPrivateDownload: "https://files.slack.com/files-pri/T123-F123/download/file.txt",
	})
	if downloadURL != "https://files.slack.com/files-pri/T123-F123/download/file.txt" {
		t.Fatalf("expected url_private_download, got %q", downloadURL)
	}

	privateURL := resolveSlackFileDownloadURL(SlackFile{
		URLPrivate: "https://files.slack.com/files-pri/T123-F123/file.txt",
	})
	if privateURL != "https://files.slack.com/files-pri/T123-F123/file.txt" {
		t.Fatalf("expected url_private, got %q", privateURL)
	}
}
