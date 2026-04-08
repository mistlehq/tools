package main

import (
	"encoding/json"
	"fmt"
	"github.com/mistlehq/tools/internal/argparse"
	"github.com/mistlehq/tools/internal/textinput"
	"io"
	"os"
	"sort"
	"strings"
)

// Version is the current slack CLI version.
var Version = "dev"

type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	env    Environment
}

func (cli CLI) slackClient() (SlackClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return SlackClient{}, err
	}

	return NewSlackClient(config), nil
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
	case "conversations":
		return cli.runConversations(args[2:])
	case "chat":
		return cli.runChat(args[2:])
	case "reactions":
		return cli.runReactions(args[2:])
	case "files":
		return cli.runFiles(args[2:])
	case "emoji":
		return cli.runEmoji(args[2:])
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
	case "test":
		sc, err := cli.slackClient()
		if err != nil {
			return err
		}

		return cli.runAuthTest(sc, args[1:])
	default:
		return fmt.Errorf("unsupported auth command: %s", args[0])
	}
}

func (cli CLI) runConversations(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printConversationsHelp()
		return nil
	}

	sc, err := cli.slackClient()
	if err != nil {
		return err
	}

	switch args[0] {
	case "list":
		return cli.runConversationsList(sc, args[1:])
	case "info":
		return cli.runConversationsInfo(sc, args[1:])
	case "history":
		return cli.runConversationsHistory(sc, args[1:])
	default:
		return fmt.Errorf("unsupported conversations command: %s", args[0])
	}
}

func (cli CLI) runChat(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printChatHelp()
		return nil
	}

	sc, err := cli.slackClient()
	if err != nil {
		return err
	}

	switch args[0] {
	case "post-message":
		return cli.runChatPostMessage(sc, args[1:])
	case "update":
		return cli.runChatUpdate(sc, args[1:])
	case "delete":
		return cli.runChatDelete(sc, args[1:])
	case "get-permalink":
		return cli.runChatGetPermalink(sc, args[1:])
	default:
		return fmt.Errorf("unsupported chat command: %s", args[0])
	}
}

func (cli CLI) runReactions(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printReactionsHelp()
		return nil
	}

	sc, err := cli.slackClient()
	if err != nil {
		return err
	}

	switch args[0] {
	case "add":
		return cli.runReactionsAdd(sc, args[1:])
	case "remove":
		return cli.runReactionsRemove(sc, args[1:])
	default:
		return fmt.Errorf("unsupported reactions command: %s", args[0])
	}
}

func (cli CLI) runFiles(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printFilesHelp()
		return nil
	}

	sc, err := cli.slackClient()
	if err != nil {
		return err
	}

	switch args[0] {
	case "upload":
		return cli.runFilesUpload(sc, args[1:])
	default:
		return fmt.Errorf("unsupported files command: %s", args[0])
	}
}

func (cli CLI) runEmoji(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printEmojiHelp()
		return nil
	}

	sc, err := cli.slackClient()
	if err != nil {
		return err
	}

	switch args[0] {
	case "list":
		return cli.runEmojiList(sc, args[1:])
	default:
		return fmt.Errorf("unsupported emoji command: %s", args[0])
	}
}

func (cli CLI) printHelp() {
	fmt.Fprintln(cli.stdout, "slack")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "CLI for Slack.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  slack help")
	fmt.Fprintln(cli.stdout, "  slack version")
	fmt.Fprintln(cli.stdout, "  slack auth help")
	fmt.Fprintln(cli.stdout, "  slack conversations help")
	fmt.Fprintln(cli.stdout, "  slack chat help")
	fmt.Fprintln(cli.stdout, "  slack reactions help")
	fmt.Fprintln(cli.stdout, "  slack files help")
	fmt.Fprintln(cli.stdout, "  slack emoji help")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Commands:")
	fmt.Fprintln(cli.stdout, "  help")
	fmt.Fprintln(cli.stdout, "  version")
	fmt.Fprintln(cli.stdout, "  auth")
	fmt.Fprintln(cli.stdout, "  conversations")
	fmt.Fprintln(cli.stdout, "  chat")
	fmt.Fprintln(cli.stdout, "  reactions")
	fmt.Fprintln(cli.stdout, "  files")
	fmt.Fprintln(cli.stdout, "  emoji")
}

