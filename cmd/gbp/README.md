# gbp

Google Business Profile CLI and local MCP server.

The binary does not read OAuth client IDs, client secrets, refresh tokens, service
account keys, or access tokens. Configure managed egress or a local proxy to
inject bearer authentication before requests reach Google.

## Configuration

Set all base URLs without trailing slashes:

```sh
GBP_ACCOUNT_MANAGEMENT_BASE_URL=https://mybusinessaccountmanagement.googleapis.com
GBP_BUSINESS_INFORMATION_BASE_URL=https://mybusinessbusinessinformation.googleapis.com
GBP_PERFORMANCE_BASE_URL=https://businessprofileperformance.googleapis.com
GBP_MYBUSINESS_BASE_URL=https://mybusiness.googleapis.com
```

Google Business Profile generally requires OAuth with:

```text
https://www.googleapis.com/auth/business.manage
```

The runtime binary only sends HTTP requests to the configured base URLs.

## CLI

Discovery and reads:

```sh
gbp auth test
gbp accounts list
gbp accounts get --account accounts/<id>
gbp locations list --account accounts/<id> --read-mask name,title,storeCode,websiteUri
gbp locations get --location locations/<id> --read-mask name,title,storeCode,websiteUri
gbp reviews list --account accounts/<id> --location locations/<id>
gbp reviews get --account accounts/<id> --location locations/<id> --review <review-id>
gbp media list --account accounts/<id> --location locations/<id>
gbp media get --media accounts/<id>/locations/<id>/media/<id>
gbp local-posts list --account accounts/<id> --location locations/<id>
gbp local-posts get --local-post accounts/<id>/locations/<id>/localPosts/<id>
gbp performance daily-metrics --location locations/<id> --request-file daily-metrics.json
gbp performance search-keywords --location locations/<id> --request-file search-keywords.json
```

Writes:

```sh
gbp locations create --account accounts/<id> --request-file location.json
gbp locations patch --location locations/<id> --update-mask title --request-file location.json
gbp locations delete --location locations/<id>
gbp reviews update-reply --account accounts/<id> --location locations/<id> --review <review-id> --request-file reply.json
gbp reviews delete-reply --account accounts/<id> --location locations/<id> --review <review-id>
gbp media create --account accounts/<id> --location locations/<id> --request-file media.json
gbp media patch --media accounts/<id>/locations/<id>/media/<id> --update-mask description --request-file media.json
gbp media delete --media accounts/<id>/locations/<id>/media/<id>
gbp media start-upload --account accounts/<id> --location locations/<id>
gbp local-posts create --account accounts/<id> --location locations/<id> --request-file local-post.json
gbp local-posts patch --local-post accounts/<id>/locations/<id>/localPosts/<id> --update-mask summary --request-file local-post.json
gbp local-posts delete --local-post accounts/<id>/locations/<id>/localPosts/<id>
gbp local-posts report-insights --account accounts/<id> --location locations/<id> --request-file local-post-insights.json
```

Report-like commands print pretty JSON. Account and location list commands have
compact text output by default and support `--json`.

Performance request files are JSON objects that map to Google's query parameter
shape. Nested objects are flattened into dotted query parameters.

Example daily metrics request:

```json
{
  "dailyMetric": "WEBSITE_CLICKS",
  "dailyRange": {
    "start_date": {
      "year": 2026,
      "month": 6,
      "day": 1
    },
    "end_date": {
      "year": 2026,
      "month": 6,
      "day": 2
    }
  }
}
```

Example monthly search keywords request:

```json
{
  "monthlyRange": {
    "start_month": {
      "year": 2026,
      "month": 5
    },
    "end_month": {
      "year": 2026,
      "month": 6
    }
  },
  "pageSize": 10
}
```

## MCP

```sh
gbp mcp serve
gbp mcp serve --addr 127.0.0.1:7351 --endpoint /mcp
```

Tools:

- `gbp_auth_test`
- `gbp_accounts_list`
- `gbp_account_get`
- `gbp_locations_list`
- `gbp_location_get`
- `gbp_location_create`
- `gbp_location_patch`
- `gbp_location_delete`
- `gbp_reviews_list`
- `gbp_review_get`
- `gbp_review_update_reply`
- `gbp_review_delete_reply`
- `gbp_media_list`
- `gbp_media_create`
- `gbp_media_get`
- `gbp_media_patch`
- `gbp_media_delete`
- `gbp_media_start_upload`
- `gbp_local_posts_list`
- `gbp_local_post_create`
- `gbp_local_post_get`
- `gbp_local_post_patch`
- `gbp_local_post_delete`
- `gbp_local_post_report_insights`
- `gbp_performance_daily_metrics`
- `gbp_performance_search_keywords`
