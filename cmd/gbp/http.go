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
	"strconv"
	"strings"
)

type GBPClient struct {
	accountManagementBaseURL   string
	businessInformationBaseURL string
	performanceBaseURL         string
	myBusinessBaseURL          string
	client                     *http.Client
}

type GBPAccountsList struct {
	Accounts      []GBPAccount `json:"accounts"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
}

type GBPAccount struct {
	Name                   string `json:"name"`
	AccountName            string `json:"accountName,omitempty"`
	Type                   string `json:"type,omitempty"`
	Role                   string `json:"role,omitempty"`
	VerificationState      string `json:"verificationState,omitempty"`
	VettedState            string `json:"vettedState,omitempty"`
	AccountNumber          string `json:"accountNumber,omitempty"`
	PermissionLevel        string `json:"permissionLevel,omitempty"`
	OrganizationInfo       any    `json:"organizationInfo,omitempty"`
	PrimaryOwner           string `json:"primaryOwner,omitempty"`
	PrimaryOwnerInviteTime string `json:"primaryOwnerInviteTime,omitempty"`
}

type GBPLocationsList struct {
	Locations     []GBPLocation `json:"locations"`
	NextPageToken string        `json:"nextPageToken,omitempty"`
	TotalSize     int           `json:"totalSize,omitempty"`
}

type GBPLocation map[string]any
type GBPReviewsList map[string]any
type GBPReview map[string]any
type GBPMediaList map[string]any
type GBPMediaItem map[string]any
type GBPLocalPostsList map[string]any
type GBPLocalPost map[string]any
type GBPPerformanceResult map[string]any
type GBPWriteResult map[string]any

func NewGBPClient(config Config) GBPClient {
	return GBPClient{
		accountManagementBaseURL:   config.AccountManagementBaseURL,
		businessInformationBaseURL: config.BusinessInformationBaseURL,
		performanceBaseURL:         config.PerformanceBaseURL,
		myBusinessBaseURL:          config.MyBusinessBaseURL,
		client:                     http.DefaultClient,
	}
}

func (gc GBPClient) AuthTest() (GBPAccountsList, error) {
	return gc.AuthTestContext(context.Background())
}

func (gc GBPClient) AuthTestContext(ctx context.Context) (GBPAccountsList, error) {
	return gc.ListAccountsContext(ctx)
}

func (gc GBPClient) ListAccounts() (GBPAccountsList, error) {
	return gc.ListAccountsContext(context.Background())
}

func (gc GBPClient) ListAccountsContext(ctx context.Context) (GBPAccountsList, error) {
	var out GBPAccountsList
	err := gc.getJSON(ctx, gc.accountManagementBaseURL, "/v1/accounts", nil, &out)
	return out, err
}

func (gc GBPClient) GetAccount(account string) (GBPAccount, error) {
	return gc.GetAccountContext(context.Background(), account)
}

func (gc GBPClient) GetAccountContext(ctx context.Context, account string) (GBPAccount, error) {
	if err := validateResourceName("account", account, "accounts"); err != nil {
		return GBPAccount{}, err
	}
	var out GBPAccount
	err := gc.getJSON(ctx, gc.accountManagementBaseURL, "/v1/"+escapeResourceName(account), nil, &out)
	return out, err
}

func (gc GBPClient) ListLocations(account string, readMask string) (GBPLocationsList, error) {
	return gc.ListLocationsContext(context.Background(), account, readMask)
}

func (gc GBPClient) ListLocationsContext(ctx context.Context, account string, readMask string) (GBPLocationsList, error) {
	if err := validateResourceName("account", account, "accounts"); err != nil {
		return GBPLocationsList{}, err
	}
	if err := validateRequired("read mask", readMask); err != nil {
		return GBPLocationsList{}, err
	}
	query := url.Values{"readMask": []string{readMask}}
	var out GBPLocationsList
	err := gc.getJSON(ctx, gc.businessInformationBaseURL, "/v1/"+escapeResourceName(account)+"/locations", query, &out)
	return out, err
}

func (gc GBPClient) GetLocation(location string, readMask string) (GBPLocation, error) {
	return gc.GetLocationContext(context.Background(), location, readMask)
}

func (gc GBPClient) GetLocationContext(ctx context.Context, location string, readMask string) (GBPLocation, error) {
	if err := validateResourceName("location", location, "locations"); err != nil {
		return nil, err
	}
	if err := validateRequired("read mask", readMask); err != nil {
		return nil, err
	}
	query := url.Values{"readMask": []string{readMask}}
	var out GBPLocation
	err := gc.getJSON(ctx, gc.businessInformationBaseURL, "/v1/"+escapeResourceName(location), query, &out)
	return out, err
}

func (gc GBPClient) CreateLocation(account string, request json.RawMessage, options locationWriteOptions) (GBPLocation, error) {
	return gc.CreateLocationContext(context.Background(), account, request, options)
}

func (gc GBPClient) CreateLocationContext(ctx context.Context, account string, request json.RawMessage, options locationWriteOptions) (GBPLocation, error) {
	if err := validateResourceName("account", account, "accounts"); err != nil {
		return nil, err
	}
	query := options.query()
	var out GBPLocation
	err := gc.methodJSON(ctx, http.MethodPost, gc.businessInformationBaseURL, "/v1/"+escapeResourceName(account)+"/locations", query, request, &out)
	return out, err
}

func (gc GBPClient) PatchLocation(location string, request json.RawMessage, options locationPatchOptions) (GBPLocation, error) {
	return gc.PatchLocationContext(context.Background(), location, request, options)
}

func (gc GBPClient) PatchLocationContext(ctx context.Context, location string, request json.RawMessage, options locationPatchOptions) (GBPLocation, error) {
	if err := validateResourceName("location", location, "locations"); err != nil {
		return nil, err
	}
	if err := validateRequired("update mask", options.UpdateMask); err != nil {
		return nil, err
	}
	var out GBPLocation
	err := gc.methodJSON(ctx, http.MethodPatch, gc.businessInformationBaseURL, "/v1/"+escapeResourceName(location), options.query(), request, &out)
	return out, err
}

func (gc GBPClient) DeleteLocation(location string) (GBPWriteResult, error) {
	return gc.DeleteLocationContext(context.Background(), location)
}

func (gc GBPClient) DeleteLocationContext(ctx context.Context, location string) (GBPWriteResult, error) {
	if err := validateResourceName("location", location, "locations"); err != nil {
		return nil, err
	}
	out := GBPWriteResult{}
	err := gc.methodJSON(ctx, http.MethodDelete, gc.businessInformationBaseURL, "/v1/"+escapeResourceName(location), nil, nil, &out)
	return out, err
}

func (gc GBPClient) ListReviews(account string, location string, options pageOptions) (GBPReviewsList, error) {
	return gc.ListReviewsContext(context.Background(), account, location, options)
}

func (gc GBPClient) ListReviewsContext(ctx context.Context, account string, location string, options pageOptions) (GBPReviewsList, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	query := options.query()
	var out GBPReviewsList
	err = gc.getJSON(ctx, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/reviews", query, &out)
	return out, err
}

func (gc GBPClient) GetReview(account string, location string, review string) (GBPReview, error) {
	return gc.GetReviewContext(context.Background(), account, location, review)
}

func (gc GBPClient) GetReviewContext(ctx context.Context, account string, location string, review string) (GBPReview, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	if err := validateRequired("review", review); err != nil {
		return nil, err
	}
	var out GBPReview
	err = gc.getJSON(ctx, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/reviews/"+escapePathPart(review), nil, &out)
	return out, err
}

func (gc GBPClient) UpdateReviewReply(account string, location string, review string, request json.RawMessage) (GBPWriteResult, error) {
	return gc.UpdateReviewReplyContext(context.Background(), account, location, review, request)
}

func (gc GBPClient) UpdateReviewReplyContext(ctx context.Context, account string, location string, review string, request json.RawMessage) (GBPWriteResult, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	if err := validateRequired("review", review); err != nil {
		return nil, err
	}
	out := GBPWriteResult{}
	err = gc.methodJSON(ctx, http.MethodPut, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/reviews/"+escapePathPart(review)+"/reply", nil, request, &out)
	return out, err
}

func (gc GBPClient) DeleteReviewReply(account string, location string, review string) (GBPWriteResult, error) {
	return gc.DeleteReviewReplyContext(context.Background(), account, location, review)
}

func (gc GBPClient) DeleteReviewReplyContext(ctx context.Context, account string, location string, review string) (GBPWriteResult, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	if err := validateRequired("review", review); err != nil {
		return nil, err
	}
	out := GBPWriteResult{}
	err = gc.methodJSON(ctx, http.MethodDelete, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/reviews/"+escapePathPart(review)+"/reply", nil, nil, &out)
	return out, err
}

func (gc GBPClient) ListMedia(account string, location string, options pageOptions) (GBPMediaList, error) {
	return gc.ListMediaContext(context.Background(), account, location, options)
}

func (gc GBPClient) ListMediaContext(ctx context.Context, account string, location string, options pageOptions) (GBPMediaList, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	var out GBPMediaList
	err = gc.getJSON(ctx, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/media", options.query(), &out)
	return out, err
}

func (gc GBPClient) CreateMedia(account string, location string, request json.RawMessage) (GBPMediaItem, error) {
	return gc.CreateMediaContext(context.Background(), account, location, request)
}

func (gc GBPClient) CreateMediaContext(ctx context.Context, account string, location string, request json.RawMessage) (GBPMediaItem, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	var out GBPMediaItem
	err = gc.methodJSON(ctx, http.MethodPost, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/media", nil, request, &out)
	return out, err
}

func (gc GBPClient) GetMedia(media string) (GBPMediaItem, error) {
	return gc.GetMediaContext(context.Background(), media)
}

func (gc GBPClient) GetMediaContext(ctx context.Context, media string) (GBPMediaItem, error) {
	if err := validateResourceName("media", media, "accounts"); err != nil {
		return nil, err
	}
	var out GBPMediaItem
	err := gc.getJSON(ctx, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(media), nil, &out)
	return out, err
}

func (gc GBPClient) PatchMedia(media string, updateMask string, request json.RawMessage) (GBPMediaItem, error) {
	return gc.PatchMediaContext(context.Background(), media, updateMask, request)
}

func (gc GBPClient) PatchMediaContext(ctx context.Context, media string, updateMask string, request json.RawMessage) (GBPMediaItem, error) {
	if err := validateResourceName("media", media, "accounts"); err != nil {
		return nil, err
	}
	if err := validateRequired("update mask", updateMask); err != nil {
		return nil, err
	}
	query := url.Values{"updateMask": []string{updateMask}}
	var out GBPMediaItem
	err := gc.methodJSON(ctx, http.MethodPatch, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(media), query, request, &out)
	return out, err
}

func (gc GBPClient) DeleteMedia(media string) (GBPWriteResult, error) {
	return gc.DeleteMediaContext(context.Background(), media)
}

func (gc GBPClient) DeleteMediaContext(ctx context.Context, media string) (GBPWriteResult, error) {
	if err := validateResourceName("media", media, "accounts"); err != nil {
		return nil, err
	}
	out := GBPWriteResult{}
	err := gc.methodJSON(ctx, http.MethodDelete, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(media), nil, nil, &out)
	return out, err
}

func (gc GBPClient) StartMediaUpload(account string, location string) (GBPWriteResult, error) {
	return gc.StartMediaUploadContext(context.Background(), account, location)
}

func (gc GBPClient) StartMediaUploadContext(ctx context.Context, account string, location string) (GBPWriteResult, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	out := GBPWriteResult{}
	err = gc.methodJSON(ctx, http.MethodPost, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/media:startUpload", nil, json.RawMessage(`{}`), &out)
	return out, err
}

func (gc GBPClient) ListLocalPosts(account string, location string, options pageOptions) (GBPLocalPostsList, error) {
	return gc.ListLocalPostsContext(context.Background(), account, location, options)
}

func (gc GBPClient) ListLocalPostsContext(ctx context.Context, account string, location string, options pageOptions) (GBPLocalPostsList, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	var out GBPLocalPostsList
	err = gc.getJSON(ctx, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/localPosts", options.query(), &out)
	return out, err
}

func (gc GBPClient) CreateLocalPost(account string, location string, request json.RawMessage) (GBPLocalPost, error) {
	return gc.CreateLocalPostContext(context.Background(), account, location, request)
}

func (gc GBPClient) CreateLocalPostContext(ctx context.Context, account string, location string, request json.RawMessage) (GBPLocalPost, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	var out GBPLocalPost
	err = gc.methodJSON(ctx, http.MethodPost, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/localPosts", nil, request, &out)
	return out, err
}

func (gc GBPClient) GetLocalPost(localPost string) (GBPLocalPost, error) {
	return gc.GetLocalPostContext(context.Background(), localPost)
}

func (gc GBPClient) GetLocalPostContext(ctx context.Context, localPost string) (GBPLocalPost, error) {
	if err := validateResourceName("local post", localPost, "accounts"); err != nil {
		return nil, err
	}
	var out GBPLocalPost
	err := gc.getJSON(ctx, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(localPost), nil, &out)
	return out, err
}

func (gc GBPClient) PatchLocalPost(localPost string, updateMask string, request json.RawMessage) (GBPLocalPost, error) {
	return gc.PatchLocalPostContext(context.Background(), localPost, updateMask, request)
}

func (gc GBPClient) PatchLocalPostContext(ctx context.Context, localPost string, updateMask string, request json.RawMessage) (GBPLocalPost, error) {
	if err := validateResourceName("local post", localPost, "accounts"); err != nil {
		return nil, err
	}
	if err := validateRequired("update mask", updateMask); err != nil {
		return nil, err
	}
	query := url.Values{"updateMask": []string{updateMask}}
	var out GBPLocalPost
	err := gc.methodJSON(ctx, http.MethodPatch, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(localPost), query, request, &out)
	return out, err
}

func (gc GBPClient) DeleteLocalPost(localPost string) (GBPWriteResult, error) {
	return gc.DeleteLocalPostContext(context.Background(), localPost)
}

func (gc GBPClient) DeleteLocalPostContext(ctx context.Context, localPost string) (GBPWriteResult, error) {
	if err := validateResourceName("local post", localPost, "accounts"); err != nil {
		return nil, err
	}
	out := GBPWriteResult{}
	err := gc.methodJSON(ctx, http.MethodDelete, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(localPost), nil, nil, &out)
	return out, err
}

func (gc GBPClient) ReportLocalPostInsights(account string, location string, request json.RawMessage) (GBPWriteResult, error) {
	return gc.ReportLocalPostInsightsContext(context.Background(), account, location, request)
}

func (gc GBPClient) ReportLocalPostInsightsContext(ctx context.Context, account string, location string, request json.RawMessage) (GBPWriteResult, error) {
	parent, err := accountLocationParent(account, location)
	if err != nil {
		return nil, err
	}
	out := GBPWriteResult{}
	err = gc.methodJSON(ctx, http.MethodPost, gc.myBusinessBaseURL, "/v4/"+escapeResourceName(parent)+"/localPosts:reportInsights", nil, request, &out)
	return out, err
}

func (gc GBPClient) GetDailyMetrics(location string, request json.RawMessage) (GBPPerformanceResult, error) {
	return gc.GetDailyMetricsContext(context.Background(), location, request)
}

func (gc GBPClient) GetDailyMetricsContext(ctx context.Context, location string, request json.RawMessage) (GBPPerformanceResult, error) {
	if err := validateResourceName("location", location, "locations"); err != nil {
		return nil, err
	}
	query, err := queryFromJSON(request)
	if err != nil {
		return nil, err
	}
	var out GBPPerformanceResult
	err = gc.getJSON(ctx, gc.performanceBaseURL, "/v1/"+escapeResourceName(location)+":getDailyMetricsTimeSeries", query, &out)
	return out, err
}

func (gc GBPClient) ListSearchKeywords(location string, request json.RawMessage) (GBPPerformanceResult, error) {
	return gc.ListSearchKeywordsContext(context.Background(), location, request)
}

func (gc GBPClient) ListSearchKeywordsContext(ctx context.Context, location string, request json.RawMessage) (GBPPerformanceResult, error) {
	if err := validateResourceName("location", location, "locations"); err != nil {
		return nil, err
	}
	query, err := queryFromJSON(request)
	if err != nil {
		return nil, err
	}
	var out GBPPerformanceResult
	err = gc.getJSON(ctx, gc.performanceBaseURL, "/v1/"+escapeResourceName(location)+"/searchkeywords/impressions/monthly", query, &out)
	return out, err
}

func (gc GBPClient) getJSON(ctx context.Context, baseURL string, path string, query url.Values, out any) error {
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

func (gc GBPClient) methodJSON(ctx context.Context, method string, baseURL string, path string, query url.Values, body json.RawMessage, out any) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("method path must start with '/': %s", path)
	}
	requestURL := baseURL + path
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
	if len(body) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}

	return gc.doJSON(request, out)
}

func (gc GBPClient) doJSON(request *http.Request, out any) error {
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
		return fmt.Errorf("google business profile api %s %s: rate limited: %s", request.Method, request.URL.Path, string(responseBody))
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("google business profile api %s %s failed with status %d: %s", request.Method, request.URL.Path, response.StatusCode, string(responseBody))
	}

	if len(responseBody) == 0 {
		return nil
	}
	return json.Unmarshal(responseBody, out)
}

type pageOptions struct {
	PageSize  string
	PageToken string
	OrderBy   string
}

type locationWriteOptions struct {
	RequestID    string
	ValidateOnly string
}

func (options locationWriteOptions) query() url.Values {
	query := url.Values{}
	if options.RequestID != "" {
		query.Set("requestId", options.RequestID)
	}
	if options.ValidateOnly != "" {
		query.Set("validateOnly", options.ValidateOnly)
	}
	return query
}

type locationPatchOptions struct {
	UpdateMask    string
	AttributeMask string
	ValidateOnly  string
}

func (options locationPatchOptions) query() url.Values {
	query := url.Values{"updateMask": []string{options.UpdateMask}}
	if options.AttributeMask != "" {
		query.Set("attributeMask", options.AttributeMask)
	}
	if options.ValidateOnly != "" {
		query.Set("validateOnly", options.ValidateOnly)
	}
	return query
}

func (options pageOptions) query() url.Values {
	query := url.Values{}
	if options.PageSize != "" {
		query.Set("pageSize", options.PageSize)
	}
	if options.PageToken != "" {
		query.Set("pageToken", options.PageToken)
	}
	if options.OrderBy != "" {
		query.Set("orderBy", options.OrderBy)
	}
	return query
}

func accountLocationParent(account string, location string) (string, error) {
	if err := validateResourceName("account", account, "accounts"); err != nil {
		return "", err
	}
	if err := validateResourceName("location", location, "locations"); err != nil {
		return "", err
	}
	return account + "/" + location, nil
}

func validateRequired(label string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", label)
	}
	return nil
}

func validateResourceName(label string, value string, prefix string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", label)
	}
	if !strings.HasPrefix(value, prefix+"/") {
		return fmt.Errorf("%s must use %s/<id> format", label, prefix)
	}
	return nil
}

func escapeResourceName(value string) string {
	parts := strings.Split(value, "/")
	for index, part := range parts {
		parts[index] = escapePathPart(part)
	}
	return strings.Join(parts, "/")
}

func escapePathPart(value string) string {
	return url.PathEscape(value)
}

func queryFromJSON(request json.RawMessage) (url.Values, error) {
	if len(request) == 0 {
		return nil, fmt.Errorf("request body must not be empty")
	}
	var parsed map[string]any
	if err := json.Unmarshal(request, &parsed); err != nil {
		return nil, err
	}
	query := url.Values{}
	flattenQueryJSON(query, "", parsed)
	return query, nil
}

func flattenQueryJSON(query url.Values, prefix string, value any) {
	switch typed := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			nextPrefix := key
			if prefix != "" {
				nextPrefix = prefix + "." + key
			}
			flattenQueryJSON(query, nextPrefix, typed[key])
		}
	case []any:
		for _, item := range typed {
			flattenQueryJSON(query, prefix, item)
		}
	case string:
		query.Add(prefix, typed)
	case float64:
		query.Add(prefix, strconv.FormatFloat(typed, 'f', -1, 64))
	case bool:
		query.Add(prefix, strconv.FormatBool(typed))
	case nil:
	default:
		query.Add(prefix, fmt.Sprint(typed))
	}
}
