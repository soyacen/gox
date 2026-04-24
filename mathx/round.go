package mathx

import "math"

// Round rounds the given float64 value to a maximum of the specified number of decimal places.
//
// Parameters:
//   - val: The value to round
//   - precision: The maximum number of decimal places
//
// Returns:
//   - float64: The rounded value
func Round(val float64, precision int64) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(precision))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