func (cli CLI) printAuthHelp() {
	fmt.Fprintln(cli.stdout, "slack auth")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Inspect Slack authentication state.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  slack auth help")
	fmt.Fprintln(cli.stdout, "  slack auth test")
}

func (cli CLI) runAuthTest(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"json": {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("auth test does not accept positional arguments")
	}

	authTest, err := sc.AuthTest()
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, authTest)
	}

	fmt.Fprintln(cli.stdout, "URL: "+authTest.URL)
	fmt.Fprintln(cli.stdout, "Team: "+authTest.Team)
	fmt.Fprintln(cli.stdout, "Team ID: "+authTest.TeamID)
	fmt.Fprintln(cli.stdout, "User: "+authTest.User)
	fmt.Fprintln(cli.stdout, "User ID: "+authTest.UserID)
	fmt.Fprintln(cli.stdout, "Bot ID: "+authTest.BotID)
	return nil
}

func (cli CLI) printConversationsHelp() {
	fmt.Fprintln(cli.stdout, "slack conversations")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Inspect Slack conversations.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  slack conversations help")
	fmt.Fprintln(cli.stdout, "  slack conversations list [--types <csv>] [--limit <n>] [--cursor <cursor>] [--exclude-archived]")
	fmt.Fprintln(cli.stdout, "  slack conversations info --channel <conversation-id> [--include-locale]")
	fmt.Fprintln(cli.stdout, "  slack conversations history --channel <conversation-id> [--cursor <cursor>] [--inclusive] [--latest <ts>] [--limit <n>] [--oldest <ts>]")
}

func (cli CLI) runConversationsList(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"types":            {TakesValue: true},
		"limit":            {TakesValue: true},
		"cursor":           {TakesValue: true},
		"exclude-archived": {},
		"json":             {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("conversations list does not accept positional arguments")
	}

	list, err := sc.ListConversations(SlackConversationsListInput{
		Types:           parsedArgs.First("types"),
		Limit:           parsedArgs.First("limit"),
		Cursor:          parsedArgs.First("cursor"),
		ExcludeArchived: parsedArgs.Has("exclude-archived"),
	})
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, list)
	}

	fmt.Fprintln(cli.stdout, "ID\tNAME\tIS_PRIVATE\tIS_ARCHIVED\tIS_MEMBER")
	for _, conversation := range list.Channels {
		fmt.Fprintf(cli.stdout, "%s\t%s\t%t\t%t\t%t\n", conversation.ID, conversation.Name, conversation.IsPrivate, conversation.IsArchived, conversation.IsMember)
	}

	if list.ResponseMetadata.NextCursor != "" {
		fmt.Fprintln(cli.stdout, "Next Cursor: "+list.ResponseMetadata.NextCursor)
		fmt.Fprintln(cli.stdout, "Next Page: "+buildConversationsListNextPageCommand(parsedArgs, list.ResponseMetadata.NextCursor))
	}

	return nil
}

func (cli CLI) runConversationsInfo(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"channel":        {TakesValue: true},
		"include-locale": {},
		"json":           {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("conversations info does not accept positional arguments")
	}

	channel := parsedArgs.First("channel")
	if channel == "" {
		return fmt.Errorf("conversations info requires --channel")
	}

	info, err := sc.GetConversationInfo(SlackConversationsInfoInput{
		Channel:       channel,
		IncludeLocale: parsedArgs.Has("include-locale"),
	})
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, info)
	}

	fmt.Fprintln(cli.stdout, "ID: "+info.Channel.ID)
	fmt.Fprintln(cli.stdout, "Name: "+info.Channel.Name)
	fmt.Fprintln(cli.stdout, "Is Private: "+formatBool(info.Channel.IsPrivate))
	fmt.Fprintln(cli.stdout, "Is Archived: "+formatBool(info.Channel.IsArchived))
	fmt.Fprintln(cli.stdout, "Is Member: "+formatBool(info.Channel.IsMember))
	if parsedArgs.Has("include-locale") {
		fmt.Fprintln(cli.stdout, "Locale: "+info.Channel.Locale)
	}
	return nil
}

