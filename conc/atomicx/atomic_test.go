package atomicx

import (
	"fmt"
	"testing"
)

func TestDecr(t *testing.T) {
	var val32 uint32 = 100
	var val64 uint64 = 100

	// 测试 SubUint32
	newVal32 := SubUint32(&val32, 10)
	fmt.Println("New value (uint32):", newVal32) // 应该输出 90

	// 测试 DecrUint32
	newVal32Decr := DecrUint32(&val32)
	fmt.Println("New value (uint32) after DecrUint32:", newVal32Decr) // 应该输出 89

	// 测试 SubUint64
	newVal64 := SubUint64(&val64, 10)
	fmt.Println("New value (uint64):", newVal64) // 应该输出 90

	// 测试 DecrUint64
	newVal64Decr := DecrUint64(&val64)
	fmt.Println("New value (uint64) after DecrUint64:", newVal64Decr) // 应该输出 89
}
