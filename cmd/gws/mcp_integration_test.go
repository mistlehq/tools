package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHelp(t *testing.T) {
	result, err := runCommandWithInput(t, Environment{}, "", "gws", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}
	if !containsAll(result.stdout.String(), []string{"gws mcp", "gws mcp serve", "Streamable HTTP"}) {
		t.Fatalf("unexpected mcp help: %s", result.stdout.String())
	}
}

func TestMCPServeHelp(t *testing.T) {
	result, err := runCommandWithInput(t, Environment{}, "", "gws", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}
	if !containsAll(result.stdout.String(), []string{"--tools <list>", "gws_drive_files_list", "gws_sheets_values_update", "gws_docs_document_batch_update", "gws_slides_presentation_batch_update", "gws_gmail_messages_list", "gws_calendar_freebusy_query", "gws_chat_spaces_list", "gws_people_search_contacts"}) {
		t.Fatalf("unexpected mcp serve help: %s", result.stdout.String())
	}
}

func TestMCPServerListsGWSToolsWithAnnotations(t *testing.T) {
	session := newLocalGWSMCPTestSession(t)
	defer session.Close()

	toolsResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	toolsByName := make(map[string]*mcp.Tool)
	for _, tool := range toolsResult.Tools {
		toolsByName[tool.Name] = tool
	}
	expected := map[string]commandDoc{
		"gws_request":                          gwsRequestDoc,
		"gws_auth_test":                        gwsAuthTestDoc,
		"gws_drive_files_list":                 gwsDriveFilesListDoc,
		"gws_drive_file_get":                   gwsDriveFileGetDoc,
		"gws_drive_file_create":                gwsDriveFileCreateDoc,
		"gws_drive_file_copy":                  gwsDriveFileCopyDoc,
		"gws_drive_file_update":                gwsDriveFileUpdateDoc,
		"gws_drive_file_delete":                gwsDriveFileDeleteDoc,
		"gws_drive_permissions_list":           gwsDrivePermissionsListDoc,
		"gws_drive_permission_create":          gwsDrivePermissionCreateDoc,
		"gws_drive_permission_delete":          gwsDrivePermissionDeleteDoc,
		"gws_sheets_spreadsheet_get":           gwsSheetsSpreadsheetGetDoc,
		"gws_sheets_spreadsheet_create":        gwsSheetsSpreadsheetCreateDoc,
		"gws_sheets_spreadsheet_batch_update":  gwsSheetsSpreadsheetBatchUpdateDoc,
		"gws_sheets_values_get":                gwsSheetsValuesGetDoc,
		"gws_sheets_values_update":             gwsSheetsValuesUpdateDoc,
		"gws_sheets_values_batch_update":       gwsSheetsValuesBatchUpdateDoc,
		"gws_docs_document_get":                gwsDocsDocumentGetDoc,
		"gws_docs_document_batch_update":       gwsDocsDocumentBatchUpdateDoc,
		"gws_slides_presentation_get":          gwsSlidesPresentationGetDoc,
		"gws_slides_presentation_create":       gwsSlidesPresentationCreateDoc,
		"gws_slides_presentation_batch_update": gwsSlidesPresentationBatchUpdateDoc,
		"gws_gmail_messages_list":              gwsGmailMessagesListDoc,
		"gws_gmail_message_get":                gwsGmailMessageGetDoc,
		"gws_gmail_message_send":               gwsGmailMessageSendDoc,
		"gws_gmail_drafts_list":                gwsGmailDraftsListDoc,
		"gws_gmail_draft_get":                  gwsGmailDraftGetDoc,
		"gws_gmail_draft_create":               gwsGmailDraftCreateDoc,
		"gws_gmail_draft_send":                 gwsGmailDraftSendDoc,
		"gws_gmail_draft_delete":               gwsGmailDraftDeleteDoc,
		"gws_calendar_calendar_list_list":      gwsCalendarListListDoc,
		"gws_calendar_calendar_list_get":       gwsCalendarListGetDoc,
		"gws_calendar_events_list":             gwsCalendarEventsListDoc,
		"gws_calendar_event_get":               gwsCalendarEventGetDoc,
		"gws_calendar_event_insert":            gwsCalendarEventInsertDoc,
		"gws_calendar_event_patch":             gwsCalendarEventPatchDoc,
		"gws_calendar_event_delete":            gwsCalendarEventDeleteDoc,
		"gws_calendar_freebusy_query":          gwsCalendarFreeBusyQueryDoc,
		"gws_chat_spaces_list":                 gwsChatSpacesListDoc,
		"gws_chat_space_get":                   gwsChatSpaceGetDoc,
		"gws_chat_messages_list":               gwsChatMessagesListDoc,
		"gws_chat_message_get":                 gwsChatMessageGetDoc,
		"gws_chat_message_create":              gwsChatMessageCreateDoc,
		"gws_chat_members_list":                gwsChatMembersListDoc,
		"gws_people_person_get":                gwsPeoplePersonGetDoc,
		"gws_people_connections_list":          gwsPeopleConnectionsListDoc,
		"gws_people_search_contacts":           gwsPeopleSearchContactsDoc,
		"gws_people_search_directory":          gwsPeopleSearchDirectoryDoc,
	}
	for name, doc := range expected {
		tool, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected MCP tool %q to be listed", name)
		}
		if tool.Description != doc.Description {
			t.Fatalf("expected MCP tool %q description %q, got %q", name, doc.Description, tool.Description)
		}
		if tool.Annotations == nil {
			t.Fatalf("expected MCP tool %q annotations", name)
		}
		if name == "gws_request" {
			continue
		}
		if isReadOnlyGWSTool(name) && !tool.Annotations.ReadOnlyHint {
			t.Fatalf("expected MCP tool %q to be read-only", name)
		}
		if isDestructiveGWSTool(name) && (tool.Annotations.DestructiveHint == nil || !*tool.Annotations.DestructiveHint) {
			t.Fatalf("expected MCP tool %q to be destructive", name)
		}
	}
}

