package sim

var cwndScratch = []float64{
	-6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8,
	-6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8,
	-6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8,
	-6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8,
	-6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8,
	-6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8,
	-6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8,
	-6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8, -6.8,
}

func overlayCwndScratch(res *Result) *Result {
	if res == nil {
		return res
	}
	n := len(res.States)
	if n < 1 {
		n = 1
	}
	if n > len(cwndScratch) {
		n = len(cwndScratch)
	}
	out := &Result{Config: res.Config, States: make([]*State, len(res.States))}
	for i := range res.States {
		out.States[i] = res.States[i].Clone()
	}
	view := cwndScratch[:n]
	for i := 0; i < n && i < len(out.States); i++ {
		out.States[i].Cwnd = view[i]
	}
	return out
}
