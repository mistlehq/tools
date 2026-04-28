# JIRA

CLI for Jira Cloud.

## Overview

`jira` is a standalone command-line interface for Jira Cloud.

The CLI covers common Jira workflows needed by Mistle users and provider integrations, including:

- identity checks
- project discovery
- issue creation
- issue lookup and search
- issue deletion
- issue comments
- comment deletion
- assignment
- workflow transitions
- summary, description, and editable field updates
- edit metadata inspection

## Usage

`jira help` is intended to work like an onboarding page for the CLI: it should give a new user enough context to start using the tool immediately, then point them at deeper namespace help for specific operations.

All namespaces and leaf commands also accept `--help`, so agent-style discovery flows like `jira issue get --help` and `jira issue search --help` are expected to work.

### Commands

The supported commands are:

- `jira help`
- `jira version`
- `jira auth help`
- `jira auth whoami`
- `jira project help`
- `jira project list`
- `jira issue help`
- `jira issue create --project-key <key> --issue-type <name> --summary <text>`
- `jira issue create --project-id <id> --issue-type-id <id> --summary <text>`
- `jira issue create --project-key <key> --issue-type <name> --summary <text> --description <text>`
- `jira issue create --project-key <key> --issue-type <name> --summary <text> --description-file <path>`
- `jira issue get <key>`
- `jira issue search '<jql query>'`
- `jira issue delete <key>`
- `jira issue comment help`
- `jira issue comment add <issue-key> --body <text>`
- `jira issue comment add <issue-key> --body-file <path>`
- `jira issue comment delete <issue-key> <comment-id>`
- `jira issue assign help`
- `jira issue assign <issue-key> --me`
- `jira issue assign <issue-key> --account-id <account-id>`
- `jira issue assign <issue-key> --unassigned`
- `jira issue transition help`
- `jira issue transition list <issue-key>`
- `jira issue transition <issue-key> --to <transition-name>`
- `jira issue transition <issue-key> --to-id <transition-id>`
- `jira issue update help`
- `jira issue update <issue-key> --summary <text>`
- `jira issue update <issue-key> --description <text>`
- `jira issue update <issue-key> --description-file <path>`
- `jira issue update <issue-key> --field <field-id=value>`
- `jira issue update <issue-key> --field-json <field-id=json>`
- `jira issue editmeta help`
- `jira issue editmeta <issue-key>`

### Discovery

The CLI groups commands under `auth`, `project`, and `issue`. Use help to discover the available command families:

```sh
jira help
jira --help
jira auth help
jira auth whoami --help
jira project help
jira project list --help
jira issue help
jira issue create --help
jira issue get --help
jira issue search --help
jira issue delete --help
jira issue comment help
jira issue comment add --help
jira issue comment delete --help
jira issue assign help
jira issue assign --help
jira issue transition help
jira issue transition list --help
jira issue update help
jira issue update --help
jira issue editmeta help
jira issue editmeta --help
```

Status changes, assignment, comments, and ordinary field edits are intentionally separate command families because Jira exposes them as separate API operations.

## Auth Scopes

For this README, permissions means Atlassian auth scopes for OAuth 2.0 and scoped API tokens. Jira project permissions like Browse projects, Edit issues, Add comments, and Transition issues are separate runtime requirements and only appear in notes when they explain a likely failure mode.

The tables below map commands to the Jira REST endpoints they call. For commands that call more than one endpoint, the scope columns show the union of those endpoints. Some Atlassian docs pages collapse long granular scope lists in the UI; those rows are marked accordingly.

### Local Commands

| Command | Purpose | Endpoint(s) | OAuth 2.0 classic | Granular / scoped API token scopes | Notes |
| --- | --- | --- | --- | --- | --- |
| `jira help` | Show top-level help. | Local only | None | None | No Jira request is made. |
| `jira version` | Show the CLI version. | Local only | None | None | No Jira request is made. |
| `jira auth help` | Show auth help. | Local only | None | None | No Jira request is made. |
| `jira project help` | Show project help. | Local only | None | None | No Jira request is made. |
| `jira issue help` | Show issue help. | Local only | None | None | No Jira request is made. |
| `jira issue create help` | Show issue create help. | Local only | None | None | No Jira request is made. |
| `jira issue delete help` | Show issue delete help. | Local only | None | None | No Jira request is made. |
| `jira issue comment help` | Show issue comment help. | Local only | None | None | No Jira request is made. |
| `jira issue comment delete help` | Show issue comment delete help. | Local only | None | None | No Jira request is made. |
| `jira issue assign help` | Show issue assign help. | Local only | None | None | No Jira request is made. |
| `jira issue transition help` | Show issue transition help. | Local only | None | None | No Jira request is made. |
| `jira issue update help` | Show issue update help. | Local only | None | None | No Jira request is made. |
| `jira issue editmeta help` | Show issue editmeta help. | Local only | None | None | No Jira request is made. |

