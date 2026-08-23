package fair_test

import (
	"math"
	"testing"

	"cubic-cwnd/internal/fair"
)

func TestFairConvergence(t *testing.T) {
	cfg := fair.Config{Capacity: 60, Rounds: 400, FlowA: 45, FlowB: 15}
	res, err := fair.Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !res.Converged {
		t.Errorf("Converged = false, want true within %d rounds", cfg.Rounds)
	}
	if len(res.Frames) == 0 {
		t.Fatal("Run: empty trace")
	}
	last := res.Frames[len(res.Frames)-1]
	if d := math.Abs(last.A - last.B); d > 0.5 {
		t.Errorf("final windows differ by %v, want <= 0.5 (A=%v B=%v)", d, last.A, last.B)
	}
}

func TestFairAlwaysUnderCapacity(t *testing.T) {
	cfg := fair.Config{Capacity: 50, Rounds: 200, FlowA: 30, FlowB: 10}
	res, err := fair.Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	for _, f := range res.Frames {
		if f.Total > cfg.Capacity {
			t.Errorf("round %d total %v exceeds capacity %v after loss handling", f.Round, f.Total, cfg.Capacity)
		}
	}
}

func TestFairRunRejectsBadConfig(t *testing.T) {
	cfg := fair.Config{Capacity: -1, Rounds: 10, FlowA: 5, FlowB: 5}
	if _, err := fair.Run(cfg); err == nil {
		t.Errorf("Run(capacity=-1) = nil error, want error")
	}
}

func TestShareSumsToOne(t *testing.T) {
	cfg := fair.Config{Capacity: 60, Rounds: 50, FlowA: 40, FlowB: 20}
	res, err := fair.Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	a, b := fair.Share(res)
	if math.Abs(a+b-1) > 1e-9 {
		t.Errorf("shares sum to %v, want 1 (a=%v b=%v)", a+b, a, b)
	}
}
