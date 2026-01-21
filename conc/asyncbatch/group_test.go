package asyncbatch

import (
	"sync"
	"testing"
	"time"
)

func TestBatchBySize(t *testing.T) {
	ch := make(chan []int, 1)
	g, err := New[int](3, time.Second, func(p any) {}, func(objs []int) { ch <- objs })
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer g.Close()

	if err := g.Submit(1); err != nil {
		t.Fatalf("Submit error: %v", err)
	}
	if err := g.Submit(2); err != nil {
		t.Fatalf("Submit error: %v", err)
	}
	if err := g.Submit(3); err != nil {
		t.Fatalf("Submit error: %v", err)
	}

	select {
	case batch := <-ch:
		if len(batch) != 3 {
			t.Fatalf("expected batch len 3, got %d", len(batch))
		}
		if batch[0] != 1 || batch[1] != 2 || batch[2] != 3 {
			t.Fatalf("unexpected batch values: %v", batch)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout waiting for batch by size")
	}
}

func TestBatchByInterval(t *testing.T) {
	ch := make(chan []int, 1)
	g, err := New[int](5, 50*time.Millisecond, func(p any) {}, func(objs []int) { ch <- objs })
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer g.Close()

	if err := g.Submit(10); err != nil {
		t.Fatalf("Submit error: %v", err)
	}
	if err := g.Submit(20); err != nil {
		t.Fatalf("Submit error: %v", err)
	}

	select {
	case batch := <-ch:
		if len(batch) != 2 {
			t.Fatalf("expected batch len 2, got %d", len(batch))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for batch by interval")
	}
}

func TestCloseFlushesRemaining(t *testing.T) {
	ch := make(chan []int, 1)
	g, err := New[int](10, time.Second, func(p any) {}, func(objs []int) { ch <- objs })
	if err != nil {
		t.Fatalf("New error: %v", err)
	}

	if err := g.Submit(1); err != nil {
		t.Fatalf("Submit error: %v", err)
	}
	if err := g.Submit(2); err != nil {
		t.Fatalf("Submit error: %v", err)
	}

	// Close should flush remaining items
	if err := g.Close(); err != nil {
		t.Fatalf("Close error: %v", err)
	}

	select {
	case batch := <-ch:
		if len(batch) != 2 {
			t.Fatalf("expected flushed batch len 2, got %d", len(batch))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for flushed batch on Close")
	}
}

func TestRecoverHandlerCalled(t *testing.T) {
	recCh := make(chan any, 1)
	g, err := New[int](1, time.Second, func(p any) { recCh <- p }, func(objs []int) { panic("boom") })
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer g.Close()

	if err := g.Submit(1); err != nil {
		t.Fatalf("Submit error: %v", err)
	}

	select {
	case p := <-recCh:
		if p == nil {
			t.Fatalf("expected non-nil panic value")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for recover handler")
	}
}

func TestConcurrentSubmitAndClose(t *testing.T) {
	total := 200
	processed := make(chan int, total)
	g, err := New[int](10, 100*time.Millisecond, func(p any) {}, func(objs []int) {
		for _, v := range objs {
			processed <- v
		}
	})
	if err != nil {
		t.Fatalf("New error: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(total)
	for i := 0; i < total; i++ {
		i := i
		go func() {
			defer wg.Done()
			if err := g.Submit(i); err != nil {
				// if closed, skip
			}
		}()
	}
	wg.Wait()

	// Close should flush any remaining and wait for loop to finish
	if err := g.Close(); err != nil {
		t.Fatalf("Close error: %v", err)
	}

	count := 0
	for count < total {
		select {
		case <-processed:
			count++
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting for processed items, got %d of %d", count, total)
		}
	}
}
