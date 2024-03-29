package command

import (
	"context"
	"errors"
	"fmt"

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

// ErrWrongType is an error thrown when a type assertion fails on a field
type ErrWrongType struct {
	field        string
	expectedType string
}

// NewErrWrongType is a constructor
func NewErrWrongType(field, expType string) ErrWrongType {
	return ErrWrongType{field, expType}
}

// Error implements the error interface
func (e ErrWrongType) Error() string {
	return fmt.Sprintf("incorrect data type for %s, should be %s", e.field, e.expectedType)
}

const (
	creatorKey = "creator"
	contentKey = "content"
	titleKey   = "title"
	idKey      = "id"
	tagsKey    = "tags"
)

// NewCreatePost is a constructor
func NewCreatePost(repo usecase.Repository) CreatePost {
	return CreatePost{repo}
}

// CreatePostEventNameV1 is self-described
const CreatePostEventNameV1 = "posts.v1.create"

// CreatePost is a command handler
type CreatePost struct {
	repo usecase.Repository
}

var errCreatePostHandler = "create post: %w"

// Handle is CommandHandler's implementation
func (c CreatePost) Handle(ctx context.Context, params eventbus.Params) error {
	creatorParam, ok := params[creatorKey]
	if !ok {
		return fmt.Errorf(errCreatePostHandler, ErrCreatorMissing)
	}
	creator, ok := creatorParam.(string)
	if !ok {
		return fmt.Errorf(
			errCreatePostHandler,
			NewErrWrongType("creator", "string"),
		)
	}
	contentParam, ok := params[contentKey]
	if !ok {
		return fmt.Errorf(errCreatePostHandler, ErrContentMissing)
	}
	content, ok := contentParam.(string)
	if !ok {
		return fmt.Errorf(
			errCreatePostHandler,
			NewErrWrongType("content", "string"),
		)
	}
	titleParam, ok := params[titleKey]
	if !ok {
		return fmt.Errorf(errCreatePostHandler, ErrTitleMissing)
	}
	title, ok := titleParam.(string)
	if !ok {
		return fmt.Errorf(
			errCreatePostHandler,
			NewErrWrongType("title", "string"),
		)
	}
	tags, err := extractTags(params)
	if err != nil {
		return fmt.Errorf(errCreatePostHandler, err)
	}
	createPost := &usecase.CreatePostDto{
		Creator: creator,
		Content: content,
		Title:   title,
		Tags:    tags,
	}
	_, err = c.repo.CreatePost(ctx, createPost)
	if err != nil {
		return fmt.Errorf(errCreatePostHandler, err)
	}
	return nil
}

// NewUpdatePost is a constructor
func NewUpdatePost(repo usecase.Repository) UpdatePost {
	return UpdatePost{repo: repo}
}

// UpdatePostEventNameV1 is self-described
const UpdatePostEventNameV1 = "posts.v1.update"

// UpdatePost is a command handler
type UpdatePost struct {
	repo usecase.Repository
}

var errUpdatePostHandler = "update post: %w"

// Handle is CommandHandler's implementation
func (u UpdatePost) Handle(ctx context.Context, params eventbus.Params) error {
	idParam, ok := params[idKey]
	if !ok {
		return fmt.Errorf(errUpdatePostHandler, ErrIDMissing)
	}
	id, ok := idParam.(string)
	if !ok {
		return fmt.Errorf(
			errUpdatePostHandler,
			NewErrWrongType("id", "string"),
		)
	}
	contentParam, ok := params[contentKey]
	if !ok {
		return fmt.Errorf(errUpdatePostHandler, ErrContentMissing)
	}
	content, ok := contentParam.(string)
	if !ok {
		return fmt.Errorf(
			errUpdatePostHandler,
			NewErrWrongType("content", "string"),
		)
	}
	titleParam, ok := params[titleKey]
	if !ok {
		return fmt.Errorf(errUpdatePostHandler, ErrTitleMissing)
	}
	title, ok := titleParam.(string)
	if !ok {
		return fmt.Errorf(
			errUpdatePostHandler,
			NewErrWrongType("title", "string"),
		)
	}
	tags, err := extractTags(params)
	if err != nil {
		return fmt.Errorf(errUpdatePostHandler, err)
	}
	updatePost := &usecase.UpdatePostDto{
		Id:      id,
		Content: content,
		Title:   title,
		Tags:    tags,
	}
	_, err = u.repo.UpdatePost(ctx, updatePost)
	if err != nil {
		return fmt.Errorf(errUpdatePostHandler, err)
	}
	return nil
}

func extractTags(params map[string]interface{}) ([]string, error) {
	var tags []string
	_, ok := params[tagsKey]
	if ok {
		tagValues, ok := params[tagsKey].([]interface{})
		if !ok {
			return nil, NewErrWrongType("tags", "array")
		}
		for i, rawTag := range tagValues {
			tag, ok := rawTag.(string)
			if !ok {
				return nil, NewErrWrongType(fmt.Sprintf("tag at %d", i), "string")
			}
			tags = append(tags, tag)
		}
	}
	return tags, nil
}
