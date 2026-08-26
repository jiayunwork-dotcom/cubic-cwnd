package check

import (
	"cubic-cwnd/internal/cubic"
	"cubic-cwnd/internal/input"
)

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

func AllPass(rs []Result) bool {
	for _, r := range rs {
		if !r.Pass {
			return false
		}
	}
	return true
}

func Verify(spec *input.Spec) ([]Result, bool) {
	rs := RunAll(spec)
	return rs, AllPass(rs)
}

func defaultParams() cubic.Params {
	return cubic.Params{WMax: 16, C: cubic.DefaultC, RTT: 0.1, T: 0.4}
}
