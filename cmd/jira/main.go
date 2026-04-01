package main

import (
	"fmt"
	"io"
	"os"
	"sort"
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

func isHelpToken(arg string) bool {
	return arg == "help" || arg == "-h" || arg == "--help"
}

func isSingleHelpArg(args []string) bool {
	return len(args) == 1 && isHelpToken(args[0])
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
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printAuthHelp()
		return nil
	}

	switch args[0] {
	case "whoami":
		if isSingleHelpArg(args[1:]) {
			cli.printAuthWhoAmIHelp()
			return nil
		}

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
	if isSingleHelpArg(args) {
		cli.printAuthWhoAmIHelp()
		return nil
	}

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
	fmt.Fprint(cli.stdout, `jira auth

Inspect Jira authentication state.

Usage:
  jira auth help
  jira auth whoami
  jira auth whoami --help

Commands:
  whoami    Show the Jira account behind the current auth context
`)
}

func (cli CLI) printAuthWhoAmIHelp() {
	fmt.Fprint(cli.stdout, `jira auth whoami

Show the Jira account behind the current auth context.

Usage:
  jira auth whoami
  jira auth whoami --help

Output:
  Prints the account ID, display name, and email returned by Jira.
`)
}

func (cli CLI) runProject(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printProjectHelp()
		return nil
	}

	switch args[0] {
	case "list":
		if isSingleHelpArg(args[1:]) {
			cli.printProjectListHelp()
			return nil
		}

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
	if isSingleHelpArg(args) {
		cli.printProjectListHelp()
		return nil
	}

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
	fmt.Fprint(cli.stdout, `jira project

Inspect Jira projects visible to the current caller.

Usage:
  jira project help
  jira project list
  jira project list --help

Commands:
  list    List visible projects with their IDs, keys, and names
`)
}

func (cli CLI) printProjectListHelp() {
	fmt.Fprint(cli.stdout, `jira project list

List Jira projects visible to the current caller.

Usage:
  jira project list
  jira project list --help

Output:
  Prints a table with project ID, project key, and project name.
`)
}

func (cli CLI) runIssue(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printIssueHelp()
		return nil
	}

	switch args[0] {
	case "get":
		if isSingleHelpArg(args[1:]) {
			cli.printIssueGetHelp()
			return nil
		}

		jc, err := cli.jiraClient()
		if err != nil {
			return err
		}

		return cli.runIssueGet(jc, args[1:])
	case "search":
		if isSingleHelpArg(args[1:]) {
			cli.printIssueSearchHelp()
			return nil
		}

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
	if len(args) == 0 || isHelpToken(args[0]) {
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
	if len(args) == 0 {
		cli.printIssueAssignHelp()
		return nil
	}

	if isSingleHelpArg(args) {
		cli.printIssueAssignCommandHelp()
		return nil
	}

	return cli.runIssueAssignSet(args)
}

func (cli CLI) runIssueTransition(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
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
	if len(args) == 0 {
		cli.printIssueUpdateHelp()
		return nil
	}

	if isSingleHelpArg(args) {
		cli.printIssueUpdateFieldsHelp()
		return nil
	}

	return cli.runIssueUpdateFields(args)
}

func (cli CLI) runIssueEditMeta(args []string) error {
	if isSingleHelpArg(args) {
		cli.printIssueEditMetaCommandHelp()
		return nil
	}

	if len(args) != 1 {
		return fmt.Errorf("issue editmeta expects exactly 1 positional argument")
	}

	jc, err := cli.jiraClient()
	if err != nil {
		return err
	}

	editMeta, err := jc.GetIssueEditMeta(args[0])
	if err != nil {
		return err
	}

	fieldIDs := make([]string, 0, len(editMeta.Fields))
	for fieldID := range editMeta.Fields {
		fieldIDs = append(fieldIDs, fieldID)
	}
	sort.Strings(fieldIDs)

	fmt.Fprintln(cli.stdout, "FIELD ID\tNAME\tREQUIRED\tTYPE")
	for _, fieldID := range fieldIDs {
		field := editMeta.Fields[fieldID]
		fmt.Fprintf(cli.stdout, "%s\t%s\t%t\t%s\n", fieldID, field.Name, field.Required, field.Schema.Type)
	}

	return nil
}

func (cli CLI) runIssueAssignSet(args []string) error {
	if isSingleHelpArg(args) {
		cli.printIssueAssignCommandHelp()
		return nil
	}

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
	if isSingleHelpArg(args) {
		cli.printIssueTransitionListHelp()
		return nil
	}

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
	if isSingleHelpArg(args) {
		cli.printIssueTransitionMoveHelp()
		return nil
	}

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

func (cli CLI) runIssueUpdateFields(args []string) error {
	if isSingleHelpArg(args) {
		cli.printIssueUpdateFieldsHelp()
		return nil
	}

	parsedArgs, err := parseArgs(args, map[string]argSpec{
		"summary": {
			takesValue: true,
		},
		"description": {
			takesValue: true,
		},
		"description-file": {
			takesValue: true,
		},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.positionals) != 1 {
		return fmt.Errorf("issue update expects exactly 1 positional argument")
	}

	summary := parsedArgs.first("summary")
	descriptionFlag := parsedArgs.first("description")
	descriptionFileFlag := parsedArgs.first("description-file")

	updatedFields := make([]string, 0, 2)
	input := UpdateIssueInput{}

	if summary != "" {
		input.Summary = &summary
		updatedFields = append(updatedFields, "summary")
	}

	if descriptionFlag != "" || descriptionFileFlag != "" {
		description, err := cli.readTextInput("description", descriptionFlag, "description-file", descriptionFileFlag)
		if err != nil {
			return err
		}

		input.Description = &description
		updatedFields = append(updatedFields, "description")
	}

	if len(updatedFields) == 0 {
		return fmt.Errorf("issue update requires at least one of --summary, --description, or --description-file")
	}

	jc, err := cli.jiraClient()
	if err != nil {
		return err
	}

	issueKey := parsedArgs.positionals[0]
	if err := jc.UpdateIssue(issueKey, input); err != nil {
		return err
	}

	fmt.Fprintln(cli.stdout, "Issue: "+issueKey)
	fmt.Fprintln(cli.stdout, "Updated: "+strings.Join(updatedFields, ", "))
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
	if isSingleHelpArg(args) {
		cli.printIssueCommentAddHelp()
		return nil
	}

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
	if isSingleHelpArg(args) {
		cli.printIssueGetHelp()
		return nil
	}

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
	if isSingleHelpArg(args) {
		cli.printIssueSearchHelp()
		return nil
	}

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
	fmt.Fprint(cli.stdout, `jira issue

Inspect and mutate Jira issues.

Usage:
  jira issue help
  jira issue get <issue-key>
  jira issue get --help
  jira issue search '<jql query>'
  jira issue search --help
  jira issue comment help
  jira issue assign help
  jira issue transition help
  jira issue update help
  jira issue editmeta help

Commands:
  get         Fetch a single issue
  search      Search issues with JQL
  comment     Add comments to issues
  assign      Change assignees
  transition  List or apply workflow transitions
  update      Edit summary and description fields
  editmeta    Show which fields are editable on an issue

Notes:
  Status changes go through 'jira issue transition', not 'jira issue update'.
  Leaf commands also accept --help for command-specific usage.
`)
}

func (cli CLI) printIssueGetHelp() {
	fmt.Fprint(cli.stdout, `jira issue get

Fetch a single Jira issue by key.

Usage:
  jira issue get <issue-key>
  jira issue get --help

Output:
  Prints the issue ID, key, summary, and current status.

Example:
  jira issue get PROJ-123
`)
}

func (cli CLI) printIssueSearchHelp() {
	fmt.Fprint(cli.stdout, `jira issue search

Search Jira issues with a JQL query.

Usage:
  jira issue search '<jql query>'
  jira issue search --help

Output:
  Prints a table with issue key, status, and summary.

Example:
  jira issue search 'project = PROJ ORDER BY updated DESC'
`)
}

func (cli CLI) printIssueCommentHelp() {
	fmt.Fprint(cli.stdout, `jira issue comment

Add comments to Jira issues.

Usage:
  jira issue comment help
  jira issue comment add <issue-key> --body <text>
  jira issue comment add <issue-key> --body-file <path>
  jira issue comment add --help

Commands:
  add    Create a new issue comment from inline text, a file, or stdin
`)
}

func (cli CLI) printIssueCommentAddHelp() {
	fmt.Fprint(cli.stdout, `jira issue comment add

Add a comment to a Jira issue.

Usage:
  jira issue comment add <issue-key> --body <text>
  jira issue comment add <issue-key> --body-file <path>
  jira issue comment add --help

Examples:
  jira issue comment add PROJ-123 --body 'Looks good'
  jira issue comment add PROJ-123 --body-file ./comment.txt
  jira issue comment add PROJ-123 --body-file -

Notes:
  Exactly one of --body or --body-file is required.
  Use --body-file - to read the comment body from stdin.
`)
}

func (cli CLI) printIssueAssignHelp() {
	fmt.Fprint(cli.stdout, `jira issue assign

Assign Jira issues.

Usage:
  jira issue assign help
  jira issue assign <issue-key> --me
  jira issue assign <issue-key> --account-id <account-id>
  jira issue assign <issue-key> --unassigned
  jira issue assign --help

Notes:
  Exactly one of --me, --account-id, or --unassigned is required.
`)
}

func (cli CLI) printIssueAssignCommandHelp() {
	fmt.Fprint(cli.stdout, `jira issue assign

Assign or clear the assignee on a Jira issue.

Usage:
  jira issue assign <issue-key> --me
  jira issue assign <issue-key> --account-id <account-id>
  jira issue assign <issue-key> --unassigned
  jira issue assign --help

Examples:
  jira issue assign PROJ-123 --me
  jira issue assign PROJ-123 --account-id 712020:abc123
  jira issue assign PROJ-123 --unassigned

Notes:
  Exactly one of --me, --account-id, or --unassigned is required.
`)
}

func (cli CLI) printIssueTransitionHelp() {
	fmt.Fprint(cli.stdout, `jira issue transition

Move Jira issues through workflow transitions.

Usage:
  jira issue transition help
  jira issue transition list <issue-key>
  jira issue transition list --help
  jira issue transition <issue-key> --to <transition-name>
  jira issue transition <issue-key> --to-id <transition-id>
  jira issue transition --help

Notes:
  Use 'jira issue update' for summary and description changes.
`)
}

func (cli CLI) printIssueTransitionListHelp() {
	fmt.Fprint(cli.stdout, `jira issue transition list

List the workflow transitions currently available for an issue.

Usage:
  jira issue transition list <issue-key>
  jira issue transition list --help

Output:
  Prints transition ID, transition name, and destination status.

Example:
  jira issue transition list PROJ-123
`)
}

func (cli CLI) printIssueTransitionMoveHelp() {
	fmt.Fprint(cli.stdout, `jira issue transition

Transition a Jira issue to a new workflow state.

Usage:
  jira issue transition <issue-key> --to <transition-name>
  jira issue transition <issue-key> --to-id <transition-id>
  jira issue transition --help

Examples:
  jira issue transition PROJ-123 --to 'In Progress'
  jira issue transition PROJ-123 --to-id 31

Notes:
  Exactly one of --to or --to-id is required.
  Use 'jira issue transition list <issue-key>' to inspect valid transitions first.
`)
}

func (cli CLI) printIssueUpdateHelp() {
	fmt.Fprint(cli.stdout, `jira issue update

Update editable Jira issue fields.

Usage:
  jira issue update help
  jira issue update <issue-key> --summary <text>
  jira issue update <issue-key> --description <text>
  jira issue update <issue-key> --description-file <path>
  jira issue update --help

Notes:
  Status changes use 'jira issue transition'.
`)
}

func (cli CLI) printIssueUpdateFieldsHelp() {
	fmt.Fprint(cli.stdout, `jira issue update

Update summary and description fields on a Jira issue.

Usage:
  jira issue update <issue-key> --summary <text>
  jira issue update <issue-key> --description <text>
  jira issue update <issue-key> --description-file <path>
  jira issue update --help

Examples:
  jira issue update PROJ-123 --summary 'Tighten validation'
  jira issue update PROJ-123 --description 'Expanded implementation notes'
  jira issue update PROJ-123 --description-file ./description.txt
  jira issue update PROJ-123 --description-file -

Notes:
  Provide at least one of --summary, --description, or --description-file.
  Use --description-file - to read the description from stdin.
  Status changes use 'jira issue transition'.
`)
}

func (cli CLI) printIssueEditMetaCommandHelp() {
	fmt.Fprint(cli.stdout, `jira issue editmeta

Show edit metadata for a Jira issue.

Usage:
  jira issue editmeta <issue-key>
  jira issue editmeta --help

Output:
  Prints editable field IDs, field names, whether they are required, and field types.

Example:
  jira issue editmeta PROJ-123
`)
}

func (cli CLI) printHelp() {
	fmt.Fprint(cli.stdout, `jira

Thin Jira Cloud CLI for shells, scripts, and agent-driven workflows.

The CLI assumes auth is already handled upstream and stays focused on direct Jira operations.

Usage:
  jira help
  jira --help
  jira version
  jira auth help
  jira project help
  jira issue help

Command Families:
  auth       Inspect the current Jira auth context
  project    Discover visible Jira projects
  issue      Read issues and perform common issue mutations

Common Starting Points:
  jira auth whoami
  jira project list
  jira issue search 'project = PROJ ORDER BY updated DESC'
  jira issue get PROJ-123

Issue Workflows:
  jira issue comment add PROJ-123 --body 'Looks good'
  jira issue assign PROJ-123 --me
  jira issue transition list PROJ-123
  jira issue transition PROJ-123 --to 'In Progress'
  jira issue update PROJ-123 --summary 'Tighten validation'
  jira issue editmeta PROJ-123

Dive Deeper:
  jira auth help
  jira project help
  jira issue help
  jira issue comment help
  jira issue assign help
  jira issue transition help
  jira issue update help
  jira issue editmeta help

Help Conventions:
  Namespaces accept help, -h, and --help.
  Leaf commands also accept --help, for example:
    jira issue get --help
    jira issue search --help
    jira issue transition list --help
`)
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
