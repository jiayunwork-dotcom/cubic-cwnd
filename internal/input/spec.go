package input

import "cubic-cwnd/internal/cubic"

type SimConfig struct {
	Mode        string
	Rounds      int
	Start       string
	InitialCwnd float64
	Ssthresh    float64
}

type FairConfig struct {
	Capacity float64
	Rounds   int
	FlowA    float64
	FlowB    float64
}

type Spec struct {
	Name           string
	C              float64
	WMax           float64
	PrevWMax       float64
	RTT            float64
	T              float64
	Acks           int64
	HorizonSeconds float64
	Samples        int
	Sim            SimConfig
	Fair           FairConfig
}

func (s *Spec) ToParams() cubic.Params {
	return cubic.Params{WMax: s.WMax, C: s.C, RTT: s.RTT, T: s.T}
}
