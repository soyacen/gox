package slicex

import (
	"cmp"
	"sort"

	"github.com/soyacen/gox/constraintx"
	"github.com/soyacen/gox/mathx"
	"github.com/soyacen/gox/randx"
	"golang.org/x/exp/slices"
)

// Max returns the maximum value in the slice.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - E: the maximum element, or the zero value if the slice is empty.
func Max[S ~[]E, E cmp.Ordered](s S) E {
	var r E
	if len(s) == 0 {
		return r
	}
	r = s[0]
	for i := 1; i < len(s); i++ {
		r = max(r, s[i])
	}
	return r
}

// Min returns the minimum value in the slice.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - E: the minimum element, or the zero value if the slice is empty.
func Min[S ~[]E, E cmp.Ordered](s S) E {
	var r E
	if len(s) == 0 {
		return r
	}
	r = s[0]
	for i := 1; i < len(s); i++ {
		r = min(r, s[i])
	}
	return r
}

// Sum returns the sum of all elements in the slice.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - E: the sum of all elements.
func Sum[S ~[]E, E constraintx.Numeric](s S) E {
	var r E
	for _, i := range s {
		r += i
	}
	return r
}

// Zip interleaves elements from two slices into a single slice.
// The length of the result is twice the length of the shorter input.
//
// Parameters:
//   - s1: the first input slice.
//   - s2: the second input slice.
//
// Returns:
//   - S: a new slice with elements interleaved as s1[0], s2[0], s1[1], s2[1], ...
func Zip[S ~[]E, E any](s1, s2 S) S {
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}

	r := make(S, minLen*2)
	for i := 0; i < minLen; i++ {
		r[i*2] = s1[i]
		r[i*2+1] = s2[i]
	}

	return r
}

// Merge concatenates multiple slices into a single slice.
//
// Parameters:
//   - ss: the slices to merge.
//
// Returns:
//   - S: a new slice containing all elements in order.
func Merge[S ~[]E, E any](ss ...S) S {
	totalLen := 0
	for _, s := range ss {
		totalLen += len(s)
	}
	r := make(S, totalLen)
	start := 0
	for _, s := range ss {
		copy(r[start:], s)
		start += len(s)
	}
	return r
}

// AppendFirst inserts an element at the beginning of the slice.
//
// Parameters:
//   - s: the input slice.
//   - e: the element to insert.
//
// Returns:
//   - S: a new slice with the element prepended.
func AppendFirst[S ~[]E, E any](s S, e E) S {
	return slices.Insert(s, 0, e)
}

// AppendIfNotContains appends an element only if it is not already present.
//
// Parameters:
//   - s: the input slice.
//   - v: the element to append.
//
// Returns:
//   - S: the slice with the element appended, or unchanged if already present.
func AppendIfNotContains[S ~[]E, E comparable](s S, v E) S {
	if slices.Contains(s, v) {
		return s
	}
	return append(s, v)
}

// Prepend inserts one or more elements at the beginning of the slice.
//
// Parameters:
//   - s: the input slice.
//   - elems: the elements to insert.
//
// Returns:
//   - S: a new slice with the elements prepended.
func Prepend[S ~[]E, E any](s S, elems ...E) S {
	return slices.Insert(s, 0, elems...)
}

// AppendUnique appends an element to the slice if it is not already present.
//
// Parameters:
//   - s: the input slice.
//   - v: the element to append.
//
// Returns:
//   - S: the slice with the element appended, or unchanged if already present.
func AppendUnique[S ~[]E, E comparable](s S, v E) S {
	return AppendIfNotContains(s, v)
}

// Chunk splits the slice into smaller slices of the given size.
//
// Parameters:
//   - s: the input slice.
//   - size: the maximum size of each chunk.
//
// Returns:
//   - []S: a slice of chunks.
func Chunk[S ~[]E, E any](s S, size int) []S {
	l := len(s)
	ss2 := make([]S, 0, (l+size)/size)
	for i := 0; i < l; i += size {
		if i+size < l {
			ss2 = append(ss2, s[i:i+size])
		} else {
			ss2 = append(ss2, s[i:l])
		}
	}
	return ss2
}

