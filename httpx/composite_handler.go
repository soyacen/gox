package httpx

import (
	"net/http"
)

// CompositeHandler is an HTTP handler that routes requests to different handlers
// based on a matching function.
type CompositeHandler struct {
	matchableHandlers []*matchableHandler
}

// AddHandler adds a handler with a matching function to the composite handler.
// The match function determines whether the handler should serve the request.
//
// Parameters:
//   - handler: The HTTP handler to add
//   - match: A function that returns true if the handler should handle the request
func (h *CompositeHandler) AddHandler(handler http.Handler, match func(request *http.Request) bool) {
	h.matchableHandlers = append(h.matchableHandlers, &matchableHandler{handler: handler, match: match})
}

// ServeHTTP implements the http.Handler interface.
// It iterates through all registered handlers and dispatches to the first one
// whose match function returns true.
//
// Parameters:
//   - writer: The response writer
//   - request: The HTTP request
func (h *CompositeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, matchableHandler := range h.matchableHandlers {
		if matchableHandler.match(request) {
			matchableHandler.handler.ServeHTTP(writer, request)
			return
		}
	}
}

// matchableHandler pairs an HTTP handler with a matching function.
type matchableHandler struct {
	handler http.Handler
	match   func(request *http.Request) bool
}

// ServeHTTP delegates to the underlying handler.
//
// Parameters:
//   - writer: The response writer
//   - request: The HTTP request
func (h *matchableHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.handler.ServeHTTP(writer, request)
}
