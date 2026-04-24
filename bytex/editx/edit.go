// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package edit implements buffered position-based editing of byte slices.
package editx

import (
	"fmt"
	"sort"
)

// Buffer accumulates text modifications to apply to a byte slice.
type Buffer struct {
	old []byte
	q   edits
}

// edit records a single text modification.
type edit struct {
	start int
	end   int
	new   string
}

// edits is a sortable list of edits by start offset, breaking ties by end offset.
type edits []edit

func (x edits) Len() int      { return len(x) }
func (x edits) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x edits) Less(i, j int) bool {
	if x[i].start != x[j].start {
		return x[i].start < x[j].start
	}
	return x[i].end < x[j].end
}

// NewBuffer creates a new Buffer to accumulate changes to a byte slice.
//
// The buffer maintains a reference to the data; do not modify old until
// after all edits have been applied.
//
// Parameters:
//   - old: initial byte slice to edit
//
// Returns:
//   - *Buffer: new edit buffer
func NewBuffer(old []byte) *Buffer {
	return &Buffer{old: old}
}

// Insert inserts a string at the specified position.
//
// Parameters:
//   - pos: byte offset where the text should be inserted
//   - new: string to insert
func (b *Buffer) Insert(pos int, new string) {
	if pos < 0 || pos > len(b.old) {
		panic("invalid edit position")
	}
	b.q = append(b.q, edit{pos, pos, new})
}

// Delete removes the text in the range [start, end).
//
// Parameters:
//   - start: byte offset of the first byte to delete
//   - end: byte offset after the last byte to delete
func (b *Buffer) Delete(start, end int) {
	if end < start || start < 0 || end > len(b.old) {
		panic("invalid edit position")
	}
	b.q = append(b.q, edit{start, end, ""})
}

// Replace replaces the text in the range [start, end) with new.
//
// Parameters:
//   - start: byte offset of the first byte to replace
//   - end: byte offset after the last byte to replace
//   - new: replacement string
func (b *Buffer) Replace(start, end int, new string) {
	if end < start || start < 0 || end > len(b.old) {
		panic("invalid edit position")
	}
	b.q = append(b.q, edit{start, end, new})
}

// Bytes applies all queued edits and returns the modified data.
//
// Returns:
//   - []byte: new byte slice with all edits applied
func (b *Buffer) Bytes() []byte {
	// Sort edits by starting position and then by ending position.
	// Breaking ties by ending position allows insertions at point x
	// to be applied before a replacement of the text at [x, y).
	sort.Stable(b.q)

	var new []byte
	offset := 0
	for i, e := range b.q {
		if e.start < offset {
			e0 := b.q[i-1]
			panic(fmt.Sprintf("overlapping edits: [%d,%d)->%q, [%d,%d)->%q", e0.start, e0.end, e0.new, e.start, e.end, e.new))
		}
		new = append(new, b.old[offset:e.start]...)
		offset = e.end
		new = append(new, e.new...)
	}
	new = append(new, b.old[offset:]...)
	return new
}

// String applies all queued edits and returns the modified data as a string.
//
// Returns:
//   - string: result with all edits applied
func (b *Buffer) String() string {
	return string(b.Bytes())
}
