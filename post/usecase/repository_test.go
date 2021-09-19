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
}

type updatePostTestCase struct {
	Name        string
	Description string
	Dto         *UpdatePostDto
	ExpErr      error
	Repo        *PostRepository
}

type filterTestCase struct {
	Name        string
	Description string
	Dto         *GeneralFilter
	ExpErr      error
}

type getOneTestCase struct {
	Name        string
	Description string
	Errored     bool
	Repo        *PostRepository
}

func TestPostRepository(t *testing.T) {
	repo := &PostRepository{
		Store:     &mockStoreNotEmpty{},
		Sanitizer: &mockSanitizer{},
		Checker:   &mockTrueChecker{},
	}
	genericError := "Got: %v; Expected: %v"
	t.Run("CreatePost", func(t *testing.T) {
		testDto := &CreatePostDto{
			Title:   "title",
			Creator: "username",
			Content: "content",
			Tags:    []string{"tag1", "tag2"},
		}
		testCases := []createPostTestCase{
			{
				Name:        "Proper Create Post",
				Description: "It should return a *Post and no error",
				Dto:         testDto,
				Repo:        repo,
			},
			{
				Name:        "User non-existent",
				Description: "It should return a UserNotFoundError",
				Dto:         testDto,
				ExpErr:      ErrUserNotFound,
				Repo: &PostRepository{
					Store:     &mockStoreNotEmpty{},
					Sanitizer: &mockSanitizer{},
					Checker:   &mockFalseChecker{},
				},
			},
			{
				Name:        "Checker errored",
				Description: "It should return a UserCheckError",
				Dto:         testDto,
				ExpErr:      ErrUserCheck,
				Repo: &PostRepository{
					Store:     &mockStoreNotEmpty{},
					Sanitizer: &mockSanitizer{},
					Checker:   &mockErrorChecker{},
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx := context.Background()
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

	t.Run("GetPost", func(t *testing.T) {
		testCases := []getOneTestCase{
			{
				Name:        "Store read error",
				Description: "It should return wrap the error returned by the store",
				Errored:     true,
				Repo: &PostRepository{
					Store:     &mockStoreReadErrored{},
					Sanitizer: &mockSanitizer{},
					Checker:   &mockErrorChecker{},
				},
			},
			{
				Name:        "Post  not found",
				Description: "It should return an error indicating the post was not found",
				Errored:     true,
				Repo: &PostRepository{
					Store:     &mockStoreEmpty{},
					Sanitizer: &mockSanitizer{},
					Checker:   &mockErrorChecker{},
				},
			},
			{
				Name:        "Can return a *Post",
				Description: "It should return the found Post, from the store",
				Repo:        repo,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx := context.Background()
				post, err := tc.Repo.GetPost(ctx, "in-bloom")
				if tc.Errored {
					require.True(t, post == nil, "Returned PostDto should be nil")
					require.True(t, err != nil, "Err should be not nil")
					return
				}
				require.True(t, err == nil, "Err should be nil")
				// This depends on the mockStore used
				require.True(t, post != nil, "Returned PostDto shouldn't be nil")
			})
		}
	})

	t.Run("UpdatePost", func(t *testing.T) {
		testDto := &UpdatePostDto{
			Id:      "id",
			Title:   "title",
			Content: "content",
			Tags:    []string{"tag1", "tag2"},
		}
		testCases := []updatePostTestCase{
			{
				Name:        "Proper Update Post",
				Description: "It should return a *Post and no error",
				Dto:         testDto,
				Repo:        repo,
			},
			{
				Name:        "Missing Id",
				Description: "It should return a MissingIdError",
				Dto: &UpdatePostDto{
					Title:   "some title",
					Content: "some content",
					Tags:    []string{"many tags... (sic)"},
				},
				ExpErr: ErrMissingID,
				Repo:   repo,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx := context.Background()
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
				Description: "It should return a []*Post and no error",
				Dto:         testDto,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Log(tc.Description)
				ctx := context.Background()
				var posts []*Post
				var err error
				if tc.Dto.Tag == "" {
					posts, err = repo.FilterByDateRange(ctx, &ByDateRangeDto{tc.Dto.From, tc.Dto.To}, 0, 1)
				} else {
					posts, err = repo.FilterByTag(ctx, &ByTagDto{tc.Dto.Tag}, 0, 1)
				}
				if tc.ExpErr != nil {
					require.True(t, posts == nil, "Returned []*Post should be nil")
					require.True(t, err != nil, "Err should be not nil")
					require.True(t, errors.Is(err, tc.ExpErr), genericError, err, tc.ExpErr)
				} else {
					require.True(t, posts != nil, "[]*Post returned nil when it shouldn't")
					require.True(t, err == nil, "Err should be nil")
					// As the mock returns a single value
					require.Len(t, posts, 1, genericError, len(posts), 1)
					post := posts[0]
					require.True(t, post.Creator == "test", genericError, post.Creator, "test")
					require.True(t, post.Content == "test", genericError, post.Content, "test")
					firstTag := post.Tags[0]
					require.True(t, firstTag == tc.Dto.Tag, genericError, firstTag, tc.Dto.Tag)
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