func (cli CLI) runConversationsHistory(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"channel":   {TakesValue: true},
		"cursor":    {TakesValue: true},
		"inclusive": {},
		"latest":    {TakesValue: true},
		"limit":     {TakesValue: true},
		"oldest":    {TakesValue: true},
		"json":      {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("conversations history does not accept positional arguments")
	}

	channel := parsedArgs.First("channel")
	if channel == "" {
		return fmt.Errorf("conversations history requires --channel")
	}

	history, err := sc.GetConversationHistory(SlackConversationsHistoryInput{
		Channel:   channel,
		Cursor:    parsedArgs.First("cursor"),
		Inclusive: parsedArgs.Has("inclusive"),
		Latest:    parsedArgs.First("latest"),
		Limit:     parsedArgs.First("limit"),
		Oldest:    parsedArgs.First("oldest"),
	})
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, history)
	}

	for index, message := range history.Messages {
		fmt.Fprintln(cli.stdout, "TS: "+message.TS)
		fmt.Fprintln(cli.stdout, "Thread TS: "+message.ThreadTS)
		fmt.Fprintln(cli.stdout, "User: "+message.User)
		fmt.Fprintln(cli.stdout, "Type: "+messageType(message))
		fmt.Fprintln(cli.stdout, "Text:")
		fmt.Fprintln(cli.stdout, message.Text)

		if index < len(history.Messages)-1 {
			fmt.Fprintln(cli.stdout, "")
		}
	}

	if history.ResponseMetadata.NextCursor != "" {
		if len(history.Messages) > 0 {
			fmt.Fprintln(cli.stdout, "")
		}

		fmt.Fprintln(cli.stdout, "Next Cursor: "+history.ResponseMetadata.NextCursor)
		fmt.Fprintln(cli.stdout, "Next Page: "+buildConversationsHistoryNextPageCommand(parsedArgs, history.ResponseMetadata.NextCursor))
	}

	return nil
}

func (cli CLI) printChatHelp() {
	fmt.Fprintln(cli.stdout, "slack chat")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Post and update Slack messages.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  slack chat help")
	fmt.Fprintln(cli.stdout, "  slack chat post-message --channel <conversation-id> --text <text>")
	fmt.Fprintln(cli.stdout, "  slack chat post-message --channel <conversation-id> --text-file <path>")
	fmt.Fprintln(cli.stdout, "  slack chat post-message --channel <conversation-id> --thread-ts <ts> --text <text>")
	fmt.Fprintln(cli.stdout, "  slack chat update --channel <conversation-id> --ts <ts> --text <text>")
	fmt.Fprintln(cli.stdout, "  slack chat update --channel <conversation-id> --ts <ts> --text-file <path>")
	fmt.Fprintln(cli.stdout, "  slack chat delete --channel <conversation-id> --ts <ts>")
	fmt.Fprintln(cli.stdout, "  slack chat get-permalink --channel <conversation-id> --message-ts <ts>")
}

func (cli CLI) runChatPostMessage(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"channel":   {TakesValue: true},
		"text":      {TakesValue: true},
		"text-file": {TakesValue: true},
		"thread-ts": {TakesValue: true},
		"json":      {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("chat post-message does not accept positional arguments")
	}

	channel := parsedArgs.First("channel")
	if channel == "" {
		return fmt.Errorf("chat post-message requires --channel")
	}

	text, err := textinput.Read(cli.stdin, "text", parsedArgs.First("text"), "text-file", parsedArgs.First("text-file"))
	if err != nil {
		return err
	}

	input := SlackChatPostMessageInput{
		Channel: channel,
		Text:    text,
	}

	if threadTS := parsedArgs.First("thread-ts"); threadTS != "" {
		input.ThreadTS = &threadTS
	}

	posted, err := sc.PostMessage(input)
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, posted)
	}

	writeMessageResult(cli.stdout, posted.Channel, posted.TS, posted.Message.ThreadTS, posted.Message.Text)
	return nil
}

