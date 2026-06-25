package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (gc GWSClient) ListGmailMessages(ctx context.Context, userID string, params map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("user-id", userID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIGmail, "/users/"+url.PathEscape(userID)+"/messages", params)
}

func (gc GWSClient) GetGmailMessage(ctx context.Context, userID string, messageID string, params map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("user-id", userID); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("message-id", messageID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIGmail, "/users/"+url.PathEscape(userID)+"/messages/"+url.PathEscape(messageID), params)
}

func (gc GWSClient) SendGmailMessage(ctx context.Context, userID string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("user-id", userID); err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GWSAPIGmail, "/users/"+url.PathEscape(userID)+"/messages/send", body, nil)
}

func (gc GWSClient) ListGmailDrafts(ctx context.Context, userID string, params map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("user-id", userID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIGmail, "/users/"+url.PathEscape(userID)+"/drafts", params)
}

func (gc GWSClient) GetGmailDraft(ctx context.Context, userID string, draftID string, params map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("user-id", userID); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("draft-id", draftID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIGmail, "/users/"+url.PathEscape(userID)+"/drafts/"+url.PathEscape(draftID), params)
}

func (gc GWSClient) CreateGmailDraft(ctx context.Context, userID string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("user-id", userID); err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GWSAPIGmail, "/users/"+url.PathEscape(userID)+"/drafts", body, nil)
}

func (gc GWSClient) SendGmailDraft(ctx context.Context, userID string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("user-id", userID); err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GWSAPIGmail, "/users/"+url.PathEscape(userID)+"/drafts/send", body, nil)
}

func (gc GWSClient) DeleteGmailDraft(ctx context.Context, userID string, draftID string) error {
	if err := requireNonEmpty("user-id", userID); err != nil {
		return err
	}
	if err := requireNonEmpty("draft-id", draftID); err != nil {
		return err
	}
	_, err := gc.RequestContext(ctx, GWSRequest{API: string(GWSAPIGmail), Method: http.MethodDelete, Path: "/users/" + url.PathEscape(userID) + "/drafts/" + url.PathEscape(draftID)})
	return err
}

func (gc GWSClient) ListCalendarList(ctx context.Context, params map[string]any) (GWSRawResult, error) {
	return gc.getJSON(ctx, GWSAPICalendar, "/users/me/calendarList", params)
}

func (gc GWSClient) GetCalendarListEntry(ctx context.Context, calendarID string) (GWSRawResult, error) {
	if err := requireNonEmpty("calendar-id", calendarID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPICalendar, "/users/me/calendarList/"+url.PathEscape(calendarID), nil)
}

func (gc GWSClient) ListCalendarEvents(ctx context.Context, calendarID string, params map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("calendar-id", calendarID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPICalendar, "/calendars/"+url.PathEscape(calendarID)+"/events", params)
}

func (gc GWSClient) GetCalendarEvent(ctx context.Context, calendarID string, eventID string) (GWSRawResult, error) {
	if err := requireNonEmpty("calendar-id", calendarID); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("event-id", eventID); err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPICalendar, "/calendars/"+url.PathEscape(calendarID)+"/events/"+url.PathEscape(eventID), nil)
}

func (gc GWSClient) InsertCalendarEvent(ctx context.Context, calendarID string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("calendar-id", calendarID); err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GWSAPICalendar, "/calendars/"+url.PathEscape(calendarID)+"/events", body, nil)
}

func (gc GWSClient) PatchCalendarEvent(ctx context.Context, calendarID string, eventID string, body map[string]any) (GWSRawResult, error) {
	if err := requireNonEmpty("calendar-id", calendarID); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("event-id", eventID); err != nil {
		return nil, err
	}
	var out GWSRawResult
	err := gc.patchTypedJSON(ctx, GWSAPICalendar, "/calendars/"+url.PathEscape(calendarID)+"/events/"+url.PathEscape(eventID), body, nil, &out)
	return out, err
}

