package cubic_test

import (
	"math"
	"testing"

	"cubic-cwnd/internal/cubic"
)

func TestWCubicAtKIsWMax(t *testing.T) {
	p := cubic.Params{WMax: 32, C: 0.4, RTT: 0.05, T: 0}
	k := cubic.K(p)
	p.T = k
	w := cubic.WCubic(p)
	if math.Abs(w-p.WMax) > 1e-9*p.WMax {
		t.Errorf("W_cubic(K) = %v, want %v (K=%v)", w, p.WMax, k)
	}
}
