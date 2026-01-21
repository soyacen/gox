package backoff

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/soyacen/gox/randx"
)

// JitterUp adds random jitter to the interval.
//
// This adds or subtracts time from the interval within a given jitter fraction.
// For example for 10s and jitter 0.1, it will return a time within [9s, 11s])
func JitterUp(backoff Func, jitter float64) Func {
	r, err := randx.NewChaCha8()
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, attempt uint) time.Duration {
		interval := backoff(ctx, attempt)
		return jitterUp(r, interval, jitter)
	}
}

func jitterUp(r *rand.Rand, interval time.Duration, jitter float64) time.Duration {
	multiplier := jitter * (r.Float64()*2 - 1)
	return time.Duration(float64(interval) * (1 + multiplier))
}

func JitterUpFactory(factory Factory, jitter float64) Factory {
	return func(delta time.Duration) Func {
		return JitterUp(factory(delta), jitter)
	}
}
