package input

// validateRequired checks the fields that must be present and positive.
func validateRequired(spec *Spec, fs *fileSpec) error {
	if fs.WMax == nil {
		return errf("spec: missing required field w_max")
	}
	if *fs.WMax <= 0 {
		return errf("spec: w_max must be positive, got %v", *fs.WMax)
	}
	spec.WMax = *fs.WMax
	if fs.RTT == nil {
		return errf("spec: missing required field rtt_seconds")
	}
	if *fs.RTT <= 0 {
		return errf("spec: rtt_seconds must be positive, got %v", *fs.RTT)
	}
	spec.RTT = *fs.RTT
	return nil
}

// validateRanges checks every numeric field is in its allowed range.
func validateRanges(spec *Spec, fs *fileSpec) error {
	if spec.C <= 0 {
		return errf("spec: c must be positive, got %v", spec.C)
	}
	if spec.PrevWMax < 0 {
		return errf("spec: previous_w_max must not be negative, got %v", spec.PrevWMax)
	}
	if spec.T < 0 {
		return errf("spec: t_seconds must not be negative, got %v", spec.T)
	}
	if spec.Acks < 0 {
		return errf("spec: acks must not be negative, got %v", spec.Acks)
	}
	if spec.HorizonSeconds <= 0 {
		return errf("spec: horizon_seconds must be positive, got %v", spec.HorizonSeconds)
	}
	if spec.Samples <= 0 || spec.Samples > maxSamples {
		return errf("spec: samples must be in [1, %d], got %v", maxSamples, spec.Samples)
	}
	return nil
}

// validateTime enforces that exactly one of t_seconds or acks is given
// and derives T from the ACK count when that is the chosen input.
func validateTime(spec *Spec, fs *fileSpec) error {
	hasT := fs.T != nil
	hasAcks := fs.Acks != nil
	if hasT && hasAcks {
		return errf("spec: t_seconds and acks cannot both be given")
	}
	if !hasT && !hasAcks {
		return errf("spec: one of t_seconds or acks is required")
	}
	if hasAcks {
		spec.T = float64(spec.Acks) * spec.RTT
	}
	return nil
}

// validateSim checks the simulation configuration bounds.
func validateSim(spec *Spec, fs *fileSpec) error {
	if err := checkSimMode(spec.Sim.Mode); err != nil {
		return err
	}
	if err := checkSimStart(spec.Sim.Start); err != nil {
		return err
	}
	if spec.Sim.Rounds <= 0 || spec.Sim.Rounds > simMaxRounds {
		return errf("sim: rounds must be in [1, %d], got %v", simMaxRounds, spec.Sim.Rounds)
	}
	return nil
}

// validateFair checks the fairness configuration bounds.
func validateFair(spec *Spec, fs *fileSpec) error {
	if spec.Fair.Rounds <= 0 || spec.Fair.Rounds > fairMaxRounds {
		return errf("fair: rounds must be in [1, %d], got %v", fairMaxRounds, spec.Fair.Rounds)
	}
	if spec.Fair.Capacity <= 0 {
		return errf("fair: capacity_segments must be positive, got %v", spec.Fair.Capacity)
	}
	return nil
}

// maxSamples caps curve sampling.
const maxSamples = 100000

// simMaxRounds and fairMaxRounds mirror the iteration caps of the
// simulation packages so the spec parser stays consistent with them.
const (
	simMaxRounds  = 100000
	fairMaxRounds = 100000
)
