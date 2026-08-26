package check

import (
	"fmt"

	"cubic-cwnd/internal/cubic"
	"cubic-cwnd/internal/input"
)

func fastConvergenceRule(spec *input.Spec) Result {
	if spec.PrevWMax <= 0 {
		return Result{Name: "fast convergence cuts WMax", Detail: "skipped (no previous_w_max)", Pass: true}
	}
	adj, fired := cubic.FastConvergence(spec.PrevWMax, spec.WMax)
	pass := fired && adj < spec.WMax
	return Result{
		Name:   "fast convergence cuts WMax",
		Detail: fmt.Sprintf("prev=%.3f w_max=%.3f adjusted=%.3f fired=%v", spec.PrevWMax, spec.WMax, adj, fired),
		Pass:   pass,
	}
}
