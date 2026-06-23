package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type MetaAdsClient struct {
	graphBaseURL string
	client       *http.Client
}

type MetaAdsRequest struct {
	Method string         `json:"method"`
	Path   string         `json:"path"`
	Params map[string]any `json:"params,omitempty"`
	Body   map[string]any `json:"body,omitempty"`
}

type MetaAdsEdgeInput struct {
	AccountID string
	Fields    string
	Limit     string
	After     string
	Params    map[string]any
}

type MetaAdsGetInput struct {
	ID     string
	Fields string
	Params map[string]any
}

type MetaAdsObjectInput struct {
	ID   string
	Body map[string]any
}

type MetaAdsInsightsInput struct {
	ID        string
	Fields    string
	Level     string
	TimeRange string
	Params    map[string]any
}

type MetaAdsTargetingSearchInput struct {
	Type   string
	Query  string
	Params map[string]any
}

func NewMetaAdsClient(config Config) MetaAdsClient {
	return MetaAdsClient{
		graphBaseURL: config.GraphBaseURL,
		client:       http.DefaultClient,
	}
}

func (mc MetaAdsClient) Request(request MetaAdsRequest) ([]byte, error) {
	return mc.RequestContext(context.Background(), request)
}

func (mc MetaAdsClient) RequestContext(ctx context.Context, request MetaAdsRequest) ([]byte, error) {
	method := strings.ToUpper(strings.TrimSpace(request.Method))
	if method == "" {
		method = http.MethodGet
	}
	if method != http.MethodGet && method != http.MethodPost && method != http.MethodDelete {
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
	if strings.TrimSpace(request.Path) == "" {
		return nil, fmt.Errorf("path is required")
	}
	if !strings.HasPrefix(request.Path, "/") {
		return nil, fmt.Errorf("path must start with '/'")
	}

	requestURL, err := url.Parse(mc.graphBaseURL + request.Path)
	if err != nil {
		return nil, err
	}
	query := requestURL.Query()
	for key, value := range request.Params {
		setQueryValue(query, key, value)
	}
	requestURL.RawQuery = query.Encode()

	var body io.Reader
	if request.Body != nil {
		encodedBody, err := json.Marshal(request.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(encodedBody)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, method, requestURL.String(), body)
	if err != nil {
		return nil, err
	}
	if request.Body != nil {
		httpRequest.Header.Set("Content-Type", "application/json")
	}

	response, err := mc.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("meta graph api request failed with status %d: %s", response.StatusCode, string(responseBody))
	}
	return responseBody, nil
}

func setQueryValue(query url.Values, key string, value any) {
	switch typedValue := value.(type) {
	case nil:
		return
	case string:
		query.Set(key, typedValue)
	default:
		encoded, err := json.Marshal(typedValue)
		if err != nil {
			query.Set(key, fmt.Sprint(typedValue))
			return
		}
		query.Set(key, string(encoded))
	}
}

func (mc MetaAdsClient) AuthTest(ctx context.Context) (map[string]any, error) {
	return mc.getJSON(ctx, MetaAdsRequest{Method: http.MethodGet, Path: "/me", Params: map[string]any{"fields": "id,name"}})
}

func (mc MetaAdsClient) ListAdAccounts(ctx context.Context, input MetaAdsEdgeInput) (map[string]any, error) {
	params := edgeParams(input, "id,name,account_id,account_status,currency,timezone_name")
	return mc.getJSON(ctx, MetaAdsRequest{Method: http.MethodGet, Path: "/me/adaccounts", Params: params})
}

func (mc MetaAdsClient) GetAdAccount(ctx context.Context, input MetaAdsGetInput) (map[string]any, error) {
	return mc.getObject(ctx, input, "id,name,account_id,account_status,currency,timezone_name")
}

func (mc MetaAdsClient) SearchCampaigns(ctx context.Context, input MetaAdsEdgeInput) (map[string]any, error) {
	return mc.accountEdge(ctx, input, "campaigns", "id,name,status,effective_status,objective,created_time,updated_time")
}

func (mc MetaAdsClient) GetCampaign(ctx context.Context, input MetaAdsGetInput) (map[string]any, error) {
	return mc.getObject(ctx, input, "id,name,status,effective_status,objective,created_time,updated_time")
}

func (mc MetaAdsClient) CreateCampaign(ctx context.Context, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.createAccountEdge(ctx, input, "campaigns")
}

func (mc MetaAdsClient) UpdateCampaign(ctx context.Context, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.updateObject(ctx, input)
}

func (mc MetaAdsClient) DeleteCampaign(ctx context.Context, id string) (map[string]any, error) {
	return mc.deleteObject(ctx, id)
}

func (mc MetaAdsClient) SearchAdSets(ctx context.Context, input MetaAdsEdgeInput) (map[string]any, error) {
	return mc.accountEdge(ctx, input, "adsets", "id,name,status,effective_status,campaign_id,daily_budget,lifetime_budget,created_time,updated_time")
}

func (mc MetaAdsClient) GetAdSet(ctx context.Context, input MetaAdsGetInput) (map[string]any, error) {
	return mc.getObject(ctx, input, "id,name,status,effective_status,campaign_id,daily_budget,lifetime_budget,created_time,updated_time")
}

func (mc MetaAdsClient) CreateAdSet(ctx context.Context, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.createAccountEdge(ctx, input, "adsets")
}

func (mc MetaAdsClient) UpdateAdSet(ctx context.Context, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.updateObject(ctx, input)
}

func (mc MetaAdsClient) DeleteAdSet(ctx context.Context, id string) (map[string]any, error) {
	return mc.deleteObject(ctx, id)
}

func (mc MetaAdsClient) SearchAds(ctx context.Context, input MetaAdsEdgeInput) (map[string]any, error) {
	return mc.accountEdge(ctx, input, "ads", "id,name,status,effective_status,campaign_id,adset_id,creative,created_time,updated_time")
}

func (mc MetaAdsClient) GetAd(ctx context.Context, input MetaAdsGetInput) (map[string]any, error) {
	return mc.getObject(ctx, input, "id,name,status,effective_status,campaign_id,adset_id,creative,created_time,updated_time")
}

func (mc MetaAdsClient) CreateAd(ctx context.Context, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.createAccountEdge(ctx, input, "ads")
}

func (mc MetaAdsClient) UpdateAd(ctx context.Context, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.updateObject(ctx, input)
}

func (mc MetaAdsClient) DeleteAd(ctx context.Context, id string) (map[string]any, error) {
	return mc.deleteObject(ctx, id)
}

func (mc MetaAdsClient) SearchCreatives(ctx context.Context, input MetaAdsEdgeInput) (map[string]any, error) {
	return mc.accountEdge(ctx, input, "adcreatives", "id,name,status,object_story_id,thumbnail_url,effective_object_story_id")
}

func (mc MetaAdsClient) GetCreative(ctx context.Context, input MetaAdsGetInput) (map[string]any, error) {
	return mc.getObject(ctx, input, "id,name,status,object_story_id,thumbnail_url,effective_object_story_id")
}

func (mc MetaAdsClient) CreateCreative(ctx context.Context, input MetaAdsObjectInput) (map[string]any, error) {
	return mc.createAccountEdge(ctx, input, "adcreatives")
}

func (mc MetaAdsClient) GetInsights(ctx context.Context, input MetaAdsInsightsInput) (map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, fmt.Errorf("id is required")
	}
	params := cloneParams(input.Params)
	if strings.TrimSpace(input.Fields) != "" {
		params["fields"] = input.Fields
	}
	if strings.TrimSpace(input.Level) != "" {
		params["level"] = input.Level
	}
	if strings.TrimSpace(input.TimeRange) != "" {
		var timeRange map[string]any
		if err := json.Unmarshal([]byte(input.TimeRange), &timeRange); err != nil {
			return nil, fmt.Errorf("time-range must be a JSON object: %w", err)
		}
		params["time_range"] = timeRange
	}
	return mc.getJSON(ctx, MetaAdsRequest{Method: http.MethodGet, Path: "/" + input.ID + "/insights", Params: params})
}

