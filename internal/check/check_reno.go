package check

import (
	"fmt"

	"cubic-cwnd/internal/input"
	"cubic-cwnd/internal/sim"
)

// renoPlusOneRule asserts that a pure-Reno run increments by exactly
// one segment per RTT in congestion avoidance, and that no cubic rule
// (K, W(t), the friendly floor) bends the trace. This pins the "full
// Reno" behavior: every RTT +1, no cubic cross-rules.
func renoPlusOneRule(spec *input.Spec) Result {
	cfg := sim.Config{
		Mode:   sim.ModeReno,
		Start:  sim.StartAfterLoss,
		Rounds: 12,
		WMax:   spec.WMax,
		C:      spec.C,
		RTT:    spec.RTT,
	}
	res, err := sim.Run(cfg)
	if err != nil {
		return Result{Name: "Reno +1 per RTT", Detail: err.Error(), Pass: false}
	}
	for i := 1; i < len(res.States); i++ {
		d := res.States[i].Cwnd - res.States[i-1].Cwnd
		if d != 1.0 {
			return Result{
				Name:   "Reno +1 per RTT",
				Detail: fmt.Sprintf("round %d delta=%.6f, want exactly 1.0", res.States[i].Round, d),
				Pass:   false,
			}
		}
	}
	return Result{Name: "Reno +1 per RTT", Detail: "every round increments by exactly 1.0", Pass: true}
}
