package cubic

import "fmt"

// Result is the outcome of one point evaluation of the CUBIC model.
type Result struct {
	// Params is the validated input.
	Params Params
	// WMaxRef is the reference window after fast convergence, if any.
	WMaxRef float64
	// K is the cubic epoch in seconds.
	K float64
	// WCubic is the raw cubic-curve value.
	WCubic float64
	// WEst is the TCP-friendly estimate.
	WEst float64
	// W is the effective window used by the flow.
	W float64
	// Slope is dW_cubic/dt at the point.
	Slope float64
	// Friendly reports whether the effective window follows the
	// TCP-friendly branch.
	Friendly bool
	// Status is the recovery status relative to WMaxRef.
	Status Status
	// FastConv reports whether fast convergence fired at this loss.
	FastConv bool
}

// Compute validates the parameters and evaluates the model.
func Compute(p Params) (Result, error) {
	if err := p.Validate(); err != nil {
		return Result{}, err
	}
	return ComputeNoValidate(p), nil
}

// ComputeNoValidate evaluates the model assuming Params is valid.
func ComputeNoValidate(p Params) Result {
	return Result{
		Params:  p,
		WMaxRef: p.WMax,
		K:       K(p),
		WCubic:  WCubic(p),
		WEst:    WEst(p),
		W:       WEffective(p),
		Slope:   SlopeAt(p),
		Friendly: IsTCPFriendly(p),
		Status:  WindowStatus(p),
	}
}

// ComputeFastConv evaluates the model after applying the RFC 8312 fast
// convergence rule. When apply is true and prevWMax is positive, WMax
// is reduced by (1+Beta)/2 if the flow was growing.
func ComputeFastConv(p Params, prevWMax float64, apply bool) (Result, error) {
	if err := p.Validate(); err != nil {
		return Result{}, err
	}
	ref := p.WMax
	fired := false
	if apply && prevWMax > 0 {
		ref, fired = FastConvergence(prevWMax, p.WMax)
	}
	q := p
	q.WMax = ref
	r := ComputeNoValidate(q)
	r.FastConv = fired
	return r, nil
}

// String renders a compact one-line summary of the result.
func (r Result) String() string {
	return fmt.Sprintf("W=%.3f K=%.4f friendly=%v status=%s fastconv=%v",
		r.W, r.K, r.Friendly, r.Status, r.FastConv)
}
