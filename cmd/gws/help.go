package main

func (cli CLI) printHelp() {
	_, _ = cli.stdout.Write([]byte(`gws

Thin Google Workspace CLI for Drive, Sheets, Docs, Slides, Gmail, Calendar, Chat, and People APIs.

Usage:
  gws help
  gws version
  gws auth help
  gws request help
  gws drive help
  gws sheets help
  gws docs help
  gws slides help
  gws gmail help
  gws calendar help
  gws chat help
  gws people help
  gws mcp help

Commands:
  auth       Check configured Workspace API access
  request    Send raw Workspace API requests
  drive      Work with Drive files and permissions
  sheets     Work with Sheets spreadsheets and values
  docs       Work with Docs documents
  slides     Work with Slides presentations
  gmail      Work with Gmail messages and drafts
  calendar   Work with Calendar lists, events, and freebusy
  chat       Work with Chat spaces, messages, and members
  people     Work with People profiles, connections, and search
  mcp        Serve Workspace tools over MCP Streamable HTTP
`))
}

func (cli CLI) printAuthHelp() {
	_, _ = cli.stdout.Write([]byte(`gws auth

Inspect Google Workspace API access.

Usage:
  gws auth help
  gws auth test
`))
}

func (cli CLI) printAuthTestHelp() {
	_, _ = cli.stdout.Write([]byte(`gws auth test

Check Google Workspace API access by fetching Drive account metadata.

Usage:
  gws auth test
`))
}

func (cli CLI) printRequestHelp() {
	_, _ = cli.stdout.Write([]byte(`gws request

Send a raw request to one configured Google Workspace API base URL.

Usage:
  gws request --api drive --method GET --path /files
  gws request --api sheets --method GET --path /spreadsheets/<id>
  gws request --api docs --method POST --path /documents/<id>:batchUpdate --request-file <json>
  gws request --api slides --method POST --path /presentations/<id>:batchUpdate --body <json>
  gws request --api gmail --method GET --path /users/me/messages

Options:
  --api <api>             One of drive, sheets, docs, slides, gmail, calendar, chat, people.
  --method <method>       GET, POST, PATCH, PUT, or DELETE. Defaults to GET.
  --path <path>           API path beginning with /.
  --body <json>           JSON object request body.
  --request-file <json>   JSON object request body file.
`))
}

func (cli CLI) printDriveHelp() {
	_, _ = cli.stdout.Write([]byte(`gws drive

Work with Google Drive files and permissions.

Usage:
  gws drive help
  gws drive files help
  gws drive permissions help
`))
}

func (cli CLI) printDriveFilesHelp() {
	_, _ = cli.stdout.Write([]byte(`gws drive files

Work with Google Drive files.

Usage:
  gws drive files list [--query <q>] [--page-size <n>] [--fields <fields>]
  gws drive files get --file-id <id> [--fields <fields>]
  gws drive files create --request-file <json> [--fields <fields>]
  gws drive files copy --file-id <id> --request-file <json> [--fields <fields>]
  gws drive files update --file-id <id> --request-file <json> [--fields <fields>]
  gws drive files delete --file-id <id>
`))
}

func (cli CLI) printDrivePermissionsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws drive permissions

Work with Google Drive file permissions.

Usage:
  gws drive permissions list --file-id <id>
  gws drive permissions create --file-id <id> --request-file <json>
  gws drive permissions delete --file-id <id> --permission-id <id>
`))
}

func (cli CLI) printSheetsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws sheets

Work with Google Sheets spreadsheets and values.

Usage:
  gws sheets help
  gws sheets spreadsheets help
  gws sheets values help
`))
}

func (cli CLI) printSheetsSpreadsheetsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws sheets spreadsheets

Work with Google Sheets spreadsheets.

Usage:
  gws sheets spreadsheets get --spreadsheet-id <id>
  gws sheets spreadsheets create --request-file <json>
  gws sheets spreadsheets batch-update --spreadsheet-id <id> --request-file <json>
`))
}

func (cli CLI) printSheetsValuesHelp() {
	_, _ = cli.stdout.Write([]byte(`gws sheets values

Work with Google Sheets values.

Usage:
  gws sheets values get --spreadsheet-id <id> --range <a1>
  gws sheets values update --spreadsheet-id <id> --range <a1> --value-input-option RAW --request-file <json>
  gws sheets values batch-update --spreadsheet-id <id> --request-file <json>
`))
}

func (cli CLI) printDocsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws docs

Work with Google Docs documents.

Usage:
  gws docs help
  gws docs documents help
`))
}

func (cli CLI) printDocsDocumentsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws docs documents

Work with Google Docs documents.

Usage:
  gws docs documents get --document-id <id>
  gws docs documents batch-update --document-id <id> --request-file <json>
`))
}

func (cli CLI) printSlidesHelp() {
	_, _ = cli.stdout.Write([]byte(`gws slides

Work with Google Slides presentations.

Usage:
  gws slides help
  gws slides presentations help
`))
}

func (cli CLI) printSlidesPresentationsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws slides presentations

Work with Google Slides presentations.

Usage:
  gws slides presentations get --presentation-id <id>
  gws slides presentations create --request-file <json>
  gws slides presentations batch-update --presentation-id <id> --request-file <json>
`))
}

