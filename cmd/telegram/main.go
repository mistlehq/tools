package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mistlehq/tools/internal/argparse"
)

// Version is the current telegram CLI version.
var Version = "dev"

type CLI struct {
	stdout io.Writer
	stderr io.Writer
	env    Environment
}

func (cli CLI) telegramClient() (TelegramClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return TelegramClient{}, err
	}

	return NewTelegramClient(config), nil
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
	case "chats":
		return cli.runChats(args[2:])
	case "messages":
		return cli.runMessages(args[2:])
	case "reactions":
		return cli.runReactions(args[2:])
	case "topics":
		return cli.runTopics(args[2:])
	case "request":
		return cli.runRequest(args[2:])
	case "mcp":
		return cli.runMCP(args[2:])
	default:
		return fmt.Errorf("unsupported command: %s", args[1])
	}
}

func (cli CLI) runAuth(args []string) error {
	if isHelp(args) {
		cli.printAuthHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printAuthHelp()
		return nil
	}
	if args[0] != "test" {
		return fmt.Errorf("unsupported auth command: %s", args[0])
	}
	if _, err := parseOutputArgs(args[1:]); err != nil {
		return err
	}
	tc, err := cli.telegramClient()
	if err != nil {
		return err
	}
	out, err := tc.AuthTestContext(cliContext())
	return writeOutput(cli, out, err, args[1:], writeTelegramUser)
}

func (cli CLI) runChats(args []string) error {
	if isHelp(args) {
		cli.printChatsHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printChatsHelp()
		return nil
	}

	switch args[0] {
	case "get":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"chat": {TakesValue: true}}, "chat")
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.GetChatContext(cliContext(), parsedArgs.First("chat"))
		return writeOutput(cli, out, err, args[1:], writeTelegramChat)
	default:
		return fmt.Errorf("unsupported chats command: %s", args[0])
	}
}

func (cli CLI) runMessages(args []string) error {
	if isHelp(args) {
		cli.printMessagesHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printMessagesHelp()
		return nil
	}

	switch args[0] {
	case "send":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":       {TakesValue: true},
			"text":       {TakesValue: true},
			"parse-mode": {TakesValue: true},
			"thread":     {TakesValue: true},
		}, "chat", "text")
		if err != nil {
			return err
		}
		threadID, err := parseOptionalPositiveIntFlag("thread", parsedArgs.First("thread"))
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.SendMessageContext(cliContext(), TelegramSendMessageInput{
			Chat:      parsedArgs.First("chat"),
			Text:      parsedArgs.First("text"),
			ParseMode: parsedArgs.First("parse-mode"),
			ThreadID:  threadID,
		})
		return writeOutput(cli, out, err, args[1:], writeTelegramMessage)
	case "edit":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":       {TakesValue: true},
			"message":    {TakesValue: true},
			"text":       {TakesValue: true},
			"parse-mode": {TakesValue: true},
		}, "chat", "message", "text")
		if err != nil {
			return err
		}
		messageID, err := parsePositiveIntFlag("message", parsedArgs.First("message"))
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.EditMessageTextContext(cliContext(), TelegramEditMessageInput{
			Chat:      parsedArgs.First("chat"),
			MessageID: messageID,
			Text:      parsedArgs.First("text"),
			ParseMode: parsedArgs.First("parse-mode"),
		})
		return writeOutput(cli, out, err, args[1:], writeTelegramMessage)
	case "delete":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":    {TakesValue: true},
			"message": {TakesValue: true},
		}, "chat", "message")
		if err != nil {
			return err
		}
		messageID, err := parsePositiveIntFlag("message", parsedArgs.First("message"))
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.DeleteMessageContext(cliContext(), TelegramMessageInput{Chat: parsedArgs.First("chat"), MessageID: messageID})
		return writeOutput(cli, out, err, args[1:], writeTelegramBool)
	case "delete-batch":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":     {TakesValue: true},
			"messages": {TakesValue: true},
		}, "chat", "messages")
		if err != nil {
			return err
		}
		messageIDs, err := parsePositiveIntCSVFlag("messages", parsedArgs.First("messages"))
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.DeleteMessagesContext(cliContext(), TelegramDeleteMessagesInput{Chat: parsedArgs.First("chat"), MessageIDs: messageIDs})
		return writeOutput(cli, out, err, args[1:], writeTelegramBool)
	default:
		return fmt.Errorf("unsupported messages command: %s", args[0])
	}
}

