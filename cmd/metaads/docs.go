package main

type commandDoc struct {
	Command     string
	Description string
}

var metaAdsAuthTestDoc = commandDoc{
	Command:     "metaads auth test",
	Description: "Call /me through the configured Meta Graph API base URL to verify the injected access token.",
}

var metaAdsGraphRequestDoc = commandDoc{
	Command:     "metaads graph request",
	Description: "Send a raw Meta Graph API request. This is the complete Meta Ads API coverage surface.",
}

var metaAdsAdAccountsListDoc = commandDoc{
	Command:     "metaads ad-accounts list",
	Description: "List ad accounts visible to the current Meta access token.",
}

var metaAdsAdAccountGetDoc = commandDoc{
	Command:     "metaads ad-accounts get",
	Description: "Get one Meta ad account by ID.",
}

var metaAdsCampaignsSearchDoc = commandDoc{
	Command:     "metaads campaigns search",
	Description: "List campaigns for an ad account.",
}

var metaAdsCampaignGetDoc = commandDoc{
	Command:     "metaads campaigns get",
	Description: "Get one Meta campaign by ID.",
}

var metaAdsCampaignCreateDoc = commandDoc{
	Command:     "metaads campaigns create",
	Description: "Create a Meta campaign using the documented Graph API request body.",
}

var metaAdsCampaignUpdateDoc = commandDoc{
	Command:     "metaads campaigns update",
	Description: "Update a Meta campaign using the documented Graph API request body.",
}

var metaAdsCampaignDeleteDoc = commandDoc{
	Command:     "metaads campaigns delete",
	Description: "Delete one Meta campaign by ID.",
}

var metaAdsAdSetsSearchDoc = commandDoc{
	Command:     "metaads adsets search",
	Description: "List ad sets for an ad account.",
}

var metaAdsAdSetGetDoc = commandDoc{
	Command:     "metaads adsets get",
	Description: "Get one Meta ad set by ID.",
}

var metaAdsAdSetCreateDoc = commandDoc{
	Command:     "metaads adsets create",
	Description: "Create a Meta ad set using the documented Graph API request body.",
}

var metaAdsAdSetUpdateDoc = commandDoc{
	Command:     "metaads adsets update",
	Description: "Update a Meta ad set using the documented Graph API request body.",
}

var metaAdsAdSetDeleteDoc = commandDoc{
	Command:     "metaads adsets delete",
	Description: "Delete one Meta ad set by ID.",
}

var metaAdsAdsSearchDoc = commandDoc{
	Command:     "metaads ads search",
	Description: "List ads for an ad account.",
}

var metaAdsAdGetDoc = commandDoc{
	Command:     "metaads ads get",
	Description: "Get one Meta ad by ID.",
}

var metaAdsAdCreateDoc = commandDoc{
	Command:     "metaads ads create",
	Description: "Create a Meta ad using the documented Graph API request body.",
}

var metaAdsAdUpdateDoc = commandDoc{
	Command:     "metaads ads update",
	Description: "Update a Meta ad using the documented Graph API request body.",
}

var metaAdsAdDeleteDoc = commandDoc{
	Command:     "metaads ads delete",
	Description: "Delete one Meta ad by ID.",
}

var metaAdsCreativesSearchDoc = commandDoc{
	Command:     "metaads creatives search",
	Description: "List ad creatives for an ad account.",
}

var metaAdsCreativeGetDoc = commandDoc{
	Command:     "metaads creatives get",
	Description: "Get one Meta ad creative by ID.",
}

var metaAdsCreativeCreateDoc = commandDoc{
	Command:     "metaads creatives create",
	Description: "Create a Meta ad creative using the documented Graph API request body.",
}

var metaAdsInsightsGetDoc = commandDoc{
	Command:     "metaads insights get",
	Description: "Get Meta Ads insights from an ad account or object insights edge.",
}

var metaAdsTargetingSearchDoc = commandDoc{
	Command:     "metaads targeting search",
	Description: "Search Meta targeting options through the targetingsearch endpoint.",
}
