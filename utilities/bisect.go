package utilities

import (
	"math"
)

// Bisect takes a known-good low and known-bad high value as the bounds
// to bisect, and a function to test each value for success or failure.
// If the function succeeds, the value is adjusted toward the maximum,
// and if the function fails, the value is adjusted toward the minimum.
// The final value is returned when the difference between the success
// and the failure is less than or equal to the acceptance threshold
// (usually 1, for integers).
func Bisect[T Number](low, high, threshold T, tryFunc func(val T) bool) T {
	for T(math.Abs(float64(high-low))) > threshold {
		currVal := low + ((high - low) / 2)
		success := tryFunc(currVal)
		if success {
			low = currVal
		} else {
			high = currVal
		}
	}

	return low
}
