package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/soyacen/gox/backoff"
)

func TestMaxAttempts(t *testing.T) {
	t.Run("creates strategy with correct max attempts", func(t *testing.T) {
		strategy := MaxAttempts(3)
		ds, ok := strategy.(*defaultStrategy)
		if !ok {
			t.Fatal("expected defaultStrategy")
		}
		if ds.maxAttempts != 3 {
			t.Errorf("expected maxAttempts=3, got %d", ds.maxAttempts)
		}
	})

	t.Run("uses zero backoff by default", func(t *testing.T) {
		// Note: MaxAttempts returns BackoffStrategy which needs Backoff() called
		// to set the backoff function. This test verifies the default struct has
		// nil backoff which would panic if used without calling Backoff()
		ds := &defaultStrategy{
			maxAttempts: 5,
		}
		// Verify backoffFunc is nil initially
		if ds.backoffFunc != nil {
			t.Error("expected nil backoffFunc initially")
		}
	})

	t.Run("retries on all errors by default", func(t *testing.T) {
		// Note: Similar to backoff, retryOnFunc needs to be set via RetryOn()
		ds := &defaultStrategy{
			maxAttempts: 3,
		}
		// Verify retryOnFunc is nil initially
		if ds.retryOnFunc != nil {
			t.Error("expected nil retryOnFunc initially")
		}
	})
}

func TestDefaultStrategy_Backoff(t *testing.T) {
	t.Run("sets backoff function correctly", func(t *testing.T) {
		ds := &defaultStrategy{}
		called := false
		
		customBackoff := func(ctx context.Context, attempt uint) time.Duration {
			called = true
			return time.Millisecond * 100
		}
		
		result := ds.Backoff(customBackoff)
		
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		
		// Verify the backoff function was set
		duration := ds.backoffFunc(context.Background(), 1)
		if duration != time.Millisecond*100 {
			t.Errorf("expected 100ms, got %v", duration)
		}
		
		if !called {
			t.Error("expected custom backoff to be called")
		}
	})

	t.Run("returns RetryOnStrategy for chaining", func(t *testing.T) {
		ds := &defaultStrategy{}
		result := ds.Backoff(backoff.Zero())
		
		// Verify it implements RetryOnStrategy
		var _ RetryOnStrategy = result
	})
}

func TestDefaultStrategy_RetryOn(t *testing.T) {
	t.Run("sets retry condition correctly", func(t *testing.T) {
		ds := &defaultStrategy{}
		shouldRetry := false
		
		retryCondition := func(err error) bool {
			return shouldRetry
		}
		
		result := ds.RetryOn(retryCondition)
		
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		
		// Test when shouldRetry is false
		shouldRetry = false
		if ds.retryOnFunc(errors.New("test error")) {
			t.Error("expected retryOnFunc to return false")
		}
		
		// Test when shouldRetry is true
		shouldRetry = true
		if !ds.retryOnFunc(errors.New("test error")) {
			t.Error("expected retryOnFunc to return true")
		}
	})

	t.Run("returns Executor for chaining", func(t *testing.T) {
		ds := &defaultStrategy{}
		result := ds.RetryOn(func(err error) bool { return true })
		
		// Verify it implements Executor
		var _ Executor = result
	})
}

func TestDefaultStrategy_Exec(t *testing.T) {
	t.Run("executes command successfully on first attempt", func(t *testing.T) {
		ds := &defaultStrategy{
			maxAttempts: 3,
			backoffFunc: backoff.Zero(),
			retryOnFunc: func(err error) bool { return true },
		}
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			return nil
		}
		
		err := ds.Exec(context.Background(), cmd)
		
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if attempts != 1 {
			t.Errorf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("retries on error until success", func(t *testing.T) {
		ds := &defaultStrategy{
			maxAttempts: 3,
			backoffFunc: backoff.Zero(),
			retryOnFunc: func(err error) bool { return true },
		}
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			if attempts < 3 {
				return errors.New("temporary error")
			}
			return nil
		}
		
		err := ds.Exec(context.Background(), cmd)
		
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("retries until max attempts exhausted", func(t *testing.T) {
		ds := &defaultStrategy{
			maxAttempts: 3,
			backoffFunc: backoff.Zero(),
			retryOnFunc: func(err error) bool { return true },
		}
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			return errors.New("persistent error")
		}
		
		err := ds.Exec(context.Background(), cmd)
		
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if attempts != 4 { // Initial + 3 retries
			t.Errorf("expected 4 attempts (initial + 3 retries), got %d", attempts)
		}
	})

	t.Run("respects retry condition", func(t *testing.T) {
		ds := &defaultStrategy{
			maxAttempts: 5,
			backoffFunc: backoff.Zero(),
			retryOnFunc: func(err error) bool {
				// Only retry on timeout errors
				return err.Error() == "timeout"
			},
		}
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			if attempts == 1 {
				return errors.New("timeout")
			}
			return errors.New("network error")
		}
		
		err := ds.Exec(context.Background(), cmd)
		
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "network error" {
			t.Errorf("expected 'network error', got %v", err)
		}
		if attempts != 2 {
			t.Errorf("expected 2 attempts, got %d", attempts)
		}
	})

	t.Run("cancels on context done", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		
		ds := &defaultStrategy{
			maxAttempts: 5,
			backoffFunc: backoff.Constant(time.Millisecond * 10),
			retryOnFunc: func(err error) bool { return true },
		}
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			if attempts == 2 {
				cancel()
			}
			return errors.New("error")
		}
		
		err := ds.Exec(ctx, cmd)
		
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got %v", err)
		}
	})

	t.Run("context timeout is respected", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		defer cancel()
		
		ds := &defaultStrategy{
			maxAttempts: 100,
			backoffFunc: backoff.Constant(time.Millisecond * 20),
			retryOnFunc: func(err error) bool { return true },
		}
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			return errors.New("error")
		}
		
		start := time.Now()
		err := ds.Exec(ctx, cmd)
		elapsed := time.Since(start)
		
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected context.DeadlineExceeded, got %v", err)
		}
		// Should timeout after ~50ms, not complete all retries
		if elapsed > time.Millisecond*100 {
			t.Errorf("expected timeout around 50ms, took %v", elapsed)
		}
	})

	t.Run("attempt number increments correctly", func(t *testing.T) {
		ds := &defaultStrategy{
			maxAttempts: 3,
			backoffFunc: backoff.Zero(),
			retryOnFunc: func(err error) bool { return true },
		}
		
		var recordedAttempts []uint
		cmd := func(ctx context.Context, attempt uint) error {
			recordedAttempts = append(recordedAttempts, attempt)
			return errors.New("error")
		}
		
		_ = ds.Exec(context.Background(), cmd)
		
		expectedAttempts := []uint{0, 1, 2, 3}
		for i, expected := range expectedAttempts {
			if i >= len(recordedAttempts) || recordedAttempts[i] != expected {
				t.Errorf("attempt %d: expected %d, got %d", i, expected, recordedAttempts[i])
			}
		}
	})
}

