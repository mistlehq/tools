package main

func (cli CLI) printHelp() {
	_, _ = cli.stdout.Write([]byte(`gws

Thin Google Workspace CLI for Drive, Sheets, Docs, and Slides.

Usage:
  gws help
  gws version
  gws auth help
  gws request help
  gws drive help
  gws sheets help
  gws docs help
  gws slides help
  gws mcp help

Commands:
  auth       Check configured Workspace API access
  request    Send raw Workspace API requests
  drive      Work with Drive files and permissions
  sheets     Work with Sheets spreadsheets and values
  docs       Work with Docs documents
  slides     Work with Slides presentations
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

Options:
  --api <api>             One of drive, sheets, docs, slides.
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
