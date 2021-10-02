// This package defines general use cases related to posts
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Post entity representation
type Post struct {
	Id        string    `json:"id"`
	Creator   string    `json:"creator"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Tags      []string  `json:"tags"`
}

// Dto for handling creation of Posts
type CreatePostDto struct {
	Title   string
	Creator string
	Content string
	Tags    []string
}

// Dto for handling update of Posts
type UpdatePostDto struct {
	Id      string
	Title   string
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
	Page     int
	PageSize int
}

// Contract for the needs of a post's repo in terms of persistence
//    The Update method should return the updated version of the post
type PostStore interface {
	Create(context.Context, *CreatePostDto) (*Post, error)
	Update(context.Context, *UpdatePostDto) (*Post, error)
	Filter(context.Context, *GeneralFilter) ([]*Post, error)
	ReadOne(context.Context, string) (*Post, error)
}

// Basic contract intended to enforce sanitizing of content to avoid
//   malicious code injection
type ContentSanitizer interface {
	SanitizeContent(string) string
}

// Defines the needed method for checking that a given creator of
//    a post indeed exists
type CreatorChecker interface {
	CheckExistence(context.Context, string) (bool, error)
}

// Repository defines the basic contract for Post's usecases
type Repository interface {
	CreatePost(context.Context, *CreatePostDto) (*Post, error)
	UpdatePost(context.Context, *UpdatePostDto) (*Post, error)
	GetPost(context.Context, string) (*Post, error)
	FilterByTag(ctx context.Context, filter *ByTagDto, page, pageSize int) ([]*Post, error)
	FilterByDateRange(ctx context.Context, filter *ByDateRangeDto, page, pageSize int) ([]*Post, error)
}

var _ Repository = &PostRepository{}

type PostRepository struct {
	Store     PostStore
	Checker   CreatorChecker
	Sanitizer ContentSanitizer
}

// Common sentinel errors
var (
	// ErrOperationCanceled returned when the context is canceled
	ErrOperationCanceled = errors.New("context of the operation was canceled")
	// ErrPostNotFound is self-described
	ErrPostNotFound = errors.New("post requested was not found")
	// ErrMissingID returned when no ID was passed
	ErrMissingID = errors.New("missing ID from the post to be updated")
	// ErrUserNotFound is self-described
	ErrUserNotFound = errors.New("creator user does not exist")
	// ErrUserCheck returned when there's an error in the upstream auth service
	ErrUserCheck = errors.New("check for user's existence")
	// ErrEmptyTags returned when tags passed is empty, on creation
	ErrEmptyTags = errors.New("tags can't be empty")
)

// Persists and return a PostDto with the data passed
func (r *PostRepository) CreatePost(
	ctx context.Context,
	post *CreatePostDto,
) (*Post, error) {
	exists, err := r.Checker.CheckExistence(ctx, post.Creator)
	if err != nil {
		return nil, logErrorAndWrap(ErrUserCheck, err.Error())
	}
	if !exists {
		return nil, logErrorAndWrap(ErrUserNotFound,
			fmt.Sprintf("User %s not found", post.Creator))
	}
	if len(post.Tags) == 0 {
		return nil, fmt.Errorf("create post: %w", ErrEmptyTags)
	}
	post.Content = r.Sanitizer.SanitizeContent(post.Content)
	return r.Store.Create(ctx, post)
}

// Updates and return a PostDto with the data passed,
//   otherwise returns no-nil error
func (r *PostRepository) UpdatePost(
	ctx context.Context,
	updated *UpdatePostDto,
) (*Post, error) {
	if updated.Id == "" {
		return nil, logErrorAndWrap(ErrMissingID, "UpdatePost")
	}
	updated.Content = r.Sanitizer.SanitizeContent(updated.Content)
	return r.Store.Update(ctx, updated)
}

// Retrieves a post by its identifier (id)
func (r *PostRepository) GetPost(ctx context.Context, id string) (*Post, error) {
	post, err := r.Store.ReadOne(ctx, id)
	if err != nil {
		return nil, logErrorAndWrap(err, "GetPost error")
	}
	if post.Id == "" {
		return nil, logErrorAndWrap(ErrPostNotFound, fmt.Sprintf("ID: %s.", id))
	}
	return post, nil
}

// Filters persisted posts by tag(s)
func (r *PostRepository) FilterByTag(
	ctx context.Context,
	filter *ByTagDto,
	page, pageSize int,
) ([]*Post, error) {
	generalFilter := &GeneralFilter{Page: page, PageSize: pageSize}
	generalFilter.Tag = filter.Tag
	return r.Store.Filter(ctx, generalFilter)
}

// Filters persisted posts by date range
func (r *PostRepository) FilterByDateRange(
	ctx context.Context,
	filter *ByDateRangeDto,
	page, pageSize int,
) ([]*Post, error) {
	generalFilter := &GeneralFilter{Page: page, PageSize: pageSize}
	generalFilter.From = filter.From
	generalFilter.To = filter.To
	return r.Store.Filter(ctx, generalFilter)
}

// TODO Remove logErrorAndWrap function as it's unnecessary, posts
func logErrorAndWrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
