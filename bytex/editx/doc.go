// Package editx implements buffered position-based editing of byte slices.
//
// It provides a Buffer type that accumulates text modifications (insertions,
// deletions, and replacements) and applies them efficiently to produce the
// edited result.
//
// This package is adapted from golang.org/x/tools/internal/edit.
package editx
