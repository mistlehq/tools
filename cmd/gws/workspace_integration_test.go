package main

import (
	"fmt"
	"testing"
)

func TestGmailCommandsWhenConfigured(t *testing.T) {
	userID := getRequiredEnv(t, "GWS_TEST_GMAIL_USER_ID")
	env := setupCommandEnvironment(t)

	listResult, err := runCommandWithInput(t, env, "", "gws", "gmail", "messages", "list", "--user-id", userID, "--max-results", "1")
	if err != nil {
		t.Fatal(err)
	}
	var list map[string]any
	decodeCommandJSON(t, listResult, &list)
	if list["resultSizeEstimate"] == nil && list["messages"] == nil {
		t.Fatalf("expected Gmail message list shape, got %#v", list)
	}

	draftsResult, err := runCommandWithInput(t, env, "", "gws", "gmail", "drafts", "list", "--user-id", userID, "--max-results", "1")
	if err != nil {
		t.Fatal(err)
	}
	var drafts map[string]any
	decodeCommandJSON(t, draftsResult, &drafts)
	if drafts["resultSizeEstimate"] == nil && drafts["drafts"] == nil {
		t.Fatalf("expected Gmail draft list shape, got %#v", drafts)
	}
}

func TestCalendarCommandsWhenConfigured(t *testing.T) {
	calendarID := getRequiredEnv(t, "GWS_TEST_CALENDAR_ID")
	env := setupCommandEnvironment(t)

	listResult, err := runCommandWithInput(t, env, "", "gws", "calendar", "calendar-list", "list", "--max-results", "10")
	if err != nil {
		t.Fatal(err)
	}
	var calendars map[string]any
	decodeCommandJSON(t, listResult, &calendars)
	if calendars["items"] == nil {
		t.Fatalf("expected calendar list items, got %#v", calendars)
	}

	eventsResult, err := runCommandWithInput(t, env, "", "gws", "calendar", "events", "list", "--calendar-id", calendarID, "--max-results", "1")
	if err != nil {
		t.Fatal(err)
	}
	var events map[string]any
	decodeCommandJSON(t, eventsResult, &events)
	if events["items"] == nil {
		t.Fatalf("expected calendar event items, got %#v", events)
	}

	freeBusyPath := writeTempJSONRequest(t, fmt.Sprintf(`{"timeMin":"2026-01-01T00:00:00Z","timeMax":"2026-01-02T00:00:00Z","items":[{"id":%q}]}`, calendarID))
	freeBusyResult, err := runCommandWithInput(t, env, "", "gws", "calendar", "freebusy", "query", "--request-file", freeBusyPath)
	if err != nil {
		t.Fatal(err)
	}
	var freeBusy map[string]any
	decodeCommandJSON(t, freeBusyResult, &freeBusy)
	if freeBusy["calendars"] == nil {
		t.Fatalf("expected freebusy calendars, got %#v", freeBusy)
	}
}

func TestChatCommandsWhenConfigured(t *testing.T) {
	spaceName := getRequiredEnv(t, "GWS_TEST_CHAT_SPACE_NAME")
	env := setupCommandEnvironment(t)

	listResult, err := runCommandWithInput(t, env, "", "gws", "chat", "spaces", "list", "--page-size", "10")
	if err != nil {
		t.Fatal(err)
	}
	var spaces map[string]any
	decodeCommandJSON(t, listResult, &spaces)
	if spaces["spaces"] == nil {
		t.Fatalf("expected chat spaces, got %#v", spaces)
	}

	membersResult, err := runCommandWithInput(t, env, "", "gws", "chat", "members", "list", "--space-name", spaceName, "--page-size", "10")
	if err != nil {
		t.Fatal(err)
	}
	var members map[string]any
	decodeCommandJSON(t, membersResult, &members)
	if members["memberships"] == nil {
		t.Fatalf("expected chat memberships, got %#v", members)
	}
}

func TestPeopleCommandsWhenConfigured(t *testing.T) {
	resourceName := getRequiredEnv(t, "GWS_TEST_PEOPLE_RESOURCE_NAME")
	env := setupCommandEnvironment(t)

	personResult, err := runCommandWithInput(t, env, "", "gws", "people", "people", "get", "--resource-name", resourceName, "--person-fields", "names,emailAddresses")
	if err != nil {
		t.Fatal(err)
	}
	var person map[string]any
	decodeCommandJSON(t, personResult, &person)
	if person["resourceName"] == nil {
		t.Fatalf("expected person resource, got %#v", person)
	}

	connectionsResult, err := runCommandWithInput(t, env, "", "gws", "people", "connections", "list", "--resource-name", resourceName, "--person-fields", "names,emailAddresses", "--page-size", "10")
	if err != nil {
		t.Fatal(err)
	}
	var connections map[string]any
	decodeCommandJSON(t, connectionsResult, &connections)
	if connections["connections"] == nil && connections["totalPeople"] == nil {
		t.Fatalf("expected people connections shape, got %#v", connections)
	}
}
