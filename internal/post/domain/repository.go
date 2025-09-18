package domain

import "context"

// PostRepository defines the interface for post data access
type PostRepository interface {
	// SearchPosts searches for posts based on criteria
	SearchPosts(ctx context.Context, criteria *SearchCriteria) ([]*Post, error)

	// CountPosts returns the total count of posts matching the criteria
	CountPosts(ctx context.Context, criteria *SearchCriteria) (int64, error)

	// GetPostByID retrieves a post by its ID
	GetPostByID(ctx context.Context, id PostID) (*Post, error)
}

// SearchCriteria represents search parameters for posts
type SearchCriteria struct {
	// Basic search
	Search string

	// Filtering
	Status     PostStatus
	Author     int64
	Categories []int64
	Tags       []int64

	// Date filtering
	Before string // ISO 8601 format
	After  string // ISO 8601 format

	// Pagination
	Page    int
	PerPage int

	// Sorting
	OrderBy string // date, relevance, id, include, title, slug
	Order   string // asc, desc
}
