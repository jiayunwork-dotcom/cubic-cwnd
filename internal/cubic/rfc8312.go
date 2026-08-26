package cubic

const (
	EqWindow    = "W(t) = C*(t-K)^3 + WMax"
	EqSlope     = "dW/dt = 3*C*(t-K)^2, zero at t=K"
	EqK         = "K = cbrt(WMax*Beta/C)"
	EqFriendly  = "W_est = WMax*Beta + 3*(1-Beta)/(1+Beta)*t/RTT"
	EqEffective = "W(t) = max(W_cubic(t), W_est(t))"
)

func Identities() []string {
	return []string{
		EqWindow,
		EqSlope,
		EqK,
		EqFriendly,
		EqEffective,
	}
}

func Coefficients() (g, fc float64) {
	return FriendlyGain(), FastConvergenceFactor()
}
