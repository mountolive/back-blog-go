// This package defines general use cases related to posts
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Tag struct
type Tag struct {
	Name string
}

// Post entity representation
type PostDto struct {
	Id        string
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
	Id      string
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
	Create(context.Context, *CreatePostDto) (*PostDto, error)
	Update(context.Context, *UpdatePostDto) (*PostDto, error)
	Filter(context.Context, *GeneralFilter) ([]*PostDto, error)
	ReadOne(context.Context, string) (*PostDto, error)
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

type PostRepository struct {
	Store     PostStore
	Checker   CreatorChecker
	Sanitizer ContentSanitizer
	Logger    Logger
}

// Common sentinel errors
var (
	OperationCanceledError = errors.New("The context of the operation was canceled")
	PostNotFoundError      = errors.New("The post requested was not found")
	UserNotFoundError      = errors.New("The creator user does not exist")
	UserCheckError         = errors.New("There was an error trying to check for the user's existence")
)

func (r *PostRepository) CreatePost(ctx context.Context, post *CreatePostDto) (*PostDto, error) {
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
		return nil, r.logErrorAndWrap(UserNotFoundError, fmt.Sprintf("User %s not found", post.Creator))
	}
	post.Content = r.Sanitizer.SanitizeContent(post.Content)
	return r.Store.Create(ctx, post)
}

func (r *PostRepository) UpdatePost(ctx context.Context, updated *UpdatePostDto) (*PostDto, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "Context error")
	}
	defer cancel()
	updated.Content = r.Sanitizer.SanitizeContent(updated.Content)
	return r.Store.Update(ctx, updated)
}

func (r *PostRepository) GetPost(ctx context.Context, id string) (*PostDto, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "Context error")
	}
	defer cancel()
	return r.Store.ReadOne(ctx, id)
}

func (r *PostRepository) FilterByTag(ctx context.Context,
	filter *ByTagDto, page, pageSize int) ([]*PostDto, error) {
	ctx, cancel, err := checkContextAndRecreate(ctx)
	if err != nil {
		return nil, r.logErrorAndWrap(err, "Context error")
	}
	defer cancel()
	generalFilter := &GeneralFilter{Page: page, PageSize: pageSize}
	generalFilter.Tag = filter.Tag
	return r.Store.Filter(ctx, generalFilter)
}

func (r *PostRepository) FilterByDateRange(ctx context.Context,
	filter *ByDateRangeDto, page, pageSize int) ([]*PostDto, error) {
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

func checkContextAndRecreate(ctx context.Context) (context.Context, context.CancelFunc, error) {
	select {
	case <-ctx.Done():
		return nil, nil, fmt.Errorf("%w: %v", OperationCanceledError, ctx.Err())
	default:
		newCtx, cancel := context.WithCancel(ctx)
		return newCtx, cancel, nil
	}
}
