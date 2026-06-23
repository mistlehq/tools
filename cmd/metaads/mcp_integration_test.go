package main

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHelp(t *testing.T) {
	result, err := runCommandWithInput(t, Environment{}, "", "metaads", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}
	if !stringsContainsAll(result.stdout.String(), []string{"metaads mcp", "metaads mcp serve", "Streamable HTTP"}) {
		t.Fatalf("unexpected mcp help: %s", result.stdout.String())
	}
}

func TestMCPServerListsMetaAdsTools(t *testing.T) {
	session := newLocalMetaAdsMCPTestSession(t)
	defer session.Close()

	toolsResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	toolsByName := make(map[string]*mcp.Tool)
	for _, tool := range toolsResult.Tools {
		toolsByName[tool.Name] = tool
	}

	expected := map[string]string{
		"metaads_graph_request":    metaAdsGraphRequestDoc.Description,
		"metaads_auth_test":        metaAdsAuthTestDoc.Description,
		"metaads_ad_accounts_list": metaAdsAdAccountsListDoc.Description,
		"metaads_ad_account_get":   metaAdsAdAccountGetDoc.Description,
		"metaads_campaigns_search": metaAdsCampaignsSearchDoc.Description,
		"metaads_campaign_get":     metaAdsCampaignGetDoc.Description,
		"metaads_campaign_create":  metaAdsCampaignCreateDoc.Description,
		"metaads_campaign_update":  metaAdsCampaignUpdateDoc.Description,
		"metaads_campaign_delete":  metaAdsCampaignDeleteDoc.Description,
		"metaads_adsets_search":    metaAdsAdSetsSearchDoc.Description,
		"metaads_adset_get":        metaAdsAdSetGetDoc.Description,
		"metaads_adset_create":     metaAdsAdSetCreateDoc.Description,
		"metaads_adset_update":     metaAdsAdSetUpdateDoc.Description,
		"metaads_adset_delete":     metaAdsAdSetDeleteDoc.Description,
		"metaads_ads_search":       metaAdsAdsSearchDoc.Description,
		"metaads_ad_get":           metaAdsAdGetDoc.Description,
		"metaads_ad_create":        metaAdsAdCreateDoc.Description,
		"metaads_ad_update":        metaAdsAdUpdateDoc.Description,
		"metaads_ad_delete":        metaAdsAdDeleteDoc.Description,
		"metaads_creatives_search": metaAdsCreativesSearchDoc.Description,
		"metaads_creative_get":     metaAdsCreativeGetDoc.Description,
		"metaads_creative_create":  metaAdsCreativeCreateDoc.Description,
		"metaads_insights_get":     metaAdsInsightsGetDoc.Description,
		"metaads_targeting_search": metaAdsTargetingSearchDoc.Description,
	}
	for name, description := range expected {
		tool, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected MCP tool %q to be listed", name)
		}
		if tool.Description != description {
			t.Fatalf("expected MCP tool %q description %q, got %q", name, description, tool.Description)
		}
		if name == "metaads_graph_request" {
			if tool.Annotations == nil || tool.Annotations.OpenWorldHint == nil || !*tool.Annotations.OpenWorldHint {
				t.Fatalf("expected MCP tool %q to be open-world", name)
			}
			continue
		}
		if name == "metaads_campaign_delete" || name == "metaads_adset_delete" || name == "metaads_ad_delete" {
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || !*tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be destructive", name)
			}
			continue
		}
	}
}

func TestMCPMetaAdsReadTools(t *testing.T) {
	_, mc := setupMetaAdsClient(t)
	session := newMetaAdsMCPTestSession(t, mc)
	defer session.Close()

	rawResult := callMetaAdsMCPTool(t, session, "metaads_graph_request", map[string]any{
		"method": "GET",
		"path":   "/me",
		"params": map[string]any{"fields": "id,name"},
	})
	var raw map[string]any
	decodeMCPStructuredContent(t, rawResult, &raw)
	if raw["id"] == nil {
		t.Fatalf("expected raw /me response, got %#v", raw)
	}

	accountsResult := callMetaAdsMCPTool(t, session, "metaads_ad_accounts_list", map[string]any{"limit": "5"})
	var accounts map[string]any
	decodeMCPStructuredContent(t, accountsResult, &accounts)
	if accounts["data"] == nil {
		t.Fatalf("expected ad accounts response, got %#v", accounts)
	}
}

func TestMCPMetaAdsToolValidation(t *testing.T) {
	session := newLocalMetaAdsMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{name: "raw missing path", tool: "metaads_graph_request", arguments: map[string]any{}},
		{name: "campaign search missing account", tool: "metaads_campaigns_search", arguments: map[string]any{}},
		{name: "campaign delete missing id", tool: "metaads_campaign_delete", arguments: map[string]any{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: tc.tool, Arguments: tc.arguments})
			if err != nil {
				t.Fatal(err)
			}
			if !result.IsError {
				t.Fatal("expected tool validation to return a tool error")
			}
		})
	}
}

func stringsContainsAll(text string, values []string) bool {
	for _, value := range values {
		if !strings.Contains(text, value) {
			return false
		}
	}
	return true
}
