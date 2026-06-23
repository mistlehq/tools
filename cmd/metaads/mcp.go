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
	defaultMCPAddr     = "127.0.0.1:7350"
	defaultMCPEndpoint = "/mcp"
)

type metaAdsMCPConfig struct {
	Addr     string
	Endpoint string
}

type metaAdsEmptyToolInput struct{}

type metaAdsGraphRequestToolInput struct {
	Method string         `json:"method,omitempty" jsonschema:"HTTP method: GET, POST, or DELETE. Defaults to GET."`
	Path   string         `json:"path" jsonschema:"Meta Graph API path, starting with '/'."`
	Params map[string]any `json:"params,omitempty" jsonschema:"Optional query parameters object."`
	Body   map[string]any `json:"body,omitempty" jsonschema:"Optional JSON request body for POST requests."`
}

type metaAdsEdgeToolInput struct {
	AccountID string         `json:"account_id,omitempty" jsonschema:"Meta ad account ID, usually act_<account-id>."`
	Fields    string         `json:"fields,omitempty" jsonschema:"Comma-separated fields parameter."`
	Limit     string         `json:"limit,omitempty" jsonschema:"Optional Meta Graph API limit parameter."`
	After     string         `json:"after,omitempty" jsonschema:"Optional Meta Graph API paging cursor."`
	Params    map[string]any `json:"params,omitempty" jsonschema:"Additional Graph API query parameters."`
}

type metaAdsGetToolInput struct {
	ID     string         `json:"id" jsonschema:"Meta object ID."`
	Fields string         `json:"fields,omitempty" jsonschema:"Comma-separated fields parameter."`
	Params map[string]any `json:"params,omitempty" jsonschema:"Additional Graph API query parameters."`
}

type metaAdsAccountBodyToolInput struct {
	AccountID string         `json:"account_id" jsonschema:"Meta ad account ID, usually act_<account-id>."`
	Body      map[string]any `json:"body" jsonschema:"Request body using the documented Meta Graph API shape."`
}

type metaAdsObjectBodyToolInput struct {
	ID   string         `json:"id" jsonschema:"Meta object ID."`
	Body map[string]any `json:"body" jsonschema:"Request body using the documented Meta Graph API shape."`
}

type metaAdsIDToolInput struct {
	ID string `json:"id" jsonschema:"Meta object ID."`
}

type metaAdsInsightsToolInput struct {
	ID        string         `json:"id" jsonschema:"Ad account or object ID whose insights edge should be queried."`
	Fields    string         `json:"fields,omitempty" jsonschema:"Comma-separated fields parameter."`
	Level     string         `json:"level,omitempty" jsonschema:"Optional insights level, such as campaign, adset, or ad."`
	TimeRange map[string]any `json:"time_range,omitempty" jsonschema:"Optional Meta time_range object."`
	Params    map[string]any `json:"params,omitempty" jsonschema:"Additional Graph API query parameters."`
}

type metaAdsTargetingSearchToolInput struct {
	Type   string         `json:"type" jsonschema:"Meta targeting search type, such as adinterest or geo_location."`
	Query  string         `json:"query,omitempty" jsonschema:"Optional q parameter."`
	Params map[string]any `json:"params,omitempty" jsonschema:"Additional Graph API query parameters."`
}

