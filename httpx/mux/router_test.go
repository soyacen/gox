package mux

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRouterMiddlewareOrderAndRoutes(t *testing.T) {
	r := New()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			w.Header().Add("X-Order", "first")
			next.ServeHTTP(w, request)
		})
	})
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			w.Header().Add("X-Order", "second")
			next.ServeHTTP(w, request)
		})
	})
	r.Handle("POST /test", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/test", nil))
	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d", recorder.Code)
	}
	values := recorder.Header().Values("X-Order")
	if len(values) != 2 || values[0] != "first" || values[1] != "second" {
		t.Fatalf("middleware order = %v", values)
	}

	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/test", nil))
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET status = %d, want 405", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/missing", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("missing route status = %d, want 404", recorder.Code)
	}
}

func TestAdapMiddlewareOrderAndRegistration(t *testing.T) {
	var trace []string
	traceFunc := func(name string) MiddlewareFunc {
		return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			trace = append(trace, name+":before")
			next(w, r)
			trace = append(trace, name+":after")
		}
	}

	middlewares := []MiddlewareFunc{traceFunc("global-first"), traceFunc("global-second")}
	r := New()
	r.Use(AdaptMiddlewareFunc(middlewares...), traceMiddleware("global-existing", &trace))
	middlewares[0] = nil

	g := r.Group("/api")
	g.Use(AdaptMiddlewareFunc(traceFunc("group")))
	child := g.Group("/v1")
	g.Use(AdaptMiddlewareFunc(traceFunc("group-late")))
	child.GET("/users", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		trace = append(trace, "handler")
		w.WriteHeader(http.StatusNoContent)
	}), AdaptMiddlewareFunc(traceFunc("route")))

	recorder := performRequest(r, http.MethodGet, "/api/v1/users")
	assertStatus(t, recorder, http.StatusNoContent)
	assertTrace(t, trace, []string{
		"global-first:before", "global-second:before", "global-existing:before", "group:before", "route:before",
		"handler",
		"route:after", "group:after", "global-existing:after", "global-second:after", "global-first:after",
	})
}

func TestAdapMiddlewareForwardsArgumentsAndCanStop(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/original", nil)
	forwardedRequest := httptest.NewRequest(http.MethodGet, "/forwarded", nil)
	forwardedWriter := httptest.NewRecorder()
	middleware := AdaptMiddlewareFunc(
		func(_ http.ResponseWriter, _ *http.Request, next http.HandlerFunc) {
			next(forwardedWriter, forwardedRequest)
		},
		func(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
			if w != forwardedWriter || r != forwardedRequest {
				t.Fatal("next received the original writer or request")
			}
			w.WriteHeader(http.StatusForbidden)
		},
	)
	handler := middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("handler called after middleware stopped the chain")
	}))
	handler.ServeHTTP(httptest.NewRecorder(), request)
	assertStatus(t, forwardedWriter, http.StatusForbidden)
}

func TestAdapMiddlewareEmptyAndNil(t *testing.T) {
	next := http.NewServeMux()
	if got := AdaptMiddlewareFunc()(next); got != next {
		t.Fatal("empty middleware chain changed the downstream handler")
	}
	assertPanics(t, func() {
		AdaptMiddlewareFunc(func(http.ResponseWriter, *http.Request, http.HandlerFunc) {}, nil)
	})
}

func TestRouterExposesMuxForInfrastructureRoutes(t *testing.T) {
	r := New()
	r.Mux().HandleFunc("GET /infra-probe", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Probe", "health")
		w.WriteHeader(http.StatusOK)
	})
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/infra-probe", nil))
	if recorder.Code != http.StatusOK || recorder.Header().Get("X-Probe") != "health" {
		t.Fatalf("health status = %d", recorder.Code)
	}
}

func TestRouterServeHTTPWithoutMiddleware(t *testing.T) {
	r := New()
	r.Handle("GET /plain", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/plain", nil))
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", recorder.Code)
	}
}

