package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MessageAPI demonstrates how your message API would use MCP client
type MessageAPI struct {
	mcpClient *mcp.Client
	session   *mcp.Session
}

// WooCommerceRequest represents a request to search WooCommerce products
type WooCommerceRequest struct {
	BaseURL        string `json:"base_url"`
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	Search         string `json:"search,omitempty"`
	Category       string `json:"category,omitempty"`
	PerPage        string `json:"per_page,omitempty"`
	Page           string `json:"page,omitempty"`
}

// NewMessageAPI creates a new message API with MCP client
func NewMessageAPI() *MessageAPI {
	return &MessageAPI{}
}

// ConnectToWooCommerceMCP establishes connection to WooCommerce MCP server
func (api *MessageAPI) ConnectToWooCommerceMCP(ctx context.Context) error {
	// Create MCP client
	api.mcpClient = mcp.NewClient(&mcp.Implementation{
		Name:    "message-api",
		Version: "v1.0.0",
	}, nil)

	// Option 1: Connect to MCP server in another Docker container via command
	// This assumes both containers can communicate and docker is available
	transport := &mcp.CommandTransport{
		Command: exec.Command("docker", "exec", "woocommerce-mcp-container", "./woocommerce-mcp"),
	}

	// Option 2: Connect to MCP server binary directly (if in same container/host)
	// transport := &mcp.CommandTransport{
	//     Command: exec.Command("./woocommerce-mcp"),
	// }

	// Option 3: Connect via stdio if you spawn the server differently
	// transport := &mcp.StdioTransport{}

	var err error
	api.session, err = api.mcpClient.Connect(ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WooCommerce MCP: %w", err)
	}

	log.Println("Successfully connected to WooCommerce MCP server")
	return nil
}

// SearchProducts searches for products using the MCP server
func (api *MessageAPI) SearchProducts(ctx context.Context, req WooCommerceRequest) (string, error) {
	if api.session == nil {
		return "", fmt.Errorf("not connected to MCP server")
	}

	// Convert request to MCP tool call parameters
	arguments := map[string]any{
		"base_url":        req.BaseURL,
		"consumer_key":    req.ConsumerKey,
		"consumer_secret": req.ConsumerSecret,
	}

	// Add optional parameters
	if req.Search != "" {
		arguments["search"] = req.Search
	}
	if req.Category != "" {
		arguments["category"] = req.Category
	}
	if req.PerPage != "" {
		arguments["per_page"] = req.PerPage
	}
	if req.Page != "" {
		arguments["page"] = req.Page
	}

	// Call the search_products tool
	params := &mcp.CallToolParams{
		Name:      "search_products",
		Arguments: arguments,
	}

	result, err := api.session.CallTool(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to call search_products tool: %w", err)
	}

	if result.IsError {
		return "", fmt.Errorf("tool returned error")
	}

	// Extract result text
	var resultText string
	for _, content := range result.Content {
		if textContent, ok := content.(*mcp.TextContent); ok {
			resultText += textContent.Text
		}
	}

	return resultText, nil
}

// ProcessMessage simulates processing a user message that requires product search
func (api *MessageAPI) ProcessMessage(ctx context.Context, userMessage string) (string, error) {
	log.Printf("Processing message: %s", userMessage)

	// Example: Parse user intent and extract search parameters
	// In a real implementation, you'd use NLP or pattern matching
	searchReq := WooCommerceRequest{
		BaseURL:        "https://mystore.com", // These would come from config
		ConsumerKey:    "ck_1234567890abcdef", // or be passed by user
		ConsumerSecret: "cs_1234567890abcdef", // securely
		Search:         "sneakers",            // Extracted from user message
		PerPage:        "5",
	}

	// Search for products
	searchResult, err := api.SearchProducts(ctx, searchReq)
	if err != nil {
		return "", fmt.Errorf("failed to search products: %w", err)
	}

	// Format response for user
	response := fmt.Sprintf("I found these products for you:\n\n%s", searchResult)
	return response, nil
}

// Close closes the MCP connection
func (api *MessageAPI) Close() error {
	if api.session != nil {
		return api.session.Close()
	}
	return nil
}

// Example usage
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create message API
	api := NewMessageAPI()

	// Connect to WooCommerce MCP server
	if err := api.ConnectToWooCommerceMCP(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer api.Close()

	// Simulate processing a user message
	response, err := api.ProcessMessage(ctx, "I'm looking for sneakers")
	if err != nil {
		log.Fatalf("Failed to process message: %v", err)
	}

	fmt.Println("Response to user:")
	fmt.Println(response)
}

// Docker Compose Example for your setup:
/*
version: '3.8'
services:
  message-api:
    build: ./chatbot-service
    depends_on:
      - woocommerce-mcp
    environment:
      - WOOCOMMERCE_MCP_CONTAINER=woocommerce-mcp
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock  # To run docker exec commands

  woocommerce-mcp:
    build: ./woocommerce_mcp
    container_name: woocommerce-mcp-container
    command: ["./woocommerce-mcp"]  # Run MCP server in stdio mode

  # Alternative: HTTP Bridge approach
  woocommerce-mcp-http:
    build: ./woocommerce_mcp
    container_name: woocommerce-mcp-http
    command: ["./http-bridge"]
    ports:
      - "8080:8080"
*/
