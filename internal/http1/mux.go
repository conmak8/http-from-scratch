package http1

import "net/http"

type ServeMux struct {
	routes map[string]http.Handler
}

func NewServerMux() *ServeMux {
	return &ServeMux{
		routes: make(map[string]http.Handler),
	}
}
