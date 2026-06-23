package main

import "fmt"

func (cli CLI) printHelp() {
	fmt.Fprint(cli.stdout, `shopify

CLI for Shopify Admin API.

Usage:
  shopify help
  shopify version
  shopify auth help
  shopify graphql help
  shopify shop help
  shopify products help
  shopify orders help
  shopify customers help
  shopify inventory help
  shopify locations help
  shopify mcp help

Commands:
  help
  version
  auth
  graphql
  shop
  products
  orders
  customers
  inventory
  locations
  mcp
`)
}

func (cli CLI) printAuthHelp() {
	fmt.Fprint(cli.stdout, `shopify auth

Inspect Shopify Admin API authentication state.

Usage:
  shopify auth help
  shopify auth test
  shopify auth test --json
`)
}

func (cli CLI) printAuthTestHelp() {
	fmt.Fprintf(cli.stdout, `shopify auth test

%s

Usage:
  shopify auth test
  shopify auth test --json
`, shopifyAuthTestDoc.Description)
}

func (cli CLI) printGraphQLHelp() {
	fmt.Fprint(cli.stdout, `shopify graphql

Send raw Shopify Admin GraphQL requests.

Usage:
  shopify graphql help
  shopify graphql request --help
`)
}

func (cli CLI) printGraphQLRequestHelp() {
	fmt.Fprintf(cli.stdout, `shopify graphql request

%s

Usage:
  shopify graphql request --query <graphql>
  shopify graphql request --query-file <path>
  shopify graphql request --query <graphql> --variables <json>
  shopify graphql request --query-file <path> --variables-file <path>

Options:
  --query <graphql>          GraphQL query or mutation text.
  --query-file <path>        File containing GraphQL query or mutation text.
  --variables <json>         GraphQL variables JSON object.
  --variables-file <path>    File containing GraphQL variables JSON object.
`, shopifyGraphQLRequestDoc.Description)
}

func (cli CLI) printShopHelp() {
	fmt.Fprint(cli.stdout, `shopify shop

Inspect the current Shopify shop.

Usage:
  shopify shop help
  shopify shop get
  shopify shop get --help
`)
}

func (cli CLI) printShopGetHelp() {
	fmt.Fprintf(cli.stdout, `shopify shop get

%s

Usage:
  shopify shop get
`, shopifyShopGetDoc.Description)
}

func (cli CLI) printProductsHelp() {
	fmt.Fprint(cli.stdout, `shopify products

Work with Shopify products.

Usage:
  shopify products help
  shopify products search --first <count>
  shopify products get --id <product-gid>
  shopify products get --handle <handle>
  shopify products create --product-json <json>
  shopify products create --product-file <path>
  shopify products update --product-json <json>
  shopify products update --product-file <path>
  shopify products delete --id <product-gid>
`)
}

func (cli CLI) printProductsSearchHelp() {
	fmt.Fprintf(cli.stdout, `shopify products search

%s

Usage:
  shopify products search --first <count>
  shopify products search --first <count> --query <search>
  shopify products search --first <count> --after <cursor>
`, shopifyProductsSearchDoc.Description)
}

func (cli CLI) printProductGetHelp() {
	fmt.Fprintf(cli.stdout, `shopify products get

%s

Usage:
  shopify products get --id <product-gid>
  shopify products get --handle <handle>
`, shopifyProductGetDoc.Description)
}

func (cli CLI) printProductCreateHelp() {
	fmt.Fprintf(cli.stdout, `shopify products create

%s

Usage:
  shopify products create --product-json <json>
  shopify products create --product-file <path>
`, shopifyProductCreateDoc.Description)
}

func (cli CLI) printProductUpdateHelp() {
	fmt.Fprintf(cli.stdout, `shopify products update

%s

Usage:
  shopify products update --product-json <json>
  shopify products update --product-file <path>
`, shopifyProductUpdateDoc.Description)
}

func (cli CLI) printProductDeleteHelp() {
	fmt.Fprintf(cli.stdout, `shopify products delete

%s

Usage:
  shopify products delete --id <product-gid>
`, shopifyProductDeleteDoc.Description)
}

func (cli CLI) printOrdersHelp() {
	fmt.Fprint(cli.stdout, `shopify orders

Work with Shopify orders.

Usage:
  shopify orders help
  shopify orders search --first <count>
  shopify orders search --first <count> --query <search>
  shopify orders get --id <order-gid>
`)
}

func (cli CLI) printOrdersSearchHelp() {
	fmt.Fprintf(cli.stdout, `shopify orders search

%s

Usage:
  shopify orders search --first <count>
  shopify orders search --first <count> --query <search>
  shopify orders search --first <count> --after <cursor>
`, shopifyOrdersSearchDoc.Description)
}

func (cli CLI) printOrderGetHelp() {
	fmt.Fprintf(cli.stdout, `shopify orders get

%s

Usage:
  shopify orders get --id <order-gid>
`, shopifyOrderGetDoc.Description)
}

func (cli CLI) printCustomersHelp() {
	fmt.Fprint(cli.stdout, `shopify customers

Work with Shopify customers.

Usage:
  shopify customers help
  shopify customers search --first <count>
  shopify customers search --first <count> --query <search>
  shopify customers get --id <customer-gid>
`)
}

func (cli CLI) printCustomersSearchHelp() {
	fmt.Fprintf(cli.stdout, `shopify customers search

%s

Usage:
  shopify customers search --first <count>
  shopify customers search --first <count> --query <search>
  shopify customers search --first <count> --after <cursor>
`, shopifyCustomersSearchDoc.Description)
}

func (cli CLI) printCustomerGetHelp() {
	fmt.Fprintf(cli.stdout, `shopify customers get

%s

Usage:
  shopify customers get --id <customer-gid>
`, shopifyCustomerGetDoc.Description)
}

func (cli CLI) printInventoryHelp() {
	fmt.Fprint(cli.stdout, `shopify inventory

Work with Shopify inventory.

Usage:
  shopify inventory help
  shopify inventory items help
  shopify inventory items search --first <count>
  shopify inventory levels help
  shopify inventory levels search --first <count>
`)
}

func (cli CLI) printInventoryItemsHelp() {
	fmt.Fprintf(cli.stdout, `shopify inventory items

%s

Usage:
  shopify inventory items help
  shopify inventory items search --first <count>
  shopify inventory items search --first <count> --query <search>
  shopify inventory items search --first <count> --after <cursor>
`, shopifyInventoryItemsSearchDoc.Description)
}

func (cli CLI) printInventoryLevelsHelp() {
	fmt.Fprintf(cli.stdout, `shopify inventory levels

%s

Usage:
  shopify inventory levels help
  shopify inventory levels search --first <count>
  shopify inventory levels search --first <count> --query <search>
  shopify inventory levels search --first <count> --after <cursor>
`, shopifyInventoryLevelsSearchDoc.Description)
}

func (cli CLI) printLocationsHelp() {
	fmt.Fprint(cli.stdout, `shopify locations

Work with Shopify locations.

Usage:
  shopify locations help
  shopify locations list --first <count>
  shopify locations list --first <count> --after <cursor>
`)
}

func (cli CLI) printLocationsListHelp() {
	fmt.Fprintf(cli.stdout, `shopify locations list

%s

Usage:
  shopify locations list --first <count>
  shopify locations list --first <count> --after <cursor>
`, shopifyLocationsListDoc.Description)
}
