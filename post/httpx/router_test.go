package httpx_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/mountolive/back-blog-go/post/httpx"
	"github.com/stretchr/testify/require"
)

type mockResponseWriter struct{}

func (mockResponseWriter) Header() http.Header {
	return nil
}

func (mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (mockResponseWriter) WriteHeader(int) {}

func recoverer(t *testing.T) {
	if err := recover(); err != nil {
		t.Fatalf("panicking go routine: %v", err)
	}
}

func TestAdd(t *testing.T) {
	require := require.New(t)

	t.Run("Errored regexp compilation", func(t *testing.T) {
		badReg := "GET /bla(/?[A-Z]"
		router := httpx.NewRouter()
		err := router.Add(badReg, func(http.ResponseWriter, *http.Request) {})
		require.Error(err)
	})

	t.Run("Correct regex", func(t *testing.T) {
		defer recoverer(t)
		pathReg := "(GET|HEAD) /bye(/?[A-Za-z0-9]*)?"
		router := httpx.NewRouter()
		calls := 0
		hh := func(http.ResponseWriter, *http.Request) {
			calls += 1
		}
		err := router.Add(pathReg, hh)
		require.Nil(err)
		mockReqGet := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/bye/123abc"}}
		router.ServeHTTP(mockResponseWriter{}, mockReqGet)
		require.Exactly(1, calls)
		mockReqHead := &http.Request{Method: http.MethodHead, URL: &url.URL{Path: "/bye"}}
		router.ServeHTTP(mockResponseWriter{}, mockReqHead)
		require.Exactly(2, calls)
	})

	t.Run("Correct path", func(t *testing.T) {
		defer recoverer(t)
		pathReg := "POST /hello"
		router := httpx.NewRouter()
		calls := 0
		hh := func(http.ResponseWriter, *http.Request) {
			calls += 1
		}
		err := router.Add(pathReg, hh)
		require.Nil(err)
		mockReqCorrect := &http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/hello"}}
		router.ServeHTTP(mockResponseWriter{}, mockReqCorrect)
		require.Exactly(1, calls)
	})
}

func TestServeHTTP(t *testing.T) {
	defer recoverer(t)
	require := require.New(t)
	totalGM := 5
	counter := make([]int, totalGM)
	gmws := make([]httpx.Middleware, totalGM)
	for i := 0; i < totalGM; i++ {
		i := i
		gmws[i] = func(next http.HandlerFunc) http.HandlerFunc {
			counter[i] += 1
			return next
		}
	}
	pathReg := "POST /server"
	router := httpx.NewRouter(gmws...)
	handlerCalls := 0
	hh := func(http.ResponseWriter, *http.Request) {
		handlerCalls += 1
	}
	totalRM := 4
	routeCounter := make([]int, totalRM)
	rmws := make([]httpx.Middleware, totalRM)
	for i := 0; i < totalRM; i++ {
		i := i
		rmws[i] = func(next http.HandlerFunc) http.HandlerFunc {
			routeCounter[i] += 1
			return next
		}
	}
	err := router.Add(pathReg, hh, rmws...)
	require.Nil(err)
	reqCorrect := &http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/server"}}
	router.ServeHTTP(mockResponseWriter{}, reqCorrect)
	for _, g := range counter {
		require.Equal(1, g)
	}
	for _, r := range routeCounter {
		require.Equal(1, r)
	}
}
