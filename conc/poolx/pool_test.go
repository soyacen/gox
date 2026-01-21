package poolx

import (
	"sync"
	"testing"
)

func TestNewPool(t *testing.T) {
	t.Run("successful creation with valid functions", func(t *testing.T) {
		newFunc := func() int {
			return 42
		}

		resetFunc := func(o int) {
			// Reset logic here if needed
		}

		pool, err := NewPool(newFunc, resetFunc)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if pool == nil {
			t.Error("expected pool to be created, got nil")
		}
	})

	t.Run("error with nil new function", func(t *testing.T) {
		var newFunc func() int = nil
		resetFunc := func(o int) {
			// Reset logic
		}

		pool, err := NewPool(newFunc, resetFunc)
		if err == nil {
			t.Error("expected error for nil new function, got nil")
		}
		if err != ErrNilNewFunction {
			t.Errorf("expected ErrNilNewFunction, got %v", err)
		}
		if pool != nil {
			t.Error("expected pool to be nil when error occurs")
		}
	})

	t.Run("error with nil reset function", func(t *testing.T) {
		newFunc := func() int {
			return 42
		}
		var resetFunc func(int) = nil

		pool, err := NewPool(newFunc, resetFunc)
		if err == nil {
			t.Error("expected error for nil reset function, got nil")
		}
		if err != ErrNilResetFunction {
			t.Errorf("expected ErrNilResetFunction, got %v", err)
		}
		if pool != nil {
			t.Error("expected pool to be nil when error occurs")
		}
	})
}

func TestPoolGetPut(t *testing.T) {
	t.Run("get and put basic types", func(t *testing.T) {
		counter := 0
		newFunc := func() int {
			counter++
			return counter
		}

		resetCalled := false
		resetFunc := func(o int) {
			resetCalled = true
		}

		pool, err := NewPool(newFunc, resetFunc)
		if err != nil {
			t.Fatalf("failed to create pool: %v", err)
		}

		// First get should create a new value
		val1 := pool.Get()
		if val1 != 1 {
			t.Errorf("expected 1, got %d", val1)
		}

		// Put value back - reset should be called
		resetCalled = false
		pool.Put(val1)
		if !resetCalled {
			t.Error("expected reset function to be called")
		}

		// Get again - should return the pooled value
		val2 := pool.Get()
		if val2 != 1 {
			t.Errorf("expected 1, got %d", val2)
		}
	})

	t.Run("get and put complex types with reset", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		counter := 0
		newFunc := func() *Person {
			counter++
			return &Person{
				Name: "Person",
				Age:  counter,
			}
		}

		// Reset function to clean up the object before putting it back in the pool
		resetFunc := func(p *Person) {
			p.Name = ""
			p.Age = 0
		}

		pool, err := NewPool(newFunc, resetFunc)
		if err != nil {
			t.Fatalf("failed to create pool: %v", err)
		}

		// Get a new person
		person1 := pool.Get()
		if person1.Name != "Person" || person1.Age != 1 {
			t.Errorf("expected Person with age 1, got %+v", person1)
		}

		// Modify the person
		person1.Name = "Modified"
		person1.Age = 99

		// Put it back - reset function should be called
		pool.Put(person1)

		// Get it again - should be reset to zero values
		person2 := pool.Get()
		if person2.Name != "" || person2.Age != 0 {
			t.Errorf("expected reset person, got %+v", person2)
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		type TestData struct {
			ID int
		}

		counter := 0
		var mutex sync.Mutex

		newFunc := func() *TestData {
			mutex.Lock()
			counter++
			id := counter
			mutex.Unlock()
			return &TestData{ID: id}
		}

		resetFunc := func(td *TestData) {
			td.ID = 0
		}

		pool, err := NewPool(newFunc, resetFunc)
		if err != nil {
			t.Fatalf("failed to create pool: %v", err)
		}

		// Test concurrent Get and Put operations
		const goroutines = 10
		const operations = 100

		var wg sync.WaitGroup
		wg.Add(goroutines)

		resetCount := 0
		var resetMutex sync.Mutex

		// Override reset function to count calls
		pool.reset = func(td *TestData) {
			resetMutex.Lock()
			resetCount++
			td.ID = 0
			resetMutex.Unlock()
		}

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()

				// Perform many Get/Put operations
				for j := 0; j < operations; j++ {
					obj := pool.Get()
					// Do some work with obj
					obj.ID = j
					pool.Put(obj)
				}
			}()
		}

		wg.Wait()

		// Verify reset was called
		resetMutex.Lock()
		if resetCount == 0 {
			t.Error("expected reset function to be called")
		}
		resetMutex.Unlock()

		// Test passes if no race conditions or panics occur
	})

	t.Run("reset function behavior", func(t *testing.T) {
		type TestStruct struct {
			Data []int
		}

		newFunc := func() *TestStruct {
			return &TestStruct{
				Data: make([]int, 0, 10),
			}
		}

		resetCalls := 0
		resetFunc := func(ts *TestStruct) {
			resetCalls++
			// Clear the slice but keep capacity
			ts.Data = ts.Data[:0]
		}

		pool, err := NewPool(newFunc, resetFunc)
		if err != nil {
			t.Fatalf("failed to create pool: %v", err)
		}

		// Get an object
		obj := pool.Get()
		obj.Data = append(obj.Data, 1, 2, 3)

		// Put it back
		pool.Put(obj)

		// Verify reset was called
		if resetCalls != 1 {
			t.Errorf("expected reset to be called once, got %d", resetCalls)
		}

		// Get it again
		obj2 := pool.Get()
		if len(obj2.Data) != 0 {
			t.Errorf("expected empty slice, got length %d", len(obj2.Data))
		}
		if cap(obj2.Data) != 10 {
			t.Errorf("expected capacity 10, got %d", cap(obj2.Data))
		}
	})
}

