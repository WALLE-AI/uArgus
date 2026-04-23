package anomaly

import "math"

// ZScore computes the Z-score: (value - mean) / stddev.
// Returns 0 if stddev is zero to avoid division by zero.
func ZScore(value, mean, stddev float64) float64 {
	if stddev == 0 {
		return 0
	}
	return (value - mean) / stddev
}

// SpikeDetect returns true if the value is more than threshold standard deviations from the mean.
func SpikeDetect(value, mean, stddev, threshold float64) bool {
	return math.Abs(ZScore(value, mean, stddev)) > threshold
}
