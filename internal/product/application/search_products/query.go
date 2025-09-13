package search_products

import (
	"context"

	"github.com/jperdior/chatbot-kit/application/query"
	kitDomain "github.com/jperdior/chatbot-kit/domain"
)

const SearchProductsQueryType query.Type = "search_products"

// SearchProductsQuery represents a query to search for products
type SearchProductsQuery struct {
	// Required authentication parameters
	BaseURL        string `json:"base_url"`
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`

	// Optional search parameters
	Search      string `json:"search,omitempty"`
	Category    string `json:"category,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Status      string `json:"status,omitempty"`
	ProductType string `json:"type,omitempty"`
	Featured    string `json:"featured,omitempty"`
	OnSale      string `json:"on_sale,omitempty"`
	MinPrice    string `json:"min_price,omitempty"`
	MaxPrice    string `json:"max_price,omitempty"`
	StockStatus string `json:"stock_status,omitempty"`
	PerPage     string `json:"per_page,omitempty"`
	Page        string `json:"page,omitempty"`
	Order       string `json:"order,omitempty"`
	OrderBy     string `json:"orderby,omitempty"`
}

// NewSearchProductsQuery creates a new SearchProductsQuery
func NewSearchProductsQuery(baseURL, consumerKey, consumerSecret string) *SearchProductsQuery {
	return &SearchProductsQuery{
		BaseURL:        baseURL,
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	}
}

// Type returns the query type
func (q *SearchProductsQuery) Type() query.Type {
	return SearchProductsQueryType
}

// SearchProductsQueryHandler handles the SearchProductsQuery
type SearchProductsQueryHandler struct {
	ProductSearcher *ProductSearcher
}

// NewSearchProductsQueryHandler creates a new SearchProductsQueryHandler
func NewSearchProductsQueryHandler(productSearcher *ProductSearcher) *SearchProductsQueryHandler {
	return &SearchProductsQueryHandler{
		ProductSearcher: productSearcher,
	}
}

// Handle processes the SearchProductsQuery
func (h *SearchProductsQueryHandler) Handle(ctx context.Context, query query.Query) (interface{}, error) {
	q, ok := query.(*SearchProductsQuery)
	if !ok {
		return nil, kitDomain.NewDomainError("INVALID_QUERY", "unexpected query type for search products")
	}

	// Convert query to search request
	request := h.queryToRequest(q)

	// Execute the search
	response, err := h.ProductSearcher.Execute(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// queryToRequest converts SearchProductsQuery to SearchRequest
func (h *SearchProductsQueryHandler) queryToRequest(q *SearchProductsQuery) *SearchRequest {
	request := NewSearchRequest(q.BaseURL, q.ConsumerKey, q.ConsumerSecret)

	if q.Search != "" {
		request.SetSearch(q.Search)
	}
	if q.Category != "" {
		request.SetCategory(q.Category)
	}
	if q.Tag != "" {
		request.SetTag(q.Tag)
	}
	if q.Status != "" {
		request.SetStatus(q.Status)
	}
	if q.ProductType != "" {
		request.SetType(q.ProductType)
	}
	if q.Featured != "" {
		request.SetFeatured(q.Featured)
	}
	if q.OnSale != "" {
		request.SetOnSale(q.OnSale)
	}
	if q.MinPrice != "" || q.MaxPrice != "" {
		request.SetPriceRange(q.MinPrice, q.MaxPrice)
	}
	if q.StockStatus != "" {
		request.SetStockStatus(q.StockStatus)
	}
	if q.Page != "" || q.PerPage != "" {
		request.SetPagination(q.Page, q.PerPage)
	}
	if q.OrderBy != "" || q.Order != "" {
		request.SetSorting(q.OrderBy, q.Order)
	}

	return request
}