func (cli CLI) runReactions(args []string) error {
	if isHelp(args) {
		cli.printReactionsHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printReactionsHelp()
		return nil
	}

	switch args[0] {
	case "set":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":            {TakesValue: true},
			"message":         {TakesValue: true},
			"emoji":           {TakesValue: true},
			"custom-emoji-id": {TakesValue: true},
			"big":             {},
		}, "chat", "message")
		if err != nil {
			return err
		}
		messageID, err := parsePositiveIntFlag("message", parsedArgs.First("message"))
		if err != nil {
			return err
		}
		reactions, err := parseReactionTypes(parsedArgs.First("emoji"), parsedArgs.First("custom-emoji-id"))
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.SetMessageReactionContext(cliContext(), TelegramSetReactionInput{
			Chat:      parsedArgs.First("chat"),
			MessageID: messageID,
			Reactions: reactions,
			IsBig:     parsedArgs.Has("big"),
		})
		return writeOutput(cli, out, err, args[1:], writeTelegramBool)
	case "clear":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":    {TakesValue: true},
			"message": {TakesValue: true},
		}, "chat", "message")
		if err != nil {
			return err
		}
		messageID, err := parsePositiveIntFlag("message", parsedArgs.First("message"))
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.SetMessageReactionContext(cliContext(), TelegramSetReactionInput{Chat: parsedArgs.First("chat"), MessageID: messageID})
		return writeOutput(cli, out, err, args[1:], writeTelegramBool)
	case "delete":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":          {TakesValue: true},
			"message":       {TakesValue: true},
			"user":          {TakesValue: true},
			"actor-chat-id": {TakesValue: true},
		}, "chat", "message")
		if err != nil {
			return err
		}
		messageID, err := parsePositiveIntFlag("message", parsedArgs.First("message"))
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.DeleteMessageReactionContext(cliContext(), TelegramDeleteReactionInput{
			Chat:        parsedArgs.First("chat"),
			MessageID:   messageID,
			UserID:      parsedArgs.First("user"),
			ActorChatID: parsedArgs.First("actor-chat-id"),
		})
		return writeOutput(cli, out, err, args[1:], writeTelegramBool)
	case "delete-all":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":          {TakesValue: true},
			"user":          {TakesValue: true},
			"actor-chat-id": {TakesValue: true},
		}, "chat")
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.DeleteAllMessageReactionsContext(cliContext(), TelegramDeleteAllReactionsInput{
			Chat:        parsedArgs.First("chat"),
			UserID:      parsedArgs.First("user"),
			ActorChatID: parsedArgs.First("actor-chat-id"),
		})
		return writeOutput(cli, out, err, args[1:], writeTelegramBool)
	default:
		return fmt.Errorf("unsupported reactions command: %s", args[0])
	}
}

func (cli CLI) runTopics(args []string) error {
	if isHelp(args) {
		cli.printTopicsHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printTopicsHelp()
		return nil
	}

	switch args[0] {
	case "create":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":                 {TakesValue: true},
			"name":                 {TakesValue: true},
			"icon-color":           {TakesValue: true},
			"icon-custom-emoji-id": {TakesValue: true},
		}, "chat", "name")
		if err != nil {
			return err
		}
		iconColor, err := parseOptionalPositiveIntFlag("icon-color", parsedArgs.First("icon-color"))
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.CreateForumTopicContext(cliContext(), TelegramCreateTopicInput{
			Chat:              parsedArgs.First("chat"),
			Name:              parsedArgs.First("name"),
			IconColor:         iconColor,
			IconCustomEmojiID: parsedArgs.First("icon-custom-emoji-id"),
		})
		return writeOutput(cli, out, err, args[1:], writeTelegramForumTopic)
	case "delete":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"chat":   {TakesValue: true},
			"thread": {TakesValue: true},
		}, "chat", "thread")
		if err != nil {
			return err
		}
		threadID, err := parsePositiveIntFlag("thread", parsedArgs.First("thread"))
		if err != nil {
			return err
		}
		tc, err := cli.telegramClient()
		if err != nil {
			return err
		}
		out, err := tc.DeleteForumTopicContext(cliContext(), TelegramTopicInput{Chat: parsedArgs.First("chat"), ThreadID: threadID})
		return writeOutput(cli, out, err, args[1:], writeTelegramBool)
	default:
		return fmt.Errorf("unsupported topics command: %s", args[0])
	}
}