// Concat concatenates multiple slices into a single slice.
//
// Parameters:
//   - ss: the slices to concatenate.
//
// Returns:
//   - S: a new slice containing all elements in order.
func Concat[S ~[]E, E any](ss ...S) S {
	var length int
	for _, s := range ss {
		length += len(s)
	}
	r := make(S, 0, length)
	for _, s := range ss {
		r = append(r, s...)
	}
	return r
}

// Remove deletes all occurrences of the specified values from the slice.
//
// Parameters:
//   - s: the input slice.
//   - vs: the values to remove.
//
// Returns:
//   - S: a new slice with the specified values removed.
func Remove[S ~[]E, E comparable](s S, vs ...E) S {
	return slices.DeleteFunc(s, func(e E) bool { return slices.Contains(vs, e) })
}

// RemoveFunc deletes elements that satisfy the predicate function.
//
// Parameters:
//   - s: the input slice.
//   - f: a function that receives the index and value; returning true removes the element.
//
// Returns:
//   - S: a new slice with matching elements removed.
func RemoveFunc[S ~[]E, E any](s S, f func(i int, v E) bool) S {
	var index int
	return slices.DeleteFunc(s, func(e E) bool {
		del := f(index, e)
		index++
		return del
	})
}

// RemoveAt deletes elements at the specified indices.
//
// Parameters:
//   - array: the input slice.
//   - is: the indices to remove.
//
// Returns:
//   - S: a new slice with the specified indices removed.
func RemoveAt[S ~[]E, E any](array S, is ...int) S {
	return RemoveFunc(array, func(i int, v E) bool {
		return slices.Contains(is, i)
	})
}

// Difference returns the symmetric difference of two slices.
//
// Deprecated: Do not use.
//
// Parameters:
//   - s1: the first input slice.
//   - s2: the second input slice.
//
// Returns:
//   - S: the symmetric difference.
func Difference[S ~[]E, E comparable](s1 S, s2 S) S {
	if len(s1) >= len(s2) {
		return difference(s1, s2)
	}
	return difference(s2, s1)
}

func difference[S ~[]E, E comparable](a S, b S) S {
	var r S
	for _, v := range a {
		if !slices.Contains(b, v) {
			r = append(r, v)
		}
	}
	return r
}

// IsEmpty reports whether the slice is nil or has length zero.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - bool: true if the slice is nil or empty.
func IsEmpty[S ~[]E, E any](s S) bool {
	return len(s) <= 0
}

// IsNotEmpty reports whether the slice has at least one element.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - bool: true if the slice is not empty.
func IsNotEmpty[S ~[]E, E any](s S) bool {
	return len(s) > 0
}

// Filter returns a new slice containing only elements that satisfy the predicate.
//
// Parameters:
//   - s: the input slice.
//   - f: a function that receives the index and value; returning true keeps the element.
//
// Returns:
//   - S: a new slice with only matching elements.
func Filter[S ~[]E, E any](s S, f func(int, E) bool) S {
	var r S
	for i, e := range s {
		if f(i, e) {
			r = append(r, e)
		}
	}
	return r
}

// FindFunc returns the first element that satisfies the predicate and a boolean indicating success.
//
// Parameters:
//   - s: the input slice.
//   - f: the predicate function.
//
// Returns:
//   - E: the found element (zero value if not found).
//   - bool: true if an element was found.
func FindFunc[E any](s []E, f func(E) bool) (E, bool) {
	if i := slices.IndexFunc(s, f); i != -1 {
		return s[i], true
	}
	var e E
	return e, false
}

// IndexOrDefault returns the element at the given index, or a default value if out of bounds.
//
// Parameters:
//   - s: the input slice.
//   - index: the index to access.
//   - d: the default value to return if index is out of bounds.
//
// Returns:
//   - E: the element at index or the default value.
func IndexOrDefault[S ~[]E, E any](s S, index int, d E) E {
	if len(s) > index {
		return s[index]
	}
	return d
}

