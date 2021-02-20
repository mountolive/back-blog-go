package handlers

import (
	"context"

	"github.com/mountolive/back-blog-go/post/usecase"
)

// Params describes a basic query or command
type Params interface {
	Validate() error
}

// CommandHandler defines the basic methods for a Command in the system
type CommandHandler interface {
	Handle(context.Context, Params) (*usecase.PostDto, error)
}

// QueryOneHandler defines the basic methods for a Query of one element in the system
type QueryOneHandler interface {
	Handle(context.Context, Params) (*usecase.PostDto, error)
}

// QueryManyHandler defines the basic methods for a Query of many elements in the system
type QueryManyHandler interface {
	Handle(context.Context, Params) ([]*usecase.PostDto, error)
}

// CreatePostCommandHandler is self-described
type CreatePostCommandHandler struct {
	repo usecase.PostRepository
}

// UpdatePostCommandHandler is self-described
type UpdatePostCommandHandler struct {
	repo usecase.PostRepository
}

var (
	_ CommandHandler = CreatePostCommandHandler{}
	_ CommandHandler = UpdatePostCommandHandler{}
)

// Handle has the requirements to create a Post, according to the Params passed
func (c CreatePostCommandHandler) Handle(ctx context.Context, p Params) (*usecase.PostDto, error) {
	// TODO Implement
	return nil, nil
}

// Handle has the requirements to update a Post, according to the Params passed
func (u UpdatePostCommandHandler) Handle(ctx context.Context, p Params) (*usecase.PostDto, error) {
	// TODO Implement
	return nil, nil
}
