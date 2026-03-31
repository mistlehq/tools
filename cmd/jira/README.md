# JIRA

CLI for Jira Cloud.

## Overview

`jira` is a standalone command-line interface for Jira Cloud.

The CLI covers the common read-oriented Jira workflows needed by Mistle users
and provider integrations, including:

- identity checks
- project discovery
- issue lookup and search

## Usage

### Commands

The supported commands are:

- `jira help`
- `jira version`
- `jira auth whoami`
- `jira project list`
- `jira issue get <key>`
- `jira issue search '<jql query>'`

### Configuration

`jira` uses one configuration shape:

```json
{
  "baseUrl": "https://example.atlassian.net"
}
```

Configuration can be supplied in either of these ways:

- `JIRA_BASE_URL=https://example.atlassian.net`

### Examples

```sh
jira auth whoami
jira project list
jira issue get PROJ-123
jira issue search 'project = PROJ ORDER BY updated DESC'
```

## Build

```sh
mkdir -p dist && go build -o dist/jira ./cmd/jira
```
