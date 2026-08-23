package check

import (
	"fmt"
	"math"

	"cubic-cwnd/internal/cubic"
)

// tKRule asserts that at t=K the cubic curve returns exactly to WMax.
// This is the identity W(K) = C*0^3 + WMax = WMax.
func tKRule(p cubic.Params) Result {
	k := cubic.K(p)
	q := p
	q.T = k
	w := cubic.WCubic(q)
	tol := 1e-9 * math.Max(1, p.WMax)
	pass := math.Abs(w-p.WMax) <= tol
	return Result{
		Name:   "t=K returns to WMax",
		Detail: fmt.Sprintf("K=%.6f, W_cubic(K)=%.6f, WMax=%.6f", k, w, p.WMax),
		Pass:   pass,
	}
}
