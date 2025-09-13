package main

import (
	"log"
	"woocommerce-mcp/cmd"
)

func main() {
	// Initialize application components
	components, err := cmd.InitializeComponents()
	if err != nil {
		log.Fatalf("Failed to initialize components: %v", err)
	}

	// Get the server from the container
	server := components.GetContainer().GetServer()

	// Start the server (this will block until shutdown)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
