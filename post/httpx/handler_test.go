package httpx_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mountolive/back-blog-go/post/httpx"
	"github.com/mountolive/back-blog-go/post/usecase"
	"github.com/stretchr/testify/require"
)

func checkHandler(
	t *testing.T, repo usecase.PostRepository, route string,
	handler http.HandlerFunc, expStatusCode int, expectedBody []byte,
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
		repo := usecase.PostRepository{
			// TODO Add erroring condition, Tag filter
		}
		server := httpx.NewServer(repo)
		// TODO Create expectedBody, Repository error, tag
		checkHandler(t, repo, "/someroute?tag=mytag", server.FilterByTag, http.StatusInternalServerError, []byte{})
	})

	t.Run("Missing tag, NotFound", func(t *testing.T) {
		defer recoverer(t)
		repo := usecase.PostRepository{
			// TODO Add mocked returned Posts, Missing tag
		}
		server := httpx.NewServer(repo)
		// TODO Create expectedBody, Missing tag
		checkHandler(t, repo, "/someroute", server.FilterByTag, http.StatusNotFound, []byte{})
	})

	t.Run("Correct, OK", func(t *testing.T) {
		defer recoverer(t)
		repo := usecase.PostRepository{
			// TODO Add mocked returned Posts, Correct, tag
		}
		server := httpx.NewServer(repo)
		// TODO Create expectedBody, Correct, tag
		checkHandler(t, repo, "/someroute?tag=sometag", server.FilterByTag, http.StatusOK, []byte{})
	})
}

func TestFilterByDateRange(t *testing.T) {
	t.Parallel()

	t.Run("Repository error, InternalServerError", func(t *testing.T) {
		defer recoverer(t)
		repo := usecase.PostRepository{
			// TODO Add erroring condition, Date range filter
		}
		server := httpx.NewServer(repo)
		// TODO Create expectedBody, Repository error, date range
		// TODO Change route to hold date range params
		checkHandler(t, repo, "/someroute?tag=mytag", server.FilterByDateRange, http.StatusInternalServerError, []byte{})
	})

	t.Run("Missing tag, NotFound", func(t *testing.T) {
		defer recoverer(t)
		repo := usecase.PostRepository{
			// TODO Add mocked returned Posts, Missing date range
		}
		server := httpx.NewServer(repo)
		// TODO Create expectedBody, Missing date range
		// TODO Change route to hold date range params
		checkHandler(t, repo, "/someroute", server.FilterByDateRange, http.StatusNotFound, []byte{})
	})

	t.Run("Correct, OK", func(t *testing.T) {
		defer recoverer(t)
		repo := usecase.PostRepository{
			// TODO Add mocked returned Posts, Correct, date range
		}
		server := httpx.NewServer(repo)
		// TODO Create expectedBody, Correct, date range
		// TODO Change route to hold date range params
		checkHandler(t, repo, "/someroute?tag=sometag", server.FilterByDateRange, http.StatusOK, []byte{})
	})
}

func TestGetOne(t *testing.T) {
	t.Parallel()

	t.Run("Repository error, InternalServerError", func(t *testing.T) {
		defer recoverer(t)
		repo := usecase.PostRepository{
			// TODO Add erroring condition, by id
		}
		server := httpx.NewServer(repo)
		// TODO Create expectedBody, Repository error, by id
		checkHandler(t, repo, "/someroute/the-id", server.GetPost, http.StatusInternalServerError, []byte{})
	})

	t.Run("Missing id, NotFound", func(t *testing.T) {
		defer recoverer(t)
		repo := usecase.PostRepository{
			// TODO Add mocked returned Posts, by id
		}
		server := httpx.NewServer(repo)
		// TODO Create expectedBody, Missing id
		checkHandler(t, repo, "/someroute", server.GetPost, http.StatusNotFound, []byte{})
	})

	t.Run("Correct, OK", func(t *testing.T) {
		defer recoverer(t)
		repo := usecase.PostRepository{
			// TODO Add mocked returned Posts, Correct, by id
		}
		server := httpx.NewServer(repo)
		// TODO Create expectedBody, Correct, by id
		checkHandler(t, repo, "/someroute/some-id", server.GetPost, http.StatusOK, []byte{})
	})
}
