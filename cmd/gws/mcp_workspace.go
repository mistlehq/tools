package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type gwsGmailUserToolInput struct {
	UserID     string `json:"userId" jsonschema:"Gmail user ID. Use me for the authenticated user."`
	Query      string `json:"query,omitempty" jsonschema:"Optional Gmail q search query."`
	LabelIDs   string `json:"labelIds,omitempty" jsonschema:"Optional comma-separated Gmail label IDs."`
	MaxResults string `json:"maxResults,omitempty" jsonschema:"Optional maxResults value."`
	PageToken  string `json:"pageToken,omitempty" jsonschema:"Optional page token."`
}

type gwsGmailMessageToolInput struct {
	UserID    string         `json:"userId" jsonschema:"Gmail user ID. Use me for the authenticated user."`
	MessageID string         `json:"messageId,omitempty" jsonschema:"Gmail message ID where required."`
	DraftID   string         `json:"draftId,omitempty" jsonschema:"Gmail draft ID where required."`
	Format    string         `json:"format,omitempty" jsonschema:"Optional Gmail format value."`
	Request   map[string]any `json:"request,omitempty" jsonschema:"Google Gmail API request body using Google's documented JSON shape."`
}

type gwsCalendarListToolInput struct {
	CalendarID string `json:"calendarId,omitempty" jsonschema:"Google Calendar calendar ID where required."`
	MaxResults string `json:"maxResults,omitempty" jsonschema:"Optional maxResults value."`
	PageToken  string `json:"pageToken,omitempty" jsonschema:"Optional page token."`
}

type gwsCalendarEventToolInput struct {
	CalendarID   string         `json:"calendarId" jsonschema:"Google Calendar calendar ID."`
	EventID      string         `json:"eventId,omitempty" jsonschema:"Google Calendar event ID where required."`
	TimeMin      string         `json:"timeMin,omitempty" jsonschema:"Optional RFC3339 lower bound."`
	TimeMax      string         `json:"timeMax,omitempty" jsonschema:"Optional RFC3339 upper bound."`
	MaxResults   string         `json:"maxResults,omitempty" jsonschema:"Optional maxResults value."`
	SingleEvents string         `json:"singleEvents,omitempty" jsonschema:"Optional singleEvents value."`
	OrderBy      string         `json:"orderBy,omitempty" jsonschema:"Optional orderBy value."`
	PageToken    string         `json:"pageToken,omitempty" jsonschema:"Optional page token."`
	Request      map[string]any `json:"request,omitempty" jsonschema:"Google Calendar API request body using Google's documented JSON shape."`
}

type gwsRequestBodyToolInput struct {
	Request map[string]any `json:"request" jsonschema:"Google Workspace API request body using Google's documented JSON shape."`
}

type gwsChatListToolInput struct {
	SpaceName   string         `json:"spaceName,omitempty" jsonschema:"Google Chat space resource name, for example spaces/AAAA."`
	MessageName string         `json:"messageName,omitempty" jsonschema:"Google Chat message resource name, for example spaces/AAAA/messages/BBBB."`
	PageSize    string         `json:"pageSize,omitempty" jsonschema:"Optional page size."`
	PageToken   string         `json:"pageToken,omitempty" jsonschema:"Optional page token."`
	Request     map[string]any `json:"request,omitempty" jsonschema:"Google Chat API request body using Google's documented JSON shape."`
}

type gwsPeopleToolInput struct {
	ResourceName string `json:"resourceName,omitempty" jsonschema:"People API resource name, for example people/me."`
	PersonFields string `json:"personFields,omitempty" jsonschema:"People API personFields selector."`
	Query        string `json:"query,omitempty" jsonschema:"Search query."`
	ReadMask     string `json:"readMask,omitempty" jsonschema:"People API readMask selector."`
	Sources      string `json:"sources,omitempty" jsonschema:"Optional People API sources value."`
	PageSize     string `json:"pageSize,omitempty" jsonschema:"Optional page size."`
	PageToken    string `json:"pageToken,omitempty" jsonschema:"Optional page token."`
}

