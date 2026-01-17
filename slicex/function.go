package slicex

import (
	"cmp"
	"sort"

	"github.com/soyacen/gox/constraintx"
	"github.com/soyacen/gox/mathx"
	"github.com/soyacen/gox/randx"
	"golang.org/x/exp/slices"
)

// Max 获取数组中的最大值
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

// Min 获取数组中的最小值
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

// Sum 数组求和
func Sum[S ~[]E, E constraintx.Numeric](s S) E {
	var r E
	for _, i := range s {
		r += i
	}
	return r
}

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

func AppendFirst[S ~[]E, E any](s S, e E) S {
	return slices.Insert(s, 0, e)
}

func AppendIfNotContains[S ~[]E, E comparable](s S, v E) S {
	if slices.Contains(s, v) {
		return s
	}
	return append(s, v)
}

func Prepend[S ~[]E, E any](s S, elems ...E) S {
	return slices.Insert(s, 0, elems...)
}

// AppendUnique appends an element to a raw, if the element is not already in the raw
func AppendUnique[S ~[]E, E comparable](s S, v E) S {
	return AppendIfNotContains(s, v)
}

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

func Remove[S ~[]E, E comparable](s S, vs ...E) S {
	return slices.DeleteFunc(s, func(e E) bool { return slices.Contains(vs, e) })
}

func RemoveFunc[S ~[]E, E any](s S, f func(i int, v E) bool) S {
	var index int
	return slices.DeleteFunc(s, func(e E) bool {
		del := f(index, e)
		index++
		return del
	})
}

func RemoveAt[S ~[]E, E any](array S, is ...int) S {
	return RemoveFunc(array, func(i int, v E) bool {
		return slices.Contains(is, i)
	})
}

// Difference 返回差集
// Deprecated: Do not use.
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

// IsEmpty Checks if an raw is nil or length equals 0
func IsEmpty[S ~[]E, E any](s S) bool {
	return len(s) <= 0
}

func IsNotEmpty[S ~[]E, E any](s S) bool {
	return len(s) > 0
}

func Filter[S ~[]E, E any](s S, f func(int, E) bool) S {
	var r S
	for i, e := range s {
		if f(i, e) {
			r = append(r, e)
		}
	}
	return r
}

func FindFunc[E any](s []E, f func(E) bool) (E, bool) {
	if i := slices.IndexFunc(s, f); i != -1 {
		return s[i], true
	}
	var e E
	return e, false
}

func IndexOrDefault[S ~[]E, E any](s S, index int, d E) E {
	if len(s) > index {
		return s[index]
	}
	return d
}

func LastIndex[E comparable](s []E, v E) int {
	for i := len(s) - 1; i > -1; i-- {
		if v == s[i] {
			return i
		}
	}
	return -1
}

func Indexes[E comparable](s []E, v E) []int {
	var r []int
	for i, vs := range s {
		if v == vs {
			r = append(r, i)
		}
	}
	return r
}

func IndexesFunc[E any](s []E, f func(E) bool) []int {
	var r []int
	for i, v := range s {
		if f(v) {
			r = append(r, i)
		}
	}
	return r
}

// Insert 在切片的指定位置插入元素
func Insert[S ~[]E, E any](s S, i int, e E) S {
	// 创建一个新的切片，长度比原切片多1
	r := make(S, len(s)+1)

	// 将原切片的前index个元素复制到新切片
	copy(r, s[:i])

	// 插入要插入的元素
	r[i] = e

	// 将原切片的index及之后的元素复制到新切片
	copy(r[i+1:], s[i:])

	return r
}

func IsSameLength[S ~[]E, E any](s1 S, s2 S) bool {
	return len(s1) == len(s2)
}

// Map 方法创建一个新数组，这个新数组由原数组中的每个元素都调用一次提供的函数后的返回值组成。
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

// MapElem 两个切片s1，s2,，获取v所有s1同位置的s2的值
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

// ContainsAny checks if any of the elem are in the given raw.
func ContainsAny[E comparable](s []E, vs ...E) bool {
	for _, v := range vs {
		if slices.Contains(s, v) {
			return true
		}
	}
	return false
}

// ContainsAll checks if all of the elem are in the given raw.
func ContainsAll[E comparable](s []E, vs ...E) bool {
	for _, v := range vs {
		if NotContains(s, v) {
			return false
		}
	}
	return true
}