// TestMuxServedDirectlyBypassesMiddleware documents that middleware is only
// applied through Router.ServeHTTP. Handlers reached via the exposed Mux served
// on its own skip the chain, which is how infrastructure probes stay cheap.
func TestMuxServedDirectlyBypassesMiddleware(t *testing.T) {
	r := New()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-MW", "applied")
			next.ServeHTTP(w, req)
		})
	})
	r.Mux().Handle("GET /direct", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	throughMux := httptest.NewRecorder()
	r.Mux().ServeHTTP(throughMux, httptest.NewRequest(http.MethodGet, "/direct", nil))
	if throughMux.Header().Get("X-MW") != "" {
		t.Fatalf("direct mux serve should bypass middleware, got %q", throughMux.Header().Get("X-MW"))
	}

	throughRouter := httptest.NewRecorder()
	r.ServeHTTP(throughRouter, httptest.NewRequest(http.MethodGet, "/direct", nil))
	if throughRouter.Header().Get("X-MW") != "applied" {
		t.Fatalf("Router.ServeHTTP should apply middleware, got %q", throughRouter.Header().Get("X-MW"))
	}
}

func TestRouterHandleFunc(t *testing.T) {
	r := New()
	r.HandleFunc("GET /func", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/func", nil))
	if recorder.Code != http.StatusTeapot {
		t.Fatalf("status = %d, want 418", recorder.Code)
	}
}

// TestRouterRoutingSemantics locks in the ServeMux behaviors the whole app
// depends on: wildcard PathValue extraction, static-over-wildcard precedence,
// AIP custom-method paths that contain a colon, and the 405/404 distinction.
func TestRouterRoutingSemantics(t *testing.T) {
	cases := []struct {
		name     string
		patterns map[string]string // method+schema -> marker
		request  [2]string         // method, path
		want     int
		wantBody string
	}{
		{
			name:     "wildcard extracts path value",
			patterns: map[string]string{"GET /users/{uid}": "uid"},
			request:  [2]string{http.MethodGet, "/users/42"},
			want:     http.StatusOK,
			wantBody: "uid=42",
		},
		{
			name: "static beats wildcard at same depth",
			patterns: map[string]string{
				"GET /quotas/me/{quota_code}": "bycode",
				"GET /quotas/me/usage":        "usage",
			},
			request:  [2]string{http.MethodGet, "/quotas/me/usage"},
			want:     http.StatusOK,
			wantBody: "usage",
		},
		{
			name: "wildcard still reachable when static sibling exists",
			patterns: map[string]string{
				"GET /quotas/me/{quota_code}": "bycode",
				"GET /quotas/me/usage":        "usage",
			},
			request:  [2]string{http.MethodGet, "/quotas/me/coins"},
			want:     http.StatusOK,
			wantBody: "bycode",
		},
		{
			name: "differing wildcard names at same position route independently",
			patterns: map[string]string{
				"GET /comments/{comment_id}":      "detail:comment_id",
				"GET /comments/{root_id}/replies": "replies:root_id",
			},
			request:  [2]string{http.MethodGet, "/comments/9/replies"},
			want:     http.StatusOK,
			wantBody: "replies:root_id=9",
		},
		{
			name:     "AIP custom method with colon matches literal path",
			patterns: map[string]string{"POST /favorites:batchGet": "batchGet"},
			request:  [2]string{http.MethodPost, "/favorites:batchGet"},
			want:     http.StatusOK,
			wantBody: "batchGet",
		},
		{
			name:     "registered path wrong method returns 405",
			patterns: map[string]string{"POST /things": "create"},
			request:  [2]string{http.MethodOptions, "/things"},
			want:     http.StatusMethodNotAllowed,
		},
		{
			name:     "absent path returns 404",
			patterns: map[string]string{"POST /things": "create"},
			request:  [2]string{http.MethodGet, "/nothing-here"},
			want:     http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := New()
			for pattern, marker := range tc.patterns {
				value := marker
				r.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					switch value {
					case "uid":
						w.Write([]byte("uid=" + req.PathValue("uid")))
					case "detail:comment_id":
						w.Write([]byte("detail:comment_id=" + req.PathValue("comment_id")))
					case "replies:root_id":
						w.Write([]byte("replies:root_id=" + req.PathValue("root_id")))
					default:
						w.Write([]byte(value))
					}
				}))
			}
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, httptest.NewRequest(tc.request[0], tc.request[1], nil))
			if recorder.Code != tc.want {
				t.Fatalf("status = %d, want %d", recorder.Code, tc.want)
			}
			if tc.wantBody != "" && recorder.Body.String() != tc.wantBody {
				t.Fatalf("body = %q, want %q", recorder.Body.String(), tc.wantBody)
			}
		})
	}
}

