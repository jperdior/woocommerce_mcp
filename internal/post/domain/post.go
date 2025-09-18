package domain

import (
	"time"
)

// PostID represents a unique identifier for a post
type PostID int64

// NewPostID creates a new PostID
func NewPostID(id int64) (PostID, error) {
	if id <= 0 {
		return 0, NewValidationError("post ID must be positive")
	}
	return PostID(id), nil
}

// Value returns the underlying int64 value
func (id PostID) Value() int64 {
	return int64(id)
}

// PostStatus represents the status of a post
type PostStatus string

const (
	PostStatusPublish PostStatus = "publish"
	PostStatusDraft   PostStatus = "draft"
	PostStatusPrivate PostStatus = "private"
	PostStatusPending PostStatus = "pending"
	PostStatusTrash   PostStatus = "trash"
)

// IsValid checks if the post status is valid
func (s PostStatus) IsValid() bool {
	switch s {
	case PostStatusPublish, PostStatusDraft, PostStatusPrivate, PostStatusPending, PostStatusTrash:
		return true
	default:
		return false
	}
}

// PostFormat represents the format of a post
type PostFormat string

const (
	PostFormatStandard PostFormat = "standard"
	PostFormatAside    PostFormat = "aside"
	PostFormatGallery  PostFormat = "gallery"
	PostFormatLink     PostFormat = "link"
	PostFormatImage    PostFormat = "image"
	PostFormatQuote    PostFormat = "quote"
	PostFormatStatus   PostFormat = "status"
	PostFormatVideo    PostFormat = "video"
	PostFormatAudio    PostFormat = "audio"
	PostFormatChat     PostFormat = "chat"
)

// Post represents a WordPress post
type Post struct {
	ID              PostID
	Title           string
	Content         string
	Excerpt         string
	Slug            string
	Status          PostStatus
	Format          PostFormat
	Type            string
	Permalink       string
	FeaturedMediaID int64
	AuthorID        int64
	DateCreated     time.Time
	DateModified    time.Time
	DateGMT         time.Time
	ModifiedGMT     time.Time
	CommentStatus   string
	PingStatus      string
	Sticky          bool
	Tags            []Tag
	Categories      []Category
	MetaData        []MetaData
}

// NewPost creates a new Post
func NewPost(id PostID, title string) *Post {
	return &Post{
		ID:     id,
		Title:  title,
		Status: PostStatusPublish,
		Format: PostFormatStandard,
		Type:   "post",
	}
}

// Tag represents a post tag
type Tag struct {
	ID   int64
	Name string
	Slug string
	Link string
}

// NewTag creates a new Tag
func NewTag(id int64, name, slug string) *Tag {
	return &Tag{
		ID:   id,
		Name: name,
		Slug: slug,
	}
}

// Category represents a post category
type Category struct {
	ID   int64
	Name string
	Slug string
	Link string
}

// NewCategory creates a new Category
func NewCategory(id int64, name, slug string) *Category {
	return &Category{
		ID:   id,
		Name: name,
		Slug: slug,
	}
}

// MetaData represents post metadata
type MetaData struct {
	ID    int64
	Key   string
	Value interface{}
}

// NewMetaData creates new MetaData
func NewMetaData(id int64, key string, value interface{}) *MetaData {
	return &MetaData{
		ID:    id,
		Key:   key,
		Value: value,
	}
}
