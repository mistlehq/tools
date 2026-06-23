package main

import (
	"strings"
	"testing"
)

func TestAccountsCommands(t *testing.T) {
	env, simulator := setupCommandEnvironment(t)

	listResult, err := runCommandWithInput(t, env, "", "gbp", "accounts", "list")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(listResult.stdout.String(), testAccount) {
		t.Fatalf("expected accounts list output to include %q, got %q", testAccount, listResult.stdout.String())
	}

	getResult, err := runCommandWithInput(t, env, "", "gbp", "accounts", "get", "--account", testAccount)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(getResult.stdout.String(), "Mistle Test") {
		t.Fatalf("expected account get output to include account name, got %q", getResult.stdout.String())
	}
	if !simulator.sawPath("/v1/accounts/123") {
		t.Fatalf("expected simulator to receive escaped account resource path, got %#v", simulator.paths)
	}
}

func TestLocationsCommands(t *testing.T) {
	env, simulator := setupCommandEnvironment(t)

	listResult, err := runCommandWithInput(t, env, "", "gbp", "locations", "list", "--account", testAccount, "--read-mask", "name,title,storeCode")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(listResult.stdout.String(), "Mistle HQ") {
		t.Fatalf("expected locations list output to include location title, got %q", listResult.stdout.String())
	}
	if !simulator.sawPath("/v1/accounts/123/locations?readMask=name%2Ctitle%2CstoreCode") {
		t.Fatalf("expected simulator to receive locations list path with readMask, got %#v", simulator.paths)
	}

	getResult, err := runCommandWithInput(t, env, "", "gbp", "locations", "get", "--location", testLocation, "--read-mask", "name,title")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(getResult.stdout.String(), testLocation) {
		t.Fatalf("expected locations get output to include location name, got %q", getResult.stdout.String())
	}

	locationRequest := writeTempJSONRequest(t, `{"title":"Mistle HQ"}`)
	createResult, err := runCommandWithInput(t, env, "", "gbp", "locations", "create", "--account", testAccount, "--request-file", locationRequest, "--request-id", "req-1", "--validate-only", "true")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(createResult.stdout.String(), "Created Mistle HQ") {
		t.Fatalf("expected create output to include created title, got %q", createResult.stdout.String())
	}
	if !simulator.sawPath("/v1/accounts/123/locations?requestId=req-1&validateOnly=true") {
		t.Fatalf("expected simulator to receive locations create path, got %#v", simulator.paths)
	}

	patchResult, err := runCommandWithInput(t, env, "", "gbp", "locations", "patch", "--location", testLocation, "--update-mask", "title", "--request-file", locationRequest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(patchResult.stdout.String(), "Patched Mistle HQ") {
		t.Fatalf("expected patch output to include patched title, got %q", patchResult.stdout.String())
	}
	if !simulator.sawPath("/v1/locations/456?updateMask=title") {
		t.Fatalf("expected simulator to receive locations patch path, got %#v", simulator.paths)
	}

	if _, err := runCommandWithInput(t, env, "", "gbp", "locations", "delete", "--location", testLocation); err != nil {
		t.Fatal(err)
	}
}

func TestReviewsMediaAndLocalPostsCommands(t *testing.T) {
	env, simulator := setupCommandEnvironment(t)

	reviewsResult, err := runCommandWithInput(t, env, "", "gbp", "reviews", "list", "--account", testAccount, "--location", testLocation, "--page-size", "10", "--order-by", "updateTime desc")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reviewsResult.stdout.String(), `"reviewId": "abc"`) {
		t.Fatalf("expected reviews output to include review ID, got %q", reviewsResult.stdout.String())
	}
	if !simulator.sawPath("/v4/accounts/123/locations/456/reviews?orderBy=updateTime+desc&pageSize=10") {
		t.Fatalf("expected simulator to receive reviews list query path, got %#v", simulator.paths)
	}

	reviewResult, err := runCommandWithInput(t, env, "", "gbp", "reviews", "get", "--account", testAccount, "--location", testLocation, "--review", "abc")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reviewResult.stdout.String(), `"starRating": "FIVE"`) {
		t.Fatalf("expected review output to include star rating, got %q", reviewResult.stdout.String())
	}

	replyRequest := writeTempJSONRequest(t, `{"comment":"Thanks!"}`)
	replyResult, err := runCommandWithInput(t, env, "", "gbp", "reviews", "update-reply", "--account", testAccount, "--location", testLocation, "--review", "abc", "--request-file", replyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(replyResult.stdout.String(), `"comment": "Thanks!"`) {
		t.Fatalf("expected review reply output to include comment, got %q", replyResult.stdout.String())
	}
	if !simulator.sawPath("/v4/accounts/123/locations/456/reviews/abc/reply") {
		t.Fatalf("expected simulator to receive review reply path, got %#v", simulator.paths)
	}

	if _, err := runCommandWithInput(t, env, "", "gbp", "reviews", "delete-reply", "--account", testAccount, "--location", testLocation, "--review", "abc"); err != nil {
		t.Fatal(err)
	}

	mediaResult, err := runCommandWithInput(t, env, "", "gbp", "media", "list", "--account", testAccount, "--location", testLocation)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(mediaResult.stdout.String(), "mediaItems") {
		t.Fatalf("expected media output to include mediaItems, got %q", mediaResult.stdout.String())
	}

	mediaRequest := writeTempJSONRequest(t, `{"mediaFormat":"PHOTO"}`)
	mediaCreateResult, err := runCommandWithInput(t, env, "", "gbp", "media", "create", "--account", testAccount, "--location", testLocation, "--request-file", mediaRequest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(mediaCreateResult.stdout.String(), "media/789") {
		t.Fatalf("expected media create output to include media name, got %q", mediaCreateResult.stdout.String())
	}
	if _, err := runCommandWithInput(t, env, "", "gbp", "media", "get", "--media", "accounts/123/locations/456/media/789"); err != nil {
		t.Fatal(err)
	}
	if _, err := runCommandWithInput(t, env, "", "gbp", "media", "patch", "--media", "accounts/123/locations/456/media/789", "--update-mask", "description", "--request-file", mediaRequest); err != nil {
		t.Fatal(err)
	}
	if _, err := runCommandWithInput(t, env, "", "gbp", "media", "delete", "--media", "accounts/123/locations/456/media/789"); err != nil {
		t.Fatal(err)
	}
	uploadResult, err := runCommandWithInput(t, env, "", "gbp", "media", "start-upload", "--account", testAccount, "--location", testLocation)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(uploadResult.stdout.String(), "upload-ref") {
		t.Fatalf("expected media start-upload output to include upload ref, got %q", uploadResult.stdout.String())
	}

	postsResult, err := runCommandWithInput(t, env, "", "gbp", "local-posts", "list", "--account", testAccount, "--location", testLocation)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(postsResult.stdout.String(), "localPosts") {
		t.Fatalf("expected local posts output to include localPosts, got %q", postsResult.stdout.String())
	}

	postRequest := writeTempJSONRequest(t, `{"summary":"Hello"}`)
	postCreateResult, err := runCommandWithInput(t, env, "", "gbp", "local-posts", "create", "--account", testAccount, "--location", testLocation, "--request-file", postRequest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(postCreateResult.stdout.String(), "localPosts/post-1") {
		t.Fatalf("expected local post create output to include post name, got %q", postCreateResult.stdout.String())
	}
	if _, err := runCommandWithInput(t, env, "", "gbp", "local-posts", "get", "--local-post", "accounts/123/locations/456/localPosts/post-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := runCommandWithInput(t, env, "", "gbp", "local-posts", "patch", "--local-post", "accounts/123/locations/456/localPosts/post-1", "--update-mask", "summary", "--request-file", postRequest); err != nil {
		t.Fatal(err)
	}
	if _, err := runCommandWithInput(t, env, "", "gbp", "local-posts", "delete", "--local-post", "accounts/123/locations/456/localPosts/post-1"); err != nil {
		t.Fatal(err)
	}
	insightsResult, err := runCommandWithInput(t, env, "", "gbp", "local-posts", "report-insights", "--account", testAccount, "--location", testLocation, "--request-file", postRequest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(insightsResult.stdout.String(), "localPostMetrics") {
		t.Fatalf("expected local post insights output to include localPostMetrics, got %q", insightsResult.stdout.String())
	}
}

func TestPerformanceCommands(t *testing.T) {
	env, simulator := setupCommandEnvironment(t)

	dailyRequest := writeTempJSONRequest(t, `{"dailyMetric":"WEBSITE_CLICKS","dailyRange":{"start_date":{"year":2026,"month":6,"day":1},"end_date":{"year":2026,"month":6,"day":2}}}`)
	dailyResult, err := runCommandWithInput(t, env, "", "gbp", "performance", "daily-metrics", "--location", testLocation, "--request-file", dailyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(dailyResult.stdout.String(), "timeSeries") {
		t.Fatalf("expected daily metrics output to include timeSeries, got %q", dailyResult.stdout.String())
	}
	if !simulator.sawPath("/v1/locations/456:getDailyMetricsTimeSeries?dailyMetric=WEBSITE_CLICKS&dailyRange.end_date.day=2&dailyRange.end_date.month=6&dailyRange.end_date.year=2026&dailyRange.start_date.day=1&dailyRange.start_date.month=6&dailyRange.start_date.year=2026") {
		t.Fatalf("expected simulator to receive flattened daily metrics query path, got %#v", simulator.paths)
	}

	keywordsRequest := writeTempJSONRequest(t, `{"monthlyRange":{"start_month":{"year":2026,"month":5},"end_month":{"year":2026,"month":6}},"pageSize":10}`)
	keywordsResult, err := runCommandWithInput(t, env, "", "gbp", "performance", "search-keywords", "--location", testLocation, "--request-file", keywordsRequest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(keywordsResult.stdout.String(), "searchKeywordsCounts") {
		t.Fatalf("expected search keywords output to include searchKeywordsCounts, got %q", keywordsResult.stdout.String())
	}
}
