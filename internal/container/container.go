package container

import (
	"woocommerce-mcp/internal/platform/server"
	"woocommerce-mcp/internal/product/application/search_products"
	"woocommerce-mcp/internal/product/infrastructure/woocommerce"

	"github.com/jperdior/chatbot-kit/application/query"
	"github.com/jperdior/chatbot-kit/infrastructure/bus/inmemory"
)

// Container holds all application dependencies
type Container struct {
	// Configuration
	ServerConfig *server.Config

	// Buses
	QueryBus query.Bus

	// Services
	ProductSearcher *search_products.ProductSearcher

	// Infrastructure
	Server *server.Server
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{}
}

// Build builds the container with all dependencies
func (c *Container) Build() error {
	// Initialize configuration
	c.ServerConfig = server.NewConfig()

	// Initialize query bus
	c.QueryBus = inmemory.NewQueryBus()

	// Initialize services
	// Create a default repository (will be overridden per request with credentials)
	defaultRepository := woocommerce.NewRepositoryFromConfig("", "", "")
	c.ProductSearcher = search_products.NewProductSearcher(defaultRepository)

	// Register query handlers
	c.registerQueryHandlers()

	// Initialize server with query bus
	c.Server = server.NewServer(c.ServerConfig, c.QueryBus)

	return nil
}

// registerQueryHandlers registers all query handlers with the query bus
func (c *Container) registerQueryHandlers() {
	// Register search products query handler
	searchProductsHandler := search_products.NewSearchProductsQueryHandler(c.ProductSearcher)
	c.QueryBus.Register(search_products.SearchProductsQueryType, searchProductsHandler)
}

// GetServer returns the HTTP server
func (c *Container) GetServer() *server.Server {
	return c.Server
}

// GetQueryBus returns the query bus
func (c *Container) GetQueryBus() query.Bus {
	return c.QueryBus
}

// GetServerConfig returns the server configuration
func (c *Container) GetServerConfig() *server.Config {
	return c.ServerConfig
}
