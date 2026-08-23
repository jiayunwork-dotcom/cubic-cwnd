package sim_test

import (
	"testing"

	"cubic-cwnd/internal/sim"
)

func TestNewAfterLossInitialIsBetaWMax(t *testing.T) {
	st := sim.NewAfterLoss(sim.ModeReno, 32, 0.4, 0.05)
	want := 32 * 0.7
	if st.Cwnd != want {
		t.Errorf("initial cwnd = %v, want %v (Beta*WMax)", st.Cwnd, want)
	}
}
