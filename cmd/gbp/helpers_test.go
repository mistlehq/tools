package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	testAccount  = "accounts/123"
	testLocation = "locations/456"
)

type commandResult struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func runCommandWithInput(t *testing.T, env Environment, input string, args ...string) (commandResult, error) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cli := CLI{
		stdin:  bytes.NewBufferString(input),
		stdout: &stdout,
		stderr: &stderr,
		env:    env,
	}

	err := cli.run(args)
	return commandResult{
		stdout: stdout,
		stderr: stderr,
	}, err
}

func setupCommandEnvironment(t *testing.T) (Environment, *simulatedGBPAPI) {
	t.Helper()

	simulator := startSimulatedGBPAPI(t)
	return Environment{
		"GBP_ACCOUNT_MANAGEMENT_BASE_URL":   simulator.URL,
		"GBP_BUSINESS_INFORMATION_BASE_URL": simulator.URL,
		"GBP_PERFORMANCE_BASE_URL":          simulator.URL,
		"GBP_MYBUSINESS_BASE_URL":           simulator.URL,
	}, simulator
}

func setupGBPClient(t *testing.T) (Environment, GBPClient, *simulatedGBPAPI) {
	t.Helper()

	env, simulator := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	return env, NewGBPClient(config), simulator
}

