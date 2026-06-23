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

// Version is the current metaads CLI version.
var Version = "dev"

type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	env    Environment
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

func isHelpToken(arg string) bool {
	return arg == "help" || arg == "-h" || arg == "--help"
}

func isSingleHelpArg(args []string) bool {
	return len(args) == 1 && isHelpToken(args[0])
}

func (cli CLI) metaAdsClient() (MetaAdsClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return MetaAdsClient{}, err
	}
	return NewMetaAdsClient(config), nil
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
	case "graph":
		return cli.runGraph(args[2:])
	case "ad-accounts":
		return cli.runAdAccounts(args[2:])
	case "campaigns":
		return cli.runCampaigns(args[2:])
	case "adsets":
		return cli.runAdSets(args[2:])
	case "ads":
		return cli.runAds(args[2:])
	case "creatives":
		return cli.runCreatives(args[2:])
	case "insights":
		return cli.runInsights(args[2:])
	case "targeting":
		return cli.runTargeting(args[2:])
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
	mc, err := cli.metaAdsClient()
	if err != nil {
		return err
	}
	out, err := mc.AuthTest(cliContext())
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runGraph(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printGraphHelp()
		return nil
	}
	if args[0] != "request" {
		return fmt.Errorf("unsupported graph command: %s", args[0])
	}
	if isSingleHelpArg(args[1:]) {
		cli.printGraphRequestHelp()
		return nil
	}
	request, err := parseGraphRequestArgs(args[1:])
	if err != nil {
		return err
	}
	mc, err := cli.metaAdsClient()
	if err != nil {
		return err
	}
	response, err := mc.Request(request)
	if err != nil {
		return err
	}
	_, err = cli.stdout.Write(append(response, '\n'))
	return err
}

func (cli CLI) runAdAccounts(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printAdAccountsHelp()
		return nil
	}
	mc, err := cli.metaAdsClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		if isSingleHelpArg(args[1:]) {
			cli.printAdAccountsListHelp()
			return nil
		}
		input, err := parseEdgeArgs(args[1:], "ad-accounts list", false)
		if err != nil {
			return err
		}
		out, err := mc.ListAdAccounts(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		if isSingleHelpArg(args[1:]) {
			cli.printAdAccountGetHelp()
			return nil
		}
		input, err := parseGetArgs(args[1:], "ad-accounts get")
		if err != nil {
			return err
		}
		out, err := mc.GetAdAccount(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported ad-accounts command: %s", args[0])
	}
}

func (cli CLI) runCampaigns(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printCampaignsHelp()
		return nil
	}
	return cli.runCrudFamily(args, crudFamily{
		name:      "campaigns",
		searchDoc: metaAdsCampaignsSearchDoc,
		getDoc:    metaAdsCampaignGetDoc,
		createDoc: metaAdsCampaignCreateDoc,
		updateDoc: metaAdsCampaignUpdateDoc,
		deleteDoc: metaAdsCampaignDeleteDoc,
		search: func(ctx context.Context, mc MetaAdsClient, input MetaAdsEdgeInput) (map[string]any, error) {
			return mc.SearchCampaigns(ctx, input)
		},
		get: func(ctx context.Context, mc MetaAdsClient, input MetaAdsGetInput) (map[string]any, error) {
			return mc.GetCampaign(ctx, input)
		},
		create: func(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
			return mc.CreateCampaign(ctx, input)
		},
		update: func(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
			return mc.UpdateCampaign(ctx, input)
		},
		deleteByID: func(ctx context.Context, mc MetaAdsClient, id string) (map[string]any, error) {
			return mc.DeleteCampaign(ctx, id)
		},
	})
}

func (cli CLI) runAdSets(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printAdSetsHelp()
		return nil
	}
	return cli.runCrudFamily(args, crudFamily{
		name:      "adsets",
		searchDoc: metaAdsAdSetsSearchDoc,
		getDoc:    metaAdsAdSetGetDoc,
		createDoc: metaAdsAdSetCreateDoc,
		updateDoc: metaAdsAdSetUpdateDoc,
		deleteDoc: metaAdsAdSetDeleteDoc,
		search: func(ctx context.Context, mc MetaAdsClient, input MetaAdsEdgeInput) (map[string]any, error) {
			return mc.SearchAdSets(ctx, input)
		},
		get: func(ctx context.Context, mc MetaAdsClient, input MetaAdsGetInput) (map[string]any, error) {
			return mc.GetAdSet(ctx, input)
		},
		create: func(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
			return mc.CreateAdSet(ctx, input)
		},
		update: func(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
			return mc.UpdateAdSet(ctx, input)
		},
		deleteByID: func(ctx context.Context, mc MetaAdsClient, id string) (map[string]any, error) {
			return mc.DeleteAdSet(ctx, id)
		},
	})
}

