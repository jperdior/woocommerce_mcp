package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func listTools(c *gin.Context) {
	tools := []Tool{
		{
			Name:        "search_products",
			Description: "Search for products in WooCommerce store. Supports various filters like search terms, categories, tags, status, and more.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"base_url": {
						Type:        "string",
						Description: "WooCommerce store base URL (e.g., https://example.com)",
					},
					"consumer_key": {
						Type:        "string",
						Description: "WooCommerce REST API consumer key",
					},
					"consumer_secret": {
						Type:        "string",
						Description: "WooCommerce REST API consumer secret",
					},
					"search": {
						Type:        "string",
						Description: "Search term to filter products by name, description, or SKU",
					},
					"category": {
						Type:        "string",
						Description: "Category ID or slug to filter products",
					},
					"tag": {
						Type:        "string",
						Description: "Tag ID or slug to filter products",
					},
					"status": {
						Type:        "string",
						Description: "Product status filter (draft, pending, private, publish)",
					},
					"type": {
						Type:        "string",
						Description: "Product type filter (simple, grouped, external, variable)",
					},
					"featured": {
						Type:        "string",
						Description: "Filter by featured products (true/false)",
					},
					"on_sale": {
						Type:        "string",
						Description: "Filter by products on sale (true/false)",
					},
					"min_price": {
						Type:        "string",
						Description: "Minimum price filter",
					},
					"max_price": {
						Type:        "string",
						Description: "Maximum price filter",
					},
					"stock_status": {
						Type:        "string",
						Description: "Stock status filter (instock, outofstock, onbackorder)",
					},
					"per_page": {
						Type:        "string",
						Description: "Number of products per page (default: 10, max: 100)",
					},
					"page": {
						Type:        "string",
						Description: "Page number for pagination (default: 1)",
					},
					"order": {
						Type:        "string",
						Description: "Sort order (asc, desc)",
					},
					"orderby": {
						Type:        "string",
						Description: "Sort by field (date, id, include, title, slug, price, popularity, rating, menu_order)",
					},
				},
				Required: []string{"base_url", "consumer_key", "consumer_secret"},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{"tools": tools})
}

func callTool(c *gin.Context) {
	var toolCall ToolCall
	if err := c.ShouldBindJSON(&toolCall); err != nil {
		c.JSON(http.StatusBadRequest, ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Invalid request format: %v", err)}},
			IsError: true,
		})
		return
	}

	switch toolCall.Name {
	case "search_products":
		handleSearchProducts(c, toolCall.Arguments)
	default:
		c.JSON(http.StatusBadRequest, ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Unknown tool: %s", toolCall.Name)}},
			IsError: true,
		})
	}
}

func handleSearchProducts(c *gin.Context, args map[string]interface{}) {
	// Convert args map to SearchProductsRequest struct
	argsJSON, err := json.Marshal(args)
	if err != nil {
		c.JSON(http.StatusBadRequest, ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Failed to process request arguments: %v", err)}},
			IsError: true,
		})
		return
	}

	var request SearchProductsRequest
	if err := json.Unmarshal(argsJSON, &request); err != nil {
		c.JSON(http.StatusBadRequest, ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Invalid request format: %v", err)}},
			IsError: true,
		})
		return
	}

	// Validate required fields
	if request.BaseURL == "" {
		c.JSON(http.StatusBadRequest, ToolResult{
			Content: []Content{{Type: "text", Text: "Missing required parameter: base_url"}},
			IsError: true,
		})
		return
	}

	if request.ConsumerKey == "" {
		c.JSON(http.StatusBadRequest, ToolResult{
			Content: []Content{{Type: "text", Text: "Missing required parameter: consumer_key"}},
			IsError: true,
		})
		return
	}

	if request.ConsumerSecret == "" {
		c.JSON(http.StatusBadRequest, ToolResult{
			Content: []Content{{Type: "text", Text: "Missing required parameter: consumer_secret"}},
			IsError: true,
		})
		return
	}

	// Perform additional validation
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Validation error: %v", err)}},
			IsError: true,
		})
		return
	}

	// Create WooCommerce client
	config := WooCommerceConfig{
		BaseURL:        request.BaseURL,
		ConsumerKey:    request.ConsumerKey,
		ConsumerSecret: request.ConsumerSecret,
	}
	client := NewWooCommerceClient(config)

	// Convert request to search parameters
	searchParams := request.ToSearchParams()

	// Search products
	products, err := client.SearchProducts(searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Failed to search products: %v", err)}},
			IsError: true,
		})
		return
	}

	// Format response
	if len(products) == 0 {
		c.JSON(http.StatusOK, ToolResult{
			Content: []Content{{Type: "text", Text: "No products found matching the search criteria."}},
		})
		return
	}

	// Convert products to JSON for response
	productsJSON, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Failed to format response: %v", err)}},
			IsError: true,
		})
		return
	}

	resultText := fmt.Sprintf("Found %d product(s):\n\n%s", len(products), string(productsJSON))

	c.JSON(http.StatusOK, ToolResult{
		Content: []Content{{Type: "text", Text: resultText}},
	})
}
