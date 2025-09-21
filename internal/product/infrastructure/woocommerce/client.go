package woocommerce

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"woocommerce-mcp/internal/product/domain"
)

// Config represents WooCommerce API configuration
type Config struct {
	BaseURL        string
	ConsumerKey    string
	ConsumerSecret string
	Timeout        time.Duration
}

// NewConfig creates a new WooCommerce configuration
func NewConfig(baseURL, consumerKey, consumerSecret string) *Config {
	return &Config{
		BaseURL:        strings.TrimSuffix(baseURL, "/"),
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		Timeout:        30 * time.Second,
	}
}

// Client represents a WooCommerce API client
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new WooCommerce client
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// SearchProducts searches for products using the WooCommerce API
func (c *Client) SearchProducts(ctx context.Context, criteria *domain.SearchCriteria) ([]*domain.Product, error) {
	// Build the API endpoint URL
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products", c.config.BaseURL)

	// Parse base URL
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, domain.NewConnectionError(endpoint, fmt.Sprintf("invalid base URL: %v", err))
	}

	// Build query parameters
	query := u.Query()
	c.addAuthParams(query)
	c.addSearchParams(query, criteria)

	u.RawQuery = query.Encode()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domain.NewConnectionError(u.String(), fmt.Sprintf("HTTP request failed: %v", err))
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleAPIError(resp.StatusCode, body)
	}

	// Parse JSON response
	var apiProducts []APIProduct
	if err := json.Unmarshal(body, &apiProducts); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Convert API products to domain products
	products := make([]*domain.Product, len(apiProducts))
	for i, apiProduct := range apiProducts {
		domainProduct, err := c.apiProductToDomain(&apiProduct)
		if err != nil {
			return nil, fmt.Errorf("failed to convert product %d: %w", apiProduct.ID, err)
		}
		products[i] = domainProduct
	}

	return products, nil
}

// CountProducts counts products matching the criteria
func (c *Client) CountProducts(ctx context.Context, criteria *domain.SearchCriteria) (int64, error) {
	// For WooCommerce API, we need to make a HEAD request or parse headers
	// Since WooCommerce doesn't provide a direct count endpoint, we'll use the X-WP-Total header
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products", c.config.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return 0, domain.NewConnectionError(endpoint, fmt.Sprintf("invalid base URL: %v", err))
	}

	// Build query parameters (same as search but we only need the count)
	query := u.Query()
	c.addAuthParams(query)
	c.addSearchParams(query, criteria)

	// Set per_page to 1 to minimize data transfer when we only need the count
	query.Set("per_page", "1")

	u.RawQuery = query.Encode()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "HEAD", u.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Make HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, domain.NewConnectionError(u.String(), fmt.Sprintf("HTTP request failed: %v", err))
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return 0, c.handleAPIError(resp.StatusCode, nil)
	}

	// Get total count from header
	totalHeader := resp.Header.Get("X-WP-Total")
	if totalHeader == "" {
		// Fallback: make a GET request and count manually
		return c.countProductsFallback(ctx, criteria)
	}

	total, err := strconv.ParseInt(totalHeader, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse total count: %w", err)
	}

	return total, nil
}

// countProductsFallback is a fallback method to count products when headers are not available
func (c *Client) countProductsFallback(ctx context.Context, criteria *domain.SearchCriteria) (int64, error) {
	// Make a request with a reasonable per_page to get actual results
	tempCriteria := *criteria
	tempCriteria.PerPage = 100 // Get up to 100 products to count
	tempCriteria.Page = 1

	products, err := c.SearchProducts(ctx, &tempCriteria)
	if err != nil {
		return 0, err
	}

	// Return the actual count of products found
	// Note: This is a simplified approach that works for small result sets
	// For large result sets, we'd need to implement proper pagination counting
	return int64(len(products)), nil
}

// addAuthParams adds authentication parameters to the query
func (c *Client) addAuthParams(query url.Values) {
	query.Set("consumer_key", c.config.ConsumerKey)
	query.Set("consumer_secret", c.config.ConsumerSecret)
}

