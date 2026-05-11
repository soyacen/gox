package distributed

import (
	"testing"
)

func TestCorput(t *testing.T) {
	testCases := []struct {
		n        int
		base     int
		expected float64
	}{
		{0, 2, 0.0},
		{1, 2, 0.5},
		{2, 2, 0.25},
		{3, 2, 0.75},
		{4, 2, 0.125},
		{10, 2, 0.3125}, // 10 is 1010 in binary, reversed is 0101 = 5, 5/16 = 0.3125
		{0, 10, 0.0},
		{1, 10, 0.1},
		{2, 10, 0.2},
		{10, 10, 0.01},
		{100, 10, 0.001},
	}

	for _, tc := range testCases {
		actual := VanDerCorputSequence(tc.n, tc.base)
		if actual != tc.expected {
			t.Errorf("VanDerCorputSequence(%d, %d) = %f, expected %f", tc.n, tc.base, actual, tc.expected)
		}
	}
}