func TestPoolPerformance(t *testing.T) {
	type TestStruct struct {
		Data [100]byte
		ID   int
	}

	allocCount := 0
	var mutex sync.Mutex

	newFunc := func() *TestStruct {
		mutex.Lock()
		allocCount++
		count := allocCount
		mutex.Unlock()
		return &TestStruct{ID: count}
	}

	resetFunc := func(ts *TestStruct) {
		ts.ID = 0
		// Clear data array
		for i := range ts.Data {
			ts.Data[i] = 0
		}
	}

	pool, err := NewPool(newFunc, resetFunc)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	// Allocate and pool several objects
	const initialAlloc = 100
	objects := make([]*TestStruct, initialAlloc)

	for i := 0; i < initialAlloc; i++ {
		objects[i] = pool.Get()
	}

	// Put them all back
	for i := 0; i < initialAlloc; i++ {
		pool.Put(objects[i])
	}

	// Reset allocation counter
	mutex.Lock()
	initialAllocCount := allocCount
	mutex.Unlock()

	// Now get them again - should come from pool
	for i := 0; i < initialAlloc; i++ {
		obj := pool.Get()
		_ = obj
		pool.Put(obj)
	}

	// Check that no new allocations occurred
	mutex.Lock()
	finalAllocCount := allocCount
	mutex.Unlock()

	if finalAllocCount != initialAllocCount {
		t.Errorf("expected no new allocations, but got %d new allocations", finalAllocCount-initialAllocCount)
	}
}

func BenchmarkPool(b *testing.B) {
	type TestStruct struct {
		Data [100]byte
		ID   int
	}

	pool, err := NewPool(func() *TestStruct {
		return &TestStruct{}
	}, func(ts *TestStruct) {
		ts.ID = 0
	})
	if err != nil {
		b.Fatalf("failed to create pool: %v", err)
	}

	b.Run("GetPut", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			obj := pool.Get()
			pool.Put(obj)
		}
	})

	b.Run("DirectAllocation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = &TestStruct{}
		}
	})
}