type metaAdsRawGraphResponse map[string]any

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
	config, err := parseMetaAdsMCPServeArgs(args)
	if err != nil {
		return err
	}
	mc, err := cli.metaAdsClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newMetaAdsMCPHTTPHandler(mc))

	fmt.Fprintf(cli.stdout, "Meta Ads MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseMetaAdsMCPServeArgs(args []string) (metaAdsMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return metaAdsMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return metaAdsMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}
	config := metaAdsMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return metaAdsMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return metaAdsMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newMetaAdsMCPHTTPHandler(mc MetaAdsClient) http.Handler {
	server := newMetaAdsMCPServer(mc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newMetaAdsMCPServer(mc MetaAdsClient) *mcp.Server {
	openWorld := true
	destructive := true
	notDestructive := false

	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	mutatingAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &notDestructive}
	destructiveAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &destructive}
	rawAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld}

	server := mcp.NewServer(&mcp.Implementation{Name: "metaads", Version: Version}, nil)

	mcp.AddTool(server, metaAdsTool("metaads_graph_request", metaAdsGraphRequestDoc, rawAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsGraphRequestToolInput) (*mcp.CallToolResult, metaAdsRawGraphResponse, error) {
		body, err := mc.RequestContext(ctx, MetaAdsRequest{Method: input.Method, Path: input.Path, Params: input.Params, Body: input.Body})
		if err != nil {
			return nil, nil, err
		}
		var out metaAdsRawGraphResponse
		if err := json.Unmarshal(body, &out); err != nil {
			return nil, nil, err
		}
		return nil, out, nil
	})
	mcp.AddTool(server, metaAdsTool("metaads_auth_test", metaAdsAuthTestDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *metaAdsEmptyToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := mc.AuthTest(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_ad_accounts_list", metaAdsAdAccountsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsEdgeToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := mc.ListAdAccounts(ctx, edgeInputFromTool(*input))
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_ad_account_get", metaAdsAdAccountGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsGetToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := mc.GetAdAccount(ctx, getInputFromTool(*input))
		return nil, out, err
	})
	addCRUDMCPTools(server, mc, "campaign", metaAdsCampaignsSearchDoc, metaAdsCampaignGetDoc, metaAdsCampaignCreateDoc, metaAdsCampaignUpdateDoc, metaAdsCampaignDeleteDoc, readOnlyAnnotations, mutatingAnnotations, destructiveAnnotations, campaignMCPHandlers{})
	addCRUDMCPTools(server, mc, "adset", metaAdsAdSetsSearchDoc, metaAdsAdSetGetDoc, metaAdsAdSetCreateDoc, metaAdsAdSetUpdateDoc, metaAdsAdSetDeleteDoc, readOnlyAnnotations, mutatingAnnotations, destructiveAnnotations, adSetMCPHandlers{})
	addCRUDMCPTools(server, mc, "ad", metaAdsAdsSearchDoc, metaAdsAdGetDoc, metaAdsAdCreateDoc, metaAdsAdUpdateDoc, metaAdsAdDeleteDoc, readOnlyAnnotations, mutatingAnnotations, destructiveAnnotations, adMCPHandlers{})
	mcp.AddTool(server, metaAdsTool("metaads_creatives_search", metaAdsCreativesSearchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsEdgeToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := mc.SearchCreatives(ctx, edgeInputFromTool(*input))
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_creative_get", metaAdsCreativeGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsGetToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := mc.GetCreative(ctx, getInputFromTool(*input))
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_creative_create", metaAdsCreativeCreateDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsAccountBodyToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := mc.CreateCreative(ctx, MetaAdsObjectInput{ID: input.AccountID, Body: input.Body})
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_insights_get", metaAdsInsightsGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsInsightsToolInput) (*mcp.CallToolResult, map[string]any, error) {
		timeRange := ""
		if input.TimeRange != nil {
			body, err := json.Marshal(input.TimeRange)
			if err != nil {
				return nil, nil, err
			}
			timeRange = string(body)
		}
		out, err := mc.GetInsights(ctx, MetaAdsInsightsInput{ID: input.ID, Fields: input.Fields, Level: input.Level, TimeRange: timeRange, Params: input.Params})
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_targeting_search", metaAdsTargetingSearchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsTargetingSearchToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := mc.SearchTargeting(ctx, MetaAdsTargetingSearchInput{Type: input.Type, Query: input.Query, Params: input.Params})
		return nil, out, err
	})

	return server
}

type crudMCPHandlers interface {
	search(context.Context, MetaAdsClient, MetaAdsEdgeInput) (map[string]any, error)
	get(context.Context, MetaAdsClient, MetaAdsGetInput) (map[string]any, error)
	create(context.Context, MetaAdsClient, MetaAdsObjectInput) (map[string]any, error)
	update(context.Context, MetaAdsClient, MetaAdsObjectInput) (map[string]any, error)
	delete(context.Context, MetaAdsClient, string) (map[string]any, error)
}

type campaignMCPHandlers struct{}
type adSetMCPHandlers struct{}
type adMCPHandlers struct{}

func (campaignMCPHandlers) search(ctx context.Context, mc MetaAdsClient, input MetaAdsEdgeInput) (map[string]any, error) {
	return mc.SearchCampaigns(ctx, input)
}
func (campaignMCPHandlers) get(ctx context.Context, mc MetaAdsClient, input MetaAdsGetInput) (map[string]any, error) {
	return mc.GetCampaign(ctx, input)
}
func (campaignMCPHandlers) create(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.CreateCampaign(ctx, input)
}
func (campaignMCPHandlers) update(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.UpdateCampaign(ctx, input)
}
func (campaignMCPHandlers) delete(ctx context.Context, mc MetaAdsClient, id string) (map[string]any, error) {
	return mc.DeleteCampaign(ctx, id)
}

func (adSetMCPHandlers) search(ctx context.Context, mc MetaAdsClient, input MetaAdsEdgeInput) (map[string]any, error) {
	return mc.SearchAdSets(ctx, input)
}
func (adSetMCPHandlers) get(ctx context.Context, mc MetaAdsClient, input MetaAdsGetInput) (map[string]any, error) {
	return mc.GetAdSet(ctx, input)
}
func (adSetMCPHandlers) create(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.CreateAdSet(ctx, input)
}
func (adSetMCPHandlers) update(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.UpdateAdSet(ctx, input)
}
func (adSetMCPHandlers) delete(ctx context.Context, mc MetaAdsClient, id string) (map[string]any, error) {
	return mc.DeleteAdSet(ctx, id)
}

func (adMCPHandlers) search(ctx context.Context, mc MetaAdsClient, input MetaAdsEdgeInput) (map[string]any, error) {
	return mc.SearchAds(ctx, input)
}
func (adMCPHandlers) get(ctx context.Context, mc MetaAdsClient, input MetaAdsGetInput) (map[string]any, error) {
	return mc.GetAd(ctx, input)
}
func (adMCPHandlers) create(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.CreateAd(ctx, input)
}
func (adMCPHandlers) update(ctx context.Context, mc MetaAdsClient, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.UpdateAd(ctx, input)
}
func (adMCPHandlers) delete(ctx context.Context, mc MetaAdsClient, id string) (map[string]any, error) {
	return mc.DeleteAd(ctx, id)
}

func addCRUDMCPTools(server *mcp.Server, mc MetaAdsClient, singular string, searchDoc commandDoc, getDoc commandDoc, createDoc commandDoc, updateDoc commandDoc, deleteDoc commandDoc, readOnlyAnnotations *mcp.ToolAnnotations, mutatingAnnotations *mcp.ToolAnnotations, destructiveAnnotations *mcp.ToolAnnotations, handlers crudMCPHandlers) {
	plural := singular + "s"
	if singular == "adset" {
		plural = "adsets"
	}
	mcp.AddTool(server, metaAdsTool("metaads_"+plural+"_search", searchDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsEdgeToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := handlers.search(ctx, mc, edgeInputFromTool(*input))
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_"+singular+"_get", getDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsGetToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := handlers.get(ctx, mc, getInputFromTool(*input))
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_"+singular+"_create", createDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsAccountBodyToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := handlers.create(ctx, mc, MetaAdsObjectInput{ID: input.AccountID, Body: input.Body})
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_"+singular+"_update", updateDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsObjectBodyToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := handlers.update(ctx, mc, MetaAdsObjectInput{ID: input.ID, Body: input.Body})
		return nil, out, err
	})
	mcp.AddTool(server, metaAdsTool("metaads_"+singular+"_delete", deleteDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *metaAdsIDToolInput) (*mcp.CallToolResult, map[string]any, error) {
		out, err := handlers.delete(ctx, mc, input.ID)
		return nil, out, err
	})
}

func edgeInputFromTool(input metaAdsEdgeToolInput) MetaAdsEdgeInput {
	return MetaAdsEdgeInput{AccountID: input.AccountID, Fields: input.Fields, Limit: input.Limit, After: input.After, Params: input.Params}
}

func getInputFromTool(input metaAdsGetToolInput) MetaAdsGetInput {
	return MetaAdsGetInput{ID: input.ID, Fields: input.Fields, Params: input.Params}
}

func metaAdsTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
	return &mcp.Tool{Name: name, Title: doc.Command, Description: doc.Description, Annotations: annotations}
}