// TestRouterMiddlewareAppliesAcrossEveryRoute confirms the global chain wraps all
// registered handlers, not just the first route, and that PathValue survives the
// wrapping so downstream handlers can still read wildcard segments.
func TestRouterMiddlewareAppliesAcrossEveryRoute(t *testing.T) {
	r := New()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Touched", "yes")
			next.ServeHTTP(w, req)
		})
	})
	r.Handle("GET /a/{id}", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(req.PathValue("id")))
	}))
	r.Handle("GET /b", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("b"))
	}))

	for path, want := range map[string]string{"/a/7": "7", "/b": "b"} {
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Header().Get("X-Touched") != "yes" {
			t.Errorf("%s: middleware not applied", path)
		}
		if recorder.Body.String() != want {
			t.Errorf("%s: body = %q, want %q", path, recorder.Body.String(), want)
		}
	}
}

// -----------------------------------------------------------------------------
// Test helpers
// -----------------------------------------------------------------------------

// performRequest executes an HTTP request against the given handler.
func performRequest(
	handler http.Handler,
	method string,
	target string,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	return recorder
}

// traceMiddleware records middleware execution order.
func traceMiddleware(name string, trace *[]string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			*trace = append(*trace, name+":before")

			next.ServeHTTP(w, r)

			*trace = append(*trace, name+":after")
		})
	}
}

// assertTrace verifies the middleware execution order.
func assertTrace(
	t *testing.T,
	got []string,
	want []string,
) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"unexpected execution trace:\n got: %v\nwant: %v",
			got,
			want,
		)
	}
}

// assertStatus verifies the HTTP response status.
func assertStatus(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	want int,
) {
	t.Helper()

	if recorder.Code != want {
		t.Errorf(
			"unexpected status code: got %d, want %d",
			recorder.Code,
			want,
		)
	}
}

// assertPanics verifies that the function panics.
func assertPanics(
	t *testing.T,
	fn func(),
) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Error("expected panic, but no panic occurred")
		}
	}()

	fn()
}

// testHandler returns a handler that responds with 204.
func testHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

// -----------------------------------------------------------------------------
// Router tests
// -----------------------------------------------------------------------------

func TestNew(t *testing.T) {
	r := New()

	if r == nil {
		t.Fatal("New() returned nil")
	}

	if r.Mux() == nil {
		t.Fatal("Router.Mux() returned nil")
	}

	recorder := performRequest(
		r,
		http.MethodGet,
		"/not-found",
	)

	assertStatus(t, recorder, http.StatusNotFound)
}

func TestRouterUseMiddlewareOrder(t *testing.T) {
	var trace []string

	r := New()

	r.Use(
		traceMiddleware("first", &trace),
		traceMiddleware("second", &trace),
	)

	r.HandleFunc("/users", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		trace = append(trace, "handler")
		w.WriteHeader(http.StatusOK)
	})

	recorder := performRequest(
		r,
		http.MethodGet,
		"/users",
	)

	assertStatus(t, recorder, http.StatusOK)

	assertTrace(t, trace, []string{
		"first:before",
		"second:before",
		"handler",
		"second:after",
		"first:after",
	})
}

func TestRouterUseMiddlewareOnNotFound(t *testing.T) {
	var trace []string

	r := New()

	r.Use(
		traceMiddleware("first", &trace),
		traceMiddleware("second", &trace),
	)

	recorder := performRequest(
		r,
		http.MethodGet,
		"/not-found",
	)

	assertStatus(t, recorder, http.StatusNotFound)

	// Global middleware wraps the entire ServeMux,
	// so it also runs for unmatched routes.
	assertTrace(t, trace, []string{
		"first:before",
		"second:before",
		"second:after",
		"first:after",
	})
}

