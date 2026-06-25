# gws

Thin CLI and MCP wrapper for Google Workspace Drive, Sheets, Docs, Slides, Gmail, Calendar, Chat, and People REST APIs.

`gws` does not mint, refresh, store, or inspect credentials. Configure API base URLs and inject Google authorization outside the binary.

## Configuration

Required environment:

- `GWS_DRIVE_BASE_URL`, for example `https://www.googleapis.com/drive/v3`
- `GWS_SHEETS_BASE_URL`, for example `https://sheets.googleapis.com/v4`
- `GWS_DOCS_BASE_URL`, for example `https://docs.googleapis.com/v1`
- `GWS_SLIDES_BASE_URL`, for example `https://slides.googleapis.com/v1`
- `GWS_GMAIL_BASE_URL`, for example `https://gmail.googleapis.com/gmail/v1`
- `GWS_CALENDAR_BASE_URL`, for example `https://www.googleapis.com/calendar/v3`
- `GWS_CHAT_BASE_URL`, for example `https://chat.googleapis.com/v1`
- `GWS_PEOPLE_BASE_URL`, for example `https://people.googleapis.com/v1`

The caller or proxy must inject:

- `Authorization: Bearer <oauth-access-token>`

## Commands

```bash
gws auth test
gws request --api drive --method GET --path /files

gws drive files list --query "'<folder-id>' in parents" --page-size 10
gws drive files get --file-id <id>
gws drive files create --request-file file.json
gws drive files copy --file-id <id> --request-file file.json
gws drive files update --file-id <id> --request-file file.json
gws drive files delete --file-id <id>
gws drive permissions list --file-id <id>
gws drive permissions create --file-id <id> --request-file permission.json
gws drive permissions delete --file-id <id> --permission-id <id>

gws sheets spreadsheets get --spreadsheet-id <id>
gws sheets spreadsheets create --request-file spreadsheet.json
gws sheets spreadsheets batch-update --spreadsheet-id <id> --request-file request.json
gws sheets values get --spreadsheet-id <id> --range Sheet1!A1:B2
gws sheets values update --spreadsheet-id <id> --range Sheet1!A1 --value-input-option RAW --request-file values.json
gws sheets values batch-update --spreadsheet-id <id> --request-file request.json

gws docs documents get --document-id <id>
gws docs documents batch-update --document-id <id> --request-file request.json

gws slides presentations get --presentation-id <id>
gws slides presentations create --request-file presentation.json
gws slides presentations batch-update --presentation-id <id> --request-file request.json

gws gmail messages list --user-id me --max-results 10
gws gmail messages get --user-id me --message-id <id>
gws gmail messages send --user-id me --request-file message.json
gws gmail drafts list --user-id me
gws gmail drafts get --user-id me --draft-id <id>
gws gmail drafts create --user-id me --request-file draft.json
gws gmail drafts send --user-id me --request-file draft-send.json
gws gmail drafts delete --user-id me --draft-id <id>

gws calendar calendar-list list
gws calendar calendar-list get --calendar-id <id>
gws calendar events list --calendar-id primary --max-results 10
gws calendar events get --calendar-id primary --event-id <id>
gws calendar events insert --calendar-id primary --request-file event.json
gws calendar events patch --calendar-id primary --event-id <id> --request-file event.json
gws calendar events delete --calendar-id primary --event-id <id>
gws calendar freebusy query --request-file freebusy.json

gws chat spaces list
gws chat spaces get --space-name spaces/<id>
gws chat messages list --space-name spaces/<id>
gws chat messages get --message-name spaces/<id>/messages/<id>
gws chat messages create --space-name spaces/<id> --request-file message.json
gws chat members list --space-name spaces/<id>

gws people people get --resource-name people/me --person-fields names,emailAddresses
gws people connections list --resource-name people/me --person-fields names,emailAddresses
gws people search-contacts --query thomas --read-mask names,emailAddresses
gws people search-directory --query thomas --read-mask names,emailAddresses

gws mcp serve
gws mcp serve --tools drive,sheets,gmail
```

