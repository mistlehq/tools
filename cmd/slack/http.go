package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

type SlackConversationsReplies struct {
	OK               bool                  `json:"ok"`
	Messages         []SlackMessage        `json:"messages"`
	ResponseMetadata SlackResponseMetadata `json:"response_metadata"`
}

type SlackChatMessage struct {
	Text     string `json:"text"`
	ThreadTS string `json:"thread_ts"`
}

type SlackChatPostMessage struct {
	OK      bool             `json:"ok"`
	Channel string           `json:"channel"`
	TS      string           `json:"ts"`
	Message SlackChatMessage `json:"message"`
}

type SlackChatUpdate struct {
	OK      bool             `json:"ok"`
	Channel string           `json:"channel"`
	TS      string           `json:"ts"`
	Text    string           `json:"text"`
	Message SlackChatMessage `json:"message"`
}

type SlackChatDelete struct {
	OK      bool   `json:"ok"`
	Channel string `json:"channel"`
	TS      string `json:"ts"`
}

type SlackChatGetPermalink struct {
	OK        bool   `json:"ok"`
	Channel   string `json:"channel"`
	Permalink string `json:"permalink"`
}

type SlackReactionResponse struct {
	OK bool `json:"ok"`
}

type SlackEmojiList struct {
	OK         bool              `json:"ok"`
	Emoji      map[string]string `json:"emoji"`
	Categories json.RawMessage   `json:"categories,omitempty"`
}

type SlackFile struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Title              string `json:"title"`
	Size               int64  `json:"size"`
	URLPrivate         string `json:"url_private"`
	URLPrivateDownload string `json:"url_private_download"`
	Permalink          string `json:"permalink"`
}

type SlackFilesInfo struct {
	OK   bool      `json:"ok"`
	File SlackFile `json:"file"`
}

type SlackFilesGetUploadURLExternal struct {
	OK        bool   `json:"ok"`
	UploadURL string `json:"upload_url"`
	FileID    string `json:"file_id"`
}

type SlackFilesCompleteUploadExternal struct {
	OK    bool        `json:"ok"`
	Files []SlackFile `json:"files"`
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

func (sc SlackClient) postMultipart(method string, values map[string]string) ([]byte, error) {
	if !strings.HasPrefix(method, "/") {
		return nil, fmt.Errorf("method path must start with '/': %s", method)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range values {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		sc.baseURL+method,
		&body,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

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

type SlackConversationsRepliesInput struct {
	Channel   string
	TS        string
	Cursor    string
	Inclusive bool
	Latest    string
	Limit     string
	Oldest    string
}

func (sc SlackClient) GetConversationReplies(input SlackConversationsRepliesInput) (SlackConversationsReplies, error) {
	query := url.Values{}
	query.Set("channel", input.Channel)
	query.Set("ts", input.TS)

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

	responseBody, err := sc.get("/conversations.replies", query)
	if err != nil {
		return SlackConversationsReplies{}, err
	}

	var replies SlackConversationsReplies
	if err := json.Unmarshal(responseBody, &replies); err != nil {
		return SlackConversationsReplies{}, err
	}

	return replies, nil
}

type SlackChatPostMessageInput struct {
	Channel  string  `json:"channel"`
	Text     string  `json:"text"`
	ThreadTS *string `json:"thread_ts,omitempty"`
}

func (sc SlackClient) PostMessage(input SlackChatPostMessageInput) (SlackChatPostMessage, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return SlackChatPostMessage{}, err
	}

	responseBody, err := sc.post("/chat.postMessage", requestBody)
	if err != nil {
		return SlackChatPostMessage{}, err
	}

	var posted SlackChatPostMessage
	if err := json.Unmarshal(responseBody, &posted); err != nil {
		return SlackChatPostMessage{}, err
	}

	return posted, nil
}

type SlackChatUpdateInput struct {
	Channel string `json:"channel"`
	TS      string `json:"ts"`
	Text    string `json:"text"`
}

func (sc SlackClient) UpdateMessage(input SlackChatUpdateInput) (SlackChatUpdate, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return SlackChatUpdate{}, err
	}

	responseBody, err := sc.post("/chat.update", requestBody)
	if err != nil {
		return SlackChatUpdate{}, err
	}

	var updated SlackChatUpdate
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return SlackChatUpdate{}, err
	}

	return updated, nil
}

type SlackChatDeleteInput struct {
	Channel string `json:"channel"`
	TS      string `json:"ts"`
}

func (sc SlackClient) DeleteMessage(input SlackChatDeleteInput) (SlackChatDelete, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return SlackChatDelete{}, err
	}

	responseBody, err := sc.post("/chat.delete", requestBody)
	if err != nil {
		return SlackChatDelete{}, err
	}

	var deleted SlackChatDelete
	if err := json.Unmarshal(responseBody, &deleted); err != nil {
		return SlackChatDelete{}, err
	}

	return deleted, nil
}

func (sc SlackClient) GetPermalink(channel string, messageTS string) (SlackChatGetPermalink, error) {
	query := url.Values{}
	query.Set("channel", channel)
	query.Set("message_ts", messageTS)

	responseBody, err := sc.get("/chat.getPermalink", query)
	if err != nil {
		return SlackChatGetPermalink{}, err
	}

	var permalink SlackChatGetPermalink
	if err := json.Unmarshal(responseBody, &permalink); err != nil {
		return SlackChatGetPermalink{}, err
	}

	return permalink, nil
}

type SlackReactionInput struct {
	Channel   string `json:"channel"`
	Timestamp string `json:"timestamp"`
	Name      string `json:"name"`
}

