package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{
		"jira mcp",
		"jira mcp serve",
		"Streamable HTTP",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp help to mention %q", want)
		}
	}
}

func TestMCPServeHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "jira", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	expected := []string{
		"--addr <addr>",
		"--endpoint <path>",
		"jira_issue_get",
		"jira_issue_search",
		"jira_status_create",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp serve help to mention %q", want)
		}
	}
}

func TestParseJiraMCPServeArgs(t *testing.T) {
	config, err := parseJiraMCPServeArgs(nil)
	if err != nil {
		t.Fatal(err)
	}
	if config.Addr != defaultMCPAddr || config.Endpoint != defaultMCPEndpoint {
		t.Fatalf("unexpected default config: %#v", config)
	}

	config, err = parseJiraMCPServeArgs([]string{
		"--addr", "127.0.0.1:9999",
		"--endpoint", "/jira-mcp",
	})
	if err != nil {
		t.Fatal(err)
	}
	if config.Addr != "127.0.0.1:9999" || config.Endpoint != "/jira-mcp" {
		t.Fatalf("unexpected config: %#v", config)
	}
}

func TestMCPServerListsJiraTools(t *testing.T) {
	jc := NewJiraClient(Config{BaseURL: "http://127.0.0.1"})
	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newJiraMCPHTTPHandler(jc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	session := connectJiraMCPTestClient(t, server.URL+defaultMCPEndpoint)
	defer session.Close()

	toolsResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	toolsByName := make(map[string]*mcp.Tool)
	for _, tool := range toolsResult.Tools {
		toolsByName[tool.Name] = tool
	}

	expected := map[string]string{
		"jira_auth_whoami":             jiraAuthWhoAmIDoc.Description,
		"jira_project_list":            jiraProjectListDoc.Description,
		"jira_issue_get":               jiraIssueGetDoc.Description,
		"jira_issue_search":            jiraIssueSearchDoc.Description,
		"jira_issue_create":            jiraIssueCreateDoc.Description,
		"jira_issue_delete":            jiraIssueDeleteDoc.Description,
		"jira_issue_comment_add":       jiraIssueCommentAddDoc.Description,
		"jira_issue_comment_delete":    jiraIssueCommentDeleteDoc.Description,
		"jira_issue_assign":            jiraIssueAssignDoc.Description,
		"jira_issue_transition_list":   jiraIssueTransitionListDoc.Description,
		"jira_issue_transition":        jiraIssueTransitionDoc.Description,
		"jira_issue_update":            jiraIssueUpdateDoc.Description,
		"jira_issue_editmeta":          jiraIssueEditMetaDoc.Description,
		"jira_status_get":              jiraStatusGetDoc.Description,
		"jira_status_search":           jiraStatusSearchDoc.Description,
		"jira_status_create":           jiraStatusCreateDoc.Description,
		"jira_status_update":           jiraStatusUpdateDoc.Description,
		"jira_status_delete":           jiraStatusDeleteDoc.Description,
		"jira_board_configuration_get": jiraBoardConfigurationGetDoc.Description,
	}
	if len(toolsByName) != len(expected) {
		t.Fatalf("expected exactly %d Jira tools, got %d: %#v", len(expected), len(toolsByName), toolsByName)
	}

	for name, description := range expected {
		tool, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected MCP tool %q to be listed", name)
		}

		if tool.Description != description {
			t.Fatalf("expected MCP tool %q description %q, got %q", name, description, tool.Description)
		}

		if strings.Contains(name, "delete") {
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || !*tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be annotated as destructive", name)
			}
			continue
		}

		if strings.Contains(name, "create") || strings.Contains(name, "update") || strings.Contains(name, "assign") || name == "jira_issue_transition" || name == "jira_issue_comment_add" {
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || *tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be annotated as non-destructive mutation", name)
			}
			continue
		}

		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Fatalf("expected MCP tool %q to be annotated as read-only", name)
		}
	}
}