func TestRouterHandleAndHandleFunc(t *testing.T) {
	r := New()

	// Global middleware should also apply to routes registered
	// through Router.Handle, Router.HandleFunc and Router.Mux.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			req *http.Request,
		) {
			w.Header().Set("X-Global-Middleware", "applied")
			next.ServeHTTP(w, req)
		})
	})

	r.Handle(
		"GET /handle",
		http.HandlerFunc(func(
			w http.ResponseWriter,
			req *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	r.HandleFunc(
		"POST /handle-func",
		func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusCreated)
		},
	)

	// Infrastructure route registered directly on ServeMux.
	r.Mux().HandleFunc(
		"GET /health",
		func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "Router.Handle",
			method:     http.MethodGet,
			path:       "/handle",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Router.HandleFunc",
			method:     http.MethodPost,
			path:       "/handle-func",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Router.Mux",
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := performRequest(
				r,
				tt.method,
				tt.path,
			)

			assertStatus(t, recorder, tt.wantStatus)

			if got := recorder.Header().Get("X-Global-Middleware"); got != "applied" {
				t.Errorf(
					"global middleware was not applied: got %q",
					got,
				)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Group prefix tests
// -----------------------------------------------------------------------------

func TestGroupPrefixNormalization(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		path   string
		target string
	}{
		{
			name:   "empty prefix",
			prefix: "",
			path:   "/users",
			target: "/users",
		},
		{
			name:   "prefix without leading slash",
			prefix: "api",
			path:   "users",
			target: "/api/users",
		},
		{
			name:   "prefix with trailing slash",
			prefix: "/api/",
			path:   "/users",
			target: "/api/users",
		},
		{
			name:   "prefix with leading and trailing whitespace",
			prefix: " /api/v1/ ",
			path:   "/users",
			target: "/api/v1/users",
		},
		{
			name:   "route without leading slash",
			prefix: "/api",
			path:   "users",
			target: "/api/users",
		},
		{
			name:   "root prefix",
			prefix: "/",
			path:   "/users",
			target: "/users",
		},
		{
			name:   "root route under prefix",
			prefix: "/api",
			path:   "/",
			target: "/api/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()

			g := r.Group(tt.prefix)
			g.GET(tt.path, testHandler())

			recorder := performRequest(
				r,
				http.MethodGet,
				tt.target,
			)

			assertStatus(t, recorder, http.StatusNoContent)
		})
	}
}

func TestGroupNestedPrefix(t *testing.T) {
	r := New()

	api := r.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")

	users.GET("/{userID}", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		userID := req.PathValue("userID")

		_, _ = w.Write([]byte(userID))
	}))

	recorder := performRequest(
		r,
		http.MethodGet,
		"/api/v1/users/user-123",
	)

	assertStatus(t, recorder, http.StatusOK)

	if got, want := recorder.Body.String(), "user-123"; got != want {
		t.Errorf("unexpected response body: got %q, want %q", got, want)
	}
}

func TestGroupRootRoute(t *testing.T) {
	r := New()

	api := r.Group("/api")

	api.GET("/", testHandler())

	// ServeMux treats a trailing slash pattern as a subtree pattern.
	tests := []struct {
		path       string
		wantStatus int
	}{
		{
			path:       "/api/",
			wantStatus: http.StatusNoContent,
		},
		{
			path:       "/api/health",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			recorder := performRequest(
				r,
				http.MethodGet,
				tt.path,
			)

			assertStatus(t, recorder, tt.wantStatus)
		})
	}
}

// -----------------------------------------------------------------------------
// Group middleware tests
// -----------------------------------------------------------------------------

func TestGroupMiddlewareOrder(t *testing.T) {
	var trace []string

	r := New()

	r.Use(
		traceMiddleware("global-1", &trace),
		traceMiddleware("global-2", &trace),
	)

	api := r.Group(
		"/api",
		traceMiddleware("api-1", &trace),
	)

	api.Use(
		traceMiddleware("api-2", &trace),
	)

	v1 := api.Group(
		"/v1",
		traceMiddleware("v1", &trace),
	)

	v1.GET("/users/{userID}", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		trace = append(trace, "handler")
		w.WriteHeader(http.StatusOK)
	}))

	recorder := performRequest(
		r,
		http.MethodGet,
		"/api/v1/users/123",
	)

	assertStatus(t, recorder, http.StatusOK)

	assertTrace(t, trace, []string{
		"global-1:before",
		"global-2:before",
		"api-1:before",
		"api-2:before",
		"v1:before",
		"handler",
		"v1:after",
		"api-2:after",
		"api-1:after",
		"global-2:after",
		"global-1:after",
	})
}

