package fair

import "fmt"

type Frame struct {
	Round int
	A     float64
	B     float64
	Total float64
	Loss  bool
}

type Result struct {
	Config        Config
	Frames        []Frame
	Converged     bool
	ConvergeRound int
}

const convergenceTolerance = 0.05

func Run(cfg Config) (*Result, error) {
	cfg = cfg.Defaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	a := NewFlow("A", cfg.FlowA)
	b := NewFlow("B", cfg.FlowB)
	res := &Result{Config: cfg}
	for i := 1; i <= cfg.Rounds; i++ {
		a.Increment()
		b.Increment()
		loss := a.Cwnd+b.Cwnd > cfg.Capacity
		if loss {
			a.Decrement()
			b.Decrement()
		}
		res.Frames = append(res.Frames, Frame{
			Round: i,
			A:     a.Cwnd,
			B:     b.Cwnd,
			Total: a.Cwnd + b.Cwnd,
			Loss:  loss,
		})
		if diff(a.Cwnd, b.Cwnd) <= convergenceTolerance {
			res.Converged = true
			res.ConvergeRound = i
			break
		}
	}
	return res, nil
}

func diff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

func (res *Result) String() string {
	if len(res.Frames) == 0 {
		return "fair: no rounds simulated"
	}
	last := res.Frames[len(res.Frames)-1]
	if res.Converged {
		return fmt.Sprintf("fair: converged at round %d (A=%.3f B=%.3f)", res.ConvergeRound, last.A, last.B)
	}
	return fmt.Sprintf("fair: not converged within %d rounds (A=%.3f B=%.3f)", len(res.Frames), last.A, last.B)
}
