package command

import (
	"context"
	"errors"

	"github.com/mountolive/back-blog-go/post/eventbus"
	"github.com/mountolive/back-blog-go/post/usecase"
)

var (
	// ErrIDMissing is self-described
	ErrIDMissing = errors.New("id missing")
	// ErrContentMissing is self-described
	ErrContentMissing = errors.New("content missing")
	// ErrCreatorMissing is self-described
	ErrCreatorMissing = errors.New("creator missing")
	// ErrTitleMissing is self-described
	ErrTitleMissing = errors.New("title missing")
)

// NewCreatePost is a constructor
func NewCreatePost(
	store usecase.PostStore,
	checker usecase.CreatorChecker,
) CreatePost {
	return CreatePost{
		store:          store,
		creatorChecker: checker,
	}
}

// CreatePost is a command handler
type CreatePost struct {
	store          usecase.PostStore
	creatorChecker usecase.CreatorChecker
}

// Handle is CommandHandler's implementation
func (CreatePost) Handle(context.Context, eventbus.Params) error {
	// TODO Implement
	return nil
}

// NewUpdatePost is a constructor
func NewUpdatePost(store usecase.PostStore) UpdatePost {
	return UpdatePost{store: store}
}

// UpdatePost is a command handler
type UpdatePost struct { // TODO Implement
	store usecase.PostStore
}

// Handle is CommandHandler's implementation
func (UpdatePost) Handle(context.Context, eventbus.Params) error {
	// TODO Implement
	return nil
}
