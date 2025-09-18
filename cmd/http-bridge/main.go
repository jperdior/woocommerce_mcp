package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"woocommerce-mcp/internal/post/application/search_posts"
	"woocommerce-mcp/internal/product/application/search_products"
	"woocommerce-mcp/internal/product/infrastructure/woocommerce"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HTTPBridge provides HTTP endpoints that internally use MCP protocol
type HTTPBridge struct {
	mcpServer *mcp.Server
	router    *gin.Engine
}

// JsonRpcRequest represents a JSON-RPC 2.0 request (compatible with chatbot-service)
type JsonRpcRequest struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      interface{} `json:"id"`
	Params  interface{} `json:"params,omitempty"`
}

// JsonRpcResponse represents a JSON-RPC 2.0 response
type JsonRpcResponse struct {
	JsonRpc string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JsonRpcError `json:"error,omitempty"`
	ID      interface{}   `json:"id"`
}

// JsonRpcError represents a JSON-RPC 2.0 error
type JsonRpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CallToolRequest represents the params for tools/call
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// SearchProductsInput defines the input structure for the search_products tool
type SearchProductsInput struct {
	BaseURL        string `json:"base_url" jsonschema:"WooCommerce store base URL (e.g., https://example.com)"`
	ConsumerKey    string `json:"consumer_key" jsonschema:"WooCommerce REST API consumer key"`
	ConsumerSecret string `json:"consumer_secret" jsonschema:"WooCommerce REST API consumer secret"`
	Search         string `json:"search,omitempty" jsonschema:"Search term to filter products by name, description, or SKU"`
	Category       string `json:"category,omitempty" jsonschema:"Category ID or slug to filter products"`
	Tag            string `json:"tag,omitempty" jsonschema:"Tag ID or slug to filter products"`
	Status         string `json:"status,omitempty" jsonschema:"Product status filter (draft, pending, private, publish)"`
	ProductType    string `json:"type,omitempty" jsonschema:"Product type filter (simple, grouped, external, variable)"`
	Featured       string `json:"featured,omitempty" jsonschema:"Filter by featured products (true/false)"`
	OnSale         string `json:"on_sale,omitempty" jsonschema:"Filter by products on sale (true/false)"`
	MinPrice       string `json:"min_price,omitempty" jsonschema:"Minimum price filter"`
	MaxPrice       string `json:"max_price,omitempty" jsonschema:"Maximum price filter"`
	StockStatus    string `json:"stock_status,omitempty" jsonschema:"Stock status filter (instock, outofstock, onbackorder)"`
	PerPage        string `json:"per_page,omitempty" jsonschema:"Number of products per page (default: 10, max: 100)"`
	Page           string `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
	Order          string `json:"order,omitempty" jsonschema:"Sort order (asc, desc)"`
	OrderBy        string `json:"orderby,omitempty" jsonschema:"Sort by field (date, id, include, title, slug, price, popularity, rating, menu_order)"`
}

// SearchProductsOutput defines the output structure for the search_products tool
type SearchProductsOutput struct {
	Message string `json:"message" jsonschema:"Human-readable message about the search results"`
	Data    string `json:"data" jsonschema:"JSON-formatted product data"`
}

// SearchPostsInput defines the input structure for the search_posts tool
type SearchPostsInput struct {
	BaseURL    string `json:"base_url" jsonschema:"WordPress site base URL (e.g., https://example.com)"`
	Search     string `json:"search,omitempty" jsonschema:"Search term to filter posts"`
	Status     string `json:"status,omitempty" jsonschema:"Post status filter (publish, draft, private, pending, trash)"`
	Author     string `json:"author,omitempty" jsonschema:"Author ID filter"`
	Categories string `json:"categories,omitempty" jsonschema:"Comma-separated category IDs"`
	Tags       string `json:"tags,omitempty" jsonschema:"Comma-separated tag IDs"`
	Before     string `json:"before,omitempty" jsonschema:"Limit response to posts published before a given date (ISO 8601 format)"`
	After      string `json:"after,omitempty" jsonschema:"Limit response to posts published after a given date (ISO 8601 format)"`
	Page       string `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
	PerPage    string `json:"per_page,omitempty" jsonschema:"Number of posts per page (default: 10, max: 100)"`
	OrderBy    string `json:"orderby,omitempty" jsonschema:"Sort by field (date, relevance, id, include, title, slug)"`
	Order      string `json:"order,omitempty" jsonschema:"Sort order (asc, desc)"`
}

