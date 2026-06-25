package main

type commandDoc struct {
	Command     string
	Summary     string
	Description string
}

var gwsRequestDoc = commandDoc{
	Command:     "gws request",
	Summary:     "Send a raw Google Workspace API request",
	Description: "Send a raw request to one configured Google Workspace API base URL.",
}

var gwsAuthTestDoc = commandDoc{
	Command:     "gws auth test",
	Summary:     "Check Google Workspace API access",
	Description: "Check Google Workspace API access by fetching Drive account metadata.",
}

var gwsDriveFilesListDoc = commandDoc{
	Command:     "gws drive files list",
	Summary:     "List Drive files",
	Description: "List Google Drive files using Drive API query parameters.",
}

var gwsDriveFileGetDoc = commandDoc{
	Command:     "gws drive files get",
	Summary:     "Get a Drive file",
	Description: "Get metadata for a Google Drive file.",
}

var gwsDriveFileCreateDoc = commandDoc{
	Command:     "gws drive files create",
	Summary:     "Create a Drive file",
	Description: "Create a Google Drive file using Google's documented file metadata shape.",
}

var gwsDriveFileCopyDoc = commandDoc{
	Command:     "gws drive files copy",
	Summary:     "Copy a Drive file",
	Description: "Copy a Google Drive file using Google's documented file metadata shape.",
}

var gwsDriveFileUpdateDoc = commandDoc{
	Command:     "gws drive files update",
	Summary:     "Update a Drive file",
	Description: "Update Google Drive file metadata using Google's documented file metadata shape.",
}

var gwsDriveFileDeleteDoc = commandDoc{
	Command:     "gws drive files delete",
	Summary:     "Delete a Drive file",
	Description: "Delete a Google Drive file.",
}

var gwsDrivePermissionsListDoc = commandDoc{
	Command:     "gws drive permissions list",
	Summary:     "List Drive permissions",
	Description: "List permissions for a Google Drive file.",
}

var gwsDrivePermissionCreateDoc = commandDoc{
	Command:     "gws drive permissions create",
	Summary:     "Create a Drive permission",
	Description: "Create a permission for a Google Drive file using Google's documented permission shape.",
}

var gwsDrivePermissionDeleteDoc = commandDoc{
	Command:     "gws drive permissions delete",
	Summary:     "Delete a Drive permission",
	Description: "Delete a permission from a Google Drive file.",
}

var gwsSheetsSpreadsheetGetDoc = commandDoc{
	Command:     "gws sheets spreadsheets get",
	Summary:     "Get a spreadsheet",
	Description: "Get a Google Sheets spreadsheet.",
}

var gwsSheetsSpreadsheetCreateDoc = commandDoc{
	Command:     "gws sheets spreadsheets create",
	Summary:     "Create a spreadsheet",
	Description: "Create a Google Sheets spreadsheet using Google's documented spreadsheet shape.",
}

var gwsSheetsSpreadsheetBatchUpdateDoc = commandDoc{
	Command:     "gws sheets spreadsheets batch-update",
	Summary:     "Batch update a spreadsheet",
	Description: "Batch update a Google Sheets spreadsheet using Google's documented request shape.",
}

var gwsSheetsValuesGetDoc = commandDoc{
	Command:     "gws sheets values get",
	Summary:     "Get spreadsheet values",
	Description: "Get values from a Google Sheets spreadsheet range.",
}

var gwsSheetsValuesUpdateDoc = commandDoc{
	Command:     "gws sheets values update",
	Summary:     "Update spreadsheet values",
	Description: "Update values in a Google Sheets spreadsheet range using Google's documented value range shape.",
}

var gwsSheetsValuesBatchUpdateDoc = commandDoc{
	Command:     "gws sheets values batch-update",
	Summary:     "Batch update spreadsheet values",
	Description: "Batch update values in a Google Sheets spreadsheet using Google's documented request shape.",
}

var gwsDocsDocumentGetDoc = commandDoc{
	Command:     "gws docs documents get",
	Summary:     "Get a document",
	Description: "Get a Google Docs document.",
}

var gwsDocsDocumentBatchUpdateDoc = commandDoc{
	Command:     "gws docs documents batch-update",
	Summary:     "Batch update a document",
	Description: "Batch update a Google Docs document using Google's documented request shape.",
}

var gwsSlidesPresentationGetDoc = commandDoc{
	Command:     "gws slides presentations get",
	Summary:     "Get a presentation",
	Description: "Get a Google Slides presentation.",
}

var gwsSlidesPresentationCreateDoc = commandDoc{
	Command:     "gws slides presentations create",
	Summary:     "Create a presentation",
	Description: "Create a Google Slides presentation using Google's documented presentation shape.",
}

var gwsSlidesPresentationBatchUpdateDoc = commandDoc{
	Command:     "gws slides presentations batch-update",
	Summary:     "Batch update a presentation",
	Description: "Batch update a Google Slides presentation using Google's documented request shape.",
}
