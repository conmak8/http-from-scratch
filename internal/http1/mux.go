package http1

import (
	"net/http"
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
