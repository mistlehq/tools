package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mistlehq/tools/internal/argparse"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultMCPAddr     = "127.0.0.1:7345"
	defaultMCPEndpoint = "/mcp"
)

type jiraMCPConfig struct {
	Addr     string
	Endpoint string
}

type jiraEmptyToolInput struct{}

type jiraIssueGetToolInput struct {
	IssueKey string `json:"issueKey" jsonschema:"Jira issue key, for example PROJ-123"`
}

type jiraIssueSearchToolInput struct {
	JQL string `json:"jql" jsonschema:"Jira Query Language expression, for example project = PROJ ORDER BY updated DESC"`
}

type jiraIssueCreateToolInput struct {
	ProjectKey    string  `json:"projectKey,omitempty" jsonschema:"Jira project key. Exactly one of projectKey or projectId is required."`
	ProjectID     string  `json:"projectId,omitempty" jsonschema:"Jira project ID. Exactly one of projectKey or projectId is required."`
	IssueTypeName string  `json:"issueType,omitempty" jsonschema:"Jira issue type name. Exactly one of issueType or issueTypeId is required."`
	IssueTypeID   string  `json:"issueTypeId,omitempty" jsonschema:"Jira issue type ID. Exactly one of issueType or issueTypeId is required."`
	Summary       string  `json:"summary" jsonschema:"Issue summary."`
	Description   *string `json:"description,omitempty" jsonschema:"Optional plain text issue description. The CLI converts it to Atlassian Document Format."`
}

type jiraIssueDeleteToolInput struct {
	IssueKey string `json:"issueKey" jsonschema:"Jira issue key, for example PROJ-123"`
}

type jiraIssueDeleteToolOutput struct {
	IssueKey string `json:"issueKey"`
	Deleted  bool   `json:"deleted"`
}

type jiraIssueCommentAddToolInput struct {
	IssueKey string `json:"issueKey" jsonschema:"Jira issue key, for example PROJ-123"`
	Body     string `json:"body" jsonschema:"Plain text comment body. The CLI converts it to Atlassian Document Format."`
}

type jiraIssueCommentAddToolOutput struct {
	IssueKey string      `json:"issueKey"`
	Comment  JiraComment `json:"comment"`
}

type jiraIssueCommentDeleteToolInput struct {
	IssueKey  string `json:"issueKey" jsonschema:"Jira issue key, for example PROJ-123"`
	CommentID string `json:"commentId" jsonschema:"Jira comment ID to delete."`
}

type jiraIssueCommentDeleteToolOutput struct {
	IssueKey  string `json:"issueKey"`
	CommentID string `json:"commentId"`
	Deleted   bool   `json:"deleted"`
}

type jiraIssueAssignToolInput struct {
	IssueKey   string  `json:"issueKey" jsonschema:"Jira issue key, for example PROJ-123"`
	Me         bool    `json:"me,omitempty" jsonschema:"Assign the issue to the current Jira user. Exactly one of me, accountId, or unassigned is required."`
	AccountID  *string `json:"accountId,omitempty" jsonschema:"Jira account ID to assign. Exactly one of me, accountId, or unassigned is required."`
	Unassigned bool    `json:"unassigned,omitempty" jsonschema:"Clear the issue assignee. Exactly one of me, accountId, or unassigned is required."`
}

type jiraIssueAssignToolOutput struct {
	Issue    JiraIssue `json:"issue"`
	Assignee string    `json:"assignee"`
}

type jiraIssueTransitionListToolInput struct {
	IssueKey string `json:"issueKey" jsonschema:"Jira issue key, for example PROJ-123"`
}

type jiraIssueTransitionToolInput struct {
	IssueKey       string `json:"issueKey" jsonschema:"Jira issue key, for example PROJ-123"`
	TransitionName string `json:"transitionName,omitempty" jsonschema:"Exact Jira transition name. Exactly one of transitionName or transitionId is required."`
	TransitionID   string `json:"transitionId,omitempty" jsonschema:"Jira transition ID. Exactly one of transitionName or transitionId is required."`
}

type jiraIssueTransitionToolOutput struct {
	Issue      JiraIssue      `json:"issue"`
	Transition JiraTransition `json:"transition"`
}

