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

type JiraClient struct {
	baseURL string
	client  *http.Client
}

func NewJiraClient(config Config) JiraClient {
	return JiraClient{
		baseURL: config.BaseURL,
		client:  http.DefaultClient,
	}
}

func (jc JiraClient) get(path string) ([]byte, error) {
	return jc.getContext(context.Background(), path)
}

func (jc JiraClient) getContext(ctx context.Context, path string) ([]byte, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("path must start with '/': %s", path)
	}

	url := jc.baseURL + path

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	response, err := jc.client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("request failed with status %d: %s", response.StatusCode, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (jc JiraClient) post(path string, body []byte) ([]byte, error) {
	return jc.postContext(context.Background(), path, body)
}

func (jc JiraClient) postContext(ctx context.Context, path string, body []byte) ([]byte, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("path must start with '/': %s", path)
	}

	url := jc.baseURL + path

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := jc.client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("request failed with status %d: %s", response.StatusCode, string(body))
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (jc JiraClient) put(path string, body []byte) ([]byte, error) {
	return jc.putContext(context.Background(), path, body)
}

func (jc JiraClient) putContext(ctx context.Context, path string, body []byte) ([]byte, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("path must start with '/': %s", path)
	}

	url := jc.baseURL + path

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := jc.client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("request failed with status %d: %s", response.StatusCode, string(body))
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (jc JiraClient) delete(path string) error {
	return jc.deleteContext(context.Background(), path)
}

func (jc JiraClient) deleteContext(ctx context.Context, path string) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("path must start with '/': %s", path)
	}

	url := jc.baseURL + path

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	response, err := jc.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("request failed with status %d: %s", response.StatusCode, string(body))
	}

	return nil
}

type JiraMyself struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
	Email       string `json:"emailAddress"`
}

type JiraUser struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
}

func (jc JiraClient) GetMyself() (JiraMyself, error) {
	return jc.GetMyselfContext(context.Background())
}

func (jc JiraClient) GetMyselfContext(ctx context.Context) (JiraMyself, error) {
	body, err := jc.getContext(ctx, "/rest/api/3/myself")

	if err != nil {
		return JiraMyself{}, err
	}

	var myself JiraMyself
	if err := json.Unmarshal(body, &myself); err != nil {
		return JiraMyself{}, err
	}

	return myself, nil
}

type JiraProject struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type JiraProjectList struct {
	Values []JiraProject `json:"values"`
}

func (jc JiraClient) ListProjects() (JiraProjectList, error) {
	return jc.ListProjectsContext(context.Background())
}

func (jc JiraClient) ListProjectsContext(ctx context.Context) (JiraProjectList, error) {
	body, err := jc.getContext(ctx, "/rest/api/3/project/search")
	if err != nil {
		return JiraProjectList{}, err
	}

	var projectList JiraProjectList
	if err := json.Unmarshal(body, &projectList); err != nil {
		return JiraProjectList{}, err
	}

	return projectList, nil
}

type JiraIssue struct {
	ID     string          `json:"id"`
	Key    string          `json:"key"`
	Fields JiraIssueFields `json:"fields"`
}

type JiraIssueFields struct {
	Summary  string          `json:"summary"`
	Status   JiraIssueStatus `json:"status"`
	Assignee *JiraUser       `json:"assignee"`
}

type JiraIssueStatus struct {
	Name string `json:"name"`
}

type JiraTransition struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	To   JiraIssueStatus `json:"to"`
}

type JiraComment struct {
	ID      string   `json:"id"`
	Author  JiraUser `json:"author"`
	Created string   `json:"created"`
}

type AddCommentInput struct {
	Body string
}

func (jc JiraClient) GetIssue(issueOrKey string) (JiraIssue, error) {
	return jc.GetIssueContext(context.Background(), issueOrKey)
}

func (jc JiraClient) GetIssueContext(ctx context.Context, issueOrKey string) (JiraIssue, error) {
	body, err := jc.getContext(ctx, fmt.Sprintf("/rest/api/3/issue/%s?fields=summary,status,assignee", issueOrKey))
	if err != nil {
		return JiraIssue{}, err
	}

	var jiraIssue JiraIssue
	if err := json.Unmarshal(body, &jiraIssue); err != nil {
		return JiraIssue{}, err
	}

	return jiraIssue, nil
}

func (jc JiraClient) DeleteIssue(issueOrKey string) error {
	return jc.DeleteIssueContext(context.Background(), issueOrKey)
}

func (jc JiraClient) DeleteIssueContext(ctx context.Context, issueOrKey string) error {
	return jc.deleteContext(ctx, fmt.Sprintf("/rest/api/3/issue/%s", issueOrKey))
}

