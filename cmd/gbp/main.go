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

func (cli CLI) gbpClient() (GBPClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return GBPClient{}, err
	}
	return NewGBPClient(config), nil
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
	case "accounts":
		return cli.runAccounts(args[2:])
	case "locations":
		return cli.runLocations(args[2:])
	case "reviews":
		return cli.runReviews(args[2:])
	case "media":
		return cli.runMedia(args[2:])
	case "local-posts":
		return cli.runLocalPosts(args[2:])
	case "performance":
		return cli.runPerformance(args[2:])
	case "mcp":
		return cli.runMCP(args[2:])
	default:
		return fmt.Errorf("unsupported command: %s", args[1])
	}
}

func (cli CLI) runAuth(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printAuthHelp()
		return nil
	}
	if args[0] != "test" {
		return fmt.Errorf("unsupported auth command: %s", args[0])
	}
	gc, err := cli.gbpClient()
	if err != nil {
		return err
	}
	return cli.runAuthTest(gc, args[1:])
}

func (cli CLI) runAccounts(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printAccountsHelp()
		return nil
	}
	gc, err := cli.gbpClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		return cli.runAccountsList(gc, args[1:])
	case "get":
		return cli.runAccountsGet(gc, args[1:])
	default:
		return fmt.Errorf("unsupported accounts command: %s", args[0])
	}
}

func (cli CLI) runLocations(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printLocationsHelp()
		return nil
	}
	gc, err := cli.gbpClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		return cli.runLocationsList(gc, args[1:])
	case "get":
		return cli.runLocationsGet(gc, args[1:])
	case "create":
		return cli.runLocationsCreate(gc, args[1:])
	case "patch":
		return cli.runLocationsPatch(gc, args[1:])
	case "delete":
		return cli.runLocationsDelete(gc, args[1:])
	default:
		return fmt.Errorf("unsupported locations command: %s", args[0])
	}
}

func (cli CLI) runReviews(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printReviewsHelp()
		return nil
	}
	gc, err := cli.gbpClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		return cli.runReviewsList(gc, args[1:])
	case "get":
		return cli.runReviewsGet(gc, args[1:])
	case "update-reply":
		return cli.runReviewsUpdateReply(gc, args[1:])
	case "delete-reply":
		return cli.runReviewsDeleteReply(gc, args[1:])
	default:
		return fmt.Errorf("unsupported reviews command: %s", args[0])
	}
}

func (cli CLI) runMedia(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printMediaHelp()
		return nil
	}
	gc, err := cli.gbpClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		return cli.runMediaList(gc, args[1:])
	case "create":
		return cli.runMediaCreate(gc, args[1:])
	case "get":
		return cli.runMediaGet(gc, args[1:])
	case "patch":
		return cli.runMediaPatch(gc, args[1:])
	case "delete":
		return cli.runMediaDelete(gc, args[1:])
	case "start-upload":
		return cli.runMediaStartUpload(gc, args[1:])
	default:
		return fmt.Errorf("unsupported media command: %s", args[0])
	}
}

func (cli CLI) runLocalPosts(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printLocalPostsHelp()
		return nil
	}
	gc, err := cli.gbpClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		return cli.runLocalPostsList(gc, args[1:])
	case "create":
		return cli.runLocalPostsCreate(gc, args[1:])
	case "get":
		return cli.runLocalPostsGet(gc, args[1:])
	case "patch":
		return cli.runLocalPostsPatch(gc, args[1:])
	case "delete":
		return cli.runLocalPostsDelete(gc, args[1:])
	case "report-insights":
		return cli.runLocalPostsReportInsights(gc, args[1:])
	default:
		return fmt.Errorf("unsupported local-posts command: %s", args[0])
	}
}

func (cli CLI) runPerformance(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printPerformanceHelp()
		return nil
	}
	gc, err := cli.gbpClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "daily-metrics":
		return cli.runPerformanceDailyMetrics(gc, args[1:])
	case "search-keywords":
		return cli.runPerformanceSearchKeywords(gc, args[1:])
	default:
		return fmt.Errorf("unsupported performance command: %s", args[0])
	}
}

