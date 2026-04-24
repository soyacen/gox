package osx

import (
	"unsafe"
)

// WordSize returns the word size of the operating system in bits.
//
// Returns:
//   - int: The word size in bits.
func WordSize() int {
	// 获取指针的大小，即计算机的字长
	wordSize := unsafe.Sizeof(new(interface{}))
	// 根据字节大小转换为位数
	return int(wordSize * 8)
}
