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

All API commands also accept `--json` for compact JSON output.

### Discovery

Use help to discover the currently supported command families:

```sh
slack help
slack auth help
slack conversations help
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
```

## Build

```sh
mkdir -p dist && go build -o dist/slack ./cmd/slack
```
