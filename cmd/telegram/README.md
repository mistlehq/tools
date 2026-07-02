# telegram

`telegram` is a standalone command-line interface for the Telegram Bot API.

It is designed to run behind an auth-injecting proxy. The CLI never accepts or stores bot tokens. Configure `TELEGRAM_BASE_URL` to point at the proxied Telegram Bot API base URL, and let the proxy inject the `bot<token>` URL path segment upstream.

## Commands

- `telegram help`
- `telegram version`
- `telegram auth help`
- `telegram auth test`
- `telegram chats help`
- `telegram chats get --chat <chat-id-or-username>`
- `telegram messages help`
- `telegram messages send --chat <chat-id-or-username> --text <text>`
- `telegram messages send --chat <chat-id-or-username> --text <text> --thread <message-thread-id>`
- `telegram messages send --chat <chat-id-or-username> --text <text> --parse-mode <mode>`
- `telegram messages edit --chat <chat-id-or-username> --message <message-id> --text <text>`
- `telegram messages delete --chat <chat-id-or-username> --message <message-id>`
- `telegram messages delete-batch --chat <chat-id-or-username> --messages <message-id-csv>`
- `telegram reactions help`
- `telegram reactions set --chat <chat-id-or-username> --message <message-id> --emoji <emoji-csv>`
- `telegram reactions clear --chat <chat-id-or-username> --message <message-id>`
- `telegram reactions delete --chat <chat-id-or-username> --message <message-id> [--user <user-id>] [--actor-chat-id <chat-id>]`
- `telegram reactions delete-all --chat <chat-id-or-username> [--user <user-id>] [--actor-chat-id <chat-id>]`
- `telegram topics help`
- `telegram topics create --chat <chat-id-or-username> --name <name>`
- `telegram topics delete --chat <chat-id-or-username> --thread <message-thread-id>`
- `telegram request help`
- `telegram request --method <telegram-method> [--body <json>]`
- `telegram mcp help`
- `telegram mcp serve`
- `telegram mcp serve --addr <addr>`
- `telegram mcp serve --endpoint <path>`

All provider-returning commands accept `--json`.

## MCP

`telegram mcp serve` runs Telegram as a local MCP server over Streamable HTTP. It reuses the same `TELEGRAM_BASE_URL` configuration as the CLI and relies on the same upstream auth-injecting proxy model.

```sh
telegram mcp serve
telegram mcp serve --addr 127.0.0.1:8080
telegram mcp serve --endpoint /mcp
telegram mcp serve --addr 127.0.0.1:8080 --endpoint /mcp
```

MCP tools:

| Tool | Description | Safety |
| --- | --- | --- |
| `telegram_auth_test` | Check Telegram bot authentication state. | Read-only |
| `telegram_chats_get` | Show details for a Telegram chat. | Read-only |
| `telegram_messages_send` | Send a Telegram text message. | Mutating, non-destructive |
| `telegram_messages_edit` | Edit a Telegram text message. | Mutating, non-destructive |
| `telegram_messages_delete` | Delete a Telegram message. | Destructive |
| `telegram_messages_delete_batch` | Delete multiple Telegram messages. | Destructive |
| `telegram_reactions_set` | Set Telegram reactions on a message. | Mutating, non-destructive |
| `telegram_reactions_clear` | Clear Telegram reactions from a message. | Mutating, non-destructive |
| `telegram_reactions_delete` | Delete a Telegram message reaction. | Destructive |
| `telegram_reactions_delete_all` | Delete Telegram message reactions in bulk. | Destructive |
| `telegram_topics_create` | Create a Telegram forum topic. | Mutating, non-destructive |
| `telegram_topics_delete` | Delete a Telegram forum topic and its messages. | Destructive |
| `telegram_request` | Call an arbitrary Telegram Bot API method. | Open-world |

## Configuration

`telegram` uses one configuration value:

```json
{
  "baseUrl": "https://api.telegram.org"
}
```

In normal CLI usage, provide it through the environment:

```sh
TELEGRAM_BASE_URL=https://api.telegram.org
```

For Mistle, this should usually be the local proxy URL. The proxy should prepend the credential-backed `bot<token>` path segment before forwarding to Telegram. The CLI itself calls tokenless method paths such as `/getMe`, `/getChat`, `/sendMessage`, `/editMessageText`, `/deleteMessage`, and `/setMessageReaction`.

## Examples

```sh
telegram help
telegram auth test
telegram auth test --json
telegram chats get --chat 123456789
telegram chats get --chat @channelusername --json
telegram messages send --chat 123456789 --text 'hello from telegram cli'
telegram messages send --chat -1001234567890 --thread 42 --text 'hello forum topic'
telegram messages send --chat @channelusername --text '<b>Deploy complete</b>' --parse-mode HTML
telegram topics create --chat -1001234567890 --name 'Mistle test topic'
telegram topics delete --chat -1001234567890 --thread 42
telegram messages edit --chat 123456789 --message 42 --text 'updated text'
telegram reactions set --chat 123456789 --message 42 --emoji '👍'
telegram reactions clear --chat 123456789 --message 42
telegram messages delete --chat 123456789 --message 42
telegram request --method getChat --body '{"chat_id":"123456789"}' --json
telegram mcp serve
telegram mcp serve --addr 127.0.0.1:8080 --endpoint /mcp
```

## Development

Build locally:

```sh
mkdir -p dist && mise exec -- go build -o dist/telegram ./cmd/telegram
```

Run focused tests:

```sh
mise exec -- go test ./cmd/telegram
```

Live Telegram integration tests require:

```sh
TELEGRAM_TEST_UPSTREAM_BASE_URL=https://api.telegram.org
TELEGRAM_TEST_BOT_TOKEN=123456:ABC...
TELEGRAM_TEST_CHAT_ID=123456789
TELEGRAM_TEST_FORUM_CHAT_ID=-1001234567890
```
