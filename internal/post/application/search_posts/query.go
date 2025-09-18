package search_posts

import (
	"strconv"
	"strings"
	"woocommerce-mcp/internal/post/domain"
)

// Query represents a search posts query
type Query struct {
	BaseURL    string
	Search     string
	Status     domain.PostStatus
	Author     int64
	Categories []int64
	Tags       []int64
	Before     string
	After      string
	Page       int
	PerPage    int
	OrderBy    string
	Order      string
}

// NewQueryFromRequest creates a new Query from a SearchRequest
func NewQueryFromRequest(req *SearchRequest) (*Query, error) {
	query := &Query{
		BaseURL: req.BaseURL,
		Search:  req.Search,
		Before:  req.Before,
		After:   req.After,
		OrderBy: req.OrderBy,
		Order:   req.Order,
	}

	// Parse status
	if req.Status != "" {
		query.Status = domain.PostStatus(req.Status)
	}

	// Parse author
	if req.Author != "" {
		if author, err := strconv.ParseInt(req.Author, 10, 64); err == nil {
			query.Author = author
		}
	}

	// Parse categories
	if req.Categories != "" {
		categoryStrs := strings.Split(req.Categories, ",")
		for _, catStr := range categoryStrs {
			if catID, err := strconv.ParseInt(strings.TrimSpace(catStr), 10, 64); err == nil {
				query.Categories = append(query.Categories, catID)
			}
		}
	}

	// Parse tags
	if req.Tags != "" {
		tagStrs := strings.Split(req.Tags, ",")
		for _, tagStr := range tagStrs {
			if tagID, err := strconv.ParseInt(strings.TrimSpace(tagStr), 10, 64); err == nil {
				query.Tags = append(query.Tags, tagID)
			}
		}
	}

	// Parse pagination
	if req.Page != "" {
		if page, err := strconv.Atoi(req.Page); err == nil && page > 0 {
			query.Page = page
		}
	}
	if query.Page == 0 {
		query.Page = 1 // Default
	}

	if req.PerPage != "" {
		if perPage, err := strconv.Atoi(req.PerPage); err == nil && perPage > 0 {
			query.PerPage = perPage
		}
	}
	if query.PerPage == 0 {
		query.PerPage = 10 // Default
	}

	// Set defaults for sorting
	if query.OrderBy == "" {
		query.OrderBy = "date"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	return query, nil
}

// ToSearchCriteria converts the query to domain search criteria
func (q *Query) ToSearchCriteria() *domain.SearchCriteria {
	return &domain.SearchCriteria{
		Search:     q.Search,
		Status:     q.Status,
		Author:     q.Author,
		Categories: q.Categories,
		Tags:       q.Tags,
		Before:     q.Before,
		After:      q.After,
		Page:       q.Page,
		PerPage:    q.PerPage,
		OrderBy:    q.OrderBy,
		Order:      q.Order,
	}
}
