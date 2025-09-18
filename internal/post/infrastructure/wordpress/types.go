package wordpress

// APIPost represents a post from the WordPress REST API
type APIPost struct {
	ID            int64                  `json:"id"`
	Date          string                 `json:"date"`
	DateGMT       string                 `json:"date_gmt"`
	GUID          GUID                   `json:"guid"`
	Modified      string                 `json:"modified"`
	ModifiedGMT   string                 `json:"modified_gmt"`
	Slug          string                 `json:"slug"`
	Status        string                 `json:"status"`
	Type          string                 `json:"type"`
	Link          string                 `json:"link"`
	Title         Title                  `json:"title"`
	Content       Content                `json:"content"`
	Excerpt       Excerpt                `json:"excerpt"`
	Author        int64                  `json:"author"`
	FeaturedMedia int64                  `json:"featured_media"`
	CommentStatus string                 `json:"comment_status"`
	PingStatus    string                 `json:"ping_status"`
	Sticky        bool                   `json:"sticky"`
	Template      string                 `json:"template"`
	Format        string                 `json:"format"`
	MetaFields    map[string]interface{} `json:"meta"`
	Categories    []int64                `json:"categories"`
	Tags          []int64                `json:"tags"`
}

// GUID represents the GUID field from WordPress API
type GUID struct {
	Rendered string `json:"rendered"`
}

// Title represents the title field from WordPress API
type Title struct {
	Rendered string `json:"rendered"`
}

// Content represents the content field from WordPress API
type Content struct {
	Rendered  string `json:"rendered"`
	Protected bool   `json:"protected"`
}

// Excerpt represents the excerpt field from WordPress API
type Excerpt struct {
	Rendered  string `json:"rendered"`
	Protected bool   `json:"protected"`
}

// MetaField represents a meta field from WordPress API
type MetaField struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// APICategory represents a category from WordPress API
type APICategory struct {
	ID          int64  `json:"id"`
	Count       int    `json:"count"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Taxonomy    string `json:"taxonomy"`
	Parent      int64  `json:"parent"`
}

// APITag represents a tag from WordPress API
type APITag struct {
	ID          int64  `json:"id"`
	Count       int    `json:"count"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Taxonomy    string `json:"taxonomy"`
}