func (cli CLI) runAuthTest(gc GBPClient, args []string) error {
	parsedArgs, err := parseNoPositionalArgs(args, "auth test", map[string]argparse.Spec{"json": {}})
	if err != nil {
		return err
	}
	out, err := gc.AuthTest()
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	writeAccountsTable(cli.stdout, out.Accounts)
	return nil
}

func (cli CLI) runAccountsList(gc GBPClient, args []string) error {
	parsedArgs, err := parseNoPositionalArgs(args, "accounts list", map[string]argparse.Spec{"json": {}})
	if err != nil {
		return err
	}
	out, err := gc.ListAccounts()
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	writeAccountsTable(cli.stdout, out.Accounts)
	return nil
}

func (cli CLI) runAccountsGet(gc GBPClient, args []string) error {
	parsedArgs, err := parseAccountArgs(args, "accounts get")
	if err != nil {
		return err
	}
	out, err := gc.GetAccount(parsedArgs.First("account"))
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	writeAccount(cli.stdout, out)
	return nil
}

func (cli CLI) runLocationsList(gc GBPClient, args []string) error {
	parsedArgs, err := parseAccountReadMaskArgs(args, "locations list")
	if err != nil {
		return err
	}
	out, err := gc.ListLocations(parsedArgs.First("account"), parsedArgs.First("read-mask"))
	if err != nil {
		return err
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, out)
	}
	writeLocationsTable(cli.stdout, out.Locations)
	return nil
}

