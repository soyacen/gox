package slicex_test

import (
	"fmt"
	"github.com/soyacen/gox/slicex"
	"testing"
)

func TestSlice(t *testing.T) {
	// 创建一个初始切片
	s := &slicex.RWMutexSlice[[]int, int]{}

	// 添加元素并检查长度
	s = s.Append(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("Expected length 3, got %d", s.Len())
	}

	// 获取指定索引处的元素
	elem := s.Index(1)
	if elem != 2 {
		t.Errorf("Expected element 2, got %d", elem)
	}

	// 遍历切片并检查元素值
	var sum int
	s.Range(func(index int, elem int) bool {
		sum += elem
		return true
	})
	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}

	// 切片操作并检查结果切片
	ns := s.Slice(1, 3)
	if ns.Len() != 2 {
		t.Errorf("Expected length 2, got %d", ns.Len())
	}
	if ns.Index(0) != 2 || ns.Index(1) != 3 {
		t.Errorf("Expected elements [2, 3], got %v", ns.Unwrap())
	}

	// 在切片开头插入元素并检查结果切片
	ns = s.Prepend(0)
	if ns.Len() != 4 {
		t.Errorf("Expected length 4, got %d", ns.Len())
	}
	if ns.Index(0) != 0 || ns.Index(1) != 1 {
		t.Errorf("Expected elements [0, 1], got %v", ns.Unwrap())
	}

	// 在切片末尾插入元素并检查结果切片
	ns = s.Append(4)
	if ns.Len() != 5 {
		t.Errorf("Expected length 4, got %d", ns.Len())
	}
	if ns.Index(4) != 4 {
		t.Errorf("Expected element 4, got %v", ns.Unwrap())
	}
}

func TestWrapSlice(t *testing.T) {
	// 创建一个初始切片
	slice := []int{1, 2, 3}

	// 使用WrapSlice函数创建包装切片
	s := slicex.WrapSlice[[]int](slice)

	// 检查切片长度和元素值
	if s.Len() != 3 {
		t.Errorf("Expected length 3, got %d", s.Len())
	}
	if s.Index(0) != 1 || s.Index(1) != 2 || s.Index(2) != 3 {
		t.Errorf("Expected elements [1, 2, 3], got %v", s.Unwrap())
	}

	// 修改原始切片并检查包装切片是否受影响
	slice[0] = 0
	if s.Index(0) != 0 {
		t.Errorf("Expected element 0, got %v", s.Unwrap())
	}
}

func ExampleSlice() {
	s := &slicex.RWMutexSlice[[]int, int]{}
	fmt.Println("len:", s.Len())
	fmt.Println("cap:", s.Cap())
	s = s.Append(1, 2, 3, 4, 5)
	fmt.Println("len:", s.Len())
	fmt.Println("cap:", s.Cap())
	fmt.Println(s.Unwrap())
	s = s.Prepend(0, 9, 8, 7, 6)
	fmt.Println("len:", s.Len())
	fmt.Println("cap:", s.Cap())
	fmt.Println(s.Unwrap())
	s.Range(func(index int, elem int) bool {
		fmt.Println(index, elem)
		return true
	})
	fmt.Println(s.Index(1))

	s1 := s.Slice(1, 7)
	fmt.Println("s1 len:", s1.Len())
	fmt.Println("s1 cap:", s1.Cap())
	fmt.Println(s1.Unwrap())

	s2 := s.Slice(1, 7, 9)
	fmt.Println("s2 len:", s2.Len())
	fmt.Println("s2 cap:", s2.Cap())
	fmt.Println(s2.Unwrap())
}