func TestMCPJiraIssueMutationTools(t *testing.T) {
	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	jc := NewJiraClient(config)
	template, err := getJiraTestIssueTemplate(jc, getJiraTestTemplateIssueKey(t))
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newJiraMCPHTTPHandler(jc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	session := connectJiraMCPTestClient(t, server.URL+defaultMCPEndpoint)
	defer session.Close()

	summary := fmt.Sprintf("mcp parity create %d", time.Now().UnixNano())
	createdResult := callJiraMCPTool(t, session, "jira_issue_create", map[string]any{
		"projectId":   template.Fields.Project.ID,
		"issueTypeId": template.Fields.IssueType.ID,
		"summary":     summary,
		"description": "created from MCP parity test",
	})

	var created JiraCreatedIssue
	decodeMCPStructuredContent(t, createdResult, &created)
	if created.Key == "" {
		t.Fatal("expected jira_issue_create to return an issue key")
	}

	deleted := false
	t.Cleanup(func() {
		if !deleted {
			if err := deleteJiraTestIssue(jc, created.Key); err != nil {
				t.Errorf("failed to delete MCP test issue %s: %v", created.Key, err)
			}
		}
	})

	updatedSummary := fmt.Sprintf("mcp parity update %d", time.Now().UnixNano())
	updateResult := callJiraMCPTool(t, session, "jira_issue_update", map[string]any{
		"issueKey": created.Key,
		"summary":  updatedSummary,
	})
	var updateOutput jiraIssueUpdateToolOutput
	decodeMCPStructuredContent(t, updateResult, &updateOutput)
	if len(updateOutput.UpdatedFields) != 1 || updateOutput.UpdatedFields[0] != "summary" {
		t.Fatalf("expected summary update output, got %#v", updateOutput.UpdatedFields)
	}

	issueResult := callJiraMCPTool(t, session, "jira_issue_get", map[string]any{
		"issueKey": created.Key,
	})
	var issue JiraIssue
	decodeMCPStructuredContent(t, issueResult, &issue)
	if issue.Fields.Summary != updatedSummary {
		t.Fatalf("expected updated summary %q, got %q", updatedSummary, issue.Fields.Summary)
	}

	descriptionResult := callJiraMCPTool(t, session, "jira_issue_update", map[string]any{
		"issueKey":    created.Key,
		"description": "description updated from MCP parity test",
	})
	var descriptionOutput jiraIssueUpdateToolOutput
	decodeMCPStructuredContent(t, descriptionResult, &descriptionOutput)
	if len(descriptionOutput.UpdatedFields) != 1 || descriptionOutput.UpdatedFields[0] != "description" {
		t.Fatalf("expected description update output, got %#v", descriptionOutput.UpdatedFields)
	}

	fieldSummary := fmt.Sprintf("mcp parity field update %d", time.Now().UnixNano())
	fieldResult := callJiraMCPTool(t, session, "jira_issue_update", map[string]any{
		"issueKey": created.Key,
		"fields": map[string]any{
			"summary": fieldSummary,
		},
	})
	var fieldOutput jiraIssueUpdateToolOutput
	decodeMCPStructuredContent(t, fieldResult, &fieldOutput)
	if len(fieldOutput.UpdatedFields) != 1 || fieldOutput.UpdatedFields[0] != "summary" {
		t.Fatalf("expected field summary update output, got %#v", fieldOutput.UpdatedFields)
	}

	fieldJSONSummary := fmt.Sprintf("mcp parity fieldJson update %d", time.Now().UnixNano())
	fieldJSONResult := callJiraMCPTool(t, session, "jira_issue_update", map[string]any{
		"issueKey": created.Key,
		"fieldJson": map[string]any{
			"summary": fieldJSONSummary,
		},
	})
	var fieldJSONOutput jiraIssueUpdateToolOutput
	decodeMCPStructuredContent(t, fieldJSONResult, &fieldJSONOutput)
	if len(fieldJSONOutput.UpdatedFields) != 1 || fieldJSONOutput.UpdatedFields[0] != "summary" {
		t.Fatalf("expected fieldJson summary update output, got %#v", fieldJSONOutput.UpdatedFields)
	}

	searchResult := callJiraMCPTool(t, session, "jira_issue_search", map[string]any{
		"jql": "issuekey = " + created.Key,
	})
	var searchOutput JiraIssueSearchResult
	decodeMCPStructuredContent(t, searchResult, &searchOutput)
	if len(searchOutput.Issues) != 1 || searchOutput.Issues[0].Key != created.Key {
		t.Fatalf("expected jira_issue_search to find %s, got %#v", created.Key, searchOutput.Issues)
	}

	editMetaResult := callJiraMCPTool(t, session, "jira_issue_editmeta", map[string]any{
		"issueKey": created.Key,
	})
	var editMeta JiraIssueEditMeta
	decodeMCPStructuredContent(t, editMetaResult, &editMeta)
	if len(editMeta.Fields) == 0 {
		t.Fatal("expected jira_issue_editmeta to return fields")
	}

	commentResult := callJiraMCPTool(t, session, "jira_issue_comment_add", map[string]any{
		"issueKey": created.Key,
		"body":     "comment from MCP parity test",
	})
	var commentOutput jiraIssueCommentAddToolOutput
	decodeMCPStructuredContent(t, commentResult, &commentOutput)
	if commentOutput.Comment.ID == "" {
		t.Fatal("expected jira_issue_comment_add to return a comment ID")
	}

	callJiraMCPTool(t, session, "jira_issue_comment_delete", map[string]any{
		"issueKey":  created.Key,
		"commentId": commentOutput.Comment.ID,
	})

	assignMeResult := callJiraMCPTool(t, session, "jira_issue_assign", map[string]any{
		"issueKey": created.Key,
		"me":       true,
	})
	var assignMeOutput jiraIssueAssignToolOutput
	decodeMCPStructuredContent(t, assignMeResult, &assignMeOutput)
	if assignMeOutput.Assignee == "" || assignMeOutput.Assignee == "Unassigned" {
		t.Fatalf("expected jira_issue_assign me to assign a user, got %q", assignMeOutput.Assignee)
	}

	authResult := callJiraMCPTool(t, session, "jira_auth_whoami", map[string]any{})
	var myself JiraMyself
	decodeMCPStructuredContent(t, authResult, &myself)
	if myself.AccountID == "" {
		t.Fatal("expected jira_auth_whoami to return an account ID")
	}

	assignAccountResult := callJiraMCPTool(t, session, "jira_issue_assign", map[string]any{
		"issueKey":  created.Key,
		"accountId": myself.AccountID,
	})
	var assignAccountOutput jiraIssueAssignToolOutput
	decodeMCPStructuredContent(t, assignAccountResult, &assignAccountOutput)
	if assignAccountOutput.Assignee == "" || assignAccountOutput.Assignee == "Unassigned" {
		t.Fatalf("expected jira_issue_assign accountId to assign a user, got %q", assignAccountOutput.Assignee)
	}

	callJiraMCPTool(t, session, "jira_issue_assign", map[string]any{
		"issueKey":   created.Key,
		"unassigned": true,
	})

	transitionsResult := callJiraMCPTool(t, session, "jira_issue_transition_list", map[string]any{
		"issueKey": created.Key,
	})
	var transitions JiraTransitionList
	decodeMCPStructuredContent(t, transitionsResult, &transitions)
	if len(transitions.Transitions) == 0 {
		t.Fatal("expected jira_issue_transition_list to return transitions")
	}

	transitionResult := callJiraMCPTool(t, session, "jira_issue_transition", map[string]any{
		"issueKey":     created.Key,
		"transitionId": transitions.Transitions[0].ID,
	})
	var transitionOutput jiraIssueTransitionToolOutput
	decodeMCPStructuredContent(t, transitionResult, &transitionOutput)
	if transitionOutput.Transition.ID != transitions.Transitions[0].ID {
		t.Fatalf("expected transition %q, got %#v", transitions.Transitions[0].ID, transitionOutput.Transition)
	}

	callJiraMCPTool(t, session, "jira_issue_delete", map[string]any{
		"issueKey": created.Key,
	})
	deleted = true

	if _, err := jc.GetIssue(created.Key); err == nil {
		t.Fatal("expected deleted issue lookup to fail")
	}
}

func TestMCPJiraProjectList(t *testing.T) {
	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	session := newJiraMCPTestSession(t, NewJiraClient(config))
	defer session.Close()

	result := callJiraMCPTool(t, session, "jira_project_list", map[string]any{})
	var output jiraProjectListToolOutput
	decodeMCPStructuredContent(t, result, &output)
	if len(output.Projects) == 0 {
		t.Fatal("expected jira_project_list to return projects")
	}
}

func TestMCPJiraIssueTransitionByName(t *testing.T) {
	env, issueKey := setupIsolatedIssue(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	session := newJiraMCPTestSession(t, NewJiraClient(config))
	defer session.Close()

	transitionsResult := callJiraMCPTool(t, session, "jira_issue_transition_list", map[string]any{
		"issueKey": issueKey,
	})
	var transitions JiraTransitionList
	decodeMCPStructuredContent(t, transitionsResult, &transitions)
	if len(transitions.Transitions) == 0 {
		t.Skip("skipping: issue has no available transitions")
	}

	transitionResult := callJiraMCPTool(t, session, "jira_issue_transition", map[string]any{
		"issueKey":       issueKey,
		"transitionName": transitions.Transitions[0].Name,
	})
	var transitionOutput jiraIssueTransitionToolOutput
	decodeMCPStructuredContent(t, transitionResult, &transitionOutput)
	if transitionOutput.Transition.Name != transitions.Transitions[0].Name {
		t.Fatalf("expected transition %q, got %#v", transitions.Transitions[0].Name, transitionOutput.Transition)
	}
}

func TestMCPJiraIssueTransitionToolValidation(t *testing.T) {
	session := newLocalJiraMCPTestSession(t)
	defer session.Close()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "jira_issue_transition",
		Arguments: map[string]any{
			"issueKey":       "PROJ-123",
			"transitionName": "Done",
			"transitionId":   "31",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected conflicting transition selectors to return a tool error")
	}
}

func TestMCPJiraToolValidation(t *testing.T) {
	session := newLocalJiraMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{
			name: "create missing project selector",
			tool: "jira_issue_create",
			arguments: map[string]any{
				"issueType": "Task",
				"summary":   "summary",
			},
		},
		{
			name: "create missing issue type selector",
			tool: "jira_issue_create",
			arguments: map[string]any{
				"projectKey": "KAN",
				"summary":    "summary",
			},
		},
		{
			name: "create missing summary",
			tool: "jira_issue_create",
			arguments: map[string]any{
				"projectKey": "KAN",
				"issueType":  "Task",
			},
		},
		{
			name: "comment add missing body",
			tool: "jira_issue_comment_add",
			arguments: map[string]any{
				"issueKey": jiraTestValidationIssueKey,
			},
		},
		{
			name: "comment delete missing comment id",
			tool: "jira_issue_comment_delete",
			arguments: map[string]any{
				"issueKey": jiraTestValidationIssueKey,
			},
		},
		{
			name: "assign missing target",
			tool: "jira_issue_assign",
			arguments: map[string]any{
				"issueKey": jiraTestValidationIssueKey,
			},
		},
		{
			name: "assign conflicting target",
			tool: "jira_issue_assign",
			arguments: map[string]any{
				"issueKey":   jiraTestValidationIssueKey,
				"me":         true,
				"unassigned": true,
			},
		},
		{
			name: "transition missing target",
			tool: "jira_issue_transition",
			arguments: map[string]any{
				"issueKey": jiraTestValidationIssueKey,
			},
		},
		{
			name: "transition conflicting target",
			tool: "jira_issue_transition",
			arguments: map[string]any{
				"issueKey":       jiraTestValidationIssueKey,
				"transitionName": "Done",
				"transitionId":   "31",
			},
		},
		{
			name: "update missing fields",
			tool: "jira_issue_update",
			arguments: map[string]any{
				"issueKey": jiraTestValidationIssueKey,
			},
		},
		{
			name: "update duplicate field",
			tool: "jira_issue_update",
			arguments: map[string]any{
				"issueKey": jiraTestValidationIssueKey,
				"summary":  "a",
				"fields": map[string]any{
					"summary": "b",
				},
			},
		},
		{
			name: "update blank field value",
			tool: "jira_issue_update",
			arguments: map[string]any{
				"issueKey": jiraTestValidationIssueKey,
				"fields": map[string]any{
					"summary": "",
				},
			},
		},
		{
			name:      "delete missing issue key",
			tool:      "jira_issue_delete",
			arguments: map[string]any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
				Name:      tc.tool,
				Arguments: tc.arguments,
			})
			if err != nil {
				t.Fatal(err)
			}
			if !result.IsError {
				t.Fatal("expected tool validation to return a tool error")
			}
		})
	}
}