// LastIndex returns the last index of the specified value, or -1 if not found.
//
// Parameters:
//   - s: the input slice.
//   - v: the value to search for.
//
// Returns:
//   - int: the last index of the value, or -1.
func LastIndex[E comparable](s []E, v E) int {
	for i := len(s) - 1; i > -1; i-- {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// Indexes returns all indices of the specified value.
//
// Parameters:
//   - s: the input slice.
//   - v: the value to search for.
//
// Returns:
//   - []int: all indices where the value occurs.
func Indexes[E comparable](s []E, v E) []int {
	var r []int
	for i, vs := range s {
		if v == vs {
			r = append(r, i)
		}
	}
	return r
}

// IndexesFunc returns all indices of elements that satisfy the predicate.
//
// Parameters:
//   - s: the input slice.
//   - f: the predicate function.
//
// Returns:
//   - []int: all indices where the predicate returns true.
func IndexesFunc[E any](s []E, f func(E) bool) []int {
	var r []int
	for i, v := range s {
		if f(v) {
			r = append(r, i)
		}
	}
	return r
}

// Insert inserts an element at the specified index.
//
// Parameters:
//   - s: the input slice.
//   - i: the index at which to insert.
//   - e: the element to insert.
//
// Returns:
//   - S: a new slice with the element inserted.
func Insert[S ~[]E, E any](s S, i int, e E) S {
	r := make(S, len(s)+1)
	copy(r, s[:i])
	r[i] = e
	copy(r[i+1:], s[i:])
	return r
}

// IsSameLength reports whether two slices have the same length.
//
// Parameters:
//   - s1: the first slice.
//   - s2: the second slice.
//
// Returns:
//   - bool: true if both slices have the same length.
func IsSameLength[S ~[]E, E any](s1 S, s2 S) bool {
	return len(s1) == len(s2)
}

// Map creates a new slice by applying a function to each element.
//
// Parameters:
//   - s: the input slice.
//   - f: a function that receives the index and value and returns a new value.
//
// Returns:
//   - S2: a new slice with transformed elements.
func Map[S1 ~[]E1, S2 ~[]E2, E1 any, E2 any](s S1, f func(int, E1) E2) S2 {
	if s == nil {
		return nil
	}
	s2 := make(S2, 0, len(s))
	for i, e := range s {
		s2 = append(s2, f(i, e))
	}
	return s2
}

// MapElem looks up a value in s1 and returns the corresponding element at the same position in s2.
//
// Parameters:
//   - v: the value to look up in s1.
//   - s1: the lookup slice.
//   - s2: the result slice.
//
// Returns:
//   - E2: the corresponding element from s2, or the zero value if not found.
func MapElem[S1 ~[]E1, S2 ~[]E2, E1 comparable, E2 any](v E1, s1 S1, s2 S2) E2 {
	var r E2
	if len(s1) != len(s2) {
		panic("length of s1 and s2 not equal")
	}
	for i := range s1 {
		if s1[i] == v {
			return s2[i]
		}
	}
	return r
}

// ContainsAny reports whether any of the values are present in the slice.
//
// Parameters:
//   - s: the input slice.
//   - vs: the values to search for.
//
// Returns:
//   - bool: true if at least one value is found.
func ContainsAny[E comparable](s []E, vs ...E) bool {
	for _, v := range vs {
		if slices.Contains(s, v) {
			return true
		}
	}
	return false
}

// ContainsAll reports whether all of the values are present in the slice.
//
// Parameters:
//   - s: the input slice.
//   - vs: the values to search for.
//
// Returns:
//   - bool: true if all values are found.
func ContainsAll[E comparable](s []E, vs ...E) bool {
	for _, v := range vs {
		if NotContains(s, v) {
			return false
		}
	}
	return true
}

// NotContains reports whether the value is not present in the slice.
//
// Parameters:
//   - s: the input slice.
//   - v: the value to check.
//
// Returns:
//   - bool: true if the value is not present.
func NotContains[E comparable](s []E, v E) bool {
	return !slices.Contains(s, v)
}

// NotContainsFunc reports whether no element satisfies the predicate.
//
// Parameters:
//   - s: the input slice.
//   - f: the predicate function.
//
// Returns:
//   - bool: true if no element satisfies the predicate.
func NotContainsFunc[E any](s []E, f func(E) bool) bool {
	return !slices.ContainsFunc(s, f)
}

// PadStart pads the slice on the left with the given value until it reaches the specified size.
//
// Parameters:
//   - s: the input slice.
//   - size: the desired length.
//   - val: the value to pad with.
//
// Returns:
//   - S: a new slice padded on the left.
func PadStart[S ~[]E, E any](s S, size int, val E) S {
	if size <= len(s) {
		return slices.Clone(s)
	}
	r := make(S, 0, size)
	for i := 0; i < (size - len(s)); i++ {
		r = append(r, val)
	}
	r = append(r, s...)
	return r
}

// PadEnd pads the slice on the right with the given value until it reaches the specified size.
//
// Parameters:
//   - s: the input slice.
//   - size: the desired length.
//   - val: the value to pad with.
//
// Returns:
//   - S: a new slice padded on the right.
func PadEnd[S ~[]E, E any](s S, size int, val E) S {
	if size <= len(s) {
		return slices.Clone(s)
	}
	r := make(S, 0, size)
	r = append(r, s...)
	for i := 0; i < (size - len(s)); i++ {
		r = append(r, val)
	}
	return r
}

// Reduce reduces the slice to a single value using an accumulator function.
//
// Parameters:
//   - s: the input slice.
//   - initValue: the initial value of the accumulator.
//   - f: the reducer function receiving accumulator, current value, index, and the slice.
//
// Returns:
//   - R: the final accumulated value.
func Reduce[S ~[]E, E any, R any](s S, initValue R, f func(previousValue R, currentValue E, currentIndex int, s S) R) R {
	r := initValue
	for i, e := range s {
		r = f(r, e, i, s)
	}
	return r
}

// Reverse reverses the elements of the slice in place.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - S: the reversed slice (same underlying array).
func Reverse[S ~[]E, E any](s S) S {
	slices.Reverse(s)
	return s
}

// ToSet converts the slice to a set represented as a map.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - M: a map with slice elements as keys.
func ToSet[S ~[]E, M ~map[E]struct{}, E comparable](s S) M {
	r := make(M)
	for _, e := range s {
		r[e] = struct{}{}
	}
	return r
}

// SetAll sets each element of the slice using the provided function.
//
// Parameters:
//   - s: the input slice.
//   - f: a function that receives the index and returns the new value.
//
// Returns:
//   - S: the modified slice.
func SetAll[S ~[]E, E any](s S, f func(int) E) S {
	for i := range s {
		s[i] = f(i)
	}
	return s
}

// Shift rotates the slice in place by the given offset.
//
// Parameters:
//   - array: the input slice.
//   - offset: the number of positions to rotate.
func Shift[S ~[]E, E any](array S, offset int) {
	if len(array) <= 0 {
		return
	}
	Shifta(array, 0, len(array), offset)
}

// Shifta rotates a subrange of the slice in place by the given offset.
//
// Parameters:
//   - array: the input slice.
//   - startIndexInclusive: the start index of the range (inclusive).
//   - endIndexExclusive: the end index of the range (exclusive).
//   - offset: the number of positions to rotate.
func Shifta[S ~[]E, E any](array S, startIndexInclusive, endIndexExclusive, offset int) {
	if len(array) <= 0 || startIndexInclusive >= len(array)-1 || endIndexExclusive <= 0 {
		return
	}
	if startIndexInclusive < 0 {
		startIndexInclusive = 0
	}
	if endIndexExclusive >= len(array) {
		endIndexExclusive = len(array)
	}
	n := endIndexExclusive - startIndexInclusive
	if n <= 1 {
		return
	}
	offset %= n
	if offset < 0 {
		offset += n
	}
	// For algorithm explanations and proof of O(n) time complexity and O(1) space complexity
	// see https://beradrian.wordpress.com/2015/04/07/shift-an-array-in-on-in-place/
	for n > 1 && offset > 0 {
		nOffset := n - offset

		if offset > nOffset {
			Swap(array, startIndexInclusive, startIndexInclusive+n-nOffset, nOffset)
			n = offset
			offset -= nOffset
		} else if offset < nOffset {
			Swap(array, startIndexInclusive, startIndexInclusive+nOffset, offset)
			startIndexInclusive += offset
			n = nOffset
		} else {
			Swap(array, startIndexInclusive, startIndexInclusive+nOffset, offset)
			break
		}
	}
}

// Swap swaps a series of elements between two positions in the slice.
//
// Parameters:
//   - array: the input slice.
//   - offset1: the starting index of the first series.
//   - offset2: the starting index of the second series.
//   - length: the number of elements to swap.
func Swap[S ~[]E, E any](array S, offset1, offset2, length int) {
	if IsEmpty(array) || offset1 >= len(array) || offset2 >= len(array) {
		return
	}
	if offset1 < 0 {
		offset1 = 0
	}
	if offset2 < 0 {
		offset2 = 0
	}
	length = mathx.Min(mathx.Min(length, len(array)-offset1), len(array)-offset2)
	for i := 0; i < length; i++ {
		aux := array[offset1]
		array[offset1] = array[offset2]
		array[offset2] = aux
		offset1++
		offset2++
	}
}

// Shuffle randomizes the order of elements in the slice.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - S: the shuffled slice (same underlying array).
func Shuffle[S ~[]E, E any](s S) S {
	r := randx.GetPCG()
	r.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	randx.PutPCG(r)
	return s
}

// ShuffleNoAdjacent shuffles the slice so that no two adjacent elements are equal.
// If it is impossible, the original slice is returned.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - S: the shuffled slice, or the original if impossible.
func ShuffleNoAdjacent[S ~[]E, E comparable](s S) S {
	if len(s) <= 1 {
		return s
	}

	count := make(map[E]int)
	for _, v := range s {
		count[v]++
	}

	maxCount := 0
	for _, c := range count {
		if c > maxCount {
			maxCount = c
		}
	}

	if maxCount > (len(s)+1)/2 {
		return s
	}

	return rearrangeNoAdjacent(s, count)
}

func rearrangeNoAdjacent[S ~[]E, E comparable](s S, count map[E]int) S {
	result := make(S, len(s))

	type elemCount struct {
		elem  E
		count int
	}

	elems := make([]elemCount, 0, len(count))
	for e, c := range count {
		elems = append(elems, elemCount{e, c})
	}

	sort.Slice(elems, func(i, j int) bool {
		return elems[i].count > elems[j].count
	})

	index := 0
	for _, ec := range elems {
		for i := 0; i < ec.count; i++ {
			result[index] = ec.elem
			index += 2
			if index >= len(result) {
				index = 1
			}
		}
	}

	return result
}

// HasAdjacentDuplicates reports whether the slice contains any adjacent duplicate elements.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - bool: true if any adjacent elements are equal.
func HasAdjacentDuplicates[S ~[]E, E comparable](s S) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			return true
		}
	}
	return false
}

