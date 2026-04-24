package chanx

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

var ErrDefaultBranch = fmt.Errorf("chanx: default branch")

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// All merges multiple input channels into a single output channel, then returns a channel
// containing a slice of all current values from the input channels.
// See: [Go Concurrency Patterns: Timing out, moving on](https://go.dev/blog/concurrency-timeouts)
//
// Parameters:
//   - ctx: context for cancellation
//   - ins: input channels to merge
//
// Returns:
//   - <-chan []T: channel containing a slice of all values
func All[T any](ctx context.Context, ins ...<-chan T) <-chan []T {
	return AsSlice(ctx, FanIn(ctx, ins...))
}

// Any selects any one value from multiple input channels and returns a channel
// containing that value.
//
// Parameters:
//   - ctx: context for cancellation
//   - ins: input channels to select from
//
// Returns:
//   - <-chan T: channel containing the selected value
func Any[T any](ctx context.Context, ins ...<-chan T) <-chan T {
	out := make(chan T, 1)
	go _any(ctx, out, ins...)
	return out
}

// _any 用于从多个输入通道中选择一个值并发送到 valueC 通道，或者在上下文取消时将错误发送到 errC 通道。
func _any[T any](ctx context.Context, out chan T, ins ...<-chan T) {
	// 确保在函数退出时关闭 out 和 errC 通道。
	defer close(out)
	for len(ins) > 0 {
		// 当 ins 列表不为空时，随机选择一个通道 ch。
		sel := r.Intn(len(ins))
		ch := ins[sel]
		select {
		case <-ctx.Done():
			return
		case v, ok := <-ch:
			if !ok {
				// ch被关闭，则从 ins 列表中删除 ch，并继续循环。
				ins = removeAt(ins, sel)
				continue
			}
			select {
			case out <- v:
				return
			case <-ctx.Done():
				return
			}
		default:
			// 防止死锁
			runtime.Gosched()
		}
	}
}

func removeFunc[S ~[]E, E any](s S, f func(i int, v E) bool) S {
	var index int
	return slices.DeleteFunc(s, func(e E) bool {
		del := f(index, e)
		index++
		return del
	})
}

func removeAt[S ~[]E, E any](array S, is ...int) S {
	return removeFunc(array, func(i int, v E) bool {
		return slices.Contains(is, i)
	})
}

// TrySend attempts to send a value to the specified channel.
//
// Parameters:
//   - in: channel to send to
//   - v: value to send
//
// Returns:
//   - error: nil on success, ErrDefaultBranch if channel is full
func TrySend[T any](in chan<- T, v T) error {
	select {
	case in <- v:
		return nil
	default:
		return ErrDefaultBranch
	}
}

// TryReceive attempts to receive data from a channel.
//
// Parameters:
//   - in: channel to receive from
//
// Returns:
//   - T: received value
//   - error: nil on success, ErrDefaultBranch if no value available
func TryReceive[T any](in <-chan T) (T, error) {
	var v T
	select {
	case v = <-in:
		return v, nil
	default:
		return v, ErrDefaultBranch
	}
}

// Emit creates a channel that asynchronously sends the provided values,
// allowing cancellation via context.
//
// Parameters:
//   - ctx: context for cancellation
//   - values: values to emit
//
// Returns:
//   - <-chan T: channel receiving the emitted values
func Emit[T any](ctx context.Context, values ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, value := range values {
			select {
			case <-ctx.Done():
				return
			case out <- value:
			}
		}
	}()
	return out
}

// Once creates a channel containing a single value.
//
// Parameters:
//   - m: value to send
//
// Returns:
//   - <-chan M: channel with a single value
func Once[M any](m M) <-chan M {
	ch := make(chan M, 1)
	ch <- m
	close(ch)
	return ch
}

// Nothing creates and immediately closes a channel.
//
// Returns:
//   - <-chan M: closed channel
func Nothing[M any]() <-chan M {
	ch := make(chan M)
	close(ch)
	return ch
}

// Append appends channels to a slice of channels.
//
// Parameters:
//   - c: base slice
//   - channels: channels to append
//
// Returns:
//   - []chan T: appended slice
func Append[T any](c []chan T, channels ...chan T) []chan T {
	return Append(c, channels...)
}

