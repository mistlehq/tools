package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type DiscordClient struct {
	baseURL string
	client  *http.Client
}

type DiscordError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Errors  json.RawMessage `json:"errors,omitempty"`
}

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	GlobalName    string `json:"global_name"`
	Bot           bool   `json:"bot"`
}

type DiscordGuild struct {
	ID                          string `json:"id"`
	Name                        string `json:"name"`
	Icon                        string `json:"icon"`
	Owner                       bool   `json:"owner,omitempty"`
	Permissions                 string `json:"permissions,omitempty"`
	ApproximateMemberCount      int    `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    int    `json:"approximate_presence_count,omitempty"`
	Description                 string `json:"description,omitempty"`
	PreferredLocale             string `json:"preferred_locale,omitempty"`
	VerificationLevel           int    `json:"verification_level,omitempty"`
	DefaultMessageNotifications int    `json:"default_message_notifications,omitempty"`
}

type DiscordChannel struct {
	ID                   string                 `json:"id"`
	Type                 int                    `json:"type"`
	GuildID              string                 `json:"guild_id,omitempty"`
	Position             int                    `json:"position,omitempty"`
	Name                 string                 `json:"name,omitempty"`
	Topic                string                 `json:"topic,omitempty"`
	NSFW                 bool                   `json:"nsfw,omitempty"`
	ParentID             string                 `json:"parent_id,omitempty"`
	RateLimitPerUser     int                    `json:"rate_limit_per_user,omitempty"`
	PermissionOverwrites []map[string]any       `json:"permission_overwrites,omitempty"`
	ThreadMetadata       *DiscordThreadMetadata `json:"thread_metadata,omitempty"`
	AppliedTags          []string               `json:"applied_tags,omitempty"`
	AvailableTags        []DiscordForumTag      `json:"available_tags,omitempty"`
	LastMessageID        string                 `json:"last_message_id,omitempty"`
	LastPinTimestamp     string                 `json:"last_pin_timestamp,omitempty"`
}

type DiscordThreadMetadata struct {
	Archived            bool   `json:"archived"`
	AutoArchiveDuration int    `json:"auto_archive_duration"`
	ArchiveTimestamp    string `json:"archive_timestamp"`
	Locked              bool   `json:"locked"`
	Invitable           bool   `json:"invitable,omitempty"`
}

type DiscordForumTag struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Moderated bool   `json:"moderated"`
	EmojiID   string `json:"emoji_id,omitempty"`
	EmojiName string `json:"emoji_name,omitempty"`
}

type DiscordMessage struct {
	ID              string              `json:"id"`
	ChannelID       string              `json:"channel_id"`
	GuildID         string              `json:"guild_id,omitempty"`
	Author          DiscordUser         `json:"author"`
	Content         string              `json:"content"`
	Timestamp       string              `json:"timestamp"`
	EditedTimestamp string              `json:"edited_timestamp"`
	TTS             bool                `json:"tts"`
	MentionEveryone bool                `json:"mention_everyone"`
	Mentions        []DiscordUser       `json:"mentions"`
	MentionRoles    []string            `json:"mention_roles"`
	Attachments     []DiscordAttachment `json:"attachments"`
	Embeds          []map[string]any    `json:"embeds"`
	Reactions       []DiscordReaction   `json:"reactions"`
	Pinned          bool                `json:"pinned"`
	Type            int                 `json:"type"`
	Thread          *DiscordChannel     `json:"thread,omitempty"`
}

type DiscordAttachment struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Description string `json:"description,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Size        int64  `json:"size"`
	URL         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
}

type DiscordReaction struct {
	Count int          `json:"count"`
	Me    bool         `json:"me"`
	Emoji DiscordEmoji `json:"emoji"`
}

type DiscordEmoji struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

type DiscordRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions string `json:"permissions"`
	Managed     bool   `json:"managed"`
	Mentionable bool   `json:"mentionable"`
}

type DiscordMember struct {
	User                       DiscordUser `json:"user"`
	Nick                       string      `json:"nick,omitempty"`
	Avatar                     string      `json:"avatar,omitempty"`
	Roles                      []string    `json:"roles"`
	JoinedAt                   string      `json:"joined_at"`
	PremiumSince               string      `json:"premium_since,omitempty"`
	Deaf                       bool        `json:"deaf"`
	Mute                       bool        `json:"mute"`
	Pending                    bool        `json:"pending,omitempty"`
	Permissions                string      `json:"permissions,omitempty"`
	CommunicationDisabledUntil string      `json:"communication_disabled_until,omitempty"`
}

type DiscordEmptyResponse struct {
	OK bool `json:"ok"`
}

type DiscordListMessagesInput struct {
	Channel string
	Limit   string
	Before  string
	After   string
	Around  string
}

type DiscordCreateMessageInput struct {
	Channel string
	Content string
}

type DiscordEditMessageInput struct {
	Channel string
	Message string
	Content string
}

