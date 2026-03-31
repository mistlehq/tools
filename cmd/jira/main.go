package main

import (
	"fmt"
	"io"
	"os"
)

// Version is the current jira CLI version.
var Version = "dev"

type CLI struct {
	stdout io.Writer
	stderr io.Writer
	env    Environment
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
	if len(args) == 0 {
		cli.printAuthHelp()
		return nil
	}

	config, err := loadConfig(cli.env)

	if err != nil {
		return err
	}

	jc := NewJiraClient(config)

	switch args[0] {
	case "whoami":
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
	fmt.Fprintln(cli.stdout, "Commands:")
	fmt.Fprintln(cli.stdout, "  whoami")
}

func (cli CLI) runProject(args []string) error {
	if len(args) == 0 {
		cli.printProjectHelp()
		return nil
	}

	config, err := loadConfig(cli.env)
	if err != nil {
		return err
	}

	jc := NewJiraClient(config)

	switch args[0] {
	case "list":
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
	fmt.Fprintln(cli.stdout, "Commands:")
	fmt.Fprintln(cli.stdout, "  list")
}

func (cli CLI) runIssue(args []string) error {
	if len(args) == 0 {
		cli.printIssueHelp()
		return nil
	}

	config, err := loadConfig(cli.env)
	if err != nil {
		return err
	}

	jc := NewJiraClient(config)

	switch args[0] {
	case "get":
		return cli.runIssueGet(jc, args[1:])
	case "search":
		return cli.runIssueSearch(jc, args[1:])
	default:
		return fmt.Errorf("unsupported issue command: %s", args[0])
	}
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
	fmt.Fprintln(cli.stdout, "Commands:")
	fmt.Fprintln(cli.stdout, "  get")
	fmt.Fprintln(cli.stdout, "  search")
}

func (cli CLI) printHelp() {
	fmt.Fprintln(cli.stdout, "jira")
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
		stdout: os.Stdout,
		stderr: os.Stderr,
		env:    loadEnvironment(),
	}
	if err := cli.run(os.Args); err != nil {
		fmt.Fprintln(cli.stderr, err)
		os.Exit(1)
	}
}
