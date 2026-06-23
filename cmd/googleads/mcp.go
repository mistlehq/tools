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
	defaultMCPAddr     = "127.0.0.1:7351"
	defaultMCPEndpoint = "/mcp"
)

type googleAdsMCPConfig struct {
	Addr     string
	Endpoint string
}

type googleAdsEmptyToolInput struct{}

type googleAdsRequestToolInput struct {
	Method string         `json:"method,omitempty" jsonschema:"HTTP method: GET, POST, PATCH, or DELETE. Defaults to GET."`
	Path   string         `json:"path" jsonschema:"Google Ads API path under the configured version base, starting with '/'."`
	Params map[string]any `json:"params,omitempty" jsonschema:"Optional query parameters object."`
	Body   map[string]any `json:"body,omitempty" jsonschema:"Optional JSON request body."`
}

type googleAdsGAQLToolInput struct {
	CustomerID string         `json:"customer_id" jsonschema:"Google Ads customer ID without dashes."`
	Query      string         `json:"query" jsonschema:"Google Ads Query Language query."`
	PageSize   string         `json:"page_size,omitempty" jsonschema:"Optional pageSize value for search."`
	PageToken  string         `json:"page_token,omitempty" jsonschema:"Optional pageToken value for search."`
	SummaryRow string         `json:"summary_row,omitempty" jsonschema:"Optional summaryRowSetting value."`
	Params     map[string]any `json:"params,omitempty" jsonschema:"Additional documented request body fields."`
}

type googleAdsFieldsSearchToolInput struct {
	Query string `json:"query" jsonschema:"Google Ads Query Language query over googleAdsFields."`
}

type googleAdsFieldGetToolInput struct {
	ResourceName string `json:"resource_name" jsonschema:"GoogleAdsField resource name, for example googleAdsFields/campaign.id."`
}

type googleAdsRawResponse map[string]any

type googleAdsSearchStreamToolOutput struct {
	Results []any `json:"results"`
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
	config, err := parseGoogleAdsMCPServeArgs(args)
	if err != nil {
		return err
	}
	gc, err := cli.googleAdsClient()
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newGoogleAdsMCPHTTPHandler(gc))
	fmt.Fprintf(cli.stdout, "Google Ads MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseGoogleAdsMCPServeArgs(args []string) (googleAdsMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return googleAdsMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return googleAdsMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}
	config := googleAdsMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return googleAdsMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return googleAdsMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newGoogleAdsMCPHTTPHandler(gc GoogleAdsClient) http.Handler {
	server := newGoogleAdsMCPServer(gc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newGoogleAdsMCPServer(gc GoogleAdsClient) *mcp.Server {
	openWorld := true
	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	rawAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld}

	server := mcp.NewServer(&mcp.Implementation{Name: "googleads", Version: Version}, nil)

	mcp.AddTool(server, googleAdsTool("googleads_request", googleAdsRequestDoc, rawAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *googleAdsRequestToolInput) (*mcp.CallToolResult, googleAdsRawResponse, error) {
		body, err := gc.RequestContext(ctx, GoogleAdsRequest{Method: input.Method, Path: input.Path, Params: input.Params, Body: input.Body})
		if err != nil {
			return nil, nil, err
		}
		var out googleAdsRawResponse
		if err := json.Unmarshal(body, &out); err != nil {
			return nil, nil, err
		}
		return nil, out, nil
	})
	mcp.AddTool(server, googleAdsTool("googleads_auth_test", googleAdsAuthTestDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *googleAdsEmptyToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := gc.AuthTest(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, googleAdsTool("googleads_customers_list_accessible", googleAdsCustomersListAccessibleDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *googleAdsEmptyToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := gc.ListAccessibleCustomers(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, googleAdsTool("googleads_gaql_search", googleAdsGAQLSearchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *googleAdsGAQLToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := gc.Search(ctx, gaqlInputFromTool(*input))
		return nil, out, err
	})
	mcp.AddTool(server, googleAdsTool("googleads_gaql_search_stream", googleAdsGAQLSearchStreamDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *googleAdsGAQLToolInput) (*mcp.CallToolResult, googleAdsSearchStreamToolOutput, error) {
		out, err := gc.SearchStream(ctx, gaqlInputFromTool(*input))
		return nil, googleAdsSearchStreamToolOutput{Results: out}, err
	})
	mcp.AddTool(server, googleAdsTool("googleads_fields_search", googleAdsFieldsSearchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *googleAdsFieldsSearchToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := gc.SearchFields(ctx, input.Query)
		return nil, out, err
	})
	mcp.AddTool(server, googleAdsTool("googleads_field_get", googleAdsFieldGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *googleAdsFieldGetToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := gc.GetField(ctx, input.ResourceName)
		return nil, out, err
	})

	return server
}

func googleAdsTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
	return &mcp.Tool{Name: name, Title: doc.Command, Description: doc.Description, Annotations: annotations}
}

func gaqlInputFromTool(input googleAdsGAQLToolInput) GoogleAdsGAQLInput {
	return GoogleAdsGAQLInput{CustomerID: input.CustomerID, Query: input.Query, PageSize: input.PageSize, PageToken: input.PageToken, SummaryRow: input.SummaryRow, Params: input.Params}
}

func (cli CLI) printMCPHelp() {
	fmt.Fprint(cli.stdout, `googleads mcp

Run Google Ads as a local MCP server.

Usage:
  googleads mcp help
  googleads mcp serve
  googleads mcp serve --help

Commands:
  serve    Serve Google Ads MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `googleads mcp serve

Serve Google Ads tools over MCP Streamable HTTP.

Usage:
  googleads mcp serve
  googleads mcp serve --addr <addr>
  googleads mcp serve --endpoint <path>
  googleads mcp serve --addr <addr> --endpoint <path>
  googleads mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.

Tools:
  googleads_request                    %s
  googleads_auth_test                  %s
  googleads_customers_list_accessible  %s
  googleads_gaql_search                %s
  googleads_gaql_search_stream         %s
  googleads_fields_search              %s
  googleads_field_get                  %s
`, defaultMCPAddr, defaultMCPEndpoint, googleAdsRequestDoc.Description, googleAdsAuthTestDoc.Description, googleAdsCustomersListAccessibleDoc.Description, googleAdsGAQLSearchDoc.Description, googleAdsGAQLSearchStreamDoc.Description, googleAdsFieldsSearchDoc.Description, googleAdsFieldGetDoc.Description)
}
