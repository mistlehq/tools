package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GoogleAdsClient struct {
	baseURL string
	client  *http.Client
}

type GoogleAdsRequest struct {
	Method          string         `json:"method"`
	Path            string         `json:"path"`
	LoginCustomerID string         `json:"loginCustomerId,omitempty"`
	Params          map[string]any `json:"params,omitempty"`
	Body            map[string]any `json:"body,omitempty"`
}

type GoogleAdsGAQLInput struct {
	CustomerID      string
	LoginCustomerID string
	Query           string
	PageSize        string
	PageToken       string
	SummaryRow      string
	Params          map[string]any
}

func NewGoogleAdsClient(config Config) GoogleAdsClient {
	return GoogleAdsClient{baseURL: config.BaseURL, client: http.DefaultClient}
}

func (gc GoogleAdsClient) Request(request GoogleAdsRequest) ([]byte, error) {
	return gc.RequestContext(context.Background(), request)
}

func (gc GoogleAdsClient) RequestContext(ctx context.Context, request GoogleAdsRequest) ([]byte, error) {
	method := strings.ToUpper(strings.TrimSpace(request.Method))
	if method == "" {
		method = http.MethodGet
	}
	if method != http.MethodGet && method != http.MethodPost && method != http.MethodDelete && method != http.MethodPatch {
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
	if strings.TrimSpace(request.Path) == "" {
		return nil, fmt.Errorf("path is required")
	}
	if !strings.HasPrefix(request.Path, "/") {
		return nil, fmt.Errorf("path must start with '/'")
	}

	requestURL, err := url.Parse(gc.baseURL + request.Path)
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
	if strings.TrimSpace(request.LoginCustomerID) != "" {
		httpRequest.Header.Set("login-customer-id", request.LoginCustomerID)
	}

	response, err := gc.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("google ads api request failed with status %d: %s", response.StatusCode, string(responseBody))
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

func (gc GoogleAdsClient) AuthTest(ctx context.Context) (map[string]any, error) {
	return gc.ListAccessibleCustomers(ctx)
}

func (gc GoogleAdsClient) ListAccessibleCustomers(ctx context.Context) (map[string]any, error) {
	return gc.getJSON(ctx, GoogleAdsRequest{Method: http.MethodGet, Path: "/customers:listAccessibleCustomers"})
}

func (gc GoogleAdsClient) Search(ctx context.Context, input GoogleAdsGAQLInput) (map[string]any, error) {
	body, err := gaqlRequestBody(input, true)
	if err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GoogleAdsRequest{Method: http.MethodPost, Path: "/customers/" + input.CustomerID + "/googleAds:search", LoginCustomerID: input.LoginCustomerID, Body: body})
}

func (gc GoogleAdsClient) SearchStream(ctx context.Context, input GoogleAdsGAQLInput) ([]any, error) {
	body, err := gaqlRequestBody(input, false)
	if err != nil {
		return nil, err
	}
	responseBody, err := gc.RequestContext(ctx, GoogleAdsRequest{Method: http.MethodPost, Path: "/customers/" + input.CustomerID + "/googleAds:searchStream", LoginCustomerID: input.LoginCustomerID, Body: body})
	if err != nil {
		return nil, err
	}
	var out []any
	if err := json.Unmarshal(responseBody, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (gc GoogleAdsClient) SearchFields(ctx context.Context, query string) (map[string]any, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query is required")
	}
	return gc.postJSON(ctx, GoogleAdsRequest{Method: http.MethodPost, Path: "/googleAdsFields:search", Body: map[string]any{"query": query}})
}

func (gc GoogleAdsClient) GetField(ctx context.Context, resourceName string) (map[string]any, error) {
	if strings.TrimSpace(resourceName) == "" {
		return nil, fmt.Errorf("resource-name is required")
	}
	return gc.getJSON(ctx, GoogleAdsRequest{Method: http.MethodGet, Path: "/" + resourceName})
}

func gaqlRequestBody(input GoogleAdsGAQLInput, includePagination bool) (map[string]any, error) {
	if strings.TrimSpace(input.CustomerID) == "" {
		return nil, fmt.Errorf("customer-id is required")
	}
	if strings.TrimSpace(input.Query) == "" {
		return nil, fmt.Errorf("query is required")
	}
	body := cloneParams(input.Params)
	body["query"] = input.Query
	if includePagination {
		if strings.TrimSpace(input.PageSize) != "" {
			body["pageSize"] = input.PageSize
		}
		if strings.TrimSpace(input.PageToken) != "" {
			body["pageToken"] = input.PageToken
		}
	}
	if strings.TrimSpace(input.SummaryRow) != "" {
		body["summaryRowSetting"] = input.SummaryRow
	}
	return body, nil
}

func (gc GoogleAdsClient) getJSON(ctx context.Context, request GoogleAdsRequest) (map[string]any, error) {
	body, err := gc.RequestContext(ctx, request)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (gc GoogleAdsClient) postJSON(ctx context.Context, request GoogleAdsRequest) (map[string]any, error) {
	responseBody, err := gc.RequestContext(ctx, request)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(responseBody, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func cloneParams(params map[string]any) map[string]any {
	out := make(map[string]any)
	for key, value := range params {
		out[key] = value
	}
	return out
}
