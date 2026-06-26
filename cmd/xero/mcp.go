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
	defaultMCPAddr     = "127.0.0.1:7355"
	defaultMCPEndpoint = "/mcp"
)

type xeroMCPConfig struct {
	Addr     string
	Endpoint string
}

type xeroEmptyToolInput struct{}

type xeroAPIReadToolInput struct {
	Family   string            `json:"family" jsonschema:"Xero API family: accounting, assets, files, or projects."`
	TenantID string            `json:"tenantId" jsonschema:"Xero tenant ID to send as the xero-tenant-id header."`
	Endpoint string            `json:"endpoint" jsonschema:"Documented API endpoint path within the family, for example /Invoices or /Files."`
	Query    map[string]string `json:"query,omitempty" jsonschema:"Optional documented query parameters."`
}

type xeroAPIWriteToolInput struct {
	Family   string            `json:"family" jsonschema:"Xero API family: accounting, assets, files, or projects."`
	TenantID string            `json:"tenantId" jsonschema:"Xero tenant ID to send as the xero-tenant-id header."`
	Endpoint string            `json:"endpoint" jsonschema:"Documented API endpoint path within the family, for example /Invoices or /Files."`
	Query    map[string]string `json:"query,omitempty" jsonschema:"Optional documented query parameters."`
	Request  map[string]any    `json:"request" jsonschema:"JSON request body using Xero's documented shape for the endpoint."`
}

func (cli CLI) runMCP(args []string) error {
	if len(args) == 0 || isHelp(args[0]) {
		cli.printMCPHelp()
		return nil
	}

	switch args[0] {
	case "serve":
		if len(args) == 2 && isHelp(args[1]) {
			cli.printMCPServeHelp()
			return nil
		}
		return cli.runMCPServe(args[1:])
	default:
		return fmt.Errorf("unsupported mcp command: %s", args[0])
	}
}

func (cli CLI) runMCPServe(args []string) error {
	config, err := parseXeroMCPServeArgs(args)
	if err != nil {
		return err
	}
	xc, err := cli.xeroClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newXeroMCPHTTPHandler(xc))

	fmt.Fprintf(cli.stdout, "Xero MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseXeroMCPServeArgs(args []string) (xeroMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return xeroMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return xeroMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}

	config := xeroMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return xeroMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return xeroMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newXeroMCPHTTPHandler(xc XeroClient) http.Handler {
	server := newXeroMCPServer(xc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newXeroMCPServer(xc XeroClient) *mcp.Server {
	openWorld := true
	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	additiveWriteAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPtr(false), OpenWorldHint: &openWorld}
	idempotentWriteAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPtr(false), IdempotentHint: true, OpenWorldHint: &openWorld}
	destructiveWriteAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPtr(true), OpenWorldHint: &openWorld}
	server := mcp.NewServer(&mcp.Implementation{Name: "xero", Version: Version}, nil)

	mcp.AddTool(server, xeroTool("xero_tenants_list", xeroTenantsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *xeroEmptyToolInput) (*mcp.CallToolResult, XeroTenantConnectionsList, error) {
		out, err := xc.ListTenantsContext(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, xeroTool("xero_api_get", xeroAPIGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *xeroAPIReadToolInput) (*mcp.CallToolResult, XeroJSONResult, error) {
		out, err := xc.GetAPIEndpointContext(ctx, XeroAPIRequest{Family: input.Family, TenantID: input.TenantID, Endpoint: input.Endpoint, Query: input.Query})
		return nil, out, err
	})
	mcp.AddTool(server, xeroTool("xero_api_post", xeroAPIPostDoc, additiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *xeroAPIWriteToolInput) (*mcp.CallToolResult, XeroJSONResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := xc.PostAPIEndpointContext(ctx, XeroAPIRequest{Family: input.Family, TenantID: input.TenantID, Endpoint: input.Endpoint, Query: input.Query, Body: body})
		return nil, out, err
	})
	mcp.AddTool(server, xeroTool("xero_api_put", xeroAPIPutDoc, idempotentWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *xeroAPIWriteToolInput) (*mcp.CallToolResult, XeroJSONResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := xc.PutAPIEndpointContext(ctx, XeroAPIRequest{Family: input.Family, TenantID: input.TenantID, Endpoint: input.Endpoint, Query: input.Query, Body: body})
		return nil, out, err
	})
	mcp.AddTool(server, xeroTool("xero_api_delete", xeroAPIDeleteDoc, destructiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *xeroAPIReadToolInput) (*mcp.CallToolResult, XeroJSONResult, error) {
		out, err := xc.DeleteAPIEndpointContext(ctx, XeroAPIRequest{Family: input.Family, TenantID: input.TenantID, Endpoint: input.Endpoint, Query: input.Query})
		return nil, out, err
	})

	return server
}

func xeroTool(name string, description string, annotations *mcp.ToolAnnotations) *mcp.Tool {
	return &mcp.Tool{Name: name, Description: description, Annotations: annotations}
}

func boolPtr(value bool) *bool {
	return &value
}

func marshalMCPRequestBody(request map[string]any) (json.RawMessage, error) {
	if len(request) == 0 {
		return nil, fmt.Errorf("request body must not be empty")
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (cli CLI) printMCPHelp() {
	fmt.Fprintln(cli.stdout, `xero mcp

Run Xero as a local MCP server over Streamable HTTP.

Commands:
  xero mcp serve
  xero mcp serve --addr <addr>
  xero mcp serve --endpoint <path>`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintln(cli.stdout, `xero mcp serve

Run the Xero MCP server over Streamable HTTP.

Options:
  --addr <addr>       Listen address. Default: 127.0.0.1:7355
  --endpoint <path>   MCP endpoint path. Default: /mcp

Tools:
  xero_tenants_list
  xero_api_get
  xero_api_post
  xero_api_put
  xero_api_delete`)
}
