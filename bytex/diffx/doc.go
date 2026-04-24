// Package diffx implements a basic diff algorithm equivalent to patience diff.
//
// It provides anchored diff functionality that compares two texts and returns
// the differences in unified diff format. The algorithm runs in O(n log n) time
// by finding unique lines that anchor matching regions.
//
// This package is adapted from golang.org/x/tools/internal/diff.
package diffx
