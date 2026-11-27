// Package mux provides HTTP route registration and middleware.
package mux

import (
	"net/http"
	"strings"
)

// Middleware represents an HTTP middleware.
type Middleware = func(http.Handler) http.Handler

var anyMethods = [...]string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
	http.MethodTrace, "QUERY",
}

// Router wraps http.ServeMux with a middleware chain.
type Router struct {
	mux        *http.ServeMux
	middleware []Middleware
}

// New creates an empty router.
func New() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Use appends global middleware.
//
// Middleware is applied in registration order,
// so the first middleware is the outermost.
func (r *Router) Use(middleware ...Middleware) {
	validateMiddlewares(middleware)
	r.middleware = append(r.middleware, middleware...)
}

// Group creates a route group with the given path prefix.
//
// The group's middleware is independent from the router's global middleware.
func (r *Router) Group(prefix string, middleware ...Middleware) *Group {
	validateMiddlewares(middleware)

	return &Group{
		router:     r,
		prefix:     normalizePrefix(prefix),
		middleware: append([]Middleware(nil), middleware...),
	}
}

// Handle registers a ServeMux pattern with optional route middleware.
func (r *Router) Handle(pattern string, handler http.Handler, middleware ...Middleware) {
	r.mux.Handle(pattern, wrapRouteHandler(handler, middleware))
}

// HandleFunc registers a handler function with optional route middleware.
func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc, middleware ...Middleware) {
	r.Handle(pattern, handler, middleware...)
}

// Any registers a route for GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE,
// CONNECT, TRACE, and QUERY.
func (r *Router) Any(path string, handler http.Handler, middleware ...Middleware) {
	r.Match(anyMethods[:], path, handler, middleware...)
}

// Match registers a route for each specified HTTP method.
func (r *Router) Match(methods []string, path string, handler http.Handler, middleware ...Middleware) {
	if len(methods) == 0 {
		return
	}
	registerMethodPatterns(r.mux, methods, joinPath("", path), wrapRouteHandler(handler, middleware))
}

// GET registers a GET route.
func (r *Router) GET(path string, handler http.Handler, middleware ...Middleware) {
	r.Match([]string{http.MethodGet}, path, handler, middleware...)
}

// HEAD registers a HEAD route.
func (r *Router) HEAD(path string, handler http.Handler, middleware ...Middleware) {
	r.Match([]string{http.MethodHead}, path, handler, middleware...)
}

// POST registers a POST route.
func (r *Router) POST(path string, handler http.Handler, middleware ...Middleware) {
	r.Match([]string{http.MethodPost}, path, handler, middleware...)
}

// PUT registers a PUT route.
func (r *Router) PUT(path string, handler http.Handler, middleware ...Middleware) {
	r.Match([]string{http.MethodPut}, path, handler, middleware...)
}

// PATCH registers a PATCH route.
func (r *Router) PATCH(path string, handler http.Handler, middleware ...Middleware) {
	r.Match([]string{http.MethodPatch}, path, handler, middleware...)
}

// DELETE registers a DELETE route.
func (r *Router) DELETE(path string, handler http.Handler, middleware ...Middleware) {
	r.Match([]string{http.MethodDelete}, path, handler, middleware...)
}

// OPTIONS registers an OPTIONS route.
func (r *Router) OPTIONS(path string, handler http.Handler, middleware ...Middleware) {
	r.Match([]string{http.MethodOptions}, path, handler, middleware...)
}

// Mux exposes the underlying ServeMux for infrastructure routes such as health
// checks.
func (r *Router) Mux() *http.ServeMux {
	return r.mux
}

// ServeHTTP applies the global middleware chain to the mux.
func (r *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	handler := applyMiddleware(r.mux, r.middleware)
	handler.ServeHTTP(w, request)
}

// Group represents a group of routes sharing a common path prefix
// and middleware chain.
type Group struct {
	router     *Router
	prefix     string
	middleware []Middleware
}

// Group creates a nested route group.
//
// The child group inherits the current group's middleware.
// Middleware added to the parent after child creation will not
// affect the already-created child group.
func (g *Group) Group(prefix string, middleware ...Middleware) *Group {
	validateMiddlewares(middleware)

	inherited := make([]Middleware, 0, len(g.middleware)+len(middleware))
	inherited = append(inherited, g.middleware...)
	inherited = append(inherited, middleware...)

	return &Group{
		router:     g.router,
		prefix:     joinGroupPrefix(g.prefix, prefix),
		middleware: inherited,
	}
}

// Use appends middleware to the group.
//
// Middleware only applies to routes registered after this call.
func (g *Group) Use(middleware ...Middleware) {
	validateMiddlewares(middleware)
	g.middleware = append(g.middleware, middleware...)
}

// Handle registers a route on the group with optional route middleware.
//
// Supports ServeMux patterns such as:
//
//	/users
//	GET /users
//	POST /users/{userID}
func (g *Group) Handle(pattern string, handler http.Handler, middleware ...Middleware) {
	fullPattern := joinPattern(g.prefix, pattern)
	g.router.mux.Handle(fullPattern, g.wrapHandler(handler, middleware))
}

// HandleFunc registers a handler function with optional route middleware.
func (g *Group) HandleFunc(pattern string, handler http.HandlerFunc, middleware ...Middleware) {
	g.Handle(pattern, handler, middleware...)
}

// HandleMethod registers a route with an HTTP method.
func (g *Group) HandleMethod(method, path string, handler http.Handler, middleware ...Middleware) {
	g.Match([]string{method}, path, handler, middleware...)
}

