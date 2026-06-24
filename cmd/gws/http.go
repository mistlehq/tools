package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GWSClient struct {
	config Config
	client *http.Client
}

type GWSAPI string

const (
	GWSAPIDrive  GWSAPI = "drive"
	GWSAPISheets GWSAPI = "sheets"
	GWSAPIDocs   GWSAPI = "docs"
	GWSAPISlides GWSAPI = "slides"
)

type GWSRequest struct {
	API    string         `json:"api"`
	Method string         `json:"method"`
	Path   string         `json:"path"`
	Params map[string]any `json:"params,omitempty"`
	Body   map[string]any `json:"body,omitempty"`
}

type GWSDriveFile struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	MimeType    string   `json:"mimeType,omitempty"`
	Parents     []string `json:"parents,omitempty"`
	WebViewLink string   `json:"webViewLink,omitempty"`
}

type GWSDriveFilesList struct {
	Files         []GWSDriveFile `json:"files"`
	NextPageToken string         `json:"nextPageToken,omitempty"`
}

type GWSDrivePermission struct {
	ID           string `json:"id,omitempty"`
	Type         string `json:"type,omitempty"`
	Role         string `json:"role,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	Domain       string `json:"domain,omitempty"`
}

type GWSDrivePermissionsList struct {
	Permissions []GWSDrivePermission `json:"permissions"`
}

type GWSRawResult = map[string]any

func NewGWSClient(config Config) GWSClient {
	return GWSClient{config: config, client: http.DefaultClient}
}

func (gc GWSClient) AuthTest(ctx context.Context) (GWSRawResult, error) {
	return gc.getJSON(ctx, GWSAPIDrive, "/about", map[string]any{"fields": "user,storageQuota"})
}

func (gc GWSClient) Request(request GWSRequest) ([]byte, error) {
	return gc.RequestContext(context.Background(), request)
}

func (gc GWSClient) RequestContext(ctx context.Context, request GWSRequest) ([]byte, error) {
	api := GWSAPI(strings.TrimSpace(request.API))
	baseURL, err := gc.baseURL(api)
	if err != nil {
		return nil, err
	}

	method := strings.ToUpper(strings.TrimSpace(request.Method))
	if method == "" {
		method = http.MethodGet
	}
	if method != http.MethodGet && method != http.MethodPost && method != http.MethodPatch && method != http.MethodPut && method != http.MethodDelete {
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
	if strings.TrimSpace(request.Path) == "" {
		return nil, fmt.Errorf("path is required")
	}
	if !strings.HasPrefix(request.Path, "/") {
		return nil, fmt.Errorf("path must start with '/'")
	}

	requestURL, err := url.Parse(baseURL + request.Path)
	if err != nil {
		return nil, err
	}
	query := requestURL.Query()
	for key, value := range request.Params {
		setQueryValue(query, key, value)
	}
	requestURL.RawQuery = query.Encode()

	var body io.Reader
	if request.Body != nil {
		encodedBody, err := json.Marshal(request.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(encodedBody)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, method, requestURL.String(), body)
	if err != nil {
		return nil, err
	}
	if request.Body != nil {
		httpRequest.Header.Set("Content-Type", "application/json")
	}

	response, err := gc.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("google workspace %s api %s %s: rate limited: %s", api, method, request.Path, string(responseBody))
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("google workspace %s api %s %s failed with status %d: %s", api, method, request.Path, response.StatusCode, string(responseBody))
	}
	return responseBody, nil
}

func (gc GWSClient) baseURL(api GWSAPI) (string, error) {
	switch api {
	case GWSAPIDrive:
		return gc.config.DriveBaseURL, nil
	case GWSAPISheets:
		return gc.config.SheetsBaseURL, nil
	case GWSAPIDocs:
		return gc.config.DocsBaseURL, nil
	case GWSAPISlides:
		return gc.config.SlidesBaseURL, nil
	default:
		return "", fmt.Errorf("unsupported api: %s", api)
	}
}

func (gc GWSClient) ListDriveFiles(ctx context.Context, params map[string]any) (GWSDriveFilesList, error) {
	var out GWSDriveFilesList
	err := gc.getTypedJSON(ctx, GWSAPIDrive, "/files", params, &out)
	return out, err
}

func (gc GWSClient) GetDriveFile(ctx context.Context, fileID string, params map[string]any) (GWSDriveFile, error) {
	if err := requireNonEmpty("file-id", fileID); err != nil {
		return GWSDriveFile{}, err
	}
	var out GWSDriveFile
	err := gc.getTypedJSON(ctx, GWSAPIDrive, "/files/"+url.PathEscape(fileID), params, &out)
	return out, err
}

func (gc GWSClient) CreateDriveFile(ctx context.Context, body map[string]any, params map[string]any) (GWSDriveFile, error) {
	var out GWSDriveFile
	err := gc.postTypedJSON(ctx, GWSAPIDrive, "/files", body, params, &out)
	return out, err
}

func (gc GWSClient) CopyDriveFile(ctx context.Context, fileID string, body map[string]any, params map[string]any) (GWSDriveFile, error) {
	if err := requireNonEmpty("file-id", fileID); err != nil {
		return GWSDriveFile{}, err
	}
	var out GWSDriveFile
	err := gc.postTypedJSON(ctx, GWSAPIDrive, "/files/"+url.PathEscape(fileID)+"/copy", body, params, &out)
	return out, err
}

func (gc GWSClient) UpdateDriveFile(ctx context.Context, fileID string, body map[string]any, params map[string]any) (GWSDriveFile, error) {
	if err := requireNonEmpty("file-id", fileID); err != nil {
		return GWSDriveFile{}, err
	}
	var out GWSDriveFile
	err := gc.patchTypedJSON(ctx, GWSAPIDrive, "/files/"+url.PathEscape(fileID), body, params, &out)
	return out, err
}

func (gc GWSClient) DeleteDriveFile(ctx context.Context, fileID string) error {
	if err := requireNonEmpty("file-id", fileID); err != nil {
		return err
	}
	_, err := gc.RequestContext(ctx, GWSRequest{API: string(GWSAPIDrive), Method: http.MethodDelete, Path: "/files/" + url.PathEscape(fileID)})
	return err
}

func (gc GWSClient) ListDrivePermissions(ctx context.Context, fileID string) (GWSDrivePermissionsList, error) {
	if err := requireNonEmpty("file-id", fileID); err != nil {
		return GWSDrivePermissionsList{}, err
	}
	var out GWSDrivePermissionsList
	err := gc.getTypedJSON(ctx, GWSAPIDrive, "/files/"+url.PathEscape(fileID)+"/permissions", nil, &out)
	return out, err
}

func (gc GWSClient) CreateDrivePermission(ctx context.Context, fileID string, body map[string]any) (GWSDrivePermission, error) {
	if err := requireNonEmpty("file-id", fileID); err != nil {
		return GWSDrivePermission{}, err
	}
	var out GWSDrivePermission
	err := gc.postTypedJSON(ctx, GWSAPIDrive, "/files/"+url.PathEscape(fileID)+"/permissions", body, nil, &out)
	return out, err
}

func (gc GWSClient) DeleteDrivePermission(ctx context.Context, fileID string, permissionID string) error {
	if err := requireNonEmpty("file-id", fileID); err != nil {
		return err
	}
	if err := requireNonEmpty("permission-id", permissionID); err != nil {
		return err
	}
	_, err := gc.RequestContext(ctx, GWSRequest{API: string(GWSAPIDrive), Method: http.MethodDelete, Path: "/files/" + url.PathEscape(fileID) + "/permissions/" + url.PathEscape(permissionID)})
	return err
}

func (gc GWSClient) GetSpreadsheet(ctx context.Context, spreadsheetID string) (GWSRawResult, error) {
	if err := requireNonEmpty("spreadsheet-id", spreadsheetID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPISheets, "/spreadsheets/"+url.PathEscape(spreadsheetID), nil)
}

func (gc GWSClient) CreateSpreadsheet(ctx context.Context, body map[string]any) (GWSRawResult, error) {
	return gc.postJSON(ctx, GWSAPISheets, "/spreadsheets", body, nil)
}

func (gc GWSClient) BatchUpdateSpreadsheet(ctx context.Context, spreadsheetID string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("spreadsheet-id", spreadsheetID); err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GWSAPISheets, "/spreadsheets/"+url.PathEscape(spreadsheetID)+":batchUpdate", body, nil)
}

func (gc GWSClient) GetSpreadsheetValues(ctx context.Context, spreadsheetID string, valueRange string) (GWSRawResult, error) {
	if err := requireNonEmpty("spreadsheet-id", spreadsheetID); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("range", valueRange); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPISheets, "/spreadsheets/"+url.PathEscape(spreadsheetID)+"/values/"+url.PathEscape(valueRange), nil)
}

func (gc GWSClient) UpdateSpreadsheetValues(ctx context.Context, spreadsheetID string, valueRange string, valueInputOption string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("spreadsheet-id", spreadsheetID); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("range", valueRange); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("value-input-option", valueInputOption); err != nil {
		return nil, err
	}
	return gc.putJSON(ctx, GWSAPISheets, "/spreadsheets/"+url.PathEscape(spreadsheetID)+"/values/"+url.PathEscape(valueRange), body, map[string]any{"valueInputOption": valueInputOption})
}

func (gc GWSClient) BatchUpdateSpreadsheetValues(ctx context.Context, spreadsheetID string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("spreadsheet-id", spreadsheetID); err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GWSAPISheets, "/spreadsheets/"+url.PathEscape(spreadsheetID)+"/values:batchUpdate", body, nil)
}

func (gc GWSClient) GetDocument(ctx context.Context, documentID string) (GWSRawResult, error) {
	if err := requireNonEmpty("document-id", documentID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIDocs, "/documents/"+url.PathEscape(documentID), nil)
}

func (gc GWSClient) BatchUpdateDocument(ctx context.Context, documentID string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("document-id", documentID); err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GWSAPIDocs, "/documents/"+url.PathEscape(documentID)+":batchUpdate", body, nil)
}

func (gc GWSClient) GetPresentation(ctx context.Context, presentationID string) (GWSRawResult, error) {
	if err := requireNonEmpty("presentation-id", presentationID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPISlides, "/presentations/"+url.PathEscape(presentationID), nil)
}

func (gc GWSClient) CreatePresentation(ctx context.Context, body map[string]any) (GWSRawResult, error) {
	return gc.postJSON(ctx, GWSAPISlides, "/presentations", body, nil)
}

func (gc GWSClient) BatchUpdatePresentation(ctx context.Context, presentationID string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("presentation-id", presentationID); err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GWSAPISlides, "/presentations/"+url.PathEscape(presentationID)+":batchUpdate", body, nil)
}

func (gc GWSClient) getJSON(ctx context.Context, api GWSAPI, path string, params map[string]any) (GWSRawResult, error) {
	var out GWSRawResult
	err := gc.getTypedJSON(ctx, api, path, params, &out)
	return out, err
}

func (gc GWSClient) postJSON(ctx context.Context, api GWSAPI, path string, body map[string]any, params map[string]any) (GWSRawResult, error) {
	var out GWSRawResult
	err := gc.postTypedJSON(ctx, api, path, body, params, &out)
	return out, err
}

func (gc GWSClient) putJSON(ctx context.Context, api GWSAPI, path string, body map[string]any, params map[string]any) (GWSRawResult, error) {
	var out GWSRawResult
	err := gc.requestTypedJSON(ctx, api, http.MethodPut, path, body, params, &out)
	return out, err
}

func (gc GWSClient) patchTypedJSON(ctx context.Context, api GWSAPI, path string, body map[string]any, params map[string]any, out any) error {
	return gc.requestTypedJSON(ctx, api, http.MethodPatch, path, body, params, out)
}

func (gc GWSClient) postTypedJSON(ctx context.Context, api GWSAPI, path string, body map[string]any, params map[string]any, out any) error {
	return gc.requestTypedJSON(ctx, api, http.MethodPost, path, body, params, out)
}

func (gc GWSClient) getTypedJSON(ctx context.Context, api GWSAPI, path string, params map[string]any, out any) error {
	return gc.requestTypedJSON(ctx, api, http.MethodGet, path, nil, params, out)
}

func (gc GWSClient) requestTypedJSON(ctx context.Context, api GWSAPI, method string, path string, body map[string]any, params map[string]any, out any) error {
	if method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut {
		if body == nil {
			return fmt.Errorf("request body is required")
		}
	}
	responseBody, err := gc.RequestContext(ctx, GWSRequest{API: string(api), Method: method, Path: path, Params: params, Body: body})
	if err != nil {
		return err
	}
	if len(responseBody) == 0 {
		return nil
	}
	return json.Unmarshal(responseBody, out)
}

func setQueryValue(query url.Values, key string, value any) {
	switch typedValue := value.(type) {
	case nil:
		return
	case string:
		query.Set(key, typedValue)
	default:
		encoded, err := json.Marshal(typedValue)
		if err != nil {
			query.Set(key, fmt.Sprint(typedValue))
			return
		}
		query.Set(key, string(encoded))
	}
}

func requireNonEmpty(label string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", label)
	}
	return nil
}