func (cli CLI) runChatUpdate(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"channel":   {TakesValue: true},
		"ts":        {TakesValue: true},
		"text":      {TakesValue: true},
		"text-file": {TakesValue: true},
		"json":      {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("chat update does not accept positional arguments")
	}

	channel := parsedArgs.First("channel")
	if channel == "" {
		return fmt.Errorf("chat update requires --channel")
	}

	ts := parsedArgs.First("ts")
	if ts == "" {
		return fmt.Errorf("chat update requires --ts")
	}

	text, err := textinput.Read(cli.stdin, "text", parsedArgs.First("text"), "text-file", parsedArgs.First("text-file"))
	if err != nil {
		return err
	}

	updated, err := sc.UpdateMessage(SlackChatUpdateInput{
		Channel: channel,
		TS:      ts,
		Text:    text,
	})
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, updated)
	}

	threadTS := updated.Message.ThreadTS
	if threadTS == "" {
		threadTS = updated.TS
	}

	writeMessageResult(cli.stdout, updated.Channel, updated.TS, threadTS, updated.Message.Text)
	return nil
}

func (cli CLI) runChatDelete(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"channel": {TakesValue: true},
		"ts":      {TakesValue: true},
		"json":    {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("chat delete does not accept positional arguments")
	}

	channel := parsedArgs.First("channel")
	if channel == "" {
		return fmt.Errorf("chat delete requires --channel")
	}

	ts := parsedArgs.First("ts")
	if ts == "" {
		return fmt.Errorf("chat delete requires --ts")
	}

	deleted, err := sc.DeleteMessage(SlackChatDeleteInput{
		Channel: channel,
		TS:      ts,
	})
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, deleted)
	}

	fmt.Fprintln(cli.stdout, "Channel: "+deleted.Channel)
	fmt.Fprintln(cli.stdout, "TS: "+deleted.TS)
	fmt.Fprintln(cli.stdout, "Deleted: true")
	return nil
}

func (cli CLI) runChatGetPermalink(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"channel":    {TakesValue: true},
		"message-ts": {TakesValue: true},
		"json":       {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("chat get-permalink does not accept positional arguments")
	}

	channel := parsedArgs.First("channel")
	if channel == "" {
		return fmt.Errorf("chat get-permalink requires --channel")
	}

	messageTS := parsedArgs.First("message-ts")
	if messageTS == "" {
		return fmt.Errorf("chat get-permalink requires --message-ts")
	}

	permalink, err := sc.GetPermalink(channel, messageTS)
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, permalink)
	}

	fmt.Fprintln(cli.stdout, "Channel: "+channel)
	fmt.Fprintln(cli.stdout, "Message TS: "+messageTS)
	fmt.Fprintln(cli.stdout, "Permalink: "+permalink.Permalink)
	return nil
}

func (cli CLI) printReactionsHelp() {
	fmt.Fprintln(cli.stdout, "slack reactions")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Add and remove Slack message reactions.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  slack reactions help")
	fmt.Fprintln(cli.stdout, "  slack reactions add --channel <conversation-id> --timestamp <ts> --name <emoji-name>")
	fmt.Fprintln(cli.stdout, "  slack reactions remove --channel <conversation-id> --timestamp <ts> --name <emoji-name>")
}

