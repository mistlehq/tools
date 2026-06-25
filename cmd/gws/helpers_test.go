package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mistlehq/tools/internal/testproxy"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/oauth2/google"
)

type commandResult struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

type gwsTestResources struct {
	DriveFileID    string
	DriveFolderID  string
	DocumentID     string
	SpreadsheetID  string
	PresentationID string
}

func runCommandWithInput(t *testing.T, env Environment, input string, args ...string) (commandResult, error) {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cli := CLI{stdin: bytes.NewBufferString(input), stdout: &stdout, stderr: &stderr, env: env}
	err := cli.run(args)
	return commandResult{stdout: stdout, stderr: stderr}, err
}

func init() {
	loadDotEnvTest()
}

func loadDotEnvTest() {
	path := findDotEnvTest()
	if path == "" {
		return
	}
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}

func findDotEnvTest() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		path := filepath.Join(dir, ".env.test")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func getRequiredEnv(t *testing.T, name string) string {
	t.Helper()
	value := os.Getenv(name)
	if value == "" {
		t.Skipf("skipping: %s is not set", name)
	}
	return value
}

func getOptionalEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func setupCommandEnvironment(t *testing.T) Environment {
	t.Helper()
	token := mintGoogleWorkspaceTestAccessToken(t)
	driveProxy := startBearerProxy(t, "drive", getOptionalEnv("GWS_TEST_DRIVE_BASE_URL", "https://www.googleapis.com/drive/v3"), token)
	sheetsProxy := startBearerProxy(t, "sheets", getOptionalEnv("GWS_TEST_SHEETS_BASE_URL", "https://sheets.googleapis.com/v4"), token)
	docsProxy := startBearerProxy(t, "docs", getOptionalEnv("GWS_TEST_DOCS_BASE_URL", "https://docs.googleapis.com/v1"), token)
	slidesProxy := startBearerProxy(t, "slides", getOptionalEnv("GWS_TEST_SLIDES_BASE_URL", "https://slides.googleapis.com/v1"), token)
	gmailProxy := startBearerProxy(t, "gmail", getOptionalEnv("GWS_TEST_GMAIL_BASE_URL", "https://gmail.googleapis.com/gmail/v1"), token)
	calendarProxy := startBearerProxy(t, "calendar", getOptionalEnv("GWS_TEST_CALENDAR_BASE_URL", "https://www.googleapis.com/calendar/v3"), token)
	chatProxy := startBearerProxy(t, "chat", getOptionalEnv("GWS_TEST_CHAT_BASE_URL", "https://chat.googleapis.com/v1"), token)
	peopleProxy := startBearerProxy(t, "people", getOptionalEnv("GWS_TEST_PEOPLE_BASE_URL", "https://people.googleapis.com/v1"), token)
	return Environment{
		"GWS_DRIVE_BASE_URL":    driveProxy.BaseURL,
		"GWS_SHEETS_BASE_URL":   sheetsProxy.BaseURL,
		"GWS_DOCS_BASE_URL":     docsProxy.BaseURL,
		"GWS_SLIDES_BASE_URL":   slidesProxy.BaseURL,
		"GWS_GMAIL_BASE_URL":    gmailProxy.BaseURL,
		"GWS_CALENDAR_BASE_URL": calendarProxy.BaseURL,
		"GWS_CHAT_BASE_URL":     chatProxy.BaseURL,
		"GWS_PEOPLE_BASE_URL":   peopleProxy.BaseURL,
	}
}

func startBearerProxy(t *testing.T, name string, upstreamBaseURL string, token string) *testproxy.Server {
	t.Helper()
	proxy, err := testproxy.Start(testproxy.Config{
		UpstreamBaseURL: upstreamBaseURL,
		AuthMode:        testproxy.AuthModeBearer,
		Token:           token,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := proxy.Close(); err != nil {
			t.Errorf("failed to close %s proxy: %v", name, err)
		}
	})
	return proxy
}

func setupGWSClient(t *testing.T) (Environment, GWSClient) {
	t.Helper()
	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	return env, NewGWSClient(config)
}