func TestGroupMiddlewareInheritanceSnapshot(t *testing.T) {
	var trace []string

	r := New()

	parent := r.Group(
		"/api",
		traceMiddleware("inherited", &trace),
	)

	// Register a route before adding another middleware.
	parent.GET("/before", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		trace = append(trace, "handler")
		w.WriteHeader(http.StatusOK)
	}))

	// Child inherits the middleware registered at creation time.
	child := parent.Group("/child")

	// Middleware added after child creation affects the parent,
	// but does not affect the already-created child.
	parent.Use(
		traceMiddleware("parent-late", &trace),
	)

	parent.GET("/after", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		trace = append(trace, "handler")
		w.WriteHeader(http.StatusOK)
	}))

	child.GET("/item", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		trace = append(trace, "handler")
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("route registered before Use", func(t *testing.T) {
		trace = nil

		recorder := performRequest(
			r,
			http.MethodGet,
			"/api/before",
		)

		assertStatus(t, recorder, http.StatusOK)

		assertTrace(t, trace, []string{
			"inherited:before",
			"handler",
			"inherited:after",
		})
	})

	t.Run("parent route registered after Use", func(t *testing.T) {
		trace = nil

		recorder := performRequest(
			r,
			http.MethodGet,
			"/api/after",
		)

		assertStatus(t, recorder, http.StatusOK)

		assertTrace(t, trace, []string{
			"inherited:before",
			"parent-late:before",
			"handler",
			"parent-late:after",
			"inherited:after",
		})
	})

	t.Run("child does not inherit parent middleware added later", func(t *testing.T) {
		trace = nil

		recorder := performRequest(
			r,
			http.MethodGet,
			"/api/child/item",
		)

		assertStatus(t, recorder, http.StatusOK)

		assertTrace(t, trace, []string{
			"inherited:before",
			"handler",
			"inherited:after",
		})
	})
}

func TestGroupUseOnlyAffectsRoutesRegisteredAfterward(t *testing.T) {
	var trace []string

	r := New()

	api := r.Group("/api")

	api.GET("/before", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		trace = append(trace, "handler")
		w.WriteHeader(http.StatusOK)
	}))

	api.Use(
		traceMiddleware("auth", &trace),
	)

	api.GET("/after", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		trace = append(trace, "handler")
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("before Use", func(t *testing.T) {
		trace = nil

		recorder := performRequest(
			r,
			http.MethodGet,
			"/api/before",
		)

		assertStatus(t, recorder, http.StatusOK)
		assertTrace(t, trace, []string{
			"handler",
		})
	})

	t.Run("after Use", func(t *testing.T) {
		trace = nil

		recorder := performRequest(
			r,
			http.MethodGet,
			"/api/after",
		)

		assertStatus(t, recorder, http.StatusOK)
		assertTrace(t, trace, []string{
			"auth:before",
			"handler",
			"auth:after",
		})
	})
}

func TestGroupMiddlewareIsIsolated(t *testing.T) {
	var trace []string

	r := New()

	public := r.Group("/public")

	private := r.Group(
		"/private",
		traceMiddleware("auth", &trace),
	)

	public.GET("/users", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		trace = append(trace, "public-handler")
		w.WriteHeader(http.StatusOK)
	}))

	private.GET("/users", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		trace = append(trace, "private-handler")
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("public group", func(t *testing.T) {
		trace = nil

		recorder := performRequest(
			r,
			http.MethodGet,
			"/public/users",
		)

		assertStatus(t, recorder, http.StatusOK)

		assertTrace(t, trace, []string{
			"public-handler",
		})
	})

	t.Run("private group", func(t *testing.T) {
		trace = nil

		recorder := performRequest(
			r,
			http.MethodGet,
			"/private/users",
		)

		assertStatus(t, recorder, http.StatusOK)

		assertTrace(t, trace, []string{
			"auth:before",
			"private-handler",
			"auth:after",
		})
	})

	t.Run("unmatched private route", func(t *testing.T) {
		trace = nil

		recorder := performRequest(
			r,
			http.MethodGet,
			"/private/not-found",
		)

		assertStatus(t, recorder, http.StatusNotFound)

		// Group middleware wraps registered handlers,
		// not the entire ServeMux.
		assertTrace(t, trace, nil)
	})
}

// -----------------------------------------------------------------------------
// HTTP method registration tests
// -----------------------------------------------------------------------------