func (cli CLI) printMCPHelp() {
	fmt.Fprint(cli.stdout, `metaads mcp

Run Meta Ads as a local MCP server.

Usage:
  metaads mcp help
  metaads mcp serve
  metaads mcp serve --help

Commands:
  serve    Serve Meta Ads MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `metaads mcp serve

Serve Meta Ads tools over MCP Streamable HTTP.

Usage:
  metaads mcp serve
  metaads mcp serve --addr <addr>
  metaads mcp serve --endpoint <path>

Options:
  --addr <addr>        Address to listen on. Defaults to %s.
  --endpoint <path>    HTTP endpoint path. Defaults to %s.

Tools:
  metaads_graph_request
  metaads_auth_test
  metaads_ad_accounts_list
  metaads_ad_account_get
  metaads_campaigns_search
  metaads_campaign_get
  metaads_campaign_create
  metaads_campaign_update
  metaads_campaign_delete
  metaads_adsets_search
  metaads_adset_get
  metaads_adset_create
  metaads_adset_update
  metaads_adset_delete
  metaads_ads_search
  metaads_ad_get
  metaads_ad_create
  metaads_ad_update
  metaads_ad_delete
  metaads_creatives_search
  metaads_creative_get
  metaads_creative_create
  metaads_insights_get
  metaads_targeting_search
`, defaultMCPAddr, defaultMCPEndpoint)
}