func (cli CLI) runRequest(args []string) error {
	if isHelp(args) {
		cli.printRequestHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printRequestHelp()
		return nil
	}

	parsedArgs, err := parseRequiredArgs(args, map[string]argparse.Spec{
		"method": {TakesValue: true},
		"body":   {TakesValue: true},
	}, "method")
	if err != nil {
		return err
	}
	body, err := parseJSONBody(parsedArgs.First("body"))
	if err != nil {
		return err
	}
	tc, err := cli.telegramClient()
	if err != nil {
		return err
	}
	out, err := tc.RequestContext(cliContext(), parsedArgs.First("method"), body)
	return writeOutput(cli, out, err, args, writeTelegramRawJSON)
}

func writeOutput[T any](cli CLI, out T, err error, args []string, writeText func(io.Writer, T)) error {
	if err != nil {
		return err
	}
	if hasJSONFlag(args) {
		encoded, err := json.Marshal(out)
		if err != nil {
			return err
		}
		fmt.Fprintln(cli.stdout, string(encoded))
		return nil
	}
	writeText(cli.stdout, out)
	return nil
}

func parseRequiredArgs(args []string, specs map[string]argparse.Spec, required ...string) (argparse.Parsed, error) {
	specs["json"] = argparse.Spec{}
	parsedArgs, err := argparse.Parse(args, specs)
	if err != nil {
		return argparse.Parsed{}, err
	}
	if len(parsedArgs.Positionals) > 0 {
		return argparse.Parsed{}, fmt.Errorf("unexpected positional argument: %s", parsedArgs.Positionals[0])
	}
	for _, name := range required {
		if strings.TrimSpace(parsedArgs.First(name)) == "" {
			return argparse.Parsed{}, fmt.Errorf("--%s is required", name)
		}
	}
	return parsedArgs, nil
}

func parseOutputArgs(args []string) (argparse.Parsed, error) {
	return parseRequiredArgs(args, map[string]argparse.Spec{})
}

func parsePositiveIntFlag(name string, value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("--%s must be a positive integer", name)
	}
	return parsed, nil
}

func parseOptionalPositiveIntFlag(name string, value string) (int, error) {
	if strings.TrimSpace(value) == "" {
		return 0, nil
	}
	return parsePositiveIntFlag(name, value)
}

func parsePositiveIntCSVFlag(name string, value string) ([]int, error) {
	parts := strings.Split(value, ",")
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		parsed, err := parsePositiveIntFlag(name, strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		out = append(out, parsed)
	}
	return out, nil
}

func parseReactionTypes(emojiCSV string, customEmojiIDCSV string) ([]map[string]string, error) {
	var reactions []map[string]string
	for _, emoji := range splitNonEmptyCSV(emojiCSV) {
		reactions = append(reactions, map[string]string{"type": "emoji", "emoji": emoji})
	}
	for _, customEmojiID := range splitNonEmptyCSV(customEmojiIDCSV) {
		reactions = append(reactions, map[string]string{"type": "custom_emoji", "custom_emoji_id": customEmojiID})
	}
	if len(reactions) == 0 {
		return nil, fmt.Errorf("--emoji or --custom-emoji-id is required")
	}
	return reactions, nil
}

func splitNonEmptyCSV(value string) []string {
	var out []string
	for _, part := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func parseJSONBody(value string) (json.RawMessage, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	var body json.RawMessage
	if err := json.Unmarshal([]byte(value), &body); err != nil {
		return nil, fmt.Errorf("--body must be valid JSON: %w", err)
	}
	if !json.Valid(body) {
		return nil, fmt.Errorf("--body must be valid JSON")
	}
	return body, nil
}

func isHelp(args []string) bool {
	return len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help"
}

func hasHelpArg(args []string) bool {
	for _, arg := range args {
		if arg == "help" || arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func hasJSONFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--json" {
			return true
		}
	}
	return false
}

