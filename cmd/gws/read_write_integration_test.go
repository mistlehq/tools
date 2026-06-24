package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestAuthAndRawRequestCommands(t *testing.T) {
	env := setupCommandEnvironment(t)

	authResult, err := runCommandWithInput(t, env, "", "gws", "auth", "test")
	if err != nil {
		t.Fatal(err)
	}
	var auth map[string]any
	decodeCommandJSON(t, authResult, &auth)
	if auth["user"] == nil {
		t.Fatalf("expected auth response to include user, got %#v", auth)
	}

	requestResult, err := runCommandWithInput(t, env, "", "gws", "request", "--api", "drive", "--method", "GET", "--path", "/about?fields=user")
	if err != nil {
		t.Fatal(err)
	}
	var request map[string]any
	decodeCommandJSON(t, requestResult, &request)
	if request["user"] == nil {
		t.Fatalf("expected raw request response to include user, got %#v", request)
	}
}

func TestDriveFileCommands(t *testing.T) {
	env, gc := setupGWSClient(t)
	resources := testGWSWorkspaceResources(t)

	listResult, err := runCommandWithInput(t, env, "", "gws", "drive", "files", "list", "--query", fmt.Sprintf("'%s' in parents and trashed = false", resources.DriveFolderID), "--page-size", "10", "--fields", "files(id,name,mimeType,parents)")
	if err != nil {
		t.Fatal(err)
	}
	var list GWSDriveFilesList
	decodeCommandJSON(t, listResult, &list)
	if !driveFilesContain(list.Files, resources.DriveFileID) {
		t.Fatalf("expected file list to include file %s, got %#v", resources.DriveFileID, list.Files)
	}

	getResult, err := runCommandWithInput(t, env, "", "gws", "drive", "files", "get", "--file-id", resources.DriveFileID, "--fields", "id,name,mimeType,parents")
	if err != nil {
		t.Fatal(err)
	}
	var file GWSDriveFile
	decodeCommandJSON(t, getResult, &file)
	if file.ID != resources.DriveFileID {
		t.Fatalf("expected file id %s, got %#v", resources.DriveFileID, file)
	}

	updatePath := writeTempJSONRequest(t, fmt.Sprintf(`{"appProperties":{"mistleGwsIntegrationTest":"%d"}}`, nowNano()))
	updateResult, err := runCommandWithInput(t, env, "", "gws", "drive", "files", "update", "--file-id", resources.DriveFileID, "--request-file", updatePath, "--fields", "id,name")
	if err != nil {
		t.Fatal(err)
	}
	var updated GWSDriveFile
	decodeCommandJSON(t, updateResult, &updated)
	if updated.ID != resources.DriveFileID {
		t.Fatalf("expected updated file %s, got %#v", resources.DriveFileID, updated)
	}

	if _, err := gc.GetDriveFile(cliContext(), resources.DriveFileID, map[string]any{"fields": "id"}); err != nil {
		t.Fatalf("expected updated file to remain readable: %v", err)
	}
}

