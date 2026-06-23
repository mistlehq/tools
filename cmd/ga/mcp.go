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
	defaultMCPAddr     = "127.0.0.1:7347"
	defaultMCPEndpoint = "/mcp"
)

type gaMCPConfig struct {
	Addr     string
	Endpoint string
}

type gaEmptyToolInput struct{}

type gaPropertyToolInput struct {
	Property string `json:"property" jsonschema:"Google Analytics property resource name, for example properties/123456789."`
}

type gaRequestToolInput struct {
	Property string         `json:"property" jsonschema:"Google Analytics property resource name, for example properties/123456789."`
	Request  map[string]any `json:"request" jsonschema:"Google Analytics API request body using the documented JSON shape for this operation."`
}

func (cli CLI) runMCP(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printMCPHelp()
		return nil
	}

	switch args[0] {
	case "serve":
		if len(args) == 2 && (args[1] == "help" || args[1] == "-h" || args[1] == "--help") {
			cli.printMCPServeHelp()
			return nil
		}
		return cli.runMCPServe(args[1:])
	default:
		return fmt.Errorf("unsupported mcp command: %s", args[0])
	}
}

func (cli CLI) runMCPServe(args []string) error {
	config, err := parseGAMCPServeArgs(args)
	if err != nil {
		return err
	}
	gc, err := cli.gaClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newGAMCPHTTPHandler(gc))

	fmt.Fprintf(cli.stdout, "Google Analytics MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseGAMCPServeArgs(args []string) (gaMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return gaMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return gaMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}

	config := gaMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return gaMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return gaMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newGAMCPHTTPHandler(gc GAClient) http.Handler {
	server := newGAMCPServer(gc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newGAMCPServer(gc GAClient) *mcp.Server {
	openWorld := true
	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	server := mcp.NewServer(&mcp.Implementation{Name: "ga", Version: Version}, nil)

	mcp.AddTool(server, gaTool("ga_auth_test", gaAuthTestDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gaPropertyToolInput) (*mcp.CallToolResult, GAProperty, error) {
		out, err := gc.AuthTestContext(ctx, input.Property)
		return nil, out, err
	})
	mcp.AddTool(server, gaTool("ga_account_summaries_list", gaAccountSummariesListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *gaEmptyToolInput) (*mcp.CallToolResult, GAAccountSummariesList, error) {
		out, err := gc.ListAccountSummariesContext(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, gaTool("ga_property_get", gaPropertyGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gaPropertyToolInput) (*mcp.CallToolResult, GAProperty, error) {
		out, err := gc.GetPropertyContext(ctx, input.Property)
		return nil, out, err
	})
	mcp.AddTool(server, gaTool("ga_metadata_get", gaMetadataGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gaPropertyToolInput) (*mcp.CallToolResult, GAMetadata, error) {
		out, err := gc.GetMetadataContext(ctx, input.Property)
		return nil, out, err
	})
	mcp.AddTool(server, gaTool("ga_compatibility_check", gaCompatibilityCheckDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gaRequestToolInput) (*mcp.CallToolResult, GACompatibilityCheckResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.CheckCompatibilityContext(ctx, input.Property, body)
		return nil, out, err
	})
	mcp.AddTool(server, gaTool("ga_report_run", gaReportRunDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gaRequestToolInput) (*mcp.CallToolResult, GAReportResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.RunReportContext(ctx, input.Property, body)
		return nil, out, err
	})
	mcp.AddTool(server, gaTool("ga_report_realtime", gaReportRealtimeDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gaRequestToolInput) (*mcp.CallToolResult, GAReportResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.RunRealtimeReportContext(ctx, input.Property, body)
		return nil, out, err
	})
	mcp.AddTool(server, gaTool("ga_report_funnel", gaReportFunnelDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gaRequestToolInput) (*mcp.CallToolResult, GAReportResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.RunFunnelReportContext(ctx, input.Property, body)
		return nil, out, err
	})
	mcp.AddTool(server, gaTool("ga_google_ads_links_list", gaGoogleAdsLinksListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gaPropertyToolInput) (*mcp.CallToolResult, GAGoogleAdsLinksList, error) {
		out, err := gc.ListGoogleAdsLinksContext(ctx, input.Property)
		return nil, out, err
	})

	return server
}

func gaTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
	return &mcp.Tool{Name: name, Title: doc.Command, Description: doc.Description, Annotations: annotations}
}

func marshalMCPRequestBody(request map[string]any) (json.RawMessage, error) {
	if request == nil {
		return nil, fmt.Errorf("request is required")
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

func (cli CLI) printMCPHelp() {
	fmt.Fprint(cli.stdout, `ga mcp

Run Google Analytics as a local MCP server.

Usage:
  ga mcp help
  ga mcp serve
  ga mcp serve --help

Commands:
  serve    Serve Google Analytics MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `ga mcp serve

Serve Google Analytics tools over MCP Streamable HTTP.

Usage:
  ga mcp serve
  ga mcp serve --addr <addr>
  ga mcp serve --endpoint <path>
  ga mcp serve --addr <addr> --endpoint <path>
  ga mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.

Tools:
  ga_auth_test                    %s
  ga_account_summaries_list       %s
  ga_property_get                 %s
  ga_metadata_get                 %s
  ga_compatibility_check          %s
  ga_report_run                   %s
  ga_report_realtime              %s
  ga_report_funnel                %s
  ga_google_ads_links_list        %s
`, defaultMCPAddr, defaultMCPEndpoint, gaAuthTestDoc.Summary, gaAccountSummariesListDoc.Summary, gaPropertyGetDoc.Summary, gaMetadataGetDoc.Summary, gaCompatibilityCheckDoc.Summary, gaReportRunDoc.Summary, gaReportRealtimeDoc.Summary, gaReportFunnelDoc.Summary, gaGoogleAdsLinksListDoc.Summary)
}
