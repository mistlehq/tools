package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mistlehq/tools/internal/argparse"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultMCPAddr     = "127.0.0.1:7353"
	defaultMCPEndpoint = "/mcp"
)

type gwsMCPConfig struct {
	Addr     string
	Endpoint string
	Tools    map[string]bool
}

type gwsEmptyToolInput struct{}

type gwsRawRequestToolInput struct {
	API    string         `json:"api" jsonschema:"Google Workspace API base to use: drive, sheets, docs, or slides."`
	Method string         `json:"method,omitempty" jsonschema:"HTTP method. Defaults to GET."`
	Path   string         `json:"path" jsonschema:"API path beginning with /."`
	Params map[string]any `json:"params,omitempty" jsonschema:"Query parameters to append to the request."`
	Body   map[string]any `json:"body,omitempty" jsonschema:"JSON object request body."`
}

type gwsFileIDToolInput struct {
	FileID string `json:"fileId" jsonschema:"Google Drive file ID."`
	Fields string `json:"fields,omitempty" jsonschema:"Optional Drive API fields selector."`
}

type gwsDriveFilesListToolInput struct {
	Query    string `json:"query,omitempty" jsonschema:"Optional Drive API q query."`
	PageSize string `json:"pageSize,omitempty" jsonschema:"Optional page size."`
	Fields   string `json:"fields,omitempty" jsonschema:"Optional Drive API fields selector."`
}

type gwsFileBodyToolInput struct {
	FileID  string         `json:"fileId,omitempty" jsonschema:"Google Drive file ID where required."`
	Request map[string]any `json:"request" jsonschema:"Google Workspace API request body using Google's documented JSON shape."`
	Fields  string         `json:"fields,omitempty" jsonschema:"Optional Drive API fields selector."`
}

type gwsPermissionDeleteToolInput struct {
	FileID       string `json:"fileId" jsonschema:"Google Drive file ID."`
	PermissionID string `json:"permissionId" jsonschema:"Google Drive permission ID."`
}

type gwsIDToolInput struct {
	SpreadsheetID  string `json:"spreadsheetId,omitempty" jsonschema:"Google Sheets spreadsheet ID."`
	DocumentID     string `json:"documentId,omitempty" jsonschema:"Google Docs document ID."`
	PresentationID string `json:"presentationId,omitempty" jsonschema:"Google Slides presentation ID."`
}

type gwsIDRequestToolInput struct {
	SpreadsheetID  string         `json:"spreadsheetId,omitempty" jsonschema:"Google Sheets spreadsheet ID."`
	DocumentID     string         `json:"documentId,omitempty" jsonschema:"Google Docs document ID."`
	PresentationID string         `json:"presentationId,omitempty" jsonschema:"Google Slides presentation ID."`
	Request        map[string]any `json:"request" jsonschema:"Google Workspace API request body using Google's documented JSON shape."`
}

type gwsSheetsValuesGetToolInput struct {
	SpreadsheetID string `json:"spreadsheetId" jsonschema:"Google Sheets spreadsheet ID."`
	Range         string `json:"range" jsonschema:"A1 notation range."`
}

type gwsSheetsValuesUpdateToolInput struct {
	SpreadsheetID    string         `json:"spreadsheetId" jsonschema:"Google Sheets spreadsheet ID."`
	Range            string         `json:"range" jsonschema:"A1 notation range."`
	ValueInputOption string         `json:"valueInputOption" jsonschema:"Google Sheets value input option, for example RAW or USER_ENTERED."`
	Request          map[string]any `json:"request" jsonschema:"Google Sheets ValueRange request body."`
}

func (cli CLI) runMCP(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printMCPHelp()
		return nil
	}
	switch args[0] {
	case "serve":
		if isSingleHelpArg(args[1:]) {
			cli.printMCPServeHelp()
			return nil
		}
		return cli.runMCPServe(args[1:])
	default:
		return fmt.Errorf("unsupported mcp command: %s", args[0])
	}
}

