package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mistlehq/tools/internal/argparse"
	"github.com/mistlehq/tools/internal/textinput"
)

// Version is the current googleads CLI version.
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

func isHelpToken(arg string) bool {
	return arg == "help" || arg == "-h" || arg == "--help"
}

func isSingleHelpArg(args []string) bool {
	return len(args) == 1 && isHelpToken(args[0])
}

func cliContext() context.Context {
	return context.Background()
}

func (cli CLI) googleAdsClient() (GoogleAdsClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return GoogleAdsClient{}, err
	}
	return NewGoogleAdsClient(config), nil
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
	case "customers":
		return cli.runCustomers(args[2:])
	case "gaql":
		return cli.runGAQL(args[2:])
	case "fields":
		return cli.runFields(args[2:])
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
	gc, err := cli.googleAdsClient()
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
	request, err := parseRequestArgs(cli.stdin, args)
	if err != nil {
		return err
	}
	gc, err := cli.googleAdsClient()
	if err != nil {
		return err
	}
	response, err := gc.Request(request)
	if err != nil {
		return err
	}
	_, err = cli.stdout.Write(append(response, '\n'))
	return err
}

func (cli CLI) runCustomers(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printCustomersHelp()
		return nil
	}
	if args[0] != "list-accessible" {
		return fmt.Errorf("unsupported customers command: %s", args[0])
	}
	if isSingleHelpArg(args[1:]) {
		cli.printCustomersListAccessibleHelp()
		return nil
	}
	if len(args[1:]) > 0 {
		return fmt.Errorf("customers list-accessible does not accept positional arguments")
	}
	gc, err := cli.googleAdsClient()
	if err != nil {
		return err
	}
	out, err := gc.ListAccessibleCustomers(cliContext())
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runGAQL(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printGAQLHelp()
		return nil
	}
	input, err := parseGAQLArgs(cli.stdin, args[1:], "gaql "+args[0])
	if err != nil {
		return err
	}
	gc, err := cli.googleAdsClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "search":
		if isSingleHelpArg(args[1:]) {
			cli.printGAQLSearchHelp()
			return nil
		}
		out, err := gc.Search(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "search-stream":
		if isSingleHelpArg(args[1:]) {
			cli.printGAQLSearchStreamHelp()
			return nil
		}
		out, err := gc.SearchStream(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported gaql command: %s", args[0])
	}
}

func (cli CLI) runFields(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printFieldsHelp()
		return nil
	}
	gc, err := cli.googleAdsClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "search":
		if isSingleHelpArg(args[1:]) {
			cli.printFieldsSearchHelp()
			return nil
		}
		query, err := parseQueryOnlyArgs(cli.stdin, args[1:], "fields search")
		if err != nil {
			return err
		}
		out, err := gc.SearchFields(cliContext(), query)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		if isSingleHelpArg(args[1:]) {
			cli.printFieldGetHelp()
			return nil
		}
		resourceName, err := parseResourceNameArgs(args[1:], "fields get")
		if err != nil {
			return err
		}
		out, err := gc.GetField(cliContext(), resourceName)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported fields command: %s", args[0])
	}
}

func parseRequestArgs(stdin io.Reader, args []string) (GoogleAdsRequest, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"method":            {TakesValue: true},
		"path":              {TakesValue: true},
		"login-customer-id": {TakesValue: true},
		"params":            {TakesValue: true},
		"params-file":       {TakesValue: true},
		"body":              {TakesValue: true},
		"body-file":         {TakesValue: true},
	})
	if err != nil {
		return GoogleAdsRequest{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return GoogleAdsRequest{}, fmt.Errorf("request does not accept positional arguments")
	}
	params, err := readJSONObject(stdin, parsedArgs.First("params"), parsedArgs.First("params-file"), "params")
	if err != nil {
		return GoogleAdsRequest{}, err
	}
	body, err := readJSONObject(stdin, parsedArgs.First("body"), parsedArgs.First("body-file"), "body")
	if err != nil {
		return GoogleAdsRequest{}, err
	}
	return GoogleAdsRequest{Method: parsedArgs.First("method"), Path: parsedArgs.First("path"), LoginCustomerID: parsedArgs.First("login-customer-id"), Params: params, Body: body}, nil
}

func parseGAQLArgs(stdin io.Reader, args []string, command string) (GoogleAdsGAQLInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"customer-id":       {TakesValue: true},
		"login-customer-id": {TakesValue: true},
		"query":             {TakesValue: true},
		"query-file":        {TakesValue: true},
		"page-size":         {TakesValue: true},
		"page-token":        {TakesValue: true},
		"summary-row":       {TakesValue: true},
		"params":            {TakesValue: true},
		"params-file":       {TakesValue: true},
	})
	if err != nil {
		return GoogleAdsGAQLInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return GoogleAdsGAQLInput{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	query, err := textinput.Read(stdin, "query", parsedArgs.First("query"), "query-file", parsedArgs.First("query-file"))
	if err != nil {
		return GoogleAdsGAQLInput{}, err
	}
	params, err := readJSONObject(stdin, parsedArgs.First("params"), parsedArgs.First("params-file"), "params")
	if err != nil {
		return GoogleAdsGAQLInput{}, err
	}
	return GoogleAdsGAQLInput{CustomerID: parsedArgs.First("customer-id"), LoginCustomerID: parsedArgs.First("login-customer-id"), Query: string(query), PageSize: parsedArgs.First("page-size"), PageToken: parsedArgs.First("page-token"), SummaryRow: parsedArgs.First("summary-row"), Params: params}, nil
}

func parseQueryOnlyArgs(stdin io.Reader, args []string, command string) (string, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"query":      {TakesValue: true},
		"query-file": {TakesValue: true},
	})
	if err != nil {
		return "", err
	}
	if len(parsedArgs.Positionals) > 0 {
		return "", fmt.Errorf("%s does not accept positional arguments", command)
	}
	query, err := textinput.Read(stdin, "query", parsedArgs.First("query"), "query-file", parsedArgs.First("query-file"))
	if err != nil {
		return "", err
	}
	return string(query), nil
}

func parseResourceNameArgs(args []string, command string) (string, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{"resource-name": {TakesValue: true}})
	if err != nil {
		return "", err
	}
	if len(parsedArgs.Positionals) > 0 {
		return "", fmt.Errorf("%s does not accept positional arguments", command)
	}
	return parsedArgs.First("resource-name"), nil
}

func readJSONObject(stdin io.Reader, inline string, filePath string, label string) (map[string]any, error) {
	if inline == "" && filePath == "" {
		return nil, nil
	}
	raw, err := textinput.Read(stdin, label, inline, label+"-file", filePath)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("%s must be a JSON object: %w", label, err)
	}
	return out, nil
}

func writeJSON(output io.Writer, value any) error {
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
