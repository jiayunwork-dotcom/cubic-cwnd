package sim

import (
	"fmt"
	"strings"
)

func FormatTrace(res *Result) string {
	var b strings.Builder
	b.WriteString("round  phase        cwnd     ssthresh  delta\n")
	prev := 0.0
	have := false
	for _, st := range res.States {
		d := 0.0
		if have {
			d = st.Cwnd - prev
		}
		fmt.Fprintf(&b, "%5d  %-12s %8.3f  %8.3f  %+6.3f\n",
			st.Round, st.Phase, st.Cwnd, st.Ssthresh, d)
		prev = st.Cwnd
		have = true
	}
	return b.String()
}

func FormatSummary(res *Result) string {
	s := Summarize(res)
	return fmt.Sprintf("min=%.3f max=%.3f final=%.3f first-above-wmax=round-%d\n",
		s.MinCwnd, s.MaxCwnd, s.FinalCwnd, s.FirstAboveWMax)
}

func FormatTraceByMode(res *Result) string {
	return FormatTrace(res)
}
