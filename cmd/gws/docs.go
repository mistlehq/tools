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

var gwsGmailMessagesListDoc = commandDoc{Command: "gws gmail messages list", Summary: "List Gmail messages", Description: "List Gmail messages for a user."}
var gwsGmailMessageGetDoc = commandDoc{Command: "gws gmail messages get", Summary: "Get a Gmail message", Description: "Get a Gmail message for a user."}
var gwsGmailMessageSendDoc = commandDoc{Command: "gws gmail messages send", Summary: "Send a Gmail message", Description: "Send a Gmail message using Google's documented message shape."}
var gwsGmailDraftsListDoc = commandDoc{Command: "gws gmail drafts list", Summary: "List Gmail drafts", Description: "List Gmail drafts for a user."}
var gwsGmailDraftGetDoc = commandDoc{Command: "gws gmail drafts get", Summary: "Get a Gmail draft", Description: "Get a Gmail draft for a user."}
var gwsGmailDraftCreateDoc = commandDoc{Command: "gws gmail drafts create", Summary: "Create a Gmail draft", Description: "Create a Gmail draft using Google's documented draft shape."}
var gwsGmailDraftSendDoc = commandDoc{Command: "gws gmail drafts send", Summary: "Send a Gmail draft", Description: "Send a Gmail draft using Google's documented draft send shape."}
var gwsGmailDraftDeleteDoc = commandDoc{Command: "gws gmail drafts delete", Summary: "Delete a Gmail draft", Description: "Delete a Gmail draft for a user."}

var gwsCalendarListListDoc = commandDoc{Command: "gws calendar calendar-list list", Summary: "List calendars", Description: "List Google Calendar calendar-list entries."}
var gwsCalendarListGetDoc = commandDoc{Command: "gws calendar calendar-list get", Summary: "Get a calendar-list entry", Description: "Get a Google Calendar calendar-list entry."}
var gwsCalendarEventsListDoc = commandDoc{Command: "gws calendar events list", Summary: "List calendar events", Description: "List events for a Google Calendar calendar."}
var gwsCalendarEventGetDoc = commandDoc{Command: "gws calendar events get", Summary: "Get a calendar event", Description: "Get an event from a Google Calendar calendar."}
var gwsCalendarEventInsertDoc = commandDoc{Command: "gws calendar events insert", Summary: "Insert a calendar event", Description: "Insert a Google Calendar event using Google's documented event shape."}
var gwsCalendarEventPatchDoc = commandDoc{Command: "gws calendar events patch", Summary: "Patch a calendar event", Description: "Patch a Google Calendar event using Google's documented event shape."}
var gwsCalendarEventDeleteDoc = commandDoc{Command: "gws calendar events delete", Summary: "Delete a calendar event", Description: "Delete a Google Calendar event."}
var gwsCalendarFreeBusyQueryDoc = commandDoc{Command: "gws calendar freebusy query", Summary: "Query calendar freebusy", Description: "Query Google Calendar free/busy information using Google's documented request shape."}

var gwsChatSpacesListDoc = commandDoc{Command: "gws chat spaces list", Summary: "List Chat spaces", Description: "List Google Chat spaces."}
var gwsChatSpaceGetDoc = commandDoc{Command: "gws chat spaces get", Summary: "Get a Chat space", Description: "Get a Google Chat space."}
var gwsChatMessagesListDoc = commandDoc{Command: "gws chat messages list", Summary: "List Chat messages", Description: "List Google Chat messages in a space."}
var gwsChatMessageGetDoc = commandDoc{Command: "gws chat messages get", Summary: "Get a Chat message", Description: "Get a Google Chat message."}
var gwsChatMessageCreateDoc = commandDoc{Command: "gws chat messages create", Summary: "Create a Chat message", Description: "Create a Google Chat message using Google's documented message shape."}
var gwsChatMembersListDoc = commandDoc{Command: "gws chat members list", Summary: "List Chat members", Description: "List Google Chat members in a space."}

var gwsPeoplePersonGetDoc = commandDoc{Command: "gws people people get", Summary: "Get a person", Description: "Get a People API person resource."}
var gwsPeopleConnectionsListDoc = commandDoc{Command: "gws people connections list", Summary: "List people connections", Description: "List People API connections for a person resource."}
var gwsPeopleSearchContactsDoc = commandDoc{Command: "gws people search-contacts", Summary: "Search contacts", Description: "Search authenticated user's contacts using People API."}
var gwsPeopleSearchDirectoryDoc = commandDoc{Command: "gws people search-directory", Summary: "Search directory people", Description: "Search Google Workspace directory people using People API."}
