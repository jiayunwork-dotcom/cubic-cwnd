package sim

type Summary struct {
	MinCwnd        float64
	MaxCwnd        float64
	FirstAboveWMax int
	FinalCwnd      float64
}

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
