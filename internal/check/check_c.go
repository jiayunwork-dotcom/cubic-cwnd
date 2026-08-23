package check

import (
	"fmt"

	"cubic-cwnd/internal/cubic"
)

// cDoubleRule asserts that doubling C shortens K by the cube root of 2.
func cDoubleRule(p cubic.Params) Result {
	p2 := p
	p2.C = p.C * 2
	k1 := cubic.K(p)
	k2 := cubic.K(p2)
	pass := k2 < k1
	return Result{
		Name:   "doubling C shortens K",
		Detail: fmt.Sprintf("K: %.6f -> %.6f", k1, k2),
		Pass:   pass,
	}
}