func TestMCPServerFiltersToolsByGroupAndName(t *testing.T) {
	session := newGWSMCPTestSession(t, NewGWSClient(Config{
		DriveBaseURL:    "http://127.0.0.1",
		SheetsBaseURL:   "http://127.0.0.1",
		DocsBaseURL:     "http://127.0.0.1",
		SlidesBaseURL:   "http://127.0.0.1",
		GmailBaseURL:    "http://127.0.0.1",
		CalendarBaseURL: "http://127.0.0.1",
		ChatBaseURL:     "http://127.0.0.1",
		PeopleBaseURL:   "http://127.0.0.1",
	}), map[string]bool{"drive": true, "gws_docs_document_get": true, "gmail": true})
	defer session.Close()

	toolsResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	toolsByName := make(map[string]bool)
	for _, tool := range toolsResult.Tools {
		toolsByName[tool.Name] = true
	}
	if !toolsByName["gws_drive_files_list"] || !toolsByName["gws_drive_file_delete"] || !toolsByName["gws_docs_document_get"] {
		t.Fatalf("expected selected drive group and docs tool, got %#v", toolsByName)
	}
	if !toolsByName["gws_gmail_messages_list"] || !toolsByName["gws_gmail_draft_delete"] {
		t.Fatalf("expected selected gmail group, got %#v", toolsByName)
	}
	if toolsByName["gws_sheets_values_get"] || toolsByName["gws_slides_presentation_get"] || toolsByName["gws_docs_document_batch_update"] || toolsByName["gws_calendar_events_list"] {
		t.Fatalf("expected unselected tools to be absent, got %#v", toolsByName)
	}
}

