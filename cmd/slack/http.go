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

type SlackClient struct {
	baseURL string
	client  *http.Client
}

type slackResponseEnvelope struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type SlackAuthTest struct {
	OK     bool   `json:"ok"`
	URL    string `json:"url"`
	Team   string `json:"team"`
	TeamID string `json:"team_id"`
	User   string `json:"user"`
	UserID string `json:"user_id"`
	BotID  string `json:"bot_id"`
}

type SlackResponseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

type SlackConversation struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsPrivate  bool   `json:"is_private"`
	IsArchived bool   `json:"is_archived"`
	IsMember   bool   `json:"is_member"`
	Locale     string `json:"locale"`
}

type SlackConversationsList struct {
	OK               bool                  `json:"ok"`
	Channels         []SlackConversation   `json:"channels"`
	ResponseMetadata SlackResponseMetadata `json:"response_metadata"`
}

type SlackConversationsInfo struct {
	OK      bool              `json:"ok"`
	Channel SlackConversation `json:"channel"`
}

type SlackMessage struct {
	TS       string `json:"ts"`
	ThreadTS string `json:"thread_ts"`
	User     string `json:"user"`
	Type     string `json:"type"`
	Subtype  string `json:"subtype"`
	Text     string `json:"text"`
}

type SlackConversationsHistory struct {
	OK               bool                  `json:"ok"`
	Messages         []SlackMessage        `json:"messages"`
	ResponseMetadata SlackResponseMetadata `json:"response_metadata"`
}

func NewSlackClient(config Config) SlackClient {
	return SlackClient{
		baseURL: config.BaseURL,
		client:  http.DefaultClient,
	}
}

func (sc SlackClient) post(method string, body []byte) ([]byte, error) {
	if !strings.HasPrefix(method, "/") {
		return nil, fmt.Errorf("method path must start with '/': %s", method)
	}

	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		sc.baseURL+method,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := sc.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusTooManyRequests {
		retryAfter := response.Header.Get("Retry-After")
		if retryAfter != "" {
			return nil, fmt.Errorf("slack api %s: rate limited, retry after %s seconds", strings.TrimPrefix(method, "/"), retryAfter)
		}

		return nil, fmt.Errorf("slack api %s: rate limited", strings.TrimPrefix(method, "/"))
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s", response.StatusCode, string(responseBody))
	}

	var envelope slackResponseEnvelope
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return nil, err
	}

	if !envelope.OK {
		return nil, fmt.Errorf("slack api %s: %s", strings.TrimPrefix(method, "/"), envelope.Error)
	}

	return responseBody, nil
}

func (sc SlackClient) get(method string, query url.Values) ([]byte, error) {
	if !strings.HasPrefix(method, "/") {
		return nil, fmt.Errorf("method path must start with '/': %s", method)
	}

	requestURL := sc.baseURL + method
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		requestURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	response, err := sc.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusTooManyRequests {
		retryAfter := response.Header.Get("Retry-After")
		if retryAfter != "" {
			return nil, fmt.Errorf("slack api %s: rate limited, retry after %s seconds", strings.TrimPrefix(method, "/"), retryAfter)
		}

		return nil, fmt.Errorf("slack api %s: rate limited", strings.TrimPrefix(method, "/"))
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s", response.StatusCode, string(responseBody))
	}

	var envelope slackResponseEnvelope
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return nil, err
	}

	if !envelope.OK {
		return nil, fmt.Errorf("slack api %s: %s", strings.TrimPrefix(method, "/"), envelope.Error)
	}

	return responseBody, nil
}

func (sc SlackClient) AuthTest() (SlackAuthTest, error) {
	responseBody, err := sc.post("/auth.test", []byte("{}"))
	if err != nil {
		return SlackAuthTest{}, err
	}

	var authTest SlackAuthTest
	if err := json.Unmarshal(responseBody, &authTest); err != nil {
		return SlackAuthTest{}, err
	}

	return authTest, nil
}

type SlackConversationsListInput struct {
	Types           string
	Limit           string
	Cursor          string
	ExcludeArchived bool
}

func (sc SlackClient) ListConversations(input SlackConversationsListInput) (SlackConversationsList, error) {
	query := url.Values{}
	if input.Types != "" {
		query.Set("types", input.Types)
	}

	if input.Limit != "" {
		query.Set("limit", input.Limit)
	}

	if input.Cursor != "" {
		query.Set("cursor", input.Cursor)
	}

	if input.ExcludeArchived {
		query.Set("exclude_archived", "true")
	}

	responseBody, err := sc.get("/conversations.list", query)
	if err != nil {
		return SlackConversationsList{}, err
	}

	var list SlackConversationsList
	if err := json.Unmarshal(responseBody, &list); err != nil {
		return SlackConversationsList{}, err
	}

	return list, nil
}

type SlackConversationsInfoInput struct {
	Channel       string
	IncludeLocale bool
}

func (sc SlackClient) GetConversationInfo(input SlackConversationsInfoInput) (SlackConversationsInfo, error) {
	query := url.Values{}
	query.Set("channel", input.Channel)

	if input.IncludeLocale {
		query.Set("include_locale", "true")
	}

	responseBody, err := sc.get("/conversations.info", query)
	if err != nil {
		return SlackConversationsInfo{}, err
	}

	var info SlackConversationsInfo
	if err := json.Unmarshal(responseBody, &info); err != nil {
		return SlackConversationsInfo{}, err
	}

	return info, nil
}

type SlackConversationsHistoryInput struct {
	Channel   string
	Cursor    string
	Inclusive bool
	Latest    string
	Limit     string
	Oldest    string
}

func (sc SlackClient) GetConversationHistory(input SlackConversationsHistoryInput) (SlackConversationsHistory, error) {
	query := url.Values{}
	query.Set("channel", input.Channel)

	if input.Cursor != "" {
		query.Set("cursor", input.Cursor)
	}

	if input.Inclusive {
		query.Set("inclusive", "true")
	}

	if input.Latest != "" {
		query.Set("latest", input.Latest)
	}

	if input.Limit != "" {
		query.Set("limit", input.Limit)
	}

	if input.Oldest != "" {
		query.Set("oldest", input.Oldest)
	}

	responseBody, err := sc.get("/conversations.history", query)
	if err != nil {
		return SlackConversationsHistory{}, err
	}

	var history SlackConversationsHistory
	if err := json.Unmarshal(responseBody, &history); err != nil {
		return SlackConversationsHistory{}, err
	}

	return history, nil
}