func (cli CLI) runAds(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printAdsHelp()
		return nil
	}
	return cli.runCrudFamily(args, crudFamily{
		name:      "ads",
		searchDoc: metaAdsAdsSearchDoc,
		getDoc:    metaAdsAdGetDoc,
		createDoc: metaAdsAdCreateDoc,
		updateDoc: metaAdsAdUpdateDoc,
		deleteDoc: metaAdsAdDeleteDoc,
		search: func(ctx context.Context, mc MetaAdsClient, input MetaAdsEdgeInput) (map[string]any, error) {
			return mc.SearchAds(ctx, input)
		},
		get: func(ctx context.Context, mc MetaAdsClient, input MetaAdsGetInput) (map[string]any, error) {
			return mc.GetAd(ctx, input)
		},
		create: func(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
			return mc.CreateAd(ctx, input)
		},
		update: func(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
			return mc.UpdateAd(ctx, input)
		},
		deleteByID: func(ctx context.Context, mc MetaAdsClient, id string) (map[string]any, error) {
			return mc.DeleteAd(ctx, id)
		},
	})
}

func (cli CLI) runCreatives(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printCreativesHelp()
		return nil
	}
	mc, err := cli.metaAdsClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "search":
		if isSingleHelpArg(args[1:]) {
			cli.printCreativesSearchHelp()
			return nil
		}
		input, err := parseEdgeArgs(args[1:], "creatives search", true)
		if err != nil {
			return err
		}
		out, err := mc.SearchCreatives(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		if isSingleHelpArg(args[1:]) {
			cli.printCreativeGetHelp()
			return nil
		}
		input, err := parseGetArgs(args[1:], "creatives get")
		if err != nil {
			return err
		}
		out, err := mc.GetCreative(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "create":
		if isSingleHelpArg(args[1:]) {
			cli.printCreativeCreateHelp()
			return nil
		}
		input, err := parseAccountBodyArgs(args[1:], "creatives create")
		if err != nil {
			return err
		}
		out, err := mc.CreateCreative(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported creatives command: %s", args[0])
	}
}

type crudFamily struct {
	name       string
	searchDoc  commandDoc
	getDoc     commandDoc
	createDoc  commandDoc
	updateDoc  commandDoc
	deleteDoc  commandDoc
	search     func(context.Context, MetaAdsClient, MetaAdsEdgeInput) (map[string]any, error)
	get        func(context.Context, MetaAdsClient, MetaAdsGetInput) (map[string]any, error)
	create     func(context.Context, MetaAdsClient, MetaAdsObjectInput) (map[string]any, error)
	update     func(context.Context, MetaAdsClient, MetaAdsObjectInput) (map[string]any, error)
	deleteByID func(context.Context, MetaAdsClient, string) (map[string]any, error)
}

func (cli CLI) runCrudFamily(args []string, family crudFamily) error {
	mc, err := cli.metaAdsClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "search":
		if isSingleHelpArg(args[1:]) {
			cli.printSearchHelp(family.name, family.searchDoc)
			return nil
		}
		input, err := parseEdgeArgs(args[1:], family.name+" search", true)
		if err != nil {
			return err
		}
		out, err := family.search(cliContext(), mc, input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		if isSingleHelpArg(args[1:]) {
			cli.printGetHelp(family.name, family.getDoc)
			return nil
		}
		input, err := parseGetArgs(args[1:], family.name+" get")
		if err != nil {
			return err
		}
		out, err := family.get(cliContext(), mc, input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "create":
		if isSingleHelpArg(args[1:]) {
			cli.printCreateHelp(family.name, family.createDoc)
			return nil
		}
		input, err := parseAccountBodyArgs(args[1:], family.name+" create")
		if err != nil {
			return err
		}
		out, err := family.create(cliContext(), mc, input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "update":
		if isSingleHelpArg(args[1:]) {
			cli.printUpdateHelp(family.name, family.updateDoc)
			return nil
		}
		input, err := parseObjectBodyArgs(args[1:], family.name+" update")
		if err != nil {
			return err
		}
		out, err := family.update(cliContext(), mc, input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "delete":
		if isSingleHelpArg(args[1:]) {
			cli.printDeleteHelp(family.name, family.deleteDoc)
			return nil
		}
		id, err := parseIDArg(args[1:], family.name+" delete")
		if err != nil {
			return err
		}
		out, err := family.deleteByID(cliContext(), mc, id)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported %s command: %s", family.name, args[0])
	}
}

func (cli CLI) runInsights(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printInsightsHelp()
		return nil
	}
	if args[0] != "get" {
		return fmt.Errorf("unsupported insights command: %s", args[0])
	}
	if isSingleHelpArg(args[1:]) {
		cli.printInsightsGetHelp()
		return nil
	}
	input, err := parseInsightsArgs(args[1:])
	if err != nil {
		return err
	}
	mc, err := cli.metaAdsClient()
	if err != nil {
		return err
	}
	out, err := mc.GetInsights(cliContext(), input)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runTargeting(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printTargetingHelp()
		return nil
	}
	if args[0] != "search" {
		return fmt.Errorf("unsupported targeting command: %s", args[0])
	}
	if isSingleHelpArg(args[1:]) {
		cli.printTargetingSearchHelp()
		return nil
	}
	input, err := parseTargetingSearchArgs(args[1:])
	if err != nil {
		return err
	}
	mc, err := cli.metaAdsClient()
	if err != nil {
		return err
	}
	out, err := mc.SearchTargeting(cliContext(), input)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func parseGraphRequestArgs(args []string) (MetaAdsRequest, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"method":      {TakesValue: true},
		"path":        {TakesValue: true},
		"params":      {TakesValue: true},
		"params-file": {TakesValue: true},
		"body":        {TakesValue: true},
		"body-file":   {TakesValue: true},
	})
	if err != nil {
		return MetaAdsRequest{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return MetaAdsRequest{}, fmt.Errorf("graph request does not accept positional arguments")
	}
	params, err := optionalObjectJSONOrFile(parsedArgs.First("params"), parsedArgs.First("params-file"), "params")
	if err != nil {
		return MetaAdsRequest{}, err
	}
	body, err := optionalObjectJSONOrFile(parsedArgs.First("body"), parsedArgs.First("body-file"), "body")
	if err != nil {
		return MetaAdsRequest{}, err
	}
	return MetaAdsRequest{
		Method: parsedArgs.First("method"),
		Path:   parsedArgs.First("path"),
		Params: params,
		Body:   body,
	}, nil
}

func parseEdgeArgs(args []string, command string, requireAccount bool) (MetaAdsEdgeInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"account-id":  {TakesValue: true},
		"fields":      {TakesValue: true},
		"limit":       {TakesValue: true},
		"after":       {TakesValue: true},
		"params":      {TakesValue: true},
		"params-file": {TakesValue: true},
	})
	if err != nil {
		return MetaAdsEdgeInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return MetaAdsEdgeInput{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	params, err := optionalObjectJSONOrFile(parsedArgs.First("params"), parsedArgs.First("params-file"), "params")
	if err != nil {
		return MetaAdsEdgeInput{}, err
	}
	input := MetaAdsEdgeInput{
		AccountID: parsedArgs.First("account-id"),
		Fields:    parsedArgs.First("fields"),
		Limit:     parsedArgs.First("limit"),
		After:     parsedArgs.First("after"),
		Params:    params,
	}
	if requireAccount && strings.TrimSpace(input.AccountID) == "" {
		return MetaAdsEdgeInput{}, fmt.Errorf("%s requires --account-id", command)
	}
	return input, nil
}

func parseGetArgs(args []string, command string) (MetaAdsGetInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"id":          {TakesValue: true},
		"fields":      {TakesValue: true},
		"params":      {TakesValue: true},
		"params-file": {TakesValue: true},
	})
	if err != nil {
		return MetaAdsGetInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return MetaAdsGetInput{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	params, err := optionalObjectJSONOrFile(parsedArgs.First("params"), parsedArgs.First("params-file"), "params")
	if err != nil {
		return MetaAdsGetInput{}, err
	}
	id := parsedArgs.First("id")
	if strings.TrimSpace(id) == "" {
		return MetaAdsGetInput{}, fmt.Errorf("%s requires --id", command)
	}
	return MetaAdsGetInput{ID: id, Fields: parsedArgs.First("fields"), Params: params}, nil
}

func parseAccountBodyArgs(args []string, command string) (MetaAdsObjectInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"account-id": {TakesValue: true},
		"body":       {TakesValue: true},
		"body-file":  {TakesValue: true},
	})
	if err != nil {
		return MetaAdsObjectInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return MetaAdsObjectInput{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	body, err := oneObjectJSONOrFile(parsedArgs.First("body"), parsedArgs.First("body-file"), "body")
	if err != nil {
		return MetaAdsObjectInput{}, err
	}
	accountID := parsedArgs.First("account-id")
	if strings.TrimSpace(accountID) == "" {
		return MetaAdsObjectInput{}, fmt.Errorf("%s requires --account-id", command)
	}
	return MetaAdsObjectInput{ID: accountID, Body: body}, nil
}

func parseObjectBodyArgs(args []string, command string) (MetaAdsObjectInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"id":        {TakesValue: true},
		"body":      {TakesValue: true},
		"body-file": {TakesValue: true},
	})
	if err != nil {
		return MetaAdsObjectInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return MetaAdsObjectInput{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	body, err := oneObjectJSONOrFile(parsedArgs.First("body"), parsedArgs.First("body-file"), "body")
	if err != nil {
		return MetaAdsObjectInput{}, err
	}
	id := parsedArgs.First("id")
	if strings.TrimSpace(id) == "" {
		return MetaAdsObjectInput{}, fmt.Errorf("%s requires --id", command)
	}
	return MetaAdsObjectInput{ID: id, Body: body}, nil
}

func parseIDArg(args []string, command string) (string, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{"id": {TakesValue: true}})
	if err != nil {
		return "", err
	}
	if len(parsedArgs.Positionals) > 0 {
		return "", fmt.Errorf("%s does not accept positional arguments", command)
	}
	id := parsedArgs.First("id")
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("%s requires --id", command)
	}
	return id, nil
}

func parseInsightsArgs(args []string) (MetaAdsInsightsInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"id":          {TakesValue: true},
		"fields":      {TakesValue: true},
		"level":       {TakesValue: true},
		"time-range":  {TakesValue: true},
		"params":      {TakesValue: true},
		"params-file": {TakesValue: true},
	})
	if err != nil {
		return MetaAdsInsightsInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return MetaAdsInsightsInput{}, fmt.Errorf("insights get does not accept positional arguments")
	}
	params, err := optionalObjectJSONOrFile(parsedArgs.First("params"), parsedArgs.First("params-file"), "params")
	if err != nil {
		return MetaAdsInsightsInput{}, err
	}
	id := parsedArgs.First("id")
	if strings.TrimSpace(id) == "" {
		return MetaAdsInsightsInput{}, fmt.Errorf("insights get requires --id")
	}
	return MetaAdsInsightsInput{
		ID:        id,
		Fields:    parsedArgs.First("fields"),
		Level:     parsedArgs.First("level"),
		TimeRange: parsedArgs.First("time-range"),
		Params:    params,
	}, nil
}