func (cli CLI) runMCPServe(args []string) error {
	config, err := parseGWSMCPServeArgs(args)
	if err != nil {
		return err
	}
	gc, err := cli.gwsClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newGWSMCPHTTPHandler(gc, config.Tools))
	fmt.Fprintf(cli.stdout, "Google Workspace MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseGWSMCPServeArgs(args []string) (gwsMCPConfig, error) {
	parsed, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
		"tools":    {TakesValue: true},
	})
	if err != nil {
		return gwsMCPConfig{}, err
	}
	if len(parsed.Positionals) > 0 {
		return gwsMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}
	config := gwsMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsed.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsed.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if tools := parsed.First("tools"); tools != "" {
		parsedTools, err := parseToolFilter(tools)
		if err != nil {
			return gwsMCPConfig{}, err
		}
		config.Tools = parsedTools
	}
	if strings.TrimSpace(config.Addr) == "" {
		return gwsMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return gwsMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func parseToolFilter(raw string) (map[string]bool, error) {
	out := make(map[string]bool)
	for _, part := range strings.Split(raw, ",") {
		value := strings.TrimSpace(part)
		if value == "" {
			return nil, fmt.Errorf("--tools contains an empty entry")
		}
		out[value] = true
	}
	return out, nil
}

func newGWSMCPHTTPHandler(gc GWSClient, tools map[string]bool) http.Handler {
	server := newGWSMCPServer(gc, tools)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newGWSMCPServer(gc GWSClient, tools map[string]bool) *mcp.Server {
	openWorld := true
	readOnly := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	mutating := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: boolPtr(false)}
	destructive := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: boolPtr(true)}
	raw := &mcp.ToolAnnotations{OpenWorldHint: &openWorld}
	server := mcp.NewServer(&mcp.Implementation{Name: "gws", Version: Version}, nil)

	addGWSTool(server, tools, "raw", gwsTool("gws_request", gwsRequestDoc, raw), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsRawRequestToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		body, err := gc.RequestContext(ctx, GWSRequest{API: input.API, Method: input.Method, Path: input.Path, Params: input.Params, Body: input.Body})
		if err != nil {
			return nil, nil, err
		}
		var out GWSRawResult
		if len(body) > 0 {
			if err := json.Unmarshal(body, &out); err != nil {
				return nil, nil, err
			}
		}
		return nil, out, nil
	})
	addGWSTool(server, tools, "auth", gwsTool("gws_auth_test", gwsAuthTestDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, _ *gwsEmptyToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.AuthTest(ctx)
		return nil, out, err
	})
	addGWSTool(server, tools, "drive", gwsTool("gws_drive_files_list", gwsDriveFilesListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsDriveFilesListToolInput) (*mcp.CallToolResult, GWSDriveFilesList, error) {
		params := map[string]any{}
		if input.Query != "" {
			params["q"] = input.Query
		}
		if input.PageSize != "" {
			params["pageSize"] = input.PageSize
		}
		if input.Fields != "" {
			params["fields"] = input.Fields
		}
		out, err := gc.ListDriveFiles(ctx, params)
		return nil, out, err
	})
	addGWSTool(server, tools, "drive", gwsTool("gws_drive_file_get", gwsDriveFileGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsFileIDToolInput) (*mcp.CallToolResult, GWSDriveFile, error) {
		out, err := gc.GetDriveFile(ctx, input.FileID, fieldsParam(input.Fields))
		return nil, out, err
	})
	addGWSTool(server, tools, "drive", gwsTool("gws_drive_file_create", gwsDriveFileCreateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsFileBodyToolInput) (*mcp.CallToolResult, GWSDriveFile, error) {
		out, err := gc.CreateDriveFile(ctx, input.Request, fieldsParam(input.Fields))
		return nil, out, err
	})
	addGWSTool(server, tools, "drive", gwsTool("gws_drive_file_copy", gwsDriveFileCopyDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsFileBodyToolInput) (*mcp.CallToolResult, GWSDriveFile, error) {
		out, err := gc.CopyDriveFile(ctx, input.FileID, input.Request, fieldsParam(input.Fields))
		return nil, out, err
	})
	addGWSTool(server, tools, "drive", gwsTool("gws_drive_file_update", gwsDriveFileUpdateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsFileBodyToolInput) (*mcp.CallToolResult, GWSDriveFile, error) {
		out, err := gc.UpdateDriveFile(ctx, input.FileID, input.Request, fieldsParam(input.Fields))
		return nil, out, err
	})
	addGWSTool(server, tools, "drive", gwsTool("gws_drive_file_delete", gwsDriveFileDeleteDoc, destructive), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsFileIDToolInput) (*mcp.CallToolResult, map[string]any, error) {
		err := gc.DeleteDriveFile(ctx, input.FileID)
		return nil, map[string]any{"deleted": err == nil, "fileId": input.FileID}, err
	})
	addGWSTool(server, tools, "drive", gwsTool("gws_drive_permissions_list", gwsDrivePermissionsListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsFileIDToolInput) (*mcp.CallToolResult, GWSDrivePermissionsList, error) {
		out, err := gc.ListDrivePermissions(ctx, input.FileID)
		return nil, out, err
	})
	addGWSTool(server, tools, "drive", gwsTool("gws_drive_permission_create", gwsDrivePermissionCreateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsFileBodyToolInput) (*mcp.CallToolResult, GWSDrivePermission, error) {
		out, err := gc.CreateDrivePermission(ctx, input.FileID, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "drive", gwsTool("gws_drive_permission_delete", gwsDrivePermissionDeleteDoc, destructive), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsPermissionDeleteToolInput) (*mcp.CallToolResult, map[string]any, error) {
		err := gc.DeleteDrivePermission(ctx, input.FileID, input.PermissionID)
		return nil, map[string]any{"deleted": err == nil, "fileId": input.FileID, "permissionId": input.PermissionID}, err
	})
	addGWSSheetsTools(server, tools, gc, readOnly, mutating)
	addGWSDocsTools(server, tools, gc, readOnly, mutating)
	addGWSSlidesTools(server, tools, gc, readOnly, mutating)
	return server
}

