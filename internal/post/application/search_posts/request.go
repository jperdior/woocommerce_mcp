package search_posts

// SearchRequest represents a request to search for posts
type SearchRequest struct {
	BaseURL string `json:"base_url"`

	// Search parameters
	Search     string `json:"search,omitempty"`
	Status     string `json:"status,omitempty"`
	Author     string `json:"author,omitempty"`
	Categories string `json:"categories,omitempty"`
	Tags       string `json:"tags,omitempty"`
	Before     string `json:"before,omitempty"`
	After      string `json:"after,omitempty"`

	// Pagination
	Page    string `json:"page,omitempty"`
	PerPage string `json:"per_page,omitempty"`

	// Sorting
	OrderBy string `json:"orderby,omitempty"`
	Order   string `json:"order,omitempty"`
}
