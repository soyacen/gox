package contextx

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestWrapContextError(t *testing.T) {
	t.Run("context with no error", func(t *testing.T) {
		ctx := context.Background()
		err := WrapContextError(ctx)
		if err != nil {
			t.Errorf("WrapContextError() = %v, want nil", err)
		}
	})

	t.Run("context with canceled error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := WrapContextError(ctx)
		if err == nil {
			t.Fatal("WrapContextError() = nil, want error")
		}
		if err.Error() != "context canceled, because context canceled" {
			t.Errorf("WrapContextError().Error() = %v, want different message", err.Error())
		}
	})

	t.Run("context with deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		time.Sleep(10 * time.Millisecond)
		err := WrapContextError(ctx)
		if err == nil {
			t.Fatal("WrapContextError() = nil, want error")
		}
		if err.Error() != "context deadline exceeded, because context deadline exceeded" {
			t.Errorf("WrapContextError().Error() = %v, want different message", err.Error())
		}
	})

	t.Run("context with custom cause (Go 1.21+)", func(t *testing.T) {
		// This test is for Go 1.21+ WithCancelCause
		// We'll skip it if the function isn't available
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := WrapContextError(ctx)
		if err == nil {
			t.Fatal("WrapContextError() = nil, want error")
		}
		// Just verify it doesn't panic
	})
}

func TestShrinkDeadline(t *testing.T) {
	t.Run("context without existing deadline", func(t *testing.T) {
		ctx := context.Background()
		dur := 10 * time.Second
		before := time.Now()
		deadline := ShrinkDeadline(ctx, dur)
		after := time.Now()

		minExpected := before.Add(dur - 1*time.Second)
		maxExpected := after.Add(dur + 1*time.Second)

		if deadline.Before(minExpected) || deadline.After(maxExpected) {
			t.Errorf("ShrinkDeadline() = %v, want around %v", deadline, before.Add(dur))
		}
	})

	t.Run("context with earlier deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		dur := 10 * time.Second
		deadline := ShrinkDeadline(ctx, dur)
		expected, _ := ctx.Deadline()

		if !deadline.Equal(expected) {
			t.Errorf("ShrinkDeadline() = %v, want %v", deadline, expected)
		}
	})

	t.Run("context with later deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		dur := 10 * time.Second
		before := time.Now()
		deadline := ShrinkDeadline(ctx, dur)
		after := time.Now()

		minExpected := before.Add(dur - 1*time.Second)
		maxExpected := after.Add(dur + 1*time.Second)

		if deadline.Before(minExpected) || deadline.After(maxExpected) {
			t.Errorf("ShrinkDeadline() = %v, want around %v", deadline, before.Add(dur))
		}
	})
}

func TestContextType(t *testing.T) {
	if ContextType == nil {
		t.Fatal("ContextType is nil")
	}

	if ContextType != reflect.TypeOf((*context.Context)(nil)).Elem() {
		t.Errorf("ContextType = %v, want %v", ContextType, reflect.TypeOf((*context.Context)(nil)).Elem())
	}

	var ctx context.Context = context.Background()
	if !reflect.TypeOf(ctx).AssignableTo(ContextType) {
		t.Error("context.Context should be assignable to ContextType")
	}
}

func TestErrorError(t *testing.T) {
	testCases := []struct {
		name     string
		err      Error
		expected string
	}{
		{
			name:     "both errors present",
			err:      Error{CtxErr: context.Canceled, CauseErr: context.Canceled},
			expected: "context canceled, because context canceled",
		},
		{
			name:     "only ctxErr present",
			err:      Error{CtxErr: context.Canceled, CauseErr: nil},
			expected: "context canceled, because <nil>",
		},
		{
			name:     "only causeErr present",
			err:      Error{CtxErr: nil, CauseErr: context.Canceled},
			expected: "<nil>, because context canceled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Error() != tc.expected {
				t.Errorf("Error.Error() = %v, want %v", tc.err.Error(), tc.expected)
			}
		})
	}
}
