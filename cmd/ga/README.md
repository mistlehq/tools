# ga

`ga` is a standalone command-line interface for Google Analytics.

It is designed to run behind an auth-injecting proxy. The binary does not manage Google OAuth credentials itself; it sends requests to configured Google Analytics API base URLs and expects the caller's proxy layer to inject authorization.

## Commands

- `ga help`
- `ga version`
- `ga auth help`
- `ga auth test --property properties/<id> [--json]`
- `ga account-summaries help`
- `ga account-summaries list [--json]`
- `ga properties help`
- `ga properties get --property properties/<id> [--json]`
- `ga metadata help`
- `ga metadata get --property properties/<id> [--json]`
- `ga compatibility help`
- `ga compatibility check --property properties/<id> --request-file <json>`
- `ga reports help`
- `ga reports run --property properties/<id> --request-file <json>`
- `ga reports realtime --property properties/<id> --request-file <json>`
- `ga reports funnel --property properties/<id> --request-file <json>`
- `ga google-ads-links help`
- `ga google-ads-links list --property properties/<id> [--json]`
- `ga mcp help`
- `ga mcp serve`

## Configuration

`ga` uses two required environment variables:

```sh
GA_ANALYTICS_DATA_BASE_URL=https://analyticsdata.googleapis.com
GA_ANALYTICS_ADMIN_BASE_URL=https://analyticsadmin.googleapis.com
```

Neither URL may end with `/`.

## MCP

`ga mcp serve` runs Google Analytics as a local MCP server over Streamable HTTP.

By default, the server listens on `127.0.0.1:7347` and serves MCP at `/mcp`:

```sh
ga mcp serve
```

The MCP tools mirror the provider-backed CLI command surface with structured inputs and outputs:

| Tool | Purpose | Annotation |
| --- | --- | --- |
| `ga_auth_test` | Check Google Analytics API access. | Read-only |
| `ga_account_summaries_list` | List account summaries. | Read-only |
| `ga_property_get` | Get property details. | Read-only |
| `ga_metadata_get` | Get dimensions and metrics metadata. | Read-only |
| `ga_compatibility_check` | Check report request compatibility. | Read-only |
| `ga_report_run` | Run a core report. | Read-only |
| `ga_report_realtime` | Run a realtime report. | Read-only |
| `ga_report_funnel` | Run a funnel report. | Read-only |
| `ga_google_ads_links_list` | List Google Ads links for a property. | Read-only |

## Examples

```sh
ga auth test --property properties/123456789
ga account-summaries list --json
ga properties get --property properties/123456789
ga metadata get --property properties/123456789
ga reports run --property properties/123456789 --request-file ./run-report.json
ga mcp serve --addr 127.0.0.1:8080 --endpoint /mcp
```