func (cli CLI) runReactionsAdd(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"channel":   {TakesValue: true},
		"timestamp": {TakesValue: true},
		"name":      {TakesValue: true},
		"json":      {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("reactions add does not accept positional arguments")
	}

	channel := parsedArgs.First("channel")
	if channel == "" {
		return fmt.Errorf("reactions add requires --channel")
	}

	timestamp := parsedArgs.First("timestamp")
	if timestamp == "" {
		return fmt.Errorf("reactions add requires --timestamp")
	}

	name := parsedArgs.First("name")
	if name == "" {
		return fmt.Errorf("reactions add requires --name")
	}

	response, err := sc.AddReaction(SlackReactionInput{
		Channel:   channel,
		Timestamp: timestamp,
		Name:      name,
	})
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, response)
	}

	fmt.Fprintln(cli.stdout, "Channel: "+channel)
	fmt.Fprintln(cli.stdout, "Timestamp: "+timestamp)
	fmt.Fprintln(cli.stdout, "Name: "+name)
	fmt.Fprintln(cli.stdout, "Action: added")
	return nil
}

func (cli CLI) runReactionsRemove(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"channel":   {TakesValue: true},
		"timestamp": {TakesValue: true},
		"name":      {TakesValue: true},
		"json":      {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("reactions remove does not accept positional arguments")
	}

	channel := parsedArgs.First("channel")
	if channel == "" {
		return fmt.Errorf("reactions remove requires --channel")
	}

	timestamp := parsedArgs.First("timestamp")
	if timestamp == "" {
		return fmt.Errorf("reactions remove requires --timestamp")
	}

	name := parsedArgs.First("name")
	if name == "" {
		return fmt.Errorf("reactions remove requires --name")
	}

	response, err := sc.RemoveReaction(SlackReactionInput{
		Channel:   channel,
		Timestamp: timestamp,
		Name:      name,
	})
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, response)
	}

	fmt.Fprintln(cli.stdout, "Channel: "+channel)
	fmt.Fprintln(cli.stdout, "Timestamp: "+timestamp)
	fmt.Fprintln(cli.stdout, "Name: "+name)
	fmt.Fprintln(cli.stdout, "Action: removed")
	return nil
}

func (cli CLI) printFilesHelp() {
	fmt.Fprintln(cli.stdout, "slack files")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Upload files to Slack.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  slack files help")
	fmt.Fprintln(cli.stdout, "  slack files upload --path <path> --channel <conversation-id>")
	fmt.Fprintln(cli.stdout, "  slack files upload --path <path> --channel <conversation-id> --initial-comment <text>")
	fmt.Fprintln(cli.stdout, "  slack files upload --path <path> --channel <conversation-id> --initial-comment-file <path>")
	fmt.Fprintln(cli.stdout, "  slack files upload --path <path> --channel <conversation-id> --thread-ts <ts>")
}

func (cli CLI) runFilesUpload(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"path":                 {TakesValue: true},
		"channel":              {TakesValue: true},
		"thread-ts":            {TakesValue: true},
		"initial-comment":      {TakesValue: true},
		"initial-comment-file": {TakesValue: true},
		"json":                 {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("files upload does not accept positional arguments")
	}

	path := parsedArgs.First("path")
	if path == "" {
		return fmt.Errorf("files upload requires --path")
	}

	channel := parsedArgs.First("channel")
	if channel == "" {
		return fmt.Errorf("files upload requires --channel")
	}

	initialComment, err := optionalTextInput(cli.stdin, "initial-comment", parsedArgs.First("initial-comment"), "initial-comment-file", parsedArgs.First("initial-comment-file"))
	if err != nil {
		return err
	}

	uploaded, err := sc.UploadFile(SlackFilesUploadInput{
		Path:           path,
		Channel:        channel,
		ThreadTS:       parsedArgs.First("thread-ts"),
		InitialComment: initialComment,
	})
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, uploaded)
	}

	file := firstUploadedFile(uploaded)
	fmt.Fprintln(cli.stdout, "Channel: "+channel)
	fmt.Fprintln(cli.stdout, "Thread TS: "+parsedArgs.First("thread-ts"))
	fmt.Fprintln(cli.stdout, "File ID: "+file.ID)
	fmt.Fprintln(cli.stdout, "Name: "+file.Name)
	fmt.Fprintln(cli.stdout, "Title: "+file.Title)
	return nil
}

