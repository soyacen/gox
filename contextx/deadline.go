package contextx

import (
	"context"
	"time"

	"github.com/soyacen/gox/operator"
)

// ShrinkDeadline calculates a new deadline based on the given duration and context.
// If no deadline is set in the context, it returns the current time plus the duration.
// If a deadline is set in the context, it returns the earlier of the two deadlines.
//
// Parameters:
//   - ctx: The context to check for an existing deadline
//   - dur: The duration to add to the current time
//
// Returns:
//   - time.Time: The calculated deadline
func ShrinkDeadline(ctx context.Context, dur time.Duration) time.Time {
	t := time.Now().Add(dur)
	dl, ok :=ctx.Deadline();
	if !ok {
		return t
	}
	return operator.Ternary(t.After(dl), dl, t)
}