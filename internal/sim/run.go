package sim

import "fmt"

type Config struct {
	Mode        Mode
	Start       StartKind
	Rounds      int
	InitialCwnd float64
	Ssthresh    float64
	WMax        float64
	C           float64
	RTT         float64
}

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

type Result struct {
	Config Config
	States []*State
}

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
