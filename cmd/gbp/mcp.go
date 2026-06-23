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

type gbpMCPConfig struct {
	Addr     string
	Endpoint string
}

type gbpEmptyToolInput struct{}

type gbpAccountToolInput struct {
	Account string `json:"account" jsonschema:"Google Business Profile account resource name, for example accounts/123."`
}

type gbpLocationToolInput struct {
	Location string `json:"location" jsonschema:"Google Business Profile location resource name, for example locations/456."`
	ReadMask string `json:"readMask" jsonschema:"Comma-separated read mask fields, for example name,title,storeCode,websiteUri."`
}

type gbpAccountReadMaskToolInput struct {
	Account  string `json:"account" jsonschema:"Google Business Profile account resource name, for example accounts/123."`
	ReadMask string `json:"readMask" jsonschema:"Comma-separated read mask fields, for example name,title,storeCode,websiteUri."`
}

type gbpLocationCreateToolInput struct {
	Account      string         `json:"account" jsonschema:"Google Business Profile account resource name, for example accounts/123."`
	Request      map[string]any `json:"request" jsonschema:"Google Business Profile Location request body using Google's documented JSON shape."`
	RequestID    string         `json:"requestId,omitempty" jsonschema:"Optional Google requestId query parameter."`
	ValidateOnly string         `json:"validateOnly,omitempty" jsonschema:"Optional Google validateOnly query parameter, usually true or false."`
}

type gbpLocationPatchToolInput struct {
	Location      string         `json:"location" jsonschema:"Google Business Profile location resource name, for example locations/456."`
	UpdateMask    string         `json:"updateMask" jsonschema:"Google updateMask query parameter."`
	AttributeMask string         `json:"attributeMask,omitempty" jsonschema:"Optional Google attributeMask query parameter."`
	ValidateOnly  string         `json:"validateOnly,omitempty" jsonschema:"Optional Google validateOnly query parameter, usually true or false."`
	Request       map[string]any `json:"request" jsonschema:"Google Business Profile Location request body using Google's documented JSON shape."`
}

type gbpLocationDeleteToolInput struct {
	Location string `json:"location" jsonschema:"Google Business Profile location resource name, for example locations/456."`
}

type gbpAccountLocationToolInput struct {
	Account   string `json:"account" jsonschema:"Google Business Profile account resource name, for example accounts/123."`
	Location  string `json:"location" jsonschema:"Google Business Profile location resource name, for example locations/456."`
	PageSize  string `json:"pageSize,omitempty" jsonschema:"Optional Google pageSize query parameter."`
	PageToken string `json:"pageToken,omitempty" jsonschema:"Optional Google pageToken query parameter."`
	OrderBy   string `json:"orderBy,omitempty" jsonschema:"Optional Google orderBy query parameter for reviews."`
}

type gbpAccountLocationRequestToolInput struct {
	Account  string         `json:"account" jsonschema:"Google Business Profile account resource name, for example accounts/123."`
	Location string         `json:"location" jsonschema:"Google Business Profile location resource name, for example locations/456."`
	Request  map[string]any `json:"request" jsonschema:"Google Business Profile API request body using the documented JSON shape for this operation."`
}

type gbpReviewToolInput struct {
	Account  string `json:"account" jsonschema:"Google Business Profile account resource name, for example accounts/123."`
	Location string `json:"location" jsonschema:"Google Business Profile location resource name, for example locations/456."`
	Review   string `json:"review" jsonschema:"Google Business Profile review ID."`
}

type gbpReviewReplyToolInput struct {
	Account  string         `json:"account" jsonschema:"Google Business Profile account resource name, for example accounts/123."`
	Location string         `json:"location" jsonschema:"Google Business Profile location resource name, for example locations/456."`
	Review   string         `json:"review" jsonschema:"Google Business Profile review ID."`
	Request  map[string]any `json:"request" jsonschema:"Google Business Profile ReviewReply request body using Google's documented JSON shape."`
}

type gbpNamedResourceToolInput struct {
	Name string `json:"name" jsonschema:"Full Google Business Profile resource name for this operation."`
}

type gbpNamedPatchToolInput struct {
	Name       string         `json:"name" jsonschema:"Full Google Business Profile resource name for this operation."`
	UpdateMask string         `json:"updateMask" jsonschema:"Google updateMask query parameter."`
	Request    map[string]any `json:"request" jsonschema:"Google Business Profile request body using the documented JSON shape for this operation."`
}

