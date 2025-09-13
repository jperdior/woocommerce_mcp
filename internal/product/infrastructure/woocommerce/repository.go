package woocommerce

import (
	"context"
	"fmt"
	"woocommerce-mcp/internal/product/domain"
	kitDomain "woocommerce-mcp/kit/domain"
)

// Repository implements the ProductRepository interface using WooCommerce API
type Repository struct {
	client *Client
}

// NewRepository creates a new WooCommerce repository
func NewRepository(client *Client) *Repository {
	return &Repository{
		client: client,
	}
}

// Search searches for products based on criteria
func (r *Repository) Search(ctx context.Context, criteria *domain.SearchCriteria) ([]*domain.Product, error) {
	if criteria == nil {
		return nil, kitDomain.NewValidationError("search criteria cannot be nil")
	}

	products, err := r.client.SearchProducts(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	return products, nil
}

// FindByID finds a product by its ID
func (r *Repository) FindByID(ctx context.Context, id *domain.ProductID) (*domain.Product, error) {
	if id == nil {
		return nil, kitDomain.NewValidationError("product ID cannot be nil")
	}

	// Create search criteria to find by ID
	criteria := domain.NewSearchCriteria()
	criteria.SetPagination(1, 1)

	// WooCommerce API doesn't have a direct "find by ID" endpoint in the search,
	// so we need to use the specific product endpoint or search with include parameter
	// For now, we'll search and filter by ID
	products, err := r.client.SearchProducts(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to find product by ID: %w", err)
	}

	// Find the product with the matching ID
	for _, product := range products {
		if product.ID.Equals(id) {
			return product, nil
		}
	}

	return nil, domain.NewProductNotFoundError(id)
}

// FindBySKU finds a product by its SKU
func (r *Repository) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	if sku == "" {
		return nil, kitDomain.NewValidationError("SKU cannot be empty")
	}

	// Create search criteria to find by SKU
	criteria := domain.NewSearchCriteria()
	criteria.SetSearch(sku)        // WooCommerce search includes SKU
	criteria.SetPagination(1, 100) // Get more results to ensure we find the exact SKU match

	products, err := r.client.SearchProducts(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to find product by SKU: %w", err)
	}

	// Find the product with the exact SKU match
	for _, product := range products {
		if product.SKU == sku {
			return product, nil
		}
	}

	return nil, kitDomain.NewNotFoundError("product", sku)
}

// Save saves a product (not implemented for read-only MCP)
func (r *Repository) Save(ctx context.Context, product *domain.Product) error {
	return kitDomain.NewDomainError("NOT_IMPLEMENTED", "save operation is not supported in read-only mode")
}

// Delete deletes a product (not implemented for read-only MCP)
func (r *Repository) Delete(ctx context.Context, id *domain.ProductID) error {
	return kitDomain.NewDomainError("NOT_IMPLEMENTED", "delete operation is not supported in read-only mode")
}

// Count returns the total count of products matching criteria
func (r *Repository) Count(ctx context.Context, criteria *domain.SearchCriteria) (int64, error) {
	if criteria == nil {
		return 0, kitDomain.NewValidationError("search criteria cannot be nil")
	}

	count, err := r.client.CountProducts(ctx, criteria)
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

// NewRepositoryFromConfig creates a new repository from configuration
func NewRepositoryFromConfig(baseURL, consumerKey, consumerSecret string) *Repository {
	config := NewConfig(baseURL, consumerKey, consumerSecret)
	client := NewClient(config)
	return NewRepository(client)
}
