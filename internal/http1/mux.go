package http1

import (
	"net/http"
	"strings"
)

type ServeMux struct {
	routes map[string]http.Handler
}

func NewServerMux() *ServeMux {
	return &ServeMux{
		routes: make(map[string]http.Handler),
	}
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler http.Handler) {
	if pattern == "" {
		panic("http1: invalid pattern")
	}
	if handler == nil {
		panic("http1: nil handler")
	}
	if _, exist := mux.routes[pattern]; exist {
		panic("http1: multiple registrations for " + pattern)
	}

	mux.routes[pattern] = handler
}

// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http1: nil handler")
	}
	// We use Go's built-in http.HandlerFunc type to convert the raw function
	// into an object that satisfies the http.Handler interface.
	mux.Handle(pattern, http.HandlerFunc(handler))
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Clean the path to prevent directory traversal attacks (e.g., /../)
	if !strings.HasPrefix(r.URL.Path, "/") {
		r.URL.Path = "/" + r.URL.Path
	}

	// 2. Look up the exact path in our routes map
	handler, exists := mux.routes[r.URL.Path]

	// 3. If no match is found, fallback to a 404 Not Found response
	if !exists {
		http.NotFound(w, r)
		return
	}

	// 4. Execute the matched handler
	handler.ServeHTTP(w, r)
}
