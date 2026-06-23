package main

import (
	"encoding/json"
	"testing"
)

func TestProductWriteCommands(t *testing.T) {
	env, sc := setupShopifyClient(t)

	title := uniqueProductTitle(t)
	createResult, err := runCommandWithInput(t, env, "", "shopify", "products", "create", "--product-json", `{"title": `+quoteJSONString(t, title)+`}`)
	if err != nil {
		t.Fatal(err)
	}
	var created ShopifyProductCreate
	decodeCommandJSON(t, createResult, &created)
	if len(created.ProductCreate.UserErrors) > 0 {
		t.Fatalf("expected product create without user errors, got %#v", created.ProductCreate.UserErrors)
	}
	if created.ProductCreate.Product == nil || created.ProductCreate.Product.ID == "" {
		t.Fatalf("expected created product, got %#v", created)
	}
	productID := created.ProductCreate.Product.ID
	deleted := false
	t.Cleanup(func() {
		if !deleted {
			_, _ = sc.DeleteProduct(cliContext(), productID)
		}
	})

	getResult, err := runCommandWithInput(t, env, "", "shopify", "products", "get", "--id", productID)
	if err != nil {
		t.Fatal(err)
	}
	var fetched ShopifyProduct
	decodeCommandJSON(t, getResult, &fetched)
	if fetched.Title != title {
		t.Fatalf("expected created product title %q, got %#v", title, fetched)
	}

	updatedTitle := title + " updated"
	updateResult, err := runCommandWithInput(t, env, "", "shopify", "products", "update", "--product-json", `{"id": `+quoteJSONString(t, productID)+`, "title": `+quoteJSONString(t, updatedTitle)+`}`)
	if err != nil {
		t.Fatal(err)
	}
	var updated ShopifyProductUpdate
	decodeCommandJSON(t, updateResult, &updated)
	if len(updated.ProductUpdate.UserErrors) > 0 {
		t.Fatalf("expected product update without user errors, got %#v", updated.ProductUpdate.UserErrors)
	}
	if updated.ProductUpdate.Product == nil || updated.ProductUpdate.Product.Title != updatedTitle {
		t.Fatalf("expected updated product title %q, got %#v", updatedTitle, updated)
	}

	deleteResult, err := runCommandWithInput(t, env, "", "shopify", "products", "delete", "--id", productID)
	if err != nil {
		t.Fatal(err)
	}
	var deletedProduct ShopifyProductDelete
	decodeCommandJSON(t, deleteResult, &deletedProduct)
	if len(deletedProduct.ProductDelete.UserErrors) > 0 {
		t.Fatalf("expected product delete without user errors, got %#v", deletedProduct.ProductDelete.UserErrors)
	}
	if deletedProduct.ProductDelete.DeletedProductID != productID {
		t.Fatalf("expected deleted product %q, got %#v", productID, deletedProduct)
	}
	deleted = true
}

func quoteJSONString(t *testing.T, value string) string {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}
