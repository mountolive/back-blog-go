package httpx_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mountolive/back-blog-go/post/httpx"
	"github.com/mountolive/back-blog-go/post/usecase"
	"github.com/stretchr/testify/require"
)

func checkHandler(
	t *testing.T,
	route string,
	handler http.HandlerFunc,
	expStatusCode int,
	expectedBody []byte,
) {
	req := httptest.NewRequest(http.MethodGet, route, nil)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	require.Equal(t, expStatusCode, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, expectedBody, body)
}

func TestFilterByTag(t *testing.T) {
	t.Parallel()

	t.Run("Repository error, InternalServerError", func(t *testing.T) {
		defer recoverer(t)
		internalErrMsg := "something bad happened"
		repo := &RepositoryMock{
			FilterByTagFunc: func(context.Context, *usecase.ByTagDto, int, int) ([]*usecase.Post, error) {
				return nil, errors.New(internalErrMsg)
			},
		}
		server := httpx.NewServer(repo)
		expectedErr := httpx.APIError{
			HTTPCode: 500,
			Error: httpx.DetailError{
				Code:    100,
				Message: internalErrMsg,
			},
		}
		serializedErr, err := json.Marshal(expectedErr)
		require.NoError(t, err)
		checkHandler(
			t,
			"/someroute?tag=mytag",
			server.FilterByTag,
			http.StatusInternalServerError,
			serializedErr,
		)
	})

	expectedPosts := []*usecase.Post{
		&usecase.Post{
			Id:      "some-id",
			Content: "some content",
			Creator: "some creator",
		},
	}
	serializedBody, err := json.Marshal(expectedPosts)
	require.NoError(t, err)

	t.Run("Missing tag, OK", func(t *testing.T) {
		defer recoverer(t)
		repo := &RepositoryMock{
			FilterByTagFunc: func(context.Context, *usecase.ByTagDto, int, int) ([]*usecase.Post, error) {
				return expectedPosts, nil
			},
		}
		server := httpx.NewServer(repo)
		require.NoError(t, err)
		checkHandler(
			t,
			"/someroute",
			server.FilterByTag,
			http.StatusOK,
			serializedBody,
		)
	})

	t.Run("Correct, OK", func(t *testing.T) {
		defer recoverer(t)
		repo := &RepositoryMock{
			FilterByTagFunc: func(context.Context, *usecase.ByTagDto, int, int) ([]*usecase.Post, error) {
				return expectedPosts, nil
			},
		}
		server := httpx.NewServer(repo)
		require.NoError(t, err)
		checkHandler(
			t,
			"/someroute?tag=some-tag",
			server.FilterByTag,
			http.StatusOK,
			serializedBody,
		)
	})
}

