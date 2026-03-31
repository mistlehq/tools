package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Version is the current jira CLI version.
var Version = "dev"

type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	env    Environment
}

func (cli CLI) jiraClient() (JiraClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return JiraClient{}, err
	}

	return NewJiraClient(config), nil
}

func (cli CLI) run(args []string) error {
	if len(args) < 2 {
		cli.printHelp()
		return nil
	}

	switch args[1] {
	case "help", "-h", "--help":
		cli.printHelp()
		return nil
	case "version", "--version":
		fmt.Fprintln(cli.stdout, Version)
		return nil
	case "auth":
		return cli.runAuth(args[2:])
	case "project":
		return cli.runProject(args[2:])
	case "issue":
		return cli.runIssue(args[2:])
	default:
		return fmt.Errorf("unsupported command: %s", args[1])
	}
}

func (cli CLI) runAuth(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printAuthHelp()
		return nil
	}

	switch args[0] {
	case "whoami":
		jc, err := cli.jiraClient()
		if err != nil {
			return err
		}

		return cli.runAuthWhoAmI(jc, args[1:])
	default:
		return fmt.Errorf("unsupported auth command: %s", args[0])
	}
}

func (cli CLI) runAuthWhoAmI(jc JiraClient, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("whoami does not accept positional arguments")
	}

	myself, err := jc.GetMyself()
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.stdout, "Account ID: "+myself.AccountID)
	fmt.Fprintln(cli.stdout, "Display name: "+myself.DisplayName)
	fmt.Fprintln(cli.stdout, "Email: "+myself.Email)
	return nil
}

func (cli CLI) printAuthHelp() {
	fmt.Fprintln(cli.stdout, "jira auth")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Inspect Jira authentication state.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  jira auth help")
	fmt.Fprintln(cli.stdout, "  jira auth whoami")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Commands:")
	fmt.Fprintln(cli.stdout, "  whoami")
}

func (cli CLI) runProject(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printProjectHelp()
		return nil
	}

	switch args[0] {
	case "list":
		jc, err := cli.jiraClient()
		if err != nil {
			return err
		}

		return cli.runProjectList(jc, args[1:])
	default:
		return fmt.Errorf("unsupported project command: %s", args[0])
	}
}

func (cli CLI) runProjectList(jc JiraClient, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("project list does not accept positional arguments")
	}

	projectList, err := jc.ListProjects()
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.stdout, "ID\tKEY\tNAME")

	for _, project := range projectList.Values {
		fmt.Fprintf(cli.stdout, "%s\t%s\t%s\n", project.ID, project.Key, project.Name)
	}

	return nil
}

func (cli CLI) printProjectHelp() {
	fmt.Fprintln(cli.stdout, "jira project")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Inspect Jira projects.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  jira project help")
	fmt.Fprintln(cli.stdout, "  jira project list")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Commands:")
	fmt.Fprintln(cli.stdout, "  list")
}

func (cli CLI) runIssue(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printIssueHelp()
		return nil
	}

	switch args[0] {
	case "get":
		jc, err := cli.jiraClient()
		if err != nil {
			return err
		}

		return cli.runIssueGet(jc, args[1:])
	case "search":
		jc, err := cli.jiraClient()
		if err != nil {
			return err
		}

		return cli.runIssueSearch(jc, args[1:])
	case "comment":
		return cli.runIssueComment(args[1:])
	case "assign":
		return cli.runIssueAssign(args[1:])
	case "transition":
		return cli.runIssueTransition(args[1:])
	case "update":
		return cli.runIssueUpdate(args[1:])
	case "editmeta":
		return cli.runIssueEditMeta(args[1:])
	default:
		return fmt.Errorf("unsupported issue command: %s", args[0])
	}
}

func (cli CLI) runIssueComment(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printIssueCommentHelp()
		return nil
	}

	switch args[0] {
	case "add":
		return cli.runIssueCommentAdd(args[1:])
	default:
		return fmt.Errorf("unsupported issue comment command: %s", args[0])
	}
}

func (cli CLI) runIssueAssign(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printIssueAssignHelp()
		return nil
	}

	return cli.runIssueAssignSet(args)
}

func (cli CLI) runIssueTransition(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printIssueTransitionHelp()
		return nil
	}

	switch args[0] {
	case "list":
		return cli.runIssueTransitionList(args[1:])
	default:
		return cli.runIssueTransitionMove(args)
	}
}

func (cli CLI) runIssueUpdate(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printIssueUpdateHelp()
		return nil
	}

	return fmt.Errorf("issue update is not implemented yet")
}

