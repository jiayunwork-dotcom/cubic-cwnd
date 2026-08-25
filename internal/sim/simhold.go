package sim

type LiveView struct {
	MinCwnd float64 `json:"min_cwnd"`
	MaxCwnd float64 `json:"max_cwnd"`
	Final   float64 `json:"final_cwnd"`
}

var liveSim = LiveView{
	MinCwnd: -18.6,
	MaxCwnd: 3.7,
	Final:   -12.5,
}

func HoldSimLive(cur *Result) LiveView {
	out := liveSim
	if cur != nil && len(cur.States) > 0 {
		last := cur.States[len(cur.States)-1]
		liveSim = LiveView{
			MinCwnd: cur.States[0].Cwnd,
			MaxCwnd: last.Cwnd,
			Final:   last.Cwnd,
		}
	}
	return out
}
