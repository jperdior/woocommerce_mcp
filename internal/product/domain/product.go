package domain

import (
	"time"
	"woocommerce-mcp/kit/domain"
)

// Product represents a WooCommerce product entity
type Product struct {
	ID                *ProductID          `json:"id"`
	Name              string              `json:"name"`
	Slug              string              `json:"slug"`
	Permalink         string              `json:"permalink"`
	DateCreated       time.Time           `json:"date_created"`
	DateModified      time.Time           `json:"date_modified"`
	Type              ProductType         `json:"type"`
	Status            ProductStatus       `json:"status"`
	Featured          bool                `json:"featured"`
	CatalogVisibility string              `json:"catalog_visibility"`
	Description       string              `json:"description"`
	ShortDescription  string              `json:"short_description"`
	SKU               string              `json:"sku"`
	Price             *Money              `json:"price"`
	RegularPrice      *Money              `json:"regular_price"`
	SalePrice         *Money              `json:"sale_price"`
	OnSale            bool                `json:"on_sale"`
	Purchasable       bool                `json:"purchasable"`
	TotalSales        int                 `json:"total_sales"`
	Virtual           bool                `json:"virtual"`
	Downloadable      bool                `json:"downloadable"`
	ExternalURL       string              `json:"external_url"`
	ButtonText        string              `json:"button_text"`
	TaxStatus         string              `json:"tax_status"`
	TaxClass          string              `json:"tax_class"`
	ManageStock       bool                `json:"manage_stock"`
	StockQuantity     *int                `json:"stock_quantity"`
	StockStatus       StockStatus         `json:"stock_status"`
	Backorders        string              `json:"backorders"`
	BackordersAllowed bool                `json:"backorders_allowed"`
	Backordered       bool                `json:"backordered"`
	Weight            string              `json:"weight"`
	Dimensions        *Dimensions         `json:"dimensions"`
	ShippingRequired  bool                `json:"shipping_required"`
	ShippingTaxable   bool                `json:"shipping_taxable"`
	ShippingClass     string              `json:"shipping_class"`
	ShippingClassID   int                 `json:"shipping_class_id"`
	ReviewsAllowed    bool                `json:"reviews_allowed"`
	AverageRating     string              `json:"average_rating"`
	RatingCount       int                 `json:"rating_count"`
	RelatedIDs        []int               `json:"related_ids"`
	UpsellIDs         []int               `json:"upsell_ids"`
	CrossSellIDs      []int               `json:"cross_sell_ids"`
	ParentID          int                 `json:"parent_id"`
	PurchaseNote      string              `json:"purchase_note"`
	Categories        []*Category         `json:"categories"`
	Tags              []*Tag              `json:"tags"`
	Images            []*Image            `json:"images"`
	Attributes        []*Attribute        `json:"attributes"`
	DefaultAttributes []*DefaultAttribute `json:"default_attributes"`
	Variations        []int               `json:"variations"`
	GroupedProducts   []int               `json:"grouped_products"`
	MenuOrder         int                 `json:"menu_order"`
	MetaData          []*MetaData         `json:"meta_data"`
}

// NewProduct creates a new product instance
func NewProduct(id *ProductID, name string) *Product {
	return &Product{
		ID:           id,
		Name:         name,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
		Type:         ProductTypeSimple,
		Status:       ProductStatusDraft,
		Categories:   make([]*Category, 0),
		Tags:         make([]*Tag, 0),
		Images:       make([]*Image, 0),
		Attributes:   make([]*Attribute, 0),
		MetaData:     make([]*MetaData, 0),
	}
}

// IsValid validates the product
func (p *Product) IsValid() error {
	if p.ID == nil {
		return domain.NewValidationError("product ID is required")
	}
	if p.Name == "" {
		return domain.NewValidationError("product name is required")
	}
	return nil
}

// UpdateName updates the product name
func (p *Product) UpdateName(name string) error {
	if name == "" {
		return domain.NewValidationError("product name cannot be empty")
	}
	p.Name = name
	p.DateModified = time.Now()
	return nil
}

// UpdatePrice updates the product price
func (p *Product) UpdatePrice(price *Money) error {
	if price == nil {
		return domain.NewValidationError("price cannot be nil")
	}
	p.Price = price
	p.DateModified = time.Now()
	return nil
}

// SetFeatured sets the product as featured or not
func (p *Product) SetFeatured(featured bool) {
	p.Featured = featured
	p.DateModified = time.Now()
}

// SetStatus updates the product status
func (p *Product) SetStatus(status ProductStatus) {
	p.Status = status
	p.DateModified = time.Now()
}

// AddCategory adds a category to the product
func (p *Product) AddCategory(category *Category) {
	if category == nil {
		return
	}

	// Check if category already exists
	for _, existingCategory := range p.Categories {
		if existingCategory.ID == category.ID {
			return
		}
	}

	p.Categories = append(p.Categories, category)
	p.DateModified = time.Now()
}

// RemoveCategory removes a category from the product
func (p *Product) RemoveCategory(categoryID int) {
	for i, category := range p.Categories {
		if category.ID == categoryID {
			p.Categories = append(p.Categories[:i], p.Categories[i+1:]...)
			p.DateModified = time.Now()
			break
		}
	}
}

// AddTag adds a tag to the product
func (p *Product) AddTag(tag *Tag) {
	if tag == nil {
		return
	}

	// Check if tag already exists
	for _, existingTag := range p.Tags {
		if existingTag.ID == tag.ID {
			return
		}
	}

	p.Tags = append(p.Tags, tag)
	p.DateModified = time.Now()
}

// RemoveTag removes a tag from the product
func (p *Product) RemoveTag(tagID int) {
	for i, tag := range p.Tags {
		if tag.ID == tagID {
			p.Tags = append(p.Tags[:i], p.Tags[i+1:]...)
			p.DateModified = time.Now()
			break
		}
	}
}
