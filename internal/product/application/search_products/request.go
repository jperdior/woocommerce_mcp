package search_products

import (
	"woocommerce-mcp/kit/domain"
)

// SearchRequest represents a request to search for products
type SearchRequest struct {
	// Required authentication parameters
	BaseURL        string `json:"base_url" binding:"required"`
	ConsumerKey    string `json:"consumer_key" binding:"required"`
	ConsumerSecret string `json:"consumer_secret" binding:"required"`

	// Optional search parameters
	Search      *string `json:"search,omitempty"`
	Category    *string `json:"category,omitempty"`
	Tag         *string `json:"tag,omitempty"`
	Status      *string `json:"status,omitempty"`
	Type        *string `json:"type,omitempty"`
	Featured    *string `json:"featured,omitempty"`
	OnSale      *string `json:"on_sale,omitempty"`
	MinPrice    *string `json:"min_price,omitempty"`
	MaxPrice    *string `json:"max_price,omitempty"`
	StockStatus *string `json:"stock_status,omitempty"`
	PerPage     *string `json:"per_page,omitempty"`
	Page        *string `json:"page,omitempty"`
	Order       *string `json:"order,omitempty"`
	OrderBy     *string `json:"orderby,omitempty"`
}

// NewSearchRequest creates a new SearchRequest
func NewSearchRequest(baseURL, consumerKey, consumerSecret string) *SearchRequest {
	return &SearchRequest{
		BaseURL:        baseURL,
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	}
}

// Validate validates the search request
func (sr *SearchRequest) Validate() error {
	// Validate required fields
	if sr.BaseURL == "" {
		return domain.NewValidationError("base_url is required")
	}

	if sr.ConsumerKey == "" {
		return domain.NewValidationError("consumer_key is required")
	}

	if sr.ConsumerSecret == "" {
		return domain.NewValidationError("consumer_secret is required")
	}

	return nil
}

// SetSearch sets the search term
func (sr *SearchRequest) SetSearch(search string) *SearchRequest {
	sr.Search = &search
	return sr
}

// SetCategory sets the category filter
func (sr *SearchRequest) SetCategory(category string) *SearchRequest {
	sr.Category = &category
	return sr
}

// SetTag sets the tag filter
func (sr *SearchRequest) SetTag(tag string) *SearchRequest {
	sr.Tag = &tag
	return sr
}

// SetStatus sets the status filter
func (sr *SearchRequest) SetStatus(status string) *SearchRequest {
	sr.Status = &status
	return sr
}

// SetType sets the type filter
func (sr *SearchRequest) SetType(productType string) *SearchRequest {
	sr.Type = &productType
	return sr
}

// SetFeatured sets the featured filter
func (sr *SearchRequest) SetFeatured(featured string) *SearchRequest {
	sr.Featured = &featured
	return sr
}

// SetOnSale sets the on sale filter
func (sr *SearchRequest) SetOnSale(onSale string) *SearchRequest {
	sr.OnSale = &onSale
	return sr
}

// SetPriceRange sets the price range filters
func (sr *SearchRequest) SetPriceRange(minPrice, maxPrice string) *SearchRequest {
	if minPrice != "" {
		sr.MinPrice = &minPrice
	}
	if maxPrice != "" {
		sr.MaxPrice = &maxPrice
	}
	return sr
}

// SetStockStatus sets the stock status filter
func (sr *SearchRequest) SetStockStatus(stockStatus string) *SearchRequest {
	sr.StockStatus = &stockStatus
	return sr
}

// SetPagination sets pagination parameters
func (sr *SearchRequest) SetPagination(page, perPage string) *SearchRequest {
	if page != "" {
		sr.Page = &page
	}
	if perPage != "" {
		sr.PerPage = &perPage
	}
	return sr
}

// SetSorting sets sorting parameters
func (sr *SearchRequest) SetSorting(orderBy, order string) *SearchRequest {
	if orderBy != "" {
		sr.OrderBy = &orderBy
	}
	if order != "" {
		sr.Order = &order
	}
	return sr
}

// GetBaseURL returns the base URL
func (sr *SearchRequest) GetBaseURL() string {
	return sr.BaseURL
}

// GetConsumerKey returns the consumer key
func (sr *SearchRequest) GetConsumerKey() string {
	return sr.ConsumerKey
}

// GetConsumerSecret returns the consumer secret
func (sr *SearchRequest) GetConsumerSecret() string {
	return sr.ConsumerSecret
}

// GetSearch returns the search term
func (sr *SearchRequest) GetSearch() string {
	if sr.Search != nil {
		return *sr.Search
	}
	return ""
}

// GetCategory returns the category filter
func (sr *SearchRequest) GetCategory() string {
	if sr.Category != nil {
		return *sr.Category
	}
	return ""
}

// GetTag returns the tag filter
func (sr *SearchRequest) GetTag() string {
	if sr.Tag != nil {
		return *sr.Tag
	}
	return ""
}

// GetStatus returns the status filter
func (sr *SearchRequest) GetStatus() string {
	if sr.Status != nil {
		return *sr.Status
	}
	return ""
}

// GetType returns the type filter
func (sr *SearchRequest) GetType() string {
	if sr.Type != nil {
		return *sr.Type
	}
	return ""
}

// GetFeatured returns the featured filter
func (sr *SearchRequest) GetFeatured() string {
	if sr.Featured != nil {
		return *sr.Featured
	}
	return ""
}

// GetOnSale returns the on sale filter
func (sr *SearchRequest) GetOnSale() string {
	if sr.OnSale != nil {
		return *sr.OnSale
	}
	return ""
}

// GetMinPrice returns the minimum price filter
func (sr *SearchRequest) GetMinPrice() string {
	if sr.MinPrice != nil {
		return *sr.MinPrice
	}
	return ""
}

// GetMaxPrice returns the maximum price filter
func (sr *SearchRequest) GetMaxPrice() string {
	if sr.MaxPrice != nil {
		return *sr.MaxPrice
	}
	return ""
}

// GetStockStatus returns the stock status filter
func (sr *SearchRequest) GetStockStatus() string {
	if sr.StockStatus != nil {
		return *sr.StockStatus
	}
	return ""
}

// GetPerPage returns the per page parameter
func (sr *SearchRequest) GetPerPage() string {
	if sr.PerPage != nil {
		return *sr.PerPage
	}
	return ""
}

// GetPage returns the page parameter
func (sr *SearchRequest) GetPage() string {
	if sr.Page != nil {
		return *sr.Page
	}
	return ""
}

// GetOrder returns the order parameter
func (sr *SearchRequest) GetOrder() string {
	if sr.Order != nil {
		return *sr.Order
	}
	return ""
}

// GetOrderBy returns the order by parameter
func (sr *SearchRequest) GetOrderBy() string {
	if sr.OrderBy != nil {
		return *sr.OrderBy
	}
	return ""
}