func TestDriveFileCreateCopyDeleteCommands(t *testing.T) {
	requireDriveCreationCapable(t)
	env, gc := setupGWSClient(t)
	resources := testGWSWorkspaceResources(t)

	createPath := writeTempJSONRequest(t, fmt.Sprintf(`{"name":"Mistle GWS CLI Created %d","mimeType":"application/vnd.google-apps.document","parents":[%q]}`, nowNano(), resources.DriveFolderID))
	createResult, err := runCommandWithInput(t, env, "", "gws", "drive", "files", "create", "--request-file", createPath, "--fields", "id,name,mimeType,parents")
	if err != nil {
		t.Fatal(err)
	}
	var created GWSDriveFile
	decodeCommandJSON(t, createResult, &created)
	if created.ID == "" {
		t.Fatalf("expected created file id, got %#v", created)
	}
	defer func() {
		if created.ID == "" {
			return
		}
		if err := gc.DeleteDriveFile(cliContext(), created.ID); err != nil {
			t.Logf("failed to clean created file %s: %v", created.ID, err)
		}
	}()

	copyPath := writeTempJSONRequest(t, fmt.Sprintf(`{"name":"Mistle GWS CLI Copy %d","parents":[%q]}`, nowNano(), resources.DriveFolderID))
	copyResult, err := runCommandWithInput(t, env, "", "gws", "drive", "files", "copy", "--file-id", created.ID, "--request-file", copyPath, "--fields", "id,name,mimeType,parents")
	if err != nil {
		t.Fatal(err)
	}
	var copied GWSDriveFile
	decodeCommandJSON(t, copyResult, &copied)
	if copied.ID == "" || copied.ID == created.ID {
		t.Fatalf("expected copied file id distinct from source, got %#v", copied)
	}
	defer func() {
		if copied.ID == "" {
			return
		}
		if err := gc.DeleteDriveFile(cliContext(), copied.ID); err != nil {
			t.Logf("failed to clean copied file %s: %v", copied.ID, err)
		}
	}()

	updatePath := writeTempJSONRequest(t, fmt.Sprintf(`{"name":"Mistle GWS CLI Updated %d"}`, nowNano()))
	updateResult, err := runCommandWithInput(t, env, "", "gws", "drive", "files", "update", "--file-id", copied.ID, "--request-file", updatePath, "--fields", "id,name")
	if err != nil {
		t.Fatal(err)
	}
	var updated GWSDriveFile
	decodeCommandJSON(t, updateResult, &updated)
	if updated.ID != copied.ID || !strings.Contains(updated.Name, "Updated") {
		t.Fatalf("expected updated copied file, got %#v", updated)
	}

	deleteResult, err := runCommandWithInput(t, env, "", "gws", "drive", "files", "delete", "--file-id", copied.ID)
	if err != nil {
		t.Fatal(err)
	}
	copied.ID = ""
	var deleted map[string]any
	decodeCommandJSON(t, deleteResult, &deleted)
	if deleted["deleted"] != true {
		t.Fatalf("expected deleted response, got %#v", deleted)
	}
}

func TestDrivePermissionCommands(t *testing.T) {
	env, _ := setupGWSClient(t)
	resources := testGWSWorkspaceResources(t)

	listResult, err := runCommandWithInput(t, env, "", "gws", "drive", "permissions", "list", "--file-id", resources.DriveFileID)
	if err != nil {
		t.Fatal(err)
	}
	var permissions GWSDrivePermissionsList
	decodeCommandJSON(t, listResult, &permissions)
	if len(permissions.Permissions) == 0 {
		t.Fatalf("expected at least one permission, got %#v", permissions)
	}

	createPath := writeTempJSONRequest(t, `{"type":"anyone","role":"reader","allowFileDiscovery":false}`)
	createResult, err := runCommandWithInput(t, env, "", "gws", "drive", "permissions", "create", "--file-id", resources.DriveFileID, "--request-file", createPath)
	if err != nil {
		t.Fatal(err)
	}
	var permission GWSDrivePermission
	decodeCommandJSON(t, createResult, &permission)
	if permission.ID == "" {
		t.Fatalf("expected created permission id, got %#v", permission)
	}

	deleteResult, err := runCommandWithInput(t, env, "", "gws", "drive", "permissions", "delete", "--file-id", resources.DriveFileID, "--permission-id", permission.ID)
	if err != nil {
		t.Fatal(err)
	}
	var deleted map[string]any
	decodeCommandJSON(t, deleteResult, &deleted)
	if deleted["deleted"] != true {
		t.Fatalf("expected permission deleted response, got %#v", deleted)
	}
}

