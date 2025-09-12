package main

import "strconv"

// MCP Protocol types
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// WooCommerce API types
type WooCommerceConfig struct {
	BaseURL        string `json:"base_url"`
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
}

type Product struct {
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	Slug              string             `json:"slug"`
	Permalink         string             `json:"permalink"`
	DateCreated       string             `json:"date_created"`
	DateModified      string             `json:"date_modified"`
	Type              string             `json:"type"`
	Status            string             `json:"status"`
	Featured          bool               `json:"featured"`
	CatalogVisibility string             `json:"catalog_visibility"`
	Description       string             `json:"description"`
	ShortDescription  string             `json:"short_description"`
	SKU               string             `json:"sku"`
	Price             string             `json:"price"`
	RegularPrice      string             `json:"regular_price"`
	SalePrice         string             `json:"sale_price"`
	OnSale            bool               `json:"on_sale"`
	Purchasable       bool               `json:"purchasable"`
	TotalSales        int                `json:"total_sales"`
	Virtual           bool               `json:"virtual"`
	Downloadable      bool               `json:"downloadable"`
	ExternalURL       string             `json:"external_url"`
	ButtonText        string             `json:"button_text"`
	TaxStatus         string             `json:"tax_status"`
	TaxClass          string             `json:"tax_class"`
	ManageStock       bool               `json:"manage_stock"`
	StockQuantity     *int               `json:"stock_quantity"`
	StockStatus       string             `json:"stock_status"`
	Backorders        string             `json:"backorders"`
	BackordersAllowed bool               `json:"backorders_allowed"`
	Backordered       bool               `json:"backordered"`
	Weight            string             `json:"weight"`
	Dimensions        Dimensions         `json:"dimensions"`
	ShippingRequired  bool               `json:"shipping_required"`
	ShippingTaxable   bool               `json:"shipping_taxable"`
	ShippingClass     string             `json:"shipping_class"`
	ShippingClassID   int                `json:"shipping_class_id"`
	ReviewsAllowed    bool               `json:"reviews_allowed"`
	AverageRating     string             `json:"average_rating"`
	RatingCount       int                `json:"rating_count"`
	RelatedIDs        []int              `json:"related_ids"`
	UpsellIDs         []int              `json:"upsell_ids"`
	CrossSellIDs      []int              `json:"cross_sell_ids"`
	ParentID          int                `json:"parent_id"`
	PurchaseNote      string             `json:"purchase_note"`
	Categories        []Category         `json:"categories"`
	Tags              []Tag              `json:"tags"`
	Images            []Image            `json:"images"`
	Attributes        []Attribute        `json:"attributes"`
	DefaultAttributes []DefaultAttribute `json:"default_attributes"`
	Variations        []int              `json:"variations"`
	GroupedProducts   []int              `json:"grouped_products"`
	MenuOrder         int                `json:"menu_order"`
	MetaData          []MetaData         `json:"meta_data"`
}

type Dimensions struct {
	Length string `json:"length"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Image struct {
	ID           int    `json:"id"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
	Src          string `json:"src"`
	Name         string `json:"name"`
	Alt          string `json:"alt"`
	Position     int    `json:"position"`
}

type Attribute struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Visible   bool     `json:"visible"`
	Variation bool     `json:"variation"`
	Options   []string `json:"options"`
}

type DefaultAttribute struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Option string `json:"option"`
}

type MetaData struct {
	ID    int         `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Search Products Request struct
type SearchProductsRequest struct {
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

// ToSearchParams converts the request to a map for the WooCommerce client
func (req *SearchProductsRequest) ToSearchParams() map[string]interface{} {
	params := make(map[string]interface{})

	// Add optional parameters only if they are provided
	if req.Search != nil && *req.Search != "" {
		params["search"] = *req.Search
	}
	if req.Category != nil && *req.Category != "" {
		params["category"] = *req.Category
	}
	if req.Tag != nil && *req.Tag != "" {
		params["tag"] = *req.Tag
	}
	if req.Status != nil && *req.Status != "" {
		params["status"] = *req.Status
	}
	if req.Type != nil && *req.Type != "" {
		params["type"] = *req.Type
	}
	if req.Featured != nil && *req.Featured != "" {
		params["featured"] = *req.Featured
	}
	if req.OnSale != nil && *req.OnSale != "" {
		params["on_sale"] = *req.OnSale
	}
	if req.MinPrice != nil && *req.MinPrice != "" {
		params["min_price"] = *req.MinPrice
	}
	if req.MaxPrice != nil && *req.MaxPrice != "" {
		params["max_price"] = *req.MaxPrice
	}
	if req.StockStatus != nil && *req.StockStatus != "" {
		params["stock_status"] = *req.StockStatus
	}
	if req.PerPage != nil && *req.PerPage != "" {
		params["per_page"] = *req.PerPage
	} else {
		// Default per_page if not provided
		params["per_page"] = "10"
	}
	if req.Page != nil && *req.Page != "" {
		params["page"] = *req.Page
	}
	if req.Order != nil && *req.Order != "" {
		params["order"] = *req.Order
	}
	if req.OrderBy != nil && *req.OrderBy != "" {
		params["orderby"] = *req.OrderBy
	}

	return params
}

// Validate performs additional validation on the request
func (req *SearchProductsRequest) Validate() error {
	// Validate per_page limit
	if req.PerPage != nil {
		if perPage, err := strconv.Atoi(*req.PerPage); err == nil && perPage > 100 {
			// Automatically cap at 100
			cappedPerPage := "100"
			req.PerPage = &cappedPerPage
		}
	}

	return nil
}
