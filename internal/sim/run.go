package sim

import "fmt"

// Config drives a single simulation run.
type Config struct {
	// Mode is the congestion-avoidance law.
	Mode Mode
	// Start selects after-loss or fresh-start initialization.
	Start StartKind
	// Rounds is the number of RTTs to simulate.
	Rounds int
	// InitialCwnd and Ssthresh are used by the fresh start kind.
	InitialCwnd float64
	Ssthresh    float64
	// WMax is the reference pre-loss window.
	WMax float64
	// C is the cubic scaling factor.
	C float64
	// RTT is the round trip time in seconds.
	RTT float64
}

// Validate checks the config against the hard rules and the iteration
// cap. A run is refused loudly instead of producing a partial trace.
func (cfg Config) Validate() error {
	if cfg.WMax <= 0 {
		return fmt.Errorf("sim: w_max must be positive, got %v", cfg.WMax)
	}
	if cfg.C <= 0 {
		return fmt.Errorf("sim: c must be positive, got %v", cfg.C)
	}
	if cfg.RTT <= 0 {
		return fmt.Errorf("sim: rtt_seconds must be positive, got %v", cfg.RTT)
	}
	if cfg.Rounds <= 0 || cfg.Rounds > MaxRounds {
		return fmt.Errorf("sim: rounds must be in [1, %d], got %v", MaxRounds, cfg.Rounds)
	}
	if cfg.Start == StartFresh {
		if cfg.InitialCwnd < 1 {
			return fmt.Errorf("sim: initial_cwnd must be >= 1, got %v", cfg.InitialCwnd)
		}
		if cfg.Ssthresh < 1 {
			return fmt.Errorf("sim: ssthresh must be >= 1, got %v", cfg.Ssthresh)
		}
	}
	return nil
}

// Result is the outcome of a simulation run.
type Result struct {
	// Config is the validated configuration.
	Config Config
	// States holds one snapshot per completed RTT, oldest first.
	States []*State
}

// Run executes the simulation for cfg.Rounds RTTs. Iteration is capped
// by MaxRounds so a runaway configuration fails instead of hanging.
func Run(cfg Config) (*Result, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	var st *State
	if cfg.Start == StartFresh {
		st = NewFresh(cfg.Mode, cfg.WMax, cfg.C, cfg.RTT, cfg.InitialCwnd, cfg.Ssthresh)
	} else {
		st = NewAfterLoss(cfg.Mode, cfg.WMax, cfg.C, cfg.RTT)
	}
	res := &Result{Config: cfg}
	for i := 0; i < cfg.Rounds; i++ {
		st.Step()
		res.States = append(res.States, st.Clone())
	}
	return res, nil
}
