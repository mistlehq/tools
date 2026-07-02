package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mistlehq/tools/internal/argparse"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultMCPAddr     = "127.0.0.1:7357"
	defaultMCPEndpoint = "/mcp"
)

type telegramMCPConfig struct {
	Addr     string
	Endpoint string
}

type telegramEmptyToolInput struct{}

type telegramChatToolInput struct {
	Chat string `json:"chat" jsonschema:"Telegram chat ID or username, such as 123456789 or @channelusername."`
}

type telegramSendMessageToolInput struct {
	Chat            string `json:"chat" jsonschema:"Telegram chat ID or username, such as 123456789 or @channelusername."`
	Text            string `json:"text" jsonschema:"Message text."`
	MessageThreadID int    `json:"message_thread_id,omitempty" jsonschema:"Optional Telegram forum topic thread ID."`
	ParseMode       string `json:"parse_mode,omitempty" jsonschema:"Optional Telegram parse mode, such as HTML or MarkdownV2."`
}

type telegramEditMessageToolInput struct {
	Chat      string `json:"chat" jsonschema:"Telegram chat ID or username, such as 123456789 or @channelusername."`
	MessageID int    `json:"message_id" jsonschema:"Telegram message ID."`
	Text      string `json:"text" jsonschema:"Replacement message text."`
	ParseMode string `json:"parse_mode,omitempty" jsonschema:"Optional Telegram parse mode, such as HTML or MarkdownV2."`
}

type telegramMessageToolInput struct {
	Chat      string `json:"chat" jsonschema:"Telegram chat ID or username, such as 123456789 or @channelusername."`
	MessageID int    `json:"message_id" jsonschema:"Telegram message ID."`
}

type telegramDeleteMessagesToolInput struct {
	Chat       string `json:"chat" jsonschema:"Telegram chat ID or username, such as 123456789 or @channelusername."`
	MessageIDs []int  `json:"message_ids" jsonschema:"Telegram message IDs."`
}

type telegramSetReactionToolInput struct {
	Chat          string `json:"chat" jsonschema:"Telegram chat ID or username, such as 123456789 or @channelusername."`
	MessageID     int    `json:"message_id" jsonschema:"Telegram message ID."`
	Emoji         string `json:"emoji,omitempty" jsonschema:"Comma-separated standard emoji reactions."`
	CustomEmojiID string `json:"custom_emoji_id,omitempty" jsonschema:"Comma-separated custom emoji IDs."`
	IsBig         bool   `json:"is_big,omitempty" jsonschema:"Whether to set a big reaction."`
}

type telegramDeleteReactionToolInput struct {
	Chat        string `json:"chat" jsonschema:"Telegram chat ID or username, such as 123456789 or @channelusername."`
	MessageID   int    `json:"message_id" jsonschema:"Telegram message ID."`
	UserID      string `json:"user_id,omitempty" jsonschema:"Identifier of the user whose reaction will be removed."`
	ActorChatID string `json:"actor_chat_id,omitempty" jsonschema:"Identifier of the chat whose reaction will be removed."`
}

type telegramDeleteAllReactionsToolInput struct {
	Chat        string `json:"chat" jsonschema:"Telegram chat ID or username, such as 123456789 or @channelusername."`
	UserID      string `json:"user_id,omitempty" jsonschema:"Identifier of the user whose reactions will be removed."`
	ActorChatID string `json:"actor_chat_id,omitempty" jsonschema:"Identifier of the chat whose reactions will be removed."`
}

type telegramCreateTopicToolInput struct {
	Chat              string `json:"chat" jsonschema:"Telegram forum supergroup chat ID or username."`
	Name              string `json:"name" jsonschema:"Forum topic name."`
	IconColor         int    `json:"icon_color,omitempty" jsonschema:"Optional Telegram-supported RGB icon color integer."`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty" jsonschema:"Optional custom emoji ID for the topic icon."`
}

type telegramTopicToolInput struct {
	Chat            string `json:"chat" jsonschema:"Telegram forum supergroup chat ID or username."`
	MessageThreadID int    `json:"message_thread_id" jsonschema:"Telegram forum topic message thread ID."`
}

type telegramRequestToolInput struct {
	Method string         `json:"method" jsonschema:"Telegram Bot API method name, such as getMe or sendMessage."`
	Body   map[string]any `json:"body,omitempty" jsonschema:"Telegram Bot API method request body."`
}

type telegramRequestToolOutput struct {
	Result any `json:"result"`
}

type requiredToolField struct {
	Name  string
	Value string
}

type requiredIntToolField struct {
	Name  string
	Value int
}

