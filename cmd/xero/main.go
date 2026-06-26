package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mistlehq/tools/internal/argparse"
)

var Version = "dev"

type CLI struct {
	stdout io.Writer
	stderr io.Writer
	env    Environment
}

func (cli CLI) xeroClient() (XeroClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return XeroClient{}, err
	}
	return NewXeroClient(config), nil
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
	case "tenants":
		return cli.runTenants(args[2:])
	case "api":
		return cli.runAPI(args[2:])
	case "mcp":
		return cli.runMCP(args[2:])
	default:
		return fmt.Errorf("unsupported command: %s", args[1])
	}
}

func (cli CLI) runTenants(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printTenantsHelp()
		return nil
	}
	if args[0] != "list" {
		return fmt.Errorf("unsupported tenants command: %s", args[0])
	}
	if len(args) == 2 && isHelp(args[1]) {
		cli.printTenantsHelp()
		return nil
	}
	xc, err := cli.xeroClient()
	if err != nil {
		return err
	}
	out, err := xc.ListTenants()
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runAPI(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printAPIHelp()
		return nil
	}

	switch args[0] {
	case "get":
		return cli.runAPIEndpoint(args[1:], false, XeroClient.GetAPIEndpoint)
	case "post":
		return cli.runAPIEndpoint(args[1:], true, XeroClient.PostAPIEndpoint)
	case "put":
		return cli.runAPIEndpoint(args[1:], true, XeroClient.PutAPIEndpoint)
	case "delete":
		return cli.runAPIEndpoint(args[1:], false, XeroClient.DeleteAPIEndpoint)
	default:
		return fmt.Errorf("unsupported api command: %s", args[0])
	}
}

func (cli CLI) runAPIEndpoint(args []string, bodyRequired bool, call func(XeroClient, XeroAPIRequest) (XeroJSONResult, error)) error {
	if len(args) == 1 && isHelp(args[0]) {
		cli.printAPIHelp()
		return nil
	}
	xc, err := cli.xeroClient()
	if err != nil {
		return err
	}
	request, err := parseAPIArgs(args, bodyRequired)
	if err != nil {
		return err
	}
	out, err := call(xc, request)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func parseAPIArgs(args []string, bodyRequired bool) (XeroAPIRequest, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"family":    {TakesValue: true},
		"tenant-id": {TakesValue: true},
		"endpoint":  {TakesValue: true},
		"query":     {TakesValue: true},
		"body":      {TakesValue: true},
	})
	if err != nil {
		return XeroAPIRequest{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return XeroAPIRequest{}, fmt.Errorf("api command does not accept positional arguments")
	}

	query, err := parseQueryFlags(parsedArgs.Flags["query"])
	if err != nil {
		return XeroAPIRequest{}, err
	}

	body := json.RawMessage(parsedArgs.First("body"))
	if bodyRequired && len(body) == 0 {
		return XeroAPIRequest{}, fmt.Errorf("--body is required")
	}
	if len(body) > 0 && !json.Valid(body) {
		return XeroAPIRequest{}, fmt.Errorf("--body must be valid JSON")
	}

	return XeroAPIRequest{
		Family:   parsedArgs.First("family"),
		TenantID: parsedArgs.First("tenant-id"),
		Endpoint: parsedArgs.First("endpoint"),
		Query:    query,
		Body:     body,
	}, nil
}

func parseQueryFlags(values []string) (map[string]string, error) {
	query := map[string]string{}
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" {
			return nil, fmt.Errorf("--query must use name=value format")
		}
		query[parts[0]] = parts[1]
	}
	return query, nil
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func isHelp(value string) bool {
	return value == "help" || value == "-h" || value == "--help"
}

func (cli CLI) printHelp() {
	fmt.Fprintln(cli.stdout, `xero

CLI for Xero API access through Mistle-managed credentials.

Commands:
  xero tenants list
  xero api get --family <family> --tenant-id <tenant-id> --endpoint <path>
  xero api post --family <family> --tenant-id <tenant-id> --endpoint <path> --body <json>
  xero api put --family <family> --tenant-id <tenant-id> --endpoint <path> --body <json>
  xero api delete --family <family> --tenant-id <tenant-id> --endpoint <path>
  xero mcp serve

API families:
  accounting, assets, files, projects

Configuration:
  XERO_API_BASE_URL must point at the Xero API origin, normally https://api.xero.com.`)
}

func (cli CLI) printTenantsHelp() {
	fmt.Fprintln(cli.stdout, `xero tenants

Commands:
  xero tenants list`)
}

func (cli CLI) printAPIHelp() {
	fmt.Fprintln(cli.stdout, `xero api

Commands:
  xero api get --family <family> --tenant-id <tenant-id> --endpoint <path>
  xero api post --family <family> --tenant-id <tenant-id> --endpoint <path> --body <json>
  xero api put --family <family> --tenant-id <tenant-id> --endpoint <path> --body <json>
  xero api delete --family <family> --tenant-id <tenant-id> --endpoint <path>

Use --query name=value for documented query parameters. Repeat --query for multiple parameters.`)
}

func main() {
	cli := CLI{
		stdout: os.Stdout,
		stderr: os.Stderr,
		env:    loadEnvironment(),
	}
	if err := cli.run(os.Args); err != nil {
		fmt.Fprintln(cli.stderr, err)
		os.Exit(1)
	}
}
