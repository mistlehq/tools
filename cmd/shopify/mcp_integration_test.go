package main

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "shopify", "mcp", "help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	for _, want := range []string{"shopify mcp", "shopify mcp serve", "Streamable HTTP"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp help to mention %q", want)
		}
	}
}

func TestMCPServeHelp(t *testing.T) {
	commandResult, err := runCommandWithInput(t, Environment{}, "", "shopify", "mcp", "serve", "--help")
	if err != nil {
		t.Fatal(err)
	}

	output := commandResult.stdout.String()
	for _, want := range []string{"--addr <addr>", "--endpoint <path>", "shopify_graphql_request", "shopify_product_delete"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected mcp serve help to mention %q", want)
		}
	}
}

func TestMCPServerListsShopifyTools(t *testing.T) {
	session := newLocalShopifyMCPTestSession(t)
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
		"shopify_graphql_request":         shopifyGraphQLRequestDoc.Description,
		"shopify_auth_test":               shopifyAuthTestDoc.Description,
		"shopify_shop_get":                shopifyShopGetDoc.Description,
		"shopify_products_search":         shopifyProductsSearchDoc.Description,
		"shopify_product_get":             shopifyProductGetDoc.Description,
		"shopify_product_create":          shopifyProductCreateDoc.Description,
		"shopify_product_update":          shopifyProductUpdateDoc.Description,
		"shopify_product_delete":          shopifyProductDeleteDoc.Description,
		"shopify_orders_search":           shopifyOrdersSearchDoc.Description,
		"shopify_order_get":               shopifyOrderGetDoc.Description,
		"shopify_customers_search":        shopifyCustomersSearchDoc.Description,
		"shopify_customer_get":            shopifyCustomerGetDoc.Description,
		"shopify_inventory_items_search":  shopifyInventoryItemsSearchDoc.Description,
		"shopify_inventory_levels_search": shopifyInventoryLevelsSearchDoc.Description,
		"shopify_locations_list":          shopifyLocationsListDoc.Description,
	}

	for name, description := range expected {
		tool, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected MCP tool %q to be listed", name)
		}
		if tool.Description != description {
			t.Fatalf("expected MCP tool %q description %q, got %q", name, description, tool.Description)
		}
		if name == "shopify_graphql_request" {
			if tool.Annotations == nil || tool.Annotations.OpenWorldHint == nil || !*tool.Annotations.OpenWorldHint {
				t.Fatalf("expected MCP tool %q to be annotated as open-world", name)
			}
			continue
		}
		if name == "shopify_product_delete" {
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || !*tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be annotated as destructive", name)
			}
			continue
		}
		if name == "shopify_product_create" || name == "shopify_product_update" {
			if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || *tool.Annotations.DestructiveHint {
				t.Fatalf("expected MCP tool %q to be annotated as non-destructive mutation", name)
			}
			continue
		}
		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Fatalf("expected MCP tool %q to be annotated as read-only", name)
		}
	}
}

func TestMCPShopifyReadTools(t *testing.T) {
	_, sc := setupShopifyClient(t)
	session := newShopifyMCPTestSession(t, sc)
	defer session.Close()

	rawResult := callShopifyMCPTool(t, session, "shopify_graphql_request", map[string]any{
		"query":     "query ProductByID($id: ID!) { product(id: $id) { id handle } }",
		"variables": map[string]any{"id": testProductID(t)},
	})
	var raw map[string]any
	decodeMCPStructuredContent(t, rawResult, &raw)
	if raw["data"] == nil {
		t.Fatalf("expected raw GraphQL response data, got %#v", raw)
	}

	shopResult := callShopifyMCPTool(t, session, "shopify_shop_get", map[string]any{})
	var shop ShopifyShop
	decodeMCPStructuredContent(t, shopResult, &shop)
	if shop.ID == "" {
		t.Fatalf("expected shop details, got %#v", shop)
	}

	productResult := callShopifyMCPTool(t, session, "shopify_product_get", map[string]any{"id": testProductID(t)})
	var product ShopifyProduct
	decodeMCPStructuredContent(t, productResult, &product)
	if product.ID != testProductID(t) {
		t.Fatalf("expected product %q, got %#v", testProductID(t), product)
	}

	ordersResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "shopify_orders_search",
		Arguments: map[string]any{"first": 5, "query": "name:" + testOrderName(t)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if ordersResult.IsError {
		if !toolErrorContains(ordersResult, "not approved to access") || !toolErrorContains(ordersResult, "protected-customer-data") {
			t.Fatalf("expected protected customer data error, got %#v", ordersResult.Content)
		}
	} else {
		var orders ShopifyOrdersSearch
		decodeMCPStructuredContent(t, ordersResult, &orders)
		if len(orders.Orders.Nodes) == 0 {
			t.Fatal("expected order search results")
		}
	}

	customerResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "shopify_customer_get",
		Arguments: map[string]any{"id": testCustomerID(t)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if customerResult.IsError {
		if !toolErrorContains(customerResult, "not approved to access") || !toolErrorContains(customerResult, "protected-customer-data") {
			t.Fatalf("expected protected customer data error, got %#v", customerResult.Content)
		}
	} else {
		var customer ShopifyCustomer
		decodeMCPStructuredContent(t, customerResult, &customer)
		if customer.ID != testCustomerID(t) {
			t.Fatalf("expected customer %q, got %#v", testCustomerID(t), customer)
		}
	}

	locationsResult := callShopifyMCPTool(t, session, "shopify_locations_list", map[string]any{"first": 5})
	var locations ShopifyLocationsList
	decodeMCPStructuredContent(t, locationsResult, &locations)
	if len(locations.Locations.Nodes) == 0 {
		t.Fatal("expected locations")
	}
}

func TestMCPShopifyToolValidation(t *testing.T) {
	session := newLocalShopifyMCPTestSession(t)
	defer session.Close()

	testCases := []struct {
		name      string
		tool      string
		arguments map[string]any
	}{
		{name: "search missing first", tool: "shopify_products_search", arguments: map[string]any{}},
		{name: "get missing id and handle", tool: "shopify_product_get", arguments: map[string]any{}},
		{name: "delete missing id", tool: "shopify_product_delete", arguments: map[string]any{}},
		{name: "raw missing query", tool: "shopify_graphql_request", arguments: map[string]any{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
				Name:      tc.tool,
				Arguments: tc.arguments,
			})
			if err != nil {
				t.Fatal(err)
			}
			if !result.IsError {
				t.Fatal("expected tool validation to return a tool error")
			}
		})
	}
}
