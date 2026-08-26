package sim

import (
	"fmt"

	"cubic-cwnd/internal/cubic"
)

type StartKind int

const (
	StartAfterLoss StartKind = iota
	StartFresh
)

func (s StartKind) String() string {
	if s == StartFresh {
		return "fresh"
	}
	return "after-loss"
}

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

type State struct {
	Mode     Mode
	Phase    Phase
	Cwnd     float64
	Ssthresh float64
	WMax     float64
	C        float64
	RTT      float64
	Round    int
	Elapsed  float64
}

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

func (s *State) Clone() *State {
	c := *s
	return &c
}

func (s *State) String() string {
	return fmt.Sprintf("round=%d phase=%s cwnd=%.3f ssthresh=%.3f",
		s.Round, s.Phase, s.Cwnd, s.Ssthresh)
}
