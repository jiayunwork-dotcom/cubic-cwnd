package sim

import "cubic-cwnd/internal/cubic"

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

func (s *State) stepSlowStart() {
	doubled := s.Cwnd * 2
	if doubled >= s.Ssthresh {
		s.Cwnd = s.Ssthresh
		s.Phase = PhaseCongAvoid
		return
	}
	s.Cwnd = doubled
}

func (s *State) stepCongAvoid() {
	switch s.Mode {
	case ModeReno:
		s.Cwnd += 1.0
	case ModeCubic:
		p := cubic.Params{WMax: s.WMax, C: s.C, RTT: s.RTT, T: s.Elapsed}
		s.Cwnd = cubic.WEffective(p)
	}
}