func (mc MetaAdsClient) SearchTargeting(ctx context.Context, input MetaAdsTargetingSearchInput) (map[string]any, error) {
	params := cloneParams(input.Params)
	if strings.TrimSpace(input.Type) == "" {
		return nil, fmt.Errorf("type is required")
	}
	params["type"] = input.Type
	if strings.TrimSpace(input.Query) != "" {
		params["q"] = input.Query
	}
	return mc.getJSON(ctx, MetaAdsRequest{Method: http.MethodGet, Path: "/search", Params: params})
}

func (mc MetaAdsClient) accountEdge(ctx context.Context, input MetaAdsEdgeInput, edge string, defaultFields string) (map[string]any, error) {
	if strings.TrimSpace(input.AccountID) == "" {
		return nil, fmt.Errorf("account id is required")
	}
	return mc.getJSON(ctx, MetaAdsRequest{
		Method: http.MethodGet,
		Path:   "/" + input.AccountID + "/" + edge,
		Params: edgeParams(input, defaultFields),
	})
}

func (mc MetaAdsClient) createAccountEdge(ctx context.Context, input MetaAdsObjectInput, edge string) (map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, fmt.Errorf("account id is required")
	}
	if input.Body == nil {
		return nil, fmt.Errorf("body is required")
	}
	return mc.getJSON(ctx, MetaAdsRequest{Method: http.MethodPost, Path: "/" + input.ID + "/" + edge, Body: input.Body})
}

