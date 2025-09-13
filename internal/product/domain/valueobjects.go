package domain

import (
	"fmt"
	"strconv"
	"strings"
	"woocommerce-mcp/kit/domain"
)

// ProductID represents a unique identifier for a product
type ProductID struct {
	value int
}

// NewProductID creates a new ProductID
func NewProductID(value int) (*ProductID, error) {
	if value <= 0 {
		return nil, domain.NewValidationError("product ID must be positive")
	}
	return &ProductID{value: value}, nil
}

// NewProductIDFromString creates a ProductID from string
func NewProductIDFromString(value string) (*ProductID, error) {
	id, err := strconv.Atoi(value)
	if err != nil {
		return nil, domain.NewValidationError("invalid product ID format")
	}
	return NewProductID(id)
}

// Value returns the underlying value
func (p *ProductID) Value() int {
	return p.value
}

// String returns string representation
func (p *ProductID) String() string {
	return strconv.Itoa(p.value)
}

// Equals checks if two ProductIDs are equal
func (p *ProductID) Equals(other *ProductID) bool {
	if other == nil {
		return false
	}
	return p.value == other.value
}

// ProductType represents the type of a product
type ProductType string

const (
	ProductTypeSimple   ProductType = "simple"
	ProductTypeGrouped  ProductType = "grouped"
	ProductTypeExternal ProductType = "external"
	ProductTypeVariable ProductType = "variable"
)

// IsValid checks if the product type is valid
func (pt ProductType) IsValid() bool {
	switch pt {
	case ProductTypeSimple, ProductTypeGrouped, ProductTypeExternal, ProductTypeVariable:
		return true
	default:
		return false
	}
}

// String returns string representation
func (pt ProductType) String() string {
	return string(pt)
}

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusDraft   ProductStatus = "draft"
	ProductStatusPending ProductStatus = "pending"
	ProductStatusPrivate ProductStatus = "private"
	ProductStatusPublish ProductStatus = "publish"
)

// IsValid checks if the product status is valid
func (ps ProductStatus) IsValid() bool {
	switch ps {
	case ProductStatusDraft, ProductStatusPending, ProductStatusPrivate, ProductStatusPublish:
		return true
	default:
		return false
	}
}

// String returns string representation
func (ps ProductStatus) String() string {
	return string(ps)
}

// StockStatus represents the stock status of a product
type StockStatus string

const (
	StockStatusInStock     StockStatus = "instock"
	StockStatusOutOfStock  StockStatus = "outofstock"
	StockStatusOnBackorder StockStatus = "onbackorder"
)

// IsValid checks if the stock status is valid
func (ss StockStatus) IsValid() bool {
	switch ss {
	case StockStatusInStock, StockStatusOutOfStock, StockStatusOnBackorder:
		return true
	default:
		return false
	}
}

// String returns string representation
func (ss StockStatus) String() string {
	return string(ss)
}

// Money represents a monetary value
type Money struct {
	amount   float64
	currency string
}

// NewMoney creates a new Money value object
func NewMoney(amount float64, currency string) (*Money, error) {
	if amount < 0 {
		return nil, domain.NewValidationError("amount cannot be negative")
	}
	if currency == "" {
		currency = "USD" // Default currency
	}
	return &Money{
		amount:   amount,
		currency: strings.ToUpper(currency),
	}, nil
}

// NewMoneyFromString creates Money from string representation
func NewMoneyFromString(value, currency string) (*Money, error) {
	if value == "" {
		return NewMoney(0, currency)
	}

	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, domain.NewValidationError("invalid money format")
	}

	return NewMoney(amount, currency)
}

// Amount returns the monetary amount
func (m *Money) Amount() float64 {
	return m.amount
}

// Currency returns the currency code
func (m *Money) Currency() string {
	return m.currency
}

// String returns string representation
func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}

// Equals checks if two Money values are equal
func (m *Money) Equals(other *Money) bool {
	if other == nil {
		return false
	}
	return m.amount == other.amount && m.currency == other.currency
}

// IsZero checks if the money value is zero
func (m *Money) IsZero() bool {
	return m.amount == 0
}

// Dimensions represents product dimensions
type Dimensions struct {
	Length string `json:"length"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

// NewDimensions creates new dimensions
func NewDimensions(length, width, height string) *Dimensions {
	return &Dimensions{
		Length: length,
		Width:  width,
		Height: height,
	}
}

// IsEmpty checks if all dimensions are empty
func (d *Dimensions) IsEmpty() bool {
	return d.Length == "" && d.Width == "" && d.Height == ""
}

// Category represents a product category
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// NewCategory creates a new category
func NewCategory(id int, name, slug string) *Category {
	return &Category{
		ID:   id,
		Name: name,
		Slug: slug,
	}
}

// Tag represents a product tag
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// NewTag creates a new tag
func NewTag(id int, name, slug string) *Tag {
	return &Tag{
		ID:   id,
		Name: name,
		Slug: slug,
	}
}

// Image represents a product image
type Image struct {
	ID           int    `json:"id"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
	Src          string `json:"src"`
	Name         string `json:"name"`
	Alt          string `json:"alt"`
	Position     int    `json:"position"`
}

// NewImage creates a new image
func NewImage(id int, src, name, alt string, position int) *Image {
	return &Image{
		ID:       id,
		Src:      src,
		Name:     name,
		Alt:      alt,
		Position: position,
	}
}

// Attribute represents a product attribute
type Attribute struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Visible   bool     `json:"visible"`
	Variation bool     `json:"variation"`
	Options   []string `json:"options"`
}

// NewAttribute creates a new attribute
func NewAttribute(id int, name string, options []string) *Attribute {
	return &Attribute{
		ID:      id,
		Name:    name,
		Options: options,
		Visible: true,
	}
}

// DefaultAttribute represents a default product attribute
type DefaultAttribute struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Option string `json:"option"`
}

// NewDefaultAttribute creates a new default attribute
func NewDefaultAttribute(id int, name, option string) *DefaultAttribute {
	return &DefaultAttribute{
		ID:     id,
		Name:   name,
		Option: option,
	}
}

// MetaData represents product metadata
type MetaData struct {
	ID    int         `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// NewMetaData creates new metadata
func NewMetaData(id int, key string, value interface{}) *MetaData {
	return &MetaData{
		ID:    id,
		Key:   key,
		Value: value,
	}
}
