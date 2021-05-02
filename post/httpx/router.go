package httpx

import (
	"net/http"
	"regexp"
)

// Router registers URL paths against corresponding http handlers
type Router struct {
	handlers map[string]http.HandlerFunc
	cache    map[string]*regexp.Regexp
}

// NewRouter is a constructor
func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]http.HandlerFunc),
		cache:    make(map[string]*regexp.Regexp),
	}
}

// Add registers a handler with the corresponding path
// path can be a regexp
func (r *Router) Add(path string, handler http.HandlerFunc) error {
	r.handlers[path] = handler
	regComp, err := regexp.Compile(path)
	if err != nil {
		return err
	}
	r.cache[path] = regComp
	return nil
}

// ServeHTTP implements the http server interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqPath := req.Method + " " + req.URL.Path
	for path, handler := range r.handlers {
		if r.cache[path].MatchString(reqPath) {
			handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
}
