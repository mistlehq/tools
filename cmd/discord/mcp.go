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
	defaultMCPAddr     = "127.0.0.1:7356"
	defaultMCPEndpoint = "/mcp"
)

type discordMCPConfig struct {
	Addr     string
	Endpoint string
}

type discordEmptyToolInput struct{}

type discordGuildToolInput struct {
	Guild string `json:"guild" jsonschema:"Discord guild ID."`
}

type discordChannelToolInput struct {
	Channel string `json:"channel" jsonschema:"Discord channel ID."`
}

type discordMessagesListToolInput struct {
	Channel string `json:"channel" jsonschema:"Discord channel ID."`
	Limit   string `json:"limit,omitempty" jsonschema:"Discord API page size."`
	Before  string `json:"before,omitempty" jsonschema:"Message ID before which to fetch messages."`
	After   string `json:"after,omitempty" jsonschema:"Message ID after which to fetch messages."`
	Around  string `json:"around,omitempty" jsonschema:"Message ID around which to fetch messages."`
}

type discordMessageContentToolInput struct {
	Channel string `json:"channel" jsonschema:"Discord channel ID."`
	Content string `json:"content" jsonschema:"Message content."`
}

type discordMessageUpdateToolInput struct {
	Channel string `json:"channel" jsonschema:"Discord channel ID."`
	Message string `json:"message" jsonschema:"Discord message ID."`
	Content string `json:"content" jsonschema:"Replacement message content."`
}

type discordMessageToolInput struct {
	Channel string `json:"channel" jsonschema:"Discord channel ID."`
	Message string `json:"message" jsonschema:"Discord message ID."`
}

type discordReactionToolInput struct {
	Channel string `json:"channel" jsonschema:"Discord channel ID."`
	Message string `json:"message" jsonschema:"Discord message ID."`
	Emoji   string `json:"emoji" jsonschema:"Unicode emoji or custom emoji formatted as name:id."`
}

type discordRoleCreateToolInput struct {
	Guild string `json:"guild" jsonschema:"Discord guild ID."`
	Name  string `json:"name" jsonschema:"Role name."`
}

type discordRoleToolInput struct {
	Guild string `json:"guild" jsonschema:"Discord guild ID."`
	Role  string `json:"role" jsonschema:"Discord role ID."`
}

type discordMembersListToolInput struct {
	Guild string `json:"guild" jsonschema:"Discord guild ID."`
	Limit string `json:"limit,omitempty" jsonschema:"Discord API page size."`
	After string `json:"after,omitempty" jsonschema:"User ID after which to fetch members."`
}

type discordMemberToolInput struct {
	Guild string `json:"guild" jsonschema:"Discord guild ID."`
	User  string `json:"user" jsonschema:"Discord user ID."`
}

type discordMemberRoleToolInput struct {
	Guild string `json:"guild" jsonschema:"Discord guild ID."`
	User  string `json:"user" jsonschema:"Discord user ID."`
	Role  string `json:"role" jsonschema:"Discord role ID."`
}

type discordGuildsListOutput struct {
	Guilds []DiscordGuild `json:"guilds"`
}

type discordChannelsListOutput struct {
	Channels []DiscordChannel `json:"channels"`
}

type discordMessagesListOutput struct {
	Messages []DiscordMessage `json:"messages"`
}

type discordRolesListOutput struct {
	Roles []DiscordRole `json:"roles"`
}

type discordMembersListOutput struct {
	Members []DiscordMember `json:"members"`
}

type requiredToolField struct {
	Name  string
	Value string
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
	config, err := parseDiscordMCPServeArgs(args)
	if err != nil {
		return err
	}

	dc, err := cli.discordClient()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, newDiscordMCPHTTPHandler(dc))

	fmt.Fprintf(cli.stdout, "Discord MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseDiscordMCPServeArgs(args []string) (discordMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return discordMCPConfig{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return discordMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}

	config := discordMCPConfig{Addr: defaultMCPAddr, Endpoint: defaultMCPEndpoint}
	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}
	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}
	if strings.TrimSpace(config.Addr) == "" {
		return discordMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}
	if !strings.HasPrefix(config.Endpoint, "/") {
		return discordMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}
	return config, nil
}

