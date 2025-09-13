package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListToolsHandler handles the list_tools MCP endpoint
// @Summary Lists available MCP tools
// @Description Returns a list of available tools for the MCP protocol
// @Accept json
// @Produce json
// @Success 200 {object} ToolsResponse "Tools listed successfully"
// @Router /list_tools [get]
func ListToolsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		response := ToolsResponse{Tools: tools}
		c.JSON(http.StatusOK, response)
	}
}

// ToolsResponse represents the response for listing tools
type ToolsResponse struct {
	Tools []Tool `json:"tools"`
}
