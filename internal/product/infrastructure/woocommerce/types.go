package woocommerce

// APIProduct represents a product as returned by the WooCommerce API
type APIProduct struct {
	ID                int                   `json:"id"`
	Name              string                `json:"name"`
	Slug              string                `json:"slug"`
	Permalink         string                `json:"permalink"`
	DateCreated       string                `json:"date_created"`
	DateModified      string                `json:"date_modified"`
	Type              string                `json:"type"`
	Status            string                `json:"status"`
	Featured          bool                  `json:"featured"`
	CatalogVisibility string                `json:"catalog_visibility"`
	Description       string                `json:"description"`
	ShortDescription  string                `json:"short_description"`
	SKU               string                `json:"sku"`
	Price             string                `json:"price"`
	RegularPrice      string                `json:"regular_price"`
	SalePrice         string                `json:"sale_price"`
	OnSale            bool                  `json:"on_sale"`
	Purchasable       bool                  `json:"purchasable"`
	TotalSales        int                   `json:"total_sales"`
	Virtual           bool                  `json:"virtual"`
	Downloadable      bool                  `json:"downloadable"`
	ExternalURL       string                `json:"external_url"`
	ButtonText        string                `json:"button_text"`
	TaxStatus         string                `json:"tax_status"`
	TaxClass          string                `json:"tax_class"`
	ManageStock       bool                  `json:"manage_stock"`
	StockQuantity     *int                  `json:"stock_quantity"`
	StockStatus       string                `json:"stock_status"`
	Backorders        string                `json:"backorders"`
	BackordersAllowed bool                  `json:"backorders_allowed"`
	Backordered       bool                  `json:"backordered"`
	Weight            string                `json:"weight"`
	Dimensions        APIDimensions         `json:"dimensions"`
	ShippingRequired  bool                  `json:"shipping_required"`
	ShippingTaxable   bool                  `json:"shipping_taxable"`
	ShippingClass     string                `json:"shipping_class"`
	ShippingClassID   int                   `json:"shipping_class_id"`
	ReviewsAllowed    bool                  `json:"reviews_allowed"`
	AverageRating     string                `json:"average_rating"`
	RatingCount       int                   `json:"rating_count"`
	RelatedIDs        []int                 `json:"related_ids"`
	UpsellIDs         []int                 `json:"upsell_ids"`
	CrossSellIDs      []int                 `json:"cross_sell_ids"`
	ParentID          int                   `json:"parent_id"`
	PurchaseNote      string                `json:"purchase_note"`
	Categories        []APICategory         `json:"categories"`
	Tags              []APITag              `json:"tags"`
	Images            []APIImage            `json:"images"`
	Attributes        []APIAttribute        `json:"attributes"`
	DefaultAttributes []APIDefaultAttribute `json:"default_attributes"`
	Variations        []int                 `json:"variations"`
	GroupedProducts   []int                 `json:"grouped_products"`
	MenuOrder         int                   `json:"menu_order"`
	MetaData          []APIMetaData         `json:"meta_data"`
}

// APIDimensions represents product dimensions from the API
type APIDimensions struct {
	Length string `json:"length"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

// IsEmpty checks if all dimensions are empty
func (d APIDimensions) IsEmpty() bool {
	return d.Length == "" && d.Width == "" && d.Height == ""
}

// APICategory represents a product category from the API
type APICategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// APITag represents a product tag from the API
type APITag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// APIImage represents a product image from the API
type APIImage struct {
	ID           int    `json:"id"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
	Src          string `json:"src"`
	Name         string `json:"name"`
	Alt          string `json:"alt"`
	Position     int    `json:"position"`
}

// APIAttribute represents a product attribute from the API
type APIAttribute struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Visible   bool     `json:"visible"`
	Variation bool     `json:"variation"`
	Options   []string `json:"options"`
}

// APIDefaultAttribute represents a default product attribute from the API
type APIDefaultAttribute struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Option string `json:"option"`
}

// APIMetaData represents product metadata from the API
type APIMetaData struct {
	ID    int         `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// APIErrorResponse represents an error response from the WooCommerce API
type APIErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Status int `json:"status"`
	} `json:"data"`
}

// APIListResponse represents a list response with pagination info
type APIListResponse struct {
	Data       []APIProduct `json:"data"`
	TotalCount int          `json:"total_count"`
	TotalPages int          `json:"total_pages"`
	Page       int          `json:"page"`
	PerPage    int          `json:"per_page"`
}

// SearchParams represents search parameters for the WooCommerce API
type SearchParams struct {
	Search      string `json:"search,omitempty"`
	Category    string `json:"category,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Status      string `json:"status,omitempty"`
	Type        string `json:"type,omitempty"`
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

// ToMap converts SearchParams to a map for URL encoding
func (sp *SearchParams) ToMap() map[string]string {
	params := make(map[string]string)

	if sp.Search != "" {
		params["search"] = sp.Search
	}
	if sp.Category != "" {
		params["category"] = sp.Category
	}
	if sp.Tag != "" {
		params["tag"] = sp.Tag
	}
	if sp.Status != "" {
		params["status"] = sp.Status
	}
	if sp.Type != "" {
		params["type"] = sp.Type
	}
	if sp.Featured != "" {
		params["featured"] = sp.Featured
	}
	if sp.OnSale != "" {
		params["on_sale"] = sp.OnSale
	}
	if sp.MinPrice != "" {
		params["min_price"] = sp.MinPrice
	}
	if sp.MaxPrice != "" {
		params["max_price"] = sp.MaxPrice
	}
	if sp.StockStatus != "" {
		params["stock_status"] = sp.StockStatus
	}
	if sp.PerPage != "" {
		params["per_page"] = sp.PerPage
	} else {
		params["per_page"] = "10"
	}
	if sp.Page != "" {
		params["page"] = sp.Page
	}
	if sp.Order != "" {
		params["order"] = sp.Order
	}
	if sp.OrderBy != "" {
		params["orderby"] = sp.OrderBy
	}

	return params
}
