package server

import (
	productPresentation "woocommerce-mcp/internal/product/presentation"

	"github.com/gin-gonic/gin"
	"github.com/jperdior/chatbot-kit/application/query"
)

// registerRoutes registers all application routes
func (s *Server) registerRoutes(queryBus query.Bus) {
	// Add CORS middleware
	s.router.Use(corsMiddleware())

	// Add recovery middleware
	s.router.Use(gin.Recovery())

	// Add logging middleware
	s.router.Use(gin.Logger())

	// Health check endpoint
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "woocommerce-mcp",
			"timestamp": "2024-01-01T00:00:00Z",
		})
	})

	// MCP protocol endpoints
	s.router.GET("/list_tools", productPresentation.ListToolsHandler())
	s.router.POST("/call_tool", productPresentation.SearchProductsHandler(queryBus))
}

// corsMiddleware returns a CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
