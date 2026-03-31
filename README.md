# mistle-tools

Open source provider CLIs for Mistle.

## Overview

`mistle-tools` is a collection of standalone provider CLIs. Each provider gets
its own executable.

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
