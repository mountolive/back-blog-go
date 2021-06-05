package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mountolive/back-blog-go/post/command"
	"github.com/mountolive/back-blog-go/post/eventbus"
	"github.com/mountolive/back-blog-go/post/pgstore"
	"github.com/mountolive/back-blog-go/post/usecase"
	"github.com/stretchr/testify/require"
)

type createTestCase struct {
	name        string
	description string
	store       usecase.PostStore
	params      eventbus.Params
	expectedErr error
	checker     usecase.CreatorChecker
}

type updateTestCase struct {
	name        string
	description string
	store       usecase.PostStore
	params      eventbus.Params
	expectedErr error
}

func TestCreatePost(t *testing.T) {
	t.Run("Canary", func(t *testing.T) {
		t.Parallel()
		var _ eventbus.CommandHandler = command.CreatePost{}
	})

	require := require.New(t)
	correctParams := eventbus.Params{
		"creator": "some-creator",
		"title":   "title",
		"content": "some content",
		"tags":    []string{"tag1", "tag2"},
	}
	createErr := errors.New("create error")
	testCases := []createTestCase{
		{
			name:        "Missing content error",
			description: "Errored execution when payload is missing the content",
			params: eventbus.Params{
				"creator": "some-creator",
				"title":   "title",
				"tags":    []string{"tag1", "tag2"},
			},
			expectedErr: command.ErrContentMissing,
		},
		{
			name:        "Missing title error",
			description: "Errored execution when payload is missing the title",
			params: eventbus.Params{
				"creator": "some-creator",
				"content": "some content",
				"tags":    []string{"tag1", "tag2"},
			},
			expectedErr: command.ErrTitleMissing,
		},
		{
			name:        "Missing creator error",
			description: "Errored execution when payload is missing the creator",
			params: eventbus.Params{
				"title":   "title",
				"content": "some content",
				"tags":    []string{"tag1", "tag2"},
			},
			expectedErr: command.ErrCreatorMissing,
		},
		{
			name:        "Wrong Creator type error",
			description: "Errored execution when payload has a creator parameter that's not a string",
			params: eventbus.Params{
				"creator": []string{},
				"title":   "title",
				"content": "content",
				"tags":    []string{"tag1"},
			},
			expectedErr: command.NewErrWrongType("creator", "string"),
		},
		{
			name:        "Wrong Content type error",
			description: "Errored execution when payload has a content parameter that's not a string",
			params: eventbus.Params{
				"creator": "creator",
				"title":   "title",
				"content": 1,
				"tags":    []string{"tag1"},
			},
			expectedErr: command.NewErrWrongType("content", "string"),
		},
		{
			name:        "Wrong Title type error",
			description: "Errored execution when payload has a title parameter that's not a string",
			params: eventbus.Params{
				"creator": "creator",
				"title":   1,
				"content": "content",
				"tags":    []string{"tag1"},
			},
			expectedErr: command.NewErrWrongType("title", "string"),
		},
		{
			name:        "Wrong Tags type error",
			description: "Errored execution when payload has a tags parameter that's not a slice",
			params: eventbus.Params{
				"creator": "some-creator",
				"title":   "title",
				"content": "some content",
				"tags":    "tag1",
			},
			expectedErr: command.NewErrWrongType("tags", "[]string"),
		},
		{
			name:        "Creator checker error",
			description: "Errored execution when trying to retrieve a creator",
			checker:     &mockErroredChecker{errors.New("checker errored")},
			store:       &mockStore{},
			params:      correctParams,
			expectedErr: usecase.UserCheckError,
		},
		{
			name:        "Creator checker not found",
			description: "Creator not found",
			checker:     &mockFalseChecker{},
			store:       &mockStore{},
			params:      correctParams,
			expectedErr: usecase.UserNotFoundError,
		},
		{
			name:        "Store Create error",
			description: "Errored execution when trying to execute store's Create",
			checker:     &mockTrueChecker{},
			store:       &mockStoreErrored{createErr},
			params:      correctParams,
			expectedErr: createErr,
		},
		{
			name:        "Correct",
			description: "Not errored execution",
			store:       &mockStore{},
			checker:     &mockTrueChecker{},
			params:      correctParams,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)
			repo := &usecase.PostRepository{
				Store:     tc.store,
				Checker:   tc.checker,
				Sanitizer: &mockSanitizer{},
				Logger:    &mockLogger{},
			}
			handler := command.NewCreatePost(repo)
			err := handler.Handle(context.Background(), tc.params)
			require.True(errors.Is(err, tc.expectedErr), "got %v, expected %v", err, tc.expectedErr)
		})
	}
}

