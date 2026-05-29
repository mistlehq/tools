package main

type commandDoc struct {
	Command     string
	Summary     string
	Description string
}

var jiraAuthWhoAmIDoc = commandDoc{
	Command:     "jira auth whoami",
	Summary:     "Show the Jira account behind the current auth context",
	Description: "Show the Jira account behind the current auth context.",
}

var jiraProjectListDoc = commandDoc{
	Command:     "jira project list",
	Summary:     "List visible projects with their IDs, keys, and names",
	Description: "List Jira projects visible to the current caller.",
}

var jiraIssueGetDoc = commandDoc{
	Command:     "jira issue get",
	Summary:     "Fetch a single issue",
	Description: "Fetch a single Jira issue by key.",
}

var jiraIssueSearchDoc = commandDoc{
	Command:     "jira issue search",
	Summary:     "Search issues with JQL",
	Description: "Search Jira issues with a JQL query.",
}
