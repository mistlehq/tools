package main

type commandDoc struct {
	Command     string
	Summary     string
	Description string
}

var shopifyAuthTestDoc = commandDoc{
	Command:     "shopify auth test",
	Summary:     "Check Shopify Admin API access",
	Description: "Check Shopify Admin API access by fetching the current shop.",
}

var shopifyGraphQLRequestDoc = commandDoc{
	Command:     "shopify graphql request",
	Summary:     "Send a raw Shopify Admin GraphQL request",
	Description: "Send a raw Shopify Admin GraphQL request. This is the full Admin API coverage surface.",
}

var shopifyShopGetDoc = commandDoc{
	Command:     "shopify shop get",
	Summary:     "Get Shopify shop details",
	Description: "Get details for the current Shopify shop.",
}

var shopifyProductsSearchDoc = commandDoc{
	Command:     "shopify products search",
	Summary:     "Search Shopify products",
	Description: "Search Shopify products with Shopify Admin API search syntax.",
}

var shopifyProductGetDoc = commandDoc{
	Command:     "shopify products get",
	Summary:     "Get a Shopify product",
	Description: "Get a Shopify product by ID or handle.",
}

var shopifyProductCreateDoc = commandDoc{
	Command:     "shopify products create",
	Summary:     "Create a Shopify product",
	Description: "Create a Shopify product using Shopify's ProductCreateInput JSON shape.",
}

var shopifyProductUpdateDoc = commandDoc{
	Command:     "shopify products update",
	Summary:     "Update a Shopify product",
	Description: "Update a Shopify product using Shopify's ProductUpdateInput JSON shape.",
}

var shopifyProductDeleteDoc = commandDoc{
	Command:     "shopify products delete",
	Summary:     "Delete a Shopify product",
	Description: "Delete a Shopify product by ID.",
}

var shopifyOrdersSearchDoc = commandDoc{
	Command:     "shopify orders search",
	Summary:     "Search Shopify orders",
	Description: "Search Shopify orders with Shopify Admin API search syntax.",
}

var shopifyOrderGetDoc = commandDoc{
	Command:     "shopify orders get",
	Summary:     "Get a Shopify order",
	Description: "Get a Shopify order by ID.",
}

var shopifyCustomersSearchDoc = commandDoc{
	Command:     "shopify customers search",
	Summary:     "Search Shopify customers",
	Description: "Search Shopify customers with Shopify Admin API search syntax.",
}

var shopifyCustomerGetDoc = commandDoc{
	Command:     "shopify customers get",
	Summary:     "Get a Shopify customer",
	Description: "Get a Shopify customer by ID.",
}

var shopifyInventoryItemsSearchDoc = commandDoc{
	Command:     "shopify inventory items search",
	Summary:     "Search Shopify inventory items",
	Description: "Search Shopify inventory items with Shopify Admin API search syntax.",
}

var shopifyInventoryLevelsSearchDoc = commandDoc{
	Command:     "shopify inventory levels search",
	Summary:     "Search Shopify inventory levels",
	Description: "Search Shopify inventory levels with Shopify Admin API search syntax.",
}

var shopifyLocationsListDoc = commandDoc{
	Command:     "shopify locations list",
	Summary:     "List Shopify locations",
	Description: "List Shopify locations.",
}
