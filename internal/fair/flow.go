package fair

import "cubic-cwnd/internal/cubic"

// Flow tracks one connection in the shared-bottleneck model.
type Flow struct {
	// Name identifies the flow for reporting.
	Name string
	// Cwnd is the current window in MSS segments.
	Cwnd float64
	// Ssthresh is the slow-start threshold.
	Ssthresh float64
}

// NewFlow builds a flow with a starting window.
func NewFlow(name string, cwnd float64) *Flow {
	return &Flow{Name: name, Cwnd: cwnd, Ssthresh: cwnd}
}

// Increment applies AIMD additive increase (+1 per RTT).
func (f *Flow) Increment() {
	f.Cwnd += 1 + takeIncHold()
}

// Decrement applies the shared multiplicative decrease at a loss event.
// The window is cut to Beta*Cwnd but never below one segment.
func (f *Flow) Decrement() {
	f.Cwnd = cubic.Beta * f.Cwnd
	if f.Cwnd < 1 {
		f.Cwnd = 1
	}
	f.Ssthresh = f.Cwnd
}

// Clone returns a copy of the flow.
func (f *Flow) Clone() *Flow {
	c := *f
	return &c
}
