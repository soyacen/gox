package pagex

import (
	"context"
)

type key struct{}

// NewContext creates a new context with the given pagination parameters.
//
// Parameters:
//   - ctx: the parent context.
//   - pageNum: page number, starting from 1.
//   - pageSize: number of items per page.
//   - opts: optional pagination options.
//
// Returns:
//   - context.Context: the new context containing the page.
//   - error: error if page creation fails.
func NewContext(ctx context.Context, pageNum uint64, pageSize uint64, opts ...Option) (context.Context, error) {
	page, err := NewPage(pageNum, pageSize, opts...)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, key{}, page), nil
}

// FromContext retrieves the Page from the context.
//
// Parameters:
//   - ctx: the context.
//
// Returns:
//   - *Page: the page stored in the context.
//   - bool: true if a page was found.
func FromContext(ctx context.Context) (*Page, bool) {
	v, ok := ctx.Value(key{}).(*Page)
	return v, ok
}

// WithoutPage returns a context with the page removed.
//
// Parameters:
//   - ctx: the context.
//
// Returns:
//   - context.Context: the context with the page removed.
func WithoutPage(ctx context.Context) context.Context {
	return context.WithValue(ctx, key{}, nil)
}
