package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestJiraAttachmentAPIAccess(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	jc := NewJiraClient(config)
	contents := "jira attachment access probe\n"
	attachment := uploadJiraTestAttachment(t, jc, issueKey, "jira-attachment-access-probe.txt", contents)

	attachments, err := jc.ListIssueAttachments(issueKey)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, candidate := range attachments {
		if candidate.ID == attachment.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected attachment %s to be listed on issue %s: %#v", attachment.ID, issueKey, attachments)
	}

	commandResult, err := runCommandWithInput(t, env, "", "jira", "issue", "attachment", "list", issueKey)
	if err != nil {
		t.Fatal(err)
	}
	output := commandResult.stdout.String()
	if !strings.Contains(output, string(attachment.ID)) || !strings.Contains(output, attachment.Filename) {
		t.Fatalf("expected CLI attachment list to include %s/%s, got %q", attachment.ID, attachment.Filename, output)
	}

	downloadPath := t.TempDir() + "/downloaded-attachment.txt"
	downloadResult, err := runCommandWithInput(t, env, "", "jira", "issue", "attachment", "download", string(attachment.ID), "--output", downloadPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(downloadResult.stdout.String(), "Attachment ID: "+string(attachment.ID)) {
		t.Fatalf("expected CLI download output to include attachment ID, got %q", downloadResult.stdout.String())
	}
	downloadedContent, err := os.ReadFile(downloadPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(downloadedContent) != contents {
		t.Fatalf("expected CLI downloaded attachment contents %q, got %q", contents, string(downloadedContent))
	}

	session := newJiraMCPTestSession(t, jc)
	defer session.Close()

	listResult := callJiraMCPTool(t, session, "jira_issue_attachment_list", map[string]any{
		"issueKey": issueKey,
	})
	var listOutput jiraIssueAttachmentListToolOutput
	decodeMCPStructuredContent(t, listResult, &listOutput)
	found = false
	for _, candidate := range listOutput.Attachments {
		if candidate.ID == attachment.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected MCP attachment list to include %s, got %#v", attachment.ID, listOutput.Attachments)
	}

	downloadMCPResult := callJiraMCPTool(t, session, "jira_issue_attachment_download", map[string]any{
		"attachmentId": string(attachment.ID),
	})
	var downloadMCPOutput jiraIssueAttachmentDownloadToolOutput
	decodeMCPStructuredContent(t, downloadMCPResult, &downloadMCPOutput)
	if downloadMCPOutput.Attachment.ID != attachment.ID {
		t.Fatalf("expected MCP attachment download metadata for %s, got %#v", attachment.ID, downloadMCPOutput.Attachment)
	}
	decodedContent, err := base64.StdEncoding.DecodeString(downloadMCPOutput.ContentBase64)
	if err != nil {
		t.Fatal(err)
	}
	if string(decodedContent) != contents {
		t.Fatalf("expected MCP downloaded attachment contents %q, got %q", contents, string(decodedContent))
	}

	metadataBody, err := jc.get(fmt.Sprintf("/rest/api/3/attachment/%s", attachment.ID))
	if err != nil {
		t.Fatal(err)
	}

	var metadata JiraAttachment
	if err := json.Unmarshal(metadataBody, &metadata); err != nil {
		t.Fatal(err)
	}

	if metadata.Filename != attachment.Filename {
		t.Fatalf("expected metadata filename %q, got %q", attachment.Filename, metadata.Filename)
	}
	if metadata.Size != len(contents) {
		t.Fatalf("expected metadata size %d, got %d", len(contents), metadata.Size)
	}

	contentPath := fmt.Sprintf("/rest/api/3/attachment/content/%s?redirect=false", attachment.ID)
	contentBody, err := jc.get(contentPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contentBody) != contents {
		t.Fatalf("expected downloaded attachment contents %q, got %q", contents, string(contentBody))
	}
}

func uploadJiraTestAttachment(t *testing.T, jc JiraClient, issueKey string, filename string, contents string) JiraAttachment {
	t.Helper()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(part, strings.NewReader(contents)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		jc.baseURL+"/rest/api/3/issue/"+issueKey+"/attachments",
		&requestBody,
	)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("X-Atlassian-Token", "no-check")

	response, err := jc.client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		t.Fatalf("attachment upload failed with status %d: %s", response.StatusCode, string(responseBody))
	}

	var attachments []JiraAttachment
	if err := json.Unmarshal(responseBody, &attachments); err != nil {
		t.Fatal(err)
	}
	if len(attachments) != 1 {
		t.Fatalf("expected one uploaded attachment, got %d: %s", len(attachments), string(responseBody))
	}

	return attachments[0]
}