func TestGroupHTTPMethodHelpers(t *testing.T) {
	r := New()
	api := r.Group("/api")

	handler := func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		w.Header().Set("X-Handler-Method", req.Method)
		w.WriteHeader(http.StatusNoContent)
	}

	tests := []struct {
		name     string
		method   string
		path     string
		register func(*Group, http.HandlerFunc)
	}{
		{
			name:   "GET",
			method: http.MethodGet,
			path:   "/get",
			register: func(g *Group, h http.HandlerFunc) {
				g.GET("/get", h)
			},
		},
		{
			name:   "HEAD",
			method: http.MethodHead,
			path:   "/head",
			register: func(g *Group, h http.HandlerFunc) {
				g.HEAD("/head", h)
			},
		},
		{
			name:   "POST",
			method: http.MethodPost,
			path:   "/post",
			register: func(g *Group, h http.HandlerFunc) {
				g.POST("/post", h)
			},
		},
		{
			name:   "PUT",
			method: http.MethodPut,
			path:   "/put",
			register: func(g *Group, h http.HandlerFunc) {
				g.PUT("/put", h)
			},
		},
		{
			name:   "PATCH",
			method: http.MethodPatch,
			path:   "/patch",
			register: func(g *Group, h http.HandlerFunc) {
				g.PATCH("/patch", h)
			},
		},
		{
			name:   "DELETE",
			method: http.MethodDelete,
			path:   "/delete",
			register: func(g *Group, h http.HandlerFunc) {
				g.DELETE("/delete", h)
			},
		},
		{
			name:   "OPTIONS",
			method: http.MethodOptions,
			path:   "/options",
			register: func(g *Group, h http.HandlerFunc) {
				g.OPTIONS("/options", h)
			},
		},
	}

	for _, tt := range tests {
		tt.register(api, handler)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := performRequest(
				r,
				tt.method,
				"/api"+tt.path,
			)

			assertStatus(t, recorder, http.StatusNoContent)

			if got := recorder.Header().Get("X-Handler-Method"); got != tt.method {
				t.Errorf(
					"unexpected handler method: got %q, want %q",
					got,
					tt.method,
				)
			}
		})
	}
}

func TestGroupHandleVariants(t *testing.T) {
	r := New()
	api := r.Group("/api")

	var handler http.HandlerFunc = func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		w.WriteHeader(http.StatusNoContent)
	}

	// ServeMux method pattern through Handle.
	api.Handle("GET /by-handle", handler)

	// ServeMux method pattern through HandleFunc.
	api.HandleFunc("POST /by-handle-func", handler)

	// Method registration through HandleMethod.
	api.HandleMethod("PUT", "/by-method", handler)

	// HandleMethod should trim surrounding whitespace.
	api.HandleMethod(" PATCH ", "/by-method-trim", handler)

	tests := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/api/by-handle",
		},
		{
			method: http.MethodPost,
			path:   "/api/by-handle-func",
		},
		{
			method: http.MethodPut,
			path:   "/api/by-method",
		},
		{
			method: http.MethodPatch,
			path:   "/api/by-method-trim",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			recorder := performRequest(
				r,
				tt.method,
				tt.path,
			)

			assertStatus(t, recorder, http.StatusNoContent)
		})
	}
}

func TestRouteMiddlewareOnRegistrationMethods(t *testing.T) {
	r := New()
	g := r.Group("/api")
	var trace []string
	middleware := traceMiddleware("route", &trace)
	for _, tt := range []struct {
		method string
		root   func(string, http.Handler, ...Middleware)
		group  func(string, http.Handler, ...Middleware)
	}{
		{http.MethodGet, r.GET, g.GET},
		{http.MethodHead, r.HEAD, g.HEAD},
		{http.MethodPost, r.POST, g.POST},
		{http.MethodDelete, r.DELETE, g.DELETE},
		{http.MethodPatch, r.PATCH, g.PATCH},
		{http.MethodPut, r.PUT, g.PUT},
		{http.MethodOptions, r.OPTIONS, g.OPTIONS},
	} {
		path := "/" + tt.method
		tt.root(path, testHandler(), middleware)
		tt.group(path, testHandler(), middleware)
		for _, target := range []string{path, "/api" + path} {
			trace = nil
			assertStatus(t, performRequest(r, tt.method, target), http.StatusNoContent)
			assertTrace(t, trace, []string{"route:before", "route:after"})
		}
	}

	r.Handle("POST /pattern", testHandler(), middleware)
	r.HandleFunc("POST /function", testHandler(), middleware)
	g.HandleFunc("POST /function", testHandler(), middleware)
	g.Group("/resource").HandleMethod(http.MethodPost, ":action", testHandler(), middleware)
	for _, path := range []string{"/pattern", "/function", "/api/function", "/api/resource:action"} {
		trace = nil
		assertStatus(t, performRequest(r, http.MethodPost, path), http.StatusNoContent)
		assertTrace(t, trace, []string{"route:before", "route:after"})
	}
}

