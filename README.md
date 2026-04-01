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

## Supported CLIs

- [`jira`](./cmd/jira/README.md)
  Jira CLI for Jira Cloud.

## Install

Build from source:

```sh
mise trust ./mise.toml
mise install
mkdir -p dist && go build -o dist/jira ./cmd/jira
```

## Usage

Each CLI has its own README with command-specific documentation.