func (cli CLI) runIssueEditMeta(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printIssueEditMetaHelp()
		return nil
	}

	return fmt.Errorf("issue editmeta is not implemented yet")
}

func (cli CLI) runIssueAssignSet(args []string) error {
	parsedArgs, err := parseArgs(args, map[string]argSpec{
		"me": {},
		"account-id": {
			takesValue: true,
		},
		"unassigned": {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.positionals) != 1 {
		return fmt.Errorf("issue assign expects exactly 1 positional argument")
	}

	selectedFlags := 0
	if parsedArgs.has("me") {
		selectedFlags++
	}

	if parsedArgs.has("account-id") {
		selectedFlags++
	}

	if parsedArgs.has("unassigned") {
		selectedFlags++
	}

	if selectedFlags != 1 {
		return fmt.Errorf("exactly one of --me, --account-id, or --unassigned is required")
	}

	jc, err := cli.jiraClient()
	if err != nil {
		return err
	}

	assignInput := AssignIssueInput{}
	if parsedArgs.has("me") {
		myself, err := jc.GetMyself()
		if err != nil {
			return err
		}

		assignInput.AccountID = &myself.AccountID
	} else if parsedArgs.has("account-id") {
		accountID := parsedArgs.first("account-id")
		assignInput.AccountID = &accountID
	}

	issueKey := parsedArgs.positionals[0]
	if err := jc.AssignIssue(issueKey, assignInput); err != nil {
		return err
	}

	issue, err := jc.GetIssue(issueKey)
	if err != nil {
		return err
	}

	assigneeName := "Unassigned"
	if issue.Fields.Assignee != nil {
		assigneeName = issue.Fields.Assignee.DisplayName
	}

	fmt.Fprintln(cli.stdout, "Issue: "+issue.Key)
	fmt.Fprintln(cli.stdout, "Assignee: "+assigneeName)
	return nil
}

func (cli CLI) runIssueTransitionList(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("issue transition list expects exactly 1 positional argument")
	}

	jc, err := cli.jiraClient()
	if err != nil {
		return err
	}

	transitionList, err := jc.ListIssueTransitions(args[0])
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.stdout, "ID\tNAME\tTO STATUS")
	for _, transition := range transitionList.Transitions {
		fmt.Fprintf(cli.stdout, "%s\t%s\t%s\n", transition.ID, transition.Name, transition.To.Name)
	}

	return nil
}

func (cli CLI) runIssueTransitionMove(args []string) error {
	parsedArgs, err := parseArgs(args, map[string]argSpec{
		"to": {
			takesValue: true,
		},
		"to-id": {
			takesValue: true,
		},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.positionals) != 1 {
		return fmt.Errorf("issue transition expects exactly 1 positional argument")
	}

	selectedFlags := 0
	if parsedArgs.has("to") {
		selectedFlags++
	}

	if parsedArgs.has("to-id") {
		selectedFlags++
	}

	if selectedFlags != 1 {
		return fmt.Errorf("exactly one of --to or --to-id is required")
	}

	issueKey := parsedArgs.positionals[0]
	jc, err := cli.jiraClient()
	if err != nil {
		return err
	}

	transitionList, err := jc.ListIssueTransitions(issueKey)
	if err != nil {
		return err
	}

	selectedTransition, err := selectTransition(issueKey, transitionList.Transitions, parsedArgs)
	if err != nil {
		return err
	}

	if err := jc.TransitionIssue(issueKey, TransitionIssueInput{
		TransitionID: selectedTransition.ID,
	}); err != nil {
		return err
	}

	issue, err := jc.GetIssue(issueKey)
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.stdout, "Issue: "+issue.Key)
	fmt.Fprintln(cli.stdout, "Transition: "+selectedTransition.Name)
	fmt.Fprintln(cli.stdout, "Status: "+issue.Fields.Status.Name)
	return nil
}

func selectTransition(issueKey string, transitions []JiraTransition, parsedArgs parsedArgs) (JiraTransition, error) {
	if parsedArgs.has("to-id") {
		transitionID := parsedArgs.first("to-id")
		for _, transition := range transitions {
			if transition.ID == transitionID {
				return transition, nil
			}
		}

		return JiraTransition{}, fmt.Errorf("no transition with id %q for issue %s", transitionID, issueKey)
	}

	transitionName := parsedArgs.first("to")
	matches := make([]JiraTransition, 0, len(transitions))
	for _, transition := range transitions {
		if transition.Name == transitionName {
			matches = append(matches, transition)
		}
	}

	switch len(matches) {
	case 0:
		return JiraTransition{}, fmt.Errorf("no transition named %q for issue %s", transitionName, issueKey)
	case 1:
		return matches[0], nil
	default:
		return JiraTransition{}, fmt.Errorf("multiple transitions named %q for issue %s; use 'jira issue transition list %s' or --to-id", transitionName, issueKey, issueKey)
	}
}

