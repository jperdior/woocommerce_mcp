package wordpress

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"woocommerce-mcp/internal/post/domain"
)

// Config represents WordPress API configuration
type Config struct {
	BaseURL string
	Timeout time.Duration
}

// NewConfig creates a new WordPress configuration
func NewConfig(baseURL string) *Config {
	return &Config{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		Timeout: 30 * time.Second,
	}
}

// Client represents a WordPress API client
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new WordPress client
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// SearchPosts searches for posts using the WordPress API
func (c *Client) SearchPosts(ctx context.Context, criteria *domain.SearchCriteria) ([]*domain.Post, error) {
	// Build the API endpoint URL
	endpoint := fmt.Sprintf("%s/wp-json/wp/v2/posts", c.config.BaseURL)

	// Parse base URL
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, domain.NewConnectionError(endpoint, fmt.Sprintf("invalid base URL: %v", err))
	}

	// Build query parameters
	query := u.Query()
	c.addSearchParams(query, criteria)

	u.RawQuery = query.Encode()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domain.NewConnectionError(u.String(), fmt.Sprintf("HTTP request failed: %v", err))
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleAPIError(resp.StatusCode, body)
	}

	// Parse JSON response
	var apiPosts []APIPost
	if err := json.Unmarshal(body, &apiPosts); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Convert API posts to domain posts
	posts := make([]*domain.Post, len(apiPosts))
	for i, apiPost := range apiPosts {
		domainPost, err := c.apiPostToDomain(&apiPost)
		if err != nil {
			return nil, fmt.Errorf("failed to convert post %d: %w", apiPost.ID, err)
		}
		posts[i] = domainPost
	}

	return posts, nil
}

// CountPosts counts posts matching the criteria
func (c *Client) CountPosts(ctx context.Context, criteria *domain.SearchCriteria) (int64, error) {
	// For WordPress API, we need to make a HEAD request or parse headers
	// Since WordPress doesn't provide a direct count endpoint, we'll use the X-WP-Total header
	endpoint := fmt.Sprintf("%s/wp-json/wp/v2/posts", c.config.BaseURL)

	u, err := url.Parse(endpoint)
	if err != nil {
		return 0, domain.NewConnectionError(endpoint, fmt.Sprintf("invalid base URL: %v", err))
	}

	// Build query parameters (same as search but we only need the count)
	query := u.Query()
	c.addSearchParams(query, criteria)

	// Set per_page to 1 to minimize data transfer when we only need the count
	query.Set("per_page", "1")

	u.RawQuery = query.Encode()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "HEAD", u.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Make HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, domain.NewConnectionError(u.String(), fmt.Sprintf("HTTP request failed: %v", err))
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return 0, c.handleAPIError(resp.StatusCode, nil)
	}

	// Get total count from header
	totalHeader := resp.Header.Get("X-WP-Total")
	if totalHeader == "" {
		// Fallback: return 0 if header is not available
		return 0, nil
	}

	total, err := strconv.ParseInt(totalHeader, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse total count: %w", err)
	}

	return total, nil
}

// addSearchParams adds search parameters to the query
func (c *Client) addSearchParams(query url.Values, criteria *domain.SearchCriteria) {
	if criteria.Search != "" {
		query.Set("search", criteria.Search)
	}
	if criteria.Status != "" {
		query.Set("status", string(criteria.Status))
	}
	if criteria.Author != 0 {
		query.Set("author", strconv.FormatInt(criteria.Author, 10))
	}
	if len(criteria.Categories) > 0 {
		categoryStrs := make([]string, len(criteria.Categories))
		for i, cat := range criteria.Categories {
			categoryStrs[i] = strconv.FormatInt(cat, 10)
		}
		query.Set("categories", strings.Join(categoryStrs, ","))
	}
	if len(criteria.Tags) > 0 {
		tagStrs := make([]string, len(criteria.Tags))
		for i, tag := range criteria.Tags {
			tagStrs[i] = strconv.FormatInt(tag, 10)
		}
		query.Set("tags", strings.Join(tagStrs, ","))
	}
	if criteria.Before != "" {
		query.Set("before", criteria.Before)
	}
	if criteria.After != "" {
		query.Set("after", criteria.After)
	}

	// Pagination
	if criteria.PerPage > 0 {
		query.Set("per_page", strconv.Itoa(criteria.PerPage))
	} else {
		query.Set("per_page", "10") // Default
	}
	if criteria.Page > 0 {
		query.Set("page", strconv.Itoa(criteria.Page))
	} else {
		query.Set("page", "1") // Default
	}

	// Sorting
	if criteria.OrderBy != "" {
		query.Set("orderby", criteria.OrderBy)
	} else {
		query.Set("orderby", "date") // Default
	}
	if criteria.Order != "" {
		query.Set("order", criteria.Order)
	} else {
		query.Set("order", "desc") // Default
	}
}

// handleAPIError handles API errors and converts them to domain errors
func (c *Client) handleAPIError(statusCode int, body []byte) error {
	message := string(body)
	if len(body) == 0 {
		message = http.StatusText(statusCode)
	}

	// Try to parse error response for more details
	var apiError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &apiError); err == nil {
			if apiError.Message != "" {
				message = apiError.Message
			}
		}
	}

	return domain.NewWordPressAPIError(statusCode, message, apiError.Code)
}

// apiPostToDomain converts an API post to a domain post
func (c *Client) apiPostToDomain(apiPost *APIPost) (*domain.Post, error) {
	// Create post ID
	postID, err := domain.NewPostID(apiPost.ID)
	if err != nil {
		return nil, err
	}

	// Create domain post
	post := domain.NewPost(postID, apiPost.Title.Rendered)

	// Set basic fields
	post.Content = apiPost.Content.Rendered
	post.Excerpt = apiPost.Excerpt.Rendered
	post.Slug = apiPost.Slug
	post.Permalink = apiPost.Link
	post.Type = apiPost.Type
	post.AuthorID = apiPost.Author
	post.FeaturedMediaID = apiPost.FeaturedMedia
	post.CommentStatus = apiPost.CommentStatus
	post.PingStatus = apiPost.PingStatus
	post.Sticky = apiPost.Sticky

	// Parse dates
	if apiPost.Date != "" {
		if dateCreated, err := time.Parse("2006-01-02T15:04:05", apiPost.Date); err == nil {
			post.DateCreated = dateCreated
		}
	}
	if apiPost.Modified != "" {
		if dateModified, err := time.Parse("2006-01-02T15:04:05", apiPost.Modified); err == nil {
			post.DateModified = dateModified
		}
	}
	if apiPost.DateGMT != "" {
		if dateGMT, err := time.Parse("2006-01-02T15:04:05", apiPost.DateGMT); err == nil {
			post.DateGMT = dateGMT
		}
	}
	if apiPost.ModifiedGMT != "" {
		if modifiedGMT, err := time.Parse("2006-01-02T15:04:05", apiPost.ModifiedGMT); err == nil {
			post.ModifiedGMT = modifiedGMT
		}
	}

	// Set post status
	if apiPost.Status != "" {
		status := domain.PostStatus(apiPost.Status)
		if status.IsValid() {
			post.Status = status
		}
	}

	// Set post format
	if apiPost.Format != "" {
		post.Format = domain.PostFormat(apiPost.Format)
	}

	// Convert meta data
	for key, value := range apiPost.MetaFields {
		metaData := domain.NewMetaData(0, key, value)
		post.MetaData = append(post.MetaData, *metaData)
	}

	return post, nil
}