func TestRouteMiddlewareOrderAndIsolation(t *testing.T) {
	var trace []string
	r := New()
	r.Use(traceMiddleware("global", &trace))
	g := r.Group("/api", traceMiddleware("group", &trace))
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		trace = append(trace, "handler")
		w.WriteHeader(http.StatusNoContent)
	})
	g.Handle("GET /specific", handler, traceMiddleware("first", &trace), traceMiddleware("second", &trace))
	g.Handle("GET /sibling", handler)
	g.Handle("GET /blocked", handler, func(http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			trace = append(trace, "blocked")
			w.WriteHeader(http.StatusForbidden)
		})
	})

	assertStatus(t, performRequest(r, http.MethodGet, "/api/specific"), http.StatusNoContent)
	assertTrace(t, trace, []string{
		"global:before", "group:before", "first:before", "second:before",
		"handler", "second:after", "first:after", "group:after", "global:after",
	})
	trace = nil
	assertStatus(t, performRequest(r, http.MethodGet, "/api/sibling"), http.StatusNoContent)
	assertTrace(t, trace, []string{"global:before", "group:before", "handler", "group:after", "global:after"})
	trace = nil
	assertStatus(t, performRequest(r, http.MethodGet, "/api/blocked"), http.StatusForbidden)
	assertTrace(t, trace, []string{"global:before", "group:before", "blocked", "group:after", "global:after"})
}

func TestAnyAndMatchRouteMiddleware(t *testing.T) {
	r := New()
	g := r.Group("/api")
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			w.Header().Set("X-Route-Middleware", request.PathValue("id"))
			next.ServeHTTP(w, request)
		})
	}
	methods := []string{http.MethodPut, http.MethodPatch, "REPORT"}
	r.Any("/any/{id}", testHandler(), middleware)
	g.Any("/any/{id}", testHandler(), middleware)
	r.Match(methods, "/match/{id}", testHandler(), middleware)
	g.Match(methods, "/match/{id}", testHandler(), middleware)

	for _, tt := range []struct {
		path       string
		allowed    []string
		disallowed string
	}{
		{"/any/record", []string{
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
			http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect, http.MethodTrace, "QUERY",
		}, "REPORT"},
		{"/match/record", methods, http.MethodGet},
	} {
		for _, path := range []string{tt.path, "/api" + tt.path} {
			for _, method := range tt.allowed {
				recorder := performRequest(r, method, path)
				assertStatus(t, recorder, http.StatusNoContent)
				if got := recorder.Header().Get("X-Route-Middleware"); got != "record" {
					t.Errorf("%s %s: X-Route-Middleware = %q, want record", method, path, got)
				}
			}
			assertStatus(t, performRequest(r, tt.disallowed, path), http.StatusMethodNotAllowed)
		}
	}
}

func TestRouteRegistrationRejectsInvalidArguments(t *testing.T) {
	r := New()
	g := r.Group("/api")
	assertPanics(t, func() { r.Handle("/nil-handler", nil) })
	assertPanics(t, func() { r.HandleFunc("/nil-handler-func", nil) })
	assertPanics(t, func() { g.Handle("/nil-middleware", testHandler(), nil) })
	assertPanics(t, func() { r.Any("/nil-middleware", testHandler(), nil) })
	for _, method := range []string{" ", "GET /invalid"} {
		assertPanics(t, func() { g.Match([]string{http.MethodGet, method}, "/invalid", testHandler()) })
		assertStatus(t, performRequest(r, http.MethodGet, "/api/invalid"), http.StatusNotFound)
	}
	g.Match(nil, "/empty", testHandler())
	assertStatus(t, performRequest(r, http.MethodGet, "/api/empty"), http.StatusNotFound)
}

