package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mistlehq/tools/internal/argparse"
)

var Version = "dev"

type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	env    Environment
}

func main() {
	cli := CLI{stdin: os.Stdin, stdout: os.Stdout, stderr: os.Stderr, env: loadEnvironment()}
	if err := cli.run(os.Args); err != nil {
		fmt.Fprintln(cli.stderr, err)
		os.Exit(1)
	}
}

func (cli CLI) gwsClient() (GWSClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return GWSClient{}, err
	}
	return NewGWSClient(config), nil
}

func cliContext() context.Context {
	return context.Background()
}

func isHelpToken(arg string) bool {
	return arg == "help" || arg == "-h" || arg == "--help"
}

func isSingleHelpArg(args []string) bool {
	return len(args) == 1 && isHelpToken(args[0])
}

func (cli CLI) run(args []string) error {
	if len(args) < 2 {
		cli.printHelp()
		return nil
	}
	switch args[1] {
	case "help", "-h", "--help":
		cli.printHelp()
		return nil
	case "version", "--version":
		fmt.Fprintln(cli.stdout, Version)
		return nil
	case "auth":
		return cli.runAuth(args[2:])
	case "request":
		return cli.runRequest(args[2:])
	case "drive":
		return cli.runDrive(args[2:])
	case "sheets":
		return cli.runSheets(args[2:])
	case "docs":
		return cli.runDocs(args[2:])
	case "slides":
		return cli.runSlides(args[2:])
	case "gmail":
		return cli.runGmail(args[2:])
	case "calendar":
		return cli.runCalendar(args[2:])
	case "chat":
		return cli.runChat(args[2:])
	case "people":
		return cli.runPeople(args[2:])
	case "mcp":
		return cli.runMCP(args[2:])
	default:
		return fmt.Errorf("unsupported command: %s", args[1])
	}
}

func (cli CLI) runAuth(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printAuthHelp()
		return nil
	}
	if args[0] != "test" {
		return fmt.Errorf("unsupported auth command: %s", args[0])
	}
	if isSingleHelpArg(args[1:]) {
		cli.printAuthTestHelp()
		return nil
	}
	if len(args[1:]) > 0 {
		return fmt.Errorf("auth test does not accept positional arguments")
	}
	gc, err := cli.gwsClient()
	if err != nil {
		return err
	}
	out, err := gc.AuthTest(cliContext())
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runRequest(args []string) error {
	if len(args) > 0 && isHelpToken(args[0]) {
		cli.printRequestHelp()
		return nil
	}
	request, err := parseRawRequestArgs(args)
	if err != nil {
		return err
	}
	gc, err := cli.gwsClient()
	if err != nil {
		return err
	}
	body, err := gc.Request(request)
	if err != nil {
		return err
	}
	_, err = cli.stdout.Write(append(body, '\n'))
	return err
}

