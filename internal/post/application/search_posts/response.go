package search_posts

import (
	"encoding/json"
	"woocommerce-mcp/internal/post/domain"
)

// SearchResponse represents a response from searching posts
type SearchResponse struct {
	Posts       []PostDTO `json:"posts"`
	TotalCount  int64     `json:"total_count"`
	CurrentPage int       `json:"current_page"`
	PerPage     int       `json:"per_page"`
	TotalPages  int       `json:"total_pages"`
	HasNext     bool      `json:"has_next"`
	HasPrev     bool      `json:"has_prev"`
}

// PostDTO represents a post data transfer object
type PostDTO struct {
	ID              int64         `json:"id"`
	Title           string        `json:"title"`
	Content         string        `json:"content"`
	Excerpt         string        `json:"excerpt"`
	Slug            string        `json:"slug"`
	Status          string        `json:"status"`
	Format          string        `json:"format"`
	Type            string        `json:"type"`
	Permalink       string        `json:"permalink"`
	FeaturedMediaID int64         `json:"featured_media_id"`
	AuthorID        int64         `json:"author_id"`
	DateCreated     string        `json:"date_created"`
	DateModified    string        `json:"date_modified"`
	CommentStatus   string        `json:"comment_status"`
	PingStatus      string        `json:"ping_status"`
	Sticky          bool          `json:"sticky"`
	Tags            []TagDTO      `json:"tags"`
	Categories      []CategoryDTO `json:"categories"`
	MetaData        []MetaDataDTO `json:"meta_data"`
}

// TagDTO represents a tag data transfer object
type TagDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Link string `json:"link"`
}

// CategoryDTO represents a category data transfer object
type CategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Link string `json:"link"`
}

// MetaDataDTO represents metadata data transfer object
type MetaDataDTO struct {
	ID    int64       `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// ToJSON converts the response to JSON string
func (r *SearchResponse) ToJSON() (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromDomainPosts converts domain posts to response DTOs
func FromDomainPosts(posts []*domain.Post, totalCount int64, currentPage, perPage int) *SearchResponse {
	postDTOs := make([]PostDTO, len(posts))
	for i, post := range posts {
		postDTOs[i] = PostDTO{
			ID:              post.ID.Value(),
			Title:           post.Title,
			Content:         post.Content,
			Excerpt:         post.Excerpt,
			Slug:            post.Slug,
			Status:          string(post.Status),
			Format:          string(post.Format),
			Type:            post.Type,
			Permalink:       post.Permalink,
			FeaturedMediaID: post.FeaturedMediaID,
			AuthorID:        post.AuthorID,
			DateCreated:     post.DateCreated.Format("2006-01-02T15:04:05"),
			DateModified:    post.DateModified.Format("2006-01-02T15:04:05"),
			CommentStatus:   post.CommentStatus,
			PingStatus:      post.PingStatus,
			Sticky:          post.Sticky,
		}

		// Convert tags
		for _, tag := range post.Tags {
			postDTOs[i].Tags = append(postDTOs[i].Tags, TagDTO{
				ID:   tag.ID,
				Name: tag.Name,
				Slug: tag.Slug,
				Link: tag.Link,
			})
		}

		// Convert categories
		for _, category := range post.Categories {
			postDTOs[i].Categories = append(postDTOs[i].Categories, CategoryDTO{
				ID:   category.ID,
				Name: category.Name,
				Slug: category.Slug,
				Link: category.Link,
			})
		}

		// Convert metadata
		for _, meta := range post.MetaData {
			postDTOs[i].MetaData = append(postDTOs[i].MetaData, MetaDataDTO{
				ID:    meta.ID,
				Key:   meta.Key,
				Value: meta.Value,
			})
		}
	}

	totalPages := int(totalCount) / perPage
	if int(totalCount)%perPage != 0 {
		totalPages++
	}

	return &SearchResponse{
		Posts:       postDTOs,
		TotalCount:  totalCount,
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     currentPage < totalPages,
		HasPrev:     currentPage > 1,
	}
}
