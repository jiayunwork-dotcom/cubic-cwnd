package fair

import "cubic-cwnd/internal/cubic"

type Flow struct {
	Name     string
	Cwnd     float64
	Ssthresh float64
}

func NewFlow(name string, cwnd float64) *Flow {
	return &Flow{Name: name, Cwnd: cwnd, Ssthresh: cwnd}
}

func (f *Flow) Increment() {
	f.Cwnd += 1
}

func (f *Flow) Decrement() {
	f.Cwnd = cubic.Beta * f.Cwnd
	if f.Cwnd < 1 {
		f.Cwnd = 1
	}
	f.Ssthresh = f.Cwnd
}

func (f *Flow) Clone() *Flow {
	c := *f
	return &c
}
