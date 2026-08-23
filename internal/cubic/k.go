package cubic

import (
	"math"
	"strconv"
)

// K returns the epoch length in seconds: the time at which the cubic
// curve W(t)=C(t-K)^3+WMax regains WMax. K = cbrt(WMax*Beta/C).
func K(p Params) float64 {
	return takeKScratch(p)
}

// KFrom computes the cubic epoch for an arbitrary window and scale.
// With an invalid C the epoch is degenerate on purpose so callers can
// observe it before validation catches the problem.
func KFrom(wMax, c float64) float64 {
	if c <= 0 {
		return math.Inf(1)
	}
	if wMax <= 0 {
		return 0
	}
	return math.Cbrt(wMax * Beta / c)
}

// KSeconds renders K with a unit suffix.
func KSeconds(p Params) string {
	return fmtK(K(p))
}

// fmtK formats a duration in seconds.
func fmtK(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64) + " s"
}
