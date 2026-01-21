package waiter

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// MockWaiter implements the Wait() interface
type MockWaiter struct {
	waitDuration time.Duration
}

func (m *MockWaiter) Wait() {
	time.Sleep(m.waitDuration)
}

// MockErrorWaiter implements the Wait() error interface
type MockErrorWaiter struct {
	waitDuration time.Duration
	shouldError  bool
}

func (m *MockErrorWaiter) Wait() error {
	time.Sleep(m.waitDuration)
	if m.shouldError {
		return errors.New("mock error")
	}
	return nil
}

// MockContextWaiter implements the Wait(context.Context) error interface
type MockContextWaiter struct {
	waitDuration time.Duration
	shouldError  bool
}

func (m *MockContextWaiter) Wait(ctx context.Context) error {
	select {
	case <-time.After(m.waitDuration):
		if m.shouldError {
			return errors.New("mock context error")
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TestWaitNotify(t *testing.T) {
	t.Run("successful wait", func(t *testing.T) {
		waiter := &MockWaiter{waitDuration: 10 * time.Millisecond}

		done := WaitNotify(waiter)

		select {
		case <-done:
			// Success - channel was closed after wait completed
		case <-time.After(100 * time.Millisecond):
			t.Error("WaitNotify did not complete within expected time")
		}
	})

	t.Run("concurrent waits", func(t *testing.T) {
		waiter := &MockWaiter{waitDuration: 20 * time.Millisecond}

		done1 := WaitNotify(waiter)
		done2 := WaitNotify(waiter)

		// Both channels should be closed after waiter completes
		timeout := time.After(100 * time.Millisecond)
		select {
		case <-done1:
		case <-timeout:
			t.Error("done1 was not closed within expected time")
			return
		}

		select {
		case <-done2:
		case <-timeout:
			t.Error("done2 was not closed within expected time")
		}
	})
}

func TestWaitNotifyE(t *testing.T) {
	t.Run("successful wait without error", func(t *testing.T) {
		waiter := &MockErrorWaiter{
			waitDuration: 10 * time.Millisecond,
			shouldError:  false,
		}

		errChan := WaitNotifyE(waiter)

		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("WaitNotifyE did not complete within expected time")
		}
	})

	t.Run("wait with error", func(t *testing.T) {
		expectedErr := errors.New("mock error")
		waiter := &MockErrorWaiter{
			waitDuration: 10 * time.Millisecond,
			shouldError:  true,
		}

		errChan := WaitNotifyE(waiter)

		select {
		case err := <-errChan:
			if err == nil {
				t.Error("expected an error, got nil")
			} else if err.Error() != expectedErr.Error() {
				t.Errorf("expected error %v, got %v", expectedErr, err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("WaitNotifyE did not complete within expected time")
		}
	})

	t.Run("channel is closed when no error", func(t *testing.T) {
		waiter := &MockErrorWaiter{
			waitDuration: 5 * time.Millisecond,
			shouldError:  false,
		}

		errChan := WaitNotifyE(waiter)

		// Channel should be closed after completion
		timeout := time.After(100 * time.Millisecond)
		done := false
		for !done {
			select {
			case _, ok := <-errChan:
				if !ok {
					// Channel is closed, which is expected
					done = true
				}
				// If we receive a nil error, that's also fine
			case <-timeout:
				t.Error("channel was not closed within expected time")
				return
			}
		}
	})
}

func TestWaitContentNotifyE(t *testing.T) {
	t.Run("successful wait without error", func(t *testing.T) {
		waiter := &MockContextWaiter{
			waitDuration: 10 * time.Millisecond,
			shouldError:  false,
		}
		ctx := context.Background()

		errChan := WaitContentNotifyE(ctx, waiter)

		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("WaitContentNotifyE did not complete within expected time")
		}
	})

	t.Run("wait with error", func(t *testing.T) {
		expectedErr := errors.New("mock context error")
		waiter := &MockContextWaiter{
			waitDuration: 10 * time.Millisecond,
			shouldError:  true,
		}
		ctx := context.Background()

		errChan := WaitContentNotifyE(ctx, waiter)

		select {
		case err := <-errChan:
			if err == nil {
				t.Error("expected an error, got nil")
			} else if err.Error() != expectedErr.Error() {
				t.Errorf("expected error %v, got %v", expectedErr, err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("WaitContentNotifyE did not complete within expected time")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		waiter := &MockContextWaiter{
			waitDuration: 100 * time.Millisecond, // Longer wait
			shouldError:  false,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		errChan := WaitContentNotifyE(ctx, waiter)

		select {
		case err := <-errChan:
			if err == nil {
				t.Error("expected context cancellation error, got nil")
			} else if err != context.DeadlineExceeded {
				t.Errorf("expected context.DeadlineExceeded, got %v", err)
			}
		case <-time.After(200 * time.Millisecond):
			t.Error("WaitContentNotifyE did not complete within expected time")
		}
	})

	t.Run("channel is closed when no error", func(t *testing.T) {
		waiter := &MockContextWaiter{
			waitDuration: 5 * time.Millisecond,
			shouldError:  false,
		}
		ctx := context.Background()

		errChan := WaitContentNotifyE(ctx, waiter)

		// Channel should be closed after completion
		timeout := time.After(100 * time.Millisecond)
		done := false
		for !done {
			select {
			case _, ok := <-errChan:
				if !ok {
					// Channel is closed, which is expected
					done = true
				}
				// If we receive a nil error, that's also fine
			case <-timeout:
				t.Error("channel was not closed within expected time")
				return
			}
		}
	})
}

func TestIntegrationWithSyncWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	// Add some work to the wait group
	wg.Add(2)

	// Start WaitNotify
	done := WaitNotify(&wg)

	// Complete the work in separate goroutines
	go func() {
		time.Sleep(10 * time.Millisecond)
		wg.Done()
	}()

	go func() {
		time.Sleep(20 * time.Millisecond)
		wg.Done()
	}()

	// Wait for completion
	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Error("WaitNotify did not complete when WaitGroup finished")
	}
}

func TestIntegrationWithErrorGroup(t *testing.T) {
	// This test simulates integration with something like golang.org/x/sync/errgroup
	waiter := &MockErrorWaiter{
		waitDuration: 15 * time.Millisecond,
		shouldError:  false,
	}

	errChan := WaitNotifyE(waiter)

	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("WaitNotifyE did not complete within expected time")
	}
}

func TestBufferedChannelBehavior(t *testing.T) {
	t.Run("WaitNotifyE buffered channel", func(t *testing.T) {
		waiter := &MockErrorWaiter{
			waitDuration: 5 * time.Millisecond,
			shouldError:  true,
		}

		errChan := WaitNotifyE(waiter)

		// We should be able to read the error multiple times until channel is closed
		err1 := <-errChan
		if err1 == nil {
			t.Error("expected an error on first read")
		}

		// Second read should either get the same error or nil (if channel is closed)
		select {
		case err2, ok := <-errChan:
			if ok && err2 != nil {
				// This shouldn't happen with current implementation since we only send once
				t.Error("unexpected second error value")
			}
			// If ok is false, channel is closed, which is expected
		case <-time.After(50 * time.Millisecond):
			t.Error("channel read timed out")
		}
	})
}
