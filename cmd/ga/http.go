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

type GAClient struct {
	analyticsDataBaseURL  string
	analyticsAdminBaseURL string
	client                *http.Client
}

type GAAccountSummariesList struct {
	AccountSummaries []GAAccountSummary `json:"accountSummaries"`
	NextPageToken    string             `json:"nextPageToken,omitempty"`
}

type GAAccountSummary struct {
	Name              string              `json:"name"`
	Account           string              `json:"account"`
	DisplayName       string              `json:"displayName"`
	PropertySummaries []GAPropertySummary `json:"propertySummaries,omitempty"`
}

type GAPropertySummary struct {
	Property     string `json:"property"`
	DisplayName  string `json:"displayName"`
	PropertyType string `json:"propertyType,omitempty"`
	Parent       string `json:"parent,omitempty"`
}

type GAProperty struct {
	Name         string `json:"name"`
	Parent       string `json:"parent,omitempty"`
	DisplayName  string `json:"displayName"`
	CurrencyCode string `json:"currencyCode,omitempty"`
	TimeZone     string `json:"timeZone,omitempty"`
	CreateTime   string `json:"createTime,omitempty"`
	UpdateTime   string `json:"updateTime,omitempty"`
}

type GAMetadata struct {
	Name       string                `json:"name"`
	Dimensions []GAMetadataDimension `json:"dimensions"`
	Metrics    []GAMetadataMetric    `json:"metrics"`
}

type GAMetadataDimension struct {
	APIName          string `json:"apiName"`
	UIName           string `json:"uiName"`
	Description      string `json:"description,omitempty"`
	Category         string `json:"category,omitempty"`
	CustomDefinition bool   `json:"customDefinition,omitempty"`
}

type GAMetadataMetric struct {
	APIName          string `json:"apiName"`
	UIName           string `json:"uiName"`
	Description      string `json:"description,omitempty"`
	Category         string `json:"category,omitempty"`
	Type             string `json:"type,omitempty"`
	Expression       string `json:"expression,omitempty"`
	CustomDefinition bool   `json:"customDefinition,omitempty"`
}

type GAGoogleAdsLinksList struct {
	GoogleAdsLinks []GAGoogleAdsLink `json:"googleAdsLinks"`
	NextPageToken  string            `json:"nextPageToken,omitempty"`
}

type GAGoogleAdsLink struct {
	Name                      string `json:"name"`
	CustomerID                string `json:"customerId"`
	CanManageClients          bool   `json:"canManageClients,omitempty"`
	AdsPersonalizationEnabled bool   `json:"adsPersonalizationEnabled,omitempty"`
	CreatorEmailAddress       string `json:"creatorEmailAddress,omitempty"`
}

type GACompatibilityCheckResult = map[string]any
type GAReportResult = map[string]any

func NewGAClient(config Config) GAClient {
	return GAClient{
		analyticsDataBaseURL:  config.AnalyticsDataBaseURL,
		analyticsAdminBaseURL: config.AnalyticsAdminBaseURL,
		client:                http.DefaultClient,
	}
}

func (gc GAClient) AuthTest(property string) (GAProperty, error) {
	return gc.AuthTestContext(context.Background(), property)
}

func (gc GAClient) AuthTestContext(ctx context.Context, property string) (GAProperty, error) {
	return gc.GetPropertyContext(ctx, property)
}

func (gc GAClient) ListAccountSummaries() (GAAccountSummariesList, error) {
	return gc.ListAccountSummariesContext(context.Background())
}

func (gc GAClient) ListAccountSummariesContext(ctx context.Context) (GAAccountSummariesList, error) {
	var out GAAccountSummariesList
	err := gc.getJSON(ctx, gc.analyticsAdminBaseURL, "/v1beta/accountSummaries", nil, &out)
	return out, err
}

func (gc GAClient) GetProperty(property string) (GAProperty, error) {
	return gc.GetPropertyContext(context.Background(), property)
}

func (gc GAClient) GetPropertyContext(ctx context.Context, property string) (GAProperty, error) {
	if err := validateGAResourceName("property", property, "properties/"); err != nil {
		return GAProperty{}, err
	}
	var out GAProperty
	err := gc.getJSON(ctx, gc.analyticsAdminBaseURL, "/v1beta/"+property, nil, &out)
	return out, err
}

func (gc GAClient) GetMetadata(property string) (GAMetadata, error) {
	return gc.GetMetadataContext(context.Background(), property)
}

