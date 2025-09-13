package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ManifestHandler serves the MCP manifest file
// @Summary Serves the MCP manifest
// @Description Returns the Model Context Protocol manifest describing server capabilities
// @Accept json
// @Produce json
// @Success 200 {object} object "MCP manifest"
// @Router /manifest.json [get]
func ManifestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File("./manifest.json")
	}
}

// ManifestInfoHandler provides information about the manifest
// @Summary Provides manifest information
// @Description Returns basic information about the MCP server and manifest location
// @Accept json
// @Produce json
// @Success 200 {object} ManifestInfoResponse "Manifest information"
// @Router /manifest [get]
func ManifestInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := ManifestInfoResponse{
			Name:        "woocommerce-mcp",
			Version:     "1.0.0",
			Description: "Model Context Protocol server for WooCommerce integration",
			ManifestURL: "/manifest.json",
			Endpoints: map[string]string{
				"health":     "/health",
				"list_tools": "/list_tools",
				"call_tool":  "/call_tool",
				"manifest":   "/manifest.json",
			},
		}
		c.JSON(http.StatusOK, response)
	}
}

// ManifestInfoResponse represents basic manifest information
type ManifestInfoResponse struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	ManifestURL string            `json:"manifest_url"`
	Endpoints   map[string]string `json:"endpoints"`
}
