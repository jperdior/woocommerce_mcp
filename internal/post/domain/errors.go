package domain

import "fmt"

// PostError represents a domain error for posts
type PostError struct {
	Code    string
	Message string
	Type    string
}

func (e *PostError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *PostError {
	return &PostError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Type:    "ValidationError",
	}
}

// NewWordPressAPIError creates a new WordPress API error
func NewWordPressAPIError(statusCode int, message, code string) *PostError {
	return &PostError{
		Code:    fmt.Sprintf("WORDPRESS_API_ERROR_%d", statusCode),
		Message: fmt.Sprintf("WordPress API error (status %d): %s", statusCode, message),
		Type:    "WordPressAPIError",
	}
}

// NewConnectionError creates a new connection error
func NewConnectionError(url, message string) *PostError {
	return &PostError{
		Code:    "CONNECTION_ERROR",
		Message: fmt.Sprintf("connection error to %s: %s", url, message),
		Type:    "ConnectionError",
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(postID PostID) *PostError {
	return &PostError{
		Code:    "POST_NOT_FOUND",
		Message: fmt.Sprintf("post with ID %d not found", postID.Value()),
		Type:    "NotFoundError",
	}
}
