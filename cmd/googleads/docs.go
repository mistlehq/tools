package main

type commandDoc struct {
	Command     string
	Description string
}

var googleAdsAuthTestDoc = commandDoc{
	Command:     "googleads auth test",
	Description: "Call customers:listAccessibleCustomers through the configured Google Ads API base URL to verify injected credentials.",
}

var googleAdsRequestDoc = commandDoc{
	Command:     "googleads request",
	Description: "Send a raw Google Ads API REST request. This is the complete Google Ads API coverage surface.",
}

var googleAdsCustomersListAccessibleDoc = commandDoc{
	Command:     "googleads customers list-accessible",
	Description: "List Google Ads customers accessible to the current OAuth token and developer token.",
}

var googleAdsGAQLSearchDoc = commandDoc{
	Command:     "googleads gaql search",
	Description: "Run a Google Ads Query Language search request for one customer.",
}

var googleAdsGAQLSearchStreamDoc = commandDoc{
	Command:     "googleads gaql search-stream",
	Description: "Run a Google Ads Query Language searchStream request for one customer.",
}

var googleAdsFieldsSearchDoc = commandDoc{
	Command:     "googleads fields search",
	Description: "Search GoogleAdsField metadata with GAQL.",
}

var googleAdsFieldGetDoc = commandDoc{
	Command:     "googleads fields get",
	Description: "Get one GoogleAdsField metadata resource by resource name.",
}
