package sim

import "cubic-cwnd/internal/cubic"

// Decrease applies the multiplicative decrease of a loss event: the
// window is cut to Beta*Cwnd and ssthresh follows, then the phase is
// re-derived (a window at or above ssthresh resumes congestion
// avoidance immediately). Returns the new ssthresh.
func (s *State) Decrease() float64 {
	s.Cwnd = cubic.Beta * s.Cwnd
	s.Ssthresh = minSsthresh(s.Cwnd)
	s.Phase = transition(s.Cwnd, s.Ssthresh)
	return s.Ssthresh
}

// Loss simulates one loss event on the state in place.
func (s *State) Loss() *State {
	s.Decrease()
	return s
}
