package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mistlehq/tools/internal/argparse"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultMCPAddr     = "127.0.0.1:7346"
	defaultMCPEndpoint = "/mcp"
)

type slackMCPConfig struct {
	Addr     string
	Endpoint string
}

type slackEmptyToolInput struct{}

type slackConversationsListToolInput struct {
	Types           string `json:"types,omitempty" jsonschema:"Comma-separated Slack conversation types, for example public_channel,private_channel."`
	Limit           string `json:"limit,omitempty" jsonschema:"Slack API page size."`
	Cursor          string `json:"cursor,omitempty" jsonschema:"Slack pagination cursor."`
	ExcludeArchived bool   `json:"excludeArchived,omitempty" jsonschema:"Exclude archived conversations."`
}

type slackConversationsInfoToolInput struct {
	Channel       string `json:"channel" jsonschema:"Slack conversation ID."`
	IncludeLocale bool   `json:"includeLocale,omitempty" jsonschema:"Include locale in the Slack API response."`
}

type slackConversationsHistoryToolInput struct {
	Channel   string `json:"channel" jsonschema:"Slack conversation ID."`
	Cursor    string `json:"cursor,omitempty" jsonschema:"Slack pagination cursor."`
	Inclusive bool   `json:"inclusive,omitempty" jsonschema:"Include messages matching latest or oldest bounds."`
	Latest    string `json:"latest,omitempty" jsonschema:"Latest message timestamp bound."`
	Limit     string `json:"limit,omitempty" jsonschema:"Slack API page size."`
	Oldest    string `json:"oldest,omitempty" jsonschema:"Oldest message timestamp bound."`
}

type slackConversationsRepliesToolInput struct {
	Channel   string `json:"channel" jsonschema:"Slack conversation ID."`
	TS        string `json:"ts" jsonschema:"Thread root message timestamp."`
	Cursor    string `json:"cursor,omitempty" jsonschema:"Slack pagination cursor."`
	Inclusive bool   `json:"inclusive,omitempty" jsonschema:"Include messages matching latest or oldest bounds."`
	Latest    string `json:"latest,omitempty" jsonschema:"Latest message timestamp bound."`
	Limit     string `json:"limit,omitempty" jsonschema:"Slack API page size."`
	Oldest    string `json:"oldest,omitempty" jsonschema:"Oldest message timestamp bound."`
}

type slackChatPostMessageToolInput struct {
	Channel  string  `json:"channel" jsonschema:"Slack conversation ID."`
	Text     string  `json:"text" jsonschema:"Message text."`
	ThreadTS *string `json:"threadTs,omitempty" jsonschema:"Optional thread root timestamp."`
}

type slackChatUpdateToolInput struct {
	Channel string `json:"channel" jsonschema:"Slack conversation ID."`
	TS      string `json:"ts" jsonschema:"Message timestamp to update."`
	Text    string `json:"text" jsonschema:"Replacement message text."`
}

type slackChatDeleteToolInput struct {
	Channel string `json:"channel" jsonschema:"Slack conversation ID."`
	TS      string `json:"ts" jsonschema:"Message timestamp to delete."`
}

type slackChatGetPermalinkToolInput struct {
	Channel   string `json:"channel" jsonschema:"Slack conversation ID."`
	MessageTS string `json:"messageTs" jsonschema:"Message timestamp."`
}

type slackReactionToolInput struct {
	Channel   string `json:"channel" jsonschema:"Slack conversation ID."`
	Timestamp string `json:"timestamp" jsonschema:"Message timestamp."`
	Name      string `json:"name" jsonschema:"Emoji reaction name without surrounding colons."`
}

type slackFilesInfoToolInput struct {
	FileID string `json:"fileId" jsonschema:"Slack file ID."`
}

type slackFilesDownloadToolInput struct {
	FileID string `json:"fileId" jsonschema:"Slack file ID."`
	Output string `json:"output" jsonschema:"Local output path where the file should be written."`
}

type slackFilesUploadToolInput struct {
	Path           string `json:"path" jsonschema:"Local path to upload."`
	Channel        string `json:"channel" jsonschema:"Slack conversation ID."`
	ThreadTS       string `json:"threadTs,omitempty" jsonschema:"Optional thread root timestamp."`
	InitialComment string `json:"initialComment,omitempty" jsonschema:"Optional initial comment for the upload."`
}

type slackEmojiListToolInput struct {
	IncludeCategories bool `json:"includeCategories,omitempty" jsonschema:"Include Slack emoji categories in the API response."`
}

func (cli CLI) runMCP(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printMCPHelp()
		return nil
	}

	switch args[0] {
	case "serve":
		if len(args) == 2 && (args[1] == "help" || args[1] == "-h" || args[1] == "--help") {
			cli.printMCPServeHelp()
			return nil
		}
		return cli.runMCPServe(args[1:])
	default:
		return fmt.Errorf("unsupported mcp command: %s", args[0])
	}
}