func TestGroupMethodNotAllowed(t *testing.T) {
	r := New()

	api := r.Group("/api")
	api.GET("/users", testHandler())

	recorder := performRequest(
		r,
		http.MethodPost,
		"/api/users",
	)

	assertStatus(t, recorder, http.StatusMethodNotAllowed)
}

// -----------------------------------------------------------------------------
// ServeMux integration tests
// -----------------------------------------------------------------------------

func TestGroupPathValue(t *testing.T) {
	r := New()

	api := r.Group("/api/v1")

	api.GET("/users/{userID}/posts/{postID}", http.HandlerFunc(func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		userID := req.PathValue("userID")
		postID := req.PathValue("postID")

		w.Header().Set("X-User-ID", userID)
		w.Header().Set("X-Post-ID", postID)

		w.WriteHeader(http.StatusOK)
	}))

	recorder := performRequest(
		r,
		http.MethodGet,
		"/api/v1/users/u123/posts/p456",
	)

	assertStatus(t, recorder, http.StatusOK)

	if got := recorder.Header().Get("X-User-ID"); got != "u123" {
		t.Errorf("unexpected userID: got %q, want %q", got, "u123")
	}

	if got := recorder.Header().Get("X-Post-ID"); got != "p456" {
		t.Errorf("unexpected postID: got %q, want %q", got, "p456")
	}
}

func TestGroupHostPattern(t *testing.T) {
	r := New()

	api := r.Group("/api")

	api.Handle(
		"GET example.com/health",
		testHandler(),
	)

	t.Run("matching host", func(t *testing.T) {
		recorder := performRequest(
			r,
			http.MethodGet,
			"http://example.com/api/health",
		)

		assertStatus(t, recorder, http.StatusNoContent)
	})

	t.Run("non-matching host", func(t *testing.T) {
		recorder := performRequest(
			r,
			http.MethodGet,
			"http://other.example/api/health",
		)

		assertStatus(t, recorder, http.StatusNotFound)
	})
}

// -----------------------------------------------------------------------------
// Validation and panic tests
// -----------------------------------------------------------------------------

func TestRouterRejectsNilMiddleware(t *testing.T) {
	r := New()

	assertPanics(t, func() {
		r.Use(nil)
	})

	assertPanics(t, func() {
		r.Group("/api", nil)
	})
}

func TestGroupRejectsNilMiddleware(t *testing.T) {
	r := New()
	g := r.Group("/api")

	assertPanics(t, func() {
		g.Use(nil)
	})

	assertPanics(t, func() {
		g.Group("/v1", nil)
	})
}

func TestGroupRejectsNilHandler(t *testing.T) {
	r := New()
	g := r.Group("/api")

	assertPanics(t, func() {
		g.Handle("/nil", nil)
	})

	assertPanics(t, func() {
		g.HandleFunc("/nil", nil)
	})

	assertPanics(t, func() {
		g.HandleMethod(
			"GET",
			"/nil",
			nil,
		)
	})
}

func TestGroupHandleMethodRejectsEmptyMethod(t *testing.T) {
	r := New()
	g := r.Group("/api")

	assertPanics(t, func() {
		g.HandleMethod(
			"",
			"/users",
			testHandler(),
		)
	})

	assertPanics(t, func() {
		g.HandleMethod(
			"   ",
			"/users",
			testHandler(),
		)
	})
}

func TestGroupHandleRejectsInvalidPattern(t *testing.T) {
	r := New()
	g := r.Group("/api")

	assertPanics(t, func() {
		g.Handle(
			"",
			testHandler(),
		)
	})

	assertPanics(t, func() {
		g.Handle(
			"GET /users extra",
			testHandler(),
		)
	})

	// Invalid ServeMux wildcard pattern.
	assertPanics(t, func() {
		g.Handle(
			"GET /users/{userID",
			testHandler(),
		)
	})
}

func TestGroupHandleRejectsDuplicatePattern(t *testing.T) {
	r := New()
	g := r.Group("/api")

	g.GET("/users", testHandler())

	assertPanics(t, func() {
		g.GET("/users", testHandler())
	})
}
