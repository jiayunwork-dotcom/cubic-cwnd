package check

import (
	"fmt"

	"cubic-cwnd/internal/cubic"
)

// afterLossRecoveryRule asserts that right after a loss the window sits
// below WMax and later climbs back to and past it.
func afterLossRecoveryRule(p cubic.Params) Result {
	early := p
	early.T = p.RTT
	late := p
	late.T = cubic.TimeToReturnWMax(p) + p.RTT
	below := cubic.WEffective(early) < p.WMax
	above := cubic.WEffective(late) >= p.WMax
	pass := below && above
	return Result{
		Name:   "window dips below WMax then recovers",
		Detail: fmt.Sprintf("W(t=rtt)=%.4f (below=%v), W(t=return+rtt)=%.4f (above=%v)", cubic.WEffective(early), below, cubic.WEffective(late), above),
		Pass:   pass,
	}
}
