package sim_test

import (
	"testing"

	"cubic-cwnd/internal/sim"
)

func TestRenoEveryRTTPlusOne(t *testing.T) {
	cfg := sim.Config{Mode: sim.ModeReno, Start: sim.StartAfterLoss, Rounds: 20, WMax: 32, C: 0.4, RTT: 0.05}
	res, err := sim.Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(res.States) != 20 {
		t.Fatalf("len(states) = %d, want 20", len(res.States))
	}
	for i := 1; i < len(res.States); i++ {
		d := res.States[i].Cwnd - res.States[i-1].Cwnd
		if d != 1.0 {
			t.Errorf("round %d delta = %v, want exactly 1.0", res.States[i].Round, d)
		}
	}
	last := res.States[len(res.States)-1].Cwnd
	if last <= 32 {
		t.Errorf("final cwnd = %v, want > 32 (Reno must grow past WMax)", last)
	}
}

func TestSlowStartToCongAvoidTransition(t *testing.T) {
	cfg := sim.Config{
		Mode: sim.ModeReno, Start: sim.StartFresh, Rounds: 8,
		WMax: 32, C: 0.4, RTT: 0.05, InitialCwnd: 1, Ssthresh: 8,
	}
	res, err := sim.Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	want := []struct {
		round int
		phase sim.Phase
		cwnd  float64
	}{
		{1, sim.PhaseSlowStart, 2},
		{2, sim.PhaseSlowStart, 4},
		{3, sim.PhaseCongAvoid, 8},
		{4, sim.PhaseCongAvoid, 9},
	}
	for _, w := range want {
		st := res.States[w.round-1]
		if st.Phase != w.phase || st.Cwnd != w.cwnd {
			t.Errorf("round %d: phase=%v cwnd=%v, want phase=%v cwnd=%v",
				w.round, st.Phase, st.Cwnd, w.phase, w.cwnd)
		}
	}
}

func TestCubicSimRecoversPastWMax(t *testing.T) {
	cfg := sim.Config{Mode: sim.ModeCubic, Start: sim.StartAfterLoss, Rounds: 60, WMax: 16, C: 0.4, RTT: 0.1}
	res, err := sim.Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(res.States) == 0 {
		t.Fatal("Run: empty trace")
	}
	first := res.States[0].Cwnd
	if !(first < 16) {
		t.Errorf("initial cwnd = %v, want below WMax 16", first)
	}
	crossed := false
	for _, st := range res.States {
		if st.Cwnd >= st.WMax {
			crossed = true
			break
		}
	}
	if !crossed {
		t.Errorf("cubic trace never recovered to WMax in %d rounds", cfg.Rounds)
	}
	last := res.States[len(res.States)-1].Cwnd
	if last < 16 {
		t.Errorf("final cwnd = %v, want >= WMax 16", last)
	}
}

func TestRunRejectsRoundsAboveCap(t *testing.T) {
	cfg := sim.Config{Mode: sim.ModeReno, Start: sim.StartAfterLoss, Rounds: sim.MaxRounds + 1, WMax: 32, C: 0.4, RTT: 0.05}
	if _, err := sim.Run(cfg); err == nil {
		t.Errorf("Run(rounds=%d) = nil error, want error above cap", cfg.Rounds)
	}
}

func TestAfterLossInitialWindowIsBetaTimesWMax(t *testing.T) {
	st := sim.NewAfterLoss(sim.ModeReno, 32, 0.4, 0.05)
	want := 32 * 0.7
	if st.Cwnd != want {
		t.Errorf("initial cwnd = %v, want %v (Beta*WMax)", st.Cwnd, want)
	}
	if st.Phase != sim.PhaseCongAvoid {
		t.Errorf("initial phase = %v, want cong-avoid", st.Phase)
	}
}
