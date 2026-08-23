package cubic

import "math"

// cbrt is a thin wrapper over the standard library cube root so callers
// outside the package never need to import math themselves.
func cbrt(x float64) float64 { return math.Cbrt(x) }

// maxFloat returns the larger of two numbers.
func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// maxFloat3 returns the largest of three numbers.
func maxFloat3(a, b, c float64) float64 {
	return maxFloat(a, maxFloat(b, c))
}

// approxEqual compares two numbers with a relative tolerance. Two zeros
// are always equal.
func approxEqual(a, b, rel float64) bool {
	if a == b {
		return true
	}
	d := math.Abs(a - b)
	m := math.Max(math.Abs(a), math.Abs(b))
	if m == 0 {
		return d == 0
	}
	return d/m <= rel
}

// round3 rounds a value to three decimals for display.
func round3(v float64) float64 {
	return math.Round(v*1000) / 1000
}
