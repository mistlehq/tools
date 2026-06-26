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

type XeroClient struct {
	apiBaseURL string
	client     *http.Client
}

type XeroTenantConnection struct {
	ID             string `json:"id,omitempty"`
	TenantID       string `json:"tenantId,omitempty"`
	TenantType     string `json:"tenantType,omitempty"`
	TenantName     string `json:"tenantName,omitempty"`
	CreatedDateUTC string `json:"createdDateUtc,omitempty"`
	UpdatedDateUTC string `json:"updatedDateUtc,omitempty"`
}

type XeroTenantConnectionsList struct {
	Connections []XeroTenantConnection `json:"connections"`
}

type XeroJSONResult map[string]any

type XeroAPIRequest struct {
	Family   string
	TenantID string
	Endpoint string
	Query    map[string]string
	Body     json.RawMessage
}

var xeroAPIFamilyBasePaths = map[string]string{
	"accounting": "/api.xro/2.0",
	"assets":     "/assets.xro/1.0",
	"files":      "/files.xro/1.0",
	"projects":   "/projects.xro/2.0",
}

func NewXeroClient(config Config) XeroClient {
	return XeroClient{
		apiBaseURL: config.APIBaseURL,
		client:     http.DefaultClient,
	}
}

func (xc XeroClient) ListTenants() (XeroTenantConnectionsList, error) {
	return xc.ListTenantsContext(context.Background())
}

func (xc XeroClient) ListTenantsContext(ctx context.Context) (XeroTenantConnectionsList, error) {
	var connections []XeroTenantConnection
	err := xc.methodJSON(ctx, http.MethodGet, "/connections", nil, "", nil, &connections)
	return XeroTenantConnectionsList{Connections: connections}, err
}

func (xc XeroClient) GetAPIEndpoint(request XeroAPIRequest) (XeroJSONResult, error) {
	return xc.GetAPIEndpointContext(context.Background(), request)
}

func (xc XeroClient) GetAPIEndpointContext(ctx context.Context, request XeroAPIRequest) (XeroJSONResult, error) {
	return xc.callAPIEndpoint(ctx, http.MethodGet, request)
}

func (xc XeroClient) PostAPIEndpoint(request XeroAPIRequest) (XeroJSONResult, error) {
	return xc.PostAPIEndpointContext(context.Background(), request)
}

func (xc XeroClient) PostAPIEndpointContext(ctx context.Context, request XeroAPIRequest) (XeroJSONResult, error) {
	return xc.callAPIEndpoint(ctx, http.MethodPost, request)
}

func (xc XeroClient) PutAPIEndpoint(request XeroAPIRequest) (XeroJSONResult, error) {
	return xc.PutAPIEndpointContext(context.Background(), request)
}

func (xc XeroClient) PutAPIEndpointContext(ctx context.Context, request XeroAPIRequest) (XeroJSONResult, error) {
	return xc.callAPIEndpoint(ctx, http.MethodPut, request)
}

func (xc XeroClient) DeleteAPIEndpoint(request XeroAPIRequest) (XeroJSONResult, error) {
	return xc.DeleteAPIEndpointContext(context.Background(), request)
}

func (xc XeroClient) DeleteAPIEndpointContext(ctx context.Context, request XeroAPIRequest) (XeroJSONResult, error) {
	return xc.callAPIEndpoint(ctx, http.MethodDelete, request)
}

func (xc XeroClient) callAPIEndpoint(ctx context.Context, method string, request XeroAPIRequest) (XeroJSONResult, error) {
	path, err := xeroAPIPath(request.Family, request.Endpoint)
	if err != nil {
		return nil, err
	}
	if err := validateRequired("tenantId", request.TenantID); err != nil {
		return nil, err
	}

	var out XeroJSONResult
	err = xc.methodJSON(ctx, method, path, queryFromMap(request.Query), request.TenantID, request.Body, &out)
	return out, err
}

func xeroAPIPath(family string, endpoint string) (string, error) {
	basePath, ok := xeroAPIFamilyBasePaths[family]
	if !ok {
		families := make([]string, 0, len(xeroAPIFamilyBasePaths))
		for name := range xeroAPIFamilyBasePaths {
			families = append(families, name)
		}
		sort.Strings(families)
		return "", fmt.Errorf("unsupported Xero API family %q; expected one of %s", family, strings.Join(families, ", "))
	}
	if err := validateEndpoint(endpoint); err != nil {
		return "", err
	}
	return basePath + endpoint, nil
}

func validateEndpoint(endpoint string) error {
	if strings.TrimSpace(endpoint) == "" {
		return fmt.Errorf("endpoint is required")
	}
	if !strings.HasPrefix(endpoint, "/") {
		return fmt.Errorf("endpoint must start with '/'")
	}
	if strings.Contains(endpoint, "://") {
		return fmt.Errorf("endpoint must be a path, not a URL")
	}
	return nil
}

func validateRequired(label string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", label)
	}
	return nil
}

func queryFromMap(values map[string]string) url.Values {
	query := url.Values{}
	for name, value := range values {
		if value != "" {
			query.Set(name, value)
		}
	}
	return query
}

func (xc XeroClient) methodJSON(ctx context.Context, method string, path string, query url.Values, tenantID string, body json.RawMessage, out any) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("method path must start with '/': %s", path)
	}
	requestURL := xc.apiBaseURL + path
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}

	request, err := http.NewRequestWithContext(ctx, method, requestURL, reader)
	if err != nil {
		return err
	}
	if tenantID != "" {
		request.Header.Set("xero-tenant-id", tenantID)
	}
	if len(body) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}
	request.Header.Set("Accept", "application/json")

	return xc.doJSON(request, out)
}

func (xc XeroClient) doJSON(request *http.Request, out any) error {
	response, err := xc.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("xero api %s %s: rate limited: %s", request.Method, request.URL.Path, string(responseBody))
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("xero api %s %s failed with status %d: %s", request.Method, request.URL.Path, response.StatusCode, string(responseBody))
	}

	if len(responseBody) == 0 {
		return nil
	}
	return json.Unmarshal(responseBody, out)
}
