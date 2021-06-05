package httpx

import (
	"net/http"

	"github.com/mountolive/back-blog-go/post/usecase"
)

const (
	// Standard Error Codes
	RepositoryErrorCode            = 100
	MissingTagErrorCode            = 200
	NotFoundErrorCode              = 300
	MissingDateParametersErrorCode = 400

	// Standard Error Messages
	MissingTagErrorMsg            = "tag parameter missing from query"
	NotFoundErrorMsg              = "post with passed id not found"
	MissingDateParametersErrorMsg = "start_date and end_date parameters missing from query"
)

// Server contains all http handlers
type Server struct {
	repo usecase.Repository
}

// APIError wraps the details of any error happened downstream to the http handler
type APIError struct {
	HTTPCode int         `json:"-"`
	Error    DetailError `json:"errors"`
}

// DetailError is self-described
type DetailError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewServer is a constructor
func NewServer(repo usecase.Repository) Server {
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
