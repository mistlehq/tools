# JIRA

CLI for Jira Cloud.

## Overview

`jira` is a standalone command-line interface for Jira Cloud.

The CLI covers the common Jira workflows needed by Mistle users and provider integrations, including:

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
- `jira issue comment help`
- `jira issue comment add <issue-key> --body <text>`
- `jira issue comment add <issue-key> --body-file <path>`
- `jira issue assign help`
- `jira issue assign <issue-key> --me`
- `jira issue assign <issue-key> --account-id <account-id>`
- `jira issue assign <issue-key> --unassigned`
- `jira issue transition help`
- `jira issue update help`
- `jira issue editmeta help`

### Discovery

The CLI groups commands under `auth`, `project`, and `issue`. Use the top-level help to discover the available command families:

```sh
jira help
jira auth
jira project
jira issue
```

Status changes, assignment, comments, and ordinary field edits are intentionally separate command families because Jira exposes them as separate API operations.

Additional issue help pages make the command families easier to discover:

```sh
jira issue comment help
jira issue assign help
jira issue transition help
jira issue update help
jira issue editmeta help
```

### Discovery

The CLI groups commands under `auth`, `project`, and `issue`. Use the top-level help to discover the available command families:

```sh
jira help
jira auth
jira project
jira issue
```

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
jira help
jira auth whoami
jira project list
jira issue get PROJ-123
jira issue search 'project = PROJ ORDER BY updated DESC'
jira issue comment add PROJ-123 --body 'Looks good'
jira issue comment add PROJ-123 --body-file ./comment.txt
jira issue assign PROJ-123 --me
jira issue assign PROJ-123 --account-id 712020:abc123
jira issue assign PROJ-123 --unassigned
```

## Build

```sh
mkdir -p dist && go build -o dist/jira ./cmd/jira
```
