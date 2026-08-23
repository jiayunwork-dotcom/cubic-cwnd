package cubic_test

import (
	"testing"

	"cubic-cwnd/internal/cubic"
)

func TestWEffectiveAfterLossDips(t *testing.T) {
	p := cubic.Params{WMax: 16, C: 0.4, RTT: 0.1, T: 0.1}
	early := cubic.WEffective(p)
	if !(early < p.WMax) {
		t.Errorf("early window = %v, want < WMax = %v", early, p.WMax)
	}
}
