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

// All 将多个输入通道合并为一个输出通道，然后返回一个包含所有输入通道当前值的切片的通道
// See: [Go Concurrency Patterns: Timing out, moving on](https://go.dev/blog/concurrency-timeouts)
func All[T any](ctx context.Context, ins ...<-chan T) <-chan []T {
	return AsSlice(ctx, FanIn(ctx, ins...))
}

// Any 从多个输入通道中任意选择一个值，返回一个包含此值的一个通道，
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

// TrySend 用于尝试将值发送到指定的通道。
func TrySend[T any](in chan<- T, v T) error {
	select {
	case in <- v:
		return nil
	default:
		return ErrDefaultBranch
	}
}

// TryReceive 用于尝试从通道接收数据。
func TryReceive[T any](in <-chan T) (T, error) {
	var v T
	select {
	case v = <-in:
		return v, nil
	default:
		return v, ErrDefaultBranch
	}
}

// Emit 创建一个通道，异步发送传入的值，并允许通过上下文取消发送过程。
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

func Once[M any](m M) <-chan M {
	ch := make(chan M, 1)
	ch <- m
	close(ch)
	return ch
}

func Nothing[M any]() <-chan M {
	ch := make(chan M)
	close(ch)
	return ch
}

func Append[T any](c []chan T, channels ...chan T) []chan T {
	return Append(c, channels...)
}

func AppendSendChannel[T any](c []<-chan T, channels ...chan T) []<-chan T {
	for _, ch := range channels {
		c = append(c, ch)
	}
	return c
}

func AppendReceiveChannel[T any](c []chan<- T, channels ...chan T) []chan<- T {
	for _, ch := range channels {
		c = append(c, ch)
	}
	return c
}

func Copy[T any](dest chan<- T, src <-chan T) {
	for v := range src {
		dest <- v
	}
	close(dest)
}

func AsyncCopy[T any](dest chan<- T, src <-chan T) {
	go func() { Copy[T](dest, src) }()
}

// Pipe copies values from src to dest.
// Deprecated: Do not use. use Copy instead.
func Pipe[T any](src <-chan T, dest chan<- T) {
	Copy[T](dest, src)
}

// AsyncPipe copies values from src to dest asynchronously.
// Deprecated: Do not use. use AsyncCopy instead.
func AsyncPipe[T any](src <-chan T, dest chan<- T) {
	AsyncCopy[T](dest, src)
}

// Skip 是一个泛型函数，用于从输入的通道in中跳过指定数量n的元素，并将剩余的元素传递到输出通道out中。
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

// Pipeline 创建一个数据处理管道。
// 它接收一个输入通道 in 和一个处理函数 f，返回一个输出通道 out。
// 在新的goroutine中，遍历输入通道中的每个元素，应用处理函数后将结果发送到输出通道。
// 如果上下文 ctx 被取消，则提前退出并关闭输出通道。
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

// Reduce 对输入channel中的元素执行归约操作，使用提供的累积函数和初始值，返回最终结果。
func Reduce[T any, R any](ctx context.Context, in <-chan T, identity R, accumulator func(R, T) R) <-chan R {
	return Pipeline[T, R](ctx, in, func(value T) R {
		identity = accumulator(identity, value)
		return identity
	})
}

// Max 用于从输入通道 in 中找到最大值，并通过输出通道 out 返回。
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

// Min 用于从输入通道 in 中找到最小值，并通过输出通道 out 返回。
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

// AllMatch 检查通道中所有元素是否满足给定条件，若全部满足则返回true，否则返回false。
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

// AnyMatch 检查通道中是否有元素满足给定条件，找到即返回 true，否则返回 false。
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

// NoneMatch 检查通道中的所有元素是否都不满足给定条件，若所有元素都不满足则返回true，否则返回false。
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

// Map 函数接收一个输入通道和一个转换函数，将输入通道的每个元素通过转换函数处理后放入新通道，并返回此新通道。
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

// AsyncMap 函数AsyncMap接收数据通道和映射函数，异步处理通道中的每个元素，并通过另一个通道输出处理结果。
// 当输入通道关闭或上下文取消时，处理停止。若输入通道为空，则直接返回空输出通道。
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

// Limit 接收泛型输入通道in和上下文ctx，以及整数maxSize作为参数。
// 当从in读取数据时，向输出通道out发送数据项，直到maxSize达到0或in关闭。
// 如果in为nil，则直接返回nil。
// 支持通过ctx.Done()提前退出以取消操作。
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

// AsyncGenerate 异步生成类型为T的值，并通过通道输出。它接收上下文和一个供应商函数，当上下文被取消时停止生成。
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

// AsyncIterate 创建一个异步生成T类型值的channel，基于初始值和迭代函数不断产生新值，直至被上下文取消。
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

// FanIn 函数将多个输入通道的数据合并到一个输出通道中。
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

// FanOut 函数将输入通道的数据分发到多个输出通道中。
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

func Discard[T any](c <-chan T) {
	for range c {
	}
}

func AsyncDiscard[T any](c <-chan T) {
	go Discard(c)
}

// AsSlice 将一个通道中的所有元素收集并转换为切片返回。
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
