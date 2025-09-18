package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"woocommerce-mcp/internal/product/application/search_products"
	"woocommerce-mcp/internal/product/infrastructure/woocommerce"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SearchProductsInput defines the input structure for the search_products tool
type SearchProductsInput struct {
	BaseURL        string `json:"base_url" jsonschema:"WooCommerce store base URL (e.g., https://example.com)"`
	ConsumerKey    string `json:"consumer_key" jsonschema:"WooCommerce REST API consumer key"`
	ConsumerSecret string `json:"consumer_secret" jsonschema:"WooCommerce REST API consumer secret"`
	Search         string `json:"search,omitempty" jsonschema:"Search term to filter products"`
	Category       string `json:"category,omitempty" jsonschema:"Category ID or slug to filter products"`
	Tag            string `json:"tag,omitempty" jsonschema:"Tag ID or slug to filter products"`
	Status         string `json:"status,omitempty" jsonschema:"Product status filter (any, draft, pending, private, publish)"`
	Type           string `json:"type,omitempty" jsonschema:"Product type filter (simple, grouped, external, variable)"`
	Featured       string `json:"featured,omitempty" jsonschema:"Limit result set to featured products (true/false)"`
	OnSale         string `json:"on_sale,omitempty" jsonschema:"Limit result set to products on sale (true/false)"`
	MinPrice       string `json:"min_price,omitempty" jsonschema:"Limit result set to products with a minimum price"`
	MaxPrice       string `json:"max_price,omitempty" jsonschema:"Limit result set to products with a maximum price"`
	StockStatus    string `json:"stock_status,omitempty" jsonschema:"Limit result set to products with specified stock status"`
	PerPage        string `json:"per_page,omitempty" jsonschema:"Number of products per page (1-100, default: 10)"`
	Page           string `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
	Order          string `json:"order,omitempty" jsonschema:"Sort order (asc, desc)"`
	OrderBy        string `json:"orderby,omitempty" jsonschema:"Sort by field (date, id, include, title, slug, price, popularity, rating, menu_order)"`
}

// SearchProductsOutput defines the output structure for the search_products tool
type SearchProductsOutput struct {
	Message string `json:"message" jsonschema:"Human-readable message about the search results"`
	Data    string `json:"data" jsonschema:"JSON-formatted product data"`
}

// SearchProductsHandler handles search_products tool calls
type SearchProductsHandler struct{}

// NewSearchProductsHandler creates a new SearchProductsHandler
func NewSearchProductsHandler() *SearchProductsHandler {
	return &SearchProductsHandler{}
}

// GetToolDefinition returns the MCP tool definition for search_products
func (h *SearchProductsHandler) GetToolDefinition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "search_products",
		Description: "Search for products in WooCommerce store. Supports various filters like search terms, categories, tags, status, and more.",
	}
}

// GetInputSchema returns the input schema for the JSON-RPC tools/list endpoint
func (h *SearchProductsHandler) GetInputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"base_url":        map[string]string{"type": "string", "description": "WooCommerce store base URL"},
			"consumer_key":    map[string]string{"type": "string", "description": "WooCommerce REST API consumer key"},
			"consumer_secret": map[string]string{"type": "string", "description": "WooCommerce REST API consumer secret"},
			"search":          map[string]string{"type": "string", "description": "Search term to filter products"},
			"category":        map[string]string{"type": "string", "description": "Category filter"},
			"tag":             map[string]string{"type": "string", "description": "Tag filter"},
			"status":          map[string]string{"type": "string", "description": "Product status filter"},
			"type":            map[string]string{"type": "string", "description": "Product type filter"},
			"featured":        map[string]string{"type": "string", "description": "Featured products filter"},
			"on_sale":         map[string]string{"type": "string", "description": "On sale products filter"},
			"min_price":       map[string]string{"type": "string", "description": "Minimum price filter"},
			"max_price":       map[string]string{"type": "string", "description": "Maximum price filter"},
			"stock_status":    map[string]string{"type": "string", "description": "Stock status filter"},
			"per_page":        map[string]string{"type": "string", "description": "Items per page"},
			"page":            map[string]string{"type": "string", "description": "Page number"},
			"order":           map[string]string{"type": "string", "description": "Sort order"},
			"orderby":         map[string]string{"type": "string", "description": "Sort field"},
		},
		"required": []string{"base_url", "consumer_key", "consumer_secret"},
	}
}

// ExecuteMCPTool implements the MCP tool execution
func (h *SearchProductsHandler) ExecuteMCPTool(ctx context.Context, req *mcp.CallToolRequest, input SearchProductsInput) (*mcp.CallToolResult, SearchProductsOutput, error) {
	// Validate required fields
	if input.BaseURL == "" {
		return nil, SearchProductsOutput{}, fmt.Errorf("base_url is required")
	}
	if input.ConsumerKey == "" {
		return nil, SearchProductsOutput{}, fmt.Errorf("consumer_key is required")
	}
	if input.ConsumerSecret == "" {
		return nil, SearchProductsOutput{}, fmt.Errorf("consumer_secret is required")
	}

	// Create WooCommerce client
	config := woocommerce.NewConfig(input.BaseURL, input.ConsumerKey, input.ConsumerSecret)
	client := woocommerce.NewClient(config)
	repo := woocommerce.NewRepository(client)

	// Create search request
	request := search_products.NewSearchRequest(input.BaseURL, input.ConsumerKey, input.ConsumerSecret)

	// Set optional parameters
	if input.Search != "" {
		request.SetSearch(input.Search)
	}
	if input.Category != "" {
		request.SetCategory(input.Category)
	}
	if input.Tag != "" {
		request.SetTag(input.Tag)
	}
	if input.Status != "" {
		request.SetStatus(input.Status)
	}
	if input.Type != "" {
		request.SetType(input.Type)
	}
	if input.Featured != "" {
		request.SetFeatured(input.Featured)
	}
	if input.OnSale != "" {
		request.SetOnSale(input.OnSale)
	}
	if input.MinPrice != "" || input.MaxPrice != "" {
		request.SetPriceRange(input.MinPrice, input.MaxPrice)
	}
	if input.StockStatus != "" {
		request.SetStockStatus(input.StockStatus)
	}
	if input.PerPage != "" || input.Page != "" {
		request.SetPagination(input.Page, input.PerPage)
	}
	if input.OrderBy != "" || input.Order != "" {
		request.SetSorting(input.OrderBy, input.Order)
	}

	// Execute search
	searcher := search_products.NewProductSearcher(repo)
	response, err := searcher.Execute(ctx, request)
	if err != nil {
		return nil, SearchProductsOutput{}, fmt.Errorf("failed to search products: %w", err)
	}

	// Convert response to JSON
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, SearchProductsOutput{}, fmt.Errorf("failed to serialize response: %w", err)
	}

	// Create human-readable message
	message := fmt.Sprintf("Found %d product(s) out of %d total (page %d of %d)",
		len(response.Products),
		response.TotalCount,
		response.CurrentPage,
		response.TotalPages,
	)

	return nil, SearchProductsOutput{
		Message: message,
		Data:    string(responseJSON),
	}, nil
}

// HandleJSONRPC handles JSON-RPC tool calls
func (h *SearchProductsHandler) HandleJSONRPC(c *gin.Context, requestID interface{}, arguments map[string]interface{}) {
	// Convert arguments to SearchProductsInput
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		h.sendJSONRPCError(c, requestID, -32602, "Invalid arguments", err.Error())
		return
	}

	var input SearchProductsInput
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		h.sendJSONRPCError(c, requestID, -32602, "Invalid input format", err.Error())
		return
	}

	// Call the MCP tool directly
	_, output, err := h.ExecuteMCPTool(c.Request.Context(), nil, input)
	if err != nil {
		h.sendJSONRPCError(c, requestID, -32603, "Tool execution failed", err.Error())
		return
	}

	// Format response as expected by the message API
	resultText := fmt.Sprintf("%s\n\n%s", output.Message, output.Data)
	content := []map[string]interface{}{
		{
			"type": "text",
			"text": resultText,
		},
	}

	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result":  map[string]interface{}{"content": content},
		"id":      requestID,
	}

	h.sendSSEResponse(c, response)
}

// HandleLegacyHTTP handles legacy HTTP tool calls
func (h *SearchProductsHandler) HandleLegacyHTTP(c *gin.Context, arguments map[string]interface{}) {
	// Convert arguments to SearchProductsInput
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": fmt.Sprintf("Invalid arguments: %v", err)}},
			"isError": true,
		})
		return
	}

	var input SearchProductsInput
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": fmt.Sprintf("Invalid input format: %v", err)}},
			"isError": true,
		})
		return
	}

	// Call the MCP tool directly
	_, output, err := h.ExecuteMCPTool(c.Request.Context(), nil, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": fmt.Sprintf("Tool execution failed: %v", err)}},
			"isError": true,
		})
		return
	}

	// Return successful result
	resultText := fmt.Sprintf("%s\n\n%s", output.Message, output.Data)
	c.JSON(http.StatusOK, map[string]interface{}{
		"content": []map[string]interface{}{{"type": "text", "text": resultText}},
	})
}

// sendSSEResponse sends a JSON-RPC response as Server-Sent Event
func (h *SearchProductsHandler) sendSSEResponse(c *gin.Context, response map[string]interface{}) {
	responseData, err := json.Marshal(response)
	if err != nil {
		h.sendJSONRPCError(c, response["id"], -32603, "Internal error", err.Error())
		return
	}

	// Send as SSE format
	c.String(http.StatusOK, "data: %s\n\n", string(responseData))
}

// sendJSONRPCError sends a JSON-RPC error response as SSE
func (h *SearchProductsHandler) sendJSONRPCError(c *gin.Context, id interface{}, code int, message, data string) {
	errorResponse := map[string]interface{}{
		"jsonrpc": "2.0",
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"data":    data,
		},
		"id": id,
	}

	responseData, _ := json.Marshal(errorResponse)
	c.String(http.StatusOK, "data: %s\n\n", string(responseData))
}
