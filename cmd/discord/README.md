# DISCORD

CLI for Discord.

## Overview

`discord` is a standalone command-line interface for Discord's REST API.

The CLI is designed for shell, script, and agent-driven workflows. It runs behind an auth-injecting proxy in Mistle, so the binary only needs a Discord API base URL and never receives or stores bot tokens directly.

The current implementation covers:

- bot authentication checks
- guild discovery
- channel discovery
- channel message listing
- message sending
- message editing
- message deletion
- message reactions
- role listing and basic role management
- member listing and basic member role/moderation operations
- local MCP serving for supported Discord tools

## Usage

### Commands

The supported commands are:

- `discord help`
- `discord version`
- `discord auth help`
- `discord auth test`
- `discord guilds help`
- `discord guilds list`
- `discord guilds get --guild <guild-id>`
- `discord channels help`
- `discord channels list --guild <guild-id>`
- `discord channels get --channel <channel-id>`
- `discord messages help`
- `discord messages list --channel <channel-id> [--limit <n>] [--before <message-id>] [--after <message-id>] [--around <message-id>]`
- `discord messages send --channel <channel-id> --content <text>`
- `discord messages edit --channel <channel-id> --message <message-id> --content <text>`
- `discord messages delete --channel <channel-id> --message <message-id>`
- `discord reactions help`
- `discord reactions add --channel <channel-id> --message <message-id> --emoji <emoji>`
- `discord reactions remove --channel <channel-id> --message <message-id> --emoji <emoji>`
- `discord roles help`
- `discord roles list --guild <guild-id>`
- `discord roles create --guild <guild-id> --name <name>`
- `discord roles delete --guild <guild-id> --role <role-id>`
- `discord members help`
- `discord members list --guild <guild-id> [--limit <n>] [--after <user-id>]`
- `discord members get --guild <guild-id> --user <user-id>`
- `discord members add-role --guild <guild-id> --user <user-id> --role <role-id>`
- `discord members remove-role --guild <guild-id> --user <user-id> --role <role-id>`
- `discord members ban --guild <guild-id> --user <user-id>`
- `discord members unban --guild <guild-id> --user <user-id>`
- `discord mcp help`
- `discord mcp serve`
- `discord mcp serve --addr <addr>`
- `discord mcp serve --endpoint <path>`

All API commands also accept `--json` for compact JSON output.

### MCP

`discord mcp serve` runs Discord as a local MCP server over Streamable HTTP. It reuses the same `DISCORD_BASE_URL` configuration as the CLI and relies on the same upstream auth-injecting proxy model.

By default, the server listens on `127.0.0.1:7356` and serves MCP at `/mcp`:

```sh
discord mcp serve
```

The MCP tools mirror the provider-backed CLI command surface with structured inputs and outputs:

| Tool | Purpose | Annotation |
| --- | --- | --- |
| `discord_auth_test` | Check Discord bot authentication state. | Read-only |
| `discord_guilds_list` | List visible Discord guilds. | Read-only |
| `discord_guilds_get` | Show details for a Discord guild. | Read-only |
| `discord_channels_list` | List Discord guild channels. | Read-only |
| `discord_channels_get` | Show details for a Discord channel. | Read-only |
| `discord_messages_list` | Fetch Discord channel messages. | Read-only |
| `discord_messages_send` | Send a Discord message. | Mutating, non-destructive |
| `discord_messages_edit` | Edit a Discord message. | Mutating, non-destructive |
| `discord_messages_delete` | Delete a Discord message. | Destructive |
| `discord_reactions_add` | Add a Discord message reaction. | Mutating, non-destructive |
| `discord_reactions_remove` | Remove the bot's Discord message reaction. | Destructive |
| `discord_roles_list` | List Discord guild roles. | Read-only |
| `discord_roles_create` | Create a Discord role. | Mutating, non-destructive |
| `discord_roles_delete` | Delete a Discord role. | Destructive |
| `discord_members_list` | List Discord guild members. | Read-only |
| `discord_members_get` | Show details for a Discord guild member. | Read-only |
| `discord_members_add_role` | Assign a role to a member. | Mutating, non-destructive |
| `discord_members_remove_role` | Remove a role from a member. | Destructive |
| `discord_members_ban` | Ban a Discord guild member. | Destructive |
| `discord_members_unban` | Unban a Discord user from a guild. | Destructive |

## Configuration

`discord` uses one configuration shape:

```json
{
  "baseUrl": "https://discord.com/api/v10"
}
```

Configuration can be supplied with:

- `DISCORD_BASE_URL=https://discord.com/api/v10`

`DISCORD_BASE_URL` must not end with a trailing slash.

When running in Mistle, `DISCORD_BASE_URL` points at an auth-injecting proxy that adds `Authorization: Bot <token>` upstream. The CLI itself does not accept Discord bot tokens.

## Build

```sh
mkdir -p dist && mise exec -- go build -o dist/discord ./cmd/discord
```
