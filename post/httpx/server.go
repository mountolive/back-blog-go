package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mountolive/back-blog-go/post/usecase"
)

const (
	// default parameters
	defaultPage     = 0
	defaultPageSize = 10
	timeFormat      = "2006-01-02"
	// standard error codes
	RepositoryErrorCode             = 100
	MissingTagErrorCode             = 200
	NotFoundErrorCode               = 300
	MissingDateParametersErrorCode  = 400
	MarshalingErrorCode             = 500
	TimeParsingErrorCode            = 600
	EndTimeBeforeStartTimeErrorCode = 700

	// standard error messages
	MissingTagErrorMsg             = "tag parameter missing from query"
	NotFoundErrorMsg               = "post with passed id not found"
	MissingDateParametersErrorMsg  = "start_date and end_date parameters missing from query"
	EndTimeBeforeStartTimeErrorMsg = "end_date can't be before start_date"
)

// Server contains all http handlers
type Server struct {
	repo usecase.Repository
}

// NewServer is a constructor
func NewServer(repo usecase.Repository) Server {
	return Server{repo}
}

// APIError wraps the details of any error happened downstream to the http handler
type APIError struct {
	Error    DetailError `json:"error"`
	HTTPCode int         `json:"-"`
}

// DetailError is self-described
type DetailError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func newNotFoundError() APIError {
	return APIError{
		HTTPCode: http.StatusNotFound,
		Error: DetailError{
			Code:    NotFoundErrorCode,
			Message: NotFoundErrorMsg,
		},
	}
}

func newEndTimeBeforeStartTimeError() APIError {
	return APIError{
		HTTPCode: http.StatusBadRequest,
		Error: DetailError{
			Code:    EndTimeBeforeStartTimeErrorCode,
			Message: EndTimeBeforeStartTimeErrorMsg,
		},
	}
}

func newInternalServerError(code int, err error) APIError {
	return APIError{
		HTTPCode: http.StatusInternalServerError,
		Error: DetailError{
			Code:    code,
			Message: err.Error(),
		},
	}
}

func newRepositoryError(err error) APIError {
	return newInternalServerError(RepositoryErrorCode, err)
}

func newMarshalingError(err error) APIError {
	return newInternalServerError(MarshalingErrorCode, err)
}

func newTimeParsingError(err error) APIError {
	return newInternalServerError(TimeParsingErrorCode, err)
}

// Filter wraps both FilterByTag and FilterByDateRange
func (s Server) Filter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	tag := query.Get("tag")
	if tag != "" {
		s.FilterByTag(w, r)
		return
	}
	s.FilterByDateRange(w, r)
}

// FilterByTag filters posts by tag
func (s Server) FilterByTag(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	tag := query.Get("tag")
	page, pageSize := calculatePageAndPageSize(query)
	posts, err := s.repo.FilterByTag(
		r.Context(),
		&usecase.ByTagDto{Tag: tag},
		page,
		pageSize,
	)
	if err != nil {
		writeError(w, newRepositoryError(err))
		return
	}
	body, err := json.Marshal(posts)
	if err != nil {
		writeError(w, newMarshalingError(err))
		return
	}
	writeResponse(w, http.StatusOK, body)
}

// FilterByDateRange filters post by date range
func (s Server) FilterByDateRange(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	startDateRaw := query.Get("start_date")
	var (
		startDate time.Time
		endDate   time.Time
		err       error
	)
	if startDateRaw != "" {
		startDate, err = time.Parse(timeFormat, startDateRaw)
		if err != nil {
			writeError(w, newTimeParsingError(err))
			return
		}
	}
	endDateRaw := query.Get("end_date")
	if endDateRaw != "" {
		endDate, err = time.Parse(timeFormat, endDateRaw)
		if err != nil {
			writeError(w, newTimeParsingError(err))
			return
		}
	}
	// TODO Date logic for FilterByDateRange should be part of the Repository and not part of the handler
	if !startDate.IsZero() && !endDate.IsZero() && endDate.Before(startDate) {
		writeError(w, newEndTimeBeforeStartTimeError())
		return
	}
	page, pageSize := calculatePageAndPageSize(query)
	posts, err := s.repo.FilterByDateRange(
		r.Context(),
		&usecase.ByDateRangeDto{
			From: startDate,
			To:   endDate,
		},
		page,
		pageSize,
	)
	if err != nil {
		writeError(w, newRepositoryError(err))
		return
	}
	body, err := json.Marshal(posts)
	if err != nil {
		writeError(w, newMarshalingError(err))
		return
	}
	writeResponse(w, http.StatusOK, body)
}

func calculatePageAndPageSize(query url.Values) (int, int) {
	page := defaultPage
	pageRaw := query.Get("page")
	if pageRaw != "" {
		parsedPage, err := strconv.Atoi(pageRaw)
		if err == nil && parsedPage > -1 {
			page = parsedPage
		}
	}
	pageSize := defaultPageSize
	pageSizeRaw := query.Get("page_size")
	if pageSizeRaw != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeRaw)
		if err == nil && parsedPageSize > -1 {
			pageSize = parsedPageSize
		}
	}
	return page, pageSize
}

// Retrieves the details of a a single post
func (s Server) GetPost(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RequestURI()
	if len(url) > 0 && string(url[0]) == "/" {
		url = url[1:]
	}
	splittedPath := strings.Split(url, "/")
	if len(splittedPath) != 2 {
		writeError(w, newNotFoundError())
		return
	}
	id := strings.Split(splittedPath[1], "?")[0]
	post, err := s.repo.GetPost(r.Context(), id)
	if err != nil {
		writeError(w, newRepositoryError(err))
		return
	}
	if post == nil {
		writeError(w, newNotFoundError())
		return
	}
	body, err := json.Marshal(post)
	if err != nil {
		writeError(w, newMarshalingError(err))
		return
	}
	writeResponse(w, http.StatusOK, body)
}

func writeError(w http.ResponseWriter, apiError APIError) {
	body, _ := json.Marshal(apiError)
	writeResponse(w, apiError.HTTPCode, body)
}

func writeResponse(w http.ResponseWriter, httpCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	fmt.Fprint(w, string(body))
}