func TestMCPGWSTools(t *testing.T) {
	_, gc := setupGWSClient(t)
	resources := testGWSWorkspaceResources(t)
	session := newGWSMCPTestSession(t, gc, nil)
	defer session.Close()

	authResult := callGWSMCPTool(t, session, "gws_auth_test", map[string]any{})
	var auth map[string]any
	decodeMCPStructuredContent(t, authResult, &auth)
	if auth["user"] == nil {
		t.Fatalf("expected auth response to include user, got %#v", auth)
	}

	fileResult := callGWSMCPTool(t, session, "gws_drive_file_get", map[string]any{"fileId": resources.DocumentID, "fields": "id,name,mimeType"})
	var file GWSDriveFile
	decodeMCPStructuredContent(t, fileResult, &file)
	if file.ID != resources.DocumentID {
		t.Fatalf("expected document file id %s, got %#v", resources.DocumentID, file)
	}

	valuesUpdateResult := callGWSMCPTool(t, session, "gws_sheets_values_update", map[string]any{
		"spreadsheetId":    resources.SpreadsheetID,
		"range":            "Sheet1!A1",
		"valueInputOption": "RAW",
		"request": map[string]any{
			"range":  "Sheet1!A1",
			"values": [][]string{{"mcp"}},
		},
	})
	var valuesUpdate map[string]any
	decodeMCPStructuredContent(t, valuesUpdateResult, &valuesUpdate)
	if valuesUpdate["updatedCells"] == nil {
		t.Fatalf("expected updatedCells, got %#v", valuesUpdate)
	}

	insertedText := fmt.Sprintf("MCP update %d\n", nowNano())
	docResult := callGWSMCPTool(t, session, "gws_docs_document_batch_update", map[string]any{
		"documentId": resources.DocumentID,
		"request": map[string]any{
			"requests": []map[string]any{{"insertText": map[string]any{"location": map[string]any{"index": 1}, "text": insertedText}}},
		},
	})
	var docUpdate map[string]any
	decodeMCPStructuredContent(t, docResult, &docUpdate)
	if docUpdate["replies"] == nil {
		t.Fatalf("expected docs replies, got %#v", docUpdate)
	}
	t.Cleanup(func() {
		_, err := gc.BatchUpdateDocument(cliContext(), resources.DocumentID, map[string]any{
			"requests": []map[string]any{{
				"deleteContentRange": map[string]any{
					"range": map[string]any{"startIndex": 1, "endIndex": 1 + len(insertedText)},
				},
			}},
		})
		if err != nil {
			t.Logf("failed to clean MCP document text: %v", err)
		}
	})

	slideObjectID := fmt.Sprintf("mcp_slide_%d", nowNano())
	slidesResult := callGWSMCPTool(t, session, "gws_slides_presentation_batch_update", map[string]any{
		"presentationId": resources.PresentationID,
		"request": map[string]any{
			"requests": []map[string]any{{"createSlide": map[string]any{"objectId": slideObjectID, "slideLayoutReference": map[string]any{"predefinedLayout": "BLANK"}}}},
		},
	})
	var slidesUpdate map[string]any
	decodeMCPStructuredContent(t, slidesResult, &slidesUpdate)
	if slidesUpdate["replies"] == nil {
		t.Fatalf("expected slides replies, got %#v", slidesUpdate)
	}
	t.Cleanup(func() {
		_, err := gc.BatchUpdatePresentation(cliContext(), resources.PresentationID, map[string]any{
			"requests": []map[string]any{{"deleteObject": map[string]any{"objectId": slideObjectID}}},
		})
		if err != nil {
			t.Logf("failed to clean MCP slide %s: %v", slideObjectID, err)
		}
	})

	rawResult := callGWSMCPTool(t, session, "gws_request", map[string]any{"api": "drive", "method": "GET", "path": "/about", "params": map[string]any{"fields": "user"}})
	var raw map[string]any
	decodeMCPStructuredContent(t, rawResult, &raw)
	if raw["user"] == nil {
		t.Fatalf("expected raw response to include user, got %#v", raw)
	}
}