func (cli CLI) runLocationsGet(gc GBPClient, args []string) error {
	parsedArgs, err := parseLocationReadMaskArgs(args, "locations get")
	if err != nil {
		return err
	}
	out, err := gc.GetLocation(parsedArgs.First("location"), parsedArgs.First("read-mask"))
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runLocationsCreate(gc GBPClient, args []string) error {
	parsedArgs, body, err := parseAccountRequestFileArgs(args, "locations create", map[string]argparse.Spec{
		"request-id":    {TakesValue: true},
		"validate-only": {TakesValue: true},
	})
	if err != nil {
		return err
	}
	out, err := gc.CreateLocation(parsedArgs.First("account"), body, locationWriteOptions{
		RequestID:    parsedArgs.First("request-id"),
		ValidateOnly: parsedArgs.First("validate-only"),
	})
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runLocationsPatch(gc GBPClient, args []string) error {
	parsedArgs, body, err := parseLocationPatchRequestArgs(args, "locations patch")
	if err != nil {
		return err
	}
	out, err := gc.PatchLocation(parsedArgs.First("location"), body, locationPatchOptions{
		UpdateMask:    parsedArgs.First("update-mask"),
		AttributeMask: parsedArgs.First("attribute-mask"),
		ValidateOnly:  parsedArgs.First("validate-only"),
	})
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runLocationsDelete(gc GBPClient, args []string) error {
	parsedArgs, err := parseLocationOnlyArgs(args, "locations delete")
	if err != nil {
		return err
	}
	out, err := gc.DeleteLocation(parsedArgs.First("location"))
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runReviewsList(gc GBPClient, args []string) error {
	parsedArgs, options, err := parseAccountLocationPageArgs(args, "reviews list", true)
	if err != nil {
		return err
	}
	out, err := gc.ListReviews(parsedArgs.First("account"), parsedArgs.First("location"), options)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runReviewsGet(gc GBPClient, args []string) error {
	parsedArgs, err := parseAccountLocationReviewArgs(args, "reviews get")
	if err != nil {
		return err
	}
	out, err := gc.GetReview(parsedArgs.First("account"), parsedArgs.First("location"), parsedArgs.First("review"))
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runReviewsUpdateReply(gc GBPClient, args []string) error {
	parsedArgs, body, err := parseAccountLocationReviewRequestArgs(args, "reviews update-reply")
	if err != nil {
		return err
	}
	out, err := gc.UpdateReviewReply(parsedArgs.First("account"), parsedArgs.First("location"), parsedArgs.First("review"), body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runReviewsDeleteReply(gc GBPClient, args []string) error {
	parsedArgs, err := parseAccountLocationReviewArgs(args, "reviews delete-reply")
	if err != nil {
		return err
	}
	out, err := gc.DeleteReviewReply(parsedArgs.First("account"), parsedArgs.First("location"), parsedArgs.First("review"))
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runMediaList(gc GBPClient, args []string) error {
	parsedArgs, options, err := parseAccountLocationPageArgs(args, "media list", false)
	if err != nil {
		return err
	}
	out, err := gc.ListMedia(parsedArgs.First("account"), parsedArgs.First("location"), options)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runMediaCreate(gc GBPClient, args []string) error {
	parsedArgs, body, err := parseAccountLocationRequestFileArgs(args, "media create")
	if err != nil {
		return err
	}
	out, err := gc.CreateMedia(parsedArgs.First("account"), parsedArgs.First("location"), body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runMediaGet(gc GBPClient, args []string) error {
	parsedArgs, err := parseNamedResourceArgs(args, "media get", "media")
	if err != nil {
		return err
	}
	out, err := gc.GetMedia(parsedArgs.First("media"))
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runMediaPatch(gc GBPClient, args []string) error {
	parsedArgs, body, err := parseNamedPatchRequestArgs(args, "media patch", "media")
	if err != nil {
		return err
	}
	out, err := gc.PatchMedia(parsedArgs.First("media"), parsedArgs.First("update-mask"), body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runMediaDelete(gc GBPClient, args []string) error {
	parsedArgs, err := parseNamedResourceArgs(args, "media delete", "media")
	if err != nil {
		return err
	}
	out, err := gc.DeleteMedia(parsedArgs.First("media"))
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runMediaStartUpload(gc GBPClient, args []string) error {
	parsedArgs, _, err := parseAccountLocationPageArgs(args, "media start-upload", false)
	if err != nil {
		return err
	}
	out, err := gc.StartMediaUpload(parsedArgs.First("account"), parsedArgs.First("location"))
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runLocalPostsList(gc GBPClient, args []string) error {
	parsedArgs, options, err := parseAccountLocationPageArgs(args, "local-posts list", false)
	if err != nil {
		return err
	}
	out, err := gc.ListLocalPosts(parsedArgs.First("account"), parsedArgs.First("location"), options)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runLocalPostsCreate(gc GBPClient, args []string) error {
	parsedArgs, body, err := parseAccountLocationRequestFileArgs(args, "local-posts create")
	if err != nil {
		return err
	}
	out, err := gc.CreateLocalPost(parsedArgs.First("account"), parsedArgs.First("location"), body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runLocalPostsGet(gc GBPClient, args []string) error {
	parsedArgs, err := parseNamedResourceArgs(args, "local-posts get", "local-post")
	if err != nil {
		return err
	}
	out, err := gc.GetLocalPost(parsedArgs.First("local-post"))
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runLocalPostsPatch(gc GBPClient, args []string) error {
	parsedArgs, body, err := parseNamedPatchRequestArgs(args, "local-posts patch", "local-post")
	if err != nil {
		return err
	}
	out, err := gc.PatchLocalPost(parsedArgs.First("local-post"), parsedArgs.First("update-mask"), body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runLocalPostsDelete(gc GBPClient, args []string) error {
	parsedArgs, err := parseNamedResourceArgs(args, "local-posts delete", "local-post")
	if err != nil {
		return err
	}
	out, err := gc.DeleteLocalPost(parsedArgs.First("local-post"))
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runLocalPostsReportInsights(gc GBPClient, args []string) error {
	parsedArgs, body, err := parseAccountLocationRequestFileArgs(args, "local-posts report-insights")
	if err != nil {
		return err
	}
	out, err := gc.ReportLocalPostInsights(parsedArgs.First("account"), parsedArgs.First("location"), body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runPerformanceDailyMetrics(gc GBPClient, args []string) error {
	location, body, err := parseLocationRequestFileArgs(args, "performance daily-metrics")
	if err != nil {
		return err
	}
	out, err := gc.GetDailyMetrics(location, body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runPerformanceSearchKeywords(gc GBPClient, args []string) error {
	location, body, err := parseLocationRequestFileArgs(args, "performance search-keywords")
	if err != nil {
		return err
	}
	out, err := gc.ListSearchKeywords(location, body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func parseNoPositionalArgs(args []string, command string, specs map[string]argparse.Spec) (argparse.Parsed, error) {
	parsedArgs, err := argparse.Parse(args, specs)
	if err != nil {
		return argparse.Parsed{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return argparse.Parsed{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	return parsedArgs, nil
}

func parseAccountArgs(args []string, command string) (argparse.Parsed, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{"account": {TakesValue: true}, "json": {}})
	if err != nil {
		return argparse.Parsed{}, err
	}
	if parsedArgs.First("account") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --account", command)
	}
	return parsedArgs, nil
}

func parseAccountReadMaskArgs(args []string, command string) (argparse.Parsed, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		"account":   {TakesValue: true},
		"read-mask": {TakesValue: true},
		"json":      {},
	})
	if err != nil {
		return argparse.Parsed{}, err
	}
	if parsedArgs.First("account") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --account", command)
	}
	if parsedArgs.First("read-mask") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --read-mask", command)
	}
	return parsedArgs, nil
}

func parseLocationReadMaskArgs(args []string, command string) (argparse.Parsed, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		"location":  {TakesValue: true},
		"read-mask": {TakesValue: true},
	})
	if err != nil {
		return argparse.Parsed{}, err
	}
	if parsedArgs.First("location") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --location", command)
	}
	if parsedArgs.First("read-mask") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --read-mask", command)
	}
	return parsedArgs, nil
}

func parseLocationOnlyArgs(args []string, command string) (argparse.Parsed, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		"location": {TakesValue: true},
	})
	if err != nil {
		return argparse.Parsed{}, err
	}
	if parsedArgs.First("location") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --location", command)
	}
	return parsedArgs, nil
}

func parseAccountLocationPageArgs(args []string, command string, includeOrderBy bool) (argparse.Parsed, pageOptions, error) {
	specs := map[string]argparse.Spec{
		"account":    {TakesValue: true},
		"location":   {TakesValue: true},
		"page-size":  {TakesValue: true},
		"page-token": {TakesValue: true},
	}
	if includeOrderBy {
		specs["order-by"] = argparse.Spec{TakesValue: true}
	}
	parsedArgs, err := parseNoPositionalArgs(args, command, specs)
	if err != nil {
		return argparse.Parsed{}, pageOptions{}, err
	}
	if parsedArgs.First("account") == "" {
		return argparse.Parsed{}, pageOptions{}, fmt.Errorf("%s requires --account", command)
	}
	if parsedArgs.First("location") == "" {
		return argparse.Parsed{}, pageOptions{}, fmt.Errorf("%s requires --location", command)
	}
	return parsedArgs, pageOptions{
		PageSize:  parsedArgs.First("page-size"),
		PageToken: parsedArgs.First("page-token"),
		OrderBy:   parsedArgs.First("order-by"),
	}, nil
}

func parseAccountRequestFileArgs(args []string, command string, extraSpecs map[string]argparse.Spec) (argparse.Parsed, json.RawMessage, error) {
	specs := map[string]argparse.Spec{
		"account":      {TakesValue: true},
		"request-file": {TakesValue: true},
	}
	for name, spec := range extraSpecs {
		specs[name] = spec
	}
	parsedArgs, err := parseNoPositionalArgs(args, command, specs)
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	if parsedArgs.First("account") == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --account", command)
	}
	body, err := readRequestFile(parsedArgs.First("request-file"), command)
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	return parsedArgs, body, nil
}

func parseAccountLocationRequestFileArgs(args []string, command string) (argparse.Parsed, json.RawMessage, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		"account":      {TakesValue: true},
		"location":     {TakesValue: true},
		"request-file": {TakesValue: true},
	})
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	if parsedArgs.First("account") == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --account", command)
	}
	if parsedArgs.First("location") == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --location", command)
	}
	body, err := readRequestFile(parsedArgs.First("request-file"), command)
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	return parsedArgs, body, nil
}

func parseLocationPatchRequestArgs(args []string, command string) (argparse.Parsed, json.RawMessage, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		"location":       {TakesValue: true},
		"update-mask":    {TakesValue: true},
		"attribute-mask": {TakesValue: true},
		"validate-only":  {TakesValue: true},
		"request-file":   {TakesValue: true},
	})
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	if parsedArgs.First("location") == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --location", command)
	}
	if parsedArgs.First("update-mask") == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --update-mask", command)
	}
	body, err := readRequestFile(parsedArgs.First("request-file"), command)
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	return parsedArgs, body, nil
}

func parseAccountLocationReviewArgs(args []string, command string) (argparse.Parsed, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		"account":  {TakesValue: true},
		"location": {TakesValue: true},
		"review":   {TakesValue: true},
	})
	if err != nil {
		return argparse.Parsed{}, err
	}
	if parsedArgs.First("account") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --account", command)
	}
	if parsedArgs.First("location") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --location", command)
	}
	if parsedArgs.First("review") == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --review", command)
	}
	return parsedArgs, nil
}

func parseAccountLocationReviewRequestArgs(args []string, command string) (argparse.Parsed, json.RawMessage, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		"account":      {TakesValue: true},
		"location":     {TakesValue: true},
		"review":       {TakesValue: true},
		"request-file": {TakesValue: true},
	})
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	if parsedArgs.First("account") == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --account", command)
	}
	if parsedArgs.First("location") == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --location", command)
	}
	if parsedArgs.First("review") == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --review", command)
	}
	body, err := readRequestFile(parsedArgs.First("request-file"), command)
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	return parsedArgs, body, nil
}

func parseNamedResourceArgs(args []string, command string, flag string) (argparse.Parsed, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		flag: {TakesValue: true},
	})
	if err != nil {
		return argparse.Parsed{}, err
	}
	if parsedArgs.First(flag) == "" {
		return argparse.Parsed{}, fmt.Errorf("%s requires --%s", command, flag)
	}
	return parsedArgs, nil
}

