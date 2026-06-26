package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mistlehq/tools/internal/argparse"
	"io"
	"os"
	"strings"
)

// Version is the current discord CLI version.
var Version = "dev"

type CLI struct {
	stdout io.Writer
	stderr io.Writer
	env    Environment
}

func (cli CLI) discordClient() (DiscordClient, error) {
	config, err := loadConfig(cli.env)
	if err != nil {
		return DiscordClient{}, err
	}

	return NewDiscordClient(config), nil
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
	case "guilds":
		return cli.runGuilds(args[2:])
	case "channels":
		return cli.runChannels(args[2:])
	case "messages":
		return cli.runMessages(args[2:])
	case "reactions":
		return cli.runReactions(args[2:])
	case "roles":
		return cli.runRoles(args[2:])
	case "members":
		return cli.runMembers(args[2:])
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
	dc, err := cli.discordClient()
	if err != nil {
		return err
	}
	out, err := dc.AuthTestContext(cliContext())
	return writeOutput(cli, out, err, args[1:], writeDiscordUser)
}

func (cli CLI) runGuilds(args []string) error {
	if isHelp(args) {
		cli.printGuildsHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printGuildsHelp()
		return nil
	}

	switch args[0] {
	case "list":
		if _, err := parseOutputArgs(args[1:]); err != nil {
			return err
		}
		dc, err := cli.discordClient()
		if err != nil {
			return err
		}
		out, err := dc.ListGuildsContext(cliContext())
		return writeOutput(cli, out, err, args[1:], writeDiscordGuilds)
	case "get":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"guild": {TakesValue: true}}, "guild")
		if err != nil {
			return err
		}
		dc, err := cli.discordClient()
		if err != nil {
			return err
		}
		out, err := dc.GetGuildContext(cliContext(), parsedArgs.First("guild"))
		return writeOutput(cli, out, err, args[1:], writeDiscordGuild)
	default:
		return fmt.Errorf("unsupported guilds command: %s", args[0])
	}
}

func (cli CLI) runChannels(args []string) error {
	if isHelp(args) {
		cli.printChannelsHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printChannelsHelp()
		return nil
	}
	dc, err := cli.discordClient()
	if err != nil {
		return err
	}

	switch args[0] {
	case "list":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"guild": {TakesValue: true}}, "guild")
		if err != nil {
			return err
		}
		out, err := dc.ListChannelsContext(cliContext(), parsedArgs.First("guild"))
		return writeOutput(cli, out, err, args[1:], writeDiscordChannels)
	case "get":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"channel": {TakesValue: true}}, "channel")
		if err != nil {
			return err
		}
		out, err := dc.GetChannelContext(cliContext(), parsedArgs.First("channel"))
		return writeOutput(cli, out, err, args[1:], writeDiscordChannel)
	default:
		return fmt.Errorf("unsupported channels command: %s", args[0])
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
	dc, err := cli.discordClient()
	if err != nil {
		return err
	}

	switch args[0] {
	case "list":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{
			"channel": {TakesValue: true},
			"limit":   {TakesValue: true},
			"before":  {TakesValue: true},
			"after":   {TakesValue: true},
			"around":  {TakesValue: true},
		}, "channel")
		if err != nil {
			return err
		}
		out, err := dc.ListMessagesContext(cliContext(), DiscordListMessagesInput{
			Channel: parsedArgs.First("channel"),
			Limit:   parsedArgs.First("limit"),
			Before:  parsedArgs.First("before"),
			After:   parsedArgs.First("after"),
			Around:  parsedArgs.First("around"),
		})
		return writeOutput(cli, out, err, args[1:], writeDiscordMessages)
	case "send":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"channel": {TakesValue: true}, "content": {TakesValue: true}}, "channel", "content")
		if err != nil {
			return err
		}
		out, err := dc.CreateMessageContext(cliContext(), DiscordCreateMessageInput{Channel: parsedArgs.First("channel"), Content: parsedArgs.First("content")})
		return writeOutput(cli, out, err, args[1:], writeDiscordMessage)
	case "edit":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"channel": {TakesValue: true}, "message": {TakesValue: true}, "content": {TakesValue: true}}, "channel", "message", "content")
		if err != nil {
			return err
		}
		out, err := dc.EditMessageContext(cliContext(), DiscordEditMessageInput{Channel: parsedArgs.First("channel"), Message: parsedArgs.First("message"), Content: parsedArgs.First("content")})
		return writeOutput(cli, out, err, args[1:], writeDiscordMessage)
	case "delete":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"channel": {TakesValue: true}, "message": {TakesValue: true}}, "channel", "message")
		if err != nil {
			return err
		}
		out, err := dc.DeleteMessageContext(cliContext(), parsedArgs.First("channel"), parsedArgs.First("message"))
		return writeOutput(cli, out, err, args[1:], writeDiscordEmpty)
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
	dc, err := cli.discordClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "add":
		input, err := parseReactionArgs(args[1:])
		if err != nil {
			return err
		}
		out, err := dc.AddReactionContext(cliContext(), input)
		return writeOutput(cli, out, err, args[1:], writeDiscordEmpty)
	case "remove":
		input, err := parseReactionArgs(args[1:])
		if err != nil {
			return err
		}
		out, err := dc.RemoveReactionContext(cliContext(), input)
		return writeOutput(cli, out, err, args[1:], writeDiscordEmpty)
	default:
		return fmt.Errorf("unsupported reactions command: %s", args[0])
	}
}

