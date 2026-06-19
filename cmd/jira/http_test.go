package main

import (
	"encoding/json"
	"testing"
)

func TestJiraStatusRequestPayloadsUseJiraFieldNames(t *testing.T) {
	createPayload, err := json.Marshal(JiraStatusCreateInput{
		Scope: JiraStatusScope{
			Type: "PROJECT",
			Project: &JiraStatusProject{
				ID: "10000",
			},
		},
		Statuses: []JiraStatusCreate{
			{
				Name:           "Ready for Review",
				Description:    "Ready for review by another engineer.",
				StatusCategory: "IN_PROGRESS",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var createDecoded struct {
		Scope struct {
			Type    string `json:"type"`
			Project struct {
				ID string `json:"id"`
			} `json:"project"`
		} `json:"scope"`
		Statuses []struct {
			Name           string `json:"name"`
			Description    string `json:"description"`
			StatusCategory string `json:"statusCategory"`
		} `json:"statuses"`
	}
	if err := json.Unmarshal(createPayload, &createDecoded); err != nil {
		t.Fatal(err)
	}
	if createDecoded.Scope.Type != "PROJECT" || createDecoded.Scope.Project.ID != "10000" {
		t.Fatalf("unexpected create scope payload: %s", string(createPayload))
	}
	if len(createDecoded.Statuses) != 1 {
		t.Fatalf("expected one status create payload, got %s", string(createPayload))
	}
	createStatus := createDecoded.Statuses[0]
	if createStatus.Name != "Ready for Review" || createStatus.Description != "Ready for review by another engineer." || createStatus.StatusCategory != "IN_PROGRESS" {
		t.Fatalf("unexpected create status payload: %#v", createStatus)
	}

	updatePayload, err := json.Marshal(JiraStatusUpdateInput{
		Statuses: []JiraStatusUpdate{
			{
				ID:             "10001",
				Name:           stringPointer("Reviewing"),
				StatusCategory: stringPointer("IN_PROGRESS"),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var updateDecoded struct {
		Statuses []struct {
			ID             string  `json:"id"`
			Name           *string `json:"name"`
			Description    *string `json:"description"`
			StatusCategory *string `json:"statusCategory"`
		} `json:"statuses"`
	}
	if err := json.Unmarshal(updatePayload, &updateDecoded); err != nil {
		t.Fatal(err)
	}
	if len(updateDecoded.Statuses) != 1 {
		t.Fatalf("expected one status update payload, got %s", string(updatePayload))
	}
	updateStatus := updateDecoded.Statuses[0]
	if updateStatus.ID != "10001" || updateStatus.Name == nil || *updateStatus.Name != "Reviewing" {
		t.Fatalf("unexpected update status payload: %#v", updateStatus)
	}
	if updateStatus.StatusCategory == nil || *updateStatus.StatusCategory != "IN_PROGRESS" {
		t.Fatalf("unexpected update status category payload: %#v", updateStatus)
	}
	if updateStatus.Description != nil {
		t.Fatalf("expected omitted description to stay omitted, got %#v", updateStatus.Description)
	}
}

func TestJiraStatusUpdatePayloadCanClearDescription(t *testing.T) {
	clearDescription := ""
	updateInput, statusIDs, err := buildJiraMCPStatusUpdateInput(&jiraStatusUpdateToolInput{
		Statuses: []JiraStatusUpdate{
			{
				ID:          "10001",
				Description: &clearDescription,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(statusIDs) != 1 || statusIDs[0] != "10001" {
		t.Fatalf("expected status id 10001, got %#v", statusIDs)
	}

	updatePayload, err := json.Marshal(updateInput)
	if err != nil {
		t.Fatal(err)
	}

	var decoded struct {
		Statuses []struct {
			Description *string `json:"description"`
		} `json:"statuses"`
	}
	if err := json.Unmarshal(updatePayload, &decoded); err != nil {
		t.Fatal(err)
	}

	if len(decoded.Statuses) != 1 {
		t.Fatalf("expected one status update, got %#v", decoded.Statuses)
	}
	if decoded.Statuses[0].Description == nil {
		t.Fatalf("expected payload %s to include an explicit description field", string(updatePayload))
	}
	if *decoded.Statuses[0].Description != "" {
		t.Fatalf("expected empty description, got %q", *decoded.Statuses[0].Description)
	}
}

func TestJiraStatusCreateResponseUsesTopLevelArray(t *testing.T) {
	var statuses []JiraStatus
	err := json.Unmarshal([]byte(`[
		{
			"id": "10001",
			"name": "Ready for Review",
			"description": "Ready for review by another engineer.",
			"statusCategory": "IN_PROGRESS",
			"scope": {
				"type": "PROJECT",
				"project": {
					"id": "10000"
				}
			}
		}
	]`), &statuses)
	if err != nil {
		t.Fatal(err)
	}

	if len(statuses) != 1 {
		t.Fatalf("expected one status, got %#v", statuses)
	}

	status := statuses[0]
	if status.ID != "10001" || status.Name != "Ready for Review" {
		t.Fatalf("unexpected status response: %#v", status)
	}
	if status.Scope.Type != "PROJECT" || status.Scope.Project == nil || status.Scope.Project.ID != "10000" {
		t.Fatalf("unexpected status scope: %#v", status.Scope)
	}
}

func stringPointer(value string) *string {
	return &value
}
