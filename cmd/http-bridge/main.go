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

	post_presentation "woocommerce-mcp/internal/post/presentation"
	product_presentation "woocommerce-mcp/internal/product/presentation"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HTTPBridge provides HTTP endpoints that internally use MCP protocol
type HTTPBridge struct {
	mcpServer      *mcp.Server
	router         *gin.Engine
	productHandler *product_presentation.SearchProductsHandler
	postHandler    *post_presentation.SearchPostsHandler
}

// JsonRpcRequest represents a JSON-RPC 2.0 request (compatible with chatbot-service)
type JsonRpcRequest struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id"`
}

// JsonRpcResponse represents a JSON-RPC 2.0 response
type JsonRpcResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// JsonRpcError represents a JSON-RPC 2.0 error
type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// CallToolRequest represents the call tool request format
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// HTTPToolCall represents the HTTP request format for tool calls
type HTTPToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// NewHTTPBridge creates a new HTTP bridge with MCP server
func NewHTTPBridge() *HTTPBridge {
	// Create handlers
	productHandler := product_presentation.NewSearchProductsHandler()
	postHandler := post_presentation.NewSearchPostsHandler()

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "woocommerce-mcp",
		Version: "1.0.0",
	}, nil)

	// Register tools using handlers
	mcp.AddTool(mcpServer, productHandler.GetToolDefinition(), func(ctx context.Context, req *mcp.CallToolRequest, input product_presentation.SearchProductsInput) (*mcp.CallToolResult, product_presentation.SearchProductsOutput, error) {
		return productHandler.ExecuteMCPTool(ctx, req, input)
	})

	mcp.AddTool(mcpServer, postHandler.GetToolDefinition(), func(ctx context.Context, req *mcp.CallToolRequest, input post_presentation.SearchPostsInput) (*mcp.CallToolResult, post_presentation.SearchPostsOutput, error) {
		return postHandler.ExecuteMCPTool(ctx, req, input)
	})

	// Create HTTP router
	router := gin.Default()

	bridge := &HTTPBridge{
		mcpServer:      mcpServer,
		router:         router,
		productHandler: productHandler,
		postHandler:    postHandler,
	}

	bridge.setupRoutes()
	return bridge
}

// setupRoutes configures the HTTP routes
func (b *HTTPBridge) setupRoutes() {
	// Health endpoint for container health checks
	b.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// JSON-RPC 2.0 endpoint (main endpoint for chatbot-service)
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
			"inputSchema": b.productHandler.GetInputSchema(),
		},
		{
			"name":        "search_posts",
			"description": "Search for blog posts in WordPress sites. Supports various filters like search terms, categories, tags, author, status, and more.",
			"inputSchema": b.postHandler.GetInputSchema(),
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

	// Handle different tools using handlers
	switch callRequest.Name {
	case "search_products":
		b.productHandler.HandleJSONRPC(c, request.ID, callRequest.Arguments)
	case "search_posts":
		b.postHandler.HandleJSONRPC(c, request.ID, callRequest.Arguments)
	default:
		b.sendJsonRpcError(c, request.ID, -32601, "Unknown tool", fmt.Sprintf("Tool '%s' not found", callRequest.Name))
	}
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
		Error: JsonRpcError{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	}

	responseData, _ := json.Marshal(errorResponse)
	c.String(http.StatusOK, "data: %s\n\n", string(responseData))
}

// handleLegacyListTools provides backward compatibility
func (b *HTTPBridge) handleLegacyListTools(c *gin.Context) {
	tools := []map[string]interface{}{
		{
			"name":        "search_products",
			"description": "Search for products in WooCommerce store. Supports various filters like search terms, categories, tags, status, and more.",
			"inputSchema": b.productHandler.GetInputSchema(),
		},
		{
			"name":        "search_posts",
			"description": "Search for blog posts in WordPress sites. Supports various filters like search terms, categories, tags, author, status, and more.",
			"inputSchema": b.postHandler.GetInputSchema(),
		},
	}
	c.JSON(http.StatusOK, map[string]interface{}{"tools": tools})
}

// handleLegacyCallTool provides backward compatibility
func (b *HTTPBridge) handleLegacyCallTool(c *gin.Context) {
	var toolCall HTTPToolCall
	if err := c.ShouldBindJSON(&toolCall); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": fmt.Sprintf("Invalid request format: %v", err)}},
			"isError": true,
		})
		return
	}

	// Handle different tools using handlers
	switch toolCall.Name {
	case "search_products":
		b.productHandler.HandleLegacyHTTP(c, toolCall.Arguments)
	case "search_posts":
		b.postHandler.HandleLegacyHTTP(c, toolCall.Arguments)
	default:
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": fmt.Sprintf("Unknown tool: %s", toolCall.Name)}},
			"isError": true,
		})
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
	return nil
}

// Run starts the HTTP bridge
func Run() error {
	bridge := NewHTTPBridge()
	return bridge.Start()
}

// main is the entry point
func main() {
	if err := Run(); err != nil {
		log.Fatalf("HTTP bridge failed: %v", err)
	}
}
