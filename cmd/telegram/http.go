package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TelegramClient struct {
	baseURL string
	client  *http.Client
}

type TelegramError struct {
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (err TelegramError) Error() string {
	if err.ErrorCode == 0 {
		return err.Description
	}
	return fmt.Sprintf("telegram api error %d: %s", err.ErrorCode, err.Description)
}

type TelegramUser struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type TelegramChat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type TelegramMessage struct {
	MessageID int           `json:"message_id"`
	Date      int64         `json:"date"`
	Chat      TelegramChat  `json:"chat"`
	From      *TelegramUser `json:"from,omitempty"`
	Text      string        `json:"text,omitempty"`
}

type TelegramForumTopic struct {
	MessageThreadID   int    `json:"message_thread_id"`
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

type TelegramSendMessageInput struct {
	Chat      string
	Text      string
	ParseMode string
	ThreadID  int
}

type TelegramEditMessageInput struct {
	Chat      string
	MessageID int
	Text      string
	ParseMode string
}

type TelegramMessageInput struct {
	Chat      string
	MessageID int
}

type TelegramDeleteMessagesInput struct {
	Chat       string
	MessageIDs []int
}

type TelegramSetReactionInput struct {
	Chat      string
	MessageID int
	Reactions []map[string]string
	IsBig     bool
}

type TelegramDeleteReactionInput struct {
	Chat        string
	MessageID   int
	UserID      string
	ActorChatID string
}

type TelegramDeleteAllReactionsInput struct {
	Chat        string
	UserID      string
	ActorChatID string
}

type TelegramCreateTopicInput struct {
	Chat              string
	Name              string
	IconColor         int
	IconCustomEmojiID string
}

type TelegramTopicInput struct {
	Chat     string
	ThreadID int
}

type TelegramBoolResponse struct {
	OK bool `json:"ok"`
}

type telegramResponse[T any] struct {
	OK          bool            `json:"ok"`
	Result      T               `json:"result"`
	ErrorCode   int             `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

func NewTelegramClient(config Config) TelegramClient {
	return TelegramClient{
		baseURL: config.BaseURL,
		client:  http.DefaultClient,
	}
}

func (tc TelegramClient) AuthTestContext(ctx context.Context) (TelegramUser, error) {
	return telegramPost[TelegramUser](ctx, tc, "/getMe", map[string]string{})
}

func (tc TelegramClient) GetChatContext(ctx context.Context, chat string) (TelegramChat, error) {
	return telegramPost[TelegramChat](ctx, tc, "/getChat", map[string]string{"chat_id": chat})
}

func (tc TelegramClient) CreateForumTopicContext(ctx context.Context, input TelegramCreateTopicInput) (TelegramForumTopic, error) {
	body := map[string]any{
		"chat_id": input.Chat,
		"name":    input.Name,
	}
	if input.IconColor > 0 {
		body["icon_color"] = input.IconColor
	}
	if strings.TrimSpace(input.IconCustomEmojiID) != "" {
		body["icon_custom_emoji_id"] = input.IconCustomEmojiID
	}
	return telegramPost[TelegramForumTopic](ctx, tc, "/createForumTopic", body)
}

func (tc TelegramClient) DeleteForumTopicContext(ctx context.Context, input TelegramTopicInput) (TelegramBoolResponse, error) {
	err := tc.requestJSON(ctx, http.MethodPost, "/deleteForumTopic", map[string]any{
		"chat_id":           input.Chat,
		"message_thread_id": input.ThreadID,
	}, nil)
	return TelegramBoolResponse{OK: err == nil}, err
}

func (tc TelegramClient) SendMessageContext(ctx context.Context, input TelegramSendMessageInput) (TelegramMessage, error) {
	body := map[string]any{
		"chat_id": input.Chat,
		"text":    input.Text,
	}
	if strings.TrimSpace(input.ParseMode) != "" {
		body["parse_mode"] = input.ParseMode
	}
	if input.ThreadID > 0 {
		body["message_thread_id"] = input.ThreadID
	}
	return telegramPost[TelegramMessage](ctx, tc, "/sendMessage", body)
}

func (tc TelegramClient) EditMessageTextContext(ctx context.Context, input TelegramEditMessageInput) (TelegramMessage, error) {
	body := map[string]any{
		"chat_id":    input.Chat,
		"message_id": input.MessageID,
		"text":       input.Text,
	}
	if strings.TrimSpace(input.ParseMode) != "" {
		body["parse_mode"] = input.ParseMode
	}
	return telegramPost[TelegramMessage](ctx, tc, "/editMessageText", body)
}

func (tc TelegramClient) DeleteMessageContext(ctx context.Context, input TelegramMessageInput) (TelegramBoolResponse, error) {
	err := tc.requestJSON(ctx, http.MethodPost, "/deleteMessage", map[string]any{
		"chat_id":    input.Chat,
		"message_id": input.MessageID,
	}, nil)
	return TelegramBoolResponse{OK: err == nil}, err
}

func (tc TelegramClient) DeleteMessagesContext(ctx context.Context, input TelegramDeleteMessagesInput) (TelegramBoolResponse, error) {
	err := tc.requestJSON(ctx, http.MethodPost, "/deleteMessages", map[string]any{
		"chat_id":     input.Chat,
		"message_ids": input.MessageIDs,
	}, nil)
	return TelegramBoolResponse{OK: err == nil}, err
}

func (tc TelegramClient) SetMessageReactionContext(ctx context.Context, input TelegramSetReactionInput) (TelegramBoolResponse, error) {
	reactions := input.Reactions
	if reactions == nil {
		reactions = []map[string]string{}
	}
	encodedBody, err := json.Marshal(map[string]any{
		"chat_id":    input.Chat,
		"message_id": input.MessageID,
		"reaction":   reactions,
	})
	if err != nil {
		return TelegramBoolResponse{}, err
	}
	if input.IsBig {
		encodedBody, err = json.Marshal(map[string]any{
			"chat_id":    input.Chat,
			"message_id": input.MessageID,
			"reaction":   reactions,
			"is_big":     true,
		})
		if err != nil {
			return TelegramBoolResponse{}, err
		}
	}
	_, err = tc.RequestContext(ctx, "setMessageReaction", json.RawMessage(encodedBody))
	return TelegramBoolResponse{OK: err == nil}, err
}

func (tc TelegramClient) DeleteMessageReactionContext(ctx context.Context, input TelegramDeleteReactionInput) (TelegramBoolResponse, error) {
	body := map[string]any{
		"chat_id":    input.Chat,
		"message_id": input.MessageID,
	}
	if strings.TrimSpace(input.UserID) != "" {
		body["user_id"] = input.UserID
	}
	if strings.TrimSpace(input.ActorChatID) != "" {
		body["actor_chat_id"] = input.ActorChatID
	}
	err := tc.requestJSON(ctx, http.MethodPost, "/deleteMessageReaction", body, nil)
	return TelegramBoolResponse{OK: err == nil}, err
}

func (tc TelegramClient) DeleteAllMessageReactionsContext(ctx context.Context, input TelegramDeleteAllReactionsInput) (TelegramBoolResponse, error) {
	body := map[string]any{
		"chat_id": input.Chat,
	}
	if strings.TrimSpace(input.UserID) != "" {
		body["user_id"] = input.UserID
	}
	if strings.TrimSpace(input.ActorChatID) != "" {
		body["actor_chat_id"] = input.ActorChatID
	}
	err := tc.requestJSON(ctx, http.MethodPost, "/deleteAllMessageReactions", body, nil)
	return TelegramBoolResponse{OK: err == nil}, err
}

func (tc TelegramClient) RequestContext(ctx context.Context, method string, body json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(method) == "" {
		return nil, fmt.Errorf("--method is required")
	}
	if strings.HasPrefix(method, "/") {
		return nil, fmt.Errorf("--method must be a Telegram method name, not a path")
	}

	requestBody := body
	if len(requestBody) == 0 {
		requestBody = json.RawMessage(`{}`)
	}

	var out json.RawMessage
	err := tc.requestJSON(ctx, http.MethodPost, "/"+method, requestBody, &out)
	return out, err
}

func telegramPost[T any](ctx context.Context, tc TelegramClient, path string, body any) (T, error) {
	var out T
	if err := tc.requestJSON(ctx, http.MethodPost, path, body, &out); err != nil {
		return out, err
	}
	return out, nil
}

func (tc TelegramClient) requestJSON(ctx context.Context, method string, path string, body any, out any) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("method path must start with '/': %s", path)
	}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, method, tc.baseURL+path, bytes.NewReader(encodedBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := tc.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("telegram api %s %s failed with status %d: %s", method, path, response.StatusCode, telegramHTTPErrorMessage(responseBody))
	}

	return decodeTelegramResponse(responseBody, out)
}

func decodeTelegramResponse(responseBody []byte, out any) error {
	var envelope telegramResponse[json.RawMessage]
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return err
	}
	if !envelope.OK {
		return TelegramError{ErrorCode: envelope.ErrorCode, Description: envelope.Description}
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(envelope.Result, out)
}

func telegramHTTPErrorMessage(responseBody []byte) string {
	var envelope telegramResponse[json.RawMessage]
	if err := json.Unmarshal(responseBody, &envelope); err == nil && envelope.Description != "" {
		return TelegramError{ErrorCode: envelope.ErrorCode, Description: envelope.Description}.Error()
	}

	return string(responseBody)
}
