package check

import (
	"cubic-cwnd/internal/cubic"
	"cubic-cwnd/internal/input"
)

// RunAll evaluates every rule against the spec, in a stable order.
func RunAll(spec *input.Spec) []Result {
	p := spec.ToParams()
	return []Result{
		tKRule(p),
		wMaxIncreaseRule(p),
		cDoubleRule(p),
		lowBDPFriendlyRule(p),
		afterLossRecoveryRule(p),
		fastConvergenceRule(spec),
		renoPlusOneRule(spec),
	}
}

// AllPass reports whether every rule passed.
func AllPass(rs []Result) bool {
	for _, r := range rs {
		if !r.Pass {
			return false
		}
	}
	return true
}

// Verify is a convenience wrapper: RunAll followed by AllPass.
func Verify(spec *input.Spec) ([]Result, bool) {
	rs := RunAll(spec)
	return rs, AllPass(rs)
}

// defaultParams is the reference point shared by the geometric rules.
func defaultParams() cubic.Params {
	return cubic.Params{WMax: 16, C: cubic.DefaultC, RTT: 0.1, T: 0.4}
}