`gws request` is the complete API coverage surface. Named commands and MCP tools provide progressive discovery for common Workspace workflows.

## MCP

`gws mcp serve` runs a local Streamable HTTP MCP server on `127.0.0.1:7353/mcp` by default.

Use `--tools` to register only selected groups or exact tool names:

```bash
gws mcp serve --tools drive,sheets
gws mcp serve --tools gws_drive_files_list,gws_docs_document_get
```

Supported groups:

- `auth`
- `raw`
- `drive`
- `sheets`
- `docs`
- `slides`
- `gmail`
- `calendar`
- `chat`
- `people`

## Integration Tests

Integration tests use the real Google Workspace APIs through `internal/testproxy` and skip when required environment variables are absent.

Required:

- `GWS_TEST_SERVICE_ACCOUNT_KEY_JSON_BASE64`
- `GWS_TEST_DRIVE_FOLDER_ID`
- `GWS_TEST_DRIVE_FILE_ID`
- `GWS_TEST_DOCUMENT_ID`
- `GWS_TEST_SPREADSHEET_ID`
- `GWS_TEST_PRESENTATION_ID`

Optional:

- `GWS_TEST_USE_WORKSPACE_USER_EMAIL=1`, enables service-account JWT subject impersonation for domain-wide delegation tests.
- `GWS_TEST_WORKSPACE_USER_EMAIL`, required when `GWS_TEST_USE_WORKSPACE_USER_EMAIL=1`.
- `GWS_TEST_RUN_CREATION_TESTS=1`, enables live tests that create and delete Google Workspace files after Drive ownership/storage has been configured for the test principal.
- `GWS_TEST_DRIVE_BASE_URL`, defaults to `https://www.googleapis.com/drive/v3`
- `GWS_TEST_SHEETS_BASE_URL`, defaults to `https://sheets.googleapis.com/v4`
- `GWS_TEST_DOCS_BASE_URL`, defaults to `https://docs.googleapis.com/v1`
- `GWS_TEST_SLIDES_BASE_URL`, defaults to `https://slides.googleapis.com/v1`
- `GWS_TEST_GMAIL_BASE_URL`, defaults to `https://gmail.googleapis.com/gmail/v1`
- `GWS_TEST_CALENDAR_BASE_URL`, defaults to `https://www.googleapis.com/calendar/v3`
- `GWS_TEST_CHAT_BASE_URL`, defaults to `https://chat.googleapis.com/v1`
- `GWS_TEST_PEOPLE_BASE_URL`, defaults to `https://people.googleapis.com/v1`
- `GWS_TEST_GMAIL_USER_ID`, enables optional Gmail list coverage, for example `me`.
- `GWS_TEST_CALENDAR_ID`, enables optional Calendar list/freebusy coverage, for example `primary`.
- `GWS_TEST_CHAT_SPACE_NAME`, enables optional Chat members coverage, for example `spaces/<id>`.
- `GWS_TEST_PEOPLE_RESOURCE_NAME`, enables optional People get/connections coverage, for example `people/me`.

The configured folder and files should be disposable. Standard live tests read and update the shared Drive file, document, spreadsheet, and presentation, then clean up inserted sheet/document/slide content where Google APIs return stable object IDs.

Google Drive file creation depends on the Google-side Drive ownership and storage configuration for the authenticated principal. Plain service accounts can read and update files shared with them, but may not be able to create new Google Workspace files in a regular shared My Drive folder. Tests for `drive files create/copy/delete`, `sheets spreadsheets create`, `slides presentations create`, and MCP create/delete tools are present but opt-in through `GWS_TEST_RUN_CREATION_TESTS=1` so local/CI results are not coupled to a specific Google Drive setup.

Setting `GWS_TEST_WORKSPACE_USER_EMAIL` alone does not change auth behavior. Set `GWS_TEST_USE_WORKSPACE_USER_EMAIL=1` only when the service account is authorized for domain-wide delegation and the test should mint tokens as that delegated subject.