func writeTempJSONRequest(t *testing.T, body string) string {
	t.Helper()
	path := fmt.Sprintf("%s/request.json", t.TempDir())
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func validUnitEnv() Environment {
	return Environment{
		"GBP_ACCOUNT_MANAGEMENT_BASE_URL":   "http://127.0.0.1",
		"GBP_BUSINESS_INFORMATION_BASE_URL": "http://127.0.0.1",
		"GBP_PERFORMANCE_BASE_URL":          "http://127.0.0.1",
		"GBP_MYBUSINESS_BASE_URL":           "http://127.0.0.1",
	}
}

type simulatedGBPAPI struct {
	*httptest.Server
	paths []string
}

func startSimulatedGBPAPI(t *testing.T) *simulatedGBPAPI {
	t.Helper()

	simulator := &simulatedGBPAPI{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", simulator.handle)
	server := httptest.NewServer(mux)
	simulator.Server = server
	t.Cleanup(server.Close)
	return simulator
}

func (api *simulatedGBPAPI) handle(writer http.ResponseWriter, request *http.Request) {
	api.paths = append(api.paths, request.URL.RequestURI())
	writer.Header().Set("Content-Type", "application/json")

	switch request.URL.Path {
	case "/v1/accounts":
		// Official method: GET https://mybusinessaccountmanagement.googleapis.com/v1/accounts
		writeSimulatedJSON(writer, GBPAccountsList{Accounts: []GBPAccount{{
			Name:        testAccount,
			AccountName: "Mistle Test",
			Type:        "LOCATION_GROUP",
			Role:        "OWNER",
		}}})
	case "/v1/accounts/123":
		// Official method: GET https://mybusinessaccountmanagement.googleapis.com/v1/{name=accounts/*}
		writeSimulatedJSON(writer, GBPAccount{Name: testAccount, AccountName: "Mistle Test", Type: "LOCATION_GROUP", Role: "OWNER"})
	case "/v1/accounts/123/locations":
		if request.Method == http.MethodPost {
			// Official method: POST https://mybusinessbusinessinformation.googleapis.com/v1/{parent=accounts/*}/locations
			writeSimulatedJSON(writer, GBPLocation{"name": testLocation, "title": "Created Mistle HQ"})
			return
		}
		if request.Method == http.MethodGet {
			// Official method: GET https://mybusinessbusinessinformation.googleapis.com/v1/{parent=accounts/*}/locations?readMask=...
			writeSimulatedJSON(writer, GBPLocationsList{Locations: []GBPLocation{{
				"name":      testLocation,
				"title":     "Mistle HQ",
				"storeCode": "sg-1",
			}}})
			return
		}
		http.Error(writer, "unsupported method", http.StatusMethodNotAllowed)
	case "/v1/locations/456":
		if request.Method == http.MethodPatch {
			// Official method: PATCH https://mybusinessbusinessinformation.googleapis.com/v1/{location.name=locations/*}
			writeSimulatedJSON(writer, GBPLocation{"name": testLocation, "title": "Patched Mistle HQ"})
			return
		}
		if request.Method == http.MethodDelete {
			// Official method: DELETE https://mybusinessbusinessinformation.googleapis.com/v1/{name=locations/*}
			writeSimulatedJSON(writer, map[string]any{})
			return
		}
		if request.Method == http.MethodGet {
			// Official method: GET https://mybusinessbusinessinformation.googleapis.com/v1/{name=locations/*}?readMask=...
			writeSimulatedJSON(writer, GBPLocation{"name": testLocation, "title": "Mistle HQ"})
			return
		}
		http.Error(writer, "unsupported method", http.StatusMethodNotAllowed)
	case "/v4/accounts/123/locations/456/reviews":
		// Official method: GET https://mybusiness.googleapis.com/v4/{parent=accounts/*/locations/*}/reviews
		writeSimulatedJSON(writer, map[string]any{"reviews": []map[string]any{{"reviewId": "abc", "starRating": "FIVE"}}, "totalReviewCount": 1})
	case "/v4/accounts/123/locations/456/reviews/abc":
		// Official method: GET https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*/reviews/*}
		writeSimulatedJSON(writer, map[string]any{"reviewId": "abc", "starRating": "FIVE"})
	case "/v4/accounts/123/locations/456/reviews/abc/reply":
		if request.Method == http.MethodPut {
			// Official method: PUT https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*/reviews/*}/reply
			writeSimulatedJSON(writer, map[string]any{"comment": "Thanks!"})
			return
		}
		if request.Method == http.MethodDelete {
			// Official method: DELETE https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*/reviews/*}/reply
			writeSimulatedJSON(writer, map[string]any{})
			return
		}
		http.Error(writer, "unsupported method", http.StatusMethodNotAllowed)
	case "/v4/accounts/123/locations/456/media":
		if request.Method == http.MethodPost {
			// Official method: POST https://mybusiness.googleapis.com/v4/{parent=accounts/*/locations/*}/media
			writeSimulatedJSON(writer, map[string]any{"name": "accounts/123/locations/456/media/789"})
			return
		}
		if request.Method == http.MethodGet {
			// Official method: GET https://mybusiness.googleapis.com/v4/{parent=accounts/*/locations/*}/media
			writeSimulatedJSON(writer, map[string]any{"mediaItems": []map[string]any{{"name": "accounts/123/locations/456/media/789"}}})
			return
		}
		http.Error(writer, "unsupported method", http.StatusMethodNotAllowed)
	case "/v4/accounts/123/locations/456/media/789":
		if request.Method == http.MethodPatch {
			// Official method: PATCH https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*/media/*}
			writeSimulatedJSON(writer, map[string]any{"name": "accounts/123/locations/456/media/789", "description": "Patched"})
			return
		}
		if request.Method == http.MethodDelete {
			// Official method: DELETE https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*/media/*}
			writeSimulatedJSON(writer, map[string]any{})
			return
		}
		if request.Method == http.MethodGet {
			// Official method: GET https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*/media/*}
			writeSimulatedJSON(writer, map[string]any{"name": "accounts/123/locations/456/media/789"})
			return
		}
		http.Error(writer, "unsupported method", http.StatusMethodNotAllowed)
	case "/v4/accounts/123/locations/456/media:startUpload":
		// Official method: POST https://mybusiness.googleapis.com/v4/{parent=accounts/*/locations/*}/media:startUpload
		writeSimulatedJSON(writer, map[string]any{"resourceName": "upload-ref"})
	case "/v4/accounts/123/locations/456/localPosts":
		if request.Method == http.MethodPost {
			// Official method: POST https://mybusiness.googleapis.com/v4/{parent=accounts/*/locations/*}/localPosts
			writeSimulatedJSON(writer, map[string]any{"name": "accounts/123/locations/456/localPosts/post-1"})
			return
		}
		if request.Method == http.MethodGet {
			// Official method: GET https://mybusiness.googleapis.com/v4/{parent=accounts/*/locations/*}/localPosts
			writeSimulatedJSON(writer, map[string]any{"localPosts": []map[string]any{{"name": "accounts/123/locations/456/localPosts/post-1"}}})
			return
		}
		http.Error(writer, "unsupported method", http.StatusMethodNotAllowed)
	case "/v4/accounts/123/locations/456/localPosts/post-1":
		if request.Method == http.MethodPatch {
			// Official method: PATCH https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*/localPosts/*}
			writeSimulatedJSON(writer, map[string]any{"name": "accounts/123/locations/456/localPosts/post-1", "summary": "Patched"})
			return
		}
		if request.Method == http.MethodDelete {
			// Official method: DELETE https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*/localPosts/*}
			writeSimulatedJSON(writer, map[string]any{})
			return
		}
		if request.Method == http.MethodGet {
			// Official method: GET https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*/localPosts/*}
			writeSimulatedJSON(writer, map[string]any{"name": "accounts/123/locations/456/localPosts/post-1"})
			return
		}
		http.Error(writer, "unsupported method", http.StatusMethodNotAllowed)
	case "/v4/accounts/123/locations/456/localPosts:reportInsights":
		// Official method: POST https://mybusiness.googleapis.com/v4/{name=accounts/*/locations/*}/localPosts:reportInsights
		writeSimulatedJSON(writer, map[string]any{"localPostMetrics": []map[string]any{{"localPostName": "accounts/123/locations/456/localPosts/post-1"}}})
	case "/v1/locations/456:getDailyMetricsTimeSeries":
		// Official method: GET https://businessprofileperformance.googleapis.com/v1/{name=locations/*}:getDailyMetricsTimeSeries
		writeSimulatedJSON(writer, map[string]any{"timeSeries": map[string]any{"datedValues": []map[string]any{{"date": map[string]any{"year": 2026, "month": 6, "day": 1}, "value": "7"}}}})
	case "/v1/locations/456/searchkeywords/impressions/monthly":
		// Official method: GET https://businessprofileperformance.googleapis.com/v1/{parent=locations/*}/searchkeywords/impressions/monthly
		writeSimulatedJSON(writer, map[string]any{"searchKeywordsCounts": []map[string]any{{"searchKeyword": "mistle", "insightsValue": map[string]any{"value": "12"}}}})
	default:
		http.Error(writer, fmt.Sprintf("unexpected path: %s", request.URL.RequestURI()), http.StatusNotFound)
	}
}

func writeSimulatedJSON(writer http.ResponseWriter, value any) {
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (api *simulatedGBPAPI) sawPath(path string) bool {
	for _, candidate := range api.paths {
		if candidate == path {
			return true
		}
	}
	return false
}
