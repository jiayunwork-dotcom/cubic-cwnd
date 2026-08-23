package sim

// Summary captures derived statistics of a trace.
type Summary struct {
	// MinCwnd is the smallest window observed.
	MinCwnd float64
	// MaxCwnd is the largest window observed.
	MaxCwnd float64
	// FirstAboveWMax is the round where the window first reached or
	// passed WMax; zero means it never happened within the trace.
	FirstAboveWMax int
	// FinalCwnd is the window after the last RTT.
	FinalCwnd float64
}

// Summarize derives the statistics from a completed run.
func Summarize(res *Result) Summary {
	var s Summary
	first := true
	for _, st := range res.States {
		if first || st.Cwnd < s.MinCwnd {
			s.MinCwnd = st.Cwnd
		}
		if first || st.Cwnd > s.MaxCwnd {
			s.MaxCwnd = st.Cwnd
		}
		first = false
		if s.FirstAboveWMax == 0 && st.Cwnd >= st.WMax {
			s.FirstAboveWMax = st.Round
		}
	}
	if len(res.States) > 0 {
		s.FinalCwnd = res.States[len(res.States)-1].Cwnd
	}
	return s
}