func addGWSSheetsTools(server *mcp.Server, tools map[string]bool, gc GWSClient, readOnly *mcp.ToolAnnotations, mutating *mcp.ToolAnnotations) {
	addGWSTool(server, tools, "sheets", gwsTool("gws_sheets_spreadsheet_get", gwsSheetsSpreadsheetGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsIDToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetSpreadsheet(ctx, input.SpreadsheetID)
		return nil, out, err
	})
	addGWSTool(server, tools, "sheets", gwsTool("gws_sheets_spreadsheet_create", gwsSheetsSpreadsheetCreateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsIDRequestToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.CreateSpreadsheet(ctx, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "sheets", gwsTool("gws_sheets_spreadsheet_batch_update", gwsSheetsSpreadsheetBatchUpdateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsIDRequestToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.BatchUpdateSpreadsheet(ctx, input.SpreadsheetID, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "sheets", gwsTool("gws_sheets_values_get", gwsSheetsValuesGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsSheetsValuesGetToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetSpreadsheetValues(ctx, input.SpreadsheetID, input.Range)
		return nil, out, err
	})
	addGWSTool(server, tools, "sheets", gwsTool("gws_sheets_values_update", gwsSheetsValuesUpdateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsSheetsValuesUpdateToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.UpdateSpreadsheetValues(ctx, input.SpreadsheetID, input.Range, input.ValueInputOption, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "sheets", gwsTool("gws_sheets_values_batch_update", gwsSheetsValuesBatchUpdateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsIDRequestToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.BatchUpdateSpreadsheetValues(ctx, input.SpreadsheetID, input.Request)
		return nil, out, err
	})
}

func addGWSDocsTools(server *mcp.Server, tools map[string]bool, gc GWSClient, readOnly *mcp.ToolAnnotations, mutating *mcp.ToolAnnotations) {
	addGWSTool(server, tools, "docs", gwsTool("gws_docs_document_get", gwsDocsDocumentGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsIDToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetDocument(ctx, input.DocumentID)
		return nil, out, err
	})
	addGWSTool(server, tools, "docs", gwsTool("gws_docs_document_batch_update", gwsDocsDocumentBatchUpdateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsIDRequestToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.BatchUpdateDocument(ctx, input.DocumentID, input.Request)
		return nil, out, err
	})
}

func addGWSSlidesTools(server *mcp.Server, tools map[string]bool, gc GWSClient, readOnly *mcp.ToolAnnotations, mutating *mcp.ToolAnnotations) {
	addGWSTool(server, tools, "slides", gwsTool("gws_slides_presentation_get", gwsSlidesPresentationGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsIDToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetPresentation(ctx, input.PresentationID)
		return nil, out, err
	})
	addGWSTool(server, tools, "slides", gwsTool("gws_slides_presentation_create", gwsSlidesPresentationCreateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsIDRequestToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.CreatePresentation(ctx, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "slides", gwsTool("gws_slides_presentation_batch_update", gwsSlidesPresentationBatchUpdateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsIDRequestToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.BatchUpdatePresentation(ctx, input.PresentationID, input.Request)
		return nil, out, err
	})
}

func addGWSTool[I any, O any](server *mcp.Server, tools map[string]bool, group string, tool *mcp.Tool, handler func(context.Context, *mcp.CallToolRequest, *I) (*mcp.CallToolResult, O, error)) {
	if !toolSelected(tools, group, tool.Name) {
		return
	}
	mcp.AddTool(server, tool, handler)
}

func toolSelected(tools map[string]bool, group string, name string) bool {
	if len(tools) == 0 {
		return true
	}
	return tools[group] || tools[name]
}

func gwsTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
	return &mcp.Tool{Name: name, Title: doc.Command, Description: doc.Description, Annotations: annotations}
}

func fieldsParam(fields string) map[string]any {
	if fields == "" {
		return nil
	}
	return map[string]any{"fields": fields}
}

func boolPtr(value bool) *bool {
	return &value
}

func (cli CLI) printMCPHelp() {
	_, _ = cli.stdout.Write([]byte(`gws mcp

Run Google Workspace as a local MCP server.

Usage:
  gws mcp help
  gws mcp serve
  gws mcp serve --help

Commands:
  serve    Serve Google Workspace MCP tools over Streamable HTTP
`))
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `gws mcp serve

Serve Google Workspace tools over MCP Streamable HTTP.

Usage:
  gws mcp serve
  gws mcp serve --addr <addr>
  gws mcp serve --endpoint <path>
  gws mcp serve --tools drive,sheets
  gws mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.
  --tools <list>       Comma-separated tool groups or tool names. Groups: auth, raw, drive, sheets, docs, slides.

Tools:
  gws_request                            %s
  gws_auth_test                          %s
  gws_drive_files_list                   %s
  gws_drive_file_get                     %s
  gws_drive_file_create                  %s
  gws_drive_file_copy                    %s
  gws_drive_file_update                  %s
  gws_drive_file_delete                  %s
  gws_drive_permissions_list             %s
  gws_drive_permission_create            %s
  gws_drive_permission_delete            %s
  gws_sheets_spreadsheet_get             %s
  gws_sheets_spreadsheet_create          %s
  gws_sheets_spreadsheet_batch_update    %s
  gws_sheets_values_get                  %s
  gws_sheets_values_update               %s
  gws_sheets_values_batch_update         %s
  gws_docs_document_get                  %s
  gws_docs_document_batch_update         %s
  gws_slides_presentation_get            %s
  gws_slides_presentation_create         %s
  gws_slides_presentation_batch_update   %s
`, defaultMCPAddr, defaultMCPEndpoint, gwsRequestDoc.Summary, gwsAuthTestDoc.Summary, gwsDriveFilesListDoc.Summary, gwsDriveFileGetDoc.Summary, gwsDriveFileCreateDoc.Summary, gwsDriveFileCopyDoc.Summary, gwsDriveFileUpdateDoc.Summary, gwsDriveFileDeleteDoc.Summary, gwsDrivePermissionsListDoc.Summary, gwsDrivePermissionCreateDoc.Summary, gwsDrivePermissionDeleteDoc.Summary, gwsSheetsSpreadsheetGetDoc.Summary, gwsSheetsSpreadsheetCreateDoc.Summary, gwsSheetsSpreadsheetBatchUpdateDoc.Summary, gwsSheetsValuesGetDoc.Summary, gwsSheetsValuesUpdateDoc.Summary, gwsSheetsValuesBatchUpdateDoc.Summary, gwsDocsDocumentGetDoc.Summary, gwsDocsDocumentBatchUpdateDoc.Summary, gwsSlidesPresentationGetDoc.Summary, gwsSlidesPresentationCreateDoc.Summary, gwsSlidesPresentationBatchUpdateDoc.Summary)
}