type JiraIssueSearchRequest struct {
	JQL    string   `json:"jql"`
	Fields []string `json:"fields"`
}

type JiraIssueSearchResult struct {
	Issues []JiraIssue `json:"issues"`
}

type JiraTransitionList struct {
	Transitions []JiraTransition `json:"transitions"`
}

type JiraFieldSchema struct {
	Type string `json:"type"`
}

type JiraEditMetaField struct {
	Name     string          `json:"name"`
	Required bool            `json:"required"`
	Schema   JiraFieldSchema `json:"schema"`
}

type JiraIssueEditMeta struct {
	Fields map[string]JiraEditMetaField `json:"fields"`
}

type CreateIssueInput struct {
	ProjectID     string
	ProjectKey    string
	IssueTypeID   string
	IssueTypeName string
	Summary       string
	Description   *string
}

type JiraCreatedIssue struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

func (jc JiraClient) CreateIssue(input CreateIssueInput) (JiraCreatedIssue, error) {
	return jc.CreateIssueContext(context.Background(), input)
}

func (jc JiraClient) CreateIssueContext(ctx context.Context, input CreateIssueInput) (JiraCreatedIssue, error) {
	fields := map[string]any{
		"summary": input.Summary,
	}

	if input.ProjectID != "" {
		fields["project"] = map[string]string{
			"id": input.ProjectID,
		}
	} else {
		fields["project"] = map[string]string{
			"key": input.ProjectKey,
		}
	}

	if input.IssueTypeID != "" {
		fields["issuetype"] = map[string]string{
			"id": input.IssueTypeID,
		}
	} else {
		fields["issuetype"] = map[string]string{
			"name": input.IssueTypeName,
		}
	}

	if input.Description != nil {
		descriptionDocument, err := NewJiraTextDocument(*input.Description)
		if err != nil {
			return JiraCreatedIssue{}, err
		}

		fields["description"] = descriptionDocument
	}

	requestBody, err := json.Marshal(struct {
		Fields map[string]any `json:"fields"`
	}{
		Fields: fields,
	})
	if err != nil {
		return JiraCreatedIssue{}, err
	}

	responseBody, err := jc.postContext(ctx, "/rest/api/3/issue", requestBody)
	if err != nil {
		return JiraCreatedIssue{}, err
	}

	var issue JiraCreatedIssue
	if err := json.Unmarshal(responseBody, &issue); err != nil {
		return JiraCreatedIssue{}, err
	}

	return issue, nil
}

func (jc JiraClient) SearchIssues(jql string) (JiraIssueSearchResult, error) {
	return jc.SearchIssuesContext(context.Background(), jql)
}

func (jc JiraClient) SearchIssuesContext(ctx context.Context, jql string) (JiraIssueSearchResult, error) {
	searchRequest := JiraIssueSearchRequest{
		JQL:    jql,
		Fields: []string{"summary", "status"},
	}

	requestBody, err := json.Marshal(searchRequest)
	if err != nil {
		return JiraIssueSearchResult{}, err
	}

	responseBody, err := jc.postContext(ctx, "/rest/api/3/search/jql", requestBody)
	if err != nil {
		return JiraIssueSearchResult{}, err
	}

	var searchResult JiraIssueSearchResult
	if err := json.Unmarshal(responseBody, &searchResult); err != nil {
		return JiraIssueSearchResult{}, err
	}

	return searchResult, nil
}

type AssignIssueInput struct {
	AccountID *string
}

func (jc JiraClient) AssignIssue(issueOrKey string, input AssignIssueInput) error {
	return jc.AssignIssueContext(context.Background(), issueOrKey, input)
}

func (jc JiraClient) AssignIssueContext(ctx context.Context, issueOrKey string, input AssignIssueInput) error {
	requestBody, err := json.Marshal(struct {
		AccountID *string `json:"accountId"`
	}{
		AccountID: input.AccountID,
	})
	if err != nil {
		return err
	}

	_, err = jc.putContext(ctx, fmt.Sprintf("/rest/api/3/issue/%s/assignee", issueOrKey), requestBody)
	return err
}

func (jc JiraClient) ListIssueTransitions(issueOrKey string) (JiraTransitionList, error) {
	return jc.ListIssueTransitionsContext(context.Background(), issueOrKey)
}

func (jc JiraClient) ListIssueTransitionsContext(ctx context.Context, issueOrKey string) (JiraTransitionList, error) {
	responseBody, err := jc.getContext(ctx, fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueOrKey))
	if err != nil {
		return JiraTransitionList{}, err
	}

	var transitionList JiraTransitionList
	if err := json.Unmarshal(responseBody, &transitionList); err != nil {
		return JiraTransitionList{}, err
	}

	return transitionList, nil
}