func parseTargetingSearchArgs(args []string) (MetaAdsTargetingSearchInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"type":        {TakesValue: true},
		"query":       {TakesValue: true},
		"params":      {TakesValue: true},
		"params-file": {TakesValue: true},
	})
	if err != nil {
		return MetaAdsTargetingSearchInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return MetaAdsTargetingSearchInput{}, fmt.Errorf("targeting search does not accept positional arguments")
	}
	params, err := optionalObjectJSONOrFile(parsedArgs.First("params"), parsedArgs.First("params-file"), "params")
	if err != nil {
		return MetaAdsTargetingSearchInput{}, err
	}
	return MetaAdsTargetingSearchInput{Type: parsedArgs.First("type"), Query: parsedArgs.First("query"), Params: params}, nil
}

func optionalObjectJSONOrFile(inline string, path string, name string) (map[string]any, error) {
	if inline == "" && path == "" {
		return nil, nil
	}
	return oneObjectJSONOrFile(inline, path, name)
}

func oneObjectJSONOrFile(inline string, path string, name string) (map[string]any, error) {
	if (inline == "") == (path == "") {
		return nil, fmt.Errorf("exactly one of --%s or --%s-file is required", name, name)
	}
	raw := inline
	if path != "" {
		body, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		raw = string(body)
	}
	var value map[string]any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return nil, fmt.Errorf("%s must be a JSON object: %w", name, err)
	}
	return value, nil
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func cliContext() context.Context {
	return context.Background()
}