func (cli CLI) runRoles(args []string) error {
	if isHelp(args) {
		cli.printRolesHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printRolesHelp()
		return nil
	}
	dc, err := cli.discordClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"guild": {TakesValue: true}}, "guild")
		if err != nil {
			return err
		}
		out, err := dc.ListRolesContext(cliContext(), parsedArgs.First("guild"))
		return writeOutput(cli, out, err, args[1:], writeDiscordRoles)
	case "create":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"guild": {TakesValue: true}, "name": {TakesValue: true}}, "guild", "name")
		if err != nil {
			return err
		}
		out, err := dc.CreateRoleContext(cliContext(), parsedArgs.First("guild"), parsedArgs.First("name"))
		return writeOutput(cli, out, err, args[1:], writeDiscordRole)
	case "delete":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"guild": {TakesValue: true}, "role": {TakesValue: true}}, "guild", "role")
		if err != nil {
			return err
		}
		out, err := dc.DeleteRoleContext(cliContext(), parsedArgs.First("guild"), parsedArgs.First("role"))
		return writeOutput(cli, out, err, args[1:], writeDiscordEmpty)
	default:
		return fmt.Errorf("unsupported roles command: %s", args[0])
	}
}

func (cli CLI) runMembers(args []string) error {
	if isHelp(args) {
		cli.printMembersHelp()
		return nil
	}
	if hasHelpArg(args[1:]) {
		cli.printMembersHelp()
		return nil
	}
	dc, err := cli.discordClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"guild": {TakesValue: true}, "limit": {TakesValue: true}, "after": {TakesValue: true}}, "guild")
		if err != nil {
			return err
		}
		out, err := dc.ListMembersContext(cliContext(), DiscordListMembersInput{Guild: parsedArgs.First("guild"), Limit: parsedArgs.First("limit"), After: parsedArgs.First("after")})
		return writeOutput(cli, out, err, args[1:], writeDiscordMembers)
	case "get":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"guild": {TakesValue: true}, "user": {TakesValue: true}}, "guild", "user")
		if err != nil {
			return err
		}
		out, err := dc.GetMemberContext(cliContext(), parsedArgs.First("guild"), parsedArgs.First("user"))
		return writeOutput(cli, out, err, args[1:], writeDiscordMember)
	case "add-role":
		guild, user, role, err := parseMemberRoleArgs(args[1:])
		if err != nil {
			return err
		}
		out, err := dc.AddMemberRoleContext(cliContext(), guild, user, role)
		return writeOutput(cli, out, err, args[1:], writeDiscordEmpty)
	case "remove-role":
		guild, user, role, err := parseMemberRoleArgs(args[1:])
		if err != nil {
			return err
		}
		out, err := dc.RemoveMemberRoleContext(cliContext(), guild, user, role)
		return writeOutput(cli, out, err, args[1:], writeDiscordEmpty)
	case "ban":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"guild": {TakesValue: true}, "user": {TakesValue: true}}, "guild", "user")
		if err != nil {
			return err
		}
		out, err := dc.BanMemberContext(cliContext(), parsedArgs.First("guild"), parsedArgs.First("user"))
		return writeOutput(cli, out, err, args[1:], writeDiscordEmpty)
	case "unban":
		parsedArgs, err := parseRequiredArgs(args[1:], map[string]argparse.Spec{"guild": {TakesValue: true}, "user": {TakesValue: true}}, "guild", "user")
		if err != nil {
			return err
		}
		out, err := dc.UnbanMemberContext(cliContext(), parsedArgs.First("guild"), parsedArgs.First("user"))
		return writeOutput(cli, out, err, args[1:], writeDiscordEmpty)
	default:
		return fmt.Errorf("unsupported members command: %s", args[0])
	}
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