// addSearchParams adds search parameters to the query
func (c *Client) addSearchParams(query url.Values, criteria *domain.SearchCriteria) {
	if criteria.Search != "" {
		query.Set("search", criteria.Search)
	}
	if criteria.Category != "" {
		query.Set("category", criteria.Category)
	}
	if criteria.Tag != "" {
		query.Set("tag", criteria.Tag)
	}
	if criteria.Status != "" {
		query.Set("status", string(criteria.Status))
	}
	if criteria.Type != "" {
		query.Set("type", string(criteria.Type))
	}
	if criteria.Featured != nil {
		query.Set("featured", strconv.FormatBool(*criteria.Featured))
	}
	if criteria.OnSale != nil {
		query.Set("on_sale", strconv.FormatBool(*criteria.OnSale))
	}
	if criteria.MinPrice != nil {
		query.Set("min_price", fmt.Sprintf("%.2f", criteria.MinPrice.Amount()))
	}
	if criteria.MaxPrice != nil {
		query.Set("max_price", fmt.Sprintf("%.2f", criteria.MaxPrice.Amount()))
	}
	if criteria.StockStatus != "" {
		query.Set("stock_status", string(criteria.StockStatus))
	}

	// Pagination
	query.Set("per_page", strconv.Itoa(criteria.PerPage))
	query.Set("page", strconv.Itoa(criteria.Page))

	// Sorting
	if criteria.OrderBy != "" {
		query.Set("orderby", criteria.OrderBy)
	}
	if criteria.Order != "" {
		query.Set("order", criteria.Order)
	}
}

// handleAPIError handles API errors and converts them to domain errors
func (c *Client) handleAPIError(statusCode int, body []byte) error {
	message := string(body)
	if len(body) == 0 {
		message = http.StatusText(statusCode)
	}

	// Try to parse error response for more details
	var apiError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &apiError); err == nil {
			if apiError.Message != "" {
				message = apiError.Message
			}
		}
	}

	return domain.NewWooCommerceAPIError(statusCode, message, apiError.Code)
}

