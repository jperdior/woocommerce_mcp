package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"woocommerce-mcp/internal/product/application/search_products"
	"woocommerce-mcp/internal/product/infrastructure/woocommerce"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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

func main() {
	// Create a server with WooCommerce integration
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "woocommerce-mcp",
		Version: "1.0.0",
	}, nil)

	// Add the search_products tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_products",
		Description: "Search for products in WooCommerce store. Supports various filters like search terms, categories, tags, status, and more.",
	}, SearchProductsTool)

	// Run the server over stdin/stdout until the client disconnects
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