func (mc MetaAdsClient) getObject(ctx context.Context, input MetaAdsGetInput, defaultFields string) (map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, fmt.Errorf("id is required")
	}
	params := cloneParams(input.Params)
	if strings.TrimSpace(input.Fields) != "" {
		params["fields"] = input.Fields
	} else if defaultFields != "" {
		params["fields"] = defaultFields
	}
	return mc.getJSON(ctx, MetaAdsRequest{Method: http.MethodGet, Path: "/" + input.ID, Params: params})
}

func (mc MetaAdsClient) updateObject(ctx context.Context, input MetaAdsObjectInput) (map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, fmt.Errorf("id is required")
	}
	if input.Body == nil {
		return nil, fmt.Errorf("body is required")
	}
	return mc.getJSON(ctx, MetaAdsRequest{Method: http.MethodPost, Path: "/" + input.ID, Body: input.Body})
}

func (mc MetaAdsClient) deleteObject(ctx context.Context, id string) (map[string]any, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("id is required")
	}
	return mc.getJSON(ctx, MetaAdsRequest{Method: http.MethodDelete, Path: "/" + id})
}

func (mc MetaAdsClient) getJSON(ctx context.Context, request MetaAdsRequest) (map[string]any, error) {
	body, err := mc.RequestContext(ctx, request)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func edgeParams(input MetaAdsEdgeInput, defaultFields string) map[string]any {
	params := cloneParams(input.Params)
	if strings.TrimSpace(input.Fields) != "" {
		params["fields"] = input.Fields
	} else if defaultFields != "" {
		params["fields"] = defaultFields
	}
	if strings.TrimSpace(input.Limit) != "" {
		params["limit"] = input.Limit
	}
	if strings.TrimSpace(input.After) != "" {
		params["after"] = input.After
	}
	return params
}

func cloneParams(params map[string]any) map[string]any {
	out := make(map[string]any)
	for key, value := range params {
		out[key] = value
	}
	return out
}

func sortedKeys(value map[string]any) []string {
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
