package main

import (
	"encoding/json"
	"fmt"
	"github.com/mistlehq/tools/internal/argparse"
	"io"
	"os"
	"strings"
)

var Version = "dev"

type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	env    Environment
}

func (cli CLI) gaClient() (GAClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return GAClient{}, err
	}
	return NewGAClient(config), nil
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
	case "account-summaries":
		return cli.runAccountSummaries(args[2:])
	case "properties":
		return cli.runProperties(args[2:])
	case "metadata":
		return cli.runMetadata(args[2:])
	case "compatibility":
		return cli.runCompatibility(args[2:])
	case "reports":
		return cli.runReports(args[2:])
	case "google-ads-links":
		return cli.runGoogleAdsLinks(args[2:])
	case "mcp":
		return cli.runMCP(args[2:])
	default:
		return fmt.Errorf("unsupported command: %s", args[1])
	}
}

func (cli CLI) runAuth(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printAuthHelp()
		return nil
	}
	if args[0] != "test" {
		return fmt.Errorf("unsupported auth command: %s", args[0])
	}
	gc, err := cli.gaClient()
	if err != nil {
		return err
	}
	return cli.runAuthTest(gc, args[1:])
}

func (cli CLI) runAccountSummaries(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printAccountSummariesHelp()
		return nil
	}
	if args[0] != "list" {
		return fmt.Errorf("unsupported account-summaries command: %s", args[0])
	}
	gc, err := cli.gaClient()
	if err != nil {
		return err
	}
	return cli.runAccountSummariesList(gc, args[1:])
}

func (cli CLI) runProperties(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printPropertiesHelp()
		return nil
	}
	if args[0] != "get" {
		return fmt.Errorf("unsupported properties command: %s", args[0])
	}
	gc, err := cli.gaClient()
	if err != nil {
		return err
	}
	return cli.runPropertiesGet(gc, args[1:])
}

func (cli CLI) runMetadata(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printMetadataHelp()
		return nil
	}
	if args[0] != "get" {
		return fmt.Errorf("unsupported metadata command: %s", args[0])
	}
	gc, err := cli.gaClient()
	if err != nil {
		return err
	}
	return cli.runMetadataGet(gc, args[1:])
}

func (cli CLI) runCompatibility(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printCompatibilityHelp()
		return nil
	}
	if args[0] != "check" {
		return fmt.Errorf("unsupported compatibility command: %s", args[0])
	}
	gc, err := cli.gaClient()
	if err != nil {
		return err
	}
	return cli.runCompatibilityCheck(gc, args[1:])
}

func (cli CLI) runReports(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printReportsHelp()
		return nil
	}
	gc, err := cli.gaClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "run":
		return cli.runReportRun(gc, args[1:])
	case "realtime":
		return cli.runReportRealtime(gc, args[1:])
	case "funnel":
		return cli.runReportFunnel(gc, args[1:])
	default:
		return fmt.Errorf("unsupported reports command: %s", args[0])
	}
}

func (cli CLI) runGoogleAdsLinks(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printGoogleAdsLinksHelp()
		return nil
	}
	if args[0] != "list" {
		return fmt.Errorf("unsupported google-ads-links command: %s", args[0])
	}
	gc, err := cli.gaClient()
	if err != nil {
		return err
	}
	return cli.runGoogleAdsLinksList(gc, args[1:])
}

func (cli CLI) runAuthTest(gc GAClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"property": {TakesValue: true},
		"json":     {},
	})
	if err != nil {
		return err
	}
	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("auth test does not accept positional arguments")
	}
	property := parsedArgs.First("property")
	if property == "" {
		return fmt.Errorf("auth test requires --property")
	}
	out, err := gc.AuthTest(property)
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	fmt.Fprintln(cli.stdout, "Property: "+out.Name)
	fmt.Fprintln(cli.stdout, "Display Name: "+out.DisplayName)
	return nil
}

func (cli CLI) runAccountSummariesList(gc GAClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{"json": {}})
	if err != nil {
		return err
	}
	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("account-summaries list does not accept positional arguments")
	}
	out, err := gc.ListAccountSummaries()
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	fmt.Fprintln(cli.stdout, "ACCOUNT\tDISPLAY_NAME\tPROPERTY_COUNT")
	for _, summary := range out.AccountSummaries {
		fmt.Fprintf(cli.stdout, "%s\t%s\t%d\n", summary.Account, summary.DisplayName, len(summary.PropertySummaries))
	}
	return nil
}

func (cli CLI) runPropertiesGet(gc GAClient, args []string) error {
	parsedArgs, err := parsePropertyCommandArgs(args, "properties get")
	if err != nil {
		return err
	}
	out, err := gc.GetProperty(parsedArgs.First("property"))
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	writeProperty(cli.stdout, out)
	return nil
}

func (cli CLI) runMetadataGet(gc GAClient, args []string) error {
	parsedArgs, err := parsePropertyCommandArgs(args, "metadata get")
	if err != nil {
		return err
	}
	out, err := gc.GetMetadata(parsedArgs.First("property"))
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	fmt.Fprintln(cli.stdout, "Name: "+out.Name)
	fmt.Fprintf(cli.stdout, "Dimensions: %d\n", len(out.Dimensions))
	fmt.Fprintf(cli.stdout, "Metrics: %d\n", len(out.Metrics))
	return nil
}

