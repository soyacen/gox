package backoff

import (
	"context"
	"github.com/soyacen/gox/mathx/randx/v2"
	"time"
)

// Random generates a function that waits for a random time in the range [alpha, beta) between calls.
func Random(alpha, beta time.Duration) Func {
	r, err := randx.NewChaCha8()
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, attempt uint) time.Duration {
		return alpha + time.Duration(r.Int64N(beta.Nanoseconds()-alpha.Nanoseconds()))
	}
}

func RandomFactory(alpha, beta time.Duration) Factory {
	return func(delta time.Duration) Func {
		return func(ctx context.Context, attempt uint) time.Duration {
			return Random(alpha, beta)(ctx, attempt)
		}
	}
}
