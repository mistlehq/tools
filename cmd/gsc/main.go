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
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	env    Environment
}

func (cli CLI) gscClient() (GSCClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return GSCClient{}, err
	}
	return NewGSCClient(config), nil
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
	case "sites":
		return cli.runSites(args[2:])
	case "searchanalytics":
		return cli.runSearchAnalytics(args[2:])
	case "sitemaps":
		return cli.runSitemaps(args[2:])
	case "url-inspection":
		return cli.runURLInspection(args[2:])
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
	gc, err := cli.gscClient()
	if err != nil {
		return err
	}
	return cli.runAuthTest(gc, args[1:])
}

func (cli CLI) runSites(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printSitesHelp()
		return nil
	}
	gc, err := cli.gscClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		return cli.runSitesList(gc, args[1:])
	case "get":
		return cli.runSitesGet(gc, args[1:])
	default:
		return fmt.Errorf("unsupported sites command: %s", args[0])
	}
}

func (cli CLI) runSearchAnalytics(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printSearchAnalyticsHelp()
		return nil
	}
	if args[0] != "query" {
		return fmt.Errorf("unsupported searchanalytics command: %s", args[0])
	}
	gc, err := cli.gscClient()
	if err != nil {
		return err
	}
	return cli.runSearchAnalyticsQuery(gc, args[1:])
}

func (cli CLI) runSitemaps(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printSitemapsHelp()
		return nil
	}
	gc, err := cli.gscClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		return cli.runSitemapsList(gc, args[1:])
	case "get":
		return cli.runSitemapsGet(gc, args[1:])
	default:
		return fmt.Errorf("unsupported sitemaps command: %s", args[0])
	}
}

func (cli CLI) runURLInspection(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printURLInspectionHelp()
		return nil
	}
	if args[0] != "inspect" {
		return fmt.Errorf("unsupported url-inspection command: %s", args[0])
	}
	gc, err := cli.gscClient()
	if err != nil {
		return err
	}
	return cli.runURLInspectionInspect(gc, args[1:])
}

func (cli CLI) runAuthTest(gc GSCClient, args []string) error {
	parsedArgs, err := parseSiteCommandArgs(args, "auth test")
	if err != nil {
		return err
	}
	out, err := gc.AuthTest(parsedArgs.First("site-url"))
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	writeSite(cli.stdout, out)
	return nil
}

func (cli CLI) runSitesList(gc GSCClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{"json": {}})
	if err != nil {
		return err
	}
	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("sites list does not accept positional arguments")
	}
	out, err := gc.ListSites()
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	fmt.Fprintln(cli.stdout, "SITE_URL\tPERMISSION_LEVEL")
	for _, site := range out.SiteEntry {
		fmt.Fprintf(cli.stdout, "%s\t%s\n", site.SiteURL, site.PermissionLevel)
	}
	return nil
}

func (cli CLI) runSitesGet(gc GSCClient, args []string) error {
	parsedArgs, err := parseSiteCommandArgs(args, "sites get")
	if err != nil {
		return err
	}
	out, err := gc.GetSite(parsedArgs.First("site-url"))
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	writeSite(cli.stdout, out)
	return nil
}

func (cli CLI) runSearchAnalyticsQuery(gc GSCClient, args []string) error {
	siteURL, body, err := parseSiteRequestFileArgs(args, "searchanalytics query")
	if err != nil {
		return err
	}
	out, err := gc.QuerySearchAnalytics(siteURL, body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runSitemapsList(gc GSCClient, args []string) error {
	parsedArgs, err := parseSiteCommandArgs(args, "sitemaps list")
	if err != nil {
		return err
	}
	out, err := gc.ListSitemaps(parsedArgs.First("site-url"))
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	fmt.Fprintln(cli.stdout, "PATH\tTYPE\tERRORS\tWARNINGS")
	for _, sitemap := range out.Sitemap {
		fmt.Fprintf(cli.stdout, "%s\t%s\t%v\t%v\n", sitemap.Path, sitemap.Type, sitemap.Errors, sitemap.Warnings)
	}
	return nil
}

func (cli CLI) runSitemapsGet(gc GSCClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"site-url":  {TakesValue: true},
		"feed-path": {TakesValue: true},
		"json":      {},
	})
	if err != nil {
		return err
	}
	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("sitemaps get does not accept positional arguments")
	}
	siteURL := parsedArgs.First("site-url")
	if siteURL == "" {
		return fmt.Errorf("sitemaps get requires --site-url")
	}
	feedPath := parsedArgs.First("feed-path")
	if feedPath == "" {
		return fmt.Errorf("sitemaps get requires --feed-path")
	}
	out, err := gc.GetSitemap(siteURL, feedPath)
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	writeSitemap(cli.stdout, out)
	return nil
}

func (cli CLI) runURLInspectionInspect(gc GSCClient, args []string) error {
	body, err := parseRequestFileArgs(args, "url-inspection inspect")
	if err != nil {
		return err
	}
	out, err := gc.InspectURL(body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func parseSiteCommandArgs(args []string, command string) (argparse.Parsed, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"site-url": {TakesValue: true},
		"json":     {},
	})
	if err != nil {
		return argparse.Parsed{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return argparse.Parsed{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	if parsedArgs.First("site-url") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --site-url", command)
	}
	return parsedArgs, nil
}

func parseSiteRequestFileArgs(args []string, command string) (string, json.RawMessage, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"site-url":     {TakesValue: true},
		"request-file": {TakesValue: true},
	})
	if err != nil {
		return "", nil, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return "", nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	siteURL := parsedArgs.First("site-url")
	if siteURL == "" {
		return "", nil, fmt.Errorf("%s requires --site-url", command)
	}
	body, err := readRequestFile(parsedArgs.First("request-file"), command)
	if err != nil {
		return "", nil, err
	}
	return siteURL, body, nil
}

func parseRequestFileArgs(args []string, command string) (json.RawMessage, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"request-file": {TakesValue: true},
	})
	if err != nil {
		return nil, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	return readRequestFile(parsedArgs.First("request-file"), command)
}

func readRequestFile(requestFile string, command string) (json.RawMessage, error) {
	if requestFile == "" {
		return nil, fmt.Errorf("%s requires --request-file", command)
	}
	body, err := os.ReadFile(requestFile)
	if err != nil {
		return nil, err
	}
	var parsed any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("%s request file must contain valid JSON: %w", command, err)
	}
	return json.RawMessage(body), nil
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func writeSite(writer io.Writer, site GSCSite) {
	fmt.Fprintln(writer, "Site URL: "+site.SiteURL)
	fmt.Fprintln(writer, "Permission Level: "+site.PermissionLevel)
}

func writeSitemap(writer io.Writer, sitemap GSCSitemap) {
	fmt.Fprintln(writer, "Path: "+sitemap.Path)
	if sitemap.Type != "" {
		fmt.Fprintln(writer, "Type: "+sitemap.Type)
	}
	if sitemap.LastSubmitted != "" {
		fmt.Fprintln(writer, "Last Submitted: "+sitemap.LastSubmitted)
	}
	if sitemap.LastDownloaded != "" {
		fmt.Fprintln(writer, "Last Downloaded: "+sitemap.LastDownloaded)
	}
	if sitemap.Errors != nil {
		fmt.Fprintf(writer, "Errors: %v\n", sitemap.Errors)
	}
	if sitemap.Warnings != nil {
		fmt.Fprintf(writer, "Warnings: %v\n", sitemap.Warnings)
	}
}

func (cli CLI) printHelp() {
	fmt.Fprint(cli.stdout, `gsc

CLI for Google Search Console.

Usage:
  gsc help
  gsc version
  gsc auth help
  gsc sites help
  gsc searchanalytics help
  gsc sitemaps help
  gsc url-inspection help
  gsc mcp help

Commands:
  help
  version
  auth
  sites
  searchanalytics
  sitemaps
  url-inspection
  mcp
`)
}

func (cli CLI) printAuthHelp() {
	fmt.Fprint(cli.stdout, `gsc auth

Inspect Google Search Console API access.

Usage:
  gsc auth help
  gsc auth test --site-url <site-url> [--json]
`)
}

func (cli CLI) printSitesHelp() {
	fmt.Fprint(cli.stdout, `gsc sites

Inspect Google Search Console sites.

Usage:
  gsc sites help
  gsc sites list [--json]
  gsc sites get --site-url <site-url> [--json]
`)
}

func (cli CLI) printSearchAnalyticsHelp() {
	fmt.Fprint(cli.stdout, `gsc searchanalytics

Query Search Console Search Analytics data.

Usage:
  gsc searchanalytics help
  gsc searchanalytics query --site-url <site-url> --request-file <json>
`)
}

func (cli CLI) printSitemapsHelp() {
	fmt.Fprint(cli.stdout, `gsc sitemaps

Inspect Search Console sitemaps.

Usage:
  gsc sitemaps help
  gsc sitemaps list --site-url <site-url> [--json]
  gsc sitemaps get --site-url <site-url> --feed-path <sitemap-url> [--json]
`)
}

func (cli CLI) printURLInspectionHelp() {
	fmt.Fprint(cli.stdout, `gsc url-inspection

Inspect URL indexing information.

Usage:
  gsc url-inspection help
  gsc url-inspection inspect --request-file <json>
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

func containsSite(sites []GSCSite, siteURL string) bool {
	for _, candidate := range sites {
		if strings.TrimSpace(candidate.SiteURL) == siteURL {
			return true
		}
	}
	return false
}