func cliContext() context.Context {
	return context.Background()
}

func writeTelegramUser(w io.Writer, user TelegramUser) {
	fmt.Fprintf(w, "id\t%d\nusername\t%s\nfirst_name\t%s\nbot\t%t\n", user.ID, user.Username, user.FirstName, user.IsBot)
}

func writeTelegramChat(w io.Writer, chat TelegramChat) {
	fmt.Fprintf(w, "id\t%d\ntype\t%s\ntitle\t%s\nusername\t%s\n", chat.ID, chat.Type, chat.Title, chat.Username)
}

func writeTelegramMessage(w io.Writer, message TelegramMessage) {
	fmt.Fprintf(w, "id\t%d\nchat\t%d\ntext\t%s\n", message.MessageID, message.Chat.ID, message.Text)
}

func writeTelegramForumTopic(w io.Writer, topic TelegramForumTopic) {
	fmt.Fprintf(w, "thread\t%d\nname\t%s\nicon_color\t%d\n", topic.MessageThreadID, topic.Name, topic.IconColor)
}

func writeTelegramBool(w io.Writer, out TelegramBoolResponse) {
	fmt.Fprintf(w, "ok\t%t\n", out.OK)
}

func writeTelegramRawJSON(w io.Writer, out json.RawMessage) {
	fmt.Fprintln(w, string(out))
}

func (cli CLI) printHelp() {
	fmt.Fprint(cli.stdout, `telegram

CLI for Telegram Bot API.

Usage:
  telegram help
  telegram version
  telegram auth help
  telegram chats help
  telegram messages help
  telegram reactions help
  telegram topics help
  telegram request help
  telegram mcp help

Commands:
  help
  version
  auth
  chats
  messages
  reactions
  topics
  request
  mcp
`)
}

func (cli CLI) printAuthHelp() {
	fmt.Fprint(cli.stdout, "telegram auth\n\nUsage:\n  telegram auth test [--json]\n")
}

func (cli CLI) printChatsHelp() {
	fmt.Fprint(cli.stdout, "telegram chats\n\nUsage:\n  telegram chats get --chat <chat-id-or-username> [--json]\n")
}

func (cli CLI) printMessagesHelp() {
	fmt.Fprint(cli.stdout, "telegram messages\n\nUsage:\n  telegram messages send --chat <chat-id-or-username> --text <text> [--thread <message-thread-id>] [--parse-mode <mode>] [--json]\n  telegram messages edit --chat <chat-id-or-username> --message <message-id> --text <text> [--parse-mode <mode>] [--json]\n  telegram messages delete --chat <chat-id-or-username> --message <message-id> [--json]\n  telegram messages delete-batch --chat <chat-id-or-username> --messages <message-id-csv> [--json]\n")
}

func (cli CLI) printReactionsHelp() {
	fmt.Fprint(cli.stdout, "telegram reactions\n\nUsage:\n  telegram reactions set --chat <chat-id-or-username> --message <message-id> --emoji <emoji-csv> [--custom-emoji-id <id-csv>] [--big] [--json]\n  telegram reactions clear --chat <chat-id-or-username> --message <message-id> [--json]\n  telegram reactions delete --chat <chat-id-or-username> --message <message-id> [--user <user-id>] [--actor-chat-id <chat-id>] [--json]\n  telegram reactions delete-all --chat <chat-id-or-username> [--user <user-id>] [--actor-chat-id <chat-id>] [--json]\n")
}

func (cli CLI) printTopicsHelp() {
	fmt.Fprint(cli.stdout, "telegram topics\n\nUsage:\n  telegram topics create --chat <chat-id-or-username> --name <name> [--icon-color <rgb-int>] [--icon-custom-emoji-id <emoji-id>] [--json]\n  telegram topics delete --chat <chat-id-or-username> --thread <message-thread-id> [--json]\n")
}

func (cli CLI) printRequestHelp() {
	fmt.Fprint(cli.stdout, "telegram request\n\nUsage:\n  telegram request --method <telegram-method> [--body <json>] [--json]\n")
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
