package usecase

import (
	"context"
	"errors"
	"time"
)

// Tag struct
type Tag struct {
	Name string
}

// Post entity representation
type PostDto struct {
	Creator   string
	Content   string
	Tags      []Tag
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Dto for handling creation of Posts
type CreatePostDto struct {
	Creator string
	Content string
	Tags    []string
}

// Dto for handling update of Posts
type UpdatePostDto struct {
	Content string
	Tags    []string
}

// Dto for handling filtering by tag
type ByTagDto struct {
	Tag string
}

// Dto for handling filtering by tag
type ByDateRangeDto struct {
	From time.Time
	To   time.Time
}

// Composed filter dto
type GeneralFilter struct {
	ByTagDto
	ByDateRangeDto
}

// Contract for the needs of a post's repo in terms of persistence
//    The Update method should return the updated version of the post
type PostStore interface {
	Create(context.Context, *CreatePostDto) (*PostDto, error)
	Update(context.Context, *UpdatePostDto) (*PostDto, error)
	Filter(context.Context, *GeneralFilter) ([]*PostDto, error)
}

// Basic contract intended to enforce sanitizing of content to avoid
//   malicious code injection
type ContentSanitizer interface {
	SanitizeContent(context.Context, string) (string, error)
}

type Logger interface {
	LogError(err error)
}

type PostRepository struct {
	Store     PostStore
	Sanitizer ContentSanitizer
	Logger    Logger
}

// Common sentinel errors
var OperationCanceledError = errors.New("The context of the operation was canceled")

func (r *PostRepository) CreatePost(ctx context.Context, post *CreatePostDto) (*PostDto, error) {
	// TODO Implement
	return nil, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, updated *UpdatePostDto) (*PostDto, error) {
	// TODO Implement
	return nil, nil
}

func (r *PostRepository) FilterByTag(ctx context.Context, filter *ByTagDto) ([]*PostDto, error) {
	// TODO Implement
	return nil, nil
}

func (r *PostRepository) FilterByDateRange(ctx context.Context, filter *ByDateRangeDto) ([]*PostDto, error) {
	// TODO Implement
	return nil, nil
}
