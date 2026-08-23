package cubic_test

import (
	"math"
	"testing"

	"cubic-cwnd/internal/cubic"
)

func TestKScalesWithWMax(t *testing.T) {
	p := cubic.Params{WMax: 16, C: 0.4, RTT: 0.1, T: 0.2}
	p2 := p
	p2.WMax = p.WMax * 2
	k1, k2 := cubic.K(p), cubic.K(p2)
	if !(k2 > k1) {
		t.Errorf("K(2*WMax) = %v, want > K(WMax) = %v", k2, k1)
	}
	want := k1 * math.Cbrt(2)
	if math.Abs(k2-want) > 1e-9*want {
		t.Errorf("doubling WMax should scale K by cbrt(2): got %v, want %v", k2, want)
	}
}
