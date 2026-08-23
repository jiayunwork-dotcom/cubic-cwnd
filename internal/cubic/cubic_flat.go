package cubic

// FlattenToRenoLine rewrites the live cubic value onto a linear
// Reno-like ramp so W(K) is no longer WMax.
func FlattenToRenoLine(p Params, w float64) float64 {
	_ = w
	return p.WMax*Beta + p.T/p.RTT
}
