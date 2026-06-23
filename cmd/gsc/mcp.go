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
	defaultMCPAddr     = "127.0.0.1:7349"
	defaultMCPEndpoint = "/mcp"
)

type gscMCPConfig struct {
	Addr     string
	Endpoint string
}

type gscEmptyToolInput struct{}

type gscSiteToolInput struct {
	SiteURL string `json:"siteUrl" jsonschema:"Google Search Console site URL, for example https://example.com/ or sc-domain:example.com."`
}

type gscSiteRequestToolInput struct {
	SiteURL string         `json:"siteUrl" jsonschema:"Google Search Console site URL, for example https://example.com/ or sc-domain:example.com."`
	Request map[string]any `json:"request" jsonschema:"Google Search Console API request body using the documented JSON shape for this operation."`
}

type gscSitemapToolInput struct {
	SiteURL  string `json:"siteUrl" jsonschema:"Google Search Console site URL, for example https://example.com/ or sc-domain:example.com."`
	FeedPath string `json:"feedPath" jsonschema:"Sitemap URL/feed path, for example https://example.com/sitemap.xml."`
}

type gscRequestToolInput struct {
	Request map[string]any `json:"request" jsonschema:"Google Search Console API request body using the documented JSON shape for this operation."`
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
	config, err := parseGSCMCPServeArgs(args)
	if err != nil {
		return err
	}
	gc, err := cli.gscClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newGSCMCPHTTPHandler(gc))

	fmt.Fprintf(cli.stdout, "Google Search Console MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseGSCMCPServeArgs(args []string) (gscMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return gscMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return gscMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}

	config := gscMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return gscMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return gscMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newGSCMCPHTTPHandler(gc GSCClient) http.Handler {
	server := newGSCMCPServer(gc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newGSCMCPServer(gc GSCClient) *mcp.Server {
	openWorld := true
	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	server := mcp.NewServer(&mcp.Implementation{Name: "gsc", Version: Version}, nil)

	mcp.AddTool(server, gscTool("gsc_auth_test", gscAuthTestDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gscSiteToolInput) (*mcp.CallToolResult, GSCSite, error) {
		out, err := gc.AuthTestContext(ctx, input.SiteURL)
		return nil, out, err
	})
	mcp.AddTool(server, gscTool("gsc_sites_list", gscSitesListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *gscEmptyToolInput) (*mcp.CallToolResult, GSCSitesList, error) {
		out, err := gc.ListSitesContext(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, gscTool("gsc_site_get", gscSiteGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gscSiteToolInput) (*mcp.CallToolResult, GSCSite, error) {
		out, err := gc.GetSiteContext(ctx, input.SiteURL)
		return nil, out, err
	})
	mcp.AddTool(server, gscTool("gsc_searchanalytics_query", gscSearchAnalyticsQueryDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gscSiteRequestToolInput) (*mcp.CallToolResult, GSCSearchAnalyticsResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.QuerySearchAnalyticsContext(ctx, input.SiteURL, body)
		return nil, out, err
	})
	mcp.AddTool(server, gscTool("gsc_sitemaps_list", gscSitemapsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gscSiteToolInput) (*mcp.CallToolResult, GSCSitemapsList, error) {
		out, err := gc.ListSitemapsContext(ctx, input.SiteURL)
		return nil, out, err
	})
	mcp.AddTool(server, gscTool("gsc_sitemap_get", gscSitemapGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gscSitemapToolInput) (*mcp.CallToolResult, GSCSitemap, error) {
		out, err := gc.GetSitemapContext(ctx, input.SiteURL, input.FeedPath)
		return nil, out, err
	})
	mcp.AddTool(server, gscTool("gsc_url_inspection_inspect", gscURLInspectionInspectDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gscRequestToolInput) (*mcp.CallToolResult, GSCURLInspectionResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.InspectURLContext(ctx, body)
		return nil, out, err
	})

	return server
}

func gscTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
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
	fmt.Fprint(cli.stdout, `gsc mcp

Run Google Search Console as a local MCP server.

Usage:
  gsc mcp help
  gsc mcp serve
  gsc mcp serve --help

Commands:
  serve    Serve Google Search Console MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `gsc mcp serve

Serve Google Search Console tools over MCP Streamable HTTP.

Usage:
  gsc mcp serve
  gsc mcp serve --addr <addr>
  gsc mcp serve --endpoint <path>
  gsc mcp serve --addr <addr> --endpoint <path>
  gsc mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.

Tools:
  gsc_auth_test                 %s
  gsc_sites_list                %s
  gsc_site_get                  %s
  gsc_searchanalytics_query     %s
  gsc_sitemaps_list             %s
  gsc_sitemap_get               %s
  gsc_url_inspection_inspect    %s
`, defaultMCPAddr, defaultMCPEndpoint, gscAuthTestDoc.Summary, gscSitesListDoc.Summary, gscSiteGetDoc.Summary, gscSearchAnalyticsQueryDoc.Summary, gscSitemapsListDoc.Summary, gscSitemapGetDoc.Summary, gscURLInspectionInspectDoc.Summary)
}
