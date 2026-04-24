package sortx

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"

	"github.com/soyacen/gox/slicex"
)

// Asc sorts a slice of any ordered type in ascending order.
//
// Parameters:
//   - x: The slice to sort
//
// Returns:
//   - []E: The sorted slice
func Asc[E constraints.Ordered](x []E) []E {
	slices.Sort(x)
	return x
}

// Desc sorts a slice of any ordered type in descending order.
//
// Parameters:
//   - x: The slice to sort
//
// Returns:
//   - []E: The sorted slice in descending order
func Desc[E constraints.Ordered](x []E) []E {
	return slicex.Reverse(Asc(x))
}

// IsAsc reports whether a slice is sorted in ascending order.
//
// Parameters:
//   - x: The slice to check
//
// Returns:
//   - bool: True if sorted in ascending order, false otherwise
func IsAsc[E constraints.Ordered](x []E) bool {
	for i := len(x) - 1; i > 0; i-- {
		if x[i] < x[i-1] {
			return false
		}
	}
	return true
}

// IsDesc reports whether a slice is sorted in descending order.
//
// Parameters:
//   - x: The slice to check
//
// Returns:
//   - bool: True if sorted in descending order, false otherwise
func IsDesc[E constraints.Ordered](x []E) bool {
	for i := len(x) - 1; i > 0; i-- {
		if x[i] > x[i-1] {
			return false
		}
	}
	return true
}

// SortFunc sorts a slice using a custom comparison function.
//
// Parameters:
//   - x: The slice to sort
//   - cmp: The comparison function
//
// Returns:
//   - S: The sorted slice
func SortFunc[S ~[]E, E any](x S, cmp func(a, b E) int) S {
	slices.SortFunc(x, cmp)
	return x
}

// SortStableFunc sorts a slice using a custom comparison function, preserving the order of equal elements.
//
// Parameters:
//   - x: The slice to sort
//   - cmp: The comparison function
//
// Returns:
//   - S: The sorted slice
func SortStableFunc[S ~[]E, E any](x S, cmp func(a, b E) int) S {
	slices.SortStableFunc(x, cmp)
	return x
}

// BubbleSort sorts a slice using the bubble sort algorithm.
//
// Parameters:
//   - x: The slice to sort
//
// Returns:
//   - []E: The sorted slice
func BubbleSort[E constraints.Ordered](x []E) []E {
	for i := 0; i < len(x)-1; i++ {
		for j := 1; j < len(x)-i; j++ {
			if x[j] < x[j-1] {
				x[j], x[j-1] = x[j-1], x[j]
			}
		}
	}
	return x
}

// SelectSort sorts a slice using the selection sort algorithm.
//
// Parameters:
//   - x: The slice to sort
//
// Returns:
//   - []E: The sorted slice
func SelectSort[E constraints.Ordered](x []E) []E {
	for i := 0; i < len(x); i++ {
		min := i
		for j := i + 1; j < len(x); j++ {
			if x[min] > x[j] {
				min = j
			}
		}
		x[i], x[min] = x[min], x[i]
	}
	return x
}

// InsertSort sorts a slice using the insertion sort algorithm.
//
// Parameters:
//   - x: The slice to sort
//
// Returns:
//   - []E: The sorted slice
func InsertSort[E constraints.Ordered](x []E) []E {
	for i := 0; i < len(x); i++ {
		for j := i - 1; j >= 0; j-- {
			if x[j+1] < x[j] {
				x[j+1], x[j] = x[j], x[j+1]
			} else {
				break
			}
		}
	}
	return x
}