func parseReactionArgs(args []string) (DiscordReactionInput, error) {
	parsedArgs, err := parseRequiredArgs(args, map[string]argparse.Spec{"channel": {TakesValue: true}, "message": {TakesValue: true}, "emoji": {TakesValue: true}}, "channel", "message", "emoji")
	if err != nil {
		return DiscordReactionInput{}, err
	}
	return DiscordReactionInput{Channel: parsedArgs.First("channel"), Message: parsedArgs.First("message"), Emoji: parsedArgs.First("emoji")}, nil
}

func parseMemberRoleArgs(args []string) (string, string, string, error) {
	parsedArgs, err := parseRequiredArgs(args, map[string]argparse.Spec{"guild": {TakesValue: true}, "user": {TakesValue: true}, "role": {TakesValue: true}}, "guild", "user", "role")
	if err != nil {
		return "", "", "", err
	}
	return parsedArgs.First("guild"), parsedArgs.First("user"), parsedArgs.First("role"), nil
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

func writeDiscordUser(w io.Writer, user DiscordUser) {
	fmt.Fprintf(w, "id\t%s\nusername\t%s\nbot\t%t\n", user.ID, user.Username, user.Bot)
}

func writeDiscordGuild(w io.Writer, guild DiscordGuild) {
	fmt.Fprintf(w, "id\t%s\nname\t%s\n", guild.ID, guild.Name)
}

func writeDiscordGuilds(w io.Writer, guilds []DiscordGuild) {
	for _, guild := range guilds {
		fmt.Fprintf(w, "%s\t%s\n", guild.ID, guild.Name)
	}
}

func writeDiscordChannel(w io.Writer, channel DiscordChannel) {
	fmt.Fprintf(w, "id\t%s\nname\t%s\ntype\t%d\n", channel.ID, channel.Name, channel.Type)
}

func writeDiscordChannels(w io.Writer, channels []DiscordChannel) {
	for _, channel := range channels {
		fmt.Fprintf(w, "%s\t%d\t%s\n", channel.ID, channel.Type, channel.Name)
	}
}

func writeDiscordMessage(w io.Writer, message DiscordMessage) {
	fmt.Fprintf(w, "id\t%s\nchannel\t%s\nauthor\t%s\ncontent\t%s\n", message.ID, message.ChannelID, message.Author.ID, message.Content)
}

func writeDiscordMessages(w io.Writer, messages []DiscordMessage) {
	for _, message := range messages {
		fmt.Fprintf(w, "%s\t%s\t%s\n", message.ID, message.Author.ID, strings.ReplaceAll(message.Content, "\n", "\\n"))
	}
}

func writeDiscordRole(w io.Writer, role DiscordRole) {
	fmt.Fprintf(w, "id\t%s\nname\t%s\nposition\t%d\n", role.ID, role.Name, role.Position)
}

func writeDiscordRoles(w io.Writer, roles []DiscordRole) {
	for _, role := range roles {
		fmt.Fprintf(w, "%s\t%d\t%s\n", role.ID, role.Position, role.Name)
	}
}

func writeDiscordMember(w io.Writer, member DiscordMember) {
	fmt.Fprintf(w, "user\t%s\nusername\t%s\nnick\t%s\nroles\t%s\n", member.User.ID, member.User.Username, member.Nick, strings.Join(member.Roles, ","))
}

func writeDiscordMembers(w io.Writer, members []DiscordMember) {
	for _, member := range members {
		fmt.Fprintf(w, "%s\t%s\t%s\n", member.User.ID, member.User.Username, member.Nick)
	}
}

func writeDiscordEmpty(w io.Writer, out DiscordEmptyResponse) {
	fmt.Fprintf(w, "ok\t%t\n", out.OK)
}

func (cli CLI) printHelp() {
	fmt.Fprint(cli.stdout, `discord

CLI for Discord.

Usage:
  discord help
  discord version
  discord auth help
  discord guilds help
  discord channels help
  discord messages help
  discord reactions help
  discord roles help
  discord members help
  discord mcp help

Commands:
  help
  version
  auth
  guilds
  channels
  messages
  reactions
  roles
  members
  mcp
`)
}

func (cli CLI) printAuthHelp() {
	fmt.Fprint(cli.stdout, "discord auth\n\nUsage:\n  discord auth test [--json]\n")
}

func (cli CLI) printGuildsHelp() {
	fmt.Fprint(cli.stdout, "discord guilds\n\nUsage:\n  discord guilds list [--json]\n  discord guilds get --guild <guild-id> [--json]\n")
}

func (cli CLI) printChannelsHelp() {
	fmt.Fprint(cli.stdout, "discord channels\n\nUsage:\n  discord channels list --guild <guild-id> [--json]\n  discord channels get --channel <channel-id> [--json]\n")
}

func (cli CLI) printMessagesHelp() {
	fmt.Fprint(cli.stdout, "discord messages\n\nUsage:\n  discord messages list --channel <channel-id> [--limit <n>] [--before <message-id>] [--after <message-id>] [--around <message-id>] [--json]\n  discord messages send --channel <channel-id> --content <text> [--json]\n  discord messages edit --channel <channel-id> --message <message-id> --content <text> [--json]\n  discord messages delete --channel <channel-id> --message <message-id> [--json]\n")
}

func (cli CLI) printReactionsHelp() {
	fmt.Fprint(cli.stdout, "discord reactions\n\nUsage:\n  discord reactions add --channel <channel-id> --message <message-id> --emoji <emoji> [--json]\n  discord reactions remove --channel <channel-id> --message <message-id> --emoji <emoji> [--json]\n")
}

func (cli CLI) printRolesHelp() {
	fmt.Fprint(cli.stdout, "discord roles\n\nUsage:\n  discord roles list --guild <guild-id> [--json]\n  discord roles create --guild <guild-id> --name <name> [--json]\n  discord roles delete --guild <guild-id> --role <role-id> [--json]\n")
}

func (cli CLI) printMembersHelp() {
	fmt.Fprint(cli.stdout, "discord members\n\nUsage:\n  discord members list --guild <guild-id> [--limit <n>] [--after <user-id>] [--json]\n  discord members get --guild <guild-id> --user <user-id> [--json]\n  discord members add-role --guild <guild-id> --user <user-id> --role <role-id> [--json]\n  discord members remove-role --guild <guild-id> --user <user-id> --role <role-id> [--json]\n  discord members ban --guild <guild-id> --user <user-id> [--json]\n  discord members unban --guild <guild-id> --user <user-id> [--json]\n")
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
