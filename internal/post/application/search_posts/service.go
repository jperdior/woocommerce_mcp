package search_posts

import (
	"context"
	"fmt"
	"woocommerce-mcp/internal/post/domain"
	"woocommerce-mcp/internal/post/infrastructure/wordpress"
)

// PostSearcher handles post search operations
type PostSearcher struct {
	repository domain.PostRepository
}

// NewPostSearcher creates a new PostSearcher
func NewPostSearcher(repository domain.PostRepository) *PostSearcher {
	return &PostSearcher{
		repository: repository,
	}
}

// SearchPosts searches for posts based on the provided request
func (s *PostSearcher) SearchPosts(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	// Validate request
	if req.BaseURL == "" {
		return nil, fmt.Errorf("base_url is required")
	}

	// Convert request to query
	query, err := NewQueryFromRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}

	// Create WordPress client and repository for this request
	config := wordpress.NewConfig(query.BaseURL)
	client := wordpress.NewClient(config)
	repository := wordpress.NewRepository(client)

	// Search for posts
	posts, err := repository.SearchPosts(ctx, query.ToSearchCriteria())
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}

	// Get total count
	totalCount, err := repository.CountPosts(ctx, query.ToSearchCriteria())
	if err != nil {
		// If count fails, we'll continue with 0 - it's not critical
		totalCount = 0
	}

	// Convert to response
	response := FromDomainPosts(posts, totalCount, query.Page, query.PerPage)

	return response, nil
}

// Execute is an alias for SearchPosts to maintain consistency with the product searcher
func (s *PostSearcher) Execute(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	return s.SearchPosts(ctx, req)
}
