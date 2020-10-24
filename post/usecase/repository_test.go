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

type mockStoreNotEmpty struct{}

func (m *mockStoreNotEmpty) Create(ctx context.Context, p *CreatePostDto) (*PostDto, error) {
	return &PostDto{Creator: p.Creator, Content: p.Content}, nil
}

func (m *mockStoreNotEmpty) Update(ctx context.Context, p *UpdatePostDto) (*PostDto, error) {
	return &PostDto{Creator: "test", Content: p.Content}, nil
}

func (m *mockStoreNotEmpty) Filter(ctx context.Context, p *GeneralFilter) ([]*PostDto, error) {
	return []*PostDto{&PostDto{Creator: "test", Content: "test", Tags: []Tag{Tag{p.Tag}}}}, nil
}

type mockStoreEmpty struct{}

func (m *mockStoreEmpty) Create(ctx context.Context, p *CreatePostDto) (*PostDto, error) {
	return nil, nil
}

func (m *mockStoreEmpty) Update(ctx context.Context, p *UpdatePostDto) (*PostDto, error) {
	return nil, nil
}

func (m *mockStoreEmpty) Filter(ctx context.Context, p *GeneralFilter) ([]*PostDto, error) {
	return nil, nil
}

type mockSanitizer struct{}

func (m *mockSanitizer) SanitizeContent(ctx context.Context, content string) (string, error) {
	return content, nil
}

func TestPostRepository(t *testing.T) {
	repo := &PostRepository{Store: &mockStoreNotEmpty{}, Sanitizer: &mockSanitizer{}}
	genericError := "Got: %v; Expected: %v"
	t.Run("CreatePost", func(t *testing.T) {
		testDto := &CreatePostDto{"username", "content", []string{"tag1", "tag2"}}
		testCases := []createPostTestCase{
			{
				Name:        "Proper Create Post",
				Description: "It should return a *PostDto and no error",
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
				post, err := repo.CreatePost(ctx, tc.Dto)
				if tc.ExpErr != nil {
					require.True(t, post == nil, "Returned PostDto should be nil")
					require.True(t, err != nil, "Err should be not nil")
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
				Name:        "Post Not Found Update Post",
				Description: "It should return a nil and a PostNotFoundError",
				Dto:         testDto,
				ExpErr:      PostNotFoundError,
				Repo:        &PostRepository{Store: &mockStoreEmpty{}, Sanitizer: repo.Sanitizer},
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
				posts, err := repo.FilterByDateRange(ctx, &ByDateRangeDto{tc.Dto.From, tc.Dto.To})
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
