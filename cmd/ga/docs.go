package main

type commandDoc struct {
	Command     string
	Summary     string
	Description string
}

var gaAuthTestDoc = commandDoc{
	Command:     "ga auth test",
	Summary:     "Check Google Analytics API access",
	Description: "Check Google Analytics API access by fetching a configured property.",
}

var gaAccountSummariesListDoc = commandDoc{
	Command:     "ga account-summaries list",
	Summary:     "List Google Analytics account summaries",
	Description: "List Google Analytics accounts and properties visible to the caller.",
}

var gaPropertyGetDoc = commandDoc{
	Command:     "ga properties get",
	Summary:     "Get a Google Analytics property",
	Description: "Get details for a Google Analytics property.",
}

var gaMetadataGetDoc = commandDoc{
	Command:     "ga metadata get",
	Summary:     "Get property metadata",
	Description: "Get Google Analytics dimensions and metrics metadata for a property.",
}

var gaCompatibilityCheckDoc = commandDoc{
	Command:     "ga compatibility check",
	Summary:     "Check report dimension and metric compatibility",
	Description: "Check whether dimensions and metrics can be used together in a report request.",
}

var gaReportRunDoc = commandDoc{
	Command:     "ga reports run",
	Summary:     "Run a core report",
	Description: "Run a Google Analytics Data API core report.",
}

var gaReportRealtimeDoc = commandDoc{
	Command:     "ga reports realtime",
	Summary:     "Run a realtime report",
	Description: "Run a Google Analytics Data API realtime report.",
}

var gaReportFunnelDoc = commandDoc{
	Command:     "ga reports funnel",
	Summary:     "Run a funnel report",
	Description: "Run a Google Analytics Data API funnel report.",
}

var gaGoogleAdsLinksListDoc = commandDoc{
	Command:     "ga google-ads-links list",
	Summary:     "List Google Ads links",
	Description: "List Google Ads links for a Google Analytics property.",
}
