// The MIT License (MIT)

// Copyright (c) 2019-present Rafał Lorenz

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package waiter

import (
	"context"
	"errors"
	"sync"
)

// ErrNoCallbacks is returned when no callbacks are provided to Race, All, or Retry.
var ErrNoCallbacks = errors.New("waiter: no callback to run")

// AsyncFunc represents an asynchronous function that can be executed by Race, All, or Retry.
// It receives a context and returns a result and an error.
//
// Parameters:
//   - ctx: the context for cancellation and deadlines
//
// Returns:
//   - interface{}: the result of the function execution
//   - error: any error that occurred during execution
type AsyncFunc func(ctx context.Context) (interface{}, error)

type response struct {
	res interface{}
	err error
}

// Race executes multiple asynchronous functions concurrently and returns the result
// from the first one that completes without an error.
// If all callbacks fail, the last error is returned.
// Will panic if context is nil.
//
// Parameters:
//   - ctx: the context for cancellation and deadlines
//   - fns: one or more asynchronous functions to execute
//
// Returns:
//   - interface{}: the result from the first successful callback
//   - error: any error that occurred, or ErrNoCallbacks if no callbacks are provided
func Race(ctx context.Context, fns ...AsyncFunc) (interface{}, error) {
	if ctx == nil {
		panic("nil context provided")
	}
	if len(fns) == 0 {
		return nil, ErrNoCallbacks
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	out := make(chan *response, 1)
	defer close(out)
	responses := make(chan *response, len(fns))
	defer close(responses)

	var responded bool
	var lock sync.Mutex
	var errCount int
	go func() {
		for {
			select {
			case <-ctx.Done():
				lock.Lock()
				if !responded {
					responded = true
					out <- &response{
						err: ctx.Err(),
					}
					lock.Unlock()
					return
				}
				lock.Unlock()
			case r, more := <-responses:
				lock.Lock()
				if !more || (!responded && r.err == nil) || (errCount == len(fns)) {
					responded = true
					out <- r
					lock.Unlock()
					return
				}
				lock.Unlock()
			}
		}
	}()

	for _, fn := range fns {
		go func(f AsyncFunc) {
			var r response
			r.res, r.err = f(ctx)

			lock.Lock()
			defer lock.Unlock()
			if r.err != nil {
				errCount++
			}
			if !responded {
				responses <- &r
			}
		}(fn)
	}

	r := <-out

	return r.res, r.err
}

// All executes all asynchronous functions concurrently and waits for all of them to complete.
// The returned responses and errors are ordered according to the callback order.
// Will panic if context is nil.
//
// Parameters:
//   - ctx: the context for cancellation and deadlines
//   - fns: one or more asynchronous functions to execute
//
// Returns:
//   - []interface{}: the results from all callbacks, ordered by callback index
//   - []error: the errors from all callbacks, ordered by callback index
func All(ctx context.Context, fns ...AsyncFunc) ([]interface{}, []error) {
	if ctx == nil {
		panic("nil context provided")
	}

	rs := make([]interface{}, len(fns))
	errs := make([]error, len(fns))

	var wg sync.WaitGroup
	wg.Add(len(fns))

	for i, fn := range fns {
		go func(index int, f AsyncFunc) {
			defer wg.Done()

			var r response
			r.res, r.err = f(ctx)

			rs[index] = r.res
			errs[index] = r.err
		}(i, fn)
	}

	wg.Wait()

	return rs, errs
}

// Retry retries the given asynchronous function until it executes without an error
// or the maximum number of retries is reached.
// When retries = 0, it will retry infinitely.
// Will panic if context is nil.
//
// Parameters:
//   - ctx: the context for cancellation and deadlines
//   - retries: the maximum number of retries, or 0 for infinite retries
//   - fn: the asynchronous function to retry
//
// Returns:
//   - interface{}: the result from the successful execution
//   - error: any error that occurred after all retries are exhausted
func Retry(ctx context.Context, retires int, fn AsyncFunc) (interface{}, error) {
	if ctx == nil {
		panic("nil context provided")
	}

	i := 1

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			var r response
			r.res, r.err = fn(ctx)

			if r.err == nil || i == retires {
				return r.res, r.err
			}

			i++
		}
	}
}
