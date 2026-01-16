package slicex_test

import (
	"github.com/soyacen/gox/mathx/randx"
	"github.com/soyacen/gox/slicex"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	type args struct {
		ss [][]int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "empty slices",
			args: args{ss: [][]int{}},
			want: []int{},
		},
		{
			name: "single raw",
			args: args{ss: [][]int{{1, 2, 3}}},
			want: []int{1, 2, 3},
		},
		{
			name: "multiple slices",
			args: args{ss: [][]int{{1, 2}, {3, 4}, {5}}},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "mixed lengths",
			args: args{ss: [][]int{{1}, {2, 3, 4}, {5, 6}}},
			want: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slicex.Merge(tt.args.ss...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge(%v) = %v, want %v", tt.args.ss, got, tt.want)
			}
		})
	}
}

func TestUniq(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		f        func(interface{}) interface{}
	}{
		{
			name:     "empty raw",
			input:    []int{},
			expected: []int{},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
		{
			name:     "raw with unique elements",
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
		{
			name:     "raw with duplicate elements",
			input:    []int{1, 2, 2, 3, 3, 3},
			expected: []int{1, 2, 3},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
		{
			name:     "string raw with duplicate elements",
			input:    []string{"a", "b", "b", "c", "c", "c"},
			expected: []string{"a", "b", "c"},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]string))
			},
		},
		{
			name:     "rune raw with duplicate elements",
			input:    []rune{'a', 'b', 'b', 'f', 'f', 'c', 'c', 'c', 'd', 'd', 'e', 'e'},
			expected: []rune{'a', 'b', 'f', 'c', 'd', 'e'},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]rune))
			},
		},

		{
			name:     "Test with duplicate elements",
			input:    []int{1, 2, 2, 3, 4, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
		{
			name:     "Test with empty raw",
			input:    []int{},
			expected: []int{},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
		{
			name:     "Test with single element",
			input:    []int{42},
			expected: []int{42},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
		{
			name:     "Test with multiple types of elements",
			input:    []int{1, 1, 2, 0, 0, -1, -1},
			expected: []int{1, 2, 0, -1},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},

		{
			name:     "Integers",
			input:    []int{1, 2, 2, 3, 4, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
		{
			name:     "Strings",
			input:    []string{"a", "b", "b", "c", "d", "d", "e"},
			expected: []string{"a", "b", "c", "d", "e"},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]string))
			},
		},
		{
			name:     "Empty raw",
			input:    []int{},
			expected: []int{},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
		{
			name:     "Nil raw",
			input:    ([]int)(nil),
			expected: ([]int)(nil),
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
		{
			name:     "Repeating elements",
			input:    []int{1, 1, 1, 1},
			expected: []int{1},
			f: func(i interface{}) interface{} {
				return slicex.Uniq(i.([]int))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.f(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Uniq() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// BenchmarkUniqV1-8   	110002	     10039 ns/op
//
//	114277	     10083 ns/op

// BenchmarkUniqV1-8   	  117631	     10064 ns/op
// BenchmarkUniqV1-8   	  103548	     10020 ns/op
// BenchmarkUniqV1-8   	  102814	     10098 ns/op
// BenchmarkUniqV1-8   	  100611	     11081 ns/op
func BenchmarkUniqV1(b *testing.B) {
	var arr []int
	for i := 0; i < 128; i++ {
		arr = append(arr, i)
	}
	for i := 0; i < b.N; i++ {
		slices.Contains(arr, randx.Intn(100000000))
	}
}

// BenchmarkUniqV2-8   	  119151	     10068 ns/op
// 112916	             10048 ns/op
// BenchmarkUniqV2-8   	  104806	     10059 ns/op
// BenchmarkUniqV2-8   	  119354	     10077 ns/op
// BenchmarkUniqV2-8   	  110120	     10404 ns/op
// enchmarkUniqV2-8   	  100659	     10100 ns/op
func BenchmarkUniqV2(b *testing.B) {
	var arr = map[int]struct{}{}
	for i := 0; i < 180; i++ {
		arr[i] = struct{}{}
	}
	for i := 0; i < b.N; i++ {
		_, _ = arr[randx.Intn(100000000)]
	}
}

func TestRemoveAt(t *testing.T) {
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 0))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 1))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 2))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 3))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 4))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 5))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 6))
	// t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, -1))
	// t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 7))

	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 0, 1))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 1, 3))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 2, 5))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 3, 1))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 4, 0))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 5, 2))
	t.Log(slicex.RemoveAt([]int{0, 1, 2, 3, 4, 5, 6}, 6, 4))
}

func TestRemove(t *testing.T) {
	ints := slicex.Remove([]int{1, 2, 3, 4, 5}, 2)
	t.Log(ints)
}

func TestReverse(t *testing.T) {
	assert.Equal(t, []int{5, 4, 3, 2, 1}, slicex.Reverse([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, []string{"5", "4", "3", "2", "1"}, slicex.Reverse([]string{"1", "2", "3", "4", "5"}))
}

func TestToSet(t *testing.T) {
	strings := []string{"1", "2", "3"}
	set := slicex.ToSet[[]string, map[string]struct{}](strings)
	t.Log(set)
}

func TestConcat(t *testing.T) {
	t.Log(slicex.Concat([]int{1, 2, 3, 4, 5}, []int{6, 7, 8, 9, 0}))
	t.Log(slicex.Concat([]float32{1, 2, 3, 4, 5}, []float32{6, 7, 8, 9, 0}))
	t.Log(slicex.Concat([]string{"1", "2", "3", "4", "5"}, []string{"6", "7", "8", "9", "0"}))
}

func TestChunk(t *testing.T) {
	tests := []struct {
		array          []int
		size           int
		expectedChunks [][]int
	}{
		{
			array:          []int{},
			size:           2,
			expectedChunks: [][]int{},
		},
		{
			array:          []int{0},
			size:           2,
			expectedChunks: [][]int{{0}},
		},
		{
			array:          []int{0, 1},
			size:           2,
			expectedChunks: [][]int{{0, 1}},
		},
		{
			array:          []int{0, 1, 2},
			size:           2,
			expectedChunks: [][]int{{0, 1}, {2}},
		},
		{
			array:          []int{0, 1, 2, 3},
			size:           2,
			expectedChunks: [][]int{{0, 1}, {2, 3}},
		},
		{
			array:          []int{0, 1, 2, 3, 4},
			size:           2,
			expectedChunks: [][]int{{0, 1}, {2, 3}, {4}},
		},
		{
			array:          []int{0, 1, 2, 3, 4, 5},
			size:           2,
			expectedChunks: [][]int{{0, 1}, {2, 3}, {4, 5}},
		},
	}
	for _, test := range tests {
		chunks := slicex.Chunk(test.array, test.size)
		if len(chunks) != len(test.expectedChunks) {
			t.Fatalf("%v expected chunks is %v, but is %v", test.array, test.expectedChunks, chunks)
		}
		for i := range chunks {
			if len(chunks[i]) != len(test.expectedChunks[i]) {
				t.Fatalf("%v expected chunks is %v, but is %v", test.array, test.expectedChunks, chunks)
			}
			for j := range chunks[i] {
				if chunks[i][j] != test.expectedChunks[i][j] {
					t.Fatalf("%v expected chunks is %v, but is %v", test.array, test.expectedChunks, chunks)
				}
			}
		}
		t.Log(chunks)
	}

}
