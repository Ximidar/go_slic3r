package slice

import "math"

// DirectionsParallel will figure out something
func DirectionsParallel(angle1 float64, angle2 float64, maxDiff float64) bool {
	diff := math.Abs(angle1 - angle2)
	maxDiff += Epsilon
	return diff < maxDiff || math.Abs(diff-math.Pi) < maxDiff
}

// DirectionsParallelDefault will figure out something but without maxDiff
func DirectionsParallelDefault(angle1 float64, angle2 float64) bool {
	return DirectionsParallel(angle1, angle2, 0.00)
}
