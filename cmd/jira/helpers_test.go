package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mistlehq/tools/internal/testproxy"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func getRequiredEnv(t *testing.T, name string) string {
	t.Helper()

	value := os.Getenv(name)
	if value == "" {
		t.Skipf("skipping: %s is not set", name)
	}

	return value
}

type commandResult struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func runCommandWithInput(t *testing.T, env Environment, input string, args ...string) (commandResult, error) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cli := CLI{
		stdin:  bytes.NewBufferString(input),
		stdout: &stdout,
		stderr: &stderr,
		env:    env,
	}

	err := cli.run(args)
	return commandResult{
		stdout: stdout,
		stderr: stderr,
	}, err
}

func setupAndRunCommandWithInput(t *testing.T, input string, args ...string) commandResult {
	env := setupCommandEnvironment(t)

	commandResult, err := runCommandWithInput(t, env, input, args...)
	if err != nil {
		t.Fatal(err)
	}

	return commandResult
}

func setupCommandEnvironment(t *testing.T) Environment {
	t.Helper()

	upstreamBaseURL := getRequiredEnv(t, "JIRA_TEST_UPSTREAM_BASE_URL")
	username := getRequiredEnv(t, "JIRA_TEST_USERNAME")
	password := getRequiredEnv(t, "JIRA_TEST_PASSWORD")

	proxy, err := testproxy.Start(testproxy.Config{
		UpstreamBaseURL: upstreamBaseURL,
		AuthMode:        testproxy.AuthModeBasic,
		Username:        username,
		Password:        password,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := proxy.Close(); err != nil {
			t.Errorf("failed to close proxy: %v", err)
		}
	})

	return Environment{
		"JIRA_BASE_URL": proxy.BaseURL,
	}
}

// The Jira test tenant is currently fixed, so tests derive project and issue
// type from a stable seed issue and then create isolated per-test issues from
// that template.
const jiraTestTemplateIssueKey = "KAN-1"

type jiraTestIssueTemplate struct {
	Fields struct {
		Project struct {
			ID string `json:"id"`
		} `json:"project"`
		IssueType struct {
			ID string `json:"id"`
		} `json:"issuetype"`
	} `json:"fields"`
}

type jiraCreatedIssue struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

func setupIsolatedIssue(t *testing.T) (Environment, string) {
	t.Helper()

	env := setupCommandEnvironment(t)
	config, err := loadConfig(env)
	if err != nil {
		t.Fatal(err)
	}

	jc := NewJiraClient(config)
	template, err := getJiraTestIssueTemplate(jc, jiraTestTemplateIssueKey)
	if err != nil {
		t.Fatal(err)
	}

	issue, err := createJiraTestIssue(jc, template, fmt.Sprintf("integration test %s %d", t.Name(), time.Now().UnixNano()))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := deleteJiraTestIssue(jc, issue.Key); err != nil {
			t.Errorf("failed to delete issue %s: %v", issue.Key, err)
		}
	})

	return env, issue.Key
}

func getJiraTestIssueTemplate(jc JiraClient, issueKey string) (jiraTestIssueTemplate, error) {
	responseBody, err := jc.get(fmt.Sprintf("/rest/api/3/issue/%s?fields=project,issuetype", issueKey))
	if err != nil {
		return jiraTestIssueTemplate{}, err
	}

	var template jiraTestIssueTemplate
	if err := json.Unmarshal(responseBody, &template); err != nil {
		return jiraTestIssueTemplate{}, err
	}

	return template, nil
}

func createJiraTestIssue(jc JiraClient, template jiraTestIssueTemplate, summary string) (jiraCreatedIssue, error) {
	requestBody, err := json.Marshal(struct {
		Fields map[string]any `json:"fields"`
	}{
		Fields: map[string]any{
			"project": map[string]string{
				"id": template.Fields.Project.ID,
			},
			"issuetype": map[string]string{
				"id": template.Fields.IssueType.ID,
			},
			"summary": summary,
		},
	})
	if err != nil {
		return jiraCreatedIssue{}, err
	}

	responseBody, err := jc.post("/rest/api/3/issue", requestBody)
	if err != nil {
		return jiraCreatedIssue{}, err
	}

	var issue jiraCreatedIssue
	if err := json.Unmarshal(responseBody, &issue); err != nil {
		return jiraCreatedIssue{}, err
	}

	return issue, nil
}

func deleteJiraTestIssue(jc JiraClient, issueKey string) error {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		jc.baseURL+fmt.Sprintf("/rest/api/3/issue/%s", issueKey),
		nil,
	)
	if err != nil {
		return err
	}

	response, err := jc.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("request failed with status %d: %s", response.StatusCode, string(body))
	}

	return nil
}
