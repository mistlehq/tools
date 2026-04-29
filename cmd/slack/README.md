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
- `slack chat post-message --channel <conversation-id> --thread-ts <ts> --text <text>`
- `slack chat update --channel <conversation-id> --ts <ts> --text <text>`
- `slack chat update --channel <conversation-id> --ts <ts> --text-file <path>`
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
```

### Output

By default, Slack API commands use stable plain text output:

- key/value lines for single-object responses
- structured message blocks for conversation history
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
slack chat update --channel C0123456789 --ts 1775060927.238849 --text 'updated text'
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
```

## Build

```sh
mkdir -p dist && go build -o dist/slack ./cmd/slack
```