func parseNamedPatchRequestArgs(args []string, command string, flag string) (argparse.Parsed, json.RawMessage, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		flag:           {TakesValue: true},
		"update-mask":  {TakesValue: true},
		"request-file": {TakesValue: true},
	})
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	if parsedArgs.First(flag) == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --%s", command, flag)
	}
	if parsedArgs.First("update-mask") == "" {
		return argparse.Parsed{}, nil, fmt.Errorf("%s requires --update-mask", command)
	}
	body, err := readRequestFile(parsedArgs.First("request-file"), command)
	if err != nil {
		return argparse.Parsed{}, nil, err
	}
	return parsedArgs, body, nil
}

func parseLocationRequestFileArgs(args []string, command string) (string, json.RawMessage, error) {
	parsedArgs, err := parseNoPositionalArgs(args, command, map[string]argparse.Spec{
		"location":     {TakesValue: true},
		"request-file": {TakesValue: true},
	})
	if err != nil {
		return "", nil, err
	}
	location := parsedArgs.First("location")
	if location == "" {
		return "", nil, fmt.Errorf("%s requires --location", command)
	}
	body, err := readRequestFile(parsedArgs.First("request-file"), command)
	if err != nil {
		return "", nil, err
	}
	return location, body, nil
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

func writeAccountsTable(writer io.Writer, accounts []GBPAccount) {
	fmt.Fprintln(writer, "NAME\tACCOUNT_NAME\tTYPE\tROLE")
	for _, account := range accounts {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", account.Name, account.AccountName, account.Type, account.Role)
	}
}

func writeAccount(writer io.Writer, account GBPAccount) {
	fmt.Fprintln(writer, "Name: "+account.Name)
	if account.AccountName != "" {
		fmt.Fprintln(writer, "Account Name: "+account.AccountName)
	}
	if account.Type != "" {
		fmt.Fprintln(writer, "Type: "+account.Type)
	}
	if account.Role != "" {
		fmt.Fprintln(writer, "Role: "+account.Role)
	}
}

func writeLocationsTable(writer io.Writer, locations []GBPLocation) {
	fmt.Fprintln(writer, "NAME\tTITLE\tSTORE_CODE")
	for _, location := range locations {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", stringField(location, "name"), stringField(location, "title"), stringField(location, "storeCode"))
	}
}

func stringField(value map[string]any, key string) string {
	field, ok := value[key].(string)
	if !ok {
		return ""
	}
	return field
}

func isHelp(arg string) bool {
	return arg == "help" || arg == "-h" || arg == "--help"
}

func (cli CLI) printHelp() {
	fmt.Fprint(cli.stdout, `gbp

CLI for Google Business Profile.

Usage:
  gbp help
  gbp version
  gbp auth help
  gbp accounts help
  gbp locations help
  gbp reviews help
  gbp media help
  gbp local-posts help
  gbp performance help
  gbp mcp help

Commands:
  help
  version
  auth
  accounts
  locations
  reviews
  media
  local-posts
  performance
  mcp
`)
}

func (cli CLI) printAuthHelp() {
	fmt.Fprint(cli.stdout, `gbp auth

Inspect Google Business Profile API access.

Usage:
  gbp auth help
  gbp auth test [--json]
`)
}

func (cli CLI) printAccountsHelp() {
	fmt.Fprint(cli.stdout, `gbp accounts

Inspect Google Business Profile accounts.

Usage:
  gbp accounts help
  gbp accounts list [--json]
  gbp accounts get --account accounts/<id> [--json]
`)
}

func (cli CLI) printLocationsHelp() {
	fmt.Fprint(cli.stdout, `gbp locations

Inspect Google Business Profile locations.

Usage:
  gbp locations help
  gbp locations list --account accounts/<id> --read-mask <fields> [--json]
  gbp locations get --location locations/<id> --read-mask <fields>
  gbp locations create --account accounts/<id> --request-file <json> [--request-id <id>] [--validate-only true]
  gbp locations patch --location locations/<id> --update-mask <fields> --request-file <json> [--attribute-mask <fields>] [--validate-only true]
  gbp locations delete --location locations/<id>
`)
}

func (cli CLI) printReviewsHelp() {
	fmt.Fprint(cli.stdout, `gbp reviews

Inspect Google Business Profile reviews.

Usage:
  gbp reviews help
  gbp reviews list --account accounts/<id> --location locations/<id> [--page-size <n>] [--page-token <token>] [--order-by <order>]
  gbp reviews get --account accounts/<id> --location locations/<id> --review <review-id>
  gbp reviews update-reply --account accounts/<id> --location locations/<id> --review <review-id> --request-file <json>
  gbp reviews delete-reply --account accounts/<id> --location locations/<id> --review <review-id>
`)
}

func (cli CLI) printMediaHelp() {
	fmt.Fprint(cli.stdout, `gbp media

Inspect Google Business Profile media.

Usage:
  gbp media help
  gbp media list --account accounts/<id> --location locations/<id> [--page-size <n>] [--page-token <token>]
  gbp media create --account accounts/<id> --location locations/<id> --request-file <json>
  gbp media get --media accounts/<id>/locations/<id>/media/<id>
  gbp media patch --media accounts/<id>/locations/<id>/media/<id> --update-mask <fields> --request-file <json>
  gbp media delete --media accounts/<id>/locations/<id>/media/<id>
  gbp media start-upload --account accounts/<id> --location locations/<id>
`)
}

func (cli CLI) printLocalPostsHelp() {
	fmt.Fprint(cli.stdout, `gbp local-posts

Inspect Google Business Profile local posts.

Usage:
  gbp local-posts help
  gbp local-posts list --account accounts/<id> --location locations/<id> [--page-size <n>] [--page-token <token>]
  gbp local-posts create --account accounts/<id> --location locations/<id> --request-file <json>
  gbp local-posts get --local-post accounts/<id>/locations/<id>/localPosts/<id>
  gbp local-posts patch --local-post accounts/<id>/locations/<id>/localPosts/<id> --update-mask <fields> --request-file <json>
  gbp local-posts delete --local-post accounts/<id>/locations/<id>/localPosts/<id>
  gbp local-posts report-insights --account accounts/<id> --location locations/<id> --request-file <json>
`)
}

func (cli CLI) printPerformanceHelp() {
	fmt.Fprint(cli.stdout, `gbp performance

Inspect Google Business Profile performance data.

Usage:
  gbp performance help
  gbp performance daily-metrics --location locations/<id> --request-file <json>
  gbp performance search-keywords --location locations/<id> --request-file <json>
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

func containsAccount(accounts []GBPAccount, account string) bool {
	for _, candidate := range accounts {
		if strings.TrimSpace(candidate.Name) == account {
			return true
		}
	}
	return false
}
