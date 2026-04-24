package httpx

import (
	"github.com/soyacen/gox/slicex"
	"net/http"
)

// Default405Body is the default response body for 405 Method Not Allowed errors.
var Default405Body = []byte("405 method not allowed")

// GetHandler wraps an HTTP handler to only allow GET requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows GET requests
func GetHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// GetHandlerFunc wraps an HTTP handler function to only allow GET requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows GET requests
func GetHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// HeadHandler wraps an HTTP handler to only allow HEAD requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows HEAD requests
func HeadHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodHead {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// HeadHandlerFunc wraps an HTTP handler function to only allow HEAD requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows HEAD requests
func HeadHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodHead {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// PostHandler wraps an HTTP handler to only allow POST requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows POST requests
func PostHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// PostHandlerFunc wraps an HTTP handler function to only allow POST requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows POST requests
func PostHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// PutHandler wraps an HTTP handler to only allow PUT requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows PUT requests
func PutHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// PutHandlerFunc wraps an HTTP handler function to only allow PUT requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows PUT requests
func PutHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// PatchHandler wraps an HTTP handler to only allow PATCH requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows PATCH requests
func PatchHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPatch {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// PatchHandlerFunc wraps an HTTP handler function to only allow PATCH requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows PATCH requests
func PatchHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPatch {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// DeleteHandler wraps an HTTP handler to only allow DELETE requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows DELETE requests
func DeleteHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// DeleteHandlerFunc wraps an HTTP handler function to only allow DELETE requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows DELETE requests
func DeleteHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// ConnectHandler wraps an HTTP handler to only allow CONNECT requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows CONNECT requests
func ConnectHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodConnect {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// ConnectHandlerFunc wraps an HTTP handler function to only allow CONNECT requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows CONNECT requests
func ConnectHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodConnect {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// OptionsHandler wraps an HTTP handler to only allow OPTIONS requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows OPTIONS requests
func OptionsHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodOptions {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// OptionsHandlerFunc wraps an HTTP handler function to only allow OPTIONS requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows OPTIONS requests
func OptionsHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodOptions {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// TraceHandler wraps an HTTP handler to only allow TRACE requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows TRACE requests
func TraceHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodTrace {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// TraceHandlerFunc wraps an HTTP handler function to only allow TRACE requests.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows TRACE requests
func TraceHandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodTrace {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// Handler wraps an HTTP handler to only allow requests with the specified method.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - method: The allowed HTTP method
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows the specified method
func Handler(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// HandlerFunc wraps an HTTP handler function to only allow requests with the specified method.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - method: The allowed HTTP method
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows the specified method
func HandlerFunc(method string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}

// MatchHandler wraps an HTTP handler to only allow requests with one of the specified methods.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - methods: The allowed HTTP methods
//   - handler: The HTTP handler to wrap
//
// Returns:
//   - http.Handler: A handler that only allows the specified methods
func MatchHandler(methods []string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if slicex.NotContains(methods, req.Method) {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler.ServeHTTP(resp, req)
	})
}

// MatchHandlerFunc wraps an HTTP handler function to only allow requests with one of the specified methods.
// Returns 405 Method Not Allowed for other HTTP methods.
//
// Parameters:
//   - methods: The allowed HTTP methods
//   - handler: The HTTP handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A handler function that only allows the specified methods
func MatchHandlerFunc(methods []string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if slicex.NotContains(methods, req.Method) {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = resp.Write(Default405Body)
		}
		handler(resp, req)
	}
}