// AppendSendChannel appends send channels to a slice of send channels.
//
// Parameters:
//   - c: base slice
//   - channels: channels to append
//
// Returns:
//   - []<-chan T: appended slice
func AppendSendChannel[T any](c []<-chan T, channels ...chan T) []<-chan T {
	for _, ch := range channels {
		c = append(c, ch)
	}
	return c
}

// AppendReceiveChannel appends receive channels to a slice of receive channels.
//
// Parameters:
//   - c: base slice
//   - channels: channels to append
//
// Returns:
//   - []chan<- T: appended slice
func AppendReceiveChannel[T any](c []chan<- T, channels ...chan T) []chan<- T {
	for _, ch := range channels {
		c = append(c, ch)
	}
	return c
}

// Copy copies values from src to dest and closes dest when done.
//
// Parameters:
//   - dest: destination channel
//   - src: source channel
func Copy[T any](dest chan<- T, src <-chan T) {
	for v := range src {
		dest <- v
	}
	close(dest)
}

// AsyncCopy copies values from src to dest asynchronously and closes dest when done.
//
// Parameters:
//   - dest: destination channel
//   - src: source channel
func AsyncCopy[T any](dest chan<- T, src <-chan T) {
	go func() { Copy[T](dest, src) }()
}

// Pipe copies values from src to dest.
// Deprecated: Do not use. Use Copy instead.
//
// Parameters:
//   - src: source channel
//   - dest: destination channel
func Pipe[T any](src <-chan T, dest chan<- T) {
	Copy[T](dest, src)
}

// AsyncPipe copies values from src to dest asynchronously.
// Deprecated: Do not use. Use AsyncCopy instead.
//
// Parameters:
//   - src: source channel
//   - dest: destination channel
func AsyncPipe[T any](src <-chan T, dest chan<- T) {
	AsyncCopy[T](dest, src)
}

// Skip skips the first n elements from the input channel and passes remaining elements to the output channel.
//
// Parameters:
//   - ctx: context for cancellation
//   - in: input channel
//   - n: number of elements to skip
//
// Returns:
//   - <-chan T: output channel with remaining elements
func Skip[T any](ctx context.Context, in <-chan T, n int) <-chan T {
	var out chan T
	if in == nil {
		return out
	}
	out = make(chan T, cap(in))
	go func() {
		defer close(out)
		for {
			var value T
			var ok bool
			select {
			case <-ctx.Done():
				return
			case value, ok = <-in:
				if !ok {
					return
				}
				n--
			}
			if n >= 0 {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- value:
			}
		}
	}()
	return out
}

// Pipeline creates a data processing pipeline.
// It receives an input channel in and a processing function f, returning an output channel out.
// In a new goroutine, it iterates over each element in the input channel, applies the processing function,
// and sends the result to the output channel.
// If context ctx is cancelled, it exits early and closes the output channel.
//
// Parameters:
//   - ctx: context for cancellation
//   - in: input channel
//   - f: processing function
//
// Returns:
//   - <-chan R: output channel with processed elements
func Pipeline[T any, R any](ctx context.Context, in <-chan T, f func(T) R) <-chan R {
	out := make(chan R, len(in))
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- f(v):
				}
			}
		}
	}()
	return out
}

// Reduce performs a reduction operation on elements from the input channel,
// using the provided accumulator function and initial value, returning the final result.
//
// Parameters:
//   - ctx: context for cancellation
//   - in: input channel
//   - identity: initial value
//   - accumulator: reduction function
//
// Returns:
//   - <-chan R: channel containing the reduced result
func Reduce[T any, R any](ctx context.Context, in <-chan T, identity R, accumulator func(R, T) R) <-chan R {
	return Pipeline[T, R](ctx, in, func(value T) R {
		identity = accumulator(identity, value)
		return identity
	})
}

// Max finds the maximum value from the input channel and returns it through the output channel.
//
// Parameters:
//   - ctx: context for cancellation
//   - in: input channel
//   - cmp: comparison function
//
// Returns:
//   - <-chan T: channel containing the maximum value
func Max[T any](ctx context.Context, in <-chan T, cmp func(a, b T) int) <-chan T {
	out := make(chan T, 1)
	go func() {
		var maxValue T
		defer func() {
			out <- maxValue
			close(out)
		}()
		var ok bool
		select {
		case <-ctx.Done():
			return
		case maxValue, ok = <-in:
			if !ok {
				return
			}
		}
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				if cmp(value, maxValue) > 0 {
					maxValue = value
				}
			}
		}
	}()
	return out
}

