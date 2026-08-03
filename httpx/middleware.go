package httpx

import (
	"net/http"
)

// ServerMiddleware defines a function type for HTTP middleware.
// It receives an http.ResponseWriter, an http.Request, and the next handler (invoker) in the chain.
//
// Parameters:
//
//	response - http.ResponseWriter to write the HTTP response
//	request  - *http.Request containing the HTTP request data
//	invoker  - http.HandlerFunc representing the next handler in the middleware chain
type ServerMiddleware func(response http.ResponseWriter, request *http.Request, invoker http.HandlerFunc)

// ClientInvoker is a function type that defines how to invoke an HTTP request.
// It takes an HTTP client, and an HTTP request, and returns an HTTP response or an error.
//
// Parameters:
//   - cli: The HTTP client to use for the request
//   - request: The HTTP request to invoke
//
// Returns:
//   - *http.Response: The HTTP response from the request
//   - error: Any error that occurred during the request, or nil if successful
type ClientInvoker func(cli *http.Client, request *http.Request) (*http.Response, error)

// ClientMiddleware is a function type that defines middleware for HTTP requests.
// It takes an HTTP client, an HTTP request, and the next invoker in the chain,
// and returns an HTTP response or an error.
//
// Parameters:
//   - cli: The HTTP client to use for the request
//   - request: The HTTP request to process
//   - invoker: The next invoker in the middleware chain
//
// Returns:
//   - *http.Response: The HTTP response from the request
//   - error: Any error that occurred during the request, or nil if successful
type ClientMiddleware func(cli *http.Client, request *http.Request, invoker ClientInvoker) (*http.Response, error)

// ServerChain combines multiple Middleware functions into a single Middleware.
// If no middlewares are provided, returns nil.
// If one middleware is provided, returns that middleware.
// If multiple middlewares are provided, returns a middleware that executes them in chain.
//
// Parameters:
//
//	middlewares - variadic list of Middleware functions to chain together
//
// Returns:
//
//	Middleware - a single middleware function representing the entire chain
func ServerChain(middlewares ...ServerMiddleware) ServerMiddleware {
	var mdw ServerMiddleware
	if len(middlewares) == 0 {
		mdw = nil
	} else if len(middlewares) == 1 {
		mdw = middlewares[0]
	} else {
		mdw = func(response http.ResponseWriter, request *http.Request, invoker http.HandlerFunc) {
			middlewares[0](response, request, getServerInvoker(middlewares, 0, invoker))
		}
	}
	return mdw
}

// getServerInvoker recursively builds the invoker chain for executing middlewares in sequence.
//
// Parameters:
//
//	interceptors   - slice of Middleware functions to be executed
//	curr           - current index in the interceptors slice
//	finalInvoker   - the final handler function to be executed after all middlewares
//
// Returns:
//
//	http.HandlerFunc - a handler function that continues the middleware chain
func getServerInvoker(interceptors []ServerMiddleware, curr int, finalInvoker http.HandlerFunc) http.HandlerFunc {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(response http.ResponseWriter, request *http.Request) {
		interceptors[curr+1](response, request, getServerInvoker(interceptors, curr+1, finalInvoker))
	}
}

// ServerInvoke wraps a Middleware and a final handler function into an http.Handler.
// If the middleware is nil, directly calls the final handler.
// Otherwise, executes the middleware chain ending with the final handler.
//
// Parameters:
//
//	middleware - Middleware function to execute (can be nil)
//	invoke - http.HandlerFunc representing the final handler
//	response - http.ResponseWriter to write the HTTP response
//	request - *http.Request representing the incoming HTTP request
//	routeInfo - *goose.RouteInfo representing the route information
func ServerInvoke(middleware ServerMiddleware, response http.ResponseWriter, request *http.Request, invoke http.HandlerFunc) {
	if middleware == nil {
		invoke(response, request)
		return
	}
	middleware(response, request, invoke)
}

// ClientChain combines multiple middleware functions into a single middleware function.
// It creates a chain where each middleware calls the next one in the sequence.
//
// Parameters:
//   - middlewares: A variadic list of middleware functions to chain together
//
// Returns:
//   - Middleware: A single middleware function that represents the entire chain
func ClientChain(middlewares ...ClientMiddleware) ClientMiddleware {
	var mdw ClientMiddleware
	if len(middlewares) == 0 {
		mdw = nil
	} else if len(middlewares) == 1 {
		mdw = middlewares[0]
	} else {
		mdw = func(cli *http.Client, request *http.Request, invoker ClientInvoker) (*http.Response, error) {
			return middlewares[0](cli, request, getClintInvoker(middlewares, 0, invoker))
		}
	}
	return mdw
}

// getClintInvoker is a helper function that creates an invoker chain from a list of middleware functions.
// It recursively builds invokers that call the next middleware in the chain.
//
// Parameters:
//   - interceptors: The list of middleware functions
//   - curr: The current index in the middleware list
//   - finalInvoker: The final invoker to call at the end of the chain
//
// Returns:
//   - Invoker: An invoker that calls the next middleware in the chain
func getClintInvoker(interceptors []ClientMiddleware, curr int, finalInvoker ClientInvoker) ClientInvoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(cli *http.Client, request *http.Request) (*http.Response, error) {
		return interceptors[curr+1](cli, request, getClintInvoker(interceptors, curr+1, finalInvoker))
	}
}

// ClientInvoke executes an HTTP request with the given middleware.
// If no middleware is provided, it directly executes the request using the HTTP client.
//
// Parameters:
//   - ctx: The context.Context for the request
//   - middleware: The middleware to apply to the request (can be nil)
//   - cli: The HTTP client to use for the request
//   - request: The HTTP request to execute
//
// Returns:
//   - *http.Response: The HTTP response from the request
//   - error: Any error that occurred during the request, or nil if successful
func ClientInvoke(middleware ClientMiddleware, cli *http.Client, request *http.Request) (*http.Response, error) {
	if middleware == nil {
		return clientInvoke(cli, request)
	}
	return middleware(cli, request, clientInvoke)
}

func clientInvoke(cli *http.Client, request *http.Request) (*http.Response, error) {
	return cli.Do(request)
}
