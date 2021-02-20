package handlers

import (
	"context"

	"github.com/mountolive/back-blog-go/post/usecase"
)

// Query is simply the name of the query, for lookup purposes
type Query string

// Command is simply the name of the command, for lookup purposes
type Command string

// CommandHandler defines the basic methods for a Command in the system
type CommandHandler interface {
	Handle(context.Context, Command) (usecase.PostDto, error)
}

// QueryOneHandler defines the basic methods for a Query of one element in the system
type QueryOneHandler interface {
	Handle(context.Context, Query) (usecase.PostDto, error)
}

// QueryManyHandler defines the basic methods for a Query of many elements in the system
type QueryManyHandler interface {
	Handle(context.Context, Query) ([]usecase.PostDto, error)
}
