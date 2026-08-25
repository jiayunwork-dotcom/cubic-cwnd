package fair

import "fmt"

type Mode int

const (
	ModeAIMD Mode = iota
)

func (m Mode) String() string { return "aimd" }

type Config struct {
	Capacity float64
	Rounds   int
	FlowA    float64
	FlowB    float64
}

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

const MaxRounds = 100000
