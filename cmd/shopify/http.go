package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ShopifyClient struct {
	adminBaseURL string
	client       *http.Client
}

type ShopifyGraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type ShopifyGraphQLError struct {
	Message    string           `json:"message"`
	Locations  []map[string]any `json:"locations,omitempty"`
	Path       []any            `json:"path,omitempty"`
	Extensions map[string]any   `json:"extensions,omitempty"`
}

type shopifyGraphQLResponse[T any] struct {
	Data       T                     `json:"data"`
	Errors     []ShopifyGraphQLError `json:"errors,omitempty"`
	Extensions map[string]any        `json:"extensions,omitempty"`
}

type ShopifyUserError struct {
	Field   []string `json:"field,omitempty"`
	Message string   `json:"message"`
}

type ShopifyShop struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	MyshopifyDomain string               `json:"myshopifyDomain"`
	PrimaryDomain   ShopifyPrimaryDomain `json:"primaryDomain"`
}

type ShopifyPrimaryDomain struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

type ShopifyProduct struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Handle      string `json:"handle"`
	Status      string `json:"status"`
	Vendor      string `json:"vendor"`
	ProductType string `json:"productType"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ShopifyOrder struct {
	ID                       string           `json:"id"`
	Name                     string           `json:"name"`
	Email                    string           `json:"email"`
	CreatedAt                string           `json:"createdAt"`
	DisplayFinancialStatus   string           `json:"displayFinancialStatus"`
	DisplayFulfillmentStatus string           `json:"displayFulfillmentStatus"`
	TotalPriceSet            ShopifyMoneyBag  `json:"totalPriceSet"`
	Customer                 *ShopifyCustomer `json:"customer"`
}

type ShopifyMoneyBag struct {
	ShopMoney ShopifyMoney `json:"shopMoney"`
}

type ShopifyMoney struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

type ShopifyCustomer struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ShopifyInventoryItem struct {
	ID        string `json:"id"`
	SKU       string `json:"sku"`
	Tracked   bool   `json:"tracked"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ShopifyInventoryLevel struct {
	ID         string                     `json:"id"`
	Quantities []ShopifyInventoryQuantity `json:"quantities"`
	Item       ShopifyInventoryItem       `json:"item"`
	Location   ShopifyLocation            `json:"location"`
}

type ShopifyInventoryQuantity struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type ShopifyLocation struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	IsActive bool                   `json:"isActive"`
	Address  ShopifyLocationAddress `json:"address"`
}

type ShopifyLocationAddress struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
	Zip      string `json:"zip"`
}

type ShopifyProductsSearch struct {
	Products ShopifyProductConnection `json:"products"`
}

type ShopifyProductConnection struct {
	Nodes    []ShopifyProduct `json:"nodes"`
	PageInfo ShopifyPageInfo  `json:"pageInfo"`
}

type ShopifyOrdersSearch struct {
	Orders ShopifyOrderConnection `json:"orders"`
}

type ShopifyOrderConnection struct {
	Nodes    []ShopifyOrder  `json:"nodes"`
	PageInfo ShopifyPageInfo `json:"pageInfo"`
}

type ShopifyCustomersSearch struct {
	Customers ShopifyCustomerConnection `json:"customers"`
}

type ShopifyCustomerConnection struct {
	Nodes    []ShopifyCustomer `json:"nodes"`
	PageInfo ShopifyPageInfo   `json:"pageInfo"`
}

type ShopifyInventoryItemsSearch struct {
	InventoryItems ShopifyInventoryItemConnection `json:"inventoryItems"`
}

type ShopifyInventoryItemConnection struct {
	Nodes    []ShopifyInventoryItem `json:"nodes"`
	PageInfo ShopifyPageInfo        `json:"pageInfo"`
}

type ShopifyInventoryLevelsSearch struct {
	InventoryItems ShopifyInventoryItemWithLevelsConnection `json:"inventoryItems"`
}

type ShopifyInventoryItemWithLevelsConnection struct {
	Nodes    []ShopifyInventoryItemWithLevels `json:"nodes"`
	PageInfo ShopifyPageInfo                  `json:"pageInfo"`
}