func addGWSGmailTools(server *mcp.Server, tools map[string]bool, gc GWSClient, readOnly *mcp.ToolAnnotations, mutating *mcp.ToolAnnotations, destructive *mcp.ToolAnnotations) {
	addGWSTool(server, tools, "gmail", gwsTool("gws_gmail_messages_list", gwsGmailMessagesListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsGmailUserToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.ListGmailMessages(ctx, input.UserID, gmailListParams(input))
		return nil, out, err
	})
	addGWSTool(server, tools, "gmail", gwsTool("gws_gmail_message_get", gwsGmailMessageGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsGmailMessageToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetGmailMessage(ctx, input.UserID, input.MessageID, optionalParam("format", input.Format))
		return nil, out, err
	})
	addGWSTool(server, tools, "gmail", gwsTool("gws_gmail_message_send", gwsGmailMessageSendDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsGmailMessageToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.SendGmailMessage(ctx, input.UserID, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "gmail", gwsTool("gws_gmail_drafts_list", gwsGmailDraftsListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsGmailUserToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.ListGmailDrafts(ctx, input.UserID, gmailDraftListParams(input))
		return nil, out, err
	})
	addGWSTool(server, tools, "gmail", gwsTool("gws_gmail_draft_get", gwsGmailDraftGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsGmailMessageToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetGmailDraft(ctx, input.UserID, input.DraftID, optionalParam("format", input.Format))
		return nil, out, err
	})
	addGWSTool(server, tools, "gmail", gwsTool("gws_gmail_draft_create", gwsGmailDraftCreateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsGmailMessageToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.CreateGmailDraft(ctx, input.UserID, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "gmail", gwsTool("gws_gmail_draft_send", gwsGmailDraftSendDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsGmailMessageToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.SendGmailDraft(ctx, input.UserID, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "gmail", gwsTool("gws_gmail_draft_delete", gwsGmailDraftDeleteDoc, destructive), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsGmailMessageToolInput) (*mcp.CallToolResult, map[string]any, error) {
		err := gc.DeleteGmailDraft(ctx, input.UserID, input.DraftID)
		return nil, map[string]any{"deleted": err == nil, "userId": input.UserID, "draftId": input.DraftID}, err
	})
}

func addGWSCalendarTools(server *mcp.Server, tools map[string]bool, gc GWSClient, readOnly *mcp.ToolAnnotations, mutating *mcp.ToolAnnotations, destructive *mcp.ToolAnnotations) {
	addGWSTool(server, tools, "calendar", gwsTool("gws_calendar_calendar_list_list", gwsCalendarListListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsCalendarListToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.ListCalendarList(ctx, calendarListParams(input))
		return nil, out, err
	})
	addGWSTool(server, tools, "calendar", gwsTool("gws_calendar_calendar_list_get", gwsCalendarListGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsCalendarListToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetCalendarListEntry(ctx, input.CalendarID)
		return nil, out, err
	})
	addGWSTool(server, tools, "calendar", gwsTool("gws_calendar_events_list", gwsCalendarEventsListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsCalendarEventToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.ListCalendarEvents(ctx, input.CalendarID, calendarEventsParams(input))
		return nil, out, err
	})
	addGWSTool(server, tools, "calendar", gwsTool("gws_calendar_event_get", gwsCalendarEventGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsCalendarEventToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetCalendarEvent(ctx, input.CalendarID, input.EventID)
		return nil, out, err
	})
	addGWSTool(server, tools, "calendar", gwsTool("gws_calendar_event_insert", gwsCalendarEventInsertDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsCalendarEventToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.InsertCalendarEvent(ctx, input.CalendarID, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "calendar", gwsTool("gws_calendar_event_patch", gwsCalendarEventPatchDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsCalendarEventToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.PatchCalendarEvent(ctx, input.CalendarID, input.EventID, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "calendar", gwsTool("gws_calendar_event_delete", gwsCalendarEventDeleteDoc, destructive), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsCalendarEventToolInput) (*mcp.CallToolResult, map[string]any, error) {
		err := gc.DeleteCalendarEvent(ctx, input.CalendarID, input.EventID)
		return nil, map[string]any{"deleted": err == nil, "calendarId": input.CalendarID, "eventId": input.EventID}, err
	})
	addGWSTool(server, tools, "calendar", gwsTool("gws_calendar_freebusy_query", gwsCalendarFreeBusyQueryDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsRequestBodyToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.QueryCalendarFreeBusy(ctx, input.Request)
		return nil, out, err
	})
}

func addGWSChatTools(server *mcp.Server, tools map[string]bool, gc GWSClient, readOnly *mcp.ToolAnnotations, mutating *mcp.ToolAnnotations) {
	addGWSTool(server, tools, "chat", gwsTool("gws_chat_spaces_list", gwsChatSpacesListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsChatListToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.ListChatSpaces(ctx, chatListParams(input))
		return nil, out, err
	})
	addGWSTool(server, tools, "chat", gwsTool("gws_chat_space_get", gwsChatSpaceGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsChatListToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetChatSpace(ctx, input.SpaceName)
		return nil, out, err
	})
	addGWSTool(server, tools, "chat", gwsTool("gws_chat_messages_list", gwsChatMessagesListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsChatListToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.ListChatMessages(ctx, input.SpaceName, chatListParams(input))
		return nil, out, err
	})
	addGWSTool(server, tools, "chat", gwsTool("gws_chat_message_get", gwsChatMessageGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsChatListToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetChatMessage(ctx, input.MessageName)
		return nil, out, err
	})
	addGWSTool(server, tools, "chat", gwsTool("gws_chat_message_create", gwsChatMessageCreateDoc, mutating), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsChatListToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.CreateChatMessage(ctx, input.SpaceName, input.Request)
		return nil, out, err
	})
	addGWSTool(server, tools, "chat", gwsTool("gws_chat_members_list", gwsChatMembersListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsChatListToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.ListChatMembers(ctx, input.SpaceName, chatListParams(input))
		return nil, out, err
	})
}

func addGWSPeopleTools(server *mcp.Server, tools map[string]bool, gc GWSClient, readOnly *mcp.ToolAnnotations) {
	addGWSTool(server, tools, "people", gwsTool("gws_people_person_get", gwsPeoplePersonGetDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsPeopleToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.GetPeoplePerson(ctx, input.ResourceName, optionalParam("personFields", input.PersonFields))
		return nil, out, err
	})
	addGWSTool(server, tools, "people", gwsTool("gws_people_connections_list", gwsPeopleConnectionsListDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsPeopleToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.ListPeopleConnections(ctx, input.ResourceName, peopleConnectionsParams(input))
		return nil, out, err
	})
	addGWSTool(server, tools, "people", gwsTool("gws_people_search_contacts", gwsPeopleSearchContactsDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsPeopleToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.SearchPeopleContacts(ctx, peopleSearchParams(input, false))
		return nil, out, err
	})
	addGWSTool(server, tools, "people", gwsTool("gws_people_search_directory", gwsPeopleSearchDirectoryDoc, readOnly), func(ctx context.Context, _ *mcp.CallToolRequest, input *gwsPeopleToolInput) (*mcp.CallToolResult, GWSRawResult, error) {
		out, err := gc.SearchPeopleDirectory(ctx, peopleSearchParams(input, true))
		return nil, out, err
	})
}

func optionalParam(name string, value string) map[string]any {
	if value == "" {
		return nil
	}
	return map[string]any{name: value}
}

func gmailListParams(input *gwsGmailUserToolInput) map[string]any {
	params := gmailDraftListParams(input)
	if input.Query != "" {
		params["q"] = input.Query
	}
	if input.LabelIDs != "" {
		params["labelIds"] = input.LabelIDs
	}
	return params
}

func gmailDraftListParams(input *gwsGmailUserToolInput) map[string]any {
	params := map[string]any{}
	if input.MaxResults != "" {
		params["maxResults"] = input.MaxResults
	}
	if input.PageToken != "" {
		params["pageToken"] = input.PageToken
	}
	return params
}

func calendarListParams(input *gwsCalendarListToolInput) map[string]any {
	params := map[string]any{}
	if input.MaxResults != "" {
		params["maxResults"] = input.MaxResults
	}
	if input.PageToken != "" {
		params["pageToken"] = input.PageToken
	}
	return params
}

func calendarEventsParams(input *gwsCalendarEventToolInput) map[string]any {
	params := map[string]any{}
	if input.TimeMin != "" {
		params["timeMin"] = input.TimeMin
	}
	if input.TimeMax != "" {
		params["timeMax"] = input.TimeMax
	}
	if input.MaxResults != "" {
		params["maxResults"] = input.MaxResults
	}
	if input.SingleEvents != "" {
		params["singleEvents"] = input.SingleEvents
	}
	if input.OrderBy != "" {
		params["orderBy"] = input.OrderBy
	}
	if input.PageToken != "" {
		params["pageToken"] = input.PageToken
	}
	return params
}

func chatListParams(input *gwsChatListToolInput) map[string]any {
	params := map[string]any{}
	if input.PageSize != "" {
		params["pageSize"] = input.PageSize
	}
	if input.PageToken != "" {
		params["pageToken"] = input.PageToken
	}
	return params
}

func peopleConnectionsParams(input *gwsPeopleToolInput) map[string]any {
	params := optionalParam("personFields", input.PersonFields)
	if params == nil {
		params = map[string]any{}
	}
	if input.PageSize != "" {
		params["pageSize"] = input.PageSize
	}
	if input.PageToken != "" {
		params["pageToken"] = input.PageToken
	}
	return params
}

func peopleSearchParams(input *gwsPeopleToolInput, includeSources bool) map[string]any {
	params := map[string]any{}
	if input.Query != "" {
		params["query"] = input.Query
	}
	if input.ReadMask != "" {
		params["readMask"] = input.ReadMask
	}
	if includeSources && input.Sources != "" {
		params["sources"] = input.Sources
	}
	if input.PageSize != "" {
		params["pageSize"] = input.PageSize
	}
	return params
}
