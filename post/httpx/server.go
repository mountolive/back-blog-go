package httpx

import (
	"net/http"

	"github.com/mountolive/back-blog-go/post/usecase"
)

// Server contains all http handlers
type Server struct {
	repo usecase.PostRepository
}

// NewServer is a constructor
func NewServer(repo usecase.PostRepository) Server {
	return Server{repo}
}

// FilterByTag filters posts by tags
func (Server) FilterByTag(http.ResponseWriter, *http.Request) {
	// TODO Implement
}

// FilterByDateRange filters post by date range
func (Server) FilterByDateRange(http.ResponseWriter, *http.Request) {
	// TODO Implement
}

// Retrieves the details of a a single post
func (Server) GetPost(http.ResponseWriter, *http.Request) {
	// TODO Implement
}