// ToMap creates a map from the slice using a key-generating function.
//
// Parameters:
//   - s: the input slice.
//   - f: a function that receives the index and value and returns the key.
//
// Returns:
//   - M: a map with keys from f and values from the slice.
func ToMap[S ~[]E, M ~map[K]E, E any, K comparable](s S, f func(int, E) K) M {
	m := make(M, len(s))
	for i, v := range s {
		m[f(i, v)] = v
	}
	return m
}

// GroupBy groups slice elements by a key-generating function.
//
// Parameters:
//   - s: the input slice.
//   - f: a function that receives the index and value and returns the group key.
//
// Returns:
//   - M: a map where keys are group keys and values are slices of elements.
func GroupBy[M ~map[K]S, S ~[]E, E any, K comparable](s S, f func(int, E) K) M {
	m := make(M, len(s))
	for i, v := range s {
		k := f(i, v)
		m[k] = append(m[k], v)
	}
	return m
}

// Uniq returns a new slice with duplicate elements removed.
// For short slices it checks each element; for long slices it uses a map for efficiency.
//
// Parameters:
//   - s: the input slice.
//
// Returns:
//   - S: a new slice with duplicates removed.
func Uniq[S ~[]E, E comparable](s S) S {
	if s == nil {
		return nil
	}
	if len(s) <= 128 {
		return uniqV1(s)
	}
	return uniqV2(s)
}

func uniqV1[S ~[]E, E comparable](s S) S {
	r := make(S, 0, len(s))
	for _, v := range s {
		if !slices.Contains(r, v) {
			r = append(r, v)
		}
	}
	return r
}

func uniqV2[S ~[]E, E comparable](s S) S {
	r := make(S, 0, len(s))
	m := make(map[E]struct{}, len(s))
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			r = append(r, v)
		}
	}
	return r
}

// SafeSlice returns a sub-slice safely, clamping bounds to the valid range.
//
// Parameters:
//   - s: the input slice.
//   - start: the start index.
//   - length: the number of elements to include.
//
// Returns:
//   - S: the sub-slice, or an empty slice if parameters are invalid.
func SafeSlice[S ~[]E, E comparable](s S, start, length int) S {
	var r S
	if start < 0 || length < 0 {
		return r
	}
	low := start
	if len(s) < low {
		return r
	}
	high := start + length
	if len(s) < high {
		high = len(s)
	}
	return s[low:high]
}
