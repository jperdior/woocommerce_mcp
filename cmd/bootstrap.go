package cmd

import (
	"woocommerce-mcp/internal/container"
)

// AppComponents holds all application components
type AppComponents struct {
	Container *container.Container
}

// InitializeComponents initializes all application components
func InitializeComponents() (*AppComponents, error) {
	// Create and build the container
	c := container.NewContainer()
	if err := c.Build(); err != nil {
		return nil, err
	}

	return &AppComponents{
		Container: c,
	}, nil
}

// GetContainer returns the dependency injection container
func (ac *AppComponents) GetContainer() *container.Container {
	return ac.Container
}
