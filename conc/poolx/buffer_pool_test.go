package poolx

import (
	"bytes"
	"testing"
)

func TestNewBufferPool(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "Valid size",
			size:    1024,
			wantErr: false,
		},
		{
			name:    "Zero size",
			size:    0,
			wantErr: false,
		},
		{
			name:    "Negative size",
			size:    -1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBufferPool(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBufferPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewBufferPool() = nil, want non-nil")
			}
		})
	}
}

func TestBufferPool_GetPut(t *testing.T) {
	pool, err := NewBufferPool(1024)
	if err != nil {
		t.Fatalf("Failed to create buffer pool: %v", err)
	}

	// Test Get
	buf := pool.Get()
	if buf == nil {
		t.Fatal("Get() returned nil")
	}

	// Write some data to buffer
	buf.WriteString("test data")
	if buf.Len() == 0 {
		t.Error("Failed to write data to buffer")
	}

	// Test Put
	pool.Put(buf)

	// Get again and verify it's reset
	buf = pool.Get()
	if buf.Len() != 0 {
		t.Error("Buffer was not reset after Put")
	}
	if buf.Cap() < 1024 {
		t.Error("Buffer capacity is less than expected")
	}
}

func TestNewBucketBufferPool(t *testing.T) {
	tests := []struct {
		name      string
		minSize   int
		maxSize   int
		wantErr   bool
		expectNil bool
	}{
		{
			name:      "Valid range",
			minSize:   128,
			maxSize:   1024,
			wantErr:   false,
			expectNil: false,
		},
		{
			name:      "Equal min and max",
			minSize:   512,
			maxSize:   512,
			wantErr:   false,
			expectNil: false,
		},
		{
			name:      "Invalid range - max < min",
			minSize:   1024,
			maxSize:   512,
			wantErr:   true,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBucketBufferPool(tt.minSize, tt.maxSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBucketBufferPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.expectNil && got != nil {
				t.Error("NewBucketBufferPool() = non-nil, want nil")
			}
			if !tt.expectNil && got == nil {
				t.Error("NewBucketBufferPool() = nil, want non-nil")
			}
		})
	}
}

func TestBucketBufferPool_GetPut(t *testing.T) {
	pool, err := NewBucketBufferPool(128, 1024)
	if err != nil {
		t.Fatalf("Failed to create bucket buffer pool: %v", err)
	}

	tests := []struct {
		name        string
		requestSize int
	}{
		{"Small buffer", 64},
		{"Exact bucket size", 128},
		{"Mid-range buffer", 256},
		{"Large buffer", 2048},
		{"Max buffer", 1024},
		{"Exceeding max", 2048},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := pool.Get(tt.requestSize)
			if buf == nil {
				t.Fatal("Get() returned nil")
			}

			// Write some data
			testData := "test data for buffer pool"
			buf.WriteString(testData)

			// Put back
			pool.Put(buf)

			// For sizes within range, we should get a reset buffer back
			if tt.requestSize <= 1024 {
				newBuf := pool.Get(tt.requestSize)
				if newBuf.Len() != 0 {
					t.Error("Buffer was not properly reset")
				}
				pool.Put(newBuf)
			}
		})
	}
}

func TestBucketBufferPool_findPool(t *testing.T) {
	pool, err := NewBucketBufferPool(128, 1024)
	if err != nil {
		t.Fatalf("Failed to create bucket buffer pool: %v", err)
	}

	tests := []struct {
		name      string
		size      int
		expectNil bool
	}{
		{"Below min size", 64, false}, // Should still find the smallest pool
		{"Exact min size", 128, false},
		{"Between buckets", 192, false},
		{"Exact bucket size", 256, false},
		{"Max size", 1024, false},
		{"Above max size", 2048, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := pool.findPool(tt.size)
			if tt.expectNil && pool != nil {
				t.Error("findPool() = non-nil, want nil")
			}
			if !tt.expectNil && pool == nil {
				t.Error("findPool() = nil, want non-nil")
			}
		})
	}
}

func TestBucketBufferPool_BucketDistribution(t *testing.T) {
	pool, err := NewBucketBufferPool(128, 1024)
	if err != nil {
		t.Fatalf("Failed to create bucket buffer pool: %v", err)
	}

	// Verify that we have the correct number of buckets (128, 256, 512, 1024)
	expectedBuckets := 4
	if len(pool.buckets) != expectedBuckets {
		t.Errorf("Expected %d buckets, got %d", expectedBuckets, len(pool.buckets))
	}

	// Test that buffers are distributed to appropriate buckets
	sizes := []int{128, 129, 255, 256, 257, 511, 512, 513, 1023, 1024}
	for _, size := range sizes {
		buf := pool.Get(size)
		if buf == nil {
			t.Errorf("Get(%d) returned nil", size)
			continue
		}

		if buf.Cap() < size {
			t.Errorf("Buffer capacity %d is less than requested size %d", buf.Cap(), size)
		}

		// Write data and put back
		buf.WriteString("test")
		pool.Put(buf)
	}
}

func BenchmarkBucketBufferPool_GetPut(b *testing.B) {
	pool, _ := NewBucketBufferPool(128, 1024)

	b.Run("Small buffers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := pool.Get(64)
			buf.WriteString("benchmark test data")
			pool.Put(buf)
		}
	})

	b.Run("Medium buffers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := pool.Get(512)
			buf.WriteString("benchmark test data")
			pool.Put(buf)
		}
	})

	b.Run("Large buffers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := pool.Get(2048) // Exceeds max, should create new each time
			buf.WriteString("benchmark test data")
			pool.Put(buf)
		}
	})
}

func BenchmarkStandardBufferCreation(b *testing.B) {
	b.Run("Standard buffer creation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := bytes.NewBuffer(make([]byte, 0, 512))
			buf.WriteString("benchmark test data")
			// No put back, simulating regular allocation
		}
	})
}
