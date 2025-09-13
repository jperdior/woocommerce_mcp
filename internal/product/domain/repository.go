package domain

import (
	"context"
	"woocommerce-mcp/kit/domain"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	// Search searches for products based on criteria
	Search(ctx context.Context, criteria *SearchCriteria) ([]*Product, error)

	// FindByID finds a product by its ID
	FindByID(ctx context.Context, id *ProductID) (*Product, error)

	// FindBySKU finds a product by its SKU
	FindBySKU(ctx context.Context, sku string) (*Product, error)

	// Save saves a product
	Save(ctx context.Context, product *Product) error

	// Delete deletes a product
	Delete(ctx context.Context, id *ProductID) error

	// Count returns the total count of products matching criteria
	Count(ctx context.Context, criteria *SearchCriteria) (int64, error)
}

// SearchCriteria represents search criteria for products
type SearchCriteria struct {
	// Search term for name, description, or SKU
	Search string

	// Category filter
	Category string

	// Tag filter
	Tag string

	// Status filter
	Status ProductStatus

	// Type filter
	Type ProductType

	// Featured filter
	Featured *bool

	// On sale filter
	OnSale *bool

	// Price range filters
	MinPrice *Money
	MaxPrice *Money

	// Stock status filter
	StockStatus StockStatus

	// Pagination
	Page    int
	PerPage int

	// Sorting
	OrderBy string
	Order   string
}

// NewSearchCriteria creates a new search criteria with defaults
func NewSearchCriteria() *SearchCriteria {
	return &SearchCriteria{
		Page:    1,
		PerPage: 10,
		OrderBy: "date",
		Order:   "desc",
	}
}

// Validate validates the search criteria
func (sc *SearchCriteria) Validate() error {
	if sc.Page < 1 {
		return domain.NewValidationError("page must be greater than 0")
	}

	if sc.PerPage < 1 {
		sc.PerPage = 10
	}

	if sc.PerPage > 100 {
		sc.PerPage = 100
	}

	// Validate status if provided
	if sc.Status != "" && !sc.Status.IsValid() {
		return domain.NewValidationError("invalid product status")
	}

	// Validate type if provided
	if sc.Type != "" && !sc.Type.IsValid() {
		return domain.NewValidationError("invalid product type")
	}

	// Validate stock status if provided
	if sc.StockStatus != "" && !sc.StockStatus.IsValid() {
		return domain.NewValidationError("invalid stock status")
	}

	// Validate order direction
	if sc.Order != "" && sc.Order != "asc" && sc.Order != "desc" {
		return domain.NewValidationError("order must be 'asc' or 'desc'")
	}

	// Validate order by field
	validOrderByFields := []string{"date", "id", "title", "slug", "price", "popularity", "rating", "menu_order"}
	if sc.OrderBy != "" {
		valid := false
		for _, field := range validOrderByFields {
			if sc.OrderBy == field {
				valid = true
				break
			}
		}
		if !valid {
			return domain.NewValidationError("invalid orderby field")
		}
	}

	return nil
}

// SetSearch sets the search term
func (sc *SearchCriteria) SetSearch(search string) *SearchCriteria {
	sc.Search = search
	return sc
}

// SetCategory sets the category filter
func (sc *SearchCriteria) SetCategory(category string) *SearchCriteria {
	sc.Category = category
	return sc
}

// SetTag sets the tag filter
func (sc *SearchCriteria) SetTag(tag string) *SearchCriteria {
	sc.Tag = tag
	return sc
}

// SetStatus sets the status filter
func (sc *SearchCriteria) SetStatus(status ProductStatus) *SearchCriteria {
	sc.Status = status
	return sc
}

// SetType sets the type filter
func (sc *SearchCriteria) SetType(productType ProductType) *SearchCriteria {
	sc.Type = productType
	return sc
}

// SetFeatured sets the featured filter
func (sc *SearchCriteria) SetFeatured(featured bool) *SearchCriteria {
	sc.Featured = &featured
	return sc
}

// SetOnSale sets the on sale filter
func (sc *SearchCriteria) SetOnSale(onSale bool) *SearchCriteria {
	sc.OnSale = &onSale
	return sc
}

// SetPriceRange sets the price range filters
func (sc *SearchCriteria) SetPriceRange(minPrice, maxPrice *Money) *SearchCriteria {
	sc.MinPrice = minPrice
	sc.MaxPrice = maxPrice
	return sc
}

// SetStockStatus sets the stock status filter
func (sc *SearchCriteria) SetStockStatus(stockStatus StockStatus) *SearchCriteria {
	sc.StockStatus = stockStatus
	return sc
}

// SetPagination sets pagination parameters
func (sc *SearchCriteria) SetPagination(page, perPage int) *SearchCriteria {
	sc.Page = page
	sc.PerPage = perPage
	return sc
}

// SetSorting sets sorting parameters
func (sc *SearchCriteria) SetSorting(orderBy, order string) *SearchCriteria {
	sc.OrderBy = orderBy
	sc.Order = order
	return sc
}