// apiProductToDomain converts an API product to a domain product
func (c *Client) apiProductToDomain(apiProduct *APIProduct) (*domain.Product, error) {
	// Create product ID
	productID, err := domain.NewProductID(apiProduct.ID)
	if err != nil {
		return nil, err
	}

	// Create domain product
	product := domain.NewProduct(productID, apiProduct.Name)

	// Set basic fields
	product.Slug = apiProduct.Slug
	product.Permalink = apiProduct.Permalink
	product.Description = apiProduct.Description
	product.ShortDescription = apiProduct.ShortDescription
	product.SKU = apiProduct.SKU
	product.Featured = apiProduct.Featured
	product.CatalogVisibility = apiProduct.CatalogVisibility
	product.OnSale = apiProduct.OnSale
	product.Purchasable = apiProduct.Purchasable
	product.TotalSales = apiProduct.TotalSales
	product.Virtual = apiProduct.Virtual
	product.Downloadable = apiProduct.Downloadable
	product.ExternalURL = apiProduct.ExternalURL
	product.ButtonText = apiProduct.ButtonText
	product.TaxStatus = apiProduct.TaxStatus
	product.TaxClass = apiProduct.TaxClass
	product.ManageStock = apiProduct.ManageStock
	product.StockQuantity = apiProduct.StockQuantity
	product.Backorders = apiProduct.Backorders
	product.BackordersAllowed = apiProduct.BackordersAllowed
	product.Backordered = apiProduct.Backordered
	product.Weight = apiProduct.Weight
	product.ShippingRequired = apiProduct.ShippingRequired
	product.ShippingTaxable = apiProduct.ShippingTaxable
	product.ShippingClass = apiProduct.ShippingClass
	product.ShippingClassID = apiProduct.ShippingClassID
	product.ReviewsAllowed = apiProduct.ReviewsAllowed
	product.AverageRating = apiProduct.AverageRating
	product.RatingCount = apiProduct.RatingCount
	product.RelatedIDs = apiProduct.RelatedIDs
	product.UpsellIDs = apiProduct.UpsellIDs
	product.CrossSellIDs = apiProduct.CrossSellIDs
	product.ParentID = apiProduct.ParentID
	product.PurchaseNote = apiProduct.PurchaseNote
	product.Variations = apiProduct.Variations
	product.GroupedProducts = apiProduct.GroupedProducts
	product.MenuOrder = apiProduct.MenuOrder

	// Parse dates
	if apiProduct.DateCreated != "" {
		if dateCreated, err := time.Parse("2006-01-02T15:04:05", apiProduct.DateCreated); err == nil {
			product.DateCreated = dateCreated
		}
	}
	if apiProduct.DateModified != "" {
		if dateModified, err := time.Parse("2006-01-02T15:04:05", apiProduct.DateModified); err == nil {
			product.DateModified = dateModified
		}
	}

	// Set product type
	if apiProduct.Type != "" {
		productType := domain.ProductType(apiProduct.Type)
		if productType.IsValid() {
			product.Type = productType
		}
	}

	// Set product status
	if apiProduct.Status != "" {
		status := domain.ProductStatus(apiProduct.Status)
		if status.IsValid() {
			product.Status = status
		}
	}

	// Set stock status
	if apiProduct.StockStatus != "" {
		stockStatus := domain.StockStatus(apiProduct.StockStatus)
		if stockStatus.IsValid() {
			product.StockStatus = stockStatus
		}
	}

	// Convert prices
	if apiProduct.Price != "" {
		if price, err := domain.NewMoneyFromString(apiProduct.Price, "USD"); err == nil {
			product.Price = price
		}
	}
	if apiProduct.RegularPrice != "" {
		if regularPrice, err := domain.NewMoneyFromString(apiProduct.RegularPrice, "USD"); err == nil {
			product.RegularPrice = regularPrice
		}
	}
	if apiProduct.SalePrice != "" {
		if salePrice, err := domain.NewMoneyFromString(apiProduct.SalePrice, "USD"); err == nil {
			product.SalePrice = salePrice
		}
	}

	// Convert dimensions
	if !apiProduct.Dimensions.IsEmpty() {
		product.Dimensions = domain.NewDimensions(
			apiProduct.Dimensions.Length,
			apiProduct.Dimensions.Width,
			apiProduct.Dimensions.Height,
		)
	}

	// Convert categories
	for _, apiCategory := range apiProduct.Categories {
		category := domain.NewCategory(apiCategory.ID, apiCategory.Name, apiCategory.Slug)
		product.Categories = append(product.Categories, category)
	}

	// Convert tags
	for _, apiTag := range apiProduct.Tags {
		tag := domain.NewTag(apiTag.ID, apiTag.Name, apiTag.Slug)
		product.Tags = append(product.Tags, tag)
	}

	// Convert images
	for _, apiImage := range apiProduct.Images {
		image := domain.NewImage(apiImage.ID, apiImage.Src, apiImage.Name, apiImage.Alt, apiImage.Position)
		image.DateCreated = apiImage.DateCreated
		image.DateModified = apiImage.DateModified
		product.Images = append(product.Images, image)
	}

	// Convert attributes
	for _, apiAttribute := range apiProduct.Attributes {
		attribute := domain.NewAttribute(apiAttribute.ID, apiAttribute.Name, apiAttribute.Options)
		attribute.Position = apiAttribute.Position
		attribute.Visible = apiAttribute.Visible
		attribute.Variation = apiAttribute.Variation
		product.Attributes = append(product.Attributes, attribute)
	}

	// Convert default attributes
	for _, apiDefaultAttr := range apiProduct.DefaultAttributes {
		defaultAttr := domain.NewDefaultAttribute(apiDefaultAttr.ID, apiDefaultAttr.Name, apiDefaultAttr.Option)
		product.DefaultAttributes = append(product.DefaultAttributes, defaultAttr)
	}

	// Convert metadata
	for _, apiMetaData := range apiProduct.MetaData {
		metaData := domain.NewMetaData(apiMetaData.ID, apiMetaData.Key, apiMetaData.Value)
		product.MetaData = append(product.MetaData, metaData)
	}

	return product, nil
}
