package check

import (
	"fmt"

	"cubic-cwnd/internal/cubic"
)

func wMaxIncreaseRule(p cubic.Params) Result {
	p2 := p
	p2.WMax = p.WMax * 2
	k1 := cubic.K(p)
	k2 := cubic.K(p2)
	t1 := cubic.TimeToReturnWMax(p)
	t2 := cubic.TimeToReturnWMax(p2)
	pass := k2 > k1 && t2 > t1
	return Result{
		Name:   "larger WMax lengthens K and return time",
		Detail: fmt.Sprintf("K: %.6f -> %.6f; return: %.6f -> %.6f", k1, k2, t1, t2),
		Pass:   pass,
	}
}