func (cli CLI) runIssueCommentAdd(args []string) error {
	parsedArgs, err := parseArgs(args, map[string]argSpec{
		"body": {
			takesValue: true,
		},
		"body-file": {
			takesValue: true,
		},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.positionals) != 1 {
		return fmt.Errorf("issue comment add expects exactly 1 positional argument")
	}

	bodyFlag := parsedArgs.first("body")
	bodyFileFlag := parsedArgs.first("body-file")

	body, err := cli.readTextInput("body", bodyFlag, "body-file", bodyFileFlag)
	if err != nil {
		return err
	}

	jc, err := cli.jiraClient()
	if err != nil {
		return err
	}

	comment, err := jc.AddIssueComment(parsedArgs.positionals[0], AddCommentInput{
		Body: body,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.stdout, "Issue: "+parsedArgs.positionals[0])
	fmt.Fprintln(cli.stdout, "Comment ID: "+comment.ID)
	fmt.Fprintln(cli.stdout, "Author: "+comment.Author.DisplayName)
	fmt.Fprintln(cli.stdout, "Created: "+comment.Created)
	return nil
}

// readTextInput normalizes a mutually exclusive `--value`/`--file` pair and
// supports `-` as stdin for file-backed input.
func (cli CLI) readTextInput(valueFlagName string, value string, fileFlagName string, filePath string) (string, error) {
	if value != "" && filePath != "" {
		return "", fmt.Errorf("--%s and --%s are mutually exclusive", valueFlagName, fileFlagName)
	}

	if value == "" && filePath == "" {
		return "", fmt.Errorf("exactly one of --%s or --%s is required", valueFlagName, fileFlagName)
	}

	if value != "" {
		if strings.TrimSpace(value) == "" {
			return "", fmt.Errorf("--%s must not be empty", valueFlagName)
		}

		return value, nil
	}

	var body []byte
	var err error

	switch filePath {
	case "-":
		body, err = io.ReadAll(cli.stdin)
	default:
		body, err = os.ReadFile(filePath)
	}

	if err != nil {
		return "", err
	}

	text := string(body)
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("--%s must not be empty", fileFlagName)
	}

	return text, nil
}

func (cli CLI) runIssueGet(jc JiraClient, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("issue get expects exactly 1 positional argument")
	}

	issue, err := jc.GetIssue(args[0])
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.stdout, "ID: "+issue.ID)
	fmt.Fprintln(cli.stdout, "Key: "+issue.Key)
	fmt.Fprintln(cli.stdout, "Summary: "+issue.Fields.Summary)
	fmt.Fprintln(cli.stdout, "Status: "+issue.Fields.Status.Name)

	return nil
}

func (cli CLI) runIssueSearch(jc JiraClient, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("issue search expects exactly 1 positional argument")
	}

	searchResult, err := jc.SearchIssues(args[0])
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.stdout, "KEY\tSTATUS\tSUMMARY")

	for _, issue := range searchResult.Issues {
		fmt.Fprintf(cli.stdout, "%s\t%s\t%s\n", issue.Key, issue.Fields.Status.Name, issue.Fields.Summary)
	}

	return nil
}

func (cli CLI) printIssueHelp() {
	fmt.Fprintln(cli.stdout, "jira issue")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Inspect and mutate Jira issues.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  jira issue help")
	fmt.Fprintln(cli.stdout, "  jira issue get <issue-key>")
	fmt.Fprintln(cli.stdout, "  jira issue search '<jql query>'")
	fmt.Fprintln(cli.stdout, "  jira issue comment help")
	fmt.Fprintln(cli.stdout, "  jira issue assign help")
	fmt.Fprintln(cli.stdout, "  jira issue transition help")
	fmt.Fprintln(cli.stdout, "  jira issue update help")
	fmt.Fprintln(cli.stdout, "  jira issue editmeta help")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Commands:")
	fmt.Fprintln(cli.stdout, "  get")
	fmt.Fprintln(cli.stdout, "  search")
	fmt.Fprintln(cli.stdout, "  comment")
	fmt.Fprintln(cli.stdout, "  assign")
	fmt.Fprintln(cli.stdout, "  transition")
	fmt.Fprintln(cli.stdout, "  update")
	fmt.Fprintln(cli.stdout, "  editmeta")
}

