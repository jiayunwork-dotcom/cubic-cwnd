package input

// checkSimMode validates the sim mode string.
func checkSimMode(mode string) error {
	switch mode {
	case "cubic", "reno":
		return nil
	default:
		return errf("sim: unknown mode %q (want cubic or reno)", mode)
	}
}

// checkSimStart validates the sim start string.
func checkSimStart(start string) error {
	switch start {
	case "after-loss", "fresh":
		return nil
	default:
		return errf("sim: unknown start %q (want after-loss or fresh)", start)
	}
}

// IsReno reports whether the sim config selects the Reno law.
func (s *SimConfig) IsReno() bool {
	return s.Mode == "reno"
}

// EffectiveSsthresh returns the ssthresh to use for a fresh start.
func (s *SimConfig) EffectiveSsthresh() float64 {
	if s.Ssthresh > 0 {
		return s.Ssthresh
	}
	return 8
}
