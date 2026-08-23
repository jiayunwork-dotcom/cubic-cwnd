// Package fair simulates two flows sharing a single bottleneck to
// demonstrate AIMD fairness and convergence: with synchronized losses
// the two windows converge to the same value even from very different
// starting points.
package fair

import "fmt"

// Mode of the fairness run. The model uses one AIMD law; the type is
// kept so future laws can be added without breaking the API.
type Mode int

const (
	// ModeAIMD grows by +1 per RTT and cuts by Beta at a loss.
	ModeAIMD Mode = iota
)

// String renders the mode name.
func (m Mode) String() string { return "aimd" }

// Config drives a fairness run.
type Config struct {
	// Capacity is the shared bottleneck in MSS segments.
	Capacity float64
	// Rounds is the maximum number of RTTs to simulate.
	Rounds int
	// FlowA and FlowB are the starting windows.
	FlowA float64
	FlowB float64
}

// Defaults fills unset fields with the reference values. Only explicit
// zeros are defaulted; negative values survive so Validate rejects them.
func (cfg Config) Defaults() Config {
	if cfg.Capacity == 0 {
		cfg.Capacity = 60
	}
	if cfg.FlowA == 0 {
		cfg.FlowA = 45
	}
	if cfg.FlowB == 0 {
		cfg.FlowB = 15
	}
	if cfg.Rounds == 0 {
		cfg.Rounds = 120
	}
	return cfg
}

// Validate checks the config and the iteration cap.
func (cfg Config) Validate() error {
	if cfg.Capacity <= 0 {
		return fmt.Errorf("fair: capacity_segments must be positive, got %v", cfg.Capacity)
	}
	if cfg.FlowA < 1 {
		return fmt.Errorf("fair: flow_a_cwnd must be >= 1, got %v", cfg.FlowA)
	}
	if cfg.FlowB < 1 {
		return fmt.Errorf("fair: flow_b_cwnd must be >= 1, got %v", cfg.FlowB)
	}
	if cfg.Rounds <= 0 || cfg.Rounds > MaxRounds {
		return fmt.Errorf("fair: rounds must be in [1, %d], got %v", MaxRounds, cfg.Rounds)
	}
	return nil
}

// MaxRounds caps the fairness iteration.
const MaxRounds = 100000
