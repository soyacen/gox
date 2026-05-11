package distributed

import (
	"testing"
)

func TestDeterministicSubSetting(t *testing.T) {
	testCases := []struct {
		name        string
		backends    []string
		clientID    int
		subsetSize  int
		expectedLen int
	}{
		{
			name:        "valid subset",
			backends:    []string{"a", "b", "c", "d", "e", "f"},
			clientID:    0,
			subsetSize:  2,
			expectedLen: 2,
		},
		{
			name:        "zero subset size",
			backends:    []string{"a", "b", "c"},
			clientID:    0,
			subsetSize:  0,
			expectedLen: 0,
		},
		{
			name:        "empty backends",
			backends:    []string{},
			clientID:    0,
			subsetSize:  2,
			expectedLen: 0,
		},
		{
			name:        "client ID out of range",
			backends:    []string{"a", "b", "c", "d"},
			clientID:    100,
			subsetSize:  2,
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make a copy of backends because DeterministicSubSetting modifies it
			backendsCopy := make([]string, len(tc.backends))
			copy(backendsCopy, tc.backends)

			result := DeterministicSubSetting(backendsCopy, tc.clientID, tc.subsetSize)
			if len(result) != tc.expectedLen {
				t.Errorf("DeterministicSubSetting() len = %d, want %d", len(result), tc.expectedLen)
			}
		})
	}

	// Test determinism: same client ID should return same subset
	backends := []string{"a", "b", "c", "d", "e", "f"}
	clientID := 5
	subsetSize := 2

	// Make copies for each test
	backends1 := make([]string, len(backends))
	copy(backends1, backends)
	result1 := DeterministicSubSetting(backends1, clientID, subsetSize)

	backends2 := make([]string, len(backends))
	copy(backends2, backends)
	result2 := DeterministicSubSetting(backends2, clientID, subsetSize)

	if len(result1) != len(result2) {
		t.Errorf("DeterministicSubSetting() not deterministic: len(result1) = %d, len(result2) = %d", len(result1), len(result2))
	}

	for i := range result1 {
		if result1[i] != result2[i] {
			t.Errorf("DeterministicSubSetting() not deterministic: result1[%d] = %s, result2[%d] = %s", i, result1[i], i, result2[i])
		}
	}
}

func TestJumpHash(t *testing.T) {
	testCases := []struct {
		name    string
		key     uint64
		buckets int
	}{
		{
			name:    "1 bucket",
			key:     12345,
			buckets: 1,
		},
		{
			name:    "10 buckets",
			key:     12345,
			buckets: 10,
		},
		{
			name:    "100 buckets",
			key:     12345,
			buckets: 100,
		},
		{
			name:    "zero buckets",
			key:     12345,
			buckets: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := JumpHash(tc.key, tc.buckets, nil)
			expectedBuckets := tc.buckets
			if expectedBuckets <= 0 {
				expectedBuckets = 1
			}
			if result < 0 || result >= expectedBuckets {
				t.Errorf("JumpHash() = %d, want between 0 and %d", result, expectedBuckets-1)
			}
		})
	}

	// Test checkAlive function
	t.Run("check alive", func(t *testing.T) {
		key := uint64(12345)
		buckets := 10
		aliveBuckets := map[int]bool{
			0: false,
			1: true,
			2: false,
			3: true,
			4: true,
			5: false,
			6: true,
			7: true,
			8: false,
			9: true,
		}
		checkAlive := func(b int) bool {
			return aliveBuckets[b]
		}

		result := JumpHash(key, buckets, checkAlive)
		if !aliveBuckets[result] {
			t.Errorf("JumpHash() returned dead bucket %d", result)
		}
	})
}
