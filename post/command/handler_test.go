package command_test

import (
	"context"
	"testing"

	"github.com/mountolive/back-blog-go/post/command"
	"github.com/mountolive/back-blog-go/post/eventbus"
	"github.com/mountolive/back-blog-go/post/pgstore"
	"github.com/mountolive/back-blog-go/post/usecase"
	"github.com/stretchr/testify/require"
)

type testCase struct {
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
	testCases := []testCase{
		{
			name:        "Creator checker error",
			description: "Errored execution when trying to retrieve a creator",
			store:       &mockStore{},
			params: map[string]interface{}{
				"creator": "some-creator",
				"title":   "title",
				"content": "some content",
				"tags":    []string{"tag1", "tag2"},
			},
		},
		{
			name:        "Correct",
			description: "Not errored execution",
			store:       &mockStore{},
			params: map[string]interface{}{
				"creator": "some-creator",
				"title":   "title",
				"content": "some content",
				"tags":    []string{"tag1", "tag2"},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)
			handler := command.NewCreatePost(tc.store)
			err := handler.Handle(context.Background(), tc.params)
			require.Equal(tc.expectedErr, err)
		})
	}

	t.Run("Integration", func(t *testing.T) {
		pgstore.CreateTestContainer(t)
	})
}

func TestUpdatePost(t *testing.T) {
	t.Run("Canary", func(t *testing.T) {
		t.Parallel()
		var _ eventbus.CommandHandler = command.UpdatePost{}
	})

	require := require.New(t)
	testCases := []testCase{
		{
			name:        "Correct",
			description: "Not errored execution",
			store:       &mockStore{},
			params: map[string]interface{}{
				"id":      "some-id",
				"title":   "title",
				"content": "some content",
				"tags":    []string{"tag1", "tag2"},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)
			handler := command.NewCreatePost(tc.store)
			err := handler.Handle(context.Background(), tc.params)
			require.Equal(tc.expectedErr, err)
		})
	}

	t.Run("Integration", func(t *testing.T) {
		pgstore.CreateTestContainer(t)
	})
}
