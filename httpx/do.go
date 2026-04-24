package httpx

import (
	"context"
	"errors"
	"fmt"
)

// DoCommand is a command that executes an HTTP request using the client and
// request stored in the context.
type DoCommand struct{}

// Execute sends the HTTP request using the client from the context.
// It expects both an *http.Client and an *http.Request to be present in the
// context. The response is stored back into the context.
//
// Parameters:
//   - ctx: The context containing the HTTP client and request
//
// Returns:
//   - context.Context: A new context containing the HTTP response
//   - error: An error if the client or request is missing, or if the request fails
func (cmd *DoCommand) Execute(ctx context.Context) (context.Context, error) {
	cli, ok := ClientFromContext(ctx)
	if !ok {
		return ctx, errors.New("http client is nil")
	}
	req, ok := RequestFromContext(ctx)
	if !ok {
		return ctx, errors.New("http request is nil")
	}
	resp, err := cli.Do(req)
	if err != nil {
		return ctx, fmt.Errorf("failed to send http request, %w", err)
	}
	return NewContextWithResponse(ctx, resp), nil
}
