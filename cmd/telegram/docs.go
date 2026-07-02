package main

type commandDoc struct {
	Command     string
	Summary     string
	Description string
}

var telegramAuthTestDoc = commandDoc{
	Command:     "telegram auth test",
	Summary:     "Check Telegram bot authentication state",
	Description: "Check Telegram bot authentication state by calling getMe.",
}

var telegramChatsGetDoc = commandDoc{
	Command:     "telegram chats get",
	Summary:     "Show details for a Telegram chat",
	Description: "Show details for a Telegram chat by ID or username.",
}

var telegramMessagesSendDoc = commandDoc{
	Command:     "telegram messages send",
	Summary:     "Send a Telegram text message",
	Description: "Send a Telegram text message to a chat.",
}

var telegramMessagesEditDoc = commandDoc{
	Command:     "telegram messages edit",
	Summary:     "Edit a Telegram text message",
	Description: "Edit a Telegram text message in a chat.",
}

var telegramMessagesDeleteDoc = commandDoc{
	Command:     "telegram messages delete",
	Summary:     "Delete a Telegram message",
	Description: "Delete a Telegram message from a chat.",
}

var telegramMessagesDeleteBatchDoc = commandDoc{
	Command:     "telegram messages delete-batch",
	Summary:     "Delete multiple Telegram messages",
	Description: "Delete multiple Telegram messages from a chat.",
}

var telegramReactionsSetDoc = commandDoc{
	Command:     "telegram reactions set",
	Summary:     "Set Telegram message reactions",
	Description: "Set Telegram reactions on a message.",
}

var telegramReactionsClearDoc = commandDoc{
	Command:     "telegram reactions clear",
	Summary:     "Clear Telegram message reactions",
	Description: "Clear Telegram reactions from a message.",
}

var telegramReactionsDeleteDoc = commandDoc{
	Command:     "telegram reactions delete",
	Summary:     "Delete a Telegram message reaction",
	Description: "Delete a reaction from a Telegram message.",
}

var telegramReactionsDeleteAllDoc = commandDoc{
	Command:     "telegram reactions delete-all",
	Summary:     "Delete Telegram message reactions in bulk",
	Description: "Delete recent Telegram reactions by user or actor chat.",
}

var telegramRequestDoc = commandDoc{
	Command:     "telegram request",
	Summary:     "Call a Telegram Bot API method",
	Description: "Call an arbitrary Telegram Bot API method with a JSON body.",
}

var telegramTopicsCreateDoc = commandDoc{
	Command:     "telegram topics create",
	Summary:     "Create a Telegram forum topic",
	Description: "Create a topic in a Telegram forum supergroup.",
}

var telegramTopicsDeleteDoc = commandDoc{
	Command:     "telegram topics delete",
	Summary:     "Delete a Telegram forum topic",
	Description: "Delete a Telegram forum topic and its messages.",
}
