package main

type commandDoc struct {
	Command     string
	Summary     string
	Description string
}

var gscAuthTestDoc = commandDoc{
	Command:     "gsc auth test",
	Summary:     "Check Google Search Console API access",
	Description: "Check Google Search Console API access by fetching a configured site.",
}

var gscSitesListDoc = commandDoc{
	Command:     "gsc sites list",
	Summary:     "List Search Console sites",
	Description: "List Search Console sites visible to the caller.",
}

var gscSiteGetDoc = commandDoc{
	Command:     "gsc sites get",
	Summary:     "Get a Search Console site",
	Description: "Get details for a Search Console site.",
}

var gscSearchAnalyticsQueryDoc = commandDoc{
	Command:     "gsc searchanalytics query",
	Summary:     "Query Search Analytics data",
	Description: "Run a Search Console Search Analytics query using Google's documented request shape.",
}

var gscSitemapsListDoc = commandDoc{
	Command:     "gsc sitemaps list",
	Summary:     "List sitemaps",
	Description: "List sitemaps for a Search Console site.",
}

var gscSitemapGetDoc = commandDoc{
	Command:     "gsc sitemaps get",
	Summary:     "Get a sitemap",
	Description: "Get details for a sitemap in a Search Console site.",
}

var gscURLInspectionInspectDoc = commandDoc{
	Command:     "gsc url-inspection inspect",
	Summary:     "Inspect a URL",
	Description: "Inspect a URL using the Search Console URL Inspection API.",
}
