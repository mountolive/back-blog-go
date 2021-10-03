package httpx

import (
	"fmt"
	"net/http"
	"regexp"
)

type Middleware func(handler http.HandlerFunc) http.HandlerFunc

type Handler struct {
	handlerFunc http.HandlerFunc
	middlewares []Middleware
}

// Router registers URL paths against corresponding http handlers
type Router struct {
	globalMiddlewares []Middleware
	handlers          map[string]Handler
	cache             map[string]*regexp.Regexp
}

// NewRouter is a constructor
// middlewares passed to the router will be applied before per-route middlewares
func NewRouter(middlewares ...Middleware) *Router {
	return &Router{
		globalMiddlewares: middlewares,
		handlers:          make(map[string]Handler),
		cache:             make(map[string]*regexp.Regexp),
	}
}

// Add registers a handler with the corresponding path
// path can be a regexp
// middlewares' order matters: they'll be applied in insertion order
// path should be of the form: `HTTP_METHOD url`
func (r *Router) Add(
	path string, handler http.HandlerFunc, mws ...Middleware,
) error {
	r.handlers[path] = Handler{handler, mws}
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
	fmt.Println(reqPath)
	for path, handler := range r.handlers {
		if r.cache[path].MatchString(reqPath) {
			composed := handler.handlerFunc
			for _, gmw := range r.globalMiddlewares {
				composed = gmw(composed)
			}
			for _, mw := range handler.middlewares {
				composed = mw(composed)
			}
			composed(w, req)
			return
		}
	}
	fmt.Println("ROUTE NOT FOUND:", reqPath)
	http.NotFound(w, req)
}