type ShopifyInventoryItemWithLevels struct {
	ID              string                          `json:"id"`
	SKU             string                          `json:"sku"`
	Tracked         bool                            `json:"tracked"`
	InventoryLevels ShopifyInventoryLevelConnection `json:"inventoryLevels"`
}

type ShopifyInventoryLevelConnection struct {
	Nodes    []ShopifyInventoryLevel `json:"nodes"`
	PageInfo ShopifyPageInfo         `json:"pageInfo"`
}

type ShopifyLocationsList struct {
	Locations ShopifyLocationConnection `json:"locations"`
}

type ShopifyLocationConnection struct {
	Nodes    []ShopifyLocation `json:"nodes"`
	PageInfo ShopifyPageInfo   `json:"pageInfo"`
}

type ShopifyPageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
}

type ShopifyProductCreate struct {
	ProductCreate ShopifyProductMutationPayload `json:"productCreate"`
}

type ShopifyProductUpdate struct {
	ProductUpdate ShopifyProductMutationPayload `json:"productUpdate"`
}

type ShopifyProductMutationPayload struct {
	Product    *ShopifyProduct    `json:"product"`
	UserErrors []ShopifyUserError `json:"userErrors"`
}

type ShopifyProductDelete struct {
	ProductDelete ShopifyProductDeletePayload `json:"productDelete"`
}

type ShopifyProductDeletePayload struct {
	DeletedProductID string             `json:"deletedProductId"`
	UserErrors       []ShopifyUserError `json:"userErrors"`
}

func NewShopifyClient(config Config) ShopifyClient {
	return ShopifyClient{
		adminBaseURL: config.AdminBaseURL,
		client:       http.DefaultClient,
	}
}

func (sc ShopifyClient) GraphQL(request ShopifyGraphQLRequest) ([]byte, error) {
	return sc.GraphQLContext(context.Background(), request)
}

func (sc ShopifyClient) GraphQLContext(ctx context.Context, request ShopifyGraphQLRequest) ([]byte, error) {
	if strings.TrimSpace(request.Query) == "" {
		return nil, fmt.Errorf("query is required")
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, sc.adminBaseURL+"/graphql.json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := sc.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("shopify admin api graphql failed with status %d: %s", response.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

func (sc ShopifyClient) Shop(ctx context.Context) (ShopifyShop, error) {
	var out struct {
		Shop ShopifyShop `json:"shop"`
	}
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyShopQuery}, &out)
	return out.Shop, err
}

func (sc ShopifyClient) SearchProducts(ctx context.Context, input ShopifySearchInput) (ShopifyProductsSearch, error) {
	var out ShopifyProductsSearch
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyProductsSearchQuery, Variables: input.variables()}, &out)
	return out, err
}

func (sc ShopifyClient) GetProduct(ctx context.Context, input ShopifyProductGetInput) (ShopifyProduct, error) {
	if strings.TrimSpace(input.ID) != "" && strings.TrimSpace(input.Handle) != "" {
		return ShopifyProduct{}, fmt.Errorf("exactly one of id or handle is required")
	}
	if strings.TrimSpace(input.ID) != "" {
		var out struct {
			Product ShopifyProduct `json:"product"`
		}
		err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyProductGetQuery, Variables: map[string]any{"id": input.ID}}, &out)
		return out.Product, err
	}
	if strings.TrimSpace(input.Handle) == "" {
		return ShopifyProduct{}, fmt.Errorf("exactly one of id or handle is required")
	}
	search, err := sc.SearchProducts(ctx, ShopifySearchInput{First: 1, Query: "handle:" + input.Handle})
	if err != nil {
		return ShopifyProduct{}, err
	}
	if len(search.Products.Nodes) == 0 {
		return ShopifyProduct{}, nil
	}
	return search.Products.Nodes[0], nil
}

func (sc ShopifyClient) CreateProduct(ctx context.Context, product map[string]any) (ShopifyProductCreate, error) {
	if product == nil {
		return ShopifyProductCreate{}, fmt.Errorf("product is required")
	}
	var out ShopifyProductCreate
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyProductCreateMutation, Variables: map[string]any{"product": product}}, &out)
	return out, err
}