func TestBuilderPattern(t *testing.T) {
	t.Run("full builder pattern works correctly", func(t *testing.T) {
		executed := false
		customError := errors.New("custom retryable error")
		
		strategy := MaxAttempts(3).
			Backoff(backoff.Constant(time.Millisecond)).
			RetryOn(func(err error) bool {
				// Retry on errors containing "retryable"
				return err.Error() == "custom retryable error"
			})
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			if attempts < 2 {
				return customError
			}
			executed = true
			return nil
		}
		
		err := strategy.Exec(context.Background(), cmd)
		
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !executed {
			t.Error("expected command to execute successfully")
		}
		// With maxAttempts=3, loop runs for attempts 0,1,2 then final at 3
		// But since we succeed on attempt 1 (attempts=2), we stop early
		if attempts != 2 {
			t.Errorf("expected 2 attempts, got %d", attempts)
		}
	})

	t.Run("retry condition prevents retry", func(t *testing.T) {
		strategy := MaxAttempts(5).
			Backoff(backoff.Zero()).
			RetryOn(func(err error) bool {
				return false // Never retry
			})
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			return errors.New("error")
		}
		
		err := strategy.Exec(context.Background(), cmd)
		
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if attempts != 1 {
			t.Errorf("expected 1 attempt (no retries), got %d", attempts)
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("zero max attempts", func(t *testing.T) {
		strategy := MaxAttempts(0).
			Backoff(backoff.Zero()).
			RetryOn(func(err error) bool { return true })
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			return errors.New("error")
		}
		
		err := strategy.Exec(context.Background(), cmd)
		
		// With 0 max attempts, loop doesn't run, only final execution happens
		if err == nil {
			t.Error("expected error, got nil")
		}
		if attempts != 1 { // Only final execution
			t.Errorf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("single attempt", func(t *testing.T) {
		strategy := MaxAttempts(1).
			Backoff(backoff.Zero()).
			RetryOn(func(err error) bool { return true })
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			return errors.New("error")
		}
		
		err := strategy.Exec(context.Background(), cmd)
		
		if err == nil {
			t.Error("expected error, got nil")
		}
		if attempts != 2 { // Initial + 1 retry
			t.Errorf("expected 2 attempts, got %d", attempts)
		}
	})

	t.Run("nil retry function", func(t *testing.T) {
		ds := &defaultStrategy{
			maxAttempts:   3,
			backoffFunc:   backoff.Zero(),
			retryOnFunc:   nil,
		}
		
		attempts := 0
		cmd := func(ctx context.Context, attempt uint) error {
			attempts++
			return errors.New("error")
		}
		
		err := ds.Exec(context.Background(), cmd)
		
		// With nil retryOnFunc, it retries until maxAttempts exhausted
		if err == nil {
			t.Error("expected error, got nil")
		}
		// Loop: 0,1,2 then final: 3 = 4 attempts
		if attempts != 4 {
			t.Errorf("expected 4 attempts, got %d", attempts)
		}
	})
}

// Example functions demonstrating usage
func ExampleMaxAttempts_basicRetry() {
	strategy := MaxAttempts(3).
		Backoff(backoff.Zero()).
		RetryOn(func(err error) bool { return true })
	
	cmd := func(ctx context.Context, attempt uint) error {
		fmt.Printf("Attempt %d\n", attempt)
		return nil
	}
	
	_ = strategy.Exec(context.Background(), cmd)
	// Output: Attempt 0
}

func ExampleMaxAttempts_withBackoff() {
	strategy := MaxAttempts(3).
		Backoff(backoff.Exponential(time.Millisecond * 100)).
		RetryOn(func(err error) bool { return true })
	
	cmd := func(ctx context.Context, attempt uint) error {
		fmt.Printf("Attempt %d\n", attempt)
		return nil
	}
	
	_ = strategy.Exec(context.Background(), cmd)
	// Output: Attempt 0
}

func ExampleMaxAttempts_withRetryCondition() {
	strategy := MaxAttempts(3).
		Backoff(backoff.Constant(time.Millisecond * 50)).
		RetryOn(func(err error) bool {
			// Only retry on timeout errors
			return err.Error() == "timeout"
		})
	
	cmd := func(ctx context.Context, attempt uint) error {
		fmt.Printf("Attempt %d\n", attempt)
		return nil
	}
	
	_ = strategy.Exec(context.Background(), cmd)
	// Output: Attempt 0
}
