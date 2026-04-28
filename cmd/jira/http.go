package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("path must start with '/': %s", path)
	}

	url := jc.baseURL + path

	request, err := http.NewRequestWithContext(
		context.Background(),
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
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("path must start with '/': %s", path)
	}

	url := jc.baseURL + path

	request, err := http.NewRequestWithContext(
		context.Background(),
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
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("path must start with '/': %s", path)
	}

	url := jc.baseURL + path

	request, err := http.NewRequestWithContext(
		context.Background(),
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
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("path must start with '/': %s", path)
	}

	url := jc.baseURL + path

	request, err := http.NewRequestWithContext(
		context.Background(),
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
	body, err := jc.get("/rest/api/3/myself")

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
	body, err := jc.get("/rest/api/3/project/search")
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
	body, err := jc.get(fmt.Sprintf("/rest/api/3/issue/%s?fields=summary,status,assignee", issueOrKey))
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
	return jc.delete(fmt.Sprintf("/rest/api/3/issue/%s", issueOrKey))
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

	responseBody, err := jc.post("/rest/api/3/issue", requestBody)
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
	searchRequest := JiraIssueSearchRequest{
		JQL:    jql,
		Fields: []string{"summary", "status"},
	}

	requestBody, err := json.Marshal(searchRequest)
	if err != nil {
		return JiraIssueSearchResult{}, err
	}

	responseBody, err := jc.post("/rest/api/3/search/jql", requestBody)
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
	requestBody, err := json.Marshal(struct {
		AccountID *string `json:"accountId"`
	}{
		AccountID: input.AccountID,
	})
	if err != nil {
		return err
	}

	_, err = jc.put(fmt.Sprintf("/rest/api/3/issue/%s/assignee", issueOrKey), requestBody)
	return err
}

func (jc JiraClient) ListIssueTransitions(issueOrKey string) (JiraTransitionList, error) {
	responseBody, err := jc.get(fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueOrKey))
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

	_, err = jc.post(fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueOrKey), requestBody)
	return err
}

type UpdateIssueInput struct {
	Summary     *string
	Description *string
}

func (jc JiraClient) UpdateIssue(issueOrKey string, input UpdateIssueInput) error {
	fields := make(map[string]any)
	if input.Summary != nil {
		fields["summary"] = *input.Summary
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

	_, err = jc.put(fmt.Sprintf("/rest/api/3/issue/%s", issueOrKey), requestBody)
	return err
}

func (jc JiraClient) GetIssueEditMeta(issueOrKey string) (JiraIssueEditMeta, error) {
	responseBody, err := jc.get(fmt.Sprintf("/rest/api/3/issue/%s/editmeta", issueOrKey))
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

	responseBody, err := jc.post(fmt.Sprintf("/rest/api/3/issue/%s/comment", issueOrKey), requestBody)
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
	return jc.delete(fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueOrKey, commentID))
}