func TestSheetsCommands(t *testing.T) {
	env, gc := setupGWSClient(t)
	resources := testGWSWorkspaceResources(t)

	getResult, err := runCommandWithInput(t, env, "", "gws", "sheets", "spreadsheets", "get", "--spreadsheet-id", resources.SpreadsheetID)
	if err != nil {
		t.Fatal(err)
	}
	var spreadsheet map[string]any
	decodeCommandJSON(t, getResult, &spreadsheet)
	if spreadsheet["spreadsheetId"] != resources.SpreadsheetID {
		t.Fatalf("expected spreadsheet id %s, got %#v", resources.SpreadsheetID, spreadsheet)
	}

	sheetTitle := fmt.Sprintf("Extra_%d", nowNano())
	batchPath := writeTempJSONRequest(t, fmt.Sprintf(`{"requests":[{"addSheet":{"properties":{"title":%q}}}]}`, sheetTitle))
	batchResult, err := runCommandWithInput(t, env, "", "gws", "sheets", "spreadsheets", "batch-update", "--spreadsheet-id", resources.SpreadsheetID, "--request-file", batchPath)
	if err != nil {
		t.Fatal(err)
	}
	var batch map[string]any
	decodeCommandJSON(t, batchResult, &batch)
	if batch["replies"] == nil {
		t.Fatalf("expected batch update replies, got %#v", batch)
	}
	cleanupAddedSheet(t, gc, resources.SpreadsheetID, batch)

	updatePath := writeTempJSONRequest(t, `{"range":"Sheet1!A1:B2","majorDimension":"ROWS","values":[["alpha","beta"],["gamma","delta"]]}`)
	updateResult, err := runCommandWithInput(t, env, "", "gws", "sheets", "values", "update", "--spreadsheet-id", resources.SpreadsheetID, "--range", "Sheet1!A1:B2", "--value-input-option", "RAW", "--request-file", updatePath)
	if err != nil {
		t.Fatal(err)
	}
	var update map[string]any
	decodeCommandJSON(t, updateResult, &update)
	if update["updatedCells"] == nil {
		t.Fatalf("expected updatedCells, got %#v", update)
	}

	getValuesResult, err := runCommandWithInput(t, env, "", "gws", "sheets", "values", "get", "--spreadsheet-id", resources.SpreadsheetID, "--range", "Sheet1!A1:B2")
	if err != nil {
		t.Fatal(err)
	}
	var values map[string]any
	decodeCommandJSON(t, getValuesResult, &values)
	if values["values"] == nil {
		t.Fatalf("expected values, got %#v", values)
	}

	batchValuesPath := writeTempJSONRequest(t, `{"valueInputOption":"RAW","data":[{"range":"Sheet1!C1","values":[["batch"]]}]}`)
	batchValuesResult, err := runCommandWithInput(t, env, "", "gws", "sheets", "values", "batch-update", "--spreadsheet-id", resources.SpreadsheetID, "--request-file", batchValuesPath)
	if err != nil {
		t.Fatal(err)
	}
	var batchValues map[string]any
	decodeCommandJSON(t, batchValuesResult, &batchValues)
	if batchValues["totalUpdatedCells"] == nil {
		t.Fatalf("expected totalUpdatedCells, got %#v", batchValues)
	}
}

func TestSheetsCreateCommand(t *testing.T) {
	requireDriveCreationCapable(t)
	env, gc := setupGWSClient(t)

	createPath := writeTempJSONRequest(t, fmt.Sprintf(`{"properties":{"title":"Mistle GWS Created Sheet %d"}}`, nowNano()))
	createResult, err := runCommandWithInput(t, env, "", "gws", "sheets", "spreadsheets", "create", "--request-file", createPath)
	if err != nil {
		t.Fatal(err)
	}
	var created map[string]any
	decodeCommandJSON(t, createResult, &created)
	createdID, _ := created["spreadsheetId"].(string)
	if createdID == "" {
		t.Fatalf("expected created spreadsheet id, got %#v", created)
	}
	defer func() {
		if err := gc.DeleteDriveFile(cliContext(), createdID); err != nil {
			t.Logf("failed to clean created spreadsheet %s: %v", createdID, err)
		}
	}()
}

func TestDocsCommands(t *testing.T) {
	env, gc := setupGWSClient(t)
	resources := testGWSWorkspaceResources(t)

	getResult, err := runCommandWithInput(t, env, "", "gws", "docs", "documents", "get", "--document-id", resources.DocumentID)
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	decodeCommandJSON(t, getResult, &document)
	if document["documentId"] != resources.DocumentID {
		t.Fatalf("expected document id %s, got %#v", resources.DocumentID, document)
	}

	insertedText := fmt.Sprintf("Hello from gws integration test %d\n", nowNano())
	batchPath := writeTempJSONRequest(t, fmt.Sprintf(`{"requests":[{"insertText":{"location":{"index":1},"text":%q}}]}`, insertedText))
	batchResult, err := runCommandWithInput(t, env, "", "gws", "docs", "documents", "batch-update", "--document-id", resources.DocumentID, "--request-file", batchPath)
	if err != nil {
		t.Fatal(err)
	}
	var batch map[string]any
	decodeCommandJSON(t, batchResult, &batch)
	if batch["replies"] == nil {
		t.Fatalf("expected docs batch update replies, got %#v", batch)
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
			t.Logf("failed to clean inserted document text: %v", err)
		}
	})
}

