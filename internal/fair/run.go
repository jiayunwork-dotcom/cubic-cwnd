package fair

import "fmt"

// Frame is the pair state at one round of the model.
type Frame struct {
	// Round is the completed RTT number.
	Round int
	// A and B are the two flow windows.
	A float64
	B float64
	// Total is A+B after the round.
	Total float64
	// Loss reports whether the round ended with a loss event.
	Loss bool
}

// Result is a complete fairness run.
type Result struct {
	// Config is the validated configuration.
	Config Config
	// Frames holds one snapshot per completed RTT.
	Frames []Frame
	// Converged reports whether the two windows converged within the
	// run.
	Converged bool
	// ConvergeRound is the round where convergence was first detected.
	ConvergeRound int
}

// convergenceTolerance is the window difference (in segments) below
// which the two flows count as converged.
const convergenceTolerance = 0.05

// Run executes the two-flow model. Each RTT both flows increment;
// when the combined window exceeds the capacity a synchronized loss
// cuts both by Beta. Iteration is capped by MaxRounds.
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

// diff returns the absolute difference of two numbers.
func diff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

// String renders the convergence verdict.
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
