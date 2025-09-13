package presentation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"woocommerce-mcp/internal/product/application/search_products"

	"github.com/gin-gonic/gin"
	"github.com/jperdior/chatbot-kit/application/query"
)

// SearchProductsRequest represents the request for searching products via MCP
type SearchProductsRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// SearchProductsHandler handles the call_tool MCP endpoint for search_products
// @Summary Search for products via MCP protocol
// @Description Searches for products in WooCommerce store using MCP tool call format
// @Accept json
// @Produce json
// @Param request body SearchProductsRequest true "MCP tool call request"
// @Success 200 {object} ToolResult "Products found successfully"
// @Failure 400 {object} ToolResult "Invalid request"
// @Failure 500 {object} ToolResult "Internal server error"
// @Router /call_tool [post]
func SearchProductsHandler(queryBus query.Bus) gin.HandlerFunc {
	return func(c *gin.Context) {
		var toolCall ToolCall
		if err := c.ShouldBindJSON(&toolCall); err != nil {
			c.JSON(http.StatusBadRequest, ToolResult{
				Content: []Content{{Type: "text", Text: fmt.Sprintf("Invalid request format: %v", err)}},
				IsError: true,
			})
			return
		}

		// Only handle search_products tool calls
		if toolCall.Name != "search_products" {
			c.JSON(http.StatusBadRequest, ToolResult{
				Content: []Content{{Type: "text", Text: fmt.Sprintf("Unknown tool: %s", toolCall.Name)}},
				IsError: true,
			})
			return
		}

		// Convert arguments to query
		query, err := argumentsToQuery(toolCall.Arguments)
		if err != nil {
			c.JSON(http.StatusBadRequest, ToolResult{
				Content: []Content{{Type: "text", Text: fmt.Sprintf("Invalid arguments: %v", err)}},
				IsError: true,
			})
			return
		}

		// Execute query via query bus
		result, err := queryBus.Ask(c.Request.Context(), query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ToolResult{
				Content: []Content{{Type: "text", Text: fmt.Sprintf("Failed to search products: %v", err)}},
				IsError: true,
			})
			return
		}

		// Convert result to MCP response
		response, ok := result.(*search_products.SearchResponse)
		if !ok {
			c.JSON(http.StatusInternalServerError, ToolResult{
				Content: []Content{{Type: "text", Text: "Invalid response type"}},
				IsError: true,
			})
			return
		}

		// Format the response
		if response.IsEmpty() {
			c.JSON(http.StatusOK, ToolResult{
				Content: []Content{{Type: "text", Text: "No products found matching the search criteria."}},
			})
			return
		}

		// Convert response to JSON for display
		responseJSON, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ToolResult{
				Content: []Content{{Type: "text", Text: fmt.Sprintf("Failed to format response: %v", err)}},
				IsError: true,
			})
			return
		}

		resultText := fmt.Sprintf("Found %d product(s) (page %d of %d):\n\n%s",
			response.GetProductCount(),
			response.CurrentPage,
			response.TotalPages,
			string(responseJSON))

		c.JSON(http.StatusOK, ToolResult{
			Content: []Content{{Type: "text", Text: resultText}},
		})
	}
}

// argumentsToQuery converts MCP tool call arguments to SearchProductsQuery
func argumentsToQuery(args map[string]interface{}) (*search_products.SearchProductsQuery, error) {
	// Convert args map to JSON and back to get proper types
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}

	var queryArgs struct {
		BaseURL        string `json:"base_url"`
		ConsumerKey    string `json:"consumer_key"`
		ConsumerSecret string `json:"consumer_secret"`
		Search         string `json:"search,omitempty"`
		Category       string `json:"category,omitempty"`
		Tag            string `json:"tag,omitempty"`
		Status         string `json:"status,omitempty"`
		Type           string `json:"type,omitempty"`
		Featured       string `json:"featured,omitempty"`
		OnSale         string `json:"on_sale,omitempty"`
		MinPrice       string `json:"min_price,omitempty"`
		MaxPrice       string `json:"max_price,omitempty"`
		StockStatus    string `json:"stock_status,omitempty"`
		PerPage        string `json:"per_page,omitempty"`
		Page           string `json:"page,omitempty"`
		Order          string `json:"order,omitempty"`
		OrderBy        string `json:"orderby,omitempty"`
	}

	if err := json.Unmarshal(argsJSON, &queryArgs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	// Validate required fields
	if queryArgs.BaseURL == "" {
		return nil, fmt.Errorf("base_url is required")
	}
	if queryArgs.ConsumerKey == "" {
		return nil, fmt.Errorf("consumer_key is required")
	}
	if queryArgs.ConsumerSecret == "" {
		return nil, fmt.Errorf("consumer_secret is required")
	}

	// Create query
	query := search_products.NewSearchProductsQuery(
		queryArgs.BaseURL,
		queryArgs.ConsumerKey,
		queryArgs.ConsumerSecret,
	)

	// Set optional parameters
	query.Search = queryArgs.Search
	query.Category = queryArgs.Category
	query.Tag = queryArgs.Tag
	query.Status = queryArgs.Status
	query.ProductType = queryArgs.Type
	query.Featured = queryArgs.Featured
	query.OnSale = queryArgs.OnSale
	query.MinPrice = queryArgs.MinPrice
	query.MaxPrice = queryArgs.MaxPrice
	query.StockStatus = queryArgs.StockStatus
	query.PerPage = queryArgs.PerPage
	query.Page = queryArgs.Page
	query.Order = queryArgs.Order
	query.OrderBy = queryArgs.OrderBy

	return query, nil
}
