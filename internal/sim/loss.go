package sim

import "cubic-cwnd/internal/cubic"

func (s *State) Decrease() float64 {
	s.Cwnd = cubic.Beta * s.Cwnd
	s.Ssthresh = minSsthresh(s.Cwnd)
	s.Phase = transition(s.Cwnd, s.Ssthresh)
	return s.Ssthresh
}

func (s *State) Loss() *State {
	s.Decrease()
	return s
}