// SearchPostsOutput defines the output structure for the search_posts tool
type SearchPostsOutput struct {
	Message string `json:"message" jsonschema:"Human-readable message about the search results"`
	Data    string `json:"data" jsonschema:"JSON-formatted post data"`
}

// HTTPToolCall represents the HTTP request format for tool calls
type HTTPToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// HTTPToolResult represents the HTTP response format for tool results
type HTTPToolResult struct {
	Content []HTTPContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// HTTPContent represents content in HTTP responses
type HTTPContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SearchProductsTool implements the search_products MCP tool
func SearchProductsTool(ctx context.Context, req *mcp.CallToolRequest, input SearchProductsInput) (*mcp.CallToolResult, SearchProductsOutput, error) {
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
	request := &search_products.SearchRequest{
		BaseURL:        input.BaseURL,
		ConsumerKey:    input.ConsumerKey,
		ConsumerSecret: input.ConsumerSecret,
	}

	// Set optional parameters
	if input.Search != "" {
		request.Search = &input.Search
	}
	if input.Category != "" {
		request.Category = &input.Category
	}
	if input.Tag != "" {
		request.Tag = &input.Tag
	}
	if input.Status != "" {
		request.Status = &input.Status
	}
	if input.ProductType != "" {
		request.Type = &input.ProductType
	}
	if input.Featured != "" {
		request.Featured = &input.Featured
	}
	if input.OnSale != "" {
		request.OnSale = &input.OnSale
	}
	if input.MinPrice != "" {
		request.MinPrice = &input.MinPrice
	}
	if input.MaxPrice != "" {
		request.MaxPrice = &input.MaxPrice
	}
	if input.StockStatus != "" {
		request.StockStatus = &input.StockStatus
	}
	if input.PerPage != "" {
		request.PerPage = &input.PerPage
	}
	if input.Page != "" {
		request.Page = &input.Page
	}
	if input.Order != "" {
		request.Order = &input.Order
	}
	if input.OrderBy != "" {
		request.OrderBy = &input.OrderBy
	}

	// Execute search
	searcher := search_products.NewProductSearcher(repo)
	response, err := searcher.Execute(ctx, request)
	if err != nil {
		return nil, SearchProductsOutput{}, fmt.Errorf("failed to search products: %w", err)
	}

	// Format response
	if len(response.Products) == 0 {
		return nil, SearchProductsOutput{
			Message: "No products found matching the search criteria.",
			Data:    "[]",
		}, nil
	}

	// Convert response to JSON
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, SearchProductsOutput{}, fmt.Errorf("failed to format response: %w", err)
	}

	message := fmt.Sprintf("Found %d product(s) (page %d of %d)",
		response.TotalCount,
		response.CurrentPage,
		response.TotalPages,
	)

	return nil, SearchProductsOutput{
		Message: message,
		Data:    string(responseJSON),
	}, nil
}

// SearchPostsTool implements the search_posts MCP tool
func SearchPostsTool(ctx context.Context, req *mcp.CallToolRequest, input SearchPostsInput) (*mcp.CallToolResult, SearchPostsOutput, error) {
	// Validate required fields
	if input.BaseURL == "" {
		return nil, SearchPostsOutput{}, fmt.Errorf("base_url is required")
	}

	// Create search request
	request := &search_posts.SearchRequest{
		BaseURL:    input.BaseURL,
		Search:     input.Search,
		Status:     input.Status,
		Author:     input.Author,
		Categories: input.Categories,
		Tags:       input.Tags,
		Before:     input.Before,
		After:      input.After,
		Page:       input.Page,
		PerPage:    input.PerPage,
		OrderBy:    input.OrderBy,
		Order:      input.Order,
	}

	// Execute search
	searcher := search_posts.NewPostSearcher(nil) // We pass nil since the searcher creates its own repository
	response, err := searcher.Execute(ctx, request)
	if err != nil {
		return nil, SearchPostsOutput{}, fmt.Errorf("failed to search posts: %w", err)
	}

	// Convert response to JSON
	jsonData, err := response.ToJSON()
	if err != nil {
		return nil, SearchPostsOutput{}, fmt.Errorf("failed to serialize response: %w", err)
	}

	// Create human-readable message
	var message string
	if len(response.Posts) == 0 {
		message = "No posts found matching the search criteria"
	} else {
		message = fmt.Sprintf("Found %d post(s) (page %d of %d)",
			len(response.Posts), response.CurrentPage, response.TotalPages)
	}

	return nil, SearchPostsOutput{
		Message: message,
		Data:    jsonData,
	}, nil
}