type gbpPerformanceToolInput struct {
	Location string         `json:"location" jsonschema:"Google Business Profile location resource name, for example locations/456."`
	Request  map[string]any `json:"request" jsonschema:"Google Business Profile Performance API query parameters using Google's documented nested JSON shape."`
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
	config, err := parseGBPMCPServeArgs(args)
	if err != nil {
		return err
	}
	gc, err := cli.gbpClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newGBPMCPHTTPHandler(gc))

	fmt.Fprintf(cli.stdout, "Google Business Profile MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseGBPMCPServeArgs(args []string) (gbpMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return gbpMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return gbpMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}

	config := gbpMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return gbpMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return gbpMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newGBPMCPHTTPHandler(gc GBPClient) http.Handler {
	server := newGBPMCPServer(gc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newGBPMCPServer(gc GBPClient) *mcp.Server {
	openWorld := true
	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	additiveWriteAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPtr(false), OpenWorldHint: &openWorld}
	idempotentWriteAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPtr(false), IdempotentHint: true, OpenWorldHint: &openWorld}
	destructiveWriteAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPtr(true), OpenWorldHint: &openWorld}
	server := mcp.NewServer(&mcp.Implementation{Name: "gbp", Version: Version}, nil)

	mcp.AddTool(server, gbpTool("gbp_auth_test", gbpAuthTestDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *gbpEmptyToolInput) (*mcp.CallToolResult, GBPAccountsList, error) {
		out, err := gc.AuthTestContext(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_accounts_list", gbpAccountsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *gbpEmptyToolInput) (*mcp.CallToolResult, GBPAccountsList, error) {
		out, err := gc.ListAccountsContext(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_account_get", gbpAccountGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpAccountToolInput) (*mcp.CallToolResult, GBPAccount, error) {
		out, err := gc.GetAccountContext(ctx, input.Account)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_locations_list", gbpLocationsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpAccountReadMaskToolInput) (*mcp.CallToolResult, GBPLocationsList, error) {
		out, err := gc.ListLocationsContext(ctx, input.Account, input.ReadMask)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_location_get", gbpLocationGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpLocationToolInput) (*mcp.CallToolResult, GBPLocation, error) {
		out, err := gc.GetLocationContext(ctx, input.Location, input.ReadMask)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_location_create", gbpLocationCreateDoc, additiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpLocationCreateToolInput) (*mcp.CallToolResult, GBPLocation, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.CreateLocationContext(ctx, input.Account, body, locationWriteOptions{RequestID: input.RequestID, ValidateOnly: input.ValidateOnly})
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_location_patch", gbpLocationPatchDoc, idempotentWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpLocationPatchToolInput) (*mcp.CallToolResult, GBPLocation, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.PatchLocationContext(ctx, input.Location, body, locationPatchOptions{UpdateMask: input.UpdateMask, AttributeMask: input.AttributeMask, ValidateOnly: input.ValidateOnly})
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_location_delete", gbpLocationDeleteDoc, destructiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpLocationDeleteToolInput) (*mcp.CallToolResult, GBPWriteResult, error) {
		out, err := gc.DeleteLocationContext(ctx, input.Location)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_reviews_list", gbpReviewsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpAccountLocationToolInput) (*mcp.CallToolResult, GBPReviewsList, error) {
		out, err := gc.ListReviewsContext(ctx, input.Account, input.Location, pageOptions{PageSize: input.PageSize, PageToken: input.PageToken, OrderBy: input.OrderBy})
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_review_get", gbpReviewGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpReviewToolInput) (*mcp.CallToolResult, GBPReview, error) {
		out, err := gc.GetReviewContext(ctx, input.Account, input.Location, input.Review)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_review_update_reply", gbpReviewUpdateReplyDoc, idempotentWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpReviewReplyToolInput) (*mcp.CallToolResult, GBPWriteResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.UpdateReviewReplyContext(ctx, input.Account, input.Location, input.Review, body)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_review_delete_reply", gbpReviewDeleteReplyDoc, destructiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpReviewToolInput) (*mcp.CallToolResult, GBPWriteResult, error) {
		out, err := gc.DeleteReviewReplyContext(ctx, input.Account, input.Location, input.Review)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_media_list", gbpMediaListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpAccountLocationToolInput) (*mcp.CallToolResult, GBPMediaList, error) {
		out, err := gc.ListMediaContext(ctx, input.Account, input.Location, pageOptions{PageSize: input.PageSize, PageToken: input.PageToken})
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_media_create", gbpMediaCreateDoc, additiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpAccountLocationRequestToolInput) (*mcp.CallToolResult, GBPMediaItem, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.CreateMediaContext(ctx, input.Account, input.Location, body)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_media_get", gbpMediaGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpNamedResourceToolInput) (*mcp.CallToolResult, GBPMediaItem, error) {
		out, err := gc.GetMediaContext(ctx, input.Name)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_media_patch", gbpMediaPatchDoc, idempotentWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpNamedPatchToolInput) (*mcp.CallToolResult, GBPMediaItem, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.PatchMediaContext(ctx, input.Name, input.UpdateMask, body)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_media_delete", gbpMediaDeleteDoc, destructiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpNamedResourceToolInput) (*mcp.CallToolResult, GBPWriteResult, error) {
		out, err := gc.DeleteMediaContext(ctx, input.Name)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_media_start_upload", gbpMediaStartUploadDoc, additiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpAccountLocationToolInput) (*mcp.CallToolResult, GBPWriteResult, error) {
		out, err := gc.StartMediaUploadContext(ctx, input.Account, input.Location)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_local_posts_list", gbpLocalPostsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpAccountLocationToolInput) (*mcp.CallToolResult, GBPLocalPostsList, error) {
		out, err := gc.ListLocalPostsContext(ctx, input.Account, input.Location, pageOptions{PageSize: input.PageSize, PageToken: input.PageToken})
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_local_post_create", gbpLocalPostCreateDoc, additiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpAccountLocationRequestToolInput) (*mcp.CallToolResult, GBPLocalPost, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.CreateLocalPostContext(ctx, input.Account, input.Location, body)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_local_post_get", gbpLocalPostGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpNamedResourceToolInput) (*mcp.CallToolResult, GBPLocalPost, error) {
		out, err := gc.GetLocalPostContext(ctx, input.Name)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_local_post_patch", gbpLocalPostPatchDoc, idempotentWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpNamedPatchToolInput) (*mcp.CallToolResult, GBPLocalPost, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.PatchLocalPostContext(ctx, input.Name, input.UpdateMask, body)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_local_post_delete", gbpLocalPostDeleteDoc, destructiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpNamedResourceToolInput) (*mcp.CallToolResult, GBPWriteResult, error) {
		out, err := gc.DeleteLocalPostContext(ctx, input.Name)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_local_post_report_insights", gbpLocalPostReportInsightsDoc, additiveWriteAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpAccountLocationRequestToolInput) (*mcp.CallToolResult, GBPWriteResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.ReportLocalPostInsightsContext(ctx, input.Account, input.Location, body)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_performance_daily_metrics", gbpPerformanceDailyMetricsDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpPerformanceToolInput) (*mcp.CallToolResult, GBPPerformanceResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.GetDailyMetricsContext(ctx, input.Location, body)
		return nil, out, err
	})
	mcp.AddTool(server, gbpTool("gbp_performance_search_keywords", gbpPerformanceSearchKeywordsDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *gbpPerformanceToolInput) (*mcp.CallToolResult, GBPPerformanceResult, error) {
		body, err := marshalMCPRequestBody(input.Request)
		if err != nil {
			return nil, nil, err
		}
		out, err := gc.ListSearchKeywordsContext(ctx, input.Location, body)
		return nil, out, err
	})

	return server
}

func gbpTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
	return &mcp.Tool{Name: name, Title: doc.Command, Description: doc.Description, Annotations: annotations}
}

func boolPtr(value bool) *bool {
	return &value
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
	fmt.Fprint(cli.stdout, `gbp mcp

Run Google Business Profile as a local MCP server.

Usage:
  gbp mcp help
  gbp mcp serve
  gbp mcp serve --help

Commands:
  serve    Serve Google Business Profile MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `gbp mcp serve

Serve Google Business Profile tools over MCP Streamable HTTP.

Usage:
  gbp mcp serve
  gbp mcp serve --addr <addr>
  gbp mcp serve --endpoint <path>
  gbp mcp serve --addr <addr> --endpoint <path>
  gbp mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.

Tools:
  gbp_auth_test                         %s
  gbp_accounts_list                     %s
  gbp_account_get                       %s
  gbp_locations_list                    %s
  gbp_location_get                      %s
  gbp_location_create                   %s
  gbp_location_patch                    %s
  gbp_location_delete                   %s
  gbp_reviews_list                      %s
  gbp_review_get                        %s
  gbp_review_update_reply               %s
  gbp_review_delete_reply               %s
  gbp_media_list                        %s
  gbp_media_create                      %s
  gbp_media_get                         %s
  gbp_media_patch                       %s
  gbp_media_delete                      %s
  gbp_media_start_upload                %s
  gbp_local_posts_list                  %s
  gbp_local_post_create                 %s
  gbp_local_post_get                    %s
  gbp_local_post_patch                  %s
  gbp_local_post_delete                 %s
  gbp_local_post_report_insights        %s
  gbp_performance_daily_metrics         %s
  gbp_performance_search_keywords       %s
`, defaultMCPAddr, defaultMCPEndpoint, gbpAuthTestDoc.Summary, gbpAccountsListDoc.Summary, gbpAccountGetDoc.Summary, gbpLocationsListDoc.Summary, gbpLocationGetDoc.Summary, gbpLocationCreateDoc.Summary, gbpLocationPatchDoc.Summary, gbpLocationDeleteDoc.Summary, gbpReviewsListDoc.Summary, gbpReviewGetDoc.Summary, gbpReviewUpdateReplyDoc.Summary, gbpReviewDeleteReplyDoc.Summary, gbpMediaListDoc.Summary, gbpMediaCreateDoc.Summary, gbpMediaGetDoc.Summary, gbpMediaPatchDoc.Summary, gbpMediaDeleteDoc.Summary, gbpMediaStartUploadDoc.Summary, gbpLocalPostsListDoc.Summary, gbpLocalPostCreateDoc.Summary, gbpLocalPostGetDoc.Summary, gbpLocalPostPatchDoc.Summary, gbpLocalPostDeleteDoc.Summary, gbpLocalPostReportInsightsDoc.Summary, gbpPerformanceDailyMetricsDoc.Summary, gbpPerformanceSearchKeywordsDoc.Summary)
}
