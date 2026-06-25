# SLACK

CLI for Slack.

## Overview

`slack` is a standalone command-line interface for Slack.

The CLI mirrors Slack Web API method families and is designed for shell, script, and agent-driven workflows.

The current implementation covers:

- authentication checks
- conversation discovery
- conversation lookup
- conversation history
- message posting
- message updates
- message deletion
- permalink lookup
- message reactions
- emoji listing
- file inspection
- file download
- file upload
- local MCP serving for supported Slack tools

## Usage

### Commands

The supported commands are:

- `slack help`
- `slack version`
- `slack auth help`
- `slack auth test`
- `slack conversations help`
- `slack conversations list [--types <csv>] [--limit <n>] [--cursor <cursor>] [--exclude-archived]`
- `slack conversations info --channel <conversation-id> [--include-locale]`
- `slack conversations history --channel <conversation-id> [--cursor <cursor>] [--inclusive] [--latest <ts>] [--limit <n>] [--oldest <ts>]`
- `slack conversations replies --channel <conversation-id> --ts <thread-root-ts> [--cursor <cursor>] [--inclusive] [--latest <ts>] [--limit <n>] [--oldest <ts>]`
- `slack chat help`
- `slack chat post-message --channel <conversation-id> --text <text>`
- `slack chat post-message --channel <conversation-id> --text-file <path>`
- `slack chat post-message --channel <conversation-id> --blocks <json-array> [--text <fallback-text>]`
- `slack chat post-message --channel <conversation-id> --blocks-file <path> [--attachments-file <path>]`
- `slack chat post-message --channel <conversation-id> --thread-ts <ts> --text <text>`
- `slack chat update --channel <conversation-id> --ts <ts> --text <text>`
- `slack chat update --channel <conversation-id> --ts <ts> --text-file <path>`
- `slack chat update --channel <conversation-id> --ts <ts> --blocks <json-array> [--text <fallback-text>]`
- `slack chat update --channel <conversation-id> --ts <ts> --attachments-file <path>`
- `slack chat delete --channel <conversation-id> --ts <ts>`
- `slack chat get-permalink --channel <conversation-id> --message-ts <ts>`
- `slack reactions help`
- `slack reactions add --channel <conversation-id> --timestamp <ts> --name <emoji-name>`
- `slack reactions remove --channel <conversation-id> --timestamp <ts> --name <emoji-name>`
- `slack files help`
- `slack files info --file <file-id>`
- `slack files download --file <file-id> --output <path>`
- `slack files upload --path <path> --channel <conversation-id>`
- `slack files upload --path <path> --channel <conversation-id> --initial-comment <text>`
- `slack files upload --path <path> --channel <conversation-id> --initial-comment-file <path>`
- `slack files upload --path <path> --channel <conversation-id> --thread-ts <ts>`
- `slack emoji help`
- `slack emoji list`
- `slack emoji list --include-categories`
- `slack mcp help`
- `slack mcp serve`
- `slack mcp serve --addr <addr>`
- `slack mcp serve --endpoint <path>`

All API commands also accept `--json` for compact JSON output.

### Discovery

Use help to discover the currently supported command families:

```sh
slack help
slack auth help
slack conversations help
slack chat help
slack reactions help
slack files help
slack emoji help
slack mcp help
slack mcp serve --help
```

### MCP

`slack mcp serve` runs Slack as a local MCP server over Streamable HTTP. It reuses the same `SLACK_BASE_URL` configuration as the CLI and relies on the same upstream auth-injecting proxy model.

By default, the server listens on `127.0.0.1:7346` and serves MCP at `/mcp`:

```sh
slack mcp serve
```

Override the listen address or endpoint when needed:

```sh
slack mcp serve --addr 127.0.0.1:8080
slack mcp serve --endpoint /mcp
slack mcp serve --addr 127.0.0.1:8080 --endpoint /mcp
```

The MCP tools mirror the provider-backed CLI command surface with structured inputs and outputs:

| Tool | Purpose | Backing command/API | Annotation |
| --- | --- | --- | --- |
| `slack_auth_test` | Check Slack authentication state. | `slack auth test` / `POST /auth.test` | Read-only |
| `slack_conversations_list` | List Slack conversations. | `slack conversations list` / `GET /conversations.list` | Read-only |
| `slack_conversations_info` | Show details for a Slack conversation. | `slack conversations info` / `GET /conversations.info` | Read-only |
| `slack_conversations_history` | Fetch Slack conversation history. | `slack conversations history` / `GET /conversations.history` | Read-only |
| `slack_conversations_replies` | Fetch replies in a Slack thread. | `slack conversations replies` / `GET /conversations.replies` | Read-only |
| `slack_chat_post_message` | Post a Slack message with text, Block Kit blocks, block elements, and attachments. | `slack chat post-message` / `POST /chat.postMessage` | Mutating, non-destructive |
| `slack_chat_update` | Update a Slack message with text, Block Kit blocks, block elements, and attachments. | `slack chat update` / `POST /chat.update` | Mutating, non-destructive |
| `slack_chat_delete` | Delete a Slack message. | `slack chat delete` / `POST /chat.delete` | Destructive |
| `slack_chat_get_permalink` | Get a permalink for a Slack message. | `slack chat get-permalink` / `GET /chat.getPermalink` | Read-only |
| `slack_reactions_add` | Add a Slack message reaction. | `slack reactions add` / `POST /reactions.add` | Mutating, non-destructive |
| `slack_reactions_remove` | Remove a Slack message reaction. | `slack reactions remove` / `POST /reactions.remove` | Destructive |
| `slack_files_info` | Show Slack file metadata. | `slack files info` / `GET /files.info` | Read-only |
| `slack_files_download` | Download a Slack file to a local path. | `slack files download` / `GET /files.info` + file download URL | Mutating local filesystem, non-destructive |
| `slack_files_upload` | Upload a local file to Slack. | `slack files upload` / `POST /files.getUploadURLExternal` + upload URL + `POST /files.completeUploadExternal` | Mutating, non-destructive |
| `slack_emoji_list` | List Slack emoji. | `slack emoji list` / `GET /emoji.list` | Read-only |

### Output

By default, Slack API commands use stable plain text output:

- key/value lines for single-object responses
- structured message blocks for conversation history
- file IDs in conversation history and thread replies when Slack returns uploaded files
- TSV tables for list responses

Use `--json` to emit compact provider JSON instead.

## Configuration

`slack` uses one configuration shape:

```json
{
  "baseUrl": "https://slack.com/api"
}
```

Configuration can be supplied in either of these ways:

- `SLACK_BASE_URL=https://slack.com/api`

### Examples

```sh
slack help
slack auth test
slack auth test --json
slack conversations list --limit 5
slack conversations list --types public_channel,private_channel --limit 5 --json
slack conversations info --channel C0123456789
slack conversations info --channel C0123456789 --include-locale
slack conversations history --channel C0123456789 --limit 2
slack conversations history --channel C0123456789 --limit 2 --json
slack conversations replies --channel C0123456789 --ts 1775661449.396699 --limit 10
slack chat post-message --channel C0123456789 --text 'hello from slack cli'
slack chat post-message --channel C0123456789 --text 'deploy status' --blocks '[{"type":"section","text":{"type":"mrkdwn","text":"*Deploy:* complete"}}]'
slack chat update --channel C0123456789 --ts 1775060927.238849 --text 'updated text'
slack chat update --channel C0123456789 --ts 1775060927.238849 --blocks-file ./blocks.json
slack chat delete --channel C0123456789 --ts 1775060927.238849
slack chat get-permalink --channel C0123456789 --message-ts 1775060927.238849
slack reactions add --channel C0123456789 --timestamp 1775060927.238849 --name eyes
slack reactions remove --channel C0123456789 --timestamp 1775060927.238849 --name eyes
slack files info --file F0123456789
slack files download --file F0123456789 --output ./downloaded-report.txt
slack files upload --path ./report.txt --channel C0123456789
slack files upload --path ./report.txt --channel C0123456789 --initial-comment 'latest report'
slack emoji list
slack emoji list --include-categories --json
slack mcp serve
slack mcp serve --addr 127.0.0.1:8080 --endpoint /mcp
```

## Build

```sh
mkdir -p dist && go build -o dist/slack ./cmd/slack
```