// NewHTTPBridge creates a new HTTP bridge with MCP server
func NewHTTPBridge() *HTTPBridge {
	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "woocommerce-mcp",
		Version: "1.0.0",
	}, nil)

	// Add the search_products tool
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "search_products",
		Description: "Search for products in WooCommerce store. Supports various filters like search terms, categories, tags, status, and more.",
	}, SearchProductsTool)

	// Add the search_posts tool
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "search_posts",
		Description: "Search for blog posts in WordPress sites. Supports various filters like search terms, categories, tags, author, status, and more.",
	}, SearchPostsTool)

	// Create HTTP router
	router := gin.Default()

	bridge := &HTTPBridge{
		mcpServer: mcpServer,
		router:    router,
	}

	bridge.setupRoutes()
	return bridge
}

// setupRoutes configures the HTTP routes
func (b *HTTPBridge) setupRoutes() {
	// Health check
	b.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// MCP manifest (for compatibility)
	b.router.GET("/manifest.json", func(c *gin.Context) {
		c.File("./manifest.json")
	})

	// JSON-RPC endpoint (compatible with chatbot-service MCP client)
	b.router.POST("/", b.handleJsonRpc)

	// Legacy endpoints for backward compatibility
	b.router.GET("/list_tools", b.handleLegacyListTools)
	b.router.POST("/call_tool", b.handleLegacyCallTool)
}

// handleJsonRpc handles JSON-RPC 2.0 requests with SSE responses
func (b *HTTPBridge) handleJsonRpc(c *gin.Context) {
	var request JsonRpcRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		b.sendJsonRpcError(c, request.ID, -32700, "Parse error", err.Error())
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	switch request.Method {
	case "tools/list":
		b.handleToolsList(c, request)
	case "tools/call":
		b.handleToolsCall(c, request)
	default:
		b.sendJsonRpcError(c, request.ID, -32601, "Method not found", fmt.Sprintf("Unknown method: %s", request.Method))
	}
}

// handleToolsList handles the tools/list JSON-RPC method
func (b *HTTPBridge) handleToolsList(c *gin.Context, request JsonRpcRequest) {
	tools := []map[string]interface{}{
		{
			"name":        "search_products",
			"description": "Search for products in WooCommerce store. Supports various filters like search terms, categories, tags, status, and more.",
			"inputSchema": map[string]interface{}{
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
			},
		},
		{
			"name":        "search_posts",
			"description": "Search for blog posts in WordPress sites. Supports various filters like search terms, categories, tags, author, status, and more.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"base_url":   map[string]string{"type": "string", "description": "WordPress site base URL"},
					"search":     map[string]string{"type": "string", "description": "Search term to filter posts"},
					"status":     map[string]string{"type": "string", "description": "Post status filter"},
					"author":     map[string]string{"type": "string", "description": "Author ID filter"},
					"categories": map[string]string{"type": "string", "description": "Comma-separated category IDs"},
					"tags":       map[string]string{"type": "string", "description": "Comma-separated tag IDs"},
					"before":     map[string]string{"type": "string", "description": "Posts published before date (ISO 8601)"},
					"after":      map[string]string{"type": "string", "description": "Posts published after date (ISO 8601)"},
					"per_page":   map[string]string{"type": "string", "description": "Number of posts per page"},
					"page":       map[string]string{"type": "string", "description": "Page number"},
					"order":      map[string]string{"type": "string", "description": "Sort order"},
					"orderby":    map[string]string{"type": "string", "description": "Sort field"},
				},
				"required": []string{"base_url"},
			},
		},
	}

	response := JsonRpcResponse{
		JsonRpc: "2.0",
		Result:  map[string]interface{}{"tools": tools},
		ID:      request.ID,
	}

	b.sendSSEResponse(c, response)
}

// handleToolsCall handles the tools/call JSON-RPC method
func (b *HTTPBridge) handleToolsCall(c *gin.Context, request JsonRpcRequest) {
	// Parse params
	paramsJSON, err := json.Marshal(request.Params)
	if err != nil {
		b.sendJsonRpcError(c, request.ID, -32602, "Invalid params", err.Error())
		return
	}

	var callRequest CallToolRequest
	if err := json.Unmarshal(paramsJSON, &callRequest); err != nil {
		b.sendJsonRpcError(c, request.ID, -32602, "Invalid params", err.Error())
		return
	}

	// Handle different tools
	switch callRequest.Name {
	case "search_products":
		b.handleSearchProducts(c, request, callRequest)
	case "search_posts":
		b.handleSearchPosts(c, request, callRequest)
	default:
		b.sendJsonRpcError(c, request.ID, -32601, "Unknown tool", fmt.Sprintf("Tool '%s' not found", callRequest.Name))
	}
}

