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
	defaultMCPAddr     = "127.0.0.1:7348"
	defaultMCPEndpoint = "/mcp"
)

type shopifyMCPConfig struct {
	Addr     string
	Endpoint string
}

type shopifyEmptyToolInput struct{}

type shopifyGraphQLRequestToolInput struct {
	Query     string         `json:"query" jsonschema:"GraphQL query or mutation text."`
	Variables map[string]any `json:"variables,omitempty" jsonschema:"Optional GraphQL variables object."`
}

type shopifySearchToolInput struct {
	First int    `json:"first" jsonschema:"Number of nodes to fetch from the connection."`
	After string `json:"after,omitempty" jsonschema:"Optional Shopify pagination cursor."`
	Query string `json:"query,omitempty" jsonschema:"Optional Shopify Admin API search query."`
}

type shopifyPaginationToolInput struct {
	First int    `json:"first" jsonschema:"Number of nodes to fetch from the connection."`
	After string `json:"after,omitempty" jsonschema:"Optional Shopify pagination cursor."`
}

type shopifyProductGetToolInput struct {
	ID     string `json:"id,omitempty" jsonschema:"Shopify product GID. Exactly one of id or handle is required."`
	Handle string `json:"handle,omitempty" jsonschema:"Shopify product handle. Exactly one of id or handle is required."`
}

type shopifyIDToolInput struct {
	ID string `json:"id" jsonschema:"Shopify resource GID."`
}

type shopifyProductInputToolInput struct {
	Product map[string]any `json:"product" jsonschema:"Shopify product input object using the documented Admin GraphQL shape for this mutation."`
}