// Min finds the minimum value from the input channel and returns it through the output channel.
//
// Parameters:
//   - ctx: context for cancellation
//   - in: input channel
//   - cmp: comparison function
//
// Returns:
//   - <-chan T: channel containing the minimum value
func Min[T any](ctx context.Context, in <-chan T, cmp func(a, b T) int) <-chan T {
	out := make(chan T, 1)
	go func() {
		var minValue T
		defer func() {
			out <- minValue
			close(out)
		}()
		var ok bool
		select {
		case <-ctx.Done():
			return
		case minValue, ok = <-in:
			if !ok {
				return
			}
		}
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				if cmp(value, minValue) < 0 {
					minValue = value
				}
			}
		}
	}()
	return out
}

// AllMatch checks if all elements in the channel satisfy the given predicate.
// Returns true if all match, false otherwise.
//
// Parameters:
//   - ctx: context for cancellation
//   - in: input channel
//   - predicate: condition function
//
// Returns:
//   - <-chan bool: channel containing the result
func AllMatch[T any](ctx context.Context, in <-chan T, predicate func(value T) bool) <-chan bool {
	var out chan bool
	if in == nil {
		return out
	}
	out = make(chan bool, 1)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-in:
				if !ok {
					out <- true
					return
				}
				if !predicate(value) {
					out <- false
					return
				}
			}
		}
	}()
	return out
}

// AnyMatch checks if any element in the channel satisfies the given predicate.
// Returns true if any match is found, false otherwise.
//
// Parameters:
//   - ctx: context for cancellation
//   - in: input channel
//   - predicate: condition function
//
// Returns:
//   - <-chan bool: channel containing the result
func AnyMatch[T any](ctx context.Context, in <-chan T, predicate func(value T) bool) <-chan bool {
	var out chan bool
	if in == nil {
		return out
	}
	out = make(chan bool, 1)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-in:
				if !ok {
					out <- false
					return
				}
				if predicate(value) {
					out <- true
					return
				}
			}
		}
	}()
	return out
}

// NoneMatch checks if no elements in the channel satisfy the given predicate.
// Returns true if none match, false otherwise.
//
// Parameters:
//   - ctx: context for cancellation
//   - in: input channel
//   - predicate: condition function
//
// Returns:
//   - <-chan bool: channel containing the result
func NoneMatch[T any](ctx context.Context, in <-chan T, predicate func(value T) bool) <-chan bool {
	var out chan bool
	if in == nil {
		return out
	}
	out = make(chan bool, 1)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-in:
				if !ok {
					out <- true
					return
				}
				if predicate(value) {
					out <- false
					return
				}
			}
		}
	}()
	return out
}

// Map receives an input channel and a transformation function, processes each element
// through the transformation function and places it into a new channel, returning this new channel.
func Map[T any, R any](in <-chan T, mapper func(T) R) <-chan R {
	var out chan R
	if in == nil {
		return out
	}
	out = make(chan R, cap(in))
	for value := range in {
		out <- mapper(value)
	}
	close(out)
	return out
}

