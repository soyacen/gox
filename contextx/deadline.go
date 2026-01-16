package contextx

import (
	"context"
	"time"

	"github.com/soyacen/gox/operator"
)

// ShrinkDeadline calculates a new deadline based on the given duration `dur` and context `ctx`. 
// If no deadline is set in `ctx`, it returns the current time plus `dur`.
// If a deadline is set in `ctx`, it returns the earlier of the two deadlines.
func ShrinkDeadline(ctx context.Context, dur time.Duration) time.Time {
	t := time.Now().Add(dur)
	dl, ok :=ctx.Deadline();
	if !ok {
		return t
	}
	return operator.Ternary(t.After(dl), dl, t)
}