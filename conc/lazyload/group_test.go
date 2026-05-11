package lazyload

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestGroup_LoadOrNew(t *testing.T) {
	t.Run("Load existing value via Load path", func(t *testing.T) {
		group := Group[int]{
			New: func(key string) (int, error) {
				return 999, nil
			},
		}

		first, err, wasLoaded := group.LoadOrNew("testKey")
		if err != nil {
			t.Fatalf("first call err: %v", err)
		}
		if first != 999 {
			t.Fatalf("first call got %v want 999", first)
		}
		if wasLoaded {
			t.Fatal("first call should report not loaded")
		}

		second, err, wasLoaded := group.LoadOrNew("testKey")
		if err != nil {
			t.Fatalf("second call err: %v", err)
		}
		if second != 999 {
			t.Fatalf("second call got %v want 999", second)
		}
		if !wasLoaded {
			t.Fatal("second call should report loaded from cache")
		}
	})

	t.Run("Create new value", func(t *testing.T) {
		group := Group[int]{
			New: func(key string) (int, error) {
				return 100, nil
			},
		}

		obj, err, wasLoaded := group.LoadOrNew("newKey")

		if err != nil {
			t.Errorf("Error should be nil, got: %v", err)
		}
		if obj != 100 {
			t.Errorf("Expected value 100, got: %v", obj)
		}
		if wasLoaded {
			t.Error("Expected value to be created, not loaded from cache")
		}
	})

	t.Run("Nil function error", func(t *testing.T) {
		group := Group[int]{}

		obj, err, wasLoaded := group.LoadOrNew("newKey")

		if obj != 0 {
			t.Errorf("Expected value to be 0, got: %v", obj)
		}
		if !errors.Is(err, ErrNilFunction) {
			t.Errorf("Expected error to be ErrNilFunction, got: %v", err)
		}
		if wasLoaded {
			t.Error("Expected value to not be loaded")
		}
	})

	t.Run("WithFactory overrides New", func(t *testing.T) {
		group := Group[int]{
			New: func(key string) (int, error) {
				return 1, nil
			},
		}

		obj, err, _ := group.LoadOrNew("factoryKey", WithFactory[int](func(key string) (int, error) {
			return 42, nil
		}))
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if obj != 42 {
			t.Fatalf("got %v want 42", obj)
		}
	})

	t.Run("Factory error propagates", func(t *testing.T) {
		group := Group[int]{}
		want := errors.New("factory failed")

		obj, err, wasLoaded := group.LoadOrNew("errKey", WithFactory[int](func(key string) (int, error) {
			return 0, want
		}))

		if obj != 0 {
			t.Fatalf("expected zero value on error, got %v", obj)
		}
		if !errors.Is(err, want) {
			t.Fatalf("expected wrapped factory error, got %v", err)
		}
		if wasLoaded {
			t.Fatal("expected wasLoaded=false on error")
		}
	})
}

func TestGroup_Load(t *testing.T) {
	group := Group[string]{
		New: func(key string) (string, error) {
			return "value-" + key, nil
		},
	}

	obj, err, wasLoaded := group.Load("k1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if obj != "value-k1" {
		t.Fatalf("got %q want %q", obj, "value-k1")
	}
	if wasLoaded {
		t.Fatal("first Load should not report loaded")
	}

	obj, _, wasLoaded = group.Load("k1")
	if !wasLoaded {
		t.Fatal("second Load should report loaded")
	}
	if obj != "value-k1" {
		t.Fatalf("got %q want %q on cached read", obj, "value-k1")
	}
}

func TestGroup_Close(t *testing.T) {
	var closed []string
	var mu sync.Mutex
	group := Group[string]{
		New: func(key string) (string, error) {
			return key, nil
		},
		Finalize: func(ctx context.Context, obj string) error {
			mu.Lock()
			defer mu.Unlock()
			closed = append(closed, obj)
			return nil
		},
	}

	for _, k := range []string{"a", "b", "c"} {
		if _, err, _ := group.LoadOrNew(k); err != nil {
			t.Fatalf("LoadOrNew(%q) err: %v", k, err)
		}
	}

	if err := group.Close(context.Background()); err != nil {
		t.Fatalf("Close err: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(closed) != 3 {
		t.Fatalf("expected 3 finalize calls, got %d", len(closed))
	}
}

func TestGroup_CloseAggregatesErrors(t *testing.T) {
	want := errors.New("finalize boom")
	group := Group[int]{
		New: func(key string) (int, error) { return 1, nil },
		Finalize: func(ctx context.Context, obj int) error {
			return want
		},
	}
	if _, err, _ := group.LoadOrNew("only"); err != nil {
		t.Fatalf("LoadOrNew err: %v", err)
	}

	if err := group.Close(context.Background()); !errors.Is(err, want) {
		t.Fatalf("expected joined error to wrap want, got %v", err)
	}
}

func TestGroup_ConcurrentLoadOrNewDeduplicates(t *testing.T) {
	var calls int
	var mu sync.Mutex
	group := Group[int]{
		New: func(key string) (int, error) {
			mu.Lock()
			calls++
			mu.Unlock()
			return 7, nil
		},
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, _ = group.LoadOrNew("same")
		}()
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("expected exactly 1 factory call, got %d", calls)
	}
}
