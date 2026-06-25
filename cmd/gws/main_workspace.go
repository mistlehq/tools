package main

import (
	"fmt"

	"github.com/mistlehq/tools/internal/argparse"
)

func (cli CLI) runGmail(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printGmailHelp()
		return nil
	}
	switch args[0] {
	case "messages":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printGmailMessagesHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runGmailMessages(gc, args[1:])
	case "drafts":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printGmailDraftsHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runGmailDrafts(gc, args[1:])
	default:
		return fmt.Errorf("unsupported gmail command: %s", args[0])
	}
}

func (cli CLI) runGmailMessages(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printGmailMessagesHelp()
		return nil
	}
	switch args[0] {
	case "list":
		userID, params, err := parseUserIDAndParamsArgs(args[1:], "messages list", []string{"query", "label-ids", "max-results", "page-token"})
		if err != nil {
			return err
		}
		out, err := gc.ListGmailMessages(cliContext(), userID, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		userID, id, params, err := parseUserIDChildIDAndParamsArgs(args[1:], "messages get", "message-id", []string{"format"})
		if err != nil {
			return err
		}
		out, err := gc.GetGmailMessage(cliContext(), userID, id, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "send":
		userID, body, err := parseUserIDRequestFileArgs(args[1:], "messages send")
		if err != nil {
			return err
		}
		out, err := gc.SendGmailMessage(cliContext(), userID, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported gmail messages command: %s", args[0])
	}
}

func (cli CLI) runGmailDrafts(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printGmailDraftsHelp()
		return nil
	}
	switch args[0] {
	case "list":
		userID, params, err := parseUserIDAndParamsArgs(args[1:], "drafts list", []string{"max-results", "page-token"})
		if err != nil {
			return err
		}
		out, err := gc.ListGmailDrafts(cliContext(), userID, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		userID, id, params, err := parseUserIDChildIDAndParamsArgs(args[1:], "drafts get", "draft-id", []string{"format"})
		if err != nil {
			return err
		}
		out, err := gc.GetGmailDraft(cliContext(), userID, id, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "create":
		userID, body, err := parseUserIDRequestFileArgs(args[1:], "drafts create")
		if err != nil {
			return err
		}
		out, err := gc.CreateGmailDraft(cliContext(), userID, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "send":
		userID, body, err := parseUserIDRequestFileArgs(args[1:], "drafts send")
		if err != nil {
			return err
		}
		out, err := gc.SendGmailDraft(cliContext(), userID, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "delete":
		userID, draftID, _, err := parseUserIDChildIDAndParamsArgs(args[1:], "drafts delete", "draft-id", nil)
		if err != nil {
			return err
		}
		if err := gc.DeleteGmailDraft(cliContext(), userID, draftID); err != nil {
			return err
		}
		return writeJSON(cli.stdout, map[string]any{"deleted": true, "userId": userID, "draftId": draftID})
	default:
		return fmt.Errorf("unsupported gmail drafts command: %s", args[0])
	}
}

func (cli CLI) runCalendar(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printCalendarHelp()
		return nil
	}
	switch args[0] {
	case "calendar-list":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printCalendarListHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runCalendarList(gc, args[1:])
	case "events":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printCalendarEventsHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runCalendarEvents(gc, args[1:])
	case "freebusy":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printCalendarFreeBusyHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runCalendarFreeBusy(gc, args[1:])
	default:
		return fmt.Errorf("unsupported calendar command: %s", args[0])
	}
}

func (cli CLI) runCalendarList(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printCalendarListHelp()
		return nil
	}
	switch args[0] {
	case "list":
		params, err := parseParamsOnlyArgs(args[1:], "calendar-list list", []string{"max-results", "page-token"})
		if err != nil {
			return err
		}
		out, err := gc.ListCalendarList(cliContext(), params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		calendarID, _, err := parseIDAndParamsArgs(args[1:], "calendar-list get", "calendar-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.GetCalendarListEntry(cliContext(), calendarID)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported calendar calendar-list command: %s", args[0])
	}
}

func (cli CLI) runCalendarEvents(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printCalendarEventsHelp()
		return nil
	}
	switch args[0] {
	case "list":
		calendarID, params, err := parseIDAndParamsArgs(args[1:], "events list", "calendar-id", []string{"time-min", "time-max", "max-results", "single-events", "order-by", "page-token"})
		if err != nil {
			return err
		}
		out, err := gc.ListCalendarEvents(cliContext(), calendarID, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		calendarID, eventID, _, err := parseTwoIDAndParamsArgs(args[1:], "events get", "calendar-id", "event-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.GetCalendarEvent(cliContext(), calendarID, eventID)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "insert":
		calendarID, body, _, err := parseIDRequestFileAndParamsArgs(args[1:], "events insert", "calendar-id", nil)
		if err != nil {
			return err
		}
		out, err := gc.InsertCalendarEvent(cliContext(), calendarID, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "patch":
		calendarID, eventID, body, err := parseTwoIDRequestFileArgs(args[1:], "events patch", "calendar-id", "event-id")
		if err != nil {
			return err
		}
		out, err := gc.PatchCalendarEvent(cliContext(), calendarID, eventID, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "delete":
		calendarID, eventID, _, err := parseTwoIDAndParamsArgs(args[1:], "events delete", "calendar-id", "event-id", nil)
		if err != nil {
			return err
		}
		if err := gc.DeleteCalendarEvent(cliContext(), calendarID, eventID); err != nil {
			return err
		}
		return writeJSON(cli.stdout, map[string]any{"deleted": true, "calendarId": calendarID, "eventId": eventID})
	default:
		return fmt.Errorf("unsupported calendar events command: %s", args[0])
	}
}

func (cli CLI) runCalendarFreeBusy(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printCalendarFreeBusyHelp()
		return nil
	}
	if args[0] != "query" {
		return fmt.Errorf("unsupported calendar freebusy command: %s", args[0])
	}
	body, _, err := parseRequestFileAndParamsArgs(args[1:], "freebusy query", nil)
	if err != nil {
		return err
	}
	out, err := gc.QueryCalendarFreeBusy(cliContext(), body)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runChat(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printChatHelp()
		return nil
	}
	switch args[0] {
	case "spaces":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printChatSpacesHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runChatSpaces(gc, args[1:])
	case "messages":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printChatMessagesHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runChatMessages(gc, args[1:])
	case "members":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printChatMembersHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runChatMembers(gc, args[1:])
	default:
		return fmt.Errorf("unsupported chat command: %s", args[0])
	}
}

func (cli CLI) runChatSpaces(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printChatSpacesHelp()
		return nil
	}
	switch args[0] {
	case "list":
		params, err := parseParamsOnlyArgs(args[1:], "spaces list", []string{"page-size", "page-token"})
		if err != nil {
			return err
		}
		out, err := gc.ListChatSpaces(cliContext(), params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		spaceName, _, err := parseIDAndParamsArgs(args[1:], "spaces get", "space-name", nil)
		if err != nil {
			return err
		}
		out, err := gc.GetChatSpace(cliContext(), spaceName)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported chat spaces command: %s", args[0])
	}
}

func (cli CLI) runChatMessages(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printChatMessagesHelp()
		return nil
	}
	switch args[0] {
	case "list":
		spaceName, params, err := parseIDAndParamsArgs(args[1:], "messages list", "space-name", []string{"page-size", "page-token"})
		if err != nil {
			return err
		}
		out, err := gc.ListChatMessages(cliContext(), spaceName, params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "get":
		messageName, _, err := parseIDAndParamsArgs(args[1:], "messages get", "message-name", nil)
		if err != nil {
			return err
		}
		out, err := gc.GetChatMessage(cliContext(), messageName)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "create":
		spaceName, body, _, err := parseIDRequestFileAndParamsArgs(args[1:], "messages create", "space-name", nil)
		if err != nil {
			return err
		}
		out, err := gc.CreateChatMessage(cliContext(), spaceName, body)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported chat messages command: %s", args[0])
	}
}

func (cli CLI) runChatMembers(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printChatMembersHelp()
		return nil
	}
	if args[0] != "list" {
		return fmt.Errorf("unsupported chat members command: %s", args[0])
	}
	spaceName, params, err := parseIDAndParamsArgs(args[1:], "members list", "space-name", []string{"page-size", "page-token"})
	if err != nil {
		return err
	}
	out, err := gc.ListChatMembers(cliContext(), spaceName, params)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runPeople(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		cli.printPeopleHelp()
		return nil
	}
	if args[0] == "search-contacts" && isSingleHelpArg(args[1:]) {
		cli.printPeopleSearchContactsHelp()
		return nil
	}
	if args[0] == "search-directory" && isSingleHelpArg(args[1:]) {
		cli.printPeopleSearchDirectoryHelp()
		return nil
	}
	switch args[0] {
	case "people":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printPeoplePeopleHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runPeoplePeople(gc, args[1:])
	case "connections":
		if len(args) < 2 || isHelpToken(args[1]) {
			cli.printPeopleConnectionsHelp()
			return nil
		}
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		return cli.runPeopleConnections(gc, args[1:])
	case "search-contacts":
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		params, err := parseRequiredQueryMaskArgs(args[1:], "search-contacts", "read-mask", []string{"page-size"})
		if err != nil {
			return err
		}
		out, err := gc.SearchPeopleContacts(cliContext(), params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	case "search-directory":
		gc, err := cli.gwsClient()
		if err != nil {
			return err
		}
		params, err := parseRequiredQueryMaskArgs(args[1:], "search-directory", "read-mask", []string{"sources", "page-size"})
		if err != nil {
			return err
		}
		out, err := gc.SearchPeopleDirectory(cliContext(), params)
		if err != nil {
			return err
		}
		return writeJSON(cli.stdout, out)
	default:
		return fmt.Errorf("unsupported people command: %s", args[0])
	}
}

func (cli CLI) runPeoplePeople(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printPeoplePeopleHelp()
		return nil
	}
	if args[0] != "get" {
		return fmt.Errorf("unsupported people people command: %s", args[0])
	}
	resourceName, params, err := parseIDAndParamsArgs(args[1:], "people get", "resource-name", []string{"person-fields"})
	if err != nil {
		return err
	}
	out, err := gc.GetPeoplePerson(cliContext(), resourceName, params)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func (cli CLI) runPeopleConnections(gc GWSClient, args []string) error {
	if isSingleHelpArg(args[1:]) {
		cli.printPeopleConnectionsHelp()
		return nil
	}
	if args[0] != "list" {
		return fmt.Errorf("unsupported people connections command: %s", args[0])
	}
	resourceName, params, err := parseIDAndParamsArgs(args[1:], "connections list", "resource-name", []string{"person-fields", "page-size", "page-token"})
	if err != nil {
		return err
	}
	out, err := gc.ListPeopleConnections(cliContext(), resourceName, params)
	if err != nil {
		return err
	}
	return writeJSON(cli.stdout, out)
}

func parseParamsOnlyArgs(args []string, command string, paramFlags []string) (map[string]any, error) {
	specs := map[string]argparse.Spec{}
	for _, flag := range paramFlags {
		specs[flag] = argparse.Spec{TakesValue: true}
	}
	parsed, err := argparse.Parse(args, specs)
	if err != nil {
		return nil, err
	}
	if len(parsed.Positionals) > 0 {
		return nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	params := map[string]any{}
	for _, flag := range paramFlags {
		copyFlag(params, parsed, flag, hyphenFlagToGoogleParam(flag))
	}
	return params, nil
}

func parseUserIDAndParamsArgs(args []string, command string, paramFlags []string) (string, map[string]any, error) {
	userID, params, err := parseIDAndParamsArgs(args, command, "user-id", paramFlags)
	if err != nil {
		return "", nil, err
	}
	return userID, normalizeGoogleParams(params), nil
}

func parseUserIDChildIDAndParamsArgs(args []string, command string, childIDFlag string, paramFlags []string) (string, string, map[string]any, error) {
	userID, childID, params, err := parseTwoIDAndParamsArgs(args, command, "user-id", childIDFlag, paramFlags)
	if err != nil {
		return "", "", nil, err
	}
	return userID, childID, normalizeGoogleParams(params), nil
}

func parseUserIDRequestFileArgs(args []string, command string) (string, map[string]any, error) {
	userID, body, _, err := parseIDRequestFileAndParamsArgs(args, command, "user-id", nil)
	return userID, body, err
}

func parseTwoIDAndParamsArgs(args []string, command string, firstIDFlag string, secondIDFlag string, paramFlags []string) (string, string, map[string]any, error) {
	specs := map[string]argparse.Spec{firstIDFlag: {TakesValue: true}, secondIDFlag: {TakesValue: true}}
	for _, flag := range paramFlags {
		specs[flag] = argparse.Spec{TakesValue: true}
	}
	parsed, err := argparse.Parse(args, specs)
	if err != nil {
		return "", "", nil, err
	}
	if len(parsed.Positionals) > 0 {
		return "", "", nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	firstID := parsed.First(firstIDFlag)
	if firstID == "" {
		return "", "", nil, fmt.Errorf("%s requires --%s", command, firstIDFlag)
	}
	secondID := parsed.First(secondIDFlag)
	if secondID == "" {
		return "", "", nil, fmt.Errorf("%s requires --%s", command, secondIDFlag)
	}
	params := map[string]any{}
	for _, flag := range paramFlags {
		copyFlag(params, parsed, flag, flag)
	}
	return firstID, secondID, normalizeGoogleParams(params), nil
}

func parseTwoIDRequestFileArgs(args []string, command string, firstIDFlag string, secondIDFlag string) (string, string, map[string]any, error) {
	specs := map[string]argparse.Spec{firstIDFlag: {TakesValue: true}, secondIDFlag: {TakesValue: true}, "request-file": {TakesValue: true}}
	parsed, err := argparse.Parse(args, specs)
	if err != nil {
		return "", "", nil, err
	}
	if len(parsed.Positionals) > 0 {
		return "", "", nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	firstID := parsed.First(firstIDFlag)
	if firstID == "" {
		return "", "", nil, fmt.Errorf("%s requires --%s", command, firstIDFlag)
	}
	secondID := parsed.First(secondIDFlag)
	if secondID == "" {
		return "", "", nil, fmt.Errorf("%s requires --%s", command, secondIDFlag)
	}
	body, err := readRequestFileAsMap(parsed.First("request-file"), command)
	return firstID, secondID, body, err
}

func parseRequiredQueryMaskArgs(args []string, command string, maskFlag string, optionalFlags []string) (map[string]any, error) {
	specs := map[string]argparse.Spec{"query": {TakesValue: true}, maskFlag: {TakesValue: true}}
	for _, flag := range optionalFlags {
		specs[flag] = argparse.Spec{TakesValue: true}
	}
	parsed, err := argparse.Parse(args, specs)
	if err != nil {
		return nil, err
	}
	if len(parsed.Positionals) > 0 {
		return nil, fmt.Errorf("%s does not accept positional arguments", command)
	}
	if parsed.First("query") == "" {
		return nil, fmt.Errorf("%s requires --query", command)
	}
	if parsed.First(maskFlag) == "" {
		return nil, fmt.Errorf("%s requires --%s", command, maskFlag)
	}
	params := map[string]any{"query": parsed.First("query"), hyphenFlagToGoogleParam(maskFlag): parsed.First(maskFlag)}
	for _, flag := range optionalFlags {
		copyFlag(params, parsed, flag, hyphenFlagToGoogleParam(flag))
	}
	return params, nil
}

func normalizeGoogleParams(params map[string]any) map[string]any {
	out := map[string]any{}
	for key, value := range params {
		out[hyphenFlagToGoogleParam(key)] = value
	}
	return out
}

func hyphenFlagToGoogleParam(flag string) string {
	switch flag {
	case "label-ids":
		return "labelIds"
	case "max-results":
		return "maxResults"
	case "page-token":
		return "pageToken"
	case "page-size":
		return "pageSize"
	case "time-min":
		return "timeMin"
	case "time-max":
		return "timeMax"
	case "single-events":
		return "singleEvents"
	case "order-by":
		return "orderBy"
	case "person-fields":
		return "personFields"
	case "read-mask":
		return "readMask"
	default:
		return flag
	}
}
