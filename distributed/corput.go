package distributed

// VanDerCorputSequence generates a van der Corput sequence.
// The sequence is a low-discrepancy sequence useful for generating evenly distributed numbers.
//
// Parameters:
//   - n: the position in the sequence, starting from 0.
//   - base: the radix of the sequence, determining its fractal nature.
//
// Returns:
//   - float64: the generated van der Corput sequence value.
func VanDerCorputSequence(n int, base int) float64 {
	// 初始化q为0，用于累加每个位的倒数权重。
	q := 0.0
	// 初始化bk为基数的倒数，用于计算每一位的权重。
	bk := 1.0 / float64(base)
	// 循环处理每一位直到n为0。
	for n > 0 {
		// 将当前位的值乘以权重bk加到q上。
		q += float64(n%base) * bk
		// 移动到下一位，更新n和bk。
		n /= base
		bk /= float64(base)
	}
	// 返回计算得到的范德科鲁普序列值。
	return q
}

// BinaryVanDerCorputSequence generates a binary van der Corput sequence.
// This function is a special case of the van der Corput sequence for base 2.
// The van der Corput sequence is a low-discrepancy sequence commonly used in
// random number generation and Monte Carlo simulations.
// In the binary case, this sequence is especially suitable for scenarios
// requiring evenly distributed binary digits.
//
// Parameters:
//   - n: the index of the nth number in the sequence.
//
// Returns:
//   - float64: the value of the nth van der Corput sequence number.
func BinaryVanDerCorputSequence(n int) float64 {
	return VanDerCorputSequence(n, 2)
}