type requiredIntSliceToolField struct {
	Name  string
	Value []int
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
	config, err := parseTelegramMCPServeArgs(args)
	if err != nil {
		return err
	}

	tc, err := cli.telegramClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newTelegramMCPHTTPHandler(tc))

	fmt.Fprintf(cli.stdout, "Telegram MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseTelegramMCPServeArgs(args []string) (telegramMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return telegramMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return telegramMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}

	config := telegramMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return telegramMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return telegramMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newTelegramMCPHTTPHandler(tc TelegramClient) http.Handler {
	server := newTelegramMCPServer(tc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newTelegramMCPServer(tc TelegramClient) *mcp.Server {
	openWorld := true
	notDestructive := false
	destructive := true

	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	mutatingAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &notDestructive}
	destructiveAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &destructive}

	server := mcp.NewServer(&mcp.Implementation{Name: "telegram", Version: Version}, nil)

	mcp.AddTool(server, telegramTool("telegram_auth_test", telegramAuthTestDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *telegramEmptyToolInput) (*mcp.CallToolResult, TelegramUser, error) {
		out, err := tc.AuthTestContext(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_chats_get", telegramChatsGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramChatToolInput) (*mcp.CallToolResult, TelegramChat, error) {
		if err := requireField("chat", input.Chat); err != nil {
			return nil, TelegramChat{}, err
		}
		out, err := tc.GetChatContext(ctx, input.Chat)
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_messages_send", telegramMessagesSendDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramSendMessageToolInput) (*mcp.CallToolResult, TelegramMessage, error) {
		if err := requireFields(field("chat", input.Chat), field("text", input.Text)); err != nil {
			return nil, TelegramMessage{}, err
		}
		out, err := tc.SendMessageContext(ctx, TelegramSendMessageInput{Chat: input.Chat, Text: input.Text, ParseMode: input.ParseMode, ThreadID: input.MessageThreadID})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_messages_edit", telegramMessagesEditDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramEditMessageToolInput) (*mcp.CallToolResult, TelegramMessage, error) {
		if err := requireFields(field("chat", input.Chat), positiveIntField("message_id", input.MessageID), field("text", input.Text)); err != nil {
			return nil, TelegramMessage{}, err
		}
		out, err := tc.EditMessageTextContext(ctx, TelegramEditMessageInput{Chat: input.Chat, MessageID: input.MessageID, Text: input.Text, ParseMode: input.ParseMode})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_messages_delete", telegramMessagesDeleteDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramMessageToolInput) (*mcp.CallToolResult, TelegramBoolResponse, error) {
		if err := requireFields(field("chat", input.Chat), positiveIntField("message_id", input.MessageID)); err != nil {
			return nil, TelegramBoolResponse{}, err
		}
		out, err := tc.DeleteMessageContext(ctx, TelegramMessageInput{Chat: input.Chat, MessageID: input.MessageID})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_messages_delete_batch", telegramMessagesDeleteBatchDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramDeleteMessagesToolInput) (*mcp.CallToolResult, TelegramBoolResponse, error) {
		if err := requireFields(field("chat", input.Chat), intSliceField("message_ids", input.MessageIDs)); err != nil {
			return nil, TelegramBoolResponse{}, err
		}
		out, err := tc.DeleteMessagesContext(ctx, TelegramDeleteMessagesInput{Chat: input.Chat, MessageIDs: input.MessageIDs})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_reactions_set", telegramReactionsSetDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramSetReactionToolInput) (*mcp.CallToolResult, TelegramBoolResponse, error) {
		if err := requireFields(field("chat", input.Chat), positiveIntField("message_id", input.MessageID)); err != nil {
			return nil, TelegramBoolResponse{}, err
		}
		reactions, err := parseReactionTypes(input.Emoji, input.CustomEmojiID)
		if err != nil {
			return nil, TelegramBoolResponse{}, err
		}
		out, err := tc.SetMessageReactionContext(ctx, TelegramSetReactionInput{Chat: input.Chat, MessageID: input.MessageID, Reactions: reactions, IsBig: input.IsBig})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_reactions_clear", telegramReactionsClearDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramMessageToolInput) (*mcp.CallToolResult, TelegramBoolResponse, error) {
		if err := requireFields(field("chat", input.Chat), positiveIntField("message_id", input.MessageID)); err != nil {
			return nil, TelegramBoolResponse{}, err
		}
		out, err := tc.SetMessageReactionContext(ctx, TelegramSetReactionInput{Chat: input.Chat, MessageID: input.MessageID})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_reactions_delete", telegramReactionsDeleteDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramDeleteReactionToolInput) (*mcp.CallToolResult, TelegramBoolResponse, error) {
		if err := requireFields(field("chat", input.Chat), positiveIntField("message_id", input.MessageID)); err != nil {
			return nil, TelegramBoolResponse{}, err
		}
		out, err := tc.DeleteMessageReactionContext(ctx, TelegramDeleteReactionInput{Chat: input.Chat, MessageID: input.MessageID, UserID: input.UserID, ActorChatID: input.ActorChatID})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_reactions_delete_all", telegramReactionsDeleteAllDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramDeleteAllReactionsToolInput) (*mcp.CallToolResult, TelegramBoolResponse, error) {
		if err := requireField("chat", input.Chat); err != nil {
			return nil, TelegramBoolResponse{}, err
		}
		out, err := tc.DeleteAllMessageReactionsContext(ctx, TelegramDeleteAllReactionsInput{Chat: input.Chat, UserID: input.UserID, ActorChatID: input.ActorChatID})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_topics_create", telegramTopicsCreateDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramCreateTopicToolInput) (*mcp.CallToolResult, TelegramForumTopic, error) {
		if err := requireFields(field("chat", input.Chat), field("name", input.Name)); err != nil {
			return nil, TelegramForumTopic{}, err
		}
		out, err := tc.CreateForumTopicContext(ctx, TelegramCreateTopicInput{Chat: input.Chat, Name: input.Name, IconColor: input.IconColor, IconCustomEmojiID: input.IconCustomEmojiID})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_topics_delete", telegramTopicsDeleteDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramTopicToolInput) (*mcp.CallToolResult, TelegramBoolResponse, error) {
		if err := requireFields(field("chat", input.Chat), positiveIntField("message_thread_id", input.MessageThreadID)); err != nil {
			return nil, TelegramBoolResponse{}, err
		}
		out, err := tc.DeleteForumTopicContext(ctx, TelegramTopicInput{Chat: input.Chat, ThreadID: input.MessageThreadID})
		return nil, out, err
	})
	mcp.AddTool(server, telegramTool("telegram_request", telegramRequestDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *telegramRequestToolInput) (*mcp.CallToolResult, telegramRequestToolOutput, error) {
		if err := requireField("method", input.Method); err != nil {
			return nil, telegramRequestToolOutput{}, err
		}
		body, err := jsonBodyFromMap(input.Body)
		if err != nil {
			return nil, telegramRequestToolOutput{}, err
		}
		raw, err := tc.RequestContext(ctx, input.Method, body)
		if err != nil {
			return nil, telegramRequestToolOutput{}, err
		}
		var out any
		if err := json.Unmarshal(raw, &out); err != nil {
			return nil, telegramRequestToolOutput{}, err
		}
		return nil, telegramRequestToolOutput{Result: out}, nil
	})

	return server
}

func telegramTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
	return &mcp.Tool{Name: name, Title: doc.Command, Description: doc.Description, Annotations: annotations}
}

func requireField(name string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

func field(name string, value string) requiredToolField {
	return requiredToolField{Name: name, Value: value}
}

func positiveIntField(name string, value int) requiredIntToolField {
	return requiredIntToolField{Name: name, Value: value}
}

func intSliceField(name string, value []int) requiredIntSliceToolField {
	return requiredIntSliceToolField{Name: name, Value: value}
}

func requireFields(fields ...any) error {
	for _, rawField := range fields {
		switch field := rawField.(type) {
		case requiredToolField:
			if err := requireField(field.Name, field.Value); err != nil {
				return err
			}
		case requiredIntToolField:
			if field.Value <= 0 {
				return fmt.Errorf("%s is required", field.Name)
			}
		case requiredIntSliceToolField:
			if len(field.Value) == 0 {
				return fmt.Errorf("%s is required", field.Name)
			}
			for _, value := range field.Value {
				if value <= 0 {
					return fmt.Errorf("%s must contain positive integers", field.Name)
				}
			}
		default:
			return fmt.Errorf("unsupported required field type")
		}
	}
	return nil
}

func jsonBodyFromMap(body map[string]any) (json.RawMessage, error) {
	if body == nil {
		return nil, nil
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(encoded), nil
}

func (cli CLI) printMCPHelp() {
	fmt.Fprint(cli.stdout, `telegram mcp

Run Telegram as a local MCP server.

Usage:
  telegram mcp help
  telegram mcp serve
  telegram mcp serve --help

Commands:
  serve    Serve Telegram MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `telegram mcp serve

Serve Telegram tools over MCP Streamable HTTP.

Usage:
  telegram mcp serve
  telegram mcp serve --addr <addr>
  telegram mcp serve --endpoint <path>
  telegram mcp serve --addr <addr> --endpoint <path>
  telegram mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.

Tools:
  telegram_auth_test                  %s
  telegram_chats_get                  %s
  telegram_messages_send              %s
  telegram_messages_edit              %s
  telegram_messages_delete            %s
  telegram_messages_delete_batch      %s
  telegram_reactions_set              %s
  telegram_reactions_clear            %s
  telegram_reactions_delete           %s
  telegram_reactions_delete_all       %s
  telegram_topics_create              %s
  telegram_topics_delete              %s
  telegram_request                    %s
`, defaultMCPAddr, defaultMCPEndpoint, telegramAuthTestDoc.Summary, telegramChatsGetDoc.Summary, telegramMessagesSendDoc.Summary, telegramMessagesEditDoc.Summary, telegramMessagesDeleteDoc.Summary, telegramMessagesDeleteBatchDoc.Summary, telegramReactionsSetDoc.Summary, telegramReactionsClearDoc.Summary, telegramReactionsDeleteDoc.Summary, telegramReactionsDeleteAllDoc.Summary, telegramTopicsCreateDoc.Summary, telegramTopicsDeleteDoc.Summary, telegramRequestDoc.Summary)
}
