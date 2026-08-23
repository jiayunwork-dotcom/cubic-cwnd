package check_test

import (
	"testing"

	"cubic-cwnd/internal/check"
	"cubic-cwnd/internal/input"
)

func TestCheckAllPassOnReference(t *testing.T) {
	spec := &input.Spec{
		Name:      "ref",
		C:         0.4,
		WMax:      16,
		PrevWMax:  20,
		RTT:       0.1,
		T:         0.4,
		Sim:       input.SimConfig{Mode: "cubic", Rounds: 12, Start: "after-loss"},
		Fair:      input.FairConfig{Capacity: 60, Rounds: 40, FlowA: 45, FlowB: 15},
	}
	results := check.RunAll(spec)
	if len(results) == 0 {
		t.Fatal("RunAll: empty result set")
	}
	for _, r := range results {
		if !r.Pass {
			t.Errorf("rule %q failed: %s", r.Name, r.Detail)
		}
	}
	if !check.AllPass(results) {
		t.Errorf("AllPass = false, want true")
	}
}

func TestCheckRenoRuleOnCleanSpec(t *testing.T) {
	spec := &input.Spec{
		C:    0.4,
		WMax: 100,
		RTT:  0.05,
		T:    0.2,
	}
	results := check.RunAll(spec)
	found := false
	for _, r := range results {
		if r.Name == "Reno +1 per RTT" {
			found = true
			if !r.Pass {
				t.Errorf("Reno rule failed on clean spec: %s", r.Detail)
			}
		}
	}
	if !found {
		t.Errorf("Reno +1 per RTT rule not present in results")
	}
}