func (cli CLI) printEmojiHelp() {
	fmt.Fprintln(cli.stdout, "slack emoji")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Inspect Slack emoji.")
	fmt.Fprintln(cli.stdout, "")
	fmt.Fprintln(cli.stdout, "Usage:")
	fmt.Fprintln(cli.stdout, "  slack emoji help")
	fmt.Fprintln(cli.stdout, "  slack emoji list")
	fmt.Fprintln(cli.stdout, "  slack emoji list --include-categories")
}

func (cli CLI) runEmojiList(sc SlackClient, args []string) error {
	parsedArgs, err := argparse.Parse(args, map[string]argparse.Spec{
		"include-categories": {},
		"json":               {},
	})
	if err != nil {
		return err
	}

	if len(parsedArgs.Positionals) > 0 {
		return fmt.Errorf("emoji list does not accept positional arguments")
	}

	list, err := sc.ListEmoji(parsedArgs.Has("include-categories"))
	if err != nil {
		return err
	}

	if parsedArgs.Has("json") {
		return writeJSON(cli.stdout, list)
	}

	names := make([]string, 0, len(list.Emoji))
	for name := range list.Emoji {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Fprintln(cli.stdout, "NAME\tVALUE")
	for _, name := range names {
		fmt.Fprintf(cli.stdout, "%s\t%s\n", name, list.Emoji[name])
	}

	return nil
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

func writeJSON(output io.Writer, value any) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(output, string(body))
	return err
}

func formatBool(value bool) string {
	if value {
		return "true"
	}

	return "false"
}

func buildConversationsListNextPageCommand(parsedArgs argparse.Parsed, nextCursor string) string {
	command := []string{"slack", "conversations", "list"}

	if types := parsedArgs.First("types"); types != "" {
		command = append(command, "--types", types)
	}

	if limit := parsedArgs.First("limit"); limit != "" {
		command = append(command, "--limit", limit)
	}

	if parsedArgs.Has("exclude-archived") {
		command = append(command, "--exclude-archived")
	}

	command = append(command, "--cursor", nextCursor)

	if parsedArgs.Has("json") {
		command = append(command, "--json")
	}

	return strings.Join(command, " ")
}

func buildConversationsHistoryNextPageCommand(parsedArgs argparse.Parsed, nextCursor string) string {
	command := []string{"slack", "conversations", "history"}

	command = append(command, "--channel", parsedArgs.First("channel"))

	if limit := parsedArgs.First("limit"); limit != "" {
		command = append(command, "--limit", limit)
	}

	if parsedArgs.Has("inclusive") {
		command = append(command, "--inclusive")
	}

	if latest := parsedArgs.First("latest"); latest != "" {
		command = append(command, "--latest", latest)
	}

	if oldest := parsedArgs.First("oldest"); oldest != "" {
		command = append(command, "--oldest", oldest)
	}

	command = append(command, "--cursor", nextCursor)

	if parsedArgs.Has("json") {
		command = append(command, "--json")
	}

	return strings.Join(command, " ")
}

func messageType(message SlackMessage) string {
	if message.Subtype != "" {
		return message.Subtype
	}

	return message.Type
}

func writeMessageResult(output io.Writer, channel string, ts string, threadTS string, text string) {
	fmt.Fprintln(output, "Channel: "+channel)
	fmt.Fprintln(output, "TS: "+ts)
	fmt.Fprintln(output, "Thread TS: "+threadTS)
	fmt.Fprintln(output, "Text:")
	fmt.Fprintln(output, text)
}

func optionalTextInput(stdin io.Reader, valueFlagName string, value string, fileFlagName string, filePath string) (string, error) {
	if value == "" && filePath == "" {
		return "", nil
	}

	return textinput.Read(stdin, valueFlagName, value, fileFlagName, filePath)
}

func firstUploadedFile(uploaded SlackFilesCompleteUploadExternal) SlackFile {
	if len(uploaded.Files) == 0 {
		return SlackFile{}
	}

	return uploaded.Files[0]
}
