# JIRA

CLI for Jira Cloud.

## Overview

`jira` is a standalone command-line interface for Jira Cloud.

The CLI covers the common read-oriented Jira workflows needed by Mistle users
and provider integrations, including:

- identity checks
- project discovery
- issue lookup and search

This command surface is expanding in a stacked series of small PRs. The goals
for that expansion are:

- preserve a small, provider-focused CLI
- mirror Jira's API model instead of hiding it behind one opaque mutation command
- keep `help` and this README aligned so coding agents can discover commands
  progressively
- document the Atlassian auth scopes required for each command

## Usage

### Commands

The supported commands are:

- `jira help`
- `jira version`
- `jira auth help`
- `jira auth whoami`
- `jira project help`
- `jira project list`
- `jira issue help`
- `jira issue get <key>`
- `jira issue search '<jql query>'`

### Progressive Discovery

Every command family should be discoverable through nested help. The target
shape for the command tree in this stack is:

```text
jira help
jira auth help
jira auth whoami

jira project help
jira project list

jira issue help
jira issue get <key>
jira issue search '<jql query>'

jira issue comment help
jira issue comment add <issue-key> --body <text>
jira issue comment add <issue-key> --body-file <path>

jira issue assign help
jira issue assign <issue-key> --me
jira issue assign <issue-key> --account-id <account-id>
jira issue assign <issue-key> --unassigned

jira issue transition help
jira issue transition list <issue-key>
jira issue transition <issue-key> --to <transition-name>
jira issue transition <issue-key> --to-id <transition-id>

jira issue update help
jira issue update <issue-key> --summary <text>
jira issue update <issue-key> --description <text>
jira issue update <issue-key> --description-file <path>

jira issue editmeta <issue-key>
```

Status changes, assignment, comments, and ordinary field edits are intentionally
separate command families because Jira exposes them as separate API operations.

### Planned Write Commands

The staged rollout for write support in this stack is:

1. nested help surfaces for the new command families
2. `jira issue comment add`
3. `jira issue assign`
4. `jira issue transition list` and `jira issue transition`
5. `jira issue update` for `summary` and `description`
6. `jira issue editmeta`

## Command Reference Format

As write commands land, this README should maintain an exhaustive command
reference table with the following columns:

- `Command`
- `Purpose`
- `Endpoint(s)`
- `OAuth 2.0 classic`
- `Granular / scoped API token scopes`
- `Notes`

Rules for that table:

- use the exact Atlassian scope strings from the REST API docs
- if a command calls more than one endpoint, list the union of scopes across
  those endpoints
- use `Notes` for runtime caveats that are easy to miss
- keep examples below the table so the table remains easy to scan

In this README, permissions should mean Atlassian auth scopes for OAuth 2.0 and
scoped API tokens. Jira project permissions like Browse projects, Edit issues,
Add comments, and Transition issues are separate runtime requirements and should
only appear as short notes when they help explain a failure mode.

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
