package main

import (
	"encoding/json"
	"fmt"
	"github.com/mistlehq/tools/internal/argparse"
	"io"
	"os"
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

	return fmt.Errorf("unsupported chat command: %s", args[0])
}

func (cli CLI) runReactions(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printReactionsHelp()
		return nil
	}

	return fmt.Errorf("unsupported reactions command: %s", args[0])
}

func (cli CLI) runFiles(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printFilesHelp()
		return nil
	}

	return fmt.Errorf("unsupported files command: %s", args[0])
}

func (cli CLI) runEmoji(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		cli.printEmojiHelp()
		return nil
	}

	return fmt.Errorf("unsupported emoji command: %s", args[0])
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