func (sc ShopifyClient) UpdateProduct(ctx context.Context, product map[string]any) (ShopifyProductUpdate, error) {
	if product == nil {
		return ShopifyProductUpdate{}, fmt.Errorf("product is required")
	}
	var out ShopifyProductUpdate
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyProductUpdateMutation, Variables: map[string]any{"product": product}}, &out)
	return out, err
}

func (sc ShopifyClient) DeleteProduct(ctx context.Context, id string) (ShopifyProductDelete, error) {
	if strings.TrimSpace(id) == "" {
		return ShopifyProductDelete{}, fmt.Errorf("id is required")
	}
	var out ShopifyProductDelete
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyProductDeleteMutation, Variables: map[string]any{"input": map[string]any{"id": id}}}, &out)
	return out, err
}

func (sc ShopifyClient) SearchOrders(ctx context.Context, input ShopifySearchInput) (ShopifyOrdersSearch, error) {
	var out ShopifyOrdersSearch
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyOrdersSearchQuery, Variables: input.variables()}, &out)
	return out, err
}

func (sc ShopifyClient) GetOrder(ctx context.Context, id string) (ShopifyOrder, error) {
	if strings.TrimSpace(id) == "" {
		return ShopifyOrder{}, fmt.Errorf("id is required")
	}
	var out struct {
		Order ShopifyOrder `json:"order"`
	}
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyOrderGetQuery, Variables: map[string]any{"id": id}}, &out)
	return out.Order, err
}

func (sc ShopifyClient) SearchCustomers(ctx context.Context, input ShopifySearchInput) (ShopifyCustomersSearch, error) {
	var out ShopifyCustomersSearch
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyCustomersSearchQuery, Variables: input.variables()}, &out)
	return out, err
}

func (sc ShopifyClient) GetCustomer(ctx context.Context, id string) (ShopifyCustomer, error) {
	if strings.TrimSpace(id) == "" {
		return ShopifyCustomer{}, fmt.Errorf("id is required")
	}
	var out struct {
		Customer ShopifyCustomer `json:"customer"`
	}
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyCustomerGetQuery, Variables: map[string]any{"id": id}}, &out)
	return out.Customer, err
}

func (sc ShopifyClient) SearchInventoryItems(ctx context.Context, input ShopifySearchInput) (ShopifyInventoryItemsSearch, error) {
	var out ShopifyInventoryItemsSearch
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyInventoryItemsSearchQuery, Variables: input.variables()}, &out)
	return out, err
}

func (sc ShopifyClient) SearchInventoryLevels(ctx context.Context, input ShopifySearchInput) (ShopifyInventoryLevelsSearch, error) {
	var out ShopifyInventoryLevelsSearch
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyInventoryLevelsSearchQuery, Variables: input.variables()}, &out)
	return out, err
}

func (sc ShopifyClient) ListLocations(ctx context.Context, input ShopifyPaginationInput) (ShopifyLocationsList, error) {
	var out ShopifyLocationsList
	err := sc.graphQLData(ctx, ShopifyGraphQLRequest{Query: shopifyLocationsListQuery, Variables: input.variables()}, &out)
	return out, err
}

func (sc ShopifyClient) graphQLData(ctx context.Context, request ShopifyGraphQLRequest, out any) error {
	body, err := sc.GraphQLContext(ctx, request)
	if err != nil {
		return err
	}

	var response shopifyGraphQLResponse[json.RawMessage]
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}
	if len(response.Errors) > 0 {
		errorsJSON, err := json.Marshal(response.Errors)
		if err != nil {
			return err
		}
		return fmt.Errorf("shopify admin api graphql errors: %s", string(errorsJSON))
	}
	if len(response.Data) == 0 {
		return fmt.Errorf("shopify admin api graphql response did not include data")
	}
	return json.Unmarshal(response.Data, out)
}

type ShopifySearchInput struct {
	First int
	After string
	Query string
}

type ShopifyPaginationInput struct {
	First int
	After string
}

type ShopifyProductGetInput struct {
	ID     string
	Handle string
}

func (input ShopifySearchInput) variables() map[string]any {
	variables := ShopifyPaginationInput{First: input.First, After: input.After}.variables()
	if strings.TrimSpace(input.Query) != "" {
		variables["query"] = input.Query
	}
	return variables
}

