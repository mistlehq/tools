# mistle-tools

Open source provider CLIs for Mistle.

## Overview

`mistle-tools` is a collection of standalone provider CLIs. Each provider gets its own executable.

These CLIs are designed to run behind a proxy that intercepts outbound requests and injects provider credentials on the caller's behalf. The CLIs therefore do not implement credential handling themselves; they only need a base URL to talk to the proxied upstream.

This matches [Mistle](https://github.com/mistlehq/mistle)'s security model:

- credentials stay outside the CLI process
- auth injection happens at the proxy boundary
- the CLI stays focused on provider commands and output

## Design Philosophy

All CLIs in this repository follow the same basic shape:

- They are thin wrappers around provider APIs rather than full local platforms or long-running services.
- They are non-interactive and are designed to work well in shells, scripts, and agent-driven workflows.
- They do not handle auth directly. They operate as if they are already authenticated and rely on the configured proxy layer to inject credentials upstream.
- Their root `help` output should work like a landing-page README for first-time users: quick orientation, major command families, common workflows, and clear next steps.
- Progressive discovery is still expected, but every namespace and leaf command should accept `--help` so users and agents can ask for local command guidance without trial-and-error.
- When a CLI exposes MCP support, the MCP server should stay local by default and expose provider API operations as structured tools rather than re-parsing CLI text output.

## Supported CLIs

- [`jira`](./cmd/jira/README.md)
  Jira CLI for Jira Cloud.
- [`slack`](./cmd/slack/README.md)
  Slack CLI.
- [`discord`](./cmd/discord/README.md)
  Discord CLI.
- [`telegram`](./cmd/telegram/README.md)
  Telegram Bot API CLI.
- [`ga`](./cmd/ga/README.md)
  Google Analytics CLI.
- [`googleads`](./cmd/googleads/README.md)
  Google Ads API CLI.
- [`gws`](./cmd/gws/README.md)
  Google Workspace Drive, Sheets, Docs, and Slides CLI.
- [`shopify`](./cmd/shopify/README.md)
  Shopify Admin API CLI.
- [`metaads`](./cmd/metaads/README.md)
  Meta Graph API / Marketing API CLI.

## Install

Build from source:

```sh
mise trust ./mise.toml
mise install
mkdir -p dist && mise exec -- go build -o dist/jira ./cmd/jira
```

## Usage

Each CLI has its own README with command-specific documentation. In general, start with the root help for the binary you are using, then drill down into namespace or leaf help with `--help` as needed.

Some CLIs also expose local MCP servers for agent clients. Start with the provider's MCP namespace help, for example:

```sh
jira mcp help
jira mcp serve --help
```

## Local Integration Tests

Some provider integration tests need real upstream credentials. Copy `.env.example` to `.env.local`, fill in the provider values, then export them before running the focused tests.

For Jira against the Mistle Atlassian tenant:

```sh
set -a
source .env.local
set +a

mise exec -- go test ./cmd/jira -run TestJiraAttachmentAPIAccess -count=1 -v
```

The Jira CLI itself does not inject auth headers. The Jira integration test harness uses `JIRA_TEST_*` values to start a local auth-injecting proxy, then points the CLI at that proxy.