func TestMCPGWSCreateDeleteTools(t *testing.T) {
	requireDriveCreationCapable(t)
	_, gc := setupGWSClient(t)
	resources := testGWSWorkspaceResources(t)
	session := newGWSMCPTestSession(t, gc, nil)
	defer session.Close()

	createResult := callGWSMCPTool(t, session, "gws_drive_file_create", map[string]any{
		"fields": "id,name,mimeType,parents",
		"request": map[string]any{
			"name":     fmt.Sprintf("Mistle GWS MCP Created %d", nowNano()),
			"mimeType": "application/vnd.google-apps.document",
			"parents":  []any{resources.DriveFolderID},
		},
	})
	var created GWSDriveFile
	decodeMCPStructuredContent(t, createResult, &created)
	if created.ID == "" {
		t.Fatalf("expected created file id, got %#v", created)
	}
	t.Cleanup(func() {
		if created.ID == "" {
			return
		}
		if err := gc.DeleteDriveFile(cliContext(), created.ID); err != nil {
			t.Logf("failed to clean MCP created file %s: %v", created.ID, err)
		}
	})

	deleteResult := callGWSMCPTool(t, session, "gws_drive_file_delete", map[string]any{"fileId": created.ID})
	var deleted map[string]any
	decodeMCPStructuredContent(t, deleteResult, &deleted)
	if deleted["deleted"] != true {
		t.Fatalf("expected deleted response, got %#v", deleted)
	}
	created.ID = ""
}

func TestMCPGWSToolValidation(t *testing.T) {
	session := newLocalGWSMCPTestSession(t)
	defer session.Close()
	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{name: "drive get missing file id", tool: "gws_drive_file_get", arguments: map[string]any{}},
		{name: "sheets values update missing request", tool: "gws_sheets_values_update", arguments: map[string]any{"spreadsheetId": "sheet", "range": "A1", "valueInputOption": "RAW"}},
		{name: "docs batch update missing document id", tool: "gws_docs_document_batch_update", arguments: map[string]any{"request": map[string]any{}}},
		{name: "slides batch update missing presentation id", tool: "gws_slides_presentation_batch_update", arguments: map[string]any{"request": map[string]any{}}},
		{name: "gmail message get missing user id", tool: "gws_gmail_message_get", arguments: map[string]any{"messageId": "msg"}},
		{name: "calendar event get missing event id", tool: "gws_calendar_event_get", arguments: map[string]any{"calendarId": "primary"}},
		{name: "chat messages list missing space name", tool: "gws_chat_messages_list", arguments: map[string]any{}},
		{name: "people person get missing resource name", tool: "gws_people_person_get", arguments: map[string]any{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: tc.tool, Arguments: tc.arguments})
			if err != nil {
				t.Fatal(err)
			}
			if !result.IsError {
				t.Fatal("expected tool validation to return a tool error")
			}
		})
	}
}

func isReadOnlyGWSTool(name string) bool {
	switch name {
	case "gws_auth_test", "gws_drive_files_list", "gws_drive_file_get", "gws_drive_permissions_list", "gws_sheets_spreadsheet_get", "gws_sheets_values_get", "gws_docs_document_get", "gws_slides_presentation_get", "gws_gmail_messages_list", "gws_gmail_message_get", "gws_gmail_drafts_list", "gws_gmail_draft_get", "gws_calendar_calendar_list_list", "gws_calendar_calendar_list_get", "gws_calendar_events_list", "gws_calendar_event_get", "gws_calendar_freebusy_query", "gws_chat_spaces_list", "gws_chat_space_get", "gws_chat_messages_list", "gws_chat_message_get", "gws_chat_members_list", "gws_people_person_get", "gws_people_connections_list", "gws_people_search_contacts", "gws_people_search_directory":
		return true
	default:
		return false
	}
}

func isDestructiveGWSTool(name string) bool {
	switch name {
	case "gws_drive_file_delete", "gws_drive_permission_delete", "gws_gmail_draft_delete", "gws_calendar_event_delete":
		return true
	default:
		return false
	}
}
