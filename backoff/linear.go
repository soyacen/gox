package backoff

import (
	"context"
	"math"
	"time"
)

// Linear it waits for "delta * attempt" time between calls.
func Linear(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return linear(delta, attempt)
	}
}

func linear(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(attempt)
}

// Linear15 1.5倍线性退避
func Linear15(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return linear15(delta, attempt)
	}
}

func linear15(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(float64(attempt)*1.5)
}

// LinearSqrt 线性平方根退避（线性与等差增量之间最完美）
func LinearSqrt(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return linearSqrt(delta, attempt)
	}
}

func linearSqrt(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(float64(attempt)*math.Sqrt(float64(attempt)))
}

// LinearPlus 线性+半量增长
func LinearPlus(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return linearPlus(delta, attempt)
	}
}

func linearPlus(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(int64(attempt)+int64(attempt)/2)
}

func LinearFactory() Factory {
	return func(delta time.Duration) Func {
		return Linear(delta)
	}
}