### Jira API Commands

| Command | Purpose | Endpoint(s) | OAuth 2.0 classic | Granular / scoped API token scopes | Notes |
| --- | --- | --- | --- | --- | --- |
| `jira auth whoami` | Show the current Jira user. | `GET /rest/api/3/myself` | `read:jira-user` | `read:application-role:jira`, `read:group:jira`, `read:user:jira`, `read:avatar:jira` | Requires permission to access Jira. |
| `jira project list` | List visible Jira projects. | `GET /rest/api/3/project/search` | `read:jira-work` | `read:issue-type:jira`, `read:project:jira`, `read:project.property:jira`, `read:user:jira`, `read:application-role:jira`, `...` | Atlassian collapses the full granular list in the docs UI for this endpoint. |
| `jira issue create --project-key <key> --issue-type <name> --summary <text>` | Create an issue by project key and issue type name. | `POST /rest/api/3/issue` | `write:jira-work` | `write:issue:jira` | Runtime Jira project permissions, required fields, and field configuration still apply. |
| `jira issue create --project-id <id> --issue-type-id <id> --summary <text>` | Create an issue by project ID and issue type ID. | `POST /rest/api/3/issue` | `write:jira-work` | `write:issue:jira` | Useful when names are ambiguous or localized. |
| `jira issue create ... --description <text>` | Create an issue with an inline description. | `POST /rest/api/3/issue` | `write:jira-work` | `write:issue:jira` | The CLI converts plain text to Atlassian Document Format internally. |
| `jira issue create ... --description-file <path>` | Create an issue with a description from a file or stdin. | `POST /rest/api/3/issue` | `write:jira-work` | `write:issue:jira` | Use `--description-file -` to read the description from stdin. |
| `jira issue get <key>` | Fetch a single issue. | `GET /rest/api/3/issue/{issueIdOrKey}` | `read:jira-work` | `read:issue-meta:jira`, `read:issue-security-level:jira`, `read:issue.vote:jira`, `read:issue.changelog:jira`, `read:avatar:jira`, `...` | Atlassian collapses the full granular list in the docs UI for this endpoint. |
| `jira issue search '<jql>'` | Search issues with JQL. | `POST /rest/api/3/search/jql` | `read:jira-work` | `read:issue-details:jira`, `read:audit-log:jira`, `read:avatar:jira`, `read:field-configuration:jira`, `read:issue-meta:jira` | Returns only issues visible to the caller. |
| `jira issue delete <key>` | Delete a single issue. | `DELETE /rest/api/3/issue/{issueIdOrKey}` | `write:jira-work` | `delete:issue:jira` | Jira may require extra project permissions or `deleteSubtasks=true` when subtasks exist. |
| `jira issue comment add <key> --body <text>` | Add a comment from inline text. | `POST /rest/api/3/issue/{issueIdOrKey}/comment` | `write:jira-work` | `read:comment:jira`, `read:comment.property:jira`, `read:group:jira`, `read:project:jira`, `read:project-role:jira`, `...` | Atlassian collapses the full granular list in the docs UI for this endpoint. |
| `jira issue comment add <key> --body-file <path>` | Add a comment from a file or stdin. | `POST /rest/api/3/issue/{issueIdOrKey}/comment` | `write:jira-work` | `read:comment:jira`, `read:comment.property:jira`, `read:group:jira`, `read:project:jira`, `read:project-role:jira`, `...` | Use `--body-file -` to read the comment body from stdin. |
| `jira issue comment delete <key> <comment-id>` | Delete a comment. | `DELETE /rest/api/3/issue/{issueIdOrKey}/comment/{id}` | `write:jira-work` | `delete:comment:jira`, `delete:comment.property:jira` | Runtime Jira comment visibility and delete permissions still apply. |
| `jira issue assign <key> --account-id <id>` | Assign the issue to a specific Jira account. | `PUT /rest/api/3/issue/{issueIdOrKey}/assignee` | `write:jira-work` | `write:issue:jira` | Runtime Jira project permissions still apply. |
| `jira issue assign <key> --unassigned` | Clear the assignee. | `PUT /rest/api/3/issue/{issueIdOrKey}/assignee` | `write:jira-work` | `write:issue:jira` | Runtime Jira project permissions still apply. |
| `jira issue assign <key> --me` | Assign the issue to the current user. | `GET /rest/api/3/myself` + `PUT /rest/api/3/issue/{issueIdOrKey}/assignee` | `read:jira-user` + `write:jira-work` | `read:application-role:jira`, `read:group:jira`, `read:user:jira`, `read:avatar:jira`, `write:issue:jira` | Uses `whoami`-style identity lookup before assignment. |
| `jira issue transition list <key>` | List available workflow transitions. | `GET /rest/api/3/issue/{issueIdOrKey}/transitions` | `read:jira-work` | `read:issue.transition:jira`, `read:status:jira`, `read:field-configuration:jira` | Available transitions depend on the issue's current workflow state. |
| `jira issue transition <key> --to <name>` | Transition by exact transition name. | `GET /rest/api/3/issue/{issueIdOrKey}/transitions` + `POST /rest/api/3/issue/{issueIdOrKey}/transitions` | `read:jira-work` + `write:jira-work` | `read:issue.transition:jira`, `read:status:jira`, `read:field-configuration:jira`, `write:issue:jira`, `write:issue.property:jira` | The CLI resolves the name against the live transition list first. |
| `jira issue transition <key> --to-id <id>` | Transition by transition id. | `GET /rest/api/3/issue/{issueIdOrKey}/transitions` + `POST /rest/api/3/issue/{issueIdOrKey}/transitions` | `read:jira-work` + `write:jira-work` | `read:issue.transition:jira`, `read:status:jira`, `read:field-configuration:jira`, `write:issue:jira`, `write:issue.property:jira` | The CLI resolves the id against the live transition list first. |
| `jira issue update <key> --summary <text>` | Update the issue summary. | `PUT /rest/api/3/issue/{issueIdOrKey}` | `write:jira-work` | `write:issue:jira` | Status changes do not go through `update`. |
| `jira issue update <key> --description <text>` | Update the issue description from inline text. | `PUT /rest/api/3/issue/{issueIdOrKey}` | `write:jira-work` | `write:issue:jira` | The CLI converts plain text to Atlassian Document Format internally. |
| `jira issue update <key> --description-file <path>` | Update the issue description from a file or stdin. | `PUT /rest/api/3/issue/{issueIdOrKey}` | `write:jira-work` | `write:issue:jira` | Use `--description-file -` to read the description from stdin. |
| `jira issue update <key> --field <field-id=value>` | Update an editable field with a string value. | `PUT /rest/api/3/issue/{issueIdOrKey}` | `write:jira-work` | `write:issue:jira` | Use `editmeta` first to inspect editable field IDs. |
| `jira issue update <key> --field-json <field-id=json>` | Update an editable field with a JSON value. | `PUT /rest/api/3/issue/{issueIdOrKey}` | `write:jira-work` | `write:issue:jira` | Use for arrays, objects, numbers, booleans, and `null`. |
| `jira issue editmeta <key>` | Show editable field metadata for the issue. | `GET /rest/api/3/issue/{issueIdOrKey}/editmeta` | `read:jira-work` | `read:issue-meta:jira`, `read:field-configuration:jira` | Useful before generic field editing. |

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
jira auth help
jira project list
jira issue help
jira issue create --project-key PROJ --issue-type Task --summary 'Tighten validation'
jira issue create --project-key PROJ --issue-type Bug --summary 'Fix auth error' --description 'Expanded implementation notes'
jira issue create --project-key PROJ --issue-type Task --summary 'Plan rollout' --description-file ./description.txt
jira issue create --project-key PROJ --issue-type Task --summary 'Plan rollout' --description-file -
jira issue get PROJ-123
jira issue search 'project = PROJ ORDER BY updated DESC'
jira issue delete PROJ-123
jira issue comment add PROJ-123 --body 'Looks good'
jira issue comment add PROJ-123 --body-file ./comment.txt
jira issue comment add PROJ-123 --body-file -
jira issue comment delete PROJ-123 10001
jira issue assign PROJ-123 --me
jira issue assign PROJ-123 --account-id 712020:abc123
jira issue assign PROJ-123 --unassigned
jira issue transition list PROJ-123
jira issue transition PROJ-123 --to 'In Progress'
jira issue transition PROJ-123 --to-id 31
jira issue update PROJ-123 --summary 'Tighten validation'
jira issue update PROJ-123 --description 'Expanded implementation notes'
jira issue update PROJ-123 --description-file ./description.txt
jira issue update PROJ-123 --description-file -
jira issue update PROJ-123 --field customfield_10010='Customer impact'
jira issue update PROJ-123 --field-json labels='["backend","urgent"]'
jira issue update PROJ-123 --field-json customfield_10020=null
jira issue editmeta PROJ-123
```

## Build

```sh
mkdir -p dist && go build -o dist/jira ./cmd/jira
```