type jiraIssueUpdateToolInput struct {
	IssueKey    string            `json:"issueKey" jsonschema:"Jira issue key, for example PROJ-123"`
	Summary     *string           `json:"summary,omitempty" jsonschema:"New issue summary."`
	Description *string           `json:"description,omitempty" jsonschema:"New plain text issue description. The CLI converts it to Atlassian Document Format."`
	Fields      map[string]string `json:"fields,omitempty" jsonschema:"Additional editable Jira fields as string values, keyed by field ID."`
	FieldJSON   map[string]any    `json:"fieldJson,omitempty" jsonschema:"Additional editable Jira fields as JSON values, keyed by field ID."`
}

type jiraIssueUpdateToolOutput struct {
	IssueKey      string   `json:"issueKey"`
	UpdatedFields []string `json:"updatedFields"`
}

type jiraIssueEditMetaToolInput struct {
	IssueKey string `json:"issueKey" jsonschema:"Jira issue key, for example PROJ-123"`
}

type jiraProjectListToolOutput struct {
	Projects []JiraProject `json:"projects"`
}

func (cli CLI) runMCP(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printMCPHelp()
		return nil
	}

	switch args[0] {
	case "serve":
		if isSingleHelpArg(args[1:]) {
			cli.printMCPServeHelp()
			return nil
		}

		return cli.runMCPServe(args[1:])
	default:
		return fmt.Errorf("unsupported mcp command: %s", args[0])
	}
}

func (cli CLI) runMCPServe(args []string) error {
	config, err := parseJiraMCPServeArgs(args)
	if err != nil {
		return err
	}

	jc, err := cli.jiraClient()
	if err != nil {
		return err
	}

	handler := newJiraMCPHTTPHandler(jc)
	mux := http.NewServeMux()
	mux.Handle(config.Endpoint, handler)

	fmt.Fprintf(cli.stdout, "Jira MCP server listening on http://%s%s\n", config.Addr, config.Endpoint)
	return http.ListenAndServe(config.Addr, mux)
}

func parseJiraMCPServeArgs(args []string) (jiraMCPConfig, error) {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"addr":     {TakesValue: true},
		"endpoint": {TakesValue: true},
	})
	if err != nil {
		return jiraMCPConfig{}, err
	}

	if len(parsedArgs.Positionals) > 0 {
		return jiraMCPConfig{}, fmt.Errorf("mcp serve does not accept positional arguments")
	}

	config := jiraMCPConfig{
		Addr:     defaultMCPAddr,
		Endpoint: defaultMCPEndpoint,
	}

	if addr := parsedArgs.First("addr"); addr != "" {
		config.Addr = addr
	}

	if endpoint := parsedArgs.First("endpoint"); endpoint != "" {
		config.Endpoint = endpoint
	}

	if strings.TrimSpace(config.Addr) == "" {
		return jiraMCPConfig{}, fmt.Errorf("--addr must not be empty")
	}

	if !strings.HasPrefix(config.Endpoint, "/") {
		return jiraMCPConfig{}, fmt.Errorf("--endpoint must start with '/'")
	}

	return config, nil
}

func newJiraMCPHTTPHandler(jc JiraClient) http.Handler {
	server := newJiraMCPServer(jc)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		CrossOriginProtection: &http.CrossOriginProtection{},
	})
}

