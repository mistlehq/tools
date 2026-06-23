package main

import "testing"

func TestShopGet(t *testing.T) {
	env := setupCommandEnvironment(t)
	commandResult, err := runCommandWithInput(t, env, "", "shopify", "shop", "get")
	if err != nil {
		t.Fatal(err)
	}

	var shop ShopifyShop
	decodeCommandJSON(t, commandResult, &shop)
	if shop.ID == "" || shop.Name == "" {
		t.Fatalf("expected shop details, got %#v", shop)
	}
}

func TestProductsReadCommands(t *testing.T) {
	env := setupCommandEnvironment(t)

	getByIDResult, err := runCommandWithInput(t, env, "", "shopify", "products", "get", "--id", testProductID(t))
	if err != nil {
		t.Fatal(err)
	}
	var productByID ShopifyProduct
	decodeCommandJSON(t, getByIDResult, &productByID)
	if productByID.ID != testProductID(t) {
		t.Fatalf("expected product %q, got %#v", testProductID(t), productByID)
	}

	getByHandleResult, err := runCommandWithInput(t, env, "", "shopify", "products", "get", "--handle", testProductHandle(t))
	if err != nil {
		t.Fatal(err)
	}
	var productByHandle ShopifyProduct
	decodeCommandJSON(t, getByHandleResult, &productByHandle)
	if productByHandle.Handle != testProductHandle(t) {
		t.Fatalf("expected product handle %q, got %#v", testProductHandle(t), productByHandle)
	}

	searchResult, err := runCommandWithInput(t, env, "", "shopify", "products", "search", "--first", "5", "--query", "handle:"+testProductHandle(t))
	if err != nil {
		t.Fatal(err)
	}
	var products ShopifyProductsSearch
	decodeCommandJSON(t, searchResult, &products)
	if len(products.Products.Nodes) == 0 {
		t.Fatal("expected product search results")
	}
}

func TestOrdersReadCommands(t *testing.T) {
	env := setupCommandEnvironment(t)

	getResult, err := runCommandWithInput(t, env, "", "shopify", "orders", "get", "--id", testOrderID(t))
	if err != nil {
		if isProtectedCustomerDataError(err) {
			return
		}
		t.Fatal(err)
	}
	var order ShopifyOrder
	decodeCommandJSON(t, getResult, &order)
	if order.ID != testOrderID(t) {
		t.Fatalf("expected order %q, got %#v", testOrderID(t), order)
	}

	searchResult, err := runCommandWithInput(t, env, "", "shopify", "orders", "search", "--first", "5", "--query", "name:"+testOrderName(t))
	if err != nil {
		if isProtectedCustomerDataError(err) {
			return
		}
		t.Fatal(err)
	}
	var orders ShopifyOrdersSearch
	decodeCommandJSON(t, searchResult, &orders)
	if len(orders.Orders.Nodes) == 0 {
		t.Fatal("expected order search results")
	}
}

func TestCustomersReadCommands(t *testing.T) {
	env := setupCommandEnvironment(t)

	getResult, err := runCommandWithInput(t, env, "", "shopify", "customers", "get", "--id", testCustomerID(t))
	if err != nil {
		if isProtectedCustomerDataError(err) {
			return
		}
		t.Fatal(err)
	}
	var customer ShopifyCustomer
	decodeCommandJSON(t, getResult, &customer)
	if customer.ID != testCustomerID(t) {
		t.Fatalf("expected customer %q, got %#v", testCustomerID(t), customer)
	}

	searchResult, err := runCommandWithInput(t, env, "", "shopify", "customers", "search", "--first", "5", "--query", "email:"+testCustomerEmail(t))
	if err != nil {
		if isProtectedCustomerDataError(err) {
			return
		}
		t.Fatal(err)
	}
	var customers ShopifyCustomersSearch
	decodeCommandJSON(t, searchResult, &customers)
	if len(customers.Customers.Nodes) == 0 {
		t.Fatal("expected customer search results")
	}
}

func TestInventoryAndLocationsReadCommands(t *testing.T) {
	env := setupCommandEnvironment(t)

	itemsResult, err := runCommandWithInput(t, env, "", "shopify", "inventory", "items", "search", "--first", "5")
	if err != nil {
		t.Fatal(err)
	}
	var items ShopifyInventoryItemsSearch
	decodeCommandJSON(t, itemsResult, &items)
	if len(items.InventoryItems.Nodes) == 0 {
		t.Fatal("expected inventory item search results")
	}

	levelsResult, err := runCommandWithInput(t, env, "", "shopify", "inventory", "levels", "search", "--first", "5")
	if err != nil {
		t.Fatal(err)
	}
	var levels ShopifyInventoryLevelsSearch
	decodeCommandJSON(t, levelsResult, &levels)
	if len(levels.InventoryItems.Nodes) == 0 {
		t.Fatal("expected inventory items with inventory levels")
	}

	locationsResult, err := runCommandWithInput(t, env, "", "shopify", "locations", "list", "--first", "5")
	if err != nil {
		t.Fatal(err)
	}
	var locations ShopifyLocationsList
	decodeCommandJSON(t, locationsResult, &locations)
	if len(locations.Locations.Nodes) == 0 {
		t.Fatal("expected locations")
	}
}
