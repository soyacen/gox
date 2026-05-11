package asyncbatch

import (
	"sync"
	"testing"
	"time"
)

func TestBatchBySize(t *testing.T) {
	ch := make(chan []int, 1)
	g, err := New[int](func(objs []int) { ch <- objs }, Size(3), Interval(time.Second))
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
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for batch by size")
	}
}

func TestBatchByInterval(t *testing.T) {
	ch := make(chan []int, 1)
	g, err := New[int](func(objs []int) { ch <- objs }, Size(5), Interval(50*time.Millisecond))
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
	g, err := New[int](func(objs []int) { ch <- objs }, Size(10), Interval(time.Second))
	if err != nil {
		t.Fatalf("New error: %v", err)
	}

	if err := g.Submit(1); err != nil {
		t.Fatalf("Submit error: %v", err)
	}
	if err := g.Submit(2); err != nil {
		t.Fatalf("Submit error: %v", err)
	}

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
	g, err := New[int](
		func(objs []int) { panic("boom") },
		Size(1),
		Interval(time.Second),
		Recover(func(p any, stack []byte) { recCh <- p }),
	)
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
	g, err := New[int](
		func(objs []int) {
			for _, v := range objs {
				processed <- v
			}
		},
		Size(10),
		Interval(100*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(total)
	for i := 0; i < total; i++ {
		i := i
		go func() {
			defer wg.Done()
			_ = g.Submit(i)
		}()
	}
	wg.Wait()

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

func TestNewNilTask(t *testing.T) {
	g, err := New[int](nil)
	if err == nil {
		t.Fatal("expected error when task is nil")
	}
	if g != nil {
		t.Fatalf("expected nil Group, got %v", g)
	}
}

func TestSubmitAfterClose(t *testing.T) {
	g, err := New[int](func(objs []int) {}, Size(10), Interval(time.Second))
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	if err := g.Close(); err != nil {
		t.Fatalf("Close error: %v", err)
	}
	if err := g.Submit(1); err == nil {
		t.Fatal("expected ErrClosed when submitting to a closed Group")
	}
}

func TestDoubleClose(t *testing.T) {
	g, err := New[int](func(objs []int) {}, Size(10), Interval(time.Second))
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	if err := g.Close(); err != nil {
		t.Fatalf("first Close error: %v", err)
	}
	if err := g.Close(); err == nil {
		t.Fatal("expected ErrClosed on second close")
	}
}
