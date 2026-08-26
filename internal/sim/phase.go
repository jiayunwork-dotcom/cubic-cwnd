package sim

import "fmt"

type Phase int

const (
	PhaseSlowStart Phase = iota
	PhaseCongAvoid
)

func (p Phase) String() string {
	if p == PhaseCongAvoid {
		return "cong-avoid"
	}
	return "slow-start"
}

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

func transition(cwnd, ssthresh float64) Phase {
	if cwnd < ssthresh {
		return PhaseSlowStart
	}
	return PhaseCongAvoid
}
