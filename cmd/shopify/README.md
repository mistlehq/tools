# Shopify

CLI for the Shopify Admin API.

## Overview

`shopify` is a standalone command-line interface for Shopify Admin GraphQL.

The CLI is a thin wrapper around the Admin GraphQL endpoint. It does not manage credentials, app scopes, or Shopify permissions. It expects `SHOPIFY_ADMIN_BASE_URL` to point at an authenticated proxy for a versioned Admin API base URL, such as:

```sh
https://example.myshopify.com/admin/api/2026-04
```

The proxy is responsible for injecting `X-Shopify-Access-Token`.

## Usage

The supported commands are:

- `shopify help`
- `shopify version`
- `shopify auth help`
- `shopify auth test`
- `shopify auth test --json`
- `shopify graphql help`
- `shopify graphql request --query <graphql>`
- `shopify graphql request --query-file <path>`
- `shopify graphql request --query <graphql> --variables <json>`
- `shopify graphql request --query-file <path> --variables-file <path>`
- `shopify shop help`
- `shopify shop get`
- `shopify products help`
- `shopify products search --first <count>`
- `shopify products search --first <count> --query <search>`
- `shopify products get --id <product-gid>`
- `shopify products get --handle <handle>`
- `shopify products create --product-json <json>`
- `shopify products create --product-file <path>`
- `shopify products update --product-json <json>`
- `shopify products update --product-file <path>`
- `shopify products delete --id <product-gid>`
- `shopify orders search --first <count>`
- `shopify orders search --first <count> --query <search>`
- `shopify orders get --id <order-gid>`
- `shopify customers search --first <count>`
- `shopify customers search --first <count> --query <search>`
- `shopify customers get --id <customer-gid>`
- `shopify inventory items search --first <count>`
- `shopify inventory levels search --first <count>`
- `shopify locations list --first <count>`
- `shopify mcp help`
- `shopify mcp serve`
- `shopify mcp serve --addr <addr>`
- `shopify mcp serve --endpoint <path>`

## Full Admin API Coverage

`shopify graphql request` is the full Shopify Admin API coverage surface. It accepts any Admin GraphQL query or mutation supported by the configured API version and allowed by the app's Shopify scopes and runtime permissions.

Named commands are convenience wrappers for common workflows and progressive discovery. They are not an exhaustive local reimplementation of every Shopify query and mutation.

## MCP

`shopify mcp serve` runs Shopify as a local MCP server over Streamable HTTP. It reuses the same `SHOPIFY_ADMIN_BASE_URL` configuration as the CLI and relies on the same upstream auth-injecting proxy model.

By default, the server listens on `127.0.0.1:7348` and serves MCP at `/mcp`:

```sh
shopify mcp serve
```

The MCP tools mirror the supported CLI command surface:

| Tool | Purpose | Annotation |
| --- | --- | --- |
| `shopify_graphql_request` | Send a raw Admin GraphQL request. | Open-world |
| `shopify_auth_test` | Check Admin API access. | Read-only |
| `shopify_shop_get` | Get current shop details. | Read-only |
| `shopify_products_search` | Search products. | Read-only |
| `shopify_product_get` | Get a product. | Read-only |
| `shopify_product_create` | Create a product. | Mutating, non-destructive |
| `shopify_product_update` | Update a product. | Mutating, non-destructive |
| `shopify_product_delete` | Delete a product. | Destructive |
| `shopify_orders_search` | Search orders. | Read-only |
| `shopify_order_get` | Get an order. | Read-only |
| `shopify_customers_search` | Search customers. | Read-only |
| `shopify_customer_get` | Get a customer. | Read-only |
| `shopify_inventory_items_search` | Search inventory items. | Read-only |
| `shopify_inventory_levels_search` | Search inventory levels. | Read-only |
| `shopify_locations_list` | List locations. | Read-only |

## Integration Test Environment

The Shopify tests use a real Shopify store and skip when required environment variables are absent:

```sh
SHOPIFY_TEST_SHOP_DOMAIN=example.myshopify.com
SHOPIFY_TEST_ADMIN_API_VERSION=2026-04
SHOPIFY_TEST_CLIENT_ID=...
SHOPIFY_TEST_CLIENT_SECRET=...
SHOPIFY_TEST_PRODUCT_ID=gid://shopify/Product/...
SHOPIFY_TEST_PRODUCT_HANDLE=...
SHOPIFY_TEST_ORDER_ID=gid://shopify/Order/...
SHOPIFY_TEST_ORDER_NAME=#1001
SHOPIFY_TEST_CUSTOMER_ID=gid://shopify/Customer/...
SHOPIFY_TEST_CUSTOMER_EMAIL=...
```