type TransitionIssueInput struct {
	TransitionID string
}

func (jc JiraClient) TransitionIssue(issueOrKey string, input TransitionIssueInput) error {
	return jc.TransitionIssueContext(context.Background(), issueOrKey, input)
}

func (jc JiraClient) TransitionIssueContext(ctx context.Context, issueOrKey string, input TransitionIssueInput) error {
	var payload struct {
		Transition struct {
			ID string `json:"id"`
		} `json:"transition"`
	}
	payload.Transition.ID = input.TransitionID

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = jc.postContext(ctx, fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueOrKey), requestBody)
	return err
}

type UpdateIssueInput struct {
	Description *string
	Fields      map[string]any
}

func (jc JiraClient) UpdateIssue(issueOrKey string, input UpdateIssueInput) error {
	return jc.UpdateIssueContext(context.Background(), issueOrKey, input)
}

func (jc JiraClient) UpdateIssueContext(ctx context.Context, issueOrKey string, input UpdateIssueInput) error {
	fields := make(map[string]any)
	for fieldID, value := range input.Fields {
		fields[fieldID] = value
	}

	if input.Description != nil {
		descriptionDocument, err := NewJiraTextDocument(*input.Description)
		if err != nil {
			return err
		}

		fields["description"] = descriptionDocument
	}

	requestBody, err := json.Marshal(struct {
		Fields map[string]any `json:"fields"`
	}{
		Fields: fields,
	})
	if err != nil {
		return err
	}

	_, err = jc.putContext(ctx, fmt.Sprintf("/rest/api/3/issue/%s", issueOrKey), requestBody)
	return err
}

func (jc JiraClient) GetIssueEditMeta(issueOrKey string) (JiraIssueEditMeta, error) {
	return jc.GetIssueEditMetaContext(context.Background(), issueOrKey)
}

func (jc JiraClient) GetIssueEditMetaContext(ctx context.Context, issueOrKey string) (JiraIssueEditMeta, error) {
	responseBody, err := jc.getContext(ctx, fmt.Sprintf("/rest/api/3/issue/%s/editmeta", issueOrKey))
	if err != nil {
		return JiraIssueEditMeta{}, err
	}

	var editMeta JiraIssueEditMeta
	if err := json.Unmarshal(responseBody, &editMeta); err != nil {
		return JiraIssueEditMeta{}, err
	}

	return editMeta, nil
}

func (jc JiraClient) AddIssueComment(issueOrKey string, input AddCommentInput) (JiraComment, error) {
	return jc.AddIssueCommentContext(context.Background(), issueOrKey, input)
}

func (jc JiraClient) AddIssueCommentContext(ctx context.Context, issueOrKey string, input AddCommentInput) (JiraComment, error) {
	bodyDocument, err := NewJiraTextDocument(input.Body)
	if err != nil {
		return JiraComment{}, err
	}

	requestBody, err := json.Marshal(struct {
		Body JiraDocument `json:"body"`
	}{
		Body: bodyDocument,
	})
	if err != nil {
		return JiraComment{}, err
	}

	responseBody, err := jc.postContext(ctx, fmt.Sprintf("/rest/api/3/issue/%s/comment", issueOrKey), requestBody)
	if err != nil {
		return JiraComment{}, err
	}

	var comment JiraComment
	if err := json.Unmarshal(responseBody, &comment); err != nil {
		return JiraComment{}, err
	}

	return comment, nil
}

func (jc JiraClient) DeleteIssueComment(issueOrKey string, commentID string) error {
	return jc.DeleteIssueCommentContext(context.Background(), issueOrKey, commentID)
}

func (jc JiraClient) DeleteIssueCommentContext(ctx context.Context, issueOrKey string, commentID string) error {
	return jc.deleteContext(ctx, fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueOrKey, commentID))
}

type JiraStatusScope struct {
	Type    string             `json:"type"`
	Project *JiraStatusProject `json:"project,omitempty"`
}

type JiraStatusProject struct {
	ID string `json:"id"`
}

type JiraStatus struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description,omitempty"`
	StatusCategory string          `json:"statusCategory"`
	Scope          JiraStatusScope `json:"scope"`
}

type JiraStatusPage struct {
	IsLast     bool         `json:"isLast"`
	MaxResults int          `json:"maxResults"`
	NextPage   string       `json:"nextPage,omitempty"`
	Self       string       `json:"self,omitempty"`
	StartAt    int          `json:"startAt"`
	Total      int          `json:"total"`
	Values     []JiraStatus `json:"values"`
}

