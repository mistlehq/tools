package main

import "fmt"

func (cli CLI) printHelp() {
	fmt.Fprint(cli.stdout, `metaads

Thin CLI for Meta Graph API / Marketing API.

Usage:
  metaads help
  metaads version
  metaads auth help
  metaads graph help
  metaads ad-accounts help
  metaads campaigns help
  metaads adsets help
  metaads ads help
  metaads creatives help
  metaads insights help
  metaads targeting help
  metaads mcp help

Commands:
  help
  version
  auth
  graph
  ad-accounts
  campaigns
  adsets
  ads
  creatives
  insights
  targeting
  mcp
`)
}

func (cli CLI) printAuthHelp() {
	fmt.Fprint(cli.stdout, `metaads auth

Inspect Meta Graph API authentication state.

Usage:
  metaads auth help
  metaads auth test
`)
}

func (cli CLI) printAuthTestHelp() {
	fmt.Fprintf(cli.stdout, `metaads auth test

%s

Usage:
  metaads auth test
`, metaAdsAuthTestDoc.Description)
}

func (cli CLI) printGraphHelp() {
	fmt.Fprint(cli.stdout, `metaads graph

Send raw Meta Graph API requests.

Usage:
  metaads graph help
  metaads graph request --help
`)
}

func (cli CLI) printGraphRequestHelp() {
	fmt.Fprintf(cli.stdout, `metaads graph request

%s

Usage:
  metaads graph request --method GET --path /me
  metaads graph request --method GET --path /act_<account-id>/campaigns --params '{"fields":"id,name"}'
  metaads graph request --method POST --path /act_<account-id>/campaigns --body '{"name":"Example"}'

Options:
  --method <method>       HTTP method: GET, POST, or DELETE. Defaults to GET.
  --path <path>           Graph API path, starting with '/'.
  --params <json>         Query parameters JSON object.
  --params-file <path>    File containing query parameters JSON object.
  --body <json>           JSON request body for POST requests.
  --body-file <path>      File containing JSON request body.
`, metaAdsGraphRequestDoc.Description)
}

func (cli CLI) printAdAccountsHelp() {
	fmt.Fprint(cli.stdout, `metaads ad-accounts

Work with Meta ad accounts.

Usage:
  metaads ad-accounts help
  metaads ad-accounts list
  metaads ad-accounts get --id act_<account-id>
`)
}

func (cli CLI) printAdAccountsListHelp() {
	fmt.Fprintf(cli.stdout, `metaads ad-accounts list

%s

Usage:
  metaads ad-accounts list
  metaads ad-accounts list --fields <fields>
  metaads ad-accounts list --limit <limit> --after <cursor>
`, metaAdsAdAccountsListDoc.Description)
}

func (cli CLI) printAdAccountGetHelp() {
	fmt.Fprintf(cli.stdout, `metaads ad-accounts get

%s

Usage:
  metaads ad-accounts get --id act_<account-id>
  metaads ad-accounts get --id act_<account-id> --fields <fields>
`, metaAdsAdAccountGetDoc.Description)
}

func (cli CLI) printCampaignsHelp() {
	cli.printCrudHelp("campaigns", "campaign")
}

func (cli CLI) printAdSetsHelp() {
	cli.printCrudHelp("adsets", "ad set")
}

func (cli CLI) printAdsHelp() {
	cli.printCrudHelp("ads", "ad")
}

func (cli CLI) printCrudHelp(name string, label string) {
	fmt.Fprintf(cli.stdout, `metaads %s

Work with Meta %ss.

Usage:
  metaads %s help
  metaads %s search --account-id act_<account-id>
  metaads %s get --id <id>
  metaads %s create --account-id act_<account-id> --body <json>
  metaads %s update --id <id> --body <json>
  metaads %s delete --id <id>
`, name, label, name, name, name, name, name, name)
}

func (cli CLI) printSearchHelp(name string, doc commandDoc) {
	fmt.Fprintf(cli.stdout, `metaads %s search

%s

Usage:
  metaads %s search --account-id act_<account-id>
  metaads %s search --account-id act_<account-id> --fields <fields>
  metaads %s search --account-id act_<account-id> --limit <limit> --after <cursor>
  metaads %s search --account-id act_<account-id> --params <json>
`, name, doc.Description, name, name, name, name)
}

func (cli CLI) printGetHelp(name string, doc commandDoc) {
	fmt.Fprintf(cli.stdout, `metaads %s get

%s

Usage:
  metaads %s get --id <id>
  metaads %s get --id <id> --fields <fields>
  metaads %s get --id <id> --params <json>
`, name, doc.Description, name, name, name)
}

func (cli CLI) printCreateHelp(name string, doc commandDoc) {
	fmt.Fprintf(cli.stdout, `metaads %s create

%s

Usage:
  metaads %s create --account-id act_<account-id> --body <json>
  metaads %s create --account-id act_<account-id> --body-file <path>
`, name, doc.Description, name, name)
}

func (cli CLI) printUpdateHelp(name string, doc commandDoc) {
	fmt.Fprintf(cli.stdout, `metaads %s update

%s

Usage:
  metaads %s update --id <id> --body <json>
  metaads %s update --id <id> --body-file <path>
`, name, doc.Description, name, name)
}

func (cli CLI) printDeleteHelp(name string, doc commandDoc) {
	fmt.Fprintf(cli.stdout, `metaads %s delete

%s

Usage:
  metaads %s delete --id <id>
`, name, doc.Description, name)
}

func (cli CLI) printCreativesHelp() {
	fmt.Fprint(cli.stdout, `metaads creatives

Work with Meta ad creatives.

Usage:
  metaads creatives help
  metaads creatives search --account-id act_<account-id>
  metaads creatives get --id <creative-id>
  metaads creatives create --account-id act_<account-id> --body <json>
`)
}

func (cli CLI) printCreativesSearchHelp() {
	cli.printSearchHelp("creatives", metaAdsCreativesSearchDoc)
}

func (cli CLI) printCreativeGetHelp() {
	cli.printGetHelp("creatives", metaAdsCreativeGetDoc)
}

func (cli CLI) printCreativeCreateHelp() {
	cli.printCreateHelp("creatives", metaAdsCreativeCreateDoc)
}

func (cli CLI) printInsightsHelp() {
	fmt.Fprint(cli.stdout, `metaads insights

Work with Meta Ads insights.

Usage:
  metaads insights help
  metaads insights get --id act_<account-id>
  metaads insights get --id <campaign-or-adset-or-ad-id>
`)
}

func (cli CLI) printInsightsGetHelp() {
	fmt.Fprintf(cli.stdout, `metaads insights get

%s

Usage:
  metaads insights get --id act_<account-id>
  metaads insights get --id act_<account-id> --fields impressions,clicks,spend
  metaads insights get --id act_<account-id> --level campaign
  metaads insights get --id act_<account-id> --time-range '{"since":"2026-06-01","until":"2026-06-23"}'
  metaads insights get --id act_<account-id> --params <json>
`, metaAdsInsightsGetDoc.Description)
}

func (cli CLI) printTargetingHelp() {
	fmt.Fprint(cli.stdout, `metaads targeting

Search Meta targeting options.

Usage:
  metaads targeting help
  metaads targeting search --type adinterest --query running
`)
}

func (cli CLI) printTargetingSearchHelp() {
	fmt.Fprintf(cli.stdout, `metaads targeting search

%s

Usage:
  metaads targeting search --type adinterest --query running
  metaads targeting search --type geo_location --params <json>
`, metaAdsTargetingSearchDoc.Description)
}
