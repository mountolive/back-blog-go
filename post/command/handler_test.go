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
			store:       &mockStoreErrored{errors.New("create error")},
			params:      correctParams,
			expectedErr: errors.New("create error"),
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
			handler := command.NewCreatePost(tc.store, tc.checker)
			err := handler.Handle(context.Background(), tc.params)
			require.Equal(tc.expectedErr, err)
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
			name:        "Store Update error",
			description: "Errored execution when trying to excute store's Update",
			store:       &mockStoreErrored{errors.New("update error")},
			params:      correctParams,
			expectedErr: errors.New("update error"),
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
			handler := command.NewUpdatePost(tc.store)
			err := handler.Handle(context.Background(), tc.params)
			require.Equal(tc.expectedErr, err)
		})
	}
}

func TestIntegration(t *testing.T) {
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
	createHandler := command.NewCreatePost(store, &mockTrueChecker{})
	err := createHandler.Handle(ctx, correctParams)
	require.NoError(err)
	filter := &usecase.GeneralFilter{}
	filter.Tag = tag1
	createdPosts, err := store.Filter(ctx, filter)
	require.NoError(err)
	require.Len(createdPosts, 1)
	correctParams["id"] = createdPosts[0].Id
	correctParams["title"] = "some-other-title"
	updateHandler := command.NewUpdatePost(store)
	err = updateHandler.Handle(ctx, correctParams)
	require.NoError(err)
}