// NotContains reports whether v is not present in s.
func NotContains[E comparable](s []E, v E) bool {
	return !slices.Contains(s, v)
}

// NotContainsFunc reports whether v is not present in s.
func NotContainsFunc[E any](s []E, f func(E) bool) bool {
	return !slices.ContainsFunc(s, f)
}

// PadStart 如果slice长度小于 length 则在左侧填充val。
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

// PadEnd 如果slice长度小于 length 则在右侧填充val。
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

func Reduce[S ~[]E, E any, R any](s S, initValue R, f func(previousValue R, currentValue E, currentIndex int, s S) R) R {
	r := initValue
	for i, e := range s {
		r = f(r, e, i, s)
	}
	return r
}

// Reverse reverses the elements of the raw in place.
func Reverse[S ~[]E, E any](s S) S {
	slices.Reverse(s)
	return s
}

func ToSet[S ~[]E, M ~map[E]struct{}, E comparable](s S) M {
	r := make(M)
	for _, e := range s {
		r[e] = struct{}{}
	}
	return r
}

func SetAll[S ~[]E, E any](s S, f func(int) E) S {
	for i := range s {
		s[i] = f(i)
	}
	return s
}

func Shift[S ~[]E, E any](array S, offset int) {
	if len(array) <= 0 {
		return
	}
	Shifta(array, 0, len(array), offset)
}

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

// Swap swaps a series of elements in the given array.
// array the array to swap.
// offset1 the index of the first element in the series to swap.
// offset2 the index of the second element in the series to swap.
// length the number of elements to swap starting with the given indices.
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

// Shuffle 打乱数组顺序
func Shuffle[S ~[]E, E any](s S) S {
	r := randx.GetPCG()
	r.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	randx.PutPCG(r)
	return s
}

// 随机打乱直到没有相邻元素相同
func ShuffleNoAdjacent[S ~[]E, E comparable](s S) S {
	if len(s) <= 1 {
		return s
	}

	// 统计每个元素的出现次数
	count := make(map[E]int)
	for _, v := range s {
		count[v]++
	}

	// 检查是否可能实现无相邻重复
	maxCount := 0
	for _, c := range count {
		if c > maxCount {
			maxCount = c
		}
	}

	// 如果某个元素出现次数超过总长度的一半+1，则无法实现无相邻重复
	if maxCount > (len(s)+1)/2 {
		return s // 返回原数组或错误处理
	}

	// 使用更智能的排列算法
	return rearrangeNoAdjacent(s, count)
}

func rearrangeNoAdjacent[S ~[]E, E comparable](s S, count map[E]int) S {
	result := make(S, len(s))

	// 按出现次数排序元素
	type elemCount struct {
		elem  E
		count int
	}

	elems := make([]elemCount, 0, len(count))
	for e, c := range count {
		elems = append(elems, elemCount{e, c})
	}

	// 按计数降序排序
	sort.Slice(elems, func(i, j int) bool {
		return elems[i].count > elems[j].count
	})

	// 先填充偶数位置，再填充奇数位置
	index := 0
	for _, ec := range elems {
		for i := 0; i < ec.count; i++ {
			result[index] = ec.elem
			index += 2
			if index >= len(result) {
				index = 1 // 切换到奇数位置
			}
		}
	}

	return result
}

// 检查是否有相邻元素相同
func HasAdjacentDuplicates[S ~[]E, E comparable](s S) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			return true
		}
	}
	return false
}

// ToMap 方法创建一个Map，这个Map由原数组中的每个元素都调用一次提供的函数后的返回值作为Key、每个元素作为Value组成。
func ToMap[S ~[]E, M ~map[K]E, E any, K comparable](s S, f func(int, E) K) M {
	m := make(M, len(s))
	for i, v := range s {
		m[f(i, v)] = v
	}
	return m
}

// GroupBy 函数将输入切片中的元素按照指定函数分组，并返回一个Map，其中键是分组的依据，值是对应元素的列表。
func GroupBy[M ~map[K]S, S ~[]E, E any, K comparable](s S, f func(int, E) K) M {
	m := make(M, len(s))
	for i, v := range s {
		k := f(i, v)
		m[k] = append(m[k], v)
	}
	return m
}

// Uniq 用于去除切片中的重复元素并返回新切片。对于短切片，它通过逐个检查元素去重；对于长切片，使用Map提高效率。
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