// handleSearchProducts handles the search_products tool call
func (b *HTTPBridge) handleSearchProducts(c *gin.Context, request JsonRpcRequest, callRequest CallToolRequest) {
	// Convert arguments to SearchProductsInput
	argsJSON, err := json.Marshal(callRequest.Arguments)
	if err != nil {
		b.sendJsonRpcError(c, request.ID, -32602, "Invalid arguments", err.Error())
		return
	}

	var input SearchProductsInput
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		b.sendJsonRpcError(c, request.ID, -32602, "Invalid input format", err.Error())
		return
	}

	// Call the MCP tool directly
	_, output, err := SearchProductsTool(c.Request.Context(), nil, input)
	if err != nil {
		b.sendJsonRpcError(c, request.ID, -32603, "Tool execution failed", err.Error())
		return
	}

	// Format response as expected by chatbot-service
	resultText := fmt.Sprintf("%s\n\n%s", output.Message, output.Data)
	content := []map[string]interface{}{
		{
			"type": "text",
			"text": resultText,
		},
	}

	response := JsonRpcResponse{
		JsonRpc: "2.0",
		Result:  map[string]interface{}{"content": content},
		ID:      request.ID,
	}

	b.sendSSEResponse(c, response)
}

// handleSearchPosts handles the search_posts tool call
func (b *HTTPBridge) handleSearchPosts(c *gin.Context, request JsonRpcRequest, callRequest CallToolRequest) {
	// Convert arguments to SearchPostsInput
	argsJSON, err := json.Marshal(callRequest.Arguments)
	if err != nil {
		b.sendJsonRpcError(c, request.ID, -32602, "Invalid arguments", err.Error())
		return
	}

	var input SearchPostsInput
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		b.sendJsonRpcError(c, request.ID, -32602, "Invalid input format", err.Error())
		return
	}

	// Call the MCP tool directly
	_, output, err := SearchPostsTool(c.Request.Context(), nil, input)
	if err != nil {
		b.sendJsonRpcError(c, request.ID, -32603, "Tool execution failed", err.Error())
		return
	}

	// Format response as expected by the message API
	content := []map[string]interface{}{
		{
			"type": "text",
			"text": output.Data,
		},
	}

	response := JsonRpcResponse{
		JsonRpc: "2.0",
		Result:  map[string]interface{}{"content": content},
		ID:      request.ID,
	}

	b.sendSSEResponse(c, response)
}

// sendSSEResponse sends a JSON-RPC response as Server-Sent Event
func (b *HTTPBridge) sendSSEResponse(c *gin.Context, response JsonRpcResponse) {
	responseData, err := json.Marshal(response)
	if err != nil {
		b.sendJsonRpcError(c, response.ID, -32603, "Internal error", err.Error())
		return
	}

	// Send as SSE format
	c.String(http.StatusOK, "data: %s\n\n", string(responseData))
}

// sendJsonRpcError sends a JSON-RPC error response as SSE
func (b *HTTPBridge) sendJsonRpcError(c *gin.Context, id interface{}, code int, message, data string) {
	errorResponse := JsonRpcResponse{
		JsonRpc: "2.0",
		Error: &JsonRpcError{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	}

	responseData, _ := json.Marshal(errorResponse)
	c.Header("Content-Type", "text/event-stream")
	c.String(http.StatusOK, "data: %s\n\n", string(responseData))
}

// handleLegacyListTools provides backward compatibility
func (b *HTTPBridge) handleLegacyListTools(c *gin.Context) {
	tools := []map[string]interface{}{
		{
			"name":        "search_products",
			"description": "Search for products in WooCommerce store. Supports various filters like search terms, categories, tags, status, and more.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"base_url":        map[string]string{"type": "string", "description": "WooCommerce store base URL"},
					"consumer_key":    map[string]string{"type": "string", "description": "WooCommerce REST API consumer key"},
					"consumer_secret": map[string]string{"type": "string", "description": "WooCommerce REST API consumer secret"},
					"search":          map[string]string{"type": "string", "description": "Search term to filter products"},
				},
				"required": []string{"base_url", "consumer_key", "consumer_secret"},
			},
		},
		{
			"name":        "search_posts",
			"description": "Search for blog posts in WordPress sites. Supports various filters like search terms, categories, tags, author, status, and more.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"base_url": map[string]string{"type": "string", "description": "WordPress site base URL"},
					"search":   map[string]string{"type": "string", "description": "Search term to filter posts"},
				},
				"required": []string{"base_url"},
			},
		},
	}
	c.JSON(http.StatusOK, map[string]interface{}{"tools": tools})
}

