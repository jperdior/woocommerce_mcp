package search_products

// SearchResponse represents the response from a product search
type SearchResponse struct {
	Products    []*ProductDTO `json:"products"`
	TotalCount  int           `json:"total_count"`
	CurrentPage int           `json:"current_page"`
	PerPage     int           `json:"per_page"`
	TotalPages  int           `json:"total_pages"`
	HasNext     bool          `json:"has_next"`
	HasPrev     bool          `json:"has_prev"`
}

// ProductDTO represents a product data transfer object
type ProductDTO struct {
	ID                int                    `json:"id"`
	Name              string                 `json:"name"`
	Slug              string                 `json:"slug"`
	Permalink         string                 `json:"permalink"`
	DateCreated       string                 `json:"date_created"`
	DateModified      string                 `json:"date_modified"`
	Type              string                 `json:"type"`
	Status            string                 `json:"status"`
	Featured          bool                   `json:"featured"`
	CatalogVisibility string                 `json:"catalog_visibility"`
	Description       string                 `json:"description"`
	ShortDescription  string                 `json:"short_description"`
	SKU               string                 `json:"sku"`
	Price             string                 `json:"price"`
	RegularPrice      string                 `json:"regular_price"`
	SalePrice         string                 `json:"sale_price"`
	OnSale            bool                   `json:"on_sale"`
	Purchasable       bool                   `json:"purchasable"`
	TotalSales        int                    `json:"total_sales"`
	Virtual           bool                   `json:"virtual"`
	Downloadable      bool                   `json:"downloadable"`
	ExternalURL       string                 `json:"external_url"`
	ButtonText        string                 `json:"button_text"`
	TaxStatus         string                 `json:"tax_status"`
	TaxClass          string                 `json:"tax_class"`
	ManageStock       bool                   `json:"manage_stock"`
	StockQuantity     *int                   `json:"stock_quantity"`
	StockStatus       string                 `json:"stock_status"`
	Backorders        string                 `json:"backorders"`
	BackordersAllowed bool                   `json:"backorders_allowed"`
	Backordered       bool                   `json:"backordered"`
	Weight            string                 `json:"weight"`
	Dimensions        *DimensionsDTO         `json:"dimensions"`
	ShippingRequired  bool                   `json:"shipping_required"`
	ShippingTaxable   bool                   `json:"shipping_taxable"`
	ShippingClass     string                 `json:"shipping_class"`
	ShippingClassID   int                    `json:"shipping_class_id"`
	ReviewsAllowed    bool                   `json:"reviews_allowed"`
	AverageRating     string                 `json:"average_rating"`
	RatingCount       int                    `json:"rating_count"`
	RelatedIDs        []int                  `json:"related_ids"`
	UpsellIDs         []int                  `json:"upsell_ids"`
	CrossSellIDs      []int                  `json:"cross_sell_ids"`
	ParentID          int                    `json:"parent_id"`
	PurchaseNote      string                 `json:"purchase_note"`
	Categories        []*CategoryDTO         `json:"categories"`
	Tags              []*TagDTO              `json:"tags"`
	Images            []*ImageDTO            `json:"images"`
	Attributes        []*AttributeDTO        `json:"attributes"`
	DefaultAttributes []*DefaultAttributeDTO `json:"default_attributes"`
	Variations        []int                  `json:"variations"`
	GroupedProducts   []int                  `json:"grouped_products"`
	MenuOrder         int                    `json:"menu_order"`
	MetaData          []*MetaDataDTO         `json:"meta_data"`
}

// DimensionsDTO represents product dimensions
type DimensionsDTO struct {
	Length string `json:"length"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

// CategoryDTO represents a product category
type CategoryDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// TagDTO represents a product tag
type TagDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ImageDTO represents a product image
type ImageDTO struct {
	ID           int    `json:"id"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
	Src          string `json:"src"`
	Name         string `json:"name"`
	Alt          string `json:"alt"`
	Position     int    `json:"position"`
}

// AttributeDTO represents a product attribute
type AttributeDTO struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Visible   bool     `json:"visible"`
	Variation bool     `json:"variation"`
	Options   []string `json:"options"`
}

// DefaultAttributeDTO represents a default product attribute
type DefaultAttributeDTO struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Option string `json:"option"`
}

// MetaDataDTO represents product metadata
type MetaDataDTO struct {
	ID    int         `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// NewSearchResponse creates a new SearchResponse
func NewSearchResponse(products []*ProductDTO, totalCount, currentPage, perPage int) *SearchResponse {
	totalPages := (totalCount + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}

	return &SearchResponse{
		Products:    products,
		TotalCount:  totalCount,
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     currentPage < totalPages,
		HasPrev:     currentPage > 1,
	}
}

// IsEmpty checks if the response has no products
func (sr *SearchResponse) IsEmpty() bool {
	return len(sr.Products) == 0
}

// GetProductCount returns the number of products in the response
func (sr *SearchResponse) GetProductCount() int {
	return len(sr.Products)
}

// GetProductByID finds a product by ID in the response
func (sr *SearchResponse) GetProductByID(id int) *ProductDTO {
	for _, product := range sr.Products {
		if product.ID == id {
			return product
		}
	}
	return nil
}

// GetProductsBySKU finds products by SKU in the response
func (sr *SearchResponse) GetProductsBySKU(sku string) []*ProductDTO {
	var products []*ProductDTO
	for _, product := range sr.Products {
		if product.SKU == sku {
			products = append(products, product)
		}
	}
	return products
}

// GetFeaturedProducts returns only featured products from the response
func (sr *SearchResponse) GetFeaturedProducts() []*ProductDTO {
	var featuredProducts []*ProductDTO
	for _, product := range sr.Products {
		if product.Featured {
			featuredProducts = append(featuredProducts, product)
		}
	}
	return featuredProducts
}

// GetProductsOnSale returns only products on sale from the response
func (sr *SearchResponse) GetProductsOnSale() []*ProductDTO {
	var saleProducts []*ProductDTO
	for _, product := range sr.Products {
		if product.OnSale {
			saleProducts = append(saleProducts, product)
		}
	}
	return saleProducts
}

// GetProductsByCategory returns products filtered by category
func (sr *SearchResponse) GetProductsByCategory(categoryID int) []*ProductDTO {
	var categoryProducts []*ProductDTO
	for _, product := range sr.Products {
		for _, category := range product.Categories {
			if category.ID == categoryID {
				categoryProducts = append(categoryProducts, product)
				break
			}
		}
	}
	return categoryProducts
}

// GetProductsByTag returns products filtered by tag
func (sr *SearchResponse) GetProductsByTag(tagID int) []*ProductDTO {
	var tagProducts []*ProductDTO
	for _, product := range sr.Products {
		for _, tag := range product.Tags {
			if tag.ID == tagID {
				tagProducts = append(tagProducts, product)
				break
			}
		}
	}
	return tagProducts
}

// GetProductsByStatus returns products filtered by status
func (sr *SearchResponse) GetProductsByStatus(status string) []*ProductDTO {
	var statusProducts []*ProductDTO
	for _, product := range sr.Products {
		if product.Status == status {
			statusProducts = append(statusProducts, product)
		}
	}
	return statusProducts
}

// GetProductsByType returns products filtered by type
func (sr *SearchResponse) GetProductsByType(productType string) []*ProductDTO {
	var typeProducts []*ProductDTO
	for _, product := range sr.Products {
		if product.Type == productType {
			typeProducts = append(typeProducts, product)
		}
	}
	return typeProducts
}
