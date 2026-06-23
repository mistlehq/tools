package main

type commandDoc struct {
	Command     string
	Summary     string
	Description string
}

var gbpAuthTestDoc = commandDoc{
	Command:     "gbp auth test",
	Summary:     "Check Google Business Profile API access",
	Description: "Check Google Business Profile API access by listing visible accounts.",
}

var gbpAccountsListDoc = commandDoc{
	Command:     "gbp accounts list",
	Summary:     "List Business Profile accounts",
	Description: "List Google Business Profile accounts visible to the caller.",
}

var gbpAccountGetDoc = commandDoc{
	Command:     "gbp accounts get",
	Summary:     "Get a Business Profile account",
	Description: "Get details for a Google Business Profile account.",
}

var gbpLocationsListDoc = commandDoc{
	Command:     "gbp locations list",
	Summary:     "List Business Profile locations",
	Description: "List locations under a Google Business Profile account using Google's documented readMask query parameter.",
}

var gbpLocationGetDoc = commandDoc{
	Command:     "gbp locations get",
	Summary:     "Get a Business Profile location",
	Description: "Get a Google Business Profile location using Google's documented readMask query parameter.",
}

var gbpLocationCreateDoc = commandDoc{
	Command:     "gbp locations create",
	Summary:     "Create a Business Profile location",
	Description: "Create a Google Business Profile location using Google's documented Location request shape.",
}

var gbpLocationPatchDoc = commandDoc{
	Command:     "gbp locations patch",
	Summary:     "Patch a Business Profile location",
	Description: "Patch a Google Business Profile location using Google's documented updateMask and Location request shape.",
}

var gbpLocationDeleteDoc = commandDoc{
	Command:     "gbp locations delete",
	Summary:     "Delete a Business Profile location",
	Description: "Delete a Google Business Profile location.",
}

var gbpReviewsListDoc = commandDoc{
	Command:     "gbp reviews list",
	Summary:     "List location reviews",
	Description: "List Google Business Profile reviews for a location.",
}

var gbpReviewGetDoc = commandDoc{
	Command:     "gbp reviews get",
	Summary:     "Get a location review",
	Description: "Get a Google Business Profile review for a location.",
}

var gbpReviewUpdateReplyDoc = commandDoc{
	Command:     "gbp reviews update-reply",
	Summary:     "Update a review reply",
	Description: "Create or update a Google Business Profile review reply using Google's documented ReviewReply request shape.",
}

var gbpReviewDeleteReplyDoc = commandDoc{
	Command:     "gbp reviews delete-reply",
	Summary:     "Delete a review reply",
	Description: "Delete a Google Business Profile review reply.",
}

var gbpMediaListDoc = commandDoc{
	Command:     "gbp media list",
	Summary:     "List location media",
	Description: "List Google Business Profile media items for a location.",
}

var gbpMediaCreateDoc = commandDoc{
	Command:     "gbp media create",
	Summary:     "Create a media item",
	Description: "Create a Google Business Profile media item using Google's documented MediaItem request shape.",
}

var gbpMediaGetDoc = commandDoc{
	Command:     "gbp media get",
	Summary:     "Get a media item",
	Description: "Get a Google Business Profile media item.",
}

var gbpMediaPatchDoc = commandDoc{
	Command:     "gbp media patch",
	Summary:     "Patch a media item",
	Description: "Patch a Google Business Profile media item using Google's documented updateMask and MediaItem request shape.",
}

var gbpMediaDeleteDoc = commandDoc{
	Command:     "gbp media delete",
	Summary:     "Delete a media item",
	Description: "Delete a Google Business Profile media item.",
}

var gbpMediaStartUploadDoc = commandDoc{
	Command:     "gbp media start-upload",
	Summary:     "Start media upload",
	Description: "Start a Google Business Profile media upload and return Google's media item data reference.",
}

var gbpLocalPostsListDoc = commandDoc{
	Command:     "gbp local-posts list",
	Summary:     "List local posts",
	Description: "List Google Business Profile local posts for a location.",
}

var gbpLocalPostCreateDoc = commandDoc{
	Command:     "gbp local-posts create",
	Summary:     "Create a local post",
	Description: "Create a Google Business Profile local post using Google's documented LocalPost request shape.",
}

var gbpLocalPostGetDoc = commandDoc{
	Command:     "gbp local-posts get",
	Summary:     "Get a local post",
	Description: "Get a Google Business Profile local post.",
}

var gbpLocalPostPatchDoc = commandDoc{
	Command:     "gbp local-posts patch",
	Summary:     "Patch a local post",
	Description: "Patch a Google Business Profile local post using Google's documented updateMask and LocalPost request shape.",
}

var gbpLocalPostDeleteDoc = commandDoc{
	Command:     "gbp local-posts delete",
	Summary:     "Delete a local post",
	Description: "Delete a Google Business Profile local post.",
}

var gbpLocalPostReportInsightsDoc = commandDoc{
	Command:     "gbp local-posts report-insights",
	Summary:     "Report local post insights",
	Description: "Report Google Business Profile local post insights using Google's documented request shape.",
}

var gbpPerformanceDailyMetricsDoc = commandDoc{
	Command:     "gbp performance daily-metrics",
	Summary:     "Fetch daily performance metrics",
	Description: "Fetch Google Business Profile daily metric time series with Google's documented query parameter shape.",
}

var gbpPerformanceSearchKeywordsDoc = commandDoc{
	Command:     "gbp performance search-keywords",
	Summary:     "Fetch monthly search keywords",
	Description: "Fetch Google Business Profile monthly search keyword impressions with Google's documented query parameter shape.",
}
