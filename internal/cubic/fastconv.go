package cubic

func FastConvergenceFactor() float64 {
	return (1 + Beta) / 2
}

func FastConvergence(prevWMax, curWMax float64) (float64, bool) {
	if prevWMax > curWMax {
		return curWMax * FastConvergenceFactor(), true
	}
	return curWMax, false
}

func ApplyFastConvergence(prevWMax, curWMax float64) float64 {
	adj, _ := FastConvergence(prevWMax, curWMax)
	return adj
}

func InFastConvergence(prevWMax, curWMax float64) bool {
	_, fired := FastConvergence(prevWMax, curWMax)
	return fired
}