func TestMCPJiraStatusToolValidation(t *testing.T) {
	session := newLocalJiraMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{
			name:      "status get missing ids",
			tool:      "jira_status_get",
			arguments: map[string]any{},
		},
		{
			name: "status search negative start",
			tool: "jira_status_search",
			arguments: map[string]any{
				"startAt": -1,
			},
		},
		{
			name: "status create invalid scope",
			tool: "jira_status_create",
			arguments: map[string]any{
				"scopeType": "TEAM",
				"statuses": []map[string]any{
					{
						"name":           "Ready",
						"statusCategory": "TODO",
					},
				},
			},
		},
		{
			name: "status create project without project id",
			tool: "jira_status_create",
			arguments: map[string]any{
				"scopeType": "PROJECT",
				"statuses": []map[string]any{
					{
						"name":           "Ready",
						"statusCategory": "TODO",
					},
				},
			},
		},
		{
			name: "status update missing changes",
			tool: "jira_status_update",
			arguments: map[string]any{
				"statuses": []map[string]any{
					{
						"id": "10000",
					},
				},
			},
		},
		{
			name:      "status delete missing ids",
			tool:      "jira_status_delete",
			arguments: map[string]any{},
		},
		{
			name:      "board configuration missing board id",
			tool:      "jira_board_configuration_get",
			arguments: map[string]any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
				Name:      tc.tool,
				Arguments: tc.arguments,
			})
			if err != nil {
				t.Fatal(err)
			}
			if !result.IsError {
				t.Fatal("expected tool validation to return a tool error")
			}
		})
	}
}