func newDiscordMCPHTTPHandler(dc DiscordClient) http.Handler {
	server := newDiscordMCPServer(dc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newDiscordMCPServer(dc DiscordClient) *mcp.Server {
	openWorld := true
	destructive := true
	notDestructive := false

	readOnlyAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &openWorld}
	mutatingAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &notDestructive}
	destructiveAnnotations := &mcp.ToolAnnotations{OpenWorldHint: &openWorld, DestructiveHint: &destructive}

	server := mcp.NewServer(&mcp.Implementation{Name: "discord", Version: Version}, nil)

	mcp.AddTool(server, discordTool("discord_auth_test", discordAuthTestDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *discordEmptyToolInput) (*mcp.CallToolResult, DiscordUser, error) {
		out, err := dc.AuthTestContext(ctx)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_guilds_list", discordGuildsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, _ *discordEmptyToolInput) (*mcp.CallToolResult, discordGuildsListOutput, error) {
		out, err := dc.ListGuildsContext(ctx)
		return nil, discordGuildsListOutput{Guilds: out}, err
	})
	mcp.AddTool(server, discordTool("discord_guilds_get", discordGuildsGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordGuildToolInput) (*mcp.CallToolResult, DiscordGuild, error) {
		if err := requireField("guild", input.Guild); err != nil {
			return nil, DiscordGuild{}, err
		}
		out, err := dc.GetGuildContext(ctx, input.Guild)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_channels_list", discordChannelsListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordGuildToolInput) (*mcp.CallToolResult, discordChannelsListOutput, error) {
		if err := requireField("guild", input.Guild); err != nil {
			return nil, discordChannelsListOutput{}, err
		}
		out, err := dc.ListChannelsContext(ctx, input.Guild)
		return nil, discordChannelsListOutput{Channels: out}, err
	})
	mcp.AddTool(server, discordTool("discord_channels_get", discordChannelsGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordChannelToolInput) (*mcp.CallToolResult, DiscordChannel, error) {
		if err := requireField("channel", input.Channel); err != nil {
			return nil, DiscordChannel{}, err
		}
		out, err := dc.GetChannelContext(ctx, input.Channel)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_messages_list", discordMessagesListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMessagesListToolInput) (*mcp.CallToolResult, discordMessagesListOutput, error) {
		if err := requireField("channel", input.Channel); err != nil {
			return nil, discordMessagesListOutput{}, err
		}
		out, err := dc.ListMessagesContext(ctx, DiscordListMessagesInput(*input))
		return nil, discordMessagesListOutput{Messages: out}, err
	})
	mcp.AddTool(server, discordTool("discord_messages_send", discordMessagesSendDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMessageContentToolInput) (*mcp.CallToolResult, DiscordMessage, error) {
		if err := requireField("channel", input.Channel); err != nil {
			return nil, DiscordMessage{}, err
		}
		if err := requireField("content", input.Content); err != nil {
			return nil, DiscordMessage{}, err
		}
		out, err := dc.CreateMessageContext(ctx, DiscordCreateMessageInput{Channel: input.Channel, Content: input.Content})
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_messages_edit", discordMessagesEditDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMessageUpdateToolInput) (*mcp.CallToolResult, DiscordMessage, error) {
		if err := requireFields(field("channel", input.Channel), field("message", input.Message), field("content", input.Content)); err != nil {
			return nil, DiscordMessage{}, err
		}
		out, err := dc.EditMessageContext(ctx, DiscordEditMessageInput{Channel: input.Channel, Message: input.Message, Content: input.Content})
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_messages_delete", discordMessagesDeleteDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMessageToolInput) (*mcp.CallToolResult, DiscordEmptyResponse, error) {
		if err := requireFields(field("channel", input.Channel), field("message", input.Message)); err != nil {
			return nil, DiscordEmptyResponse{}, err
		}
		out, err := dc.DeleteMessageContext(ctx, input.Channel, input.Message)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_reactions_add", discordReactionsAddDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordReactionToolInput) (*mcp.CallToolResult, DiscordEmptyResponse, error) {
		if err := requireFields(field("channel", input.Channel), field("message", input.Message), field("emoji", input.Emoji)); err != nil {
			return nil, DiscordEmptyResponse{}, err
		}
		out, err := dc.AddReactionContext(ctx, DiscordReactionInput(*input))
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_reactions_remove", discordReactionsRemoveDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordReactionToolInput) (*mcp.CallToolResult, DiscordEmptyResponse, error) {
		if err := requireFields(field("channel", input.Channel), field("message", input.Message), field("emoji", input.Emoji)); err != nil {
			return nil, DiscordEmptyResponse{}, err
		}
		out, err := dc.RemoveReactionContext(ctx, DiscordReactionInput(*input))
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_roles_list", discordRolesListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordGuildToolInput) (*mcp.CallToolResult, discordRolesListOutput, error) {
		if err := requireField("guild", input.Guild); err != nil {
			return nil, discordRolesListOutput{}, err
		}
		out, err := dc.ListRolesContext(ctx, input.Guild)
		return nil, discordRolesListOutput{Roles: out}, err
	})
	mcp.AddTool(server, discordTool("discord_roles_create", discordRolesCreateDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordRoleCreateToolInput) (*mcp.CallToolResult, DiscordRole, error) {
		if err := requireFields(field("guild", input.Guild), field("name", input.Name)); err != nil {
			return nil, DiscordRole{}, err
		}
		out, err := dc.CreateRoleContext(ctx, input.Guild, input.Name)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_roles_delete", discordRolesDeleteDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordRoleToolInput) (*mcp.CallToolResult, DiscordEmptyResponse, error) {
		if err := requireFields(field("guild", input.Guild), field("role", input.Role)); err != nil {
			return nil, DiscordEmptyResponse{}, err
		}
		out, err := dc.DeleteRoleContext(ctx, input.Guild, input.Role)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_members_list", discordMembersListDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMembersListToolInput) (*mcp.CallToolResult, discordMembersListOutput, error) {
		if err := requireField("guild", input.Guild); err != nil {
			return nil, discordMembersListOutput{}, err
		}
		out, err := dc.ListMembersContext(ctx, DiscordListMembersInput(*input))
		return nil, discordMembersListOutput{Members: out}, err
	})
	mcp.AddTool(server, discordTool("discord_members_get", discordMembersGetDoc, readOnlyAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMemberToolInput) (*mcp.CallToolResult, DiscordMember, error) {
		if err := requireFields(field("guild", input.Guild), field("user", input.User)); err != nil {
			return nil, DiscordMember{}, err
		}
		out, err := dc.GetMemberContext(ctx, input.Guild, input.User)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_members_add_role", discordMembersAddRoleDoc, mutatingAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMemberRoleToolInput) (*mcp.CallToolResult, DiscordEmptyResponse, error) {
		if err := requireFields(field("guild", input.Guild), field("user", input.User), field("role", input.Role)); err != nil {
			return nil, DiscordEmptyResponse{}, err
		}
		out, err := dc.AddMemberRoleContext(ctx, input.Guild, input.User, input.Role)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_members_remove_role", discordMembersRemoveRoleDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMemberRoleToolInput) (*mcp.CallToolResult, DiscordEmptyResponse, error) {
		if err := requireFields(field("guild", input.Guild), field("user", input.User), field("role", input.Role)); err != nil {
			return nil, DiscordEmptyResponse{}, err
		}
		out, err := dc.RemoveMemberRoleContext(ctx, input.Guild, input.User, input.Role)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_members_ban", discordMembersBanDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMemberToolInput) (*mcp.CallToolResult, DiscordEmptyResponse, error) {
		if err := requireFields(field("guild", input.Guild), field("user", input.User)); err != nil {
			return nil, DiscordEmptyResponse{}, err
		}
		out, err := dc.BanMemberContext(ctx, input.Guild, input.User)
		return nil, out, err
	})
	mcp.AddTool(server, discordTool("discord_members_unban", discordMembersUnbanDoc, destructiveAnnotations), func(ctx context.Context, _ *mcp.CallToolRequest, input *discordMemberToolInput) (*mcp.CallToolResult, DiscordEmptyResponse, error) {
		if err := requireFields(field("guild", input.Guild), field("user", input.User)); err != nil {
			return nil, DiscordEmptyResponse{}, err
		}
		out, err := dc.UnbanMemberContext(ctx, input.Guild, input.User)
		return nil, out, err
	})

	return server
}

