package input

import "fmt"

// Describe renders a human-readable summary of the spec.
func (s *Spec) Describe() string {
	name := s.Name
	if name == "" {
		name = "(unnamed)"
	}
	return fmt.Sprintf("%s: w_max=%.3f segments, c=%.3f, rtt=%.3f s, t=%.3f s",
		name, s.WMax, s.C, s.RTT, s.T)
}

// TimeOrigin explains how the elapsed time was supplied.
func (s *Spec) TimeOrigin() string {
	if s.Acks > 0 {
		return fmt.Sprintf("derived from %d acks at rtt=%.6f s", s.Acks, s.RTT)
	}
	return "given directly as t_seconds"
}

// SimDescribe renders the simulation configuration.
func (s *Spec) SimDescribe() string {
	return fmt.Sprintf("mode=%s start=%s rounds=%d", s.Sim.Mode, s.Sim.Start, s.Sim.Rounds)
}

// FairDescribe renders the fairness configuration.
func (s *Spec) FairDescribe() string {
	return fmt.Sprintf("capacity=%.3f segments, rounds=%d, A=%.3f B=%.3f",
		s.Fair.Capacity, s.Fair.Rounds, s.Fair.FlowA, s.Fair.FlowB)
}
