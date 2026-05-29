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

var jiraIssueCreateDoc = commandDoc{
	Command:     "jira issue create",
	Summary:     "Create a Jira issue",
	Description: "Create a Jira issue.",
}

var jiraIssueDeleteDoc = commandDoc{
	Command:     "jira issue delete",
	Summary:     "Delete a Jira issue by key",
	Description: "Delete a Jira issue by key.",
}

var jiraIssueCommentAddDoc = commandDoc{
	Command:     "jira issue comment add",
	Summary:     "Add a comment to a Jira issue",
	Description: "Add a comment to a Jira issue.",
}

var jiraIssueCommentDeleteDoc = commandDoc{
	Command:     "jira issue comment delete",
	Summary:     "Delete a comment from a Jira issue",
	Description: "Delete a comment from a Jira issue.",
}

var jiraIssueAssignDoc = commandDoc{
	Command:     "jira issue assign",
	Summary:     "Assign or clear the assignee on a Jira issue",
	Description: "Assign or clear the assignee on a Jira issue.",
}

var jiraIssueTransitionListDoc = commandDoc{
	Command:     "jira issue transition list",
	Summary:     "List available issue workflow transitions",
	Description: "List the workflow transitions currently available for an issue.",
}

var jiraIssueTransitionDoc = commandDoc{
	Command:     "jira issue transition",
	Summary:     "Transition a Jira issue to a new workflow state",
	Description: "Transition a Jira issue to a new workflow state.",
}

var jiraIssueUpdateDoc = commandDoc{
	Command:     "jira issue update",
	Summary:     "Update editable Jira issue fields",
	Description: "Update summary, description, and other editable fields on a Jira issue.",
}

var jiraIssueEditMetaDoc = commandDoc{
	Command:     "jira issue editmeta",
	Summary:     "Show edit metadata for a Jira issue",
	Description: "Show edit metadata for a Jira issue.",
}
