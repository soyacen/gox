package brave

import (
	"errors"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	t.Run("no panic", func(t *testing.T) {
		executed := false
		Do(func() {
			executed = true
		})
		if !executed {
			t.Error("function should be executed")
		}
	})

	t.Run("with panic and default recovery", func(t *testing.T) {
		Do(func() {
			panic("test panic")
		})
		// Should not panic, test passes if no panic occurs
	})

	t.Run("with panic and custom recovery", func(t *testing.T) {
		recovered := false
		var recoveredValue any
		Do(func() {
			panic("test panic")
		}, func(p any, stack []byte) {
			recovered = true
			recoveredValue = p
		})

		if !recovered {
			t.Error("custom recovery function should be called")
		}
		if recoveredValue != "test panic" {
			t.Errorf("expected recovered value to be 'test panic', got %v", recoveredValue)
		}
	})
}

func TestDoE(t *testing.T) {
	t.Run("no panic, no error", func(t *testing.T) {
		err := DoE(func() error {
			return nil
		})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("no panic, with error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		err := DoE(func() error {
			return expectedErr
		})
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("with panic and default recovery", func(t *testing.T) {
		err := DoE(func() error {
			panic("test panic")
		})
		if err == nil {
			t.Error("expected error from panic, got nil")
		}
	})

	t.Run("with panic and custom recovery", func(t *testing.T) {
		expectedErr := errors.New("custom error")
		err := DoE(func() error {
			panic("test panic")
		}, func(p any, stack []byte) error {
			return expectedErr
		})
		if err != expectedErr {
			t.Errorf("expected custom error, got %v", err)
		}
	})
}

func TestDoRE(t *testing.T) {
	t.Run("no panic, no error", func(t *testing.T) {
		result, err := DoRE(func() (string, error) {
			return "success", nil
		})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "success" {
			t.Errorf("expected result 'success', got %v", result)
		}
	})

	t.Run("no panic, with error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		result, err := DoRE(func() (string, error) {
			return "", expectedErr
		})
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if result != "" {
			t.Errorf("expected empty result, got %v", result)
		}
	})

	t.Run("with panic and default recovery", func(t *testing.T) {
		result, err := DoRE(func() (string, error) {
			panic("test panic")
		})
		if err == nil {
			t.Error("expected error from panic, got nil")
		}
		if result != "" {
			t.Errorf("expected empty result, got %v", result)
		}
	})

	t.Run("with panic and custom recovery", func(t *testing.T) {
		expectedErr := errors.New("custom error")
		result, err := DoRE(func() (string, error) {
			panic("test panic")
		}, func(p any, stack []byte) error {
			return expectedErr
		})
		if err != expectedErr {
			t.Errorf("expected custom error, got %v", err)
		}
		if result != "" {
			t.Errorf("expected empty result, got %v", result)
		}
	})
}

func TestGo(t *testing.T) {
	t.Run("async execution without panic", func(t *testing.T) {
		executed := make(chan bool, 1)
		Go(func() {
			executed <- true
		})
		select {
		case <-executed:
			// Success
		case <-time.After(1 * time.Second):
			t.Error("function was not executed within timeout")
		}
	})

	t.Run("async execution with panic", func(t *testing.T) {
		Go(func() {
			panic("test panic")
		})
		// Should not panic the test goroutine
		time.Sleep(100 * time.Millisecond) // Give time for goroutine to execute
	})
}

func TestGoE(t *testing.T) {
	t.Run("async execution without panic, no error", func(t *testing.T) {
		errCh := GoE(func() error {
			return nil
		})
		select {
		case err := <-errCh:
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		case <-time.After(1 * time.Second):
			t.Error("error was not received within timeout")
		}
	})

	t.Run("async execution without panic, with error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		errCh := GoE(func() error {
			return expectedErr
		})
		select {
		case err := <-errCh:
			if err != expectedErr {
				t.Errorf("expected error %v, got %v", expectedErr, err)
			}
		case <-time.After(1 * time.Second):
			t.Error("error was not received within timeout")
		}
	})

	t.Run("async execution with panic", func(t *testing.T) {
		errCh := GoE(func() error {
			panic("test panic")
		})
		select {
		case err := <-errCh:
			if err == nil {
				t.Error("expected error from panic, got nil")
			}
		case <-time.After(1 * time.Second):
			t.Error("error was not received within timeout")
		}
	})
}

func TestGoRE(t *testing.T) {
	t.Run("async execution without panic, no error", func(t *testing.T) {
		retCh, errCh := GoRE(func() (string, error) {
			return "success", nil
		})

		var result string
		var err error

		// Wait for result or error
		select {
		case result = <-retCh:
		case <-time.After(1 * time.Second):
			t.Error("result was not received within timeout")
			return
		}

		select {
		case err = <-errCh:
		case <-time.After(1 * time.Second):
			t.Error("error channel was not closed within timeout")
			return
		}

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "success" {
			t.Errorf("expected result 'success', got %v", result)
		}
	})

	t.Run("async execution without panic, with error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		retCh, errCh := GoRE(func() (string, error) {
			return "", expectedErr
		})

		var result string
		var err error

		// Should receive error first
		select {
		case err = <-errCh:
			if err != expectedErr {
				t.Errorf("expected error %v, got %v", expectedErr, err)
			}
		case <-time.After(1 * time.Second):
			t.Error("error was not received within timeout")
			return
		}

		// Result channel should be closed
		select {
		case result = <-retCh:
			if result != "" {
				t.Errorf("expected empty result, got %v", result)
			}
		case <-time.After(1 * time.Second):
			t.Error("result channel was not closed within timeout")
		}
	})

	t.Run("async execution with panic", func(t *testing.T) {
		retCh, errCh := GoRE(func() (string, error) {
			panic("test panic")
		})

		var result string
		var err error

		// Should receive error from panic
		select {
		case err = <-errCh:
			if err == nil {
				t.Error("expected error from panic, got nil")
			}
		case <-time.After(1 * time.Second):
			t.Error("error was not received within timeout")
			return
		}

		// Result channel should be closed
		select {
		case result = <-retCh:
			if result != "" {
				t.Errorf("expected empty result, got %v", result)
			}
		case <-time.After(1 * time.Second):
			t.Error("result channel was not closed within timeout")
		}
	})
}