func (input ShopifyPaginationInput) variables() map[string]any {
	variables := map[string]any{"first": input.First}
	if strings.TrimSpace(input.After) != "" {
		variables["after"] = input.After
	}
	return variables
}

const shopifyShopQuery = `
query ShopifyShopGet {
  shop {
    id
    name
    myshopifyDomain
    primaryDomain {
      host
      url
    }
  }
}`

const shopifyProductFields = `
id
title
handle
status
vendor
productType
createdAt
updatedAt`

const shopifyProductsSearchQuery = `
query ShopifyProductsSearch($first: Int!, $after: String, $query: String) {
  products(first: $first, after: $after, query: $query) {
    nodes {
      ` + shopifyProductFields + `
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
  }
}`

const shopifyProductGetQuery = `
query ShopifyProductGet($id: ID!) {
  product(id: $id) {
    ` + shopifyProductFields + `
  }
}`

const shopifyProductCreateMutation = `
mutation ShopifyProductCreate($product: ProductCreateInput!) {
  productCreate(product: $product) {
    product {
      ` + shopifyProductFields + `
    }
    userErrors {
      field
      message
    }
  }
}`

const shopifyProductUpdateMutation = `
mutation ShopifyProductUpdate($product: ProductUpdateInput!) {
  productUpdate(product: $product) {
    product {
      ` + shopifyProductFields + `
    }
    userErrors {
      field
      message
    }
  }
}`

const shopifyProductDeleteMutation = `
mutation ShopifyProductDelete($input: ProductDeleteInput!) {
  productDelete(input: $input) {
    deletedProductId
    userErrors {
      field
      message
    }
  }
}`

const shopifyOrderFields = `
id
name
email
createdAt
displayFinancialStatus
displayFulfillmentStatus
totalPriceSet {
  shopMoney {
    amount
    currencyCode
  }
}
customer {
  id
  email
  firstName
  lastName
  displayName
}`

const shopifyOrdersSearchQuery = `
query ShopifyOrdersSearch($first: Int!, $after: String, $query: String) {
  orders(first: $first, after: $after, query: $query) {
    nodes {
      ` + shopifyOrderFields + `
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
  }
}`

const shopifyOrderGetQuery = `
query ShopifyOrderGet($id: ID!) {
  order(id: $id) {
    ` + shopifyOrderFields + `
  }
}`

const shopifyCustomerFields = `
id
email
firstName
lastName
displayName
createdAt
updatedAt`

const shopifyCustomersSearchQuery = `
query ShopifyCustomersSearch($first: Int!, $after: String, $query: String) {
  customers(first: $first, after: $after, query: $query) {
    nodes {
      ` + shopifyCustomerFields + `
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
  }
}`

const shopifyCustomerGetQuery = `
query ShopifyCustomerGet($id: ID!) {
  customer(id: $id) {
    ` + shopifyCustomerFields + `
  }
}`

const shopifyInventoryItemFields = `
id
sku
tracked
createdAt
updatedAt`

const shopifyInventoryItemsSearchQuery = `
query ShopifyInventoryItemsSearch($first: Int!, $after: String, $query: String) {
  inventoryItems(first: $first, after: $after, query: $query) {
    nodes {
      ` + shopifyInventoryItemFields + `
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
  }
}`

const shopifyInventoryLevelsSearchQuery = `
query ShopifyInventoryLevelsSearch($first: Int!, $after: String, $query: String) {
  inventoryItems(first: $first, after: $after, query: $query) {
    nodes {
      id
      sku
      tracked
      inventoryLevels(first: 10) {
        nodes {
          id
          quantities(names: ["available"]) {
            name
            quantity
          }
          item {
            ` + shopifyInventoryItemFields + `
          }
          location {
            id
            name
            isActive
          }
        }
        pageInfo {
          hasNextPage
          hasPreviousPage
          startCursor
          endCursor
        }
      }
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
  }
}`

const shopifyLocationsListQuery = `
query ShopifyLocationsList($first: Int!, $after: String) {
  locations(first: $first, after: $after) {
    nodes {
      id
      name
      isActive
      address {
        address1
        address2
        city
        province
        country
        zip
      }
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
  }
}`