// handleLegacyCallTool provides backward compatibility
func (b *HTTPBridge) handleLegacyCallTool(c *gin.Context) {
	var toolCall HTTPToolCall
	if err := c.ShouldBindJSON(&toolCall); err != nil {
		c.JSON(http.StatusBadRequest, HTTPToolResult{
			Content: []HTTPContent{{Type: "text", Text: fmt.Sprintf("Invalid request format: %v", err)}},
			IsError: true,
		})
		return
	}

	// Handle different tools
	switch toolCall.Name {
	case "search_products":
		b.handleLegacySearchProducts(c, toolCall)
	case "search_posts":
		b.handleLegacySearchPosts(c, toolCall)
	default:
		c.JSON(http.StatusBadRequest, HTTPToolResult{
			Content: []HTTPContent{{Type: "text", Text: fmt.Sprintf("Unknown tool: %s", toolCall.Name)}},
			IsError: true,
		})
	}
}

// handleLegacySearchProducts handles legacy search_products calls
func (b *HTTPBridge) handleLegacySearchProducts(c *gin.Context, toolCall HTTPToolCall) {

	// Convert arguments to SearchProductsInput
	argsJSON, err := json.Marshal(toolCall.Arguments)
	if err != nil {
		c.JSON(http.StatusBadRequest, HTTPToolResult{
			Content: []HTTPContent{{Type: "text", Text: fmt.Sprintf("Invalid arguments: %v", err)}},
			IsError: true,
		})
		return
	}

	var input SearchProductsInput
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		c.JSON(http.StatusBadRequest, HTTPToolResult{
			Content: []HTTPContent{{Type: "text", Text: fmt.Sprintf("Invalid input format: %v", err)}},
			IsError: true,
		})
		return
	}

	// Call the MCP tool directly
	_, output, err := SearchProductsTool(c.Request.Context(), nil, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HTTPToolResult{
			Content: []HTTPContent{{Type: "text", Text: fmt.Sprintf("Tool execution failed: %v", err)}},
			IsError: true,
		})
		return
	}

	// Return successful result
	resultText := fmt.Sprintf("%s\n\n%s", output.Message, output.Data)
	c.JSON(http.StatusOK, HTTPToolResult{
		Content: []HTTPContent{{Type: "text", Text: resultText}},
	})
}

// handleLegacySearchPosts handles legacy search_posts calls
func (b *HTTPBridge) handleLegacySearchPosts(c *gin.Context, toolCall HTTPToolCall) {
	// Convert arguments to SearchPostsInput
	argsJSON, err := json.Marshal(toolCall.Arguments)
	if err != nil {
		c.JSON(http.StatusBadRequest, HTTPToolResult{
			Content: []HTTPContent{{Type: "text", Text: fmt.Sprintf("Invalid arguments: %v", err)}},
			IsError: true,
		})
		return
	}

	var input SearchPostsInput
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		c.JSON(http.StatusBadRequest, HTTPToolResult{
			Content: []HTTPContent{{Type: "text", Text: fmt.Sprintf("Invalid input format: %v", err)}},
			IsError: true,
		})
		return
	}

	// Call the MCP tool directly
	_, output, err := SearchPostsTool(c.Request.Context(), nil, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HTTPToolResult{
			Content: []HTTPContent{{Type: "text", Text: fmt.Sprintf("Tool execution failed: %v", err)}},
			IsError: true,
		})
		return
	}

	// Return successful result
	resultText := fmt.Sprintf("%s\n\n%s", output.Message, output.Data)
	c.JSON(http.StatusOK, HTTPToolResult{
		Content: []HTTPContent{{Type: "text", Text: resultText}},
	})
}

// Start starts the HTTP bridge server
func (b *HTTPBridge) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: b.router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting WooCommerce MCP HTTP Bridge on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}

func main() {
	bridge := NewHTTPBridge()
	if err := bridge.Start(); err != nil {
		log.Fatal(err)
	}
}