func TestFilterByDateRange(t *testing.T) {
	t.Parallel()

	t.Run("Repository error, InternalServerError", func(t *testing.T) {
		defer recoverer(t)
		internalErrMsg := "something bad happened"
		repo := &RepositoryMock{
			FilterByDateRangeFunc: func(context.Context, *usecase.ByDateRangeDto, int, int) ([]*usecase.Post, error) {
				return nil, errors.New(internalErrMsg)
			},
		}
		server := httpx.NewServer(repo)
		expectedErr := httpx.APIError{
			HTTPCode: 500,
			Error: httpx.DetailError{
				Code:    100,
				Message: internalErrMsg,
			},
		}
		serializedErr, err := json.Marshal(expectedErr)
		require.NoError(t, err)
		checkHandler(
			t,
			"/someroute?start_date=2020-05-20&end_date=2020-05-21",
			server.FilterByDateRange,
			http.StatusInternalServerError,
			serializedErr,
		)
	})

	expectedPosts := []*usecase.Post{
		&usecase.Post{
			Id:      "some-id",
			Content: "some content",
			Creator: "some creator",
		},
	}
	serializedBody, err := json.Marshal(expectedPosts)
	require.NoError(t, err)

	t.Run("Missing start_date and end_date, OK", func(t *testing.T) {
		defer recoverer(t)
		repo := &RepositoryMock{
			FilterByDateRangeFunc: func(context.Context, *usecase.ByDateRangeDto, int, int) ([]*usecase.Post, error) {
				return expectedPosts, nil
			},
		}
		server := httpx.NewServer(repo)
		checkHandler(
			t,
			"/someroute",
			server.FilterByDateRange,
			http.StatusOK,
			serializedBody,
		)
	})

	t.Run("Correct only start_date parameter, OK", func(t *testing.T) {
		defer recoverer(t)
		repo := &RepositoryMock{
			FilterByDateRangeFunc: func(context.Context, *usecase.ByDateRangeDto, int, int) ([]*usecase.Post, error) {
				return expectedPosts, nil
			},
		}
		server := httpx.NewServer(repo)
		checkHandler(
			t,
			"/someroute?start_date=2020-05-21",
			server.FilterByDateRange,
			http.StatusOK,
			serializedBody,
		)
	})

	t.Run("Correct only end_date parameter, OK", func(t *testing.T) {
		defer recoverer(t)
		repo := &RepositoryMock{
			FilterByDateRangeFunc: func(context.Context, *usecase.ByDateRangeDto, int, int) ([]*usecase.Post, error) {
				return expectedPosts, nil
			},
		}
		server := httpx.NewServer(repo)
		checkHandler(
			t,
			"/someroute?end_date=2020-05-21",
			server.FilterByDateRange,
			http.StatusOK,
			serializedBody,
		)
	})

	t.Run("Correct both start_date and end_date parameter, OK", func(t *testing.T) {
		defer recoverer(t)
		repo := &RepositoryMock{
			FilterByDateRangeFunc: func(context.Context, *usecase.ByDateRangeDto, int, int) ([]*usecase.Post, error) {
				return expectedPosts, nil
			},
		}
		server := httpx.NewServer(repo)
		checkHandler(
			t,
			"/someroute?start_date=2020-05-19&end_date=2020-05-21",
			server.FilterByDateRange,
			http.StatusOK,
			serializedBody,
		)
	})
}

func TestGetOne(t *testing.T) {
	t.Parallel()

	t.Run("Repository error, InternalServerError", func(t *testing.T) {
		defer recoverer(t)
		internalErrMsg := "something bad happened"
		repo := &RepositoryMock{
			GetPostFunc: func(context.Context, string) (*usecase.Post, error) {
				return nil, errors.New(internalErrMsg)
			},
		}
		server := httpx.NewServer(repo)
		expectedErr := httpx.APIError{
			HTTPCode: 500,
			Error: httpx.DetailError{
				Code:    100,
				Message: internalErrMsg,
			},
		}
		serializedErr, err := json.Marshal(expectedErr)
		require.NoError(t, err)
		checkHandler(
			t,
			"/someroute/some-id",
			server.GetPost,
			http.StatusInternalServerError,
			serializedErr,
		)
	})

	t.Run("Unexistent Post, NotFound", func(t *testing.T) {
		defer recoverer(t)
		server := httpx.NewServer(nil)
		expectedErr := httpx.APIError{
			HTTPCode: 404,
			Error: httpx.DetailError{
				Code:    300,
				Message: "post with passed id not found",
			},
		}
		serializedErr, err := json.Marshal(expectedErr)
		require.NoError(t, err)
		checkHandler(
			t,
			"/someroute/not-found-id",
			server.GetPost,
			http.StatusNotFound,
			serializedErr,
		)
	})

	t.Run("Correct, OK", func(t *testing.T) {
		defer recoverer(t)
		expectedPost := &usecase.Post{
			Id:      "some-id",
			Content: "some content",
			Creator: "some creator",
		}
		repo := &RepositoryMock{
			GetPostFunc: func(context.Context, string) (*usecase.Post, error) {
				return expectedPost, nil
			},
		}
		server := httpx.NewServer(repo)
		serializedBody, err := json.Marshal(expectedPost)
		require.NoError(t, err)
		checkHandler(
			t,
			"/someroute/some-id",
			server.GetPost,
			http.StatusOK,
			serializedBody,
		)
	})
}
