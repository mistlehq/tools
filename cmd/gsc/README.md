# gsc

`gsc` is a standalone command-line interface for Google Search Console.

It is designed to run behind an auth-injecting proxy. The binary does not manage Google OAuth credentials itself; it sends requests to a configured Search Console API base URL and expects the caller's proxy layer to inject authorization.

## Commands

- `gsc help`
- `gsc version`
- `gsc auth help`
- `gsc auth test --site-url <site-url> [--json]`
- `gsc sites help`
- `gsc sites list [--json]`
- `gsc sites get --site-url <site-url> [--json]`
- `gsc searchanalytics help`
- `gsc searchanalytics query --site-url <site-url> --request-file <json>`
- `gsc sitemaps help`
- `gsc sitemaps list --site-url <site-url> [--json]`
- `gsc sitemaps get --site-url <site-url> --feed-path <sitemap-url> [--json]`
- `gsc url-inspection help`
- `gsc url-inspection inspect --request-file <json>`
- `gsc mcp help`
- `gsc mcp serve`

## Configuration

`gsc` uses one required environment variable:

```sh
GSC_SEARCH_CONSOLE_BASE_URL=https://searchconsole.googleapis.com
```

The URL may not end with `/`.

## MCP

`gsc mcp serve` runs Google Search Console as a local MCP server over Streamable HTTP.

By default, the server listens on `127.0.0.1:7349` and serves MCP at `/mcp`:

```sh
gsc mcp serve
```

The MCP tools mirror the provider-backed CLI command surface with structured inputs and outputs:

| Tool | Purpose | Annotation |
| --- | --- | --- |
| `gsc_auth_test` | Check Search Console API access. | Read-only |
| `gsc_sites_list` | List visible sites. | Read-only |
| `gsc_site_get` | Get site details. | Read-only |
| `gsc_searchanalytics_query` | Query Search Analytics data. | Read-only |
| `gsc_sitemaps_list` | List sitemaps for a site. | Read-only |
| `gsc_sitemap_get` | Get sitemap details. | Read-only |
| `gsc_url_inspection_inspect` | Inspect URL indexing information. | Read-only |

## Examples

```sh
gsc auth test --site-url https://example.com/
gsc sites list --json
gsc sites get --site-url sc-domain:example.com
gsc searchanalytics query --site-url https://example.com/ --request-file ./searchanalytics-query.json
gsc url-inspection inspect --request-file ./url-inspection.json
gsc mcp serve --addr 127.0.0.1:8080 --endpoint /mcp
```