func (gc GAClient) GetMetadataContext(ctx context.Context, property string) (GAMetadata, error) {
	if err := validateGAResourceName("property", property, "properties/"); err != nil {
		return GAMetadata{}, err
	}
	var out GAMetadata
	err := gc.getJSON(ctx, gc.analyticsDataBaseURL, "/v1beta/"+property+"/metadata", nil, &out)
	return out, err
}

func (gc GAClient) CheckCompatibility(property string, request json.RawMessage) (GACompatibilityCheckResult, error) {
	return gc.CheckCompatibilityContext(context.Background(), property, request)
}

func (gc GAClient) CheckCompatibilityContext(ctx context.Context, property string, request json.RawMessage) (GACompatibilityCheckResult, error) {
	if err := validateGAResourceName("property", property, "properties/"); err != nil {
		return nil, err
	}
	var out GACompatibilityCheckResult
	err := gc.postJSON(ctx, gc.analyticsDataBaseURL, "/v1beta/"+property+":checkCompatibility", request, &out)
	return out, err
}

func (gc GAClient) RunReport(property string, request json.RawMessage) (GAReportResult, error) {
	return gc.RunReportContext(context.Background(), property, request)
}

func (gc GAClient) RunReportContext(ctx context.Context, property string, request json.RawMessage) (GAReportResult, error) {
	if err := validateGAResourceName("property", property, "properties/"); err != nil {
		return nil, err
	}
	var out GAReportResult
	err := gc.postJSON(ctx, gc.analyticsDataBaseURL, "/v1beta/"+property+":runReport", request, &out)
	return out, err
}

func (gc GAClient) RunRealtimeReport(property string, request json.RawMessage) (GAReportResult, error) {
	return gc.RunRealtimeReportContext(context.Background(), property, request)
}

func (gc GAClient) RunRealtimeReportContext(ctx context.Context, property string, request json.RawMessage) (GAReportResult, error) {
	if err := validateGAResourceName("property", property, "properties/"); err != nil {
		return nil, err
	}
	var out GAReportResult
	err := gc.postJSON(ctx, gc.analyticsDataBaseURL, "/v1beta/"+property+":runRealtimeReport", request, &out)
	return out, err
}

func (gc GAClient) RunFunnelReport(property string, request json.RawMessage) (GAReportResult, error) {
	return gc.RunFunnelReportContext(context.Background(), property, request)
}

func (gc GAClient) RunFunnelReportContext(ctx context.Context, property string, request json.RawMessage) (GAReportResult, error) {
	if err := validateGAResourceName("property", property, "properties/"); err != nil {
		return nil, err
	}
	var out GAReportResult
	err := gc.postJSON(ctx, gc.analyticsDataBaseURL, "/v1alpha/"+property+":runFunnelReport", request, &out)
	return out, err
}

func (gc GAClient) ListGoogleAdsLinks(property string) (GAGoogleAdsLinksList, error) {
	return gc.ListGoogleAdsLinksContext(context.Background(), property)
}

func (gc GAClient) ListGoogleAdsLinksContext(ctx context.Context, property string) (GAGoogleAdsLinksList, error) {
	if err := validateGAResourceName("property", property, "properties/"); err != nil {
		return GAGoogleAdsLinksList{}, err
	}
	var out GAGoogleAdsLinksList
	err := gc.getJSON(ctx, gc.analyticsAdminBaseURL, "/v1beta/"+property+"/googleAdsLinks", nil, &out)
	return out, err
}

func (gc GAClient) getJSON(ctx context.Context, baseURL string, path string, query url.Values, out any) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("method path must start with '/': %s", path)
	}

	requestURL := baseURL + path
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	return gc.doJSON(request, out)
}

func (gc GAClient) postJSON(ctx context.Context, baseURL string, path string, body json.RawMessage, out any) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("method path must start with '/': %s", path)
	}
	if len(body) == 0 {
		return fmt.Errorf("request body must not be empty")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	return gc.doJSON(request, out)
}

func (gc GAClient) doJSON(request *http.Request, out any) error {
	response, err := gc.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("google analytics api %s %s: rate limited: %s", request.Method, request.URL.Path, string(responseBody))
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("google analytics api %s %s failed with status %d: %s", request.Method, request.URL.Path, response.StatusCode, string(responseBody))
	}

	if len(responseBody) == 0 {
		return nil
	}
	return json.Unmarshal(responseBody, out)
}

func validateGAResourceName(label string, value string, prefix string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", label)
	}
	if !strings.HasPrefix(value, prefix) {
		return fmt.Errorf("%s must start with %s", label, prefix)
	}
	return nil
}
