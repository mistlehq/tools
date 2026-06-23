package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mistlehq/tools/internal/argparse"
	"io"
	"os"
	"strings"
)

// Version is the current shopify CLI version.
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

func (cli CLI) shopifyClient() (ShopifyClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return ShopifyClient{}, err
	}

	return NewShopifyClient(config), nil
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
	case "graphql":
		return cli.runGraphQL(args[2:])
	case "shop":
		return cli.runShop(args[2:])
	case "products":
		return cli.runProducts(args[2:])
	case "orders":
		return cli.runOrders(args[2:])
	case "customers":
		return cli.runCustomers(args[2:])
	case "inventory":
		return cli.runInventory(args[2:])
	case "locations":
		return cli.runLocations(args[2:])
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

	sc, err := cli.shopifyClient()
	if err != nil {
		return err
	}

	shop, err := sc.Shop(cliContext())
	if err != nil {
		return err
	}

	parsedArgs, err := argparse.Parse(args[1:], map[string]argparse.Spec{"json": {}})
	if err != nil {
		return err
	}
	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("auth test does not accept positional arguments")
	}
	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, shop)
	}

	fmt.Fprintln(cli.stdout, "Name: "+shop.Name)
	fmt.Fprintln(cli.stdout, "Myshopify domain: "+shop.MyshopifyDomain)
	return nil
}

func (cli CLI) runGraphQL(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printGraphQLHelp()
		return nil
	}
	if args[0] != "request" {
		return fmt.Errorf("unsupported graphql command: %s", args[0])
	}
	if isSingleHelpArg(args[1:]) {
		cli.printGraphQLRequestHelp()
		return nil
	}

	sc, err := cli.shopifyClient()
	if err != nil {
		return err
	}
	request, err := parseGraphQLRequestArgs(args[1:])
	if err != nil {
		return err
	}
	response, err := sc.GraphQL(request)
	if err != nil {
		return err
	}
	_, err = cli.stdout.Write(append(response, '\n'))
	return err
}

func (cli CLI) runShop(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printShopHelp()
		return nil
	}
	if args[0] != "get" {
		return fmt.Errorf("unsupported shop command: %s", args[0])
	}
	if isSingleHelpArg(args[1:]) {
		cli.printShopGetHelp()
		return nil
	}
	if len(args[1:]) > 0 {
		return fmt.Errorf("shop get does not accept positional arguments")
	}
	sc, err := cli.shopifyClient()
	if err != nil {
		return err
	}
	shop, err := sc.Shop(cliContext())
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, shop)
}

func (cli CLI) runProducts(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printProductsHelp()
		return nil
	}

	sc, err := cli.shopifyClient()
	if err != nil {
		return err
	}

	switch args[0] {
	case "search":
		if isSingleHelpArg(args[1:]) {
			cli.printProductsSearchHelp()
			return nil
		}
		input, err := parseSearchArgs(args[1:], "products search")
		if err != nil {
			return err
		}
		out, err := sc.SearchProducts(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		if isSingleHelpArg(args[1:]) {
			cli.printProductGetHelp()
			return nil
		}
		input, err := parseProductGetArgs(args[1:])
		if err != nil {
			return err
		}
		out, err := sc.GetProduct(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "create":
		if isSingleHelpArg(args[1:]) {
			cli.printProductCreateHelp()
			return nil
		}
		product, err := parseObjectInputArgs(args[1:], "products create", "product")
		if err != nil {
			return err
		}
		out, err := sc.CreateProduct(cliContext(), product)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "update":
		if isSingleHelpArg(args[1:]) {
			cli.printProductUpdateHelp()
			return nil
		}
		product, err := parseObjectInputArgs(args[1:], "products update", "product")
		if err != nil {
			return err
		}
		out, err := sc.UpdateProduct(cliContext(), product)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "delete":
		if isSingleHelpArg(args[1:]) {
			cli.printProductDeleteHelp()
			return nil
		}
		id, err := parseIDArg(args[1:], "products delete")
		if err != nil {
			return err
		}
		out, err := sc.DeleteProduct(cliContext(), id)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported products command: %s", args[0])
	}
}

func (cli CLI) runOrders(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printOrdersHelp()
		return nil
	}
	sc, err := cli.shopifyClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "search":
		if isSingleHelpArg(args[1:]) {
			cli.printOrdersSearchHelp()
			return nil
		}
		input, err := parseSearchArgs(args[1:], "orders search")
		if err != nil {
			return err
		}
		out, err := sc.SearchOrders(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		if isSingleHelpArg(args[1:]) {
			cli.printOrderGetHelp()
			return nil
		}
		id, err := parseIDArg(args[1:], "orders get")
		if err != nil {
			return err
		}
		out, err := sc.GetOrder(cliContext(), id)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported orders command: %s", args[0])
	}
}

func (cli CLI) runCustomers(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printCustomersHelp()
		return nil
	}
	sc, err := cli.shopifyClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "search":
		if isSingleHelpArg(args[1:]) {
			cli.printCustomersSearchHelp()
			return nil
		}
		input, err := parseSearchArgs(args[1:], "customers search")
		if err != nil {
			return err
		}
		out, err := sc.SearchCustomers(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		if isSingleHelpArg(args[1:]) {
			cli.printCustomerGetHelp()
			return nil
		}
		id, err := parseIDArg(args[1:], "customers get")
		if err != nil {
			return err
		}
		out, err := sc.GetCustomer(cliContext(), id)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported customers command: %s", args[0])
	}
}