func (cli CLI) runCompatibilityCheck(gc GAClient, args []string) error {
	property, body, err := parsePropertyRequestFileArgs(args, "compatibility check")
	if err != nil {
		return err
	}
	out, err := gc.CheckCompatibility(property, body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runReportRun(gc GAClient, args []string) error {
	property, body, err := parsePropertyRequestFileArgs(args, "reports run")
	if err != nil {
		return err
	}
	out, err := gc.RunReport(property, body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runReportRealtime(gc GAClient, args []string) error {
	property, body, err := parsePropertyRequestFileArgs(args, "reports realtime")
	if err != nil {
		return err
	}
	out, err := gc.RunRealtimeReport(property, body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runReportFunnel(gc GAClient, args []string) error {
	property, body, err := parsePropertyRequestFileArgs(args, "reports funnel")
	if err != nil {
		return err
	}
	out, err := gc.RunFunnelReport(property, body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runGoogleAdsLinksList(gc GAClient, args []string) error {
	parsedArgs, err := parsePropertyCommandArgs(args, "google-ads-links list")
	if err != nil {
		return err
	}
	out, err := gc.ListGoogleAdsLinks(parsedArgs.First("property"))
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	fmt.Fprintln(cli.stdout, "NAME\tCUSTOMER_ID")
	for _, link := range out.GoogleAdsLinks {
		fmt.Fprintf(cli.stdout, "%s\t%s\n", link.Name, link.CustomerID)
	}
	return nil
}

func parsePropertyCommandArgs(args []string, command string) (argparse.Parsed, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"property": {TakesValue: true},
		"json":     {},
	})
	if err != nil {
		return argparse.Parsed{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return argparse.Parsed{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	if parsedArgs.First("property") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --property", command)
	}
	return parsedArgs, nil
}

func parsePropertyRequestFileArgs(args []string, command string) (string, json.RawMessage, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"property":     {TakesValue: true},
		"request-file": {TakesValue: true},
	})
	if err != nil {
		return "", nil, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return "", nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	property := parsedArgs.First("property")
	if property == "" {
		return "", nil, fmt.Errorf("%s requires --property", command)
	}
	requestFile := parsedArgs.First("request-file")
	if requestFile == "" {
		return "", nil, fmt.Errorf("%s requires --request-file", command)
	}
	body, err := os.ReadFile(requestFile)
	if err != nil {
		return "", nil, err
	}
	var parsed any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", nil, fmt.Errorf("%s request file must contain valid JSON: %w", command, err)
	}
	return property, json.RawMessage(body), nil
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func writeProperty(writer io.Writer, property GAProperty) {
	fmt.Fprintln(writer, "Name: "+property.Name)
	fmt.Fprintln(writer, "Display Name: "+property.DisplayName)
	if property.Parent != "" {
		fmt.Fprintln(writer, "Parent: "+property.Parent)
	}
	if property.TimeZone != "" {
		fmt.Fprintln(writer, "Time Zone: "+property.TimeZone)
	}
	if property.CurrencyCode != "" {
		fmt.Fprintln(writer, "Currency Code: "+property.CurrencyCode)
	}
}

func (cli CLI) printHelp() {
	fmt.Fprint(cli.stdout, `ga

CLI for Google Analytics.

Usage:
  ga help
  ga version
  ga auth help
  ga account-summaries help
  ga properties help
  ga metadata help
  ga compatibility help
  ga reports help
  ga google-ads-links help
  ga mcp help

Commands:
  help
  version
  auth
  account-summaries
  properties
  metadata
  compatibility
  reports
  google-ads-links
  mcp
`)
}

func (cli CLI) printAuthHelp() {
	fmt.Fprint(cli.stdout, `ga auth

Inspect Google Analytics API access.

Usage:
  ga auth help
  ga auth test --property properties/<id> [--json]
`)
}

func (cli CLI) printAccountSummariesHelp() {
	fmt.Fprint(cli.stdout, `ga account-summaries

Inspect Google Analytics account summaries.

Usage:
  ga account-summaries help
  ga account-summaries list [--json]
`)
}

func (cli CLI) printPropertiesHelp() {
	fmt.Fprint(cli.stdout, `ga properties

Inspect Google Analytics properties.

Usage:
  ga properties help
  ga properties get --property properties/<id> [--json]
`)
}

func (cli CLI) printMetadataHelp() {
	fmt.Fprint(cli.stdout, `ga metadata

Inspect Google Analytics property metadata.

Usage:
  ga metadata help
  ga metadata get --property properties/<id> [--json]
`)
}

func (cli CLI) printCompatibilityHelp() {
	fmt.Fprint(cli.stdout, `ga compatibility

Check Google Analytics report compatibility.

Usage:
  ga compatibility help
  ga compatibility check --property properties/<id> --request-file <json>
`)
}

func (cli CLI) printReportsHelp() {
	fmt.Fprint(cli.stdout, `ga reports

Run Google Analytics reports.

Usage:
  ga reports help
  ga reports run --property properties/<id> --request-file <json>
  ga reports realtime --property properties/<id> --request-file <json>
  ga reports funnel --property properties/<id> --request-file <json>
`)
}

func (cli CLI) printGoogleAdsLinksHelp() {
	fmt.Fprint(cli.stdout, `ga google-ads-links

Inspect Google Ads links for a Google Analytics property.

Usage:
  ga google-ads-links help
  ga google-ads-links list --property properties/<id> [--json]
`)
}

func main() {
	cli := CLI{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		env:    loadEnvironment(),
	}

	if err := cli.run(os.Args); err != nil {
		fmt.Fprintln(cli.stderr, err)
		os.Exit(1)
	}
}

func containsResource(resources []string, resource string) bool {
	for _, candidate := range resources {
		if strings.TrimSpace(candidate) == resource {
			return true
		}
	}
	return false
}
