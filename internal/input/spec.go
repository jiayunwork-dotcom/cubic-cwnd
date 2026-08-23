// Package input loads and validates the JSON problem specs that drive
// the cubic-cwnd CLI. It performs strict decoding (unknown fields are
// rejected) and cross-field checks such as "exactly one of t_seconds
// or acks" and positive window/RTT/C.
package input

import "cubic-cwnd/internal/cubic"

// SimConfig controls the RTT-by-RTT simulation subcommand.
type SimConfig struct {
	// Mode is "cubic" or "reno".
	Mode string
	// Rounds is the number of RTTs to simulate.
	Rounds int
	// Start is "after-loss" or "fresh".
	Start string
	// InitialCwnd and Ssthresh are used by the fresh start kind.
	InitialCwnd float64
	Ssthresh    float64
}

// FairConfig controls the two-flow convergence subcommand.
type FairConfig struct {
	// Capacity is the shared bottleneck in MSS segments.
	Capacity float64
	// Rounds is the maximum number of RTTs.
	Rounds int
	// FlowA and FlowB are the starting windows.
	FlowA float64
	FlowB float64
}

// Spec is the fully validated top-level problem.
type Spec struct {
	// Name is an optional human-readable label.
	Name string
	// C is the CUBIC scaling factor (defaults to 0.4).
	C float64
	// WMax is the pre-loss window in MSS segments.
	WMax float64
	// PrevWMax is the previous loss event's window, if known.
	PrevWMax float64
	// RTT is the round trip time in seconds.
	RTT float64
	// T is the elapsed time since the loss in seconds. When derived
	// from an ACK count this is acks*rtt.
	T float64
	// Acks is the ACK count used to derive T (0 when T was given).
	Acks int64
	// HorizonSeconds bounds the curve sampling range.
	HorizonSeconds float64
	// Samples is the number of points in a curve.
	Samples int
	// Sim is the simulation configuration.
	Sim SimConfig
	// Fair is the fairness configuration.
	Fair FairConfig
}

// ToParams converts the spec into a cubic.Params for point evaluation.
func (s *Spec) ToParams() cubic.Params {
	return cubic.Params{WMax: s.WMax, C: s.C, RTT: s.RTT, T: s.T}
}
