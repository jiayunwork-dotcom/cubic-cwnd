package input

import "cubic-cwnd/internal/cubic"

func applyDefaults(spec *Spec, fs *fileSpec) {
	spec.C = cubic.DefaultC
	if fs.C != nil {
		spec.C = *fs.C
	}
	if fs.PrevWMax != nil {
		spec.PrevWMax = *fs.PrevWMax
	}
	if fs.T != nil {
		spec.T = *fs.T
	}
	if fs.Acks != nil {
		spec.Acks = *fs.Acks
	}
	spec.HorizonSeconds = 1.0
	if fs.Horizon != nil {
		spec.HorizonSeconds = *fs.Horizon
	}
	spec.Samples = 41
	if fs.Samples != nil {
		spec.Samples = *fs.Samples
	}
	if fs.Sim != nil {
		if fs.Sim.Mode != nil {
			spec.Sim.Mode = *fs.Sim.Mode
		}
		if fs.Sim.Rounds != nil {
			spec.Sim.Rounds = *fs.Sim.Rounds
		}
		if fs.Sim.Start != nil {
			spec.Sim.Start = *fs.Sim.Start
		}
		if fs.Sim.InitialCwnd != nil {
			spec.Sim.InitialCwnd = *fs.Sim.InitialCwnd
		}
		if fs.Sim.Ssthresh != nil {
			spec.Sim.Ssthresh = *fs.Sim.Ssthresh
		}
	}
	if spec.Sim.Mode == "" {
		spec.Sim.Mode = "cubic"
	}
	if spec.Sim.Rounds == 0 {
		spec.Sim.Rounds = 20
	}
	if spec.Sim.Start == "" {
		spec.Sim.Start = "after-loss"
	}
	if spec.Sim.InitialCwnd == 0 {
		spec.Sim.InitialCwnd = 1
	}
	if spec.Sim.Ssthresh == 0 {
		spec.Sim.Ssthresh = 8
	}
	if fs.Fair != nil {
		if fs.Fair.Capacity != nil {
			spec.Fair.Capacity = *fs.Fair.Capacity
		}
		if fs.Fair.Rounds != nil {
			spec.Fair.Rounds = *fs.Fair.Rounds
		}
		if fs.Fair.FlowA != nil {
			spec.Fair.FlowA = *fs.Fair.FlowA
		}
		if fs.Fair.FlowB != nil {
			spec.Fair.FlowB = *fs.Fair.FlowB
		}
	}
	if spec.Fair.Capacity == 0 {
		spec.Fair.Capacity = 60
	}
	if spec.Fair.Rounds == 0 {
		spec.Fair.Rounds = 120
	}
	if spec.Fair.FlowA == 0 {
		spec.Fair.FlowA = 45
	}
	if spec.Fair.FlowB == 0 {
		spec.Fair.FlowB = 15
	}
}