type shopifyRawGraphQLResponse map[string]any

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
	config, err := parseShopifyMCPServeArgs(args)
	if err != nil {
		return err
	}

	sc, err := cli.shopifyClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newShopifyMCPHTTPHandler(sc))

	fmt.Fprintf(cli.stdout, "Shopify MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseShopifyMCPServeArgs(args []string) (shopifyMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return shopifyMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return shopifyMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}

	config := shopifyMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return shopifyMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return shopifyMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newShopifyMCPHTTPHandler(sc ShopifyClient) http.Handler {
	server := newShopifyMCPServer(sc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newShopifyMCPServer(sc ShopifyClient) *mcp.Server {
	openWorld := true
	destructive := true
	notDestructive := false

	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	mutatingAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &notDestructive}
	destructiveAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &destructive}
	rawGraphQLAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld}

	server := mcp.NewServer(&mcp.Implementation{Name: "shopify", Version: Version}, nil)

	mcp.AddTool(server, shopifyTool("shopify_graphql_request", shopifyGraphQLRequestDoc, rawGraphQLAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifyGraphQLRequestToolInput) (*mcp.CallToolResult, shopifyRawGraphQLResponse, error) {
		body, err := sc.GraphQLContext(ctx, ShopifyGraphQLRequest{Query: input.Query, Variables: input.Variables})
		if err != nil {
			return nil, nil, err
		}
		var out shopifyRawGraphQLResponse
		if err := json.Unmarshal(body, &out); err != nil {
			return nil, nil, err
		}
		return nil, out, nil
	})
	mcp.AddTool(server, shopifyTool("shopify_auth_test", shopifyAuthTestDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *shopifyEmptyToolInput) (*mcp.CallToolResult, ShopifyShop, error) {
		out, err := sc.Shop(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_shop_get", shopifyShopGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *shopifyEmptyToolInput) (*mcp.CallToolResult, ShopifyShop, error) {
		out, err := sc.Shop(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_products_search", shopifyProductsSearchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifySearchToolInput) (*mcp.CallToolResult, ShopifyProductsSearch, error) {
		searchInput, err := mcpSearchInput(*input)
		if err != nil {
			return nil, ShopifyProductsSearch{}, err
		}
		out, err := sc.SearchProducts(ctx, searchInput)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_product_get", shopifyProductGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifyProductGetToolInput) (*mcp.CallToolResult, ShopifyProduct, error) {
		out, err := sc.GetProduct(ctx, ShopifyProductGetInput{ID: input.ID, Handle: input.Handle})
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_product_create", shopifyProductCreateDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifyProductInputToolInput) (*mcp.CallToolResult, ShopifyProductCreate, error) {
		out, err := sc.CreateProduct(ctx, input.Product)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_product_update", shopifyProductUpdateDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifyProductInputToolInput) (*mcp.CallToolResult, ShopifyProductUpdate, error) {
		out, err := sc.UpdateProduct(ctx, input.Product)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_product_delete", shopifyProductDeleteDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifyIDToolInput) (*mcp.CallToolResult, ShopifyProductDelete, error) {
		out, err := sc.DeleteProduct(ctx, input.ID)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_orders_search", shopifyOrdersSearchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifySearchToolInput) (*mcp.CallToolResult, ShopifyOrdersSearch, error) {
		searchInput, err := mcpSearchInput(*input)
		if err != nil {
			return nil, ShopifyOrdersSearch{}, err
		}
		out, err := sc.SearchOrders(ctx, searchInput)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_order_get", shopifyOrderGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifyIDToolInput) (*mcp.CallToolResult, ShopifyOrder, error) {
		out, err := sc.GetOrder(ctx, input.ID)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_customers_search", shopifyCustomersSearchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifySearchToolInput) (*mcp.CallToolResult, ShopifyCustomersSearch, error) {
		searchInput, err := mcpSearchInput(*input)
		if err != nil {
			return nil, ShopifyCustomersSearch{}, err
		}
		out, err := sc.SearchCustomers(ctx, searchInput)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_customer_get", shopifyCustomerGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifyIDToolInput) (*mcp.CallToolResult, ShopifyCustomer, error) {
		out, err := sc.GetCustomer(ctx, input.ID)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_inventory_items_search", shopifyInventoryItemsSearchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifySearchToolInput) (*mcp.CallToolResult, ShopifyInventoryItemsSearch, error) {
		searchInput, err := mcpSearchInput(*input)
		if err != nil {
			return nil, ShopifyInventoryItemsSearch{}, err
		}
		out, err := sc.SearchInventoryItems(ctx, searchInput)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_inventory_levels_search", shopifyInventoryLevelsSearchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifySearchToolInput) (*mcp.CallToolResult, ShopifyInventoryLevelsSearch, error) {
		searchInput, err := mcpSearchInput(*input)
		if err != nil {
			return nil, ShopifyInventoryLevelsSearch{}, err
		}
		out, err := sc.SearchInventoryLevels(ctx, searchInput)
		return nil, out, err
	})
	mcp.AddTool(server, shopifyTool("shopify_locations_list", shopifyLocationsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *shopifyPaginationToolInput) (*mcp.CallToolResult, ShopifyLocationsList, error) {
		paginationInput, err := mcpPaginationInput(*input)
		if err != nil {
			return nil, ShopifyLocationsList{}, err
		}
		out, err := sc.ListLocations(ctx, paginationInput)
		return nil, out, err
	})

	return server
}

func shopifyTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
	return &mcp.Tool{Name: name, Title: doc.Command, Description: doc.Description, Annotations: annotations}
}

func mcpSearchInput(input shopifySearchToolInput) (ShopifySearchInput, error) {
	pagination, err := mcpPaginationInput(shopifyPaginationToolInput{First: input.First, After: input.After})
	if err != nil {
		return ShopifySearchInput{}, err
	}
	return ShopifySearchInput{First: pagination.First, After: pagination.After, Query: input.Query}, nil
}

func mcpPaginationInput(input shopifyPaginationToolInput) (ShopifyPaginationInput, error) {
	if input.First <= 0 {
		return ShopifyPaginationInput{}, fmt.Errorf("first must be greater than zero")
	}
	return ShopifyPaginationInput{First: input.First, After: input.After}, nil
}

func (cli CLI) printMCPHelp() {
	fmt.Fprint(cli.stdout, `shopify mcp

Run Shopify as a local MCP server.

Usage:
  shopify mcp help
  shopify mcp serve
  shopify mcp serve --help

Commands:
  serve    Serve Shopify MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `shopify mcp serve

Serve Shopify tools over MCP Streamable HTTP.

Usage:
  shopify mcp serve
  shopify mcp serve --addr <addr>
  shopify mcp serve --endpoint <path>
  shopify mcp serve --addr <addr> --endpoint <path>
  shopify mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.

Tools:
  shopify_graphql_request             %s
  shopify_auth_test                   %s
  shopify_shop_get                    %s
  shopify_products_search             %s
  shopify_product_get                 %s
  shopify_product_create              %s
  shopify_product_update              %s
  shopify_product_delete              %s
  shopify_orders_search               %s
  shopify_order_get                   %s
  shopify_customers_search            %s
  shopify_customer_get                %s
  shopify_inventory_items_search      %s
  shopify_inventory_levels_search     %s
  shopify_locations_list              %s
`, defaultMCPAddr, defaultMCPEndpoint, shopifyGraphQLRequestDoc.Summary, shopifyAuthTestDoc.Summary, shopifyShopGetDoc.Summary, shopifyProductsSearchDoc.Summary, shopifyProductGetDoc.Summary, shopifyProductCreateDoc.Summary, shopifyProductUpdateDoc.Summary, shopifyProductDeleteDoc.Summary, shopifyOrdersSearchDoc.Summary, shopifyOrderGetDoc.Summary, shopifyCustomersSearchDoc.Summary, shopifyCustomerGetDoc.Summary, shopifyInventoryItemsSearchDoc.Summary, shopifyInventoryLevelsSearchDoc.Summary, shopifyLocationsListDoc.Summary)
}
