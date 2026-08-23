// Package cubic implements the TCP CUBIC congestion-window kernel:
// the W(t)=C(t-K)^3+WMax curve, the epoch K=cbrt(WMax*Beta/C), the
// RFC 8312 TCP-friendly estimate, and the fast-convergence rule.
//
// All windows are expressed in MSS segments, matching how CUBIC is
// described in the RFC. The constants Beta=0.7 and DefaultC=0.4 are
// pinned project-wide.
package cubic

// Beta is the multiplicative decrease factor after a loss event.
// RFC 8312 pins beta=0.7; the window is cut to Beta*WMax and the
// cubic epoch is derived from it.
const Beta = 0.7

// DefaultC is the default CUBIC scaling factor from RFC 8312 (C=0.4).
const DefaultC = 0.4

// MinWindow is the floor applied to the effective window, in MSS
// segments. A congestion window below one full segment has no meaning.
const MinWindow = 1.0

// Params describes a single point evaluation of the CUBIC model:
// the pre-loss window, the scaling factor, the RTT and the elapsed
// time since the loss.
type Params struct {
	// WMax is the window just before the loss, in MSS segments.
	WMax float64
	// C is the CUBIC scaling factor.
	C float64
	// RTT is the round trip time in seconds.
	RTT float64
	// T is the elapsed time since the loss, in seconds.
	T float64
}

// NewParams builds a Params value, defaulting C when it is not positive.
func NewParams(wMax, c, rtt, t float64) Params {
	if c <= 0 {
		c = DefaultC
	}
	return Params{WMax: wMax, C: c, RTT: rtt, T: t}
}
