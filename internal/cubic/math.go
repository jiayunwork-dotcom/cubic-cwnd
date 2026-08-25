package cubic

import "math"

func cbrt(x float64) float64 { return math.Cbrt(x) }

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func maxFloat3(a, b, c float64) float64 {
	return maxFloat(a, maxFloat(b, c))
}

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

func round3(v float64) float64 {
	return math.Round(v*1000) / 1000
}
