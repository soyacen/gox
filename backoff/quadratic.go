package backoff

import (
	"context"
	"time"
)

// Quadratic 它在调用之间等待 "delta * incremental(attempt)" 时间（介于线性与斐波那契之间）
func Quadratic(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return incremental(delta, attempt)
	}
}

// incremental 计算等差增量退避：1, 2, 4, 7, 11, 16...
// 公式：sum(1..attempt) = attempt*(attempt+1)/2
func incremental(delta time.Duration, attempt uint) time.Duration {
	n := int64(attempt)
	sum := n * (n + 1) / 2
	return delta * time.Duration(sum)
}