type DiscordReactionInput struct {
	Channel string
	Message string
	Emoji   string
}

type DiscordListMembersInput struct {
	Guild string
	Limit string
	After string
}

func NewDiscordClient(config Config) DiscordClient {
	return DiscordClient{
		baseURL: config.BaseURL,
		client:  http.DefaultClient,
	}
}

func (dc DiscordClient) AuthTestContext(ctx context.Context) (DiscordUser, error) {
	var out DiscordUser
	err := dc.getJSON(ctx, "/users/@me", nil, &out)
	return out, err
}

func (dc DiscordClient) ListGuildsContext(ctx context.Context) ([]DiscordGuild, error) {
	var out []DiscordGuild
	err := dc.getJSON(ctx, "/users/@me/guilds", nil, &out)
	return out, err
}

func (dc DiscordClient) GetGuildContext(ctx context.Context, guildID string) (DiscordGuild, error) {
	var out DiscordGuild
	err := dc.getJSON(ctx, "/guilds/"+url.PathEscape(guildID), nil, &out)
	return out, err
}

func (dc DiscordClient) ListChannelsContext(ctx context.Context, guildID string) ([]DiscordChannel, error) {
	var out []DiscordChannel
	err := dc.getJSON(ctx, "/guilds/"+url.PathEscape(guildID)+"/channels", nil, &out)
	return out, err
}

func (dc DiscordClient) GetChannelContext(ctx context.Context, channelID string) (DiscordChannel, error) {
	var out DiscordChannel
	err := dc.getJSON(ctx, "/channels/"+url.PathEscape(channelID), nil, &out)
	return out, err
}

func (dc DiscordClient) ListMessagesContext(ctx context.Context, input DiscordListMessagesInput) ([]DiscordMessage, error) {
	query := url.Values{}
	addQueryValue(query, "limit", input.Limit)
	addQueryValue(query, "before", input.Before)
	addQueryValue(query, "after", input.After)
	addQueryValue(query, "around", input.Around)

	var out []DiscordMessage
	err := dc.getJSON(ctx, "/channels/"+url.PathEscape(input.Channel)+"/messages", query, &out)
	return out, err
}

func (dc DiscordClient) CreateMessageContext(ctx context.Context, input DiscordCreateMessageInput) (DiscordMessage, error) {
	var out DiscordMessage
	err := dc.requestJSON(ctx, http.MethodPost, "/channels/"+url.PathEscape(input.Channel)+"/messages", map[string]string{
		"content": input.Content,
	}, &out)
	return out, err
}

func (dc DiscordClient) EditMessageContext(ctx context.Context, input DiscordEditMessageInput) (DiscordMessage, error) {
	var out DiscordMessage
	err := dc.requestJSON(ctx, http.MethodPatch, "/channels/"+url.PathEscape(input.Channel)+"/messages/"+url.PathEscape(input.Message), map[string]string{
		"content": input.Content,
	}, &out)
	return out, err
}

func (dc DiscordClient) DeleteMessageContext(ctx context.Context, channelID string, messageID string) (DiscordEmptyResponse, error) {
	err := dc.requestJSON(ctx, http.MethodDelete, "/channels/"+url.PathEscape(channelID)+"/messages/"+url.PathEscape(messageID), nil, nil)
	return DiscordEmptyResponse{OK: err == nil}, err
}

func (dc DiscordClient) AddReactionContext(ctx context.Context, input DiscordReactionInput) (DiscordEmptyResponse, error) {
	err := dc.requestJSON(ctx, http.MethodPut, "/channels/"+url.PathEscape(input.Channel)+"/messages/"+url.PathEscape(input.Message)+"/reactions/"+url.PathEscape(input.Emoji)+"/@me", nil, nil)
	return DiscordEmptyResponse{OK: err == nil}, err
}

func (dc DiscordClient) RemoveReactionContext(ctx context.Context, input DiscordReactionInput) (DiscordEmptyResponse, error) {
	err := dc.requestJSON(ctx, http.MethodDelete, "/channels/"+url.PathEscape(input.Channel)+"/messages/"+url.PathEscape(input.Message)+"/reactions/"+url.PathEscape(input.Emoji)+"/@me", nil, nil)
	return DiscordEmptyResponse{OK: err == nil}, err
}

func (dc DiscordClient) ListRolesContext(ctx context.Context, guildID string) ([]DiscordRole, error) {
	var out []DiscordRole
	err := dc.getJSON(ctx, "/guilds/"+url.PathEscape(guildID)+"/roles", nil, &out)
	return out, err
}

func (dc DiscordClient) CreateRoleContext(ctx context.Context, guildID string, name string) (DiscordRole, error) {
	var out DiscordRole
	err := dc.requestJSON(ctx, http.MethodPost, "/guilds/"+url.PathEscape(guildID)+"/roles", map[string]string{
		"name": name,
	}, &out)
	return out, err
}