func (cli CLI) runMCPServe(args []string) error {
	config, err := parseSlackMCPServeArgs(args)
	if err != nil {
		return err
	}

	sc, err := cli.slackClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newSlackMCPHTTPHandler(sc))

	fmt.Fprintf(cli.stdout, "Slack MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseSlackMCPServeArgs(args []string) (slackMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return slackMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return slackMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}

	config := slackMCPConfig{
		Addr:     defaultMCPAddr,
		Endpoint: defaultMCPEndpoint,
	}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return slackMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return slackMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newSlackMCPHTTPHandler(sc SlackClient) http.Handler {
	server := newSlackMCPServer(sc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newSlackMCPServer(sc SlackClient) *mcp.Server {
	openWorld := true
	destructive := true
	notDestructive := false

	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	mutatingAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &notDestructive}
	destructiveAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &destructive}

	server := mcp.NewServer(&mcp.Implementation{Name: "slack", Version: Version}, nil)

	mcp.AddTool(server, slackTool("slack_auth_test", slackAuthTestDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *slackEmptyToolInput) (*mcp.CallToolResult, SlackAuthTest, error) {
		out, err := sc.AuthTestContext(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_conversations_list", slackConversationsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackConversationsListToolInput) (*mcp.CallToolResult, SlackConversationsList, error) {
		out, err := sc.ListConversationsContext(ctx, SlackConversationsListInput(*input))
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_conversations_info", slackConversationsInfoDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackConversationsInfoToolInput) (*mcp.CallToolResult, SlackConversationsInfo, error) {
		if strings.TrimSpace(input.Channel) == "" {
			return nil, SlackConversationsInfo{}, fmt.Errorf("channel is required")
		}
		out, err := sc.GetConversationInfoContext(ctx, SlackConversationsInfoInput(*input))
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_conversations_history", slackConversationsHistoryDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackConversationsHistoryToolInput) (*mcp.CallToolResult, SlackConversationsHistory, error) {
		if strings.TrimSpace(input.Channel) == "" {
			return nil, SlackConversationsHistory{}, fmt.Errorf("channel is required")
		}
		out, err := sc.GetConversationHistoryContext(ctx, SlackConversationsHistoryInput(*input))
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_conversations_replies", slackConversationsRepliesDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackConversationsRepliesToolInput) (*mcp.CallToolResult, SlackConversationsReplies, error) {
		if strings.TrimSpace(input.Channel) == "" {
			return nil, SlackConversationsReplies{}, fmt.Errorf("channel is required")
		}
		if strings.TrimSpace(input.TS) == "" {
			return nil, SlackConversationsReplies{}, fmt.Errorf("ts is required")
		}
		out, err := sc.GetConversationRepliesContext(ctx, SlackConversationsRepliesInput(*input))
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_chat_post_message", slackChatPostMessageDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackChatPostMessageToolInput) (*mcp.CallToolResult, SlackChatPostMessage, error) {
		if strings.TrimSpace(input.Channel) == "" {
			return nil, SlackChatPostMessage{}, fmt.Errorf("channel is required")
		}
		if strings.TrimSpace(input.Text) == "" {
			return nil, SlackChatPostMessage{}, fmt.Errorf("text is required and must not be empty")
		}
		out, err := sc.PostMessageContext(ctx, SlackChatPostMessageInput{Channel: input.Channel, Text: input.Text, ThreadTS: input.ThreadTS})
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_chat_update", slackChatUpdateDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackChatUpdateToolInput) (*mcp.CallToolResult, SlackChatUpdate, error) {
		if err := validateSlackChannelAndTS(input.Channel, input.TS); err != nil {
			return nil, SlackChatUpdate{}, err
		}
		if strings.TrimSpace(input.Text) == "" {
			return nil, SlackChatUpdate{}, fmt.Errorf("text is required and must not be empty")
		}
		out, err := sc.UpdateMessageContext(ctx, SlackChatUpdateInput(*input))
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_chat_delete", slackChatDeleteDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackChatDeleteToolInput) (*mcp.CallToolResult, SlackChatDelete, error) {
		if err := validateSlackChannelAndTS(input.Channel, input.TS); err != nil {
			return nil, SlackChatDelete{}, err
		}
		out, err := sc.DeleteMessageContext(ctx, SlackChatDeleteInput(*input))
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_chat_get_permalink", slackChatGetPermalinkDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackChatGetPermalinkToolInput) (*mcp.CallToolResult, SlackChatGetPermalink, error) {
		if strings.TrimSpace(input.Channel) == "" {
			return nil, SlackChatGetPermalink{}, fmt.Errorf("channel is required")
		}
		if strings.TrimSpace(input.MessageTS) == "" {
			return nil, SlackChatGetPermalink{}, fmt.Errorf("messageTs is required")
		}
		out, err := sc.GetPermalinkContext(ctx, input.Channel, input.MessageTS)
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_reactions_add", slackReactionsAddDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackReactionToolInput) (*mcp.CallToolResult, SlackReactionResponse, error) {
		reactionInput, err := buildSlackReactionInput(input)
		if err != nil {
			return nil, SlackReactionResponse{}, err
		}
		out, err := sc.AddReactionContext(ctx, reactionInput)
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_reactions_remove", slackReactionsRemoveDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackReactionToolInput) (*mcp.CallToolResult, SlackReactionResponse, error) {
		reactionInput, err := buildSlackReactionInput(input)
		if err != nil {
			return nil, SlackReactionResponse{}, err
		}
		out, err := sc.RemoveReactionContext(ctx, reactionInput)
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_files_info", slackFilesInfoDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackFilesInfoToolInput) (*mcp.CallToolResult, SlackFilesInfo, error) {
		if strings.TrimSpace(input.FileID) == "" {
			return nil, SlackFilesInfo{}, fmt.Errorf("fileId is required")
		}
		out, err := sc.GetFileInfoContext(ctx, input.FileID)
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_files_download", slackFilesDownloadDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackFilesDownloadToolInput) (*mcp.CallToolResult, SlackFilesDownload, error) {
		if strings.TrimSpace(input.FileID) == "" {
			return nil, SlackFilesDownload{}, fmt.Errorf("fileId is required")
		}
		if strings.TrimSpace(input.Output) == "" {
			return nil, SlackFilesDownload{}, fmt.Errorf("output is required")
		}
		out, err := sc.DownloadFileContext(ctx, SlackFilesDownloadInput{FileID: input.FileID, Output: input.Output})
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_files_upload", slackFilesUploadDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackFilesUploadToolInput) (*mcp.CallToolResult, SlackFilesCompleteUploadExternal, error) {
		if strings.TrimSpace(input.Path) == "" {
			return nil, SlackFilesCompleteUploadExternal{}, fmt.Errorf("path is required")
		}
		if strings.TrimSpace(input.Channel) == "" {
			return nil, SlackFilesCompleteUploadExternal{}, fmt.Errorf("channel is required")
		}
		out, err := sc.UploadFileContext(ctx, SlackFilesUploadInput(*input))
		return nil, out, err
	})
	mcp.AddTool(server, slackTool("slack_emoji_list", slackEmojiListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *slackEmojiListToolInput) (*mcp.CallToolResult, any, error) {
		out, err := sc.ListEmojiContext(ctx, input.IncludeCategories)
		return nil, out, err
	})

	return server
}

func slackTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
	return &mcp.Tool{Name: name, Title: doc.Command, Description: doc.Description, Annotations: annotations}
}

func validateSlackChannelAndTS(channel string, ts string) error {
	if strings.TrimSpace(channel) == "" {
		return fmt.Errorf("channel is required")
	}
	if strings.TrimSpace(ts) == "" {
		return fmt.Errorf("ts is required")
	}
	return nil
}

func buildSlackReactionInput(input *slackReactionToolInput) (SlackReactionInput, error) {
	if strings.TrimSpace(input.Channel) == "" {
		return SlackReactionInput{}, fmt.Errorf("channel is required")
	}
	if strings.TrimSpace(input.Timestamp) == "" {
		return SlackReactionInput{}, fmt.Errorf("timestamp is required")
	}
	if strings.TrimSpace(input.Name) == "" {
		return SlackReactionInput{}, fmt.Errorf("name is required")
	}
	return SlackReactionInput{Channel: input.Channel, Timestamp: input.Timestamp, Name: input.Name}, nil
}

func (cli CLI) printMCPHelp() {
	fmt.Fprint(cli.stdout, `slack mcp

Run Slack as a local MCP server.

Usage:
  slack mcp help
  slack mcp serve
  slack mcp serve --help

Commands:
  serve    Serve Slack MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `slack mcp serve

Serve Slack tools over MCP Streamable HTTP.

Usage:
  slack mcp serve
  slack mcp serve --addr <addr>
  slack mcp serve --endpoint <path>
  slack mcp serve --addr <addr> --endpoint <path>
  slack mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.

Tools:
  slack_auth_test               %s
  slack_conversations_list      %s
  slack_conversations_info      %s
  slack_conversations_history   %s
  slack_conversations_replies   %s
  slack_chat_post_message       %s
  slack_chat_update             %s
  slack_chat_delete             %s
  slack_chat_get_permalink      %s
  slack_reactions_add           %s
  slack_reactions_remove        %s
  slack_files_info              %s
  slack_files_download          %s
  slack_files_upload            %s
  slack_emoji_list              %s
`, defaultMCPAddr, defaultMCPEndpoint, slackAuthTestDoc.Summary, slackConversationsListDoc.Summary, slackConversationsInfoDoc.Summary, slackConversationsHistoryDoc.Summary, slackConversationsRepliesDoc.Summary, slackChatPostMessageDoc.Summary, slackChatUpdateDoc.Summary, slackChatDeleteDoc.Summary, slackChatGetPermalinkDoc.Summary, slackReactionsAddDoc.Summary, slackReactionsRemoveDoc.Summary, slackFilesInfoDoc.Summary, slackFilesDownloadDoc.Summary, slackFilesUploadDoc.Summary, slackEmojiListDoc.Summary)
}