func newJiraMCPServer(jc JiraClient) *mcp.Server {
	openWorld := true
	destructive := true
	notDestructive := false

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "jira",
		Version: Version,
	}, nil)

	readOnlyAnnotations := &mcp.ToolAnnotations{
		ReadOnlyHint:  true,
		OpenWorldHint: &openWorld,
	}
	mutatingAnnotations := &mcp.ToolAnnotations{
		OpenWorldHint:   &openWorld,
		DestructiveHint: &notDestructive,
	}
	destructiveAnnotations := &mcp.ToolAnnotations{
		OpenWorldHint:   &openWorld,
		DestructiveHint: &destructive,
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_auth_whoami",
		Title:       jiraAuthWhoAmIDoc.Command,
		Description: jiraAuthWhoAmIDoc.Description,
		Annotations: readOnlyAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ *jiraEmptyToolInput) (*mcp.CallToolResult, JiraMyself, error) {
		myself, err := jc.GetMyselfContext(ctx)
		return nil, myself, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_project_list",
		Title:       jiraProjectListDoc.Command,
		Description: jiraProjectListDoc.Description,
		Annotations: readOnlyAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ *jiraEmptyToolInput) (*mcp.CallToolResult, jiraProjectListToolOutput, error) {
		projectList, err := jc.ListProjectsContext(ctx)
		if err != nil {
			return nil, jiraProjectListToolOutput{}, err
		}

		return nil, jiraProjectListToolOutput{
			Projects: projectList.Values,
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_get",
		Title:       jiraIssueGetDoc.Command,
		Description: jiraIssueGetDoc.Description,
		Annotations: readOnlyAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueGetToolInput) (*mcp.CallToolResult, JiraIssue, error) {
		if strings.TrimSpace(input.IssueKey) == "" {
			return nil, JiraIssue{}, fmt.Errorf("issueKey is required")
		}

		issue, err := jc.GetIssueContext(ctx, input.IssueKey)
		return nil, issue, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_search",
		Title:       jiraIssueSearchDoc.Command,
		Description: jiraIssueSearchDoc.Description,
		Annotations: readOnlyAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueSearchToolInput) (*mcp.CallToolResult, JiraIssueSearchResult, error) {
		if strings.TrimSpace(input.JQL) == "" {
			return nil, JiraIssueSearchResult{}, fmt.Errorf("jql is required")
		}

		searchResult, err := jc.SearchIssuesContext(ctx, input.JQL)
		return nil, searchResult, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_create",
		Title:       jiraIssueCreateDoc.Command,
		Description: jiraIssueCreateDoc.Description,
		Annotations: mutatingAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueCreateToolInput) (*mcp.CallToolResult, JiraCreatedIssue, error) {
		if err := validateJiraIssueCreateToolInput(input); err != nil {
			return nil, JiraCreatedIssue{}, err
		}

		issue, err := jc.CreateIssueContext(ctx, CreateIssueInput{
			ProjectID:     input.ProjectID,
			ProjectKey:    input.ProjectKey,
			IssueTypeID:   input.IssueTypeID,
			IssueTypeName: input.IssueTypeName,
			Summary:       input.Summary,
			Description:   input.Description,
		})
		return nil, issue, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_delete",
		Title:       jiraIssueDeleteDoc.Command,
		Description: jiraIssueDeleteDoc.Description,
		Annotations: destructiveAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueDeleteToolInput) (*mcp.CallToolResult, jiraIssueDeleteToolOutput, error) {
		if strings.TrimSpace(input.IssueKey) == "" {
			return nil, jiraIssueDeleteToolOutput{}, fmt.Errorf("issueKey is required")
		}

		if err := jc.DeleteIssueContext(ctx, input.IssueKey); err != nil {
			return nil, jiraIssueDeleteToolOutput{}, err
		}

		return nil, jiraIssueDeleteToolOutput{
			IssueKey: input.IssueKey,
			Deleted:  true,
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_comment_add",
		Title:       jiraIssueCommentAddDoc.Command,
		Description: jiraIssueCommentAddDoc.Description,
		Annotations: mutatingAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueCommentAddToolInput) (*mcp.CallToolResult, jiraIssueCommentAddToolOutput, error) {
		if strings.TrimSpace(input.IssueKey) == "" {
			return nil, jiraIssueCommentAddToolOutput{}, fmt.Errorf("issueKey is required")
		}
		if strings.TrimSpace(input.Body) == "" {
			return nil, jiraIssueCommentAddToolOutput{}, fmt.Errorf("body is required and must not be empty")
		}

		comment, err := jc.AddIssueCommentContext(ctx, input.IssueKey, AddCommentInput{
			Body: input.Body,
		})
		if err != nil {
			return nil, jiraIssueCommentAddToolOutput{}, err
		}

		return nil, jiraIssueCommentAddToolOutput{
			IssueKey: input.IssueKey,
			Comment:  comment,
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_comment_delete",
		Title:       jiraIssueCommentDeleteDoc.Command,
		Description: jiraIssueCommentDeleteDoc.Description,
		Annotations: destructiveAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueCommentDeleteToolInput) (*mcp.CallToolResult, jiraIssueCommentDeleteToolOutput, error) {
		if strings.TrimSpace(input.IssueKey) == "" {
			return nil, jiraIssueCommentDeleteToolOutput{}, fmt.Errorf("issueKey is required")
		}
		if strings.TrimSpace(input.CommentID) == "" {
			return nil, jiraIssueCommentDeleteToolOutput{}, fmt.Errorf("commentId is required")
		}

		if err := jc.DeleteIssueCommentContext(ctx, input.IssueKey, input.CommentID); err != nil {
			return nil, jiraIssueCommentDeleteToolOutput{}, err
		}

		return nil, jiraIssueCommentDeleteToolOutput{
			IssueKey:  input.IssueKey,
			CommentID: input.CommentID,
			Deleted:   true,
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_assign",
		Title:       jiraIssueAssignDoc.Command,
		Description: jiraIssueAssignDoc.Description,
		Annotations: mutatingAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueAssignToolInput) (*mcp.CallToolResult, jiraIssueAssignToolOutput, error) {
		assignInput, err := buildJiraMCPAssignInput(ctx, jc, input)
		if err != nil {
			return nil, jiraIssueAssignToolOutput{}, err
		}

		if err := jc.AssignIssueContext(ctx, input.IssueKey, assignInput); err != nil {
			return nil, jiraIssueAssignToolOutput{}, err
		}

		issue, err := jc.GetIssueContext(ctx, input.IssueKey)
		if err != nil {
			return nil, jiraIssueAssignToolOutput{}, err
		}

		assignee := "Unassigned"
		if issue.Fields.Assignee != nil {
			assignee = issue.Fields.Assignee.DisplayName
		}

		return nil, jiraIssueAssignToolOutput{
			Issue:    issue,
			Assignee: assignee,
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_transition_list",
		Title:       jiraIssueTransitionListDoc.Command,
		Description: jiraIssueTransitionListDoc.Description,
		Annotations: readOnlyAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueTransitionListToolInput) (*mcp.CallToolResult, JiraTransitionList, error) {
		if strings.TrimSpace(input.IssueKey) == "" {
			return nil, JiraTransitionList{}, fmt.Errorf("issueKey is required")
		}

		transitionList, err := jc.ListIssueTransitionsContext(ctx, input.IssueKey)
		return nil, transitionList, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_transition",
		Title:       jiraIssueTransitionDoc.Command,
		Description: jiraIssueTransitionDoc.Description,
		Annotations: mutatingAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueTransitionToolInput) (*mcp.CallToolResult, jiraIssueTransitionToolOutput, error) {
		if strings.TrimSpace(input.IssueKey) == "" {
			return nil, jiraIssueTransitionToolOutput{}, fmt.Errorf("issueKey is required")
		}

		selectedTransition, err := selectJiraMCPTransition(ctx, jc, input)
		if err != nil {
			return nil, jiraIssueTransitionToolOutput{}, err
		}

		if err := jc.TransitionIssueContext(ctx, input.IssueKey, TransitionIssueInput{
			TransitionID: selectedTransition.ID,
		}); err != nil {
			return nil, jiraIssueTransitionToolOutput{}, err
		}

		issue, err := jc.GetIssueContext(ctx, input.IssueKey)
		if err != nil {
			return nil, jiraIssueTransitionToolOutput{}, err
		}

		return nil, jiraIssueTransitionToolOutput{
			Issue:      issue,
			Transition: selectedTransition,
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_update",
		Title:       jiraIssueUpdateDoc.Command,
		Description: jiraIssueUpdateDoc.Description,
		Annotations: mutatingAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueUpdateToolInput) (*mcp.CallToolResult, jiraIssueUpdateToolOutput, error) {
		updateInput, updatedFields, err := buildJiraMCPUpdateInput(input)
		if err != nil {
			return nil, jiraIssueUpdateToolOutput{}, err
		}

		if err := jc.UpdateIssueContext(ctx, input.IssueKey, updateInput); err != nil {
			return nil, jiraIssueUpdateToolOutput{}, err
		}

		return nil, jiraIssueUpdateToolOutput{
			IssueKey:      input.IssueKey,
			UpdatedFields: updatedFields,
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "jira_issue_editmeta",
		Title:       jiraIssueEditMetaDoc.Command,
		Description: jiraIssueEditMetaDoc.Description,
		Annotations: readOnlyAnnotations,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input *jiraIssueEditMetaToolInput) (*mcp.CallToolResult, JiraIssueEditMeta, error) {
		if strings.TrimSpace(input.IssueKey) == "" {
			return nil, JiraIssueEditMeta{}, fmt.Errorf("issueKey is required")
		}

		editMeta, err := jc.GetIssueEditMetaContext(ctx, input.IssueKey)
		return nil, editMeta, err
	})

	return server
}

func validateJiraIssueCreateToolInput(input *jiraIssueCreateToolInput) error {
	projectSelectors := 0
	if input.ProjectKey != "" {
		projectSelectors++
	}
	if input.ProjectID != "" {
		projectSelectors++
	}
	if projectSelectors != 1 {
		return fmt.Errorf("exactly one of projectKey or projectId is required")
	}
	if strings.TrimSpace(input.ProjectKey) == "" && input.ProjectKey != "" {
		return fmt.Errorf("projectKey must not be empty")
	}
	if strings.TrimSpace(input.ProjectID) == "" && input.ProjectID != "" {
		return fmt.Errorf("projectId must not be empty")
	}

	issueTypeSelectors := 0
	if input.IssueTypeName != "" {
		issueTypeSelectors++
	}
	if input.IssueTypeID != "" {
		issueTypeSelectors++
	}
	if issueTypeSelectors != 1 {
		return fmt.Errorf("exactly one of issueType or issueTypeId is required")
	}
	if strings.TrimSpace(input.IssueTypeName) == "" && input.IssueTypeName != "" {
		return fmt.Errorf("issueType must not be empty")
	}
	if strings.TrimSpace(input.IssueTypeID) == "" && input.IssueTypeID != "" {
		return fmt.Errorf("issueTypeId must not be empty")
	}

	if strings.TrimSpace(input.Summary) == "" {
		return fmt.Errorf("summary is required and must not be empty")
	}

	return nil
}

func buildJiraMCPAssignInput(ctx context.Context, jc JiraClient, input *jiraIssueAssignToolInput) (AssignIssueInput, error) {
	if strings.TrimSpace(input.IssueKey) == "" {
		return AssignIssueInput{}, fmt.Errorf("issueKey is required")
	}

	selectedTargets := 0
	if input.Me {
		selectedTargets++
	}
	if input.AccountID != nil {
		selectedTargets++
	}
	if input.Unassigned {
		selectedTargets++
	}
	if selectedTargets != 1 {
		return AssignIssueInput{}, fmt.Errorf("exactly one of me, accountId, or unassigned is required")
	}

	assignInput := AssignIssueInput{}
	if input.Me {
		myself, err := jc.GetMyselfContext(ctx)
		if err != nil {
			return AssignIssueInput{}, err
		}

		assignInput.AccountID = &myself.AccountID
	} else if input.AccountID != nil {
		accountID := *input.AccountID
		if strings.TrimSpace(accountID) == "" {
			return AssignIssueInput{}, fmt.Errorf("accountId must not be empty")
		}
		assignInput.AccountID = &accountID
	}

	return assignInput, nil
}

func selectJiraMCPTransition(ctx context.Context, jc JiraClient, input *jiraIssueTransitionToolInput) (JiraTransition, error) {
	selectedTargets := 0
	if input.TransitionName != "" {
		selectedTargets++
	}
	if input.TransitionID != "" {
		selectedTargets++
	}
	if selectedTargets != 1 {
		return JiraTransition{}, fmt.Errorf("exactly one of transitionName or transitionId is required")
	}

	transitionList, err := jc.ListIssueTransitionsContext(ctx, input.IssueKey)
	if err != nil {
		return JiraTransition{}, err
	}

	parsedArgs := argparse.Parsed{
		Flags: make(map[string][]string),
	}
	if input.TransitionID != "" {
		parsedArgs.Flags["to-id"] = []string{input.TransitionID}
	} else {
		parsedArgs.Flags["to"] = []string{input.TransitionName}
	}

	return selectTransition(input.IssueKey, transitionList.Transitions, parsedArgs)
}

func buildJiraMCPUpdateInput(input *jiraIssueUpdateToolInput) (UpdateIssueInput, []string, error) {
	if strings.TrimSpace(input.IssueKey) == "" {
		return UpdateIssueInput{}, nil, fmt.Errorf("issueKey is required")
	}

	updateInput := UpdateIssueInput{
		Fields: make(map[string]any),
	}
	updatedFields := make([]string, 0, 2+len(input.Fields)+len(input.FieldJSON))
	seenFields := make(map[string]struct{})
	addField := func(fieldID string, value any) error {
		if strings.TrimSpace(fieldID) == "" {
			return fmt.Errorf("field id must not be empty")
		}
		if _, ok := seenFields[fieldID]; ok {
			return fmt.Errorf("field %q was provided more than once", fieldID)
		}

		seenFields[fieldID] = struct{}{}
		updateInput.Fields[fieldID] = value
		updatedFields = append(updatedFields, fieldID)
		return nil
	}

	if input.Summary != nil {
		if strings.TrimSpace(*input.Summary) == "" {
			return UpdateIssueInput{}, nil, fmt.Errorf("summary must not be empty")
		}
		if err := addField("summary", *input.Summary); err != nil {
			return UpdateIssueInput{}, nil, err
		}
	}

	if input.Description != nil {
		updateInput.Description = input.Description
		seenFields["description"] = struct{}{}
		updatedFields = append(updatedFields, "description")
	}

	for fieldID, value := range input.Fields {
		if strings.TrimSpace(value) == "" {
			return UpdateIssueInput{}, nil, fmt.Errorf("field value for %s must not be empty; use fieldJson with null to clear a field", fieldID)
		}
		if err := addField(fieldID, value); err != nil {
			return UpdateIssueInput{}, nil, err
		}
	}

	for fieldID, value := range input.FieldJSON {
		if err := addField(fieldID, value); err != nil {
			return UpdateIssueInput{}, nil, err
		}
	}

	if len(updatedFields) == 0 {
		return UpdateIssueInput{}, nil, fmt.Errorf("issue update requires at least one of summary, description, fields, or fieldJson")
	}

	return updateInput, updatedFields, nil
}

func (cli CLI) printMCPHelp() {
	fmt.Fprint(cli.stdout, `jira mcp

Run Jira as a local MCP server.

Usage:
  jira mcp help
  jira mcp serve
  jira mcp serve --help

Commands:
  serve    Serve Jira MCP tools over Streamable HTTP
`)
}

func (cli CLI) printMCPServeHelp() {
	fmt.Fprintf(cli.stdout, `jira mcp serve

Serve Jira tools over MCP Streamable HTTP.

Usage:
  jira mcp serve
  jira mcp serve --addr <addr>
  jira mcp serve --endpoint <path>
  jira mcp serve --addr <addr> --endpoint <path>
  jira mcp serve --help

Options:
  --addr <addr>        Listen address. Defaults to %s.
  --endpoint <path>    MCP HTTP endpoint. Defaults to %s.

Tools:
  jira_auth_whoami              %s
  jira_project_list             %s
  jira_issue_create             %s
  jira_issue_get                %s
  jira_issue_search             %s
  jira_issue_delete             %s
  jira_issue_comment_add        %s
  jira_issue_comment_delete     %s
  jira_issue_assign             %s
  jira_issue_transition_list    %s
  jira_issue_transition         %s
  jira_issue_update             %s
  jira_issue_editmeta           %s
`, defaultMCPAddr, defaultMCPEndpoint, jiraAuthWhoAmIDoc.Summary, jiraProjectListDoc.Summary, jiraIssueCreateDoc.Summary, jiraIssueGetDoc.Summary, jiraIssueSearchDoc.Summary, jiraIssueDeleteDoc.Summary, jiraIssueCommentAddDoc.Summary, jiraIssueCommentDeleteDoc.Summary, jiraIssueAssignDoc.Summary, jiraIssueTransitionListDoc.Summary, jiraIssueTransitionDoc.Summary, jiraIssueUpdateDoc.Summary, jiraIssueEditMetaDoc.Summary)
}
