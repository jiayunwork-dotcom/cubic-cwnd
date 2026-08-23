package sim

import "cubic-cwnd/internal/cubic"

// Step advances the state by one RTT.
func (s *State) Step() {
	s.Round++
	s.Elapsed = float64(s.Round) * s.RTT
	switch s.Phase {
	case PhaseSlowStart:
		s.stepSlowStart()
	case PhaseCongAvoid:
		s.stepCongAvoid()
	}
}

// stepSlowStart doubles the window every RTT. The moment the doubled
// window reaches ssthresh, the connection hands over to congestion
// avoidance with the window clamped at ssthresh.
func (s *State) stepSlowStart() {
	doubled := s.Cwnd * 2
	if doubled >= s.Ssthresh {
		s.Cwnd = s.Ssthresh
		s.Phase = PhaseCongAvoid
		return
	}
	s.Cwnd = doubled
}

// stepCongAvoid grows the window according to the selected law:
//   - Reno: exactly +1 segment per RTT.
//   - CUBIC: the effective window W(t)=max(W_cubic, W_est) at the
//     elapsed time since the trace start.
func (s *State) stepCongAvoid() {
	switch s.Mode {
	case ModeReno:
		s.Cwnd += 1.0
	case ModeCubic:
		p := cubic.Params{WMax: s.WMax, C: s.C, RTT: s.RTT, T: s.Elapsed}
		s.Cwnd = cubic.WEffective(p)
	}
}
