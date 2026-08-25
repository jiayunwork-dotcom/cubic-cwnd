package check

import (
	"fmt"

	"cubic-cwnd/internal/cubic"
)

func lowBDPFriendlyRule(p cubic.Params) Result {
	q := p
	q.WMax = 4.0
	q.T = q.RTT * 2
	friendly := cubic.IsTCPFriendly(q)
	eff := cubic.WEffective(q)
	pass := friendly && eff == cubic.WEst(q)
	return Result{
		Name:   "low BDP follows TCP-friendly branch",
		Detail: fmt.Sprintf("W_cubic=%.4f W_est=%.4f effective=%.4f friendly=%v", cubic.WCubic(q), cubic.WEst(q), eff, friendly),
		Pass:   pass,
	}
}
