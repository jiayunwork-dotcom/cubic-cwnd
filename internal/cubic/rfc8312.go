package cubic

// RFC8312 collects the defining identities and constants of the CUBIC
// kernel. Every evaluation must satisfy them simultaneously.
const (
	// EqWindow is the cubic window curve.
	EqWindow = "W(t) = C*(t-K)^3 + WMax"
	// EqSlope is the slope of the curve, zero at t=K.
	EqSlope = "dW/dt = 3*C*(t-K)^2, zero at t=K"
	// EqK is the epoch relation.
	EqK = "K = cbrt(WMax*Beta/C)"
	// EqFriendly is the TCP-friendly estimate.
	EqFriendly = "W_est = WMax*Beta + 3*(1-Beta)/(1+Beta)*t/RTT"
	// EqEffective is how the two curves are combined.
	EqEffective = "W(t) = max(W_cubic(t), W_est(t))"
)

// Identities returns the list of equations the kernel guarantees, in
// display order. Used by the CLI and documentation.
func Identities() []string {
	return []string{
		EqWindow,
		EqSlope,
		EqK,
		EqFriendly,
		EqEffective,
	}
}

// Coefficients returns the derived constants that depend only on Beta.
// g is the friendly-estimate slope per RTT; fc is the fast-convergence
// factor.
func Coefficients() (g, fc float64) {
	return FriendlyGain(), FastConvergenceFactor()
}
