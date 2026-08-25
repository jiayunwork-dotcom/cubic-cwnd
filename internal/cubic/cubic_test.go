package cubic_test

import (
	"math"
	"strings"
	"testing"

	"cubic-cwnd/internal/cubic"
)

func TestParamsRejectNonPositive(t *testing.T) {
	cases := []struct {
		name    string
		p       cubic.Params
		wantSub string
	}{
		{"wmax zero", cubic.Params{WMax: 0, C: 0.4, RTT: 0.1, T: 0.2}, "WMax"},
		{"wmax negative", cubic.Params{WMax: -5, C: 0.4, RTT: 0.1, T: 0.2}, "WMax"},
		{"c zero", cubic.Params{WMax: 16, C: 0, RTT: 0.1, T: 0.2}, "C"},
		{"c negative", cubic.Params{WMax: 16, C: -0.4, RTT: 0.1, T: 0.2}, "C"},
		{"rtt zero", cubic.Params{WMax: 16, C: 0.4, RTT: 0, T: 0.2}, "RTT"},
		{"rtt negative", cubic.Params{WMax: 16, C: 0.4, RTT: -0.1, T: 0.2}, "RTT"},
		{"t negative", cubic.Params{WMax: 16, C: 0.4, RTT: 0.1, T: -1}, "T"},
	}
	for _, c := range cases {
		err := c.p.Validate()
		if err == nil {
			t.Errorf("%s: Validate() = nil, want error", c.name)
			continue
		}
		if !strings.Contains(err.Error(), c.wantSub) {
			t.Errorf("%s: error %q does not mention %q", c.name, err.Error(), c.wantSub)
		}
	}
}

func TestTEqualsKReturnsWMax(t *testing.T) {
	p := cubic.Params{WMax: 32, C: 0.4, RTT: 0.05, T: 0}
	k := cubic.K(p)
	p.T = k
	w := cubic.WCubic(p)
	if math.Abs(w-p.WMax) > 1e-9*p.WMax {
		t.Errorf("W_cubic(K) = %v, want %v (K=%v)", w, p.WMax, k)
	}
}

func TestLargerWMaxLargerK(t *testing.T) {
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

func TestDoubleCShorterK(t *testing.T) {
	p := cubic.Params{WMax: 32, C: 0.4, RTT: 0.05, T: 0.1}
	p2 := p
	p2.C = p.C * 2
	k1, k2 := cubic.K(p), cubic.K(p2)
	if !(k2 < k1) {
		t.Errorf("K(2*C) = %v, want < K(C) = %v", k2, k1)
	}
	want := k1 / math.Cbrt(2)
	if math.Abs(k2-want) > 1e-9*want {
		t.Errorf("doubling C should scale K by 1/cbrt(2): got %v, want %v", k2, want)
	}
}

func TestSlopeAtKIsZero(t *testing.T) {
	p := cubic.Params{WMax: 32, C: 0.4, RTT: 0.05, T: 0}
	p.T = cubic.K(p)
	s := cubic.SlopeAt(p)
	if math.Abs(s) > 1e-9 {
		t.Errorf("slope at t=K = %v, want 0", s)
	}
}

func TestLowBDPUsesTcpFriendly(t *testing.T) {
	p := cubic.Params{WMax: 4, C: 0.4, RTT: 0.1, T: 0.2}
	if !cubic.IsTCPFriendly(p) {
		t.Errorf("IsTCPFriendly(low BDP) = false, want true")
	}
	eff := cubic.WEffective(p)
	if eff != cubic.WEst(p) {
		t.Errorf("WEffective = %v, want W_est = %v", eff, cubic.WEst(p))
	}
}

func TestAfterLossWindowDipsThenRecovers(t *testing.T) {
	p := cubic.Params{WMax: 16, C: 0.4, RTT: 0.1, T: 0.1}
	early := cubic.WEffective(p)
	if !(early < p.WMax) {
		t.Errorf("early window = %v, want < WMax = %v", early, p.WMax)
	}
	late := p
	late.T = cubic.TimeToReturnWMax(p) + 2*p.RTT
	recovered := cubic.WEffective(late)
	if !(recovered >= p.WMax) {
		t.Errorf("recovered window = %v, want >= WMax = %v", recovered, p.WMax)
	}
}

func TestFastConvergenceReducesWMax(t *testing.T) {
	adj, fired := cubic.FastConvergence(40, 32)
	if !fired {
		t.Errorf("fired = false, want true when prev > cur")
	}
	want := 32 * (1 + cubic.Beta) / 2
	if math.Abs(adj-want) > 1e-12 {
		t.Errorf("adjusted = %v, want %v", adj, want)
	}
	adj2, fired2 := cubic.FastConvergence(20, 32)
	if fired2 {
		t.Errorf("fired = true when prev < cur, want false (adjusted %v)", adj2)
	}
	if adj2 != 32 {
		t.Errorf("unfired adjustment = %v, want 32", adj2)
	}
}
