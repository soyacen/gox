package outgoing

import (
	"context"
	"net/http"
)

type Invoker func(ctx context.Context, req *http.Request, cli *http.Client) (*http.Response, error)

type Middleware func(ctx context.Context, req *http.Request, cli *http.Client, invoker Invoker) (*http.Response, error)

func chainMiddlewares(middlewares ...Middleware) Middleware {
	var chainedInt Middleware
	if len(middlewares) == 0 {
		chainedInt = nil
	} else if len(middlewares) == 1 {
		chainedInt = middlewares[0]
	} else {
		chainedInt = func(ctx context.Context, req *http.Request, cli *http.Client, invoker Invoker) (*http.Response, error) {
			return middlewares[0](ctx, req, cli, getInvoker(middlewares, 0, invoker))
		}
	}
	return chainedInt
}

func getInvoker(interceptors []Middleware, curr int, finalInvoker Invoker) Invoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(ctx context.Context, req *http.Request, cli *http.Client) (*http.Response, error) {
		return interceptors[curr+1](ctx, req, cli, getInvoker(interceptors, curr+1, finalInvoker))
	}
}

func Invoke(ctx context.Context, middleware Middleware, cli *http.Client, request *http.Request) (*http.Response, error) {
	if middleware == nil {
		return invoke(ctx, request, cli)
	}
	return middleware(ctx, request, cli, invoke)
}

func invoke(_ context.Context, req *http.Request, cli *http.Client) (*http.Response, error) {
	return cli.Do(req)
}
