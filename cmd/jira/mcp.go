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

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "jira",
		Version: Version,
	}, nil)

	readOnlyAnnotations := &mcp.ToolAnnotations{
		ReadOnlyHint:  true,
		OpenWorldHint: &openWorld,
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

	return server
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
  jira_auth_whoami     %s
  jira_project_list    %s
  jira_issue_get       %s
  jira_issue_search    %s
`, defaultMCPAddr, defaultMCPEndpoint, jiraAuthWhoAmIDoc.Summary, jiraProjectListDoc.Summary, jiraIssueGetDoc.Summary, jiraIssueSearchDoc.Summary)
}
