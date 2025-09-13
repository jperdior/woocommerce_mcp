package domain

import "fmt"

// ValidationError represents a domain validation error
type ValidationError struct {
	Message string
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

// Error returns the error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

// Is checks if the error is of the same type
func (e *ValidationError) Is(target error) bool {
	_, ok := target.(*ValidationError)
	return ok
}

// DomainError represents a generic domain error
type DomainError struct {
	Message string
	Code    string
}

// NewDomainError creates a new DomainError
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// Error returns the error message
func (e *DomainError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("domain error [%s]: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("domain error: %s", e.Message)
}

// Is checks if the error is of the same type
func (e *DomainError) Is(target error) bool {
	_, ok := target.(*DomainError)
	return ok
}

// GetCode returns the error code
func (e *DomainError) GetCode() string {
	return e.Code
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Resource string
	ID       string
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// Error returns the error message
func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// Is checks if the error is of the same type
func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

// ConflictError represents a conflict error
type ConflictError struct {
	Message string
}

// NewConflictError creates a new ConflictError
func NewConflictError(message string) *ConflictError {
	return &ConflictError{
		Message: message,
	}
}

// Error returns the error message
func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict error: %s", e.Message)
}

// Is checks if the error is of the same type
func (e *ConflictError) Is(target error) bool {
	_, ok := target.(*ConflictError)
	return ok
}
