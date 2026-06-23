package main

import "fmt"

func (cli CLI) printHelp() {
	fmt.Fprint(cli.stdout, `googleads

Thin CLI for Google Ads API REST endpoints.

Usage:
  googleads help
  googleads version
  googleads auth help
  googleads request --help
  googleads customers help
  googleads gaql help
  googleads fields help
  googleads mcp help

Commands:
  help
  version
  auth
  request
  customers
  gaql
  fields
  mcp
`)
}

func (cli CLI) printAuthHelp() {
	fmt.Fprint(cli.stdout, `googleads auth

Inspect Google Ads API authentication state.

Usage:
  googleads auth help
  googleads auth test
`)
}

func (cli CLI) printAuthTestHelp() {
	fmt.Fprintf(cli.stdout, `googleads auth test

%s

Usage:
  googleads auth test
`, googleAdsAuthTestDoc.Description)
}

func (cli CLI) printRequestHelp() {
	fmt.Fprintf(cli.stdout, `googleads request

%s

Usage:
  googleads request --method GET --path /customers:listAccessibleCustomers
  googleads request --method POST --path /customers/<customer-id>/googleAds:search --body '{"query":"SELECT customer.id FROM customer LIMIT 1"}'
  googleads request --method POST --path /customers/<customer-id>/campaigns:mutate --body <json>

Options:
  --method <method>       HTTP method: GET, POST, PATCH, or DELETE. Defaults to GET.
  --path <path>           Google Ads API path under the configured version base, starting with '/'.
  --params <json>         Query parameters JSON object.
  --params-file <path>    File containing query parameters JSON object.
  --body <json>           JSON request body.
  --body-file <path>      File containing JSON request body.
`, googleAdsRequestDoc.Description)
}

func (cli CLI) printCustomersHelp() {
	fmt.Fprint(cli.stdout, `googleads customers

Work with Google Ads customers.

Usage:
  googleads customers help
  googleads customers list-accessible
`)
}

func (cli CLI) printCustomersListAccessibleHelp() {
	fmt.Fprintf(cli.stdout, `googleads customers list-accessible

%s

Usage:
  googleads customers list-accessible
`, googleAdsCustomersListAccessibleDoc.Description)
}

func (cli CLI) printGAQLHelp() {
	fmt.Fprint(cli.stdout, `googleads gaql

Run Google Ads Query Language requests.

Usage:
  googleads gaql help
  googleads gaql search --customer-id <customer-id> --query <gaql>
  googleads gaql search-stream --customer-id <customer-id> --query <gaql>
`)
}

func (cli CLI) printGAQLSearchHelp() {
	fmt.Fprintf(cli.stdout, `googleads gaql search

%s

Usage:
  googleads gaql search --customer-id <customer-id> --query 'SELECT customer.id FROM customer LIMIT 1'
  googleads gaql search --customer-id <customer-id> --query-file <path>
  googleads gaql search --customer-id <customer-id> --query <gaql> --page-size <n> --page-token <token>
`, googleAdsGAQLSearchDoc.Description)
}

func (cli CLI) printGAQLSearchStreamHelp() {
	fmt.Fprintf(cli.stdout, `googleads gaql search-stream

%s

Usage:
  googleads gaql search-stream --customer-id <customer-id> --query 'SELECT customer.id FROM customer LIMIT 1'
  googleads gaql search-stream --customer-id <customer-id> --query-file <path>
`, googleAdsGAQLSearchStreamDoc.Description)
}

func (cli CLI) printFieldsHelp() {
	fmt.Fprint(cli.stdout, `googleads fields

Work with Google Ads API field metadata.

Usage:
  googleads fields help
  googleads fields search --query <gaql>
  googleads fields get --resource-name googleAdsFields/<field-name>
`)
}

func (cli CLI) printFieldsSearchHelp() {
	fmt.Fprintf(cli.stdout, `googleads fields search

%s

Usage:
  googleads fields search --query 'SELECT name, category, data_type WHERE name = "campaign.id"'
  googleads fields search --query-file <path>
`, googleAdsFieldsSearchDoc.Description)
}

func (cli CLI) printFieldGetHelp() {
	fmt.Fprintf(cli.stdout, `googleads fields get

%s

Usage:
  googleads fields get --resource-name googleAdsFields/campaign.id
`, googleAdsFieldGetDoc.Description)
}