func mintGoogleWorkspaceTestAccessToken(t *testing.T) string {
	t.Helper()
	keyJSONBase64 := getRequiredEnv(t, "GWS_TEST_SERVICE_ACCOUNT_KEY_JSON_BASE64")
	keyJSON, err := base64.StdEncoding.DecodeString(keyJSONBase64)
	if err != nil {
		t.Fatalf("failed to decode GWS_TEST_SERVICE_ACCOUNT_KEY_JSON_BASE64: %v", err)
	}
	config, err := google.JWTConfigFromJSON(
		keyJSON,
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/spreadsheets",
		"https://www.googleapis.com/auth/documents",
		"https://www.googleapis.com/auth/presentations",
		"https://www.googleapis.com/auth/gmail.readonly",
		"https://www.googleapis.com/auth/gmail.compose",
		"https://www.googleapis.com/auth/calendar.calendarlist.readonly",
		"https://www.googleapis.com/auth/calendar.events.readonly",
		"https://www.googleapis.com/auth/calendar.events.freebusy",
		"https://www.googleapis.com/auth/chat.spaces.readonly",
		"https://www.googleapis.com/auth/chat.memberships.readonly",
		"https://www.googleapis.com/auth/chat.messages.readonly",
		"https://www.googleapis.com/auth/chat.messages.create",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/contacts.readonly",
		"https://www.googleapis.com/auth/directory.readonly",
	)
	if err != nil {
		t.Fatalf("failed to parse service account JSON key: %v", err)
	}
	if os.Getenv("GWS_TEST_USE_WORKSPACE_USER_EMAIL") == "1" {
		subject := os.Getenv("GWS_TEST_WORKSPACE_USER_EMAIL")
		if subject == "" {
			t.Fatal("GWS_TEST_WORKSPACE_USER_EMAIL is required when GWS_TEST_USE_WORKSPACE_USER_EMAIL=1")
		}
		config.Subject = subject
	}
	token, err := config.TokenSource(context.Background()).Token()
	if err != nil {
		t.Fatalf("failed to mint Google Workspace test access token: %v", err)
	}
	if token.AccessToken == "" {
		t.Fatal("Google Workspace test access token was empty")
	}
	return token.AccessToken
}

func testGWSWorkspaceResources(t *testing.T) gwsTestResources {
	t.Helper()
	return gwsTestResources{
		DriveFileID:    getRequiredEnv(t, "GWS_TEST_DRIVE_FILE_ID"),
		DriveFolderID:  getRequiredEnv(t, "GWS_TEST_DRIVE_FOLDER_ID"),
		DocumentID:     getRequiredEnv(t, "GWS_TEST_DOCUMENT_ID"),
		SpreadsheetID:  getRequiredEnv(t, "GWS_TEST_SPREADSHEET_ID"),
		PresentationID: getRequiredEnv(t, "GWS_TEST_PRESENTATION_ID"),
	}
}

func requireDriveCreationCapable(t *testing.T) {
	t.Helper()
	if os.Getenv("GWS_TEST_RUN_CREATION_TESTS") != "1" {
		t.Skip("skipping Drive file creation coverage: set GWS_TEST_RUN_CREATION_TESTS=1 only after Google Drive ownership/storage configuration supports creating Workspace files")
	}
}

func writeTempJSONRequest(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func decodeCommandJSON(t *testing.T, result commandResult, out any) {
	t.Helper()
	if err := json.Unmarshal(result.stdout.Bytes(), out); err != nil {
		t.Fatalf("expected valid JSON output: %v\nstdout: %s", err, result.stdout.String())
	}
}

func newLocalGWSMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	return newGWSMCPTestSession(t, NewGWSClient(Config{
		DriveBaseURL:    "http://127.0.0.1",
		SheetsBaseURL:   "http://127.0.0.1",
		DocsBaseURL:     "http://127.0.0.1",
		SlidesBaseURL:   "http://127.0.0.1",
		GmailBaseURL:    "http://127.0.0.1",
		CalendarBaseURL: "http://127.0.0.1",
		ChatBaseURL:     "http://127.0.0.1",
		PeopleBaseURL:   "http://127.0.0.1",
	}), nil)
}

func newGWSMCPTestSession(t *testing.T, gc GWSClient, tools map[string]bool) *mcp.ClientSession {
	t.Helper()
	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newGWSMCPHTTPHandler(gc, tools))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return connectGWSMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func connectGWSMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()
	client := mcp.NewClient(&mcp.Implementation{Name: "gws-test-client", Version: "dev"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func callGWSMCPTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
	t.Helper()
	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("expected %s to succeed, got tool error: %#v", name, result.Content)
	}
	return result
}

func decodeMCPStructuredContent(t *testing.T, result *mcp.CallToolResult, output any) {
	t.Helper()
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, output); err != nil {
		t.Fatalf("expected structured content to decode: %v\ncontent: %s", err, string(raw))
	}
}
