package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type createPostTestCase struct {
	Name        string
	Description string
	Dto         *CreatePostDto
	ExpErr      error
	Repo        *PostRepository
	CtxCancel   bool
}

type updatePostTestCase struct {
	Name        string
	Description string
	Dto         *UpdatePostDto
	ExpErr      error
	CtxCancel   bool
	Repo        *PostRepository
}

type filterTestCase struct {
	Name        string
	Description string
	Dto         *GeneralFilter
	ExpErr      error
	CtxCancel   bool
}

func TestPostRepository(t *testing.T) {
	logger := &mockLogger{}
	repo := &PostRepository{
		Store:     &mockStoreNotEmpty{},
		Sanitizer: &mockSanitizer{},
		Checker:   &mockTrueChecker{},
		Logger:    logger,
	}
	genericError := "Got: %v; Expected: %v"
	t.Run("CreatePost", func(t *testing.T) {
		testDto := &CreatePostDto{"username", "content", []string{"tag1", "tag2"}}
		testCases := []createPostTestCase{
			{
				Name:        "Proper Create Post",
				Description: "It should return a *PostDto and no error",
				Dto:         testDto,
				Repo:        repo,
			},
			{
				Name:        "User non-existent",
				Description: "It should return a UserNotFoundError",
				Dto:         testDto,
				ExpErr:      UserNotFoundError,
				Repo: &PostRepository{
					Store:     &mockStoreNotEmpty{},
					Sanitizer: &mockSanitizer{},
					Checker:   &mockFalseChecker{},
					Logger:    logger,
				},
			},
			{
				Name:        "Checker errored",
				Description: "It should return a UserCheckError",
				Dto:         testDto,
				ExpErr:      UserCheckError,
				Repo: &PostRepository{
					Store:     &mockStoreNotEmpty{},
					Sanitizer: &mockSanitizer{},
					Checker:   &mockErrorChecker{},
					Logger:    logger,
				},
			},
			{
				Name:        "Context Canceled",
				Description: "It should return an OperationCanceledError",
				Dto:         testDto,
				ExpErr:      OperationCanceledError,
				CtxCancel:   true,
				Repo:        repo,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tc.CtxCancel {
					cancel()
				}
				post, err := tc.Repo.CreatePost(ctx, tc.Dto)
				if tc.ExpErr != nil {
					require.True(t, post == nil, "Returned PostDto should be nil")
					require.True(t, err != nil, "Err should be not nil")
					require.True(t, errors.Is(err, tc.ExpErr), genericError, err, tc.ExpErr)
				} else {
					require.True(t, err == nil, "Err should be nil")
					creator := tc.Dto.Creator
					require.Equal(t, creator, post.Creator, genericError, post.Creator, creator)
					content := tc.Dto.Content
					require.Equal(t, content, post.Content, genericError, post.Content, content)
				}
			})
		}
	})

	t.Run("UpdatePost", func(t *testing.T) {
		testDto := &UpdatePostDto{"id", "content", []string{"tag1", "tag2"}}
		testCases := []updatePostTestCase{
			{
				Name:        "Proper Update Post",
				Description: "It should return a *PostDto and no error",
				Dto:         testDto,
				Repo:        repo,
			},
			{
				Name:        "Context Canceled",
				Description: "It should return an OperationCanceledError",
				Dto:         testDto,
				ExpErr:      OperationCanceledError,
				CtxCancel:   true,
				Repo:        repo,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tc.CtxCancel {
					cancel()
				}
				post, err := tc.Repo.UpdatePost(ctx, tc.Dto)
				if tc.ExpErr != nil {
					require.True(t, post == nil, "Returned PostDto should be nil")
					require.True(t, err != nil, "Err should be not nil")
					require.True(t, errors.Is(err, tc.ExpErr), genericError, err, tc.ExpErr)
				} else {
					id := tc.Dto.Id
					require.Equal(t, id, post.Id, genericError, post.Id, id)
					content := tc.Dto.Content
					require.Equal(t, content, post.Content, genericError, post.Content, content)
				}
			})
		}
	})

	testFilter := func(t *testing.T, testDto *GeneralFilter) {
		testCases := []filterTestCase{
			{
				Name:        "Proper Filtering",
				Description: "It should return a []*PostDto and no error",
				Dto:         testDto,
			},
			{
				Name:        "Context Canceled",
				Description: "It should return an OperationCanceledError",
				Dto:         testDto,
				ExpErr:      OperationCanceledError,
				CtxCancel:   true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tc.CtxCancel {
					cancel()
				}
				var posts []*PostDto
				var err error
				if tc.Dto.Tag == "" {
					posts, err = repo.FilterByDateRange(ctx, &ByDateRangeDto{tc.Dto.From, tc.Dto.To})
				} else {
					posts, err = repo.FilterByTag(ctx, &ByTagDto{tc.Dto.Tag})
				}
				if tc.ExpErr != nil {
					require.True(t, posts == nil, "Returned []*PostDto should be nil")
					require.True(t, err != nil, "Err should be not nil")
					require.True(t, errors.Is(err, tc.ExpErr), genericError, err, tc.ExpErr)
				} else {
					require.True(t, posts != nil, "[]*PostDto returned nil when it shouldn't")
					require.True(t, err == nil, "Err should be nil")
					// As the mock returns a single value
					length := len(posts)
					require.True(t, length == 1, genericError, length, 1)
					post := posts[0]
					require.True(t, post.Creator == "test", genericError, post.Creator, "test")
					require.True(t, post.Content == "test", genericError, post.Content, "test")
					firstTag := post.Tags[0]
					require.True(t, firstTag.Name == tc.Dto.Tag, genericError, firstTag.Name, tc.Dto.Tag)
				}
			})
		}
	}

	t.Run("FilterByDateRange", func(t *testing.T) {
		testDto := &GeneralFilter{}
		testDto.From = time.Now()
		testDto.To = testDto.From.Add(5 * time.Hour)
		testFilter(t, testDto)
	})

	t.Run("FilterByTag", func(t *testing.T) {
		testDto := &GeneralFilter{}
		testDto.Tag = "test"
		testFilter(t, testDto)
	})
}
