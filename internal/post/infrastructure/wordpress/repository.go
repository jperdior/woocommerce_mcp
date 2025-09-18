package wordpress

import (
	"context"
	"woocommerce-mcp/internal/post/domain"
)

// Repository implements the domain PostRepository interface
type Repository struct {
	client *Client
}

// NewRepository creates a new WordPress repository
func NewRepository(client *Client) *Repository {
	return &Repository{
		client: client,
	}
}

// SearchPosts searches for posts using the WordPress API
func (r *Repository) SearchPosts(ctx context.Context, criteria *domain.SearchCriteria) ([]*domain.Post, error) {
	return r.client.SearchPosts(ctx, criteria)
}

// CountPosts returns the total count of posts matching the criteria
func (r *Repository) CountPosts(ctx context.Context, criteria *domain.SearchCriteria) (int64, error) {
	return r.client.CountPosts(ctx, criteria)
}

// GetPostByID retrieves a post by its ID (for future implementation)
func (r *Repository) GetPostByID(ctx context.Context, id domain.PostID) (*domain.Post, error) {
	// For now, we'll use search with a specific criteria
	// In a real implementation, you might want a dedicated endpoint
	criteria := &domain.SearchCriteria{
		PerPage: 1,
		Page:    1,
	}

	posts, err := r.client.SearchPosts(ctx, criteria)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		if post.ID == id {
			return post, nil
		}
	}

	return nil, domain.NewNotFoundError(id)
}
