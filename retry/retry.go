package retry

import (
	"context"
	"time"

	"github.com/soyacen/gox/backoff"
)

// Strategy 接口定义了一个 Backoff 方法，用于设置重试间隔策略。
type Strategy interface {

	// Backoff 方法设置重试间隔策略函数。
	Backoff(backoffFunc backoff.Func) Strategy

	// RetryOn 方法判断是否重试。
	RetryOn(retryOnFunc func(err error) bool) Strategy

	// Exec 方法执行一个命令，并允许传递上下文和当前尝试次数。
	Exec(ctx context.Context, cmd func(ctx context.Context, attempt uint) error) error
}

// defaultStrategy 结构体实现了 Retry 和 Executor 接口，并持有最大尝试次数和重试间隔策略函数。
type defaultStrategy struct {
	maxAttempts uint
	backoffFunc backoff.Func
	retryOnFunc func(err error) bool
}

func (r *defaultStrategy) Backoff(backoffFunc backoff.Func) Strategy {
	r.backoffFunc = backoffFunc
	return r
}

func (r *defaultStrategy) RetryOn(retryOnFunc func(err error) bool) Strategy {
	r.retryOnFunc = retryOnFunc
	return r
}

func (r *defaultStrategy) Exec(ctx context.Context, cmd func(ctx context.Context, attempt uint) error) error {
	var attempt uint
	for attempt < r.maxAttempts {
		// execute cmd
		err := cmd(ctx, attempt)
		if err == nil {
			// return if err is nil.
			return nil
		}
		if r.retryOnFunc != nil && !r.retryOnFunc(err) {
			// return if retryOnFunc is not nil and retryOnFunc returns false
			return err
		}
		// increase the number of attempts
		attempt++
		select {
		case <-ctx.Done(): // return if context is done, return context error
			return ctx.Err()
		case <-time.After(r.backoffFunc(ctx, attempt)): // sleep and wait retry
			continue
		}
	}
	// perform the execution
	return cmd(ctx, attempt)
}

// MaxAttempts 函数创建一个具有指定最大重试次数的 defaultStrategy 实例。
func MaxAttempts(maxAttempts uint) Strategy {
	return &defaultStrategy{
		maxAttempts: maxAttempts,
		backoffFunc: backoff.Zero(),
		retryOnFunc: func(err error) bool {
			return true
		},
	}
}

// Call 函数执行一个命令，并允许传递上下文、最大尝试次数和重试间隔策略函数。
// Deprecated: Do not use. use MaxAttempts
func Call(ctx context.Context, maxAttempts uint, backoffFunc backoff.Func, method func(attemptTime int) error) error {
	return MaxAttempts(maxAttempts).Backoff(backoffFunc).Exec(ctx, func(ctx context.Context, attempt uint) error {
		return method(int(attempt))
	})
}