func (cli CLI) printGmailHelp() {
	_, _ = cli.stdout.Write([]byte(`gws gmail

Work with Gmail messages and drafts.

Usage:
  gws gmail help
  gws gmail messages help
  gws gmail drafts help
`))
}

func (cli CLI) printGmailMessagesHelp() {
	_, _ = cli.stdout.Write([]byte(`gws gmail messages

Work with Gmail messages.

Usage:
  gws gmail messages list --user-id me [--query <q>] [--label-ids <ids>] [--max-results <n>] [--page-token <token>]
  gws gmail messages get --user-id me --message-id <id> [--format <format>]
  gws gmail messages send --user-id me --request-file <json>
`))
}

func (cli CLI) printGmailDraftsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws gmail drafts

Work with Gmail drafts.

Usage:
  gws gmail drafts list --user-id me [--max-results <n>] [--page-token <token>]
  gws gmail drafts get --user-id me --draft-id <id> [--format <format>]
  gws gmail drafts create --user-id me --request-file <json>
  gws gmail drafts send --user-id me --request-file <json>
  gws gmail drafts delete --user-id me --draft-id <id>
`))
}

func (cli CLI) printCalendarHelp() {
	_, _ = cli.stdout.Write([]byte(`gws calendar

Work with Google Calendar lists, events, and freebusy.

Usage:
  gws calendar help
  gws calendar calendar-list help
  gws calendar events help
  gws calendar freebusy help
`))
}

func (cli CLI) printCalendarListHelp() {
	_, _ = cli.stdout.Write([]byte(`gws calendar calendar-list

Work with Google Calendar calendar-list entries.

Usage:
  gws calendar calendar-list list [--max-results <n>] [--page-token <token>]
  gws calendar calendar-list get --calendar-id <id>
`))
}

func (cli CLI) printCalendarEventsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws calendar events

Work with Google Calendar events.

Usage:
  gws calendar events list --calendar-id <id> [--time-min <rfc3339>] [--time-max <rfc3339>] [--max-results <n>] [--single-events <bool>] [--order-by <field>] [--page-token <token>]
  gws calendar events get --calendar-id <id> --event-id <id>
  gws calendar events insert --calendar-id <id> --request-file <json>
  gws calendar events patch --calendar-id <id> --event-id <id> --request-file <json>
  gws calendar events delete --calendar-id <id> --event-id <id>
`))
}

func (cli CLI) printCalendarFreeBusyHelp() {
	_, _ = cli.stdout.Write([]byte(`gws calendar freebusy

Query Google Calendar free/busy information.

Usage:
  gws calendar freebusy query --request-file <json>
`))
}

func (cli CLI) printChatHelp() {
	_, _ = cli.stdout.Write([]byte(`gws chat

Work with Google Chat spaces, messages, and members.

Usage:
  gws chat help
  gws chat spaces help
  gws chat messages help
  gws chat members help
`))
}

func (cli CLI) printChatSpacesHelp() {
	_, _ = cli.stdout.Write([]byte(`gws chat spaces

Work with Google Chat spaces.

Usage:
  gws chat spaces list [--page-size <n>] [--page-token <token>]
  gws chat spaces get --space-name spaces/<id>
`))
}

func (cli CLI) printChatMessagesHelp() {
	_, _ = cli.stdout.Write([]byte(`gws chat messages

Work with Google Chat messages.

Usage:
  gws chat messages list --space-name spaces/<id> [--page-size <n>] [--page-token <token>]
  gws chat messages get --message-name spaces/<id>/messages/<id>
  gws chat messages create --space-name spaces/<id> --request-file <json>
`))
}

func (cli CLI) printChatMembersHelp() {
	_, _ = cli.stdout.Write([]byte(`gws chat members

Work with Google Chat members.

Usage:
  gws chat members list --space-name spaces/<id> [--page-size <n>] [--page-token <token>]
`))
}

func (cli CLI) printPeopleHelp() {
	_, _ = cli.stdout.Write([]byte(`gws people

Work with Google People profiles, connections, and search.

Usage:
  gws people help
  gws people people help
  gws people connections help
  gws people search-contacts --query <query> --read-mask <fields>
  gws people search-directory --query <query> --read-mask <fields>
`))
}

func (cli CLI) printPeoplePeopleHelp() {
	_, _ = cli.stdout.Write([]byte(`gws people people

Work with People API person resources.

Usage:
  gws people people get --resource-name people/me --person-fields <fields>
`))
}

func (cli CLI) printPeopleConnectionsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws people connections

Work with People API connections.

Usage:
  gws people connections list --resource-name people/me --person-fields <fields> [--page-size <n>] [--page-token <token>]
`))
}

func (cli CLI) printPeopleSearchContactsHelp() {
	_, _ = cli.stdout.Write([]byte(`gws people search-contacts

Search authenticated user's contacts.

Usage:
  gws people search-contacts --query <query> --read-mask <fields> [--page-size <n>]
`))
}

func (cli CLI) printPeopleSearchDirectoryHelp() {
	_, _ = cli.stdout.Write([]byte(`gws people search-directory

Search Google Workspace directory people.

Usage:
  gws people search-directory --query <query> --read-mask <fields> [--sources <sources>] [--page-size <n>]
`))
}