func (cli CLI) runDrive(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printDriveHelp()
		return nil
	}
	switch args[0] {
	case "files":
		if len(args[1:]) == 0 || isHelpToken(args[1]) {
			cli.printDriveFilesHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runDriveFiles(gc, args[1:])
	case "permissions":
		if len(args[1:]) == 0 || isHelpToken(args[1]) {
			cli.printDrivePermissionsHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runDrivePermissions(gc, args[1:])
	default:
		return fmt.Errorf("unsupported drive command: %s", args[0])
	}
}

func (cli CLI) runDriveFiles(gc GWSClient, args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printDriveFilesHelp()
		return nil
	}
	if isSingleHelpArg(args[1:]) {
		cli.printDriveFilesHelp()
		return nil
	}
	switch args[0] {
	case "list":
		parsed, err := parseDriveFilesListArgs(args[1:])
		if err != nil {
			return err
		}
		out, err := gc.ListDriveFiles(cliContext(), parsed)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		fileID, params, err := parseIDAndParamsArgs(args[1:], "files get", "file-id", []string{"fields"})
		if err != nil {
			return err
		}
		out, err := gc.GetDriveFile(cliContext(), fileID, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "create":
		body, params, err := parseRequestFileAndParamsArgs(args[1:], "files create", []string{"fields"})
		if err != nil {
			return err
		}
		out, err := gc.CreateDriveFile(cliContext(), body, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "copy":
		fileID, body, params, err := parseIDRequestFileAndParamsArgs(args[1:], "files copy", "file-id", []string{"fields"})
		if err != nil {
			return err
		}
		out, err := gc.CopyDriveFile(cliContext(), fileID, body, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "update":
		fileID, body, params, err := parseIDRequestFileAndParamsArgs(args[1:], "files update", "file-id", []string{"fields"})
		if err != nil {
			return err
		}
		out, err := gc.UpdateDriveFile(cliContext(), fileID, body, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "delete":
		fileID, _, err := parseIDAndParamsArgs(args[1:], "files delete", "file-id", nil)
		if err != nil {
			return err
		}
		if err := gc.DeleteDriveFile(cliContext(), fileID); err != nil {
			return err
		}
		return writeJSON(cli.stdout, map[string]any{"deleted": true, "fileId": fileID})
	default:
		return fmt.Errorf("unsupported drive files command: %s", args[0])
	}
}

func (cli CLI) runDrivePermissions(gc GWSClient, args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printDrivePermissionsHelp()
		return nil
	}
	if isSingleHelpArg(args[1:]) {
		cli.printDrivePermissionsHelp()
		return nil
	}
	switch args[0] {
	case "list":
		fileID, _, err := parseIDAndParamsArgs(args[1:], "permissions list", "file-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.ListDrivePermissions(cliContext(), fileID)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "create":
		fileID, body, _, err := parseIDRequestFileAndParamsArgs(args[1:], "permissions create", "file-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.CreateDrivePermission(cliContext(), fileID, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "delete":
		parsed, err := argparse.Parse(args[1:], map[string]argparse.Spec{"file-id": {TakesValue: true}, "permission-id": {TakesValue: true}})
		if err != nil {
			return err
		}
		if len(parsed.Positionals) > 0 {
			return fmt.Errorf("permissions delete does not accept positional arguments")
		}
		fileID := parsed.First("file-id")
		permissionID := parsed.First("permission-id")
		if fileID == "" {
			return fmt.Errorf("permissions delete requires --file-id")
		}
		if permissionID == "" {
			return fmt.Errorf("permissions delete requires --permission-id")
		}
		if err := gc.DeleteDrivePermission(cliContext(), fileID, permissionID); err != nil {
			return err
		}
		return writeJSON(cli.stdout, map[string]any{"deleted": true, "fileId": fileID, "permissionId": permissionID})
	default:
		return fmt.Errorf("unsupported drive permissions command: %s", args[0])
	}
}

func (cli CLI) runSheets(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printSheetsHelp()
		return nil
	}
	switch args[0] {
	case "spreadsheets":
		if len(args[1:]) == 0 || isHelpToken(args[1]) {
			cli.printSheetsSpreadsheetsHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runSheetsSpreadsheets(gc, args[1:])
	case "values":
		if len(args[1:]) == 0 || isHelpToken(args[1]) {
			cli.printSheetsValuesHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runSheetsValues(gc, args[1:])
	default:
		return fmt.Errorf("unsupported sheets command: %s", args[0])
	}
}

func (cli CLI) runSheetsSpreadsheets(gc GWSClient, args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printSheetsSpreadsheetsHelp()
		return nil
	}
	if isSingleHelpArg(args[1:]) {
		cli.printSheetsSpreadsheetsHelp()
		return nil
	}
	switch args[0] {
	case "get":
		id, _, err := parseIDAndParamsArgs(args[1:], "spreadsheets get", "spreadsheet-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.GetSpreadsheet(cliContext(), id)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "create":
		body, _, err := parseRequestFileAndParamsArgs(args[1:], "spreadsheets create", nil)
		if err != nil {
			return err
		}
		out, err := gc.CreateSpreadsheet(cliContext(), body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "batch-update":
		id, body, _, err := parseIDRequestFileAndParamsArgs(args[1:], "spreadsheets batch-update", "spreadsheet-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.BatchUpdateSpreadsheet(cliContext(), id, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported sheets spreadsheets command: %s", args[0])
	}
}

func (cli CLI) runSheetsValues(gc GWSClient, args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printSheetsValuesHelp()
		return nil
	}
	if isSingleHelpArg(args[1:]) {
		cli.printSheetsValuesHelp()
		return nil
	}
	switch args[0] {
	case "get":
		spreadsheetID, valueRange, err := parseSpreadsheetRangeArgs(args[1:], "values get", false)
		if err != nil {
			return err
		}
		out, err := gc.GetSpreadsheetValues(cliContext(), spreadsheetID, valueRange)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "update":
		spreadsheetID, valueRange, valueInputOption, body, err := parseSpreadsheetValuesUpdateArgs(args[1:], "values update")
		if err != nil {
			return err
		}
		out, err := gc.UpdateSpreadsheetValues(cliContext(), spreadsheetID, valueRange, valueInputOption, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "batch-update":
		id, body, _, err := parseIDRequestFileAndParamsArgs(args[1:], "values batch-update", "spreadsheet-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.BatchUpdateSpreadsheetValues(cliContext(), id, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported sheets values command: %s", args[0])
	}
}

func (cli CLI) runDocs(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printDocsHelp()
		return nil
	}
	if args[0] != "documents" {
		return fmt.Errorf("unsupported docs command: %s", args[0])
	}
	if len(args) < 2 || isHelpToken(args[1]) {
		cli.printDocsDocumentsHelp()
		return nil
	}
	gc, err := cli.gwsClient()
	if err != nil {
		return err
	}
	switch args[1] {
	case "get":
		id, _, err := parseIDAndParamsArgs(args[2:], "documents get", "document-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.GetDocument(cliContext(), id)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "batch-update":
		id, body, _, err := parseIDRequestFileAndParamsArgs(args[2:], "documents batch-update", "document-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.BatchUpdateDocument(cliContext(), id, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported docs documents command: %s", args[1])
	}
}

func (cli CLI) runSlides(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printSlidesHelp()
		return nil
	}
	if args[0] != "presentations" {
		return fmt.Errorf("unsupported slides command: %s", args[0])
	}
	if len(args) < 2 || isHelpToken(args[1]) {
		cli.printSlidesPresentationsHelp()
		return nil
	}
	gc, err := cli.gwsClient()
	if err != nil {
		return err
	}
	switch args[1] {
	case "get":
		id, _, err := parseIDAndParamsArgs(args[2:], "presentations get", "presentation-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.GetPresentation(cliContext(), id)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "create":
		body, _, err := parseRequestFileAndParamsArgs(args[2:], "presentations create", nil)
		if err != nil {
			return err
		}
		out, err := gc.CreatePresentation(cliContext(), body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "batch-update":
		id, body, _, err := parseIDRequestFileAndParamsArgs(args[2:], "presentations batch-update", "presentation-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.BatchUpdatePresentation(cliContext(), id, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported slides presentations command: %s", args[1])
	}
}

func parseRawRequestArgs(args []string) (GWSRequest, error) {
	parsed, err := argparse.Parse(args, map[string]argparse.Spec{
		"api":          {TakesValue: true},
		"method":       {TakesValue: true},
		"path":         {TakesValue: true},
		"body":         {TakesValue: true},
		"request-file": {TakesValue: true},
	})
	if err != nil {
		return GWSRequest{}, err
	}
	if len(parsed.Positionals) > 0 {
		return GWSRequest{}, fmt.Errorf("request does not accept positional arguments")
	}
	request := GWSRequest{API: parsed.First("api"), Method: parsed.First("method"), Path: parsed.First("path")}
	if request.API == "" {
		return GWSRequest{}, fmt.Errorf("request requires --api")
	}
	if request.Path == "" {
		return GWSRequest{}, fmt.Errorf("request requires --path")
	}
	body, err := readOptionalBody(parsed.First("body"), parsed.First("request-file"), "request")
	if err != nil {
		return GWSRequest{}, err
	}
	request.Body = body
	return request, nil
}

func parseDriveFilesListArgs(args []string) (map[string]any, error) {
	parsed, err := argparse.Parse(args, map[string]argparse.Spec{
		"query":     {TakesValue: true},
		"page-size": {TakesValue: true},
		"fields":    {TakesValue: true},
	})
	if err != nil {
		return nil, err
	}
	if len(parsed.Positionals) > 0 {
		return nil, fmt.Errorf("files list does not accept positional arguments")
	}
	params := map[string]any{}
	copyFlag(params, parsed, "query", "q")
	copyFlag(params, parsed, "page-size", "pageSize")
	copyFlag(params, parsed, "fields", "fields")
	return params, nil
}

func parseIDAndParamsArgs(args []string, command string, idFlag string, paramFlags []string) (string, map[string]any, error) {
	specs := map[string]argparse.Spec{idFlag: {TakesValue: true}}
	for _, flag := range paramFlags {
		specs[flag] = argparse.Spec{TakesValue: true}
	}
	parsed, err := argparse.Parse(args, specs)
	if err != nil {
		return "", nil, err
	}
	if len(parsed.Positionals) > 0 {
		return "", nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	id := parsed.First(idFlag)
	if id == "" {
		return "", nil, fmt.Errorf("%s requires --%s", command, idFlag)
	}
	params := map[string]any{}
	for _, flag := range paramFlags {
		copyFlag(params, parsed, flag, flag)
	}
	return id, params, nil
}

func parseRequestFileAndParamsArgs(args []string, command string, paramFlags []string) (map[string]any, map[string]any, error) {
	specs := map[string]argparse.Spec{"request-file": {TakesValue: true}}
	for _, flag := range paramFlags {
		specs[flag] = argparse.Spec{TakesValue: true}
	}
	parsed, err := argparse.Parse(args, specs)
	if err != nil {
		return nil, nil, err
	}
	if len(parsed.Positionals) > 0 {
		return nil, nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	body, err := readRequestFileAsMap(parsed.First("request-file"), command)
	if err != nil {
		return nil, nil, err
	}
	params := map[string]any{}
	for _, flag := range paramFlags {
		copyFlag(params, parsed, flag, flag)
	}
	return body, params, nil
}

func parseIDRequestFileAndParamsArgs(args []string, command string, idFlag string, paramFlags []string) (string, map[string]any, map[string]any, error) {
	specs := map[string]argparse.Spec{idFlag: {TakesValue: true}, "request-file": {TakesValue: true}}
	for _, flag := range paramFlags {
		specs[flag] = argparse.Spec{TakesValue: true}
	}
	parsed, err := argparse.Parse(args, specs)
	if err != nil {
		return "", nil, nil, err
	}
	if len(parsed.Positionals) > 0 {
		return "", nil, nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	id := parsed.First(idFlag)
	if id == "" {
		return "", nil, nil, fmt.Errorf("%s requires --%s", command, idFlag)
	}
	body, err := readRequestFileAsMap(parsed.First("request-file"), command)
	if err != nil {
		return "", nil, nil, err
	}
	params := map[string]any{}
	for _, flag := range paramFlags {
		copyFlag(params, parsed, flag, flag)
	}
	return id, body, params, nil
}

func parseSpreadsheetRangeArgs(args []string, command string, includeValueInputOption bool) (string, string, error) {
	specs := map[string]argparse.Spec{"spreadsheet-id": {TakesValue: true}, "range": {TakesValue: true}}
	if includeValueInputOption {
		specs["value-input-option"] = argparse.Spec{TakesValue: true}
	}
	parsed, err := argparse.Parse(args, specs)
	if err != nil {
		return "", "", err
	}
	if len(parsed.Positionals) > 0 {
		return "", "", fmt.Errorf("%s does not accept positional arguments", command)
	}
	spreadsheetID := parsed.First("spreadsheet-id")
	if spreadsheetID == "" {
		return "", "", fmt.Errorf("%s requires --spreadsheet-id", command)
	}
	valueRange := parsed.First("range")
	if valueRange == "" {
		return "", "", fmt.Errorf("%s requires --range", command)
	}
	return spreadsheetID, valueRange, nil
}

func parseSpreadsheetValuesUpdateArgs(args []string, command string) (string, string, string, map[string]any, error) {
	parsed, err := argparse.Parse(args, map[string]argparse.Spec{
		"spreadsheet-id":     {TakesValue: true},
		"range":              {TakesValue: true},
		"value-input-option": {TakesValue: true},
		"request-file":       {TakesValue: true},
	})
	if err != nil {
		return "", "", "", nil, err
	}
	if len(parsed.Positionals) > 0 {
		return "", "", "", nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	spreadsheetID := parsed.First("spreadsheet-id")
	valueRange := parsed.First("range")
	valueInputOption := parsed.First("value-input-option")
	if spreadsheetID == "" {
		return "", "", "", nil, fmt.Errorf("%s requires --spreadsheet-id", command)
	}
	if valueRange == "" {
		return "", "", "", nil, fmt.Errorf("%s requires --range", command)
	}
	if valueInputOption == "" {
		return "", "", "", nil, fmt.Errorf("%s requires --value-input-option", command)
	}
	body, err := readRequestFileAsMap(parsed.First("request-file"), command)
	return spreadsheetID, valueRange, valueInputOption, body, err
}

func copyFlag(params map[string]any, parsed argparse.Parsed, flag string, name string) {
	if value := parsed.First(flag); value != "" {
		params[name] = value
	}
}

func readOptionalBody(body string, requestFile string, command string) (map[string]any, error) {
	if body != "" && requestFile != "" {
		return nil, fmt.Errorf("%s accepts only one of --body or --request-file", command)
	}
	if body != "" {
		return parseJSONMap([]byte(body), command)
	}
	if requestFile != "" {
		return readRequestFileAsMap(requestFile, command)
	}
	return nil, nil
}

func readRequestFileAsMap(requestFile string, command string) (map[string]any, error) {
	if requestFile == "" {
		return nil, fmt.Errorf("%s requires --request-file", command)
	}
	body, err := os.ReadFile(requestFile)
	if err != nil {
		return nil, err
	}
	return parseJSONMap(body, command)
}

func parseJSONMap(body []byte, command string) (map[string]any, error) {
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("%s body must contain a valid JSON object: %w", command, err)
	}
	if out == nil {
		return nil, fmt.Errorf("%s body must contain a JSON object", command)
	}
	return out, nil
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func containsAll(text string, values []string) bool {
	for _, value := range values {
		if !strings.Contains(text, value) {
			return false
		}
	}
	return true
}
