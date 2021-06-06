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
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
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

type Logger interface {
	LogError(err error)
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
	Logger    Logger
}

// Common sentinel errors
// TODO Fix naming convention on errors: avoid first uppercase letter
// TODO Use patter Err* for errors' names
var (
	OperationCanceledError = errors.New("The context of the operation was canceled")
	PostNotFoundError      = errors.New("The post requested was not found")
	MissingIdError         = errors.New("Please pass the Id from the post to be updated")
	UserNotFoundError      = errors.New("The creator user does not exist")
	UserCheckError         = errors.New("Error trying to check for the user's existence")
)

// Persists and return a PostDto with the data passed
func (r *PostRepository) CreatePost(ctx context.Context,
	post *CreatePostDto) (*Post, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "Context error")
	}
	defer cancel()
	exists, err := r.Checker.CheckExistence(ctx, post.Creator)
	if err != nil {
		return nil, r.logErrorAndWrap(UserCheckError, err.Error())
	}
	if !exists {
		return nil, r.logErrorAndWrap(UserNotFoundError,
			fmt.Sprintf("User %s not found", post.Creator))
	}
	post.Content = r.Sanitizer.SanitizeContent(post.Content)
	return r.Store.Create(ctx, post)
}

// Updates and return a PostDto with the data passed,
//   otherwise returns no-nil error
func (r *PostRepository) UpdatePost(ctx context.Context,
	updated *UpdatePostDto) (*Post, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "Context error")
	}
	defer cancel()
	if updated.Id == "" {
		return nil, r.logErrorAndWrap(MissingIdError, "UpdatePost")
	}
	updated.Content = r.Sanitizer.SanitizeContent(updated.Content)
	return r.Store.Update(ctx, updated)
}

// Retrieves a post by its identifier (id)
func (r *PostRepository) GetPost(ctx context.Context, id string) (*Post, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "Context error")
	}
	defer cancel()
	post, err := r.Store.ReadOne(ctx, id)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "GetPost error")
	}
	if post.Id == "" {
		return nil, r.logErrorAndWrap(PostNotFoundError, fmt.Sprintf("ID: %s.", id))
	}
	return post, nil
}

// Filters persisted posts by tag(s)
func (r *PostRepository) FilterByTag(ctx context.Context,
	filter *ByTagDto, page, pageSize int) ([]*Post, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "Context error")
	}
	defer cancel()
	generalFilter := &GeneralFilter{Page: page, PageSize: pageSize}
	generalFilter.Tag = filter.Tag
	return r.Store.Filter(ctx, generalFilter)
}

// Filters persisted posts by date range
func (r *PostRepository) FilterByDateRange(ctx context.Context,
	filter *ByDateRangeDto, page, pageSize int) ([]*Post, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "Context error")
	}
	defer cancel()
	generalFilter := &GeneralFilter{Page: page, PageSize: pageSize}
	generalFilter.From = filter.From
	generalFilter.To = filter.To
	return r.Store.Filter(ctx, generalFilter)
}

func (r *PostRepository) logErrorAndWrap(err error, msg string) error {
	r.Logger.LogError(err)
	return fmt.Errorf("%s: %w \n", msg, err)
}

// TODO Remove context recreation on Post's repository, unnecessary and error prone
func checkContextAndRecreate(
	ctx context.Context) (context.Context, context.CancelFunc, error) {
	select {
	case <-ctx.Done():
		return nil, nil, fmt.Errorf("%w: %v", OperationCanceledError, ctx.Err())
	default:
		newCtx, cancel := context.WithCancel(ctx)
		return newCtx, cancel, nil
	}
}