func (sc SlackClient) AddReaction(input SlackReactionInput) (SlackReactionResponse, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return SlackReactionResponse{}, err
	}

	responseBody, err := sc.post("/reactions.add", requestBody)
	if err != nil {
		return SlackReactionResponse{}, err
	}

	var response SlackReactionResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return SlackReactionResponse{}, err
	}

	return response, nil
}

func (sc SlackClient) RemoveReaction(input SlackReactionInput) (SlackReactionResponse, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return SlackReactionResponse{}, err
	}

	responseBody, err := sc.post("/reactions.remove", requestBody)
	if err != nil {
		return SlackReactionResponse{}, err
	}

	var response SlackReactionResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return SlackReactionResponse{}, err
	}

	return response, nil
}

func (sc SlackClient) ListEmoji(includeCategories bool) (SlackEmojiList, error) {
	query := url.Values{}
	if includeCategories {
		query.Set("include_categories", "true")
	}

	responseBody, err := sc.get("/emoji.list", query)
	if err != nil {
		return SlackEmojiList{}, err
	}

	var list SlackEmojiList
	if err := json.Unmarshal(responseBody, &list); err != nil {
		return SlackEmojiList{}, err
	}

	return list, nil
}

type SlackFilesUploadInput struct {
	Path           string
	Channel        string
	ThreadTS       string
	InitialComment string
}

type SlackFilesDownloadInput struct {
	FileID string
	Output string
}

type SlackFilesDownload struct {
	FileID string `json:"file_id"`
	Output string `json:"output"`
	Bytes  int    `json:"bytes"`
}

func (sc SlackClient) GetFileInfo(fileID string) (SlackFilesInfo, error) {
	query := url.Values{}
	query.Set("file", fileID)

	responseBody, err := sc.get("/files.info", query)
	if err != nil {
		return SlackFilesInfo{}, err
	}

	var info SlackFilesInfo
	if err := json.Unmarshal(responseBody, &info); err != nil {
		return SlackFilesInfo{}, err
	}

	return info, nil
}

func (sc SlackClient) DownloadFile(input SlackFilesDownloadInput) (SlackFilesDownload, error) {
	info, err := sc.GetFileInfo(input.FileID)
	if err != nil {
		return SlackFilesDownload{}, err
	}

	downloadURL := resolveSlackFileDownloadURL(info.File)
	if downloadURL == "" {
		return SlackFilesDownload{}, fmt.Errorf("slack file %s does not include url_private_download or url_private", input.FileID)
	}

	fileBytes, err := sc.downloadFileBytes(downloadURL)
	if err != nil {
		return SlackFilesDownload{}, err
	}

	if err := os.WriteFile(input.Output, fileBytes, 0o600); err != nil {
		return SlackFilesDownload{}, err
	}

	return SlackFilesDownload{
		FileID: input.FileID,
		Output: input.Output,
		Bytes:  len(fileBytes),
	}, nil
}

func resolveSlackFileDownloadURL(file SlackFile) string {
	if file.URLPrivateDownload != "" {
		return file.URLPrivateDownload
	}

	return file.URLPrivate
}

func (sc SlackClient) downloadFileBytes(downloadURL string) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		downloadURL,
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

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("file download failed with status %d: %s", response.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

func (sc SlackClient) UploadFile(input SlackFilesUploadInput) (SlackFilesCompleteUploadExternal, error) {
	fileBytes, err := os.ReadFile(input.Path)
	if err != nil {
		return SlackFilesCompleteUploadExternal{}, err
	}

	filename := filepath.Base(input.Path)
	getUpload, err := sc.getUploadURLExternal(filename, len(fileBytes))
	if err != nil {
		return SlackFilesCompleteUploadExternal{}, err
	}

	if err := sc.uploadFileBytes(getUpload.UploadURL, fileBytes); err != nil {
		return SlackFilesCompleteUploadExternal{}, err
	}

	return sc.completeUploadExternal(getUpload.FileID, filename, input)
}

func (sc SlackClient) getUploadURLExternal(filename string, length int) (SlackFilesGetUploadURLExternal, error) {
	responseBody, err := sc.postMultipart("/files.getUploadURLExternal", map[string]string{
		"filename": filename,
		"length":   fmt.Sprintf("%d", length),
	})
	if err != nil {
		return SlackFilesGetUploadURLExternal{}, err
	}

	var response SlackFilesGetUploadURLExternal
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return SlackFilesGetUploadURLExternal{}, err
	}

	return response, nil
}

func (sc SlackClient) uploadFileBytes(uploadURL string, fileBytes []byte) error {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		uploadURL,
		bytes.NewReader(fileBytes),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/octet-stream")

	response, err := sc.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("file upload failed with status %d: %s", response.StatusCode, string(body))
	}

	return nil
}

func (sc SlackClient) completeUploadExternal(fileID string, filename string, input SlackFilesUploadInput) (SlackFilesCompleteUploadExternal, error) {
	filesValue, err := json.Marshal([]map[string]string{{
		"id":    fileID,
		"title": filename,
	}})
	if err != nil {
		return SlackFilesCompleteUploadExternal{}, err
	}

	values := map[string]string{
		"files":      string(filesValue),
		"channel_id": input.Channel,
	}

	if input.InitialComment != "" {
		values["initial_comment"] = input.InitialComment
	}

	if input.ThreadTS != "" {
		values["thread_ts"] = input.ThreadTS
	}

	responseBody, err := sc.postMultipart("/files.completeUploadExternal", values)
	if err != nil {
		return SlackFilesCompleteUploadExternal{}, err
	}

	var response SlackFilesCompleteUploadExternal
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return SlackFilesCompleteUploadExternal{}, err
	}

	return response, nil
}

func (sc SlackClient) DeleteFile(fileID string) error {
	_, err := sc.postMultipart("/files.delete", map[string]string{
		"file": fileID,
	})
	return err
}
