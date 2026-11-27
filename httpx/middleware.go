package httpx

import (
	"net/http"
)

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
type ClientInvoker = func(cli *http.Client, request *http.Request) (*http.Response, error)

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
type ClientMiddleware = func(cli *http.Client, request *http.Request, invoker ClientInvoker) (*http.Response, error)

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