func TestSlidesCommands(t *testing.T) {
	env, gc := setupGWSClient(t)
	resources := testGWSWorkspaceResources(t)

	getResult, err := runCommandWithInput(t, env, "", "gws", "slides", "presentations", "get", "--presentation-id", resources.PresentationID)
	if err != nil {
		t.Fatal(err)
	}
	var presentation map[string]any
	decodeCommandJSON(t, getResult, &presentation)
	if presentation["presentationId"] != resources.PresentationID {
		t.Fatalf("expected presentation id %s, got %#v", resources.PresentationID, presentation)
	}

	slideObjectID := fmt.Sprintf("slide_%d", nowNano())
	batchPath := writeTempJSONRequest(t, fmt.Sprintf(`{"requests":[{"createSlide":{"objectId":%q,"slideLayoutReference":{"predefinedLayout":"BLANK"}}}]}`, slideObjectID))
	batchResult, err := runCommandWithInput(t, env, "", "gws", "slides", "presentations", "batch-update", "--presentation-id", resources.PresentationID, "--request-file", batchPath)
	if err != nil {
		t.Fatal(err)
	}
	var batch map[string]any
	decodeCommandJSON(t, batchResult, &batch)
	if batch["replies"] == nil {
		t.Fatalf("expected slides batch update replies, got %#v", batch)
	}
	t.Cleanup(func() {
		_, err := gc.BatchUpdatePresentation(cliContext(), resources.PresentationID, map[string]any{
			"requests": []map[string]any{{"deleteObject": map[string]any{"objectId": slideObjectID}}},
		})
		if err != nil {
			t.Logf("failed to clean created slide %s: %v", slideObjectID, err)
		}
	})
}

func TestSlidesCreateCommand(t *testing.T) {
	requireDriveCreationCapable(t)
	env, gc := setupGWSClient(t)

	createPath := writeTempJSONRequest(t, fmt.Sprintf(`{"title":"Mistle GWS Created Slides %d"}`, nowNano()))
	createResult, err := runCommandWithInput(t, env, "", "gws", "slides", "presentations", "create", "--request-file", createPath)
	if err != nil {
		t.Fatal(err)
	}
	var created map[string]any
	decodeCommandJSON(t, createResult, &created)
	createdID, _ := created["presentationId"].(string)
	if createdID == "" {
		t.Fatalf("expected created presentation id, got %#v", created)
	}
	defer func() {
		if err := gc.DeleteDriveFile(cliContext(), createdID); err != nil {
			t.Logf("failed to clean created presentation %s: %v", createdID, err)
		}
	}()
}

func driveFilesContain(files []GWSDriveFile, id string) bool {
	for _, file := range files {
		if file.ID == id {
			return true
		}
	}
	return false
}

func nowNano() int64 {
	return time.Now().UnixNano()
}

func cleanupAddedSheet(t *testing.T, gc GWSClient, spreadsheetID string, batch map[string]any) {
	t.Helper()
	sheetID, ok := addedSheetID(batch)
	if !ok {
		t.Fatalf("expected addSheet reply with sheetId, got %#v", batch)
	}
	t.Cleanup(func() {
		_, err := gc.BatchUpdateSpreadsheet(cliContext(), spreadsheetID, map[string]any{
			"requests": []map[string]any{{"deleteSheet": map[string]any{"sheetId": sheetID}}},
		})
		if err != nil {
			t.Logf("failed to clean added sheet %v: %v", sheetID, err)
		}
	})
}

func addedSheetID(batch map[string]any) (float64, bool) {
	replies, ok := batch["replies"].([]any)
	if !ok || len(replies) == 0 {
		return 0, false
	}
	reply, ok := replies[0].(map[string]any)
	if !ok {
		return 0, false
	}
	addSheet, ok := reply["addSheet"].(map[string]any)
	if !ok {
		return 0, false
	}
	properties, ok := addSheet["properties"].(map[string]any)
	if !ok {
		return 0, false
	}
	sheetID, ok := properties["sheetId"].(float64)
	return sheetID, ok
}
