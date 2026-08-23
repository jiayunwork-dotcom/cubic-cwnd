package sim

import (
	"fmt"

	"cubic-cwnd/internal/cubic"
)

// StartKind selects how the simulation begins.
type StartKind int

const (
	// StartAfterLoss starts right after a loss event: the window is
	// cut to Beta*WMax and ssthresh follows the same decrease.
	StartAfterLoss StartKind = iota
	// StartFresh starts from an explicit initial window and ssthresh,
	// giving slow start room to run before congestion avoidance.
	StartFresh
)

// String renders the start kind.
func (s StartKind) String() string {
	if s == StartFresh {
		return "fresh"
	}
	return "after-loss"
}

// ParseStart converts a user string into a StartKind.
func ParseStart(s string) (StartKind, error) {
	switch s {
	case "", "after-loss":
		return StartAfterLoss, nil
	case "fresh":
		return StartFresh, nil
	default:
		return StartAfterLoss, fmt.Errorf("sim: unknown start %q (want after-loss or fresh)", s)
	}
}

// State is a single connection's window state at one point in time.
type State struct {
	// Mode is the congestion-avoidance law.
	Mode Mode
	// Phase is slow start or congestion avoidance.
	Phase Phase
	// Cwnd is the current window in MSS segments.
	Cwnd float64
	// Ssthresh is the slow-start threshold in MSS segments.
	Ssthresh float64
	// WMax is the reference pre-loss window for the cubic curve.
	WMax float64
	// C is the cubic scaling factor.
	C float64
	// RTT is the round trip time in seconds.
	RTT float64
	// Round is the number of completed RTTs.
	Round int
	// Elapsed is the time in seconds since the start of the trace.
	Elapsed float64
}

// NewAfterLoss builds a state right after a loss event.
func NewAfterLoss(mode Mode, wMax, c, rtt float64) *State {
	s := &State{
		Mode: mode,
		WMax: wMax,
		C:    c,
		RTT:  rtt,
	}
	s.Cwnd = cubic.Beta * wMax
	s.Ssthresh = minSsthresh(s.Cwnd)
	s.Phase = transition(s.Cwnd, s.Ssthresh)
	return s
}

// NewFresh builds a state that starts from an explicit window.
func NewFresh(mode Mode, wMax, c, rtt, initialCwnd, ssthresh float64) *State {
	s := &State{
		Mode:     mode,
		WMax:     wMax,
		C:        c,
		RTT:      rtt,
		Cwnd:     initialCwnd,
		Ssthresh: ssthresh,
	}
	s.Phase = transition(s.Cwnd, s.Ssthresh)
	return s
}

// Clone returns a copy of the state.
func (s *State) Clone() *State {
	c := *s
	return &c
}

// String renders a compact summary of the state.
func (s *State) String() string {
	return fmt.Sprintf("round=%d phase=%s cwnd=%.3f ssthresh=%.3f",
		s.Round, s.Phase, s.Cwnd, s.Ssthresh)
}