type JiraStatusCreate struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	StatusCategory string `json:"statusCategory"`
}

type JiraStatusUpdate struct {
	ID             string  `json:"id"`
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	StatusCategory *string `json:"statusCategory,omitempty"`
}

type JiraStatusSearchInput struct {
	ProjectID      string
	StartAt        *int
	MaxResults     *int
	SearchString   string
	StatusCategory string
}

type JiraStatusCreateInput struct {
	Scope    JiraStatusScope    `json:"scope"`
	Statuses []JiraStatusCreate `json:"statuses"`
}

type JiraStatusUpdateInput struct {
	Statuses []JiraStatusUpdate `json:"statuses"`
}

type JiraBoardConfiguration map[string]any

func (jc JiraClient) GetStatuses(ids []string) ([]JiraStatus, error) {
	return jc.GetStatusesContext(context.Background(), ids)
}

func (jc JiraClient) GetStatusesContext(ctx context.Context, ids []string) ([]JiraStatus, error) {
	query := url.Values{}
	for _, id := range ids {
		query.Add("id", id)
	}

	body, err := jc.getContext(ctx, "/rest/api/3/statuses?"+query.Encode())
	if err != nil {
		return nil, err
	}

	var statuses []JiraStatus
	if err := json.Unmarshal(body, &statuses); err != nil {
		return nil, err
	}

	return statuses, nil
}

func (jc JiraClient) SearchStatuses(input JiraStatusSearchInput) (JiraStatusPage, error) {
	return jc.SearchStatusesContext(context.Background(), input)
}

func (jc JiraClient) SearchStatusesContext(ctx context.Context, input JiraStatusSearchInput) (JiraStatusPage, error) {
	query := url.Values{}
	if input.ProjectID != "" {
		query.Set("projectId", input.ProjectID)
	}
	if input.StartAt != nil {
		query.Set("startAt", fmt.Sprintf("%d", *input.StartAt))
	}
	if input.MaxResults != nil {
		query.Set("maxResults", fmt.Sprintf("%d", *input.MaxResults))
	}
	if input.SearchString != "" {
		query.Set("searchString", input.SearchString)
	}
	if input.StatusCategory != "" {
		query.Set("statusCategory", input.StatusCategory)
	}

	path := "/rest/api/3/statuses/search"
	if encodedQuery := query.Encode(); encodedQuery != "" {
		path += "?" + encodedQuery
	}

	body, err := jc.getContext(ctx, path)
	if err != nil {
		return JiraStatusPage{}, err
	}

	var statusPage JiraStatusPage
	if err := json.Unmarshal(body, &statusPage); err != nil {
		return JiraStatusPage{}, err
	}

	return statusPage, nil
}

func (jc JiraClient) CreateStatuses(input JiraStatusCreateInput) ([]JiraStatus, error) {
	return jc.CreateStatusesContext(context.Background(), input)
}

func (jc JiraClient) CreateStatusesContext(ctx context.Context, input JiraStatusCreateInput) ([]JiraStatus, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	responseBody, err := jc.postContext(ctx, "/rest/api/3/statuses", requestBody)
	if err != nil {
		return nil, err
	}

	var statuses []JiraStatus
	if err := json.Unmarshal(responseBody, &statuses); err != nil {
		return nil, err
	}

	return statuses, nil
}

func (jc JiraClient) UpdateStatuses(input JiraStatusUpdateInput) error {
	return jc.UpdateStatusesContext(context.Background(), input)
}

func (jc JiraClient) UpdateStatusesContext(ctx context.Context, input JiraStatusUpdateInput) error {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = jc.putContext(ctx, "/rest/api/3/statuses", requestBody)
	return err
}

func (jc JiraClient) DeleteStatuses(ids []string) error {
	return jc.DeleteStatusesContext(context.Background(), ids)
}

func (jc JiraClient) DeleteStatusesContext(ctx context.Context, ids []string) error {
	query := url.Values{}
	for _, id := range ids {
		query.Add("id", id)
	}

	return jc.deleteContext(ctx, "/rest/api/3/statuses?"+query.Encode())
}

func (jc JiraClient) GetBoardConfiguration(boardID string) (JiraBoardConfiguration, error) {
	return jc.GetBoardConfigurationContext(context.Background(), boardID)
}

func (jc JiraClient) GetBoardConfigurationContext(ctx context.Context, boardID string) (JiraBoardConfiguration, error) {
	body, err := jc.getContext(ctx, fmt.Sprintf("/rest/agile/1.0/board/%s/configuration", boardID))
	if err != nil {
		return nil, err
	}

	var configuration JiraBoardConfiguration
	if err := json.Unmarshal(body, &configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}
