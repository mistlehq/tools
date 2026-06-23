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

type GSCClient struct {
	searchConsoleBaseURL string
	client               *http.Client
}

type GSCSitesList struct {
	SiteEntry []GSCSite `json:"siteEntry"`
}

type GSCSite struct {
	SiteURL         string `json:"siteUrl"`
	PermissionLevel string `json:"permissionLevel"`
}

type GSCSitemapsList struct {
	Sitemap []GSCSitemap `json:"sitemap"`
}

type GSCSitemap struct {
	Path              string   `json:"path"`
	LastSubmitted     string   `json:"lastSubmitted,omitempty"`
	IsPending         bool     `json:"isPending,omitempty"`
	IsSitemapsIndex   bool     `json:"isSitemapsIndex,omitempty"`
	Type              string   `json:"type,omitempty"`
	LastDownloaded    string   `json:"lastDownloaded,omitempty"`
	Warnings          any      `json:"warnings,omitempty"`
	Errors            any      `json:"errors,omitempty"`
	Contents          []any    `json:"contents,omitempty"`
	SitemapIndex      string   `json:"sitemapIndex,omitempty"`
	Children          []string `json:"children,omitempty"`
	Submitted         string   `json:"submitted,omitempty"`
	Indexed           string   `json:"indexed,omitempty"`
	SubmittedSitemaps string   `json:"submittedSitemaps,omitempty"`
}

type GSCSearchAnalyticsResult = map[string]any
type GSCURLInspectionResult = map[string]any

func NewGSCClient(config Config) GSCClient {
	return GSCClient{
		searchConsoleBaseURL: config.SearchConsoleBaseURL,
		client:               http.DefaultClient,
	}
}

func (gc GSCClient) AuthTest(siteURL string) (GSCSite, error) {
	return gc.AuthTestContext(context.Background(), siteURL)
}

func (gc GSCClient) AuthTestContext(ctx context.Context, siteURL string) (GSCSite, error) {
	return gc.GetSiteContext(ctx, siteURL)
}

func (gc GSCClient) ListSites() (GSCSitesList, error) {
	return gc.ListSitesContext(context.Background())
}

func (gc GSCClient) ListSitesContext(ctx context.Context) (GSCSitesList, error) {
	var out GSCSitesList
	err := gc.getJSON(ctx, "/webmasters/v3/sites", &out)
	return out, err
}

func (gc GSCClient) GetSite(siteURL string) (GSCSite, error) {
	return gc.GetSiteContext(context.Background(), siteURL)
}

func (gc GSCClient) GetSiteContext(ctx context.Context, siteURL string) (GSCSite, error) {
	if err := validateGSCRequired("site URL", siteURL); err != nil {
		return GSCSite{}, err
	}
	var out GSCSite
	err := gc.getJSON(ctx, "/webmasters/v3/sites/"+escapePathPart(siteURL), &out)
	return out, err
}

func (gc GSCClient) QuerySearchAnalytics(siteURL string, request json.RawMessage) (GSCSearchAnalyticsResult, error) {
	return gc.QuerySearchAnalyticsContext(context.Background(), siteURL, request)
}

func (gc GSCClient) QuerySearchAnalyticsContext(ctx context.Context, siteURL string, request json.RawMessage) (GSCSearchAnalyticsResult, error) {
	if err := validateGSCRequired("site URL", siteURL); err != nil {
		return nil, err
	}
	var out GSCSearchAnalyticsResult
	err := gc.postJSON(ctx, "/webmasters/v3/sites/"+escapePathPart(siteURL)+"/searchAnalytics/query", request, &out)
	return out, err
}

func (gc GSCClient) ListSitemaps(siteURL string) (GSCSitemapsList, error) {
	return gc.ListSitemapsContext(context.Background(), siteURL)
}

func (gc GSCClient) ListSitemapsContext(ctx context.Context, siteURL string) (GSCSitemapsList, error) {
	if err := validateGSCRequired("site URL", siteURL); err != nil {
		return GSCSitemapsList{}, err
	}
	var out GSCSitemapsList
	err := gc.getJSON(ctx, "/webmasters/v3/sites/"+escapePathPart(siteURL)+"/sitemaps", &out)
	return out, err
}

func (gc GSCClient) GetSitemap(siteURL string, feedPath string) (GSCSitemap, error) {
	return gc.GetSitemapContext(context.Background(), siteURL, feedPath)
}

func (gc GSCClient) GetSitemapContext(ctx context.Context, siteURL string, feedPath string) (GSCSitemap, error) {
	if err := validateGSCRequired("site URL", siteURL); err != nil {
		return GSCSitemap{}, err
	}
	if err := validateGSCRequired("feed path", feedPath); err != nil {
		return GSCSitemap{}, err
	}
	var out GSCSitemap
	err := gc.getJSON(ctx, "/webmasters/v3/sites/"+escapePathPart(siteURL)+"/sitemaps/"+escapePathPart(feedPath), &out)
	return out, err
}

func (gc GSCClient) InspectURL(request json.RawMessage) (GSCURLInspectionResult, error) {
	return gc.InspectURLContext(context.Background(), request)
}

func (gc GSCClient) InspectURLContext(ctx context.Context, request json.RawMessage) (GSCURLInspectionResult, error) {
	var out GSCURLInspectionResult
	err := gc.postJSON(ctx, "/v1/urlInspection/index:inspect", request, &out)
	return out, err
}

func (gc GSCClient) getJSON(ctx context.Context, path string, out any) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("method path must start with '/': %s", path)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, gc.searchConsoleBaseURL+path, nil)
	if err != nil {
		return err
	}

	return gc.doJSON(request, out)
}

func (gc GSCClient) postJSON(ctx context.Context, path string, body json.RawMessage, out any) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("method path must start with '/': %s", path)
	}
	if len(body) == 0 {
		return fmt.Errorf("request body must not be empty")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, gc.searchConsoleBaseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	return gc.doJSON(request, out)
}

func (gc GSCClient) doJSON(request *http.Request, out any) error {
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
		return fmt.Errorf("google search console api %s %s: rate limited: %s", request.Method, request.URL.Path, string(responseBody))
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("google search console api %s %s failed with status %d: %s", request.Method, request.URL.Path, response.StatusCode, string(responseBody))
	}

	if len(responseBody) == 0 {
		return nil
	}
	return json.Unmarshal(responseBody, out)
}

func validateGSCRequired(label string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", label)
	}
	return nil
}

func escapePathPart(value string) string {
	return url.PathEscape(value)
}
