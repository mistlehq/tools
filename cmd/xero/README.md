# Xero

CLI for Xero API access through Mistle-managed credentials.

## Overview

`xero` is a standalone command-line interface and local MCP server for Xero.
It calls Xero API endpoints directly and relies on Mistle managed egress to
inject OAuth bearer credentials.

The first surface is intentionally endpoint-oriented so agents can use Xero API
families that are not covered by Xero's official MCP server.

## Configuration

Set:

```sh
XERO_API_BASE_URL=https://api.xero.com
```

## Commands

- `xero help`
- `xero version`
- `xero tenants list`
- `xero api get --family <family> --tenant-id <tenant-id> --endpoint <path>`
- `xero api post --family <family> --tenant-id <tenant-id> --endpoint <path> --body <json>`
- `xero api put --family <family> --tenant-id <tenant-id> --endpoint <path> --body <json>`
- `xero api delete --family <family> --tenant-id <tenant-id> --endpoint <path>`
- `xero mcp help`
- `xero mcp serve`

Supported API families:

- `accounting` maps to `/api.xro/2.0`
- `assets` maps to `/assets.xro/1.0`
- `files` maps to `/files.xro/1.0`
- `projects` maps to `/projects.xro/2.0`

## MCP

`xero mcp serve` runs Xero as a local MCP server over Streamable HTTP. By
default, it listens on `127.0.0.1:7355` and serves MCP at `/mcp`.

The MCP tools expose documented Xero API endpoints with structured inputs:

| Tool | Purpose | Annotation |
| --- | --- | --- |
| `xero_tenants_list` | List Xero tenant connections for the current token. | Read-only |
| `xero_api_get` | Call a documented GET endpoint for an API family. | Read-only |
| `xero_api_post` | Call a documented POST endpoint for an API family. | Mutating, non-destructive |
| `xero_api_put` | Call a documented PUT endpoint for an API family. | Mutating, idempotent |
| `xero_api_delete` | Call a documented DELETE endpoint for an API family. | Destructive |