func (gc GWSClient) DeleteCalendarEvent(ctx context.Context, calendarID string, eventID string) error {
	if err := requireNonEmpty("calendar-id", calendarID); err != nil {
		return err
	}
	if err := requireNonEmpty("event-id", eventID); err != nil {
		return err
	}
	_, err := gc.RequestContext(ctx, GWSRequest{API: string(GWSAPICalendar), Method: http.MethodDelete, Path: "/calendars/" + url.PathEscape(calendarID) + "/events/" + url.PathEscape(eventID)})
	return err
}

func (gc GWSClient) QueryCalendarFreeBusy(ctx context.Context, body map[string]any) (GWSRawResult, error) {
	return gc.postJSON(ctx, GWSAPICalendar, "/freeBusy", body, nil)
}

func (gc GWSClient) ListChatSpaces(ctx context.Context, params map[string]any) (GWSRawResult, error) {
	return gc.getJSON(ctx, GWSAPIChat, "/spaces", params)
}

func (gc GWSClient) GetChatSpace(ctx context.Context, spaceName string) (GWSRawResult, error) {
	path, err := googleResourcePath("space-name", spaceName, "spaces/")
	if err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIChat, "/"+path, nil)
}

func (gc GWSClient) ListChatMessages(ctx context.Context, spaceName string, params map[string]any) (GWSRawResult, error) {
	path, err := googleResourcePath("space-name", spaceName, "spaces/")
	if err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIChat, "/"+path+"/messages", params)
}

func (gc GWSClient) GetChatMessage(ctx context.Context, messageName string) (GWSRawResult, error) {
	path, err := googleResourcePath("message-name", messageName, "spaces/")
	if err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIChat, "/"+path, nil)
}

func (gc GWSClient) CreateChatMessage(ctx context.Context, spaceName string, body map[string]any) (GWSRawResult, error) {
	path, err := googleResourcePath("space-name", spaceName, "spaces/")
	if err != nil {
		return nil, err
	}
	return gc.postJSON(ctx, GWSAPIChat, "/"+path+"/messages", body, nil)
}

func (gc GWSClient) ListChatMembers(ctx context.Context, spaceName string, params map[string]any) (GWSRawResult, error) {
	path, err := googleResourcePath("space-name", spaceName, "spaces/")
	if err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIChat, "/"+path+"/members", params)
}

func (gc GWSClient) GetPeoplePerson(ctx context.Context, resourceName string, params map[string]any) (GWSRawResult, error) {
	path, err := googleResourcePath("resource-name", resourceName, "people/")
	if err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIPeople, "/"+path, params)
}

func (gc GWSClient) ListPeopleConnections(ctx context.Context, resourceName string, params map[string]any) (GWSRawResult, error) {
	path, err := googleResourcePath("resource-name", resourceName, "people/")
	if err != nil {
		return nil, err
	}
	return gc.getJSON(ctx, GWSAPIPeople, "/"+path+"/connections", params)
}

func (gc GWSClient) SearchPeopleContacts(ctx context.Context, params map[string]any) (GWSRawResult, error) {
	return gc.getJSON(ctx, GWSAPIPeople, "/people:searchContacts", params)
}

func (gc GWSClient) SearchPeopleDirectory(ctx context.Context, params map[string]any) (GWSRawResult, error) {
	return gc.getJSON(ctx, GWSAPIPeople, "/people:searchDirectoryPeople", params)
}

func googleResourcePath(label string, value string, prefix string) (string, error) {
	trimmed := strings.Trim(value, "/")
	if err := requireNonEmpty(label, trimmed); err != nil {
		return "", err
	}
	if !strings.HasPrefix(trimmed, prefix) {
		return "", fmt.Errorf("%s must start with %s", label, prefix)
	}
	if strings.ContainsAny(trimmed, "?#") {
		return "", fmt.Errorf("%s must not contain query or fragment characters", label)
	}
	return trimmed, nil
}
