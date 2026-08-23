// Package sim runs RTT-by-RTT congestion-window traces with explicit
// slow-start / congestion-avoidance phase transitions. Two congestion
// laws are supported: CUBIC (the cubic curve plus the TCP-friendly
// floor) and a pure Reno AIMD law that increments by exactly one
// segment per RTT in congestion avoidance.
package sim

import "fmt"

// Mode selects the congestion-avoidance law.
type Mode int

const (
	// ModeCubic uses the CUBIC curve with the RFC 8312 friendly floor.
	ModeCubic Mode = iota
	// ModeReno uses pure additive increase: +1 per RTT.
	ModeReno
)

// String renders the mode name.
func (m Mode) String() string {
	if m == ModeReno {
		return "reno"
	}
	return "cubic"
}

// ParseMode converts a user string into a Mode.
func ParseMode(s string) (Mode, error) {
	switch s {
	case "", "cubic":
		return ModeCubic, nil
	case "reno":
		return ModeReno, nil
	default:
		return ModeCubic, fmt.Errorf("sim: unknown mode %q (want cubic or reno)", s)
	}
}
