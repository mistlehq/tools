package main

type commandDoc struct {
	Command     string
	Summary     string
	Description string
}

var slackAuthTestDoc = commandDoc{
	Command:     "slack auth test",
	Summary:     "Check Slack authentication state",
	Description: "Check Slack authentication state.",
}

var slackConversationsListDoc = commandDoc{
	Command:     "slack conversations list",
	Summary:     "List Slack conversations",
	Description: "List Slack conversations.",
}

var slackConversationsInfoDoc = commandDoc{
	Command:     "slack conversations info",
	Summary:     "Show details for a Slack conversation",
	Description: "Show details for a Slack conversation.",
}

var slackConversationsHistoryDoc = commandDoc{
	Command:     "slack conversations history",
	Summary:     "Fetch Slack conversation history",
	Description: "Fetch Slack conversation history.",
}

var slackConversationsRepliesDoc = commandDoc{
	Command:     "slack conversations replies",
	Summary:     "Fetch replies in a Slack thread",
	Description: "Fetch replies in a Slack thread.",
}

var slackChatPostMessageDoc = commandDoc{
	Command:     "slack chat post-message",
	Summary:     "Post a Slack message",
	Description: "Post a Slack message.",
}

var slackChatUpdateDoc = commandDoc{
	Command:     "slack chat update",
	Summary:     "Update a Slack message",
	Description: "Update a Slack message.",
}

var slackChatDeleteDoc = commandDoc{
	Command:     "slack chat delete",
	Summary:     "Delete a Slack message",
	Description: "Delete a Slack message.",
}

var slackChatGetPermalinkDoc = commandDoc{
	Command:     "slack chat get-permalink",
	Summary:     "Get a permalink for a Slack message",
	Description: "Get a permalink for a Slack message.",
}

var slackReactionsAddDoc = commandDoc{
	Command:     "slack reactions add",
	Summary:     "Add a Slack message reaction",
	Description: "Add a Slack message reaction.",
}

var slackReactionsRemoveDoc = commandDoc{
	Command:     "slack reactions remove",
	Summary:     "Remove a Slack message reaction",
	Description: "Remove a Slack message reaction.",
}

var slackFilesInfoDoc = commandDoc{
	Command:     "slack files info",
	Summary:     "Show Slack file metadata",
	Description: "Show Slack file metadata.",
}

var slackFilesDownloadDoc = commandDoc{
	Command:     "slack files download",
	Summary:     "Download a Slack file to a local path",
	Description: "Download a Slack file to a local path.",
}

var slackFilesUploadDoc = commandDoc{
	Command:     "slack files upload",
	Summary:     "Upload a local file to Slack",
	Description: "Upload a local file to Slack.",
}

var slackEmojiListDoc = commandDoc{
	Command:     "slack emoji list",
	Summary:     "List Slack emoji",
	Description: "List Slack emoji.",
}
