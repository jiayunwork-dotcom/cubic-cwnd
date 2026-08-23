package sim

import "fmt"

// Phase is the state-machine phase of a connection.
type Phase int

const (
	// PhaseSlowStart doubles the window every RTT until ssthresh.
	PhaseSlowStart Phase = iota
	// PhaseCongAvoid grows the window per the congestion-avoidance law.
	PhaseCongAvoid
)

// String renders the phase name.
func (p Phase) String() string {
	if p == PhaseCongAvoid {
		return "cong-avoid"
	}
	return "slow-start"
}

// ParsePhase converts a user string into a Phase.
func ParsePhase(s string) (Phase, error) {
	switch s {
	case "slow-start", "ss", "slowstart":
		return PhaseSlowStart, nil
	case "cong-avoid", "ca", "congavoid":
		return PhaseCongAvoid, nil
	default:
		return PhaseSlowStart, fmt.Errorf("sim: unknown phase %q", s)
	}
}

// transition decides whether slow start hands over to congestion
// avoidance at a given window. The transition happens the moment the
// doubled window reaches ssthresh.
func transition(cwnd, ssthresh float64) Phase {
	if cwnd < ssthresh {
		return PhaseSlowStart
	}
	return PhaseCongAvoid
}
