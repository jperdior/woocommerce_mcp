package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"woocommerce-mcp/internal/post/application/search_posts"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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

// SearchPostsHandler handles search_posts tool calls
type SearchPostsHandler struct{}

// NewSearchPostsHandler creates a new SearchPostsHandler
func NewSearchPostsHandler() *SearchPostsHandler {
	return &SearchPostsHandler{}
}

// GetToolDefinition returns the MCP tool definition for search_posts
func (h *SearchPostsHandler) GetToolDefinition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "search_posts",
		Description: "Search for blog posts in WordPress sites. Supports various filters like search terms, categories, tags, author, status, and more.",
	}
}

// GetInputSchema returns the input schema for the JSON-RPC tools/list endpoint
func (h *SearchPostsHandler) GetInputSchema() map[string]interface{} {
	return map[string]interface{}{
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
	}
}

// ExecuteMCPTool implements the MCP tool execution
func (h *SearchPostsHandler) ExecuteMCPTool(ctx context.Context, req *mcp.CallToolRequest, input SearchPostsInput) (*mcp.CallToolResult, SearchPostsOutput, error) {
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

// HandleJSONRPC handles JSON-RPC tool calls
func (h *SearchPostsHandler) HandleJSONRPC(c *gin.Context, requestID interface{}, arguments map[string]interface{}) {
	// Convert arguments to SearchPostsInput
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		h.sendJSONRPCError(c, requestID, -32602, "Invalid arguments", err.Error())
		return
	}

	var input SearchPostsInput
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
	content := []map[string]interface{}{
		{
			"type": "text",
			"text": output.Data,
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
func (h *SearchPostsHandler) HandleLegacyHTTP(c *gin.Context, arguments map[string]interface{}) {
	// Convert arguments to SearchPostsInput
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": fmt.Sprintf("Invalid arguments: %v", err)}},
			"isError": true,
		})
		return
	}

	var input SearchPostsInput
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
func (h *SearchPostsHandler) sendSSEResponse(c *gin.Context, response map[string]interface{}) {
	responseData, err := json.Marshal(response)
	if err != nil {
		h.sendJSONRPCError(c, response["id"], -32603, "Internal error", err.Error())
		return
	}

	// Send as SSE format
	c.String(http.StatusOK, "data: %s\n\n", string(responseData))
}

// sendJSONRPCError sends a JSON-RPC error response as SSE
func (h *SearchPostsHandler) sendJSONRPCError(c *gin.Context, id interface{}, code int, message, data string) {
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
