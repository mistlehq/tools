package main

type commandDoc struct {
	Command     string
	Summary     string
	Description string
}

var discordAuthTestDoc = commandDoc{
	Command:     "discord auth test",
	Summary:     "Check Discord bot authentication state",
	Description: "Check Discord bot authentication state by fetching the current bot user.",
}

var discordGuildsListDoc = commandDoc{
	Command:     "discord guilds list",
	Summary:     "List Discord guilds visible to the bot",
	Description: "List Discord guilds visible to the bot.",
}

var discordGuildsGetDoc = commandDoc{
	Command:     "discord guilds get",
	Summary:     "Show details for a Discord guild",
	Description: "Show details for a Discord guild.",
}

var discordChannelsListDoc = commandDoc{
	Command:     "discord channels list",
	Summary:     "List Discord guild channels",
	Description: "List Discord channels in a guild.",
}

var discordChannelsGetDoc = commandDoc{
	Command:     "discord channels get",
	Summary:     "Show details for a Discord channel",
	Description: "Show details for a Discord channel.",
}

var discordMessagesListDoc = commandDoc{
	Command:     "discord messages list",
	Summary:     "Fetch Discord channel messages",
	Description: "Fetch recent Discord messages from a channel.",
}

var discordMessagesSendDoc = commandDoc{
	Command:     "discord messages send",
	Summary:     "Send a Discord message",
	Description: "Send a Discord message to a channel.",
}

var discordMessagesEditDoc = commandDoc{
	Command:     "discord messages edit",
	Summary:     "Edit a Discord message",
	Description: "Edit a Discord message.",
}

var discordMessagesDeleteDoc = commandDoc{
	Command:     "discord messages delete",
	Summary:     "Delete a Discord message",
	Description: "Delete a Discord message.",
}

var discordReactionsAddDoc = commandDoc{
	Command:     "discord reactions add",
	Summary:     "Add a Discord message reaction",
	Description: "Add a Discord reaction to a message.",
}

var discordReactionsRemoveDoc = commandDoc{
	Command:     "discord reactions remove",
	Summary:     "Remove the bot's Discord message reaction",
	Description: "Remove the bot's Discord reaction from a message.",
}

var discordRolesListDoc = commandDoc{
	Command:     "discord roles list",
	Summary:     "List Discord guild roles",
	Description: "List Discord roles in a guild.",
}

var discordRolesCreateDoc = commandDoc{
	Command:     "discord roles create",
	Summary:     "Create a Discord guild role",
	Description: "Create a Discord role in a guild.",
}

var discordRolesDeleteDoc = commandDoc{
	Command:     "discord roles delete",
	Summary:     "Delete a Discord guild role",
	Description: "Delete a Discord role in a guild.",
}

var discordMembersListDoc = commandDoc{
	Command:     "discord members list",
	Summary:     "List Discord guild members",
	Description: "List Discord guild members visible to the bot.",
}

var discordMembersGetDoc = commandDoc{
	Command:     "discord members get",
	Summary:     "Show details for a Discord guild member",
	Description: "Show details for a Discord guild member.",
}

var discordMembersAddRoleDoc = commandDoc{
	Command:     "discord members add-role",
	Summary:     "Assign a Discord role to a member",
	Description: "Assign a Discord guild role to a member.",
}

var discordMembersRemoveRoleDoc = commandDoc{
	Command:     "discord members remove-role",
	Summary:     "Remove a Discord role from a member",
	Description: "Remove a Discord guild role from a member.",
}

var discordMembersBanDoc = commandDoc{
	Command:     "discord members ban",
	Summary:     "Ban a Discord guild member",
	Description: "Ban a Discord guild member.",
}

var discordMembersUnbanDoc = commandDoc{
	Command:     "discord members unban",
	Summary:     "Unban a Discord user from a guild",
	Description: "Unban a Discord user from a guild.",
}