// Any registers a route for GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE,
// CONNECT, TRACE, and QUERY.
func (g *Group) Any(path string, handler http.Handler, middleware ...Middleware) {
	g.Match(anyMethods[:], path, handler, middleware...)
}

// Match registers a route for each specified HTTP method.
func (g *Group) Match(methods []string, path string, handler http.Handler, middleware ...Middleware) {
	if len(methods) == 0 {
		return
	}
	registerMethodPatterns(g.router.mux, methods, joinPath(g.prefix, path), g.wrapHandler(handler, middleware))
}

func (g *Group) wrapHandler(handler http.Handler, middleware []Middleware) http.Handler {
	return applyMiddleware(wrapRouteHandler(handler, middleware), g.middleware)
}

// GET registers a GET route.
func (g *Group) GET(path string, handler http.Handler, middleware ...Middleware) {
	g.HandleMethod(http.MethodGet, path, handler, middleware...)
}

// HEAD registers a HEAD route.
func (g *Group) HEAD(path string, handler http.Handler, middleware ...Middleware) {
	g.HandleMethod(http.MethodHead, path, handler, middleware...)
}

// POST registers a POST route.
func (g *Group) POST(path string, handler http.Handler, middleware ...Middleware) {
	g.HandleMethod(http.MethodPost, path, handler, middleware...)
}

// PUT registers a PUT route.
func (g *Group) PUT(path string, handler http.Handler, middleware ...Middleware) {
	g.HandleMethod(http.MethodPut, path, handler, middleware...)
}

// PATCH registers a PATCH route.
func (g *Group) PATCH(path string, handler http.Handler, middleware ...Middleware) {
	g.HandleMethod(http.MethodPatch, path, handler, middleware...)
}

// DELETE registers a DELETE route.
func (g *Group) DELETE(path string, handler http.Handler, middleware ...Middleware) {
	g.HandleMethod(http.MethodDelete, path, handler, middleware...)
}

// OPTIONS registers an OPTIONS route.
func (g *Group) OPTIONS(path string, handler http.Handler, middleware ...Middleware) {
	g.HandleMethod(http.MethodOptions, path, handler, middleware...)
}

func registerMethodPatterns(mux *http.ServeMux, methods []string, path string, handler http.Handler) {
	patterns := make([]string, 0, len(methods))
	for _, method := range methods {
		method = strings.TrimSpace(method)
		if method == "" {
			panic("router: empty HTTP method")
		}
		if strings.ContainsAny(method, " \t") {
			panic("router: invalid HTTP method: " + method)
		}
		patterns = append(patterns, method+" "+path)
	}
	for _, pattern := range patterns {
		mux.Handle(pattern, handler)
	}
}

func wrapRouteHandler(handler http.Handler, middleware []Middleware) http.Handler {
	if handler == nil {
		panic("router: nil handler")
	}
	if handlerFunc, ok := handler.(http.HandlerFunc); ok && handlerFunc == nil {
		panic("router: nil handler")
	}
	validateMiddlewares(middleware)
	return applyMiddleware(handler, middleware)
}

// joinGroupPrefix combines parent and child group prefixes.
func joinGroupPrefix(parent, child string) string {
	return normalizePrefix(joinPath(parent, child))
}

// normalizePrefix normalizes a group path prefix.
func normalizePrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	prefix = strings.Trim(prefix, "/")

	if prefix == "" {
		return ""
	}

	return "/" + prefix
}

// joinPath combines a group prefix and a route path.
//
// It preserves a trailing slash because ServeMux uses it
// to distinguish subtree routes.
func joinPath(prefix, path string) string {
	prefix = normalizePrefix(prefix)
	path = strings.TrimSpace(path)

	if path == "" {
		if prefix == "" {
			return "/"
		}
		return prefix
	}

	// A path beginning with ":" is an action suffix of the group resource,
	// for example /comments:batchStats. It must not introduce a slash.
	if strings.HasPrefix(path, ":") {
		return prefix + path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if prefix == "" {
		return path
	}

	if path == "/" {
		return prefix + "/"
	}

	return prefix + path
}

// joinPattern adds the group prefix to a ServeMux pattern.
//
// Supported forms:
//
//	/path
//	GET /path
//	example.com/path
//	GET example.com/path
func joinPattern(prefix, pattern string) string {
	parts := strings.Fields(strings.TrimSpace(pattern))

	if len(parts) == 0 || len(parts) > 2 {
		panic("router: invalid ServeMux pattern: " + pattern)
	}

	var method string
	var target string

	if len(parts) == 1 {
		target = parts[0]
	} else {
		method = parts[0]
		target = parts[1]
	}

	// Preserve optional host patterns supported by ServeMux.
	var host string
	if index := strings.IndexByte(target, '/'); index > 0 {
		host = target[:index]
		target = target[index:]
	}

	target = joinPath(prefix, target)

	if host != "" {
		target = host + target
	}

	if method != "" {
		return method + " " + target
	}

	return target
}

// applyMiddleware applies middleware in registration order.
// The first middleware becomes the outermost middleware.
func applyMiddleware(
	handler http.Handler,
	middleware []Middleware,
) http.Handler {
	for index := len(middleware) - 1; index >= 0; index-- {
		handler = middleware[index](handler)
	}

	return handler
}

// validateMiddlewares validates middleware registration.
func validateMiddlewares(middleware []Middleware) {
	for _, m := range middleware {
		if m == nil {
			panic("router: nil middleware")
		}
	}
}