func discordTool(name string, doc commandDoc, annotations *mcp.ToolAnnotations) *mcp.Tool {
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

func requireFields(fields ...requiredToolField) error {
	for _, field := range fields {
		if err := requireField(field.Name, field.Value); err != nil {
			return err
		}
	}
	return nil
}

func (cli CLI) printMCPHelp() {
	fmt.Fprint(cli.stdout, `discord mcp

Run Discord as a local MCP server.

Usage:
  discord mcp help
  discord mcp serve
  discord mcp serve --help

Commands:
  serve    Serve Discord MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `discord mcp serve

Serve Discord tools over MCP Streamable HTTP.

Usage:
  discord mcp serve
  discord mcp serve --addr <addr>
  discord mcp serve --endpoint <path>
  discord mcp serve --addr <addr> --endpoint <path>
  discord mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.

Tools:
  discord_auth_test             %s
  discord_guilds_list           %s
  discord_guilds_get            %s
  discord_channels_list         %s
  discord_channels_get          %s
  discord_messages_list         %s
  discord_messages_send         %s
  discord_messages_edit         %s
  discord_messages_delete       %s
  discord_reactions_add         %s
  discord_reactions_remove      %s
  discord_roles_list            %s
  discord_roles_create          %s
  discord_roles_delete          %s
  discord_members_list          %s
  discord_members_get           %s
  discord_members_add_role      %s
  discord_members_remove_role   %s
  discord_members_ban           %s
  discord_members_unban         %s
`, defaultMCPAddr, defaultMCPEndpoint, discordAuthTestDoc.Summary, discordGuildsListDoc.Summary, discordGuildsGetDoc.Summary, discordChannelsListDoc.Summary, discordChannelsGetDoc.Summary, discordMessagesListDoc.Summary, discordMessagesSendDoc.Summary, discordMessagesEditDoc.Summary, discordMessagesDeleteDoc.Summary, discordReactionsAddDoc.Summary, discordReactionsRemoveDoc.Summary, discordRolesListDoc.Summary, discordRolesCreateDoc.Summary, discordRolesDeleteDoc.Summary, discordMembersListDoc.Summary, discordMembersGetDoc.Summary, discordMembersAddRoleDoc.Summary, discordMembersRemoveRoleDoc.Summary, discordMembersBanDoc.Summary, discordMembersUnbanDoc.Summary)
}
