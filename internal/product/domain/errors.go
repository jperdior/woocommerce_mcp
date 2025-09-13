package domain

import (
	"fmt"
	"woocommerce-mcp/kit/domain"
)

// ProductNotFoundError represents an error when a product is not found
type ProductNotFoundError struct {
	ProductID *ProductID
}

// NewProductNotFoundError creates a new ProductNotFoundError
func NewProductNotFoundError(productID *ProductID) *ProductNotFoundError {
	return &ProductNotFoundError{
		ProductID: productID,
	}
}

// Error returns the error message
func (e *ProductNotFoundError) Error() string {
	if e.ProductID != nil {
		return fmt.Sprintf("product with ID %s not found", e.ProductID.String())
	}
	return "product not found"
}

// Is checks if the error is of the same type
func (e *ProductNotFoundError) Is(target error) bool {
	_, ok := target.(*ProductNotFoundError)
	return ok
}

// ProductValidationError represents a product validation error
type ProductValidationError struct {
	Field   string
	Message string
}

// NewProductValidationError creates a new ProductValidationError
func NewProductValidationError(field, message string) *ProductValidationError {
	return &ProductValidationError{
		Field:   field,
		Message: message,
	}
}

// Error returns the error message
func (e *ProductValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// Is checks if the error is of the same type
func (e *ProductValidationError) Is(target error) bool {
	_, ok := target.(*ProductValidationError)
	return ok
}

// WooCommerceAPIError represents an error from the WooCommerce API
type WooCommerceAPIError struct {
	StatusCode int
	Message    string
	Code       string
}

// NewWooCommerceAPIError creates a new WooCommerceAPIError
func NewWooCommerceAPIError(statusCode int, message, code string) *WooCommerceAPIError {
	return &WooCommerceAPIError{
		StatusCode: statusCode,
		Message:    message,
		Code:       code,
	}
}

// Error returns the error message
func (e *WooCommerceAPIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("WooCommerce API error (status %d, code %s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("WooCommerce API error (status %d): %s", e.StatusCode, e.Message)
}

// Is checks if the error is of the same type
func (e *WooCommerceAPIError) Is(target error) bool {
	_, ok := target.(*WooCommerceAPIError)
	return ok
}

// IsNotFound checks if the error represents a not found error
func (e *WooCommerceAPIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsUnauthorized checks if the error represents an unauthorized error
func (e *WooCommerceAPIError) IsUnauthorized() bool {
	return e.StatusCode == 401 || e.StatusCode == 403
}

// IsBadRequest checks if the error represents a bad request error
func (e *WooCommerceAPIError) IsBadRequest() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError checks if the error represents a server error
func (e *WooCommerceAPIError) IsServerError() bool {
	return e.StatusCode >= 500
}

// SearchCriteriaError represents an error in search criteria
type SearchCriteriaError struct {
	Field   string
	Message string
}

// NewSearchCriteriaError creates a new SearchCriteriaError
func NewSearchCriteriaError(field, message string) *SearchCriteriaError {
	return &SearchCriteriaError{
		Field:   field,
		Message: message,
	}
}

// Error returns the error message
func (e *SearchCriteriaError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("search criteria error for field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("search criteria error: %s", e.Message)
}

// Is checks if the error is of the same type
func (e *SearchCriteriaError) Is(target error) bool {
	_, ok := target.(*SearchCriteriaError)
	return ok
}

// ConnectionError represents a connection error to WooCommerce
type ConnectionError struct {
	URL     string
	Message string
}

// NewConnectionError creates a new ConnectionError
func NewConnectionError(url, message string) *ConnectionError {
	return &ConnectionError{
		URL:     url,
		Message: message,
	}
}

// Error returns the error message
func (e *ConnectionError) Error() string {
	return fmt.Sprintf("connection error to %s: %s", e.URL, e.Message)
}

// Is checks if the error is of the same type
func (e *ConnectionError) Is(target error) bool {
	_, ok := target.(*ConnectionError)
	return ok
}

// AuthenticationError represents an authentication error
type AuthenticationError struct {
	Message string
}

// NewAuthenticationError creates a new AuthenticationError
func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{
		Message: message,
	}
}

// Error returns the error message
func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error: %s", e.Message)
}

// Is checks if the error is of the same type
func (e *AuthenticationError) Is(target error) bool {
	_, ok := target.(*AuthenticationError)
	return ok
}

// Helper functions to create common domain errors

// NewInvalidProductIDError creates a validation error for invalid product ID
func NewInvalidProductIDError() error {
	return domain.NewValidationError("invalid product ID")
}

// NewEmptyProductNameError creates a validation error for empty product name
func NewEmptyProductNameError() error {
	return domain.NewValidationError("product name cannot be empty")
}

// NewInvalidPriceError creates a validation error for invalid price
func NewInvalidPriceError() error {
	return domain.NewValidationError("invalid price value")
}

// NewInvalidProductTypeError creates a validation error for invalid product type
func NewInvalidProductTypeError(productType string) error {
	return domain.NewValidationError(fmt.Sprintf("invalid product type: %s", productType))
}

// NewInvalidProductStatusError creates a validation error for invalid product status
func NewInvalidProductStatusError(status string) error {
	return domain.NewValidationError(fmt.Sprintf("invalid product status: %s", status))
}

// NewInvalidStockStatusError creates a validation error for invalid stock status
func NewInvalidStockStatusError(stockStatus string) error {
	return domain.NewValidationError(fmt.Sprintf("invalid stock status: %s", stockStatus))
}
