# googleads

Thin CLI and MCP wrapper for Google Ads API REST endpoints.

## Configuration

`googleads` does not mint, refresh, store, or inspect credentials. Configure the API base URL and inject Google Ads credentials outside the binary.

Required environment:

- `GOOGLEADS_BASE_URL`, for example `https://googleads.googleapis.com/v24`

The caller or proxy must inject:

- `Authorization: Bearer <oauth-access-token>`
- `developer-token: <google-ads-developer-token>`
- `login-customer-id: <manager-customer-id>` when required by the Google Ads account hierarchy

## Commands

```bash
googleads auth test
googleads customers list-accessible
googleads gaql search --customer-id 1234567890 --query 'SELECT customer.id FROM customer LIMIT 1'
googleads gaql search-stream --customer-id 1234567890 --query 'SELECT customer.id FROM customer LIMIT 1'
googleads fields search --query 'SELECT name, category, data_type WHERE name = "campaign.id"'
googleads fields get --resource-name googleAdsFields/campaign.id
googleads request --method POST --path /customers/1234567890/googleAds:search --body '{"query":"SELECT customer.id FROM customer LIMIT 1"}'
googleads mcp serve
```

`googleads request` is the complete API coverage surface. Named commands and MCP tools provide progressive discovery for common Google Ads workflows.

## Integration Tests

Integration tests use the real Google Ads API through `internal/testproxy` and skip when required environment variables are absent.

Required:

- `GOOGLEADS_TEST_ACCESS_TOKEN`
- `GOOGLEADS_TEST_DEVELOPER_TOKEN`
- `GOOGLEADS_TEST_CUSTOMER_ID`

Optional:

- `GOOGLEADS_TEST_API_VERSION`, defaults to `v24`
- `GOOGLEADS_TEST_LOGIN_CUSTOMER_ID`