func (dc DiscordClient) DeleteRoleContext(ctx context.Context, guildID string, roleID string) (DiscordEmptyResponse, error) {
	err := dc.requestJSON(ctx, http.MethodDelete, "/guilds/"+url.PathEscape(guildID)+"/roles/"+url.PathEscape(roleID), nil, nil)
	return DiscordEmptyResponse{OK: err == nil}, err
}

func (dc DiscordClient) ListMembersContext(ctx context.Context, input DiscordListMembersInput) ([]DiscordMember, error) {
	query := url.Values{}
	addQueryValue(query, "limit", input.Limit)
	addQueryValue(query, "after", input.After)

	var out []DiscordMember
	err := dc.getJSON(ctx, "/guilds/"+url.PathEscape(input.Guild)+"/members", query, &out)
	return out, err
}

func (dc DiscordClient) GetMemberContext(ctx context.Context, guildID string, userID string) (DiscordMember, error) {
	var out DiscordMember
	err := dc.getJSON(ctx, "/guilds/"+url.PathEscape(guildID)+"/members/"+url.PathEscape(userID), nil, &out)
	return out, err
}

func (dc DiscordClient) AddMemberRoleContext(ctx context.Context, guildID string, userID string, roleID string) (DiscordEmptyResponse, error) {
	err := dc.requestJSON(ctx, http.MethodPut, "/guilds/"+url.PathEscape(guildID)+"/members/"+url.PathEscape(userID)+"/roles/"+url.PathEscape(roleID), nil, nil)
	return DiscordEmptyResponse{OK: err == nil}, err
}

func (dc DiscordClient) RemoveMemberRoleContext(ctx context.Context, guildID string, userID string, roleID string) (DiscordEmptyResponse, error) {
	err := dc.requestJSON(ctx, http.MethodDelete, "/guilds/"+url.PathEscape(guildID)+"/members/"+url.PathEscape(userID)+"/roles/"+url.PathEscape(roleID), nil, nil)
	return DiscordEmptyResponse{OK: err == nil}, err
}

func (dc DiscordClient) BanMemberContext(ctx context.Context, guildID string, userID string) (DiscordEmptyResponse, error) {
	err := dc.requestJSON(ctx, http.MethodPut, "/guilds/"+url.PathEscape(guildID)+"/bans/"+url.PathEscape(userID), map[string]int{
		"delete_message_seconds": 0,
	}, nil)
	return DiscordEmptyResponse{OK: err == nil}, err
}

func (dc DiscordClient) UnbanMemberContext(ctx context.Context, guildID string, userID string) (DiscordEmptyResponse, error) {
	err := dc.requestJSON(ctx, http.MethodDelete, "/guilds/"+url.PathEscape(guildID)+"/bans/"+url.PathEscape(userID), nil, nil)
	return DiscordEmptyResponse{OK: err == nil}, err
}

func (dc DiscordClient) getJSON(ctx context.Context, path string, query url.Values, out any) error {
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	return dc.requestJSON(ctx, http.MethodGet, path, nil, out)
}

func (dc DiscordClient) requestJSON(ctx context.Context, method string, path string, body any, out any) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("method path must start with '/': %s", path)
	}

	var requestBody io.Reader
	if body != nil {
		encodedBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		requestBody = bytes.NewReader(encodedBody)
	}

	request, err := http.NewRequestWithContext(ctx, method, dc.baseURL+path, requestBody)
	if err != nil {
		return err
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := dc.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("discord api %s %s: rate limited%s", method, path, discordRetryAfterMessage(response, responseBody))
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("discord api %s %s failed with status %d: %s", method, path, response.StatusCode, discordErrorMessage(responseBody))
	}
	if out == nil || len(responseBody) == 0 {
		return nil
	}

	return json.Unmarshal(responseBody, out)
}

func addQueryValue(query url.Values, key string, value string) {
	if strings.TrimSpace(value) != "" {
		query.Set(key, value)
	}
}

func discordRetryAfterMessage(response *http.Response, responseBody []byte) string {
	if retryAfterHeader := response.Header.Get("Retry-After"); retryAfterHeader != "" {
		return ", retry after " + retryAfterHeader + " seconds"
	}

	var payload struct {
		RetryAfter float64 `json:"retry_after"`
	}
	if err := json.Unmarshal(responseBody, &payload); err == nil && payload.RetryAfter > 0 {
		return ", retry after " + strconv.FormatFloat(payload.RetryAfter, 'f', -1, 64) + " seconds"
	}

	return ""
}

func discordErrorMessage(responseBody []byte) string {
	var discordError DiscordError
	if err := json.Unmarshal(responseBody, &discordError); err == nil && discordError.Message != "" {
		if len(discordError.Errors) > 0 {
			return fmt.Sprintf("%s: %s", discordError.Message, string(discordError.Errors))
		}
		return discordError.Message
	}

	return string(responseBody)
}