func TestMCPAuthWhoAmI(t *testing.T) {
	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newJiraMCPHTTPHandler(NewJiraClient(config)))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	session := connectJiraMCPTestClient(t, server.URL+defaultMCPEndpoint)
	defer session.Close()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "jira_auth_whoami",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.IsError {
		t.Fatalf("expected jira_auth_whoami to succeed, got tool error: %#v", result.Content)
	}

	var myself JiraMyself
	rawMyself, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(rawMyself, &myself); err != nil {
		t.Fatal(err)
	}

	if myself.AccountID == "" {
		t.Fatal("expected jira_auth_whoami structured output to include accountId")
	}
}

func connectJiraMCPTestClient(t *testing.T, endpoint string) *mcp.ClientSession {
	t.Helper()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "jira-test-client",
		Version: "dev",
	}, nil)

	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{
		Endpoint: endpoint,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	return session
}

func newLocalJiraMCPTestSession(t *testing.T) *mcp.ClientSession {
	t.Helper()

	jc := NewJiraClient(Config{BaseURL: "http://127.0.0.1"})
	return newJiraMCPTestSession(t, jc)
}

func newJiraMCPTestSession(t *testing.T, jc JiraClient) *mcp.ClientSession {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(defaultMCPEndpoint, newJiraMCPHTTPHandler(jc))
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return connectJiraMCPTestClient(t, server.URL+defaultMCPEndpoint)
}

func callJiraMCPTool(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
	t.Helper()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      name,
		Arguments: arguments,
	})
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

	rawOutput, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(rawOutput, output); err != nil {
		t.Fatal(err)
	}
}
