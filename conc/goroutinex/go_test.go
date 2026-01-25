package goroutinex

import (
	"sync"
	"testing"
)

// TestGoID tests the GoID function to ensure it returns a valid goroutine ID.
func TestGoID(t *testing.T) {
	// Capture the current goroutine ID using the GoID function.
	goroutineID := GoID()

	// Check that the returned value is not zero, assuming that a valid goroutine ID is non-zero.
	if goroutineID == 0 {
		t.Errorf("GoID() returned an invalid goroutine ID: %d", goroutineID)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			goroutineID := GoID()
			t.Log("goroutineID:", goroutineID)
		}()
	}
	wg.Wait()
}