// AsyncMap receives a data channel and a mapping function, asynchronously processes each element
// in the channel and outputs the processed results through another channel.
// Processing stops when the input channel is closed or context is cancelled.
// If the input channel is nil, returns an empty output channel.
func AsyncMap[T any, R any](ctx context.Context, in <-chan T, mapper func(T) R) <-chan R {
	var out chan R
	if in == nil {
		return out
	}
	out = make(chan R, cap(in))
	go func() {
		defer close(out)
		for {
			select {
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- mapper(value)
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Limit reads from the input channel and sends items to the output channel
// until maxSize reaches 0 or the input channel is closed.
// If in is nil, returns nil.
// Supports early exit via ctx.Done() to cancel the operation.
func Limit[T any](ctx context.Context, in <-chan T, maxSize int) <-chan T {
	var out chan T
	if in == nil {
		return out
	}
	out = make(chan T, cap(in))
	go func() {
		defer close(out)
		for {
			if maxSize < 0 {
				return
			}
			var value T
			var ok bool
			select {
			case <-ctx.Done():
				return
			case value, ok = <-in:
				if !ok {
					return
				}
			}
			select {
			case <-ctx.Done():
				return
			case out <- value:
				maxSize--
			}
		}
	}()
	return out
}

// AsyncGenerate asynchronously generates values of type T and outputs them through a channel.
// It receives a context and a supplier function, and stops generating when the context is cancelled.
func AsyncGenerate[T any](ctx context.Context, supplier func(context.Context) T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			value := supplier(ctx)
			select {
			case out <- value:
				// 发送value
			case <-ctx.Done():
				// 被下游打断
				return
			}
		}
	}()
	return out
}

// AsyncIterate creates a channel that asynchronously generates values of type T,
// continuously producing new values based on an initial seed and iteration function
// until the context is cancelled.
func AsyncIterate[T any](ctx context.Context, seed T, f func(context.Context, T) T) <-chan T {
	out := make(chan T)
	go func(seed T) {
		defer close(out)
		pre := seed
		for {
			value := f(ctx, pre)
			select {
			case out <- value:
				// 发送value
			case <-ctx.Done():
				// 被下游打断
				return
			}
			pre = value
		}
	}(seed)
	return out
}

// Filter filters elements from the input channel based on a predicate function.
//
// Parameters:
//   - in: input channel
//   - predicate: filter function
//
// Returns:
//   - <-chan T: filtered output channel
func Filter[T any](in <-chan T, predicate func(value T) bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for value := range in {
			if !predicate(value) {
				continue
			}
			out <- value
		}
	}()
	return out
}

// FanIn merges data from multiple input channels into a single output channel.
// See: [concurrency](https://go.dev/talks/2012/concurrency.slide#28)
func FanIn[T any](ctx context.Context, ins ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	for _, ch := range ins {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- v:
					}
				}
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// FanOut distributes data from the input channel to multiple output channels.
func FanOut[T any](ctx context.Context, in <-chan T, length int) []chan<- T {
	outs := make([]chan<- T, length)
	for i := range outs {
		outs[i] = make(chan<- T)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				for _, out := range outs {
					select {
					case <-ctx.Done():
						return
					case out <- v:
					}
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		for i := 0; i < len(outs); i++ {
			close(outs[i])
		}
	}()
	return outs
}

func _FanIn[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		// 构造SelectCase slice
		var cases []reflect.SelectCase
		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 循环，从cases中选择一个可用的
		for len(cases) > 0 {
			chosen, recv, ok := reflect.Select(cases)
			if !ok { // 此channel已经close
				cases = append(cases[:chosen], cases[chosen+1:]...)
				continue
			}
			v, ok := recv.Interface().(T)
			if ok {
				out <- v
			}
		}
	}()
	return out
}

// Distinct filters out duplicate elements from the input channel.
//
// Parameters:
//   - in: input channel
//   - cmp: comparison function to determine duplicates
//
// Returns:
//   - <-chan T: output channel with distinct elements
func Distinct[T any](in <-chan T, cmp func(value, stored T) bool) <-chan T {
	out := make(chan T, cap(in))
	go func() {
		defer close(out)
		values := make([]T, 0, cap(in))
		for value := range in {
			// 判断是否已经发送过
			if slices.ContainsFunc[[]T](values, func(stored T) bool { return cmp(value, stored) }) {
				continue
			}
			out <- value
			values = append(values, value)
		}
	}()
	return out
}

// Discard discards all values from the channel until it is closed.
//
// Parameters:
//   - c: channel to discard from
func Discard[T any](c <-chan T) {
	for range c {
	}
}

// AsyncDiscard asynchronously discards all values from the channel until it is closed.
//
// Parameters:
//   - c: channel to discard from
func AsyncDiscard[T any](c <-chan T) {
	go Discard(c)
}

// AsSlice collects all elements from a channel and returns them as a slice.
//
// Parameters:
//   - ctx: context for cancellation
//   - in: input channel
//
// Returns:
//   - <-chan []T: channel containing the collected slice
func AsSlice[T any](ctx context.Context, in <-chan T) <-chan []T {
	out := make(chan []T, 1)
	go func() {
		defer close(out)
		s := make([]T, 0)
		for {
			select {
			case <-ctx.Done():
				out <- s
				return
			case t, ok := <-in:
				if ok {
					s = append(s, t)
					continue
				}
				out <- s
				return
			}
		}
	}()
	return out
}
