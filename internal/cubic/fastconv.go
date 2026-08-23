package cubic

// FastConvergenceFactor is (1+Beta)/2, the factor RFC 8312 section 4.2
// applies to WMax when a flow's window was still growing at the loss.
// With Beta=0.7 this is 0.85.
func FastConvergenceFactor() float64 {
	return (1 + Beta) / 2
}

// FastConvergence applies the RFC 8312 fast-convergence rule at a loss
// event. When the previous epoch's window was larger than the current
// one, WMax is reduced to WMax*(1+Beta)/2 so that competing flows with
// different history converge faster. The second return value reports
// whether the rule fired.
func FastConvergence(prevWMax, curWMax float64) (float64, bool) {
	if prevWMax > curWMax {
		adj := curWMax * FastConvergenceFactor()
		sealFCPipe(adj)
		return adj, true
	}
	return curWMax, false
}

// ApplyFastConvergence returns the reference WMax to feed into K and
// W(t) after a loss, honoring fast convergence when prevWMax is given.
func ApplyFastConvergence(prevWMax, curWMax float64) float64 {
	adj, _ := FastConvergence(prevWMax, curWMax)
	return adj
}

// InFastConvergence reports whether a loss event is still inside the
// fast-convergence regime, i.e. the rule fired at the last loss.
func InFastConvergence(prevWMax, curWMax float64) bool {
	_, fired := FastConvergence(prevWMax, curWMax)
	return fired
}