func TestUpdatePost(t *testing.T) {
	t.Run("Canary", func(t *testing.T) {
		t.Parallel()
		var _ eventbus.CommandHandler = command.UpdatePost{}
	})

	correctParams := eventbus.Params{
		"id":      "some-id",
		"title":   "title",
		"content": "some content",
		"tags":    []string{"tag1", "tag2"},
	}
	updateErr := errors.New("update error")
	require := require.New(t)
	testCases := []updateTestCase{
		{
			name:        "Missing id error",
			description: "Errored execution when payload is missing the id",
			store:       &mockStore{},
			params: eventbus.Params{
				"title":   "title",
				"content": "some content",
				"tags":    []string{"tag1", "tag2"},
			},
			expectedErr: command.ErrIDMissing,
		},
		{
			name:        "Missing content error",
			description: "Errored execution when payload is missing the content",
			params: eventbus.Params{
				"id":    "some-id",
				"title": "title",
				"tags":  []string{"tag1", "tag2"},
			},
			expectedErr: command.ErrContentMissing,
		},
		{
			name:        "Missing title error",
			description: "Errored execution when payload is missing the title",
			params: eventbus.Params{
				"id":      "some-id",
				"content": "some content",
				"tags":    []string{"tag1", "tag2"},
			},
			expectedErr: command.ErrTitleMissing,
		},
		{
			name:        "Wrong Id type error",
			description: "Errored execution when payload has a id parameter that's not a string",
			params: eventbus.Params{
				"id":      []string{},
				"title":   "title",
				"content": "content",
				"tags":    []string{"tag1"},
			},
			expectedErr: command.NewErrWrongType("id", "string"),
		},
		{
			name:        "Wrong Content type error",
			description: "Errored execution when payload has a content parameter that's not a string",
			params: eventbus.Params{
				"id":      "some-id",
				"title":   "title",
				"content": 1,
				"tags":    []string{"tag1"},
			},
			expectedErr: command.NewErrWrongType("content", "string"),
		},
		{
			name:        "Wrong Title type error",
			description: "Errored execution when payload has a title parameter that's not a string",
			params: eventbus.Params{
				"id":      "some-id",
				"title":   1,
				"content": "content",
				"tags":    []string{"tag1"},
			},
			expectedErr: command.NewErrWrongType("title", "string"),
		},
		{
			name:        "Wrong Tags type error",
			description: "Errored execution when payload has a tags parameter that's not a slice",
			params: eventbus.Params{
				"id":      "some-id",
				"title":   "title",
				"content": "some content",
				"tags":    "tag1",
			},
			expectedErr: command.NewErrWrongType("tags", "[]string"),
		},
		{
			name:        "Store Update error",
			description: "Errored execution when trying to excute store's Update",
			store:       &mockStoreErrored{updateErr},
			params:      correctParams,
			expectedErr: updateErr,
		},
		{
			name:        "Correct",
			description: "Not errored execution",
			store:       &mockStore{},
			params:      correctParams,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)
			repo := &usecase.PostRepository{
				Store:     tc.store,
				Checker:   &mockTrueChecker{},
				Sanitizer: &mockSanitizer{},
				Logger:    &mockLogger{},
			}
			handler := command.NewUpdatePost(repo)
			err := handler.Handle(context.Background(), tc.params)
			require.True(errors.Is(err, tc.expectedErr), "got %v, expected %v", err, tc.expectedErr)
		})
	}
}

func TestStoreIntegration(t *testing.T) {
	require := require.New(t)
	store := pgstore.CreateTestContainer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tag1 := "tag1"
	correctParams := eventbus.Params{
		"creator": "some-creator",
		"title":   "title",
		"content": "some content",
		"tags":    []string{tag1, "tag2"},
	}
	repo := &usecase.PostRepository{
		Store:     store,
		Checker:   &mockTrueChecker{},
		Sanitizer: &mockSanitizer{},
		Logger:    &mockLogger{},
	}
	createHandler := command.NewCreatePost(repo)
	err := createHandler.Handle(ctx, correctParams)
	require.NoError(err)
	filter := &usecase.GeneralFilter{PageSize: 1}
	filter.Tag = tag1
	createdPosts, err := store.Filter(ctx, filter)
	require.NoError(err)
	require.Len(createdPosts, 1)
	correctParams["id"] = createdPosts[0].Id
	correctParams["title"] = "some-other-title"
	updateHandler := command.NewUpdatePost(repo)
	err = updateHandler.Handle(ctx, correctParams)
	require.NoError(err)
}