func (cli CLI) runInventory(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printInventoryHelp()
		return nil
	}
	sc, err := cli.shopifyClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "items":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printInventoryItemsHelp()
			return nil
		}
		if args[1] != "search" {
			return fmt.Errorf("unsupported inventory items command: %s", args[1])
		}
		input, err := parseSearchArgs(args[2:], "inventory items search")
		if err != nil {
			return err
		}
		out, err := sc.SearchInventoryItems(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "levels":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printInventoryLevelsHelp()
			return nil
		}
		if args[1] != "search" {
			return fmt.Errorf("unsupported inventory levels command: %s", args[1])
		}
		input, err := parseSearchArgs(args[2:], "inventory levels search")
		if err != nil {
			return err
		}
		out, err := sc.SearchInventoryLevels(cliContext(), input)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported inventory command: %s", args[0])
	}
}

func (cli CLI) runLocations(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printLocationsHelp()
		return nil
	}
	if args[0] != "list" {
		return fmt.Errorf("unsupported locations command: %s", args[0])
	}
	if isSingleHelpArg(args[1:]) {
		cli.printLocationsListHelp()
		return nil
	}
	input, err := parsePaginationArgs(args[1:], "locations list")
	if err != nil {
		return err
	}
	sc, err := cli.shopifyClient()
	if err != nil {
		return err
	}
	out, err := sc.ListLocations(cliContext(), input)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func parseGraphQLRequestArgs(args []string) (ShopifyGraphQLRequest, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"query":          {TakesValue: true},
		"query-file":     {TakesValue: true},
		"variables":      {TakesValue: true},
		"variables-file": {TakesValue: true},
	})
	if err != nil {
		return ShopifyGraphQLRequest{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return ShopifyGraphQLRequest{}, fmt.Errorf("graphql request does not accept positional arguments")
	}

	query, err := oneStringOrFile(parsedArgs.First("query"), parsedArgs.First("query-file"), "query")
	if err != nil {
		return ShopifyGraphQLRequest{}, err
	}

	variables, err := optionalObjectJSONOrFile(parsedArgs.First("variables"), parsedArgs.First("variables-file"), "variables")
	if err != nil {
		return ShopifyGraphQLRequest{}, err
	}

	return ShopifyGraphQLRequest{Query: query, Variables: variables}, nil
}

func parseSearchArgs(args []string, command string) (ShopifySearchInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"first": {TakesValue: true},
		"after": {TakesValue: true},
		"query": {TakesValue: true},
	})
	if err != nil {
		return ShopifySearchInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return ShopifySearchInput{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	pagination, err := paginationFromParsed(parsedArgs, command)
	if err != nil {
		return ShopifySearchInput{}, err
	}
	return ShopifySearchInput{First: pagination.First, After: pagination.After, Query: parsedArgs.First("query")}, nil
}

func parsePaginationArgs(args []string, command string) (ShopifyPaginationInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"first": {TakesValue: true},
		"after": {TakesValue: true},
	})
	if err != nil {
		return ShopifyPaginationInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return ShopifyPaginationInput{}, fmt.Errorf("%s does not accept positional arguments", command)
	}
	return paginationFromParsed(parsedArgs, command)
}

func paginationFromParsed(parsedArgs argparse.Parsed, command string) (ShopifyPaginationInput, error) {
	firstString := parsedArgs.First("first")
	if firstString == "" {
		return ShopifyPaginationInput{}, fmt.Errorf("%s requires --first", command)
	}
	var first int
	if _, err := fmt.Sscanf(firstString, "%d", &first); err != nil {
		return ShopifyPaginationInput{}, fmt.Errorf("--first must be an integer")
	}
	if first <= 0 {
		return ShopifyPaginationInput{}, fmt.Errorf("--first must be greater than zero")
	}
	return ShopifyPaginationInput{First: first, After: parsedArgs.First("after")}, nil
}

func parseProductGetArgs(args []string) (ShopifyProductGetInput, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"id":     {TakesValue: true},
		"handle": {TakesValue: true},
	})
	if err != nil {
		return ShopifyProductGetInput{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return ShopifyProductGetInput{}, fmt.Errorf("products get does not accept positional arguments")
	}
	input := ShopifyProductGetInput{ID: parsedArgs.First("id"), Handle: parsedArgs.First("handle")}
	if (input.ID == "") == (input.Handle == "") {
		return ShopifyProductGetInput{}, fmt.Errorf("products get requires exactly one of --id or --handle")
	}
	return input, nil
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

func parseObjectInputArgs(args []string, command string, name string) (map[string]any, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		name + "-json": {TakesValue: true},
		name + "-file": {TakesValue: true},
	})
	if err != nil {
		return nil, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	return oneObjectJSONOrFile(parsedArgs.First(name+"-json"), parsedArgs.First(name+"-file"), name)
}

func oneStringOrFile(inline string, path string, name string) (string, error) {
	if (inline == "") == (path == "") {
		return "", fmt.Errorf("exactly one of --%s or --%s-file is required", name, name)
	}
	if inline != "" {
		return inline, nil
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func optionalObjectJSONOrFile(inline string, path string, name string) (map[string]any, error) {
	if inline == "" && path == "" {
		return nil, nil
	}
	return oneObjectJSONOrFile(inline, path, name)
}

func oneObjectJSONOrFile(inline string, path string, name string) (map[string]any, error) {
	raw, err := oneStringOrFile(inline, path, name)
	if err != nil {
		return nil, err
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
