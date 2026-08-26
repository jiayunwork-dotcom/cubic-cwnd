package cubic

import "fmt"

type Result struct {
	Params   Params
	WMaxRef  float64
	K        float64
	WCubic   float64
	WEst     float64
	W        float64
	Slope    float64
	Friendly bool
	Status   Status
	FastConv bool
}

func Compute(p Params) (Result, error) {
	if err := p.Validate(); err != nil {
		return Result{}, err
	}
	return ComputeNoValidate(p), nil
}

func ComputeNoValidate(p Params) Result {
	return Result{
		Params:   p,
		WMaxRef:  p.WMax,
		K:        K(p),
		WCubic:   WCubic(p),
		WEst:     WEst(p),
		W:        WEffective(p),
		Slope:    SlopeAt(p),
		Friendly: IsTCPFriendly(p),
		Status:   WindowStatus(p),
	}
}

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

func (r Result) String() string {
	return fmt.Sprintf("W=%.3f K=%.4f friendly=%v status=%s fastconv=%v",
		r.W, r.K, r.Friendly, r.Status, r.FastConv)
}