func (cli CLI) printIssueCommentHelp() {
	fmt.Fprintln(cli.stdout, "jira issue comment")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Add comments to Jira issues.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  jira issue comment help")
	fmt.Fprintln(cli.stdout, "  jira issue comment add <issue-key> --body <text>")
	fmt.Fprintln(cli.stdout, "  jira issue comment add <issue-key> --body-file <path>")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Commands:")
	fmt.Fprintln(cli.stdout, "  add")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Examples:")
	fmt.Fprintln(cli.stdout, "  jira issue comment add PROJ-123 --body 'Looks good'")
	fmt.Fprintln(cli.stdout, "  jira issue comment add PROJ-123 --body-file ./comment.txt")
}

func (cli CLI) printIssueAssignHelp() {
	fmt.Fprintln(cli.stdout, "jira issue assign")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Assign Jira issues.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  jira issue assign help")
	fmt.Fprintln(cli.stdout, "  jira issue assign <issue-key> --me")
	fmt.Fprintln(cli.stdout, "  jira issue assign <issue-key> --account-id <account-id>")
	fmt.Fprintln(cli.stdout, "  jira issue assign <issue-key> --unassigned")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Examples:")
	fmt.Fprintln(cli.stdout, "  jira issue assign PROJ-123 --me")
	fmt.Fprintln(cli.stdout, "  jira issue assign PROJ-123 --account-id 712020:abc123")
}

func (cli CLI) printIssueTransitionHelp() {
	fmt.Fprintln(cli.stdout, "jira issue transition")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Move Jira issues through workflow transitions.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  jira issue transition help")
	fmt.Fprintln(cli.stdout, "  jira issue transition list <issue-key>")
	fmt.Fprintln(cli.stdout, "  jira issue transition <issue-key> --to <transition-name>")
	fmt.Fprintln(cli.stdout, "  jira issue transition <issue-key> --to-id <transition-id>")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Examples:")
	fmt.Fprintln(cli.stdout, "  jira issue transition list PROJ-123")
	fmt.Fprintln(cli.stdout, "  jira issue transition PROJ-123 --to 'In Progress'")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Notes:")
	fmt.Fprintln(cli.stdout, "  Use 'jira issue update' for summary and description changes.")
}

func (cli CLI) printIssueUpdateHelp() {
	fmt.Fprintln(cli.stdout, "jira issue update")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Update editable Jira issue fields.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  jira issue update help")
	fmt.Fprintln(cli.stdout, "  jira issue update <issue-key> --summary <text>")
	fmt.Fprintln(cli.stdout, "  jira issue update <issue-key> --description <text>")
	fmt.Fprintln(cli.stdout, "  jira issue update <issue-key> --description-file <path>")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Examples:")
	fmt.Fprintln(cli.stdout, "  jira issue update PROJ-123 --summary 'Tighten validation'")
	fmt.Fprintln(cli.stdout, "  jira issue update PROJ-123 --description-file ./description.txt")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Notes:")
	fmt.Fprintln(cli.stdout, "  Status changes use 'jira issue transition'.")
}

func (cli CLI) printIssueEditMetaHelp() {
	fmt.Fprintln(cli.stdout, "jira issue editmeta")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Show which Jira issue fields are editable.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  jira issue editmeta help")
	fmt.Fprintln(cli.stdout, "  jira issue editmeta <issue-key>")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Examples:")
	fmt.Fprintln(cli.stdout, "  jira issue editmeta PROJ-123")
}

func (cli CLI) printHelp() {
	fmt.Fprintln(cli.stdout, "jira")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "CLI for Jira Cloud.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  jira help")
	fmt.Fprintln(cli.stdout, "  jira version")
	fmt.Fprintln(cli.stdout, "  jira auth help")
	fmt.Fprintln(cli.stdout, "  jira project help")
	fmt.Fprintln(cli.stdout, "  jira issue help")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Commands:")
	fmt.Fprintln(cli.stdout, "  help")
	fmt.Fprintln(cli.stdout, "  version")
	fmt.Fprintln(cli.stdout, "  auth")
	fmt.Fprintln(cli.stdout, "  project")
	fmt.Fprintln(cli.stdout, "  issue")
}

func main() {
	cli := CLI{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		env:    loadEnvironment(),
	}
	if err := cli.run(os.Args); err != nil {
		fmt.Fprintln(cli.stderr, err)
		os.Exit(1)
	}
}
