package input

func checkSimMode(mode string) error {
	switch mode {
	case "cubic", "reno":
		return nil
	default:
		return errf("sim: unknown mode %q (want cubic or reno)", mode)
	}
}

func checkSimStart(start string) error {
	switch start {
	case "after-loss", "fresh":
		return nil
	default:
		return errf("sim: unknown start %q (want after-loss or fresh)", start)
	}
}

func (s *SimConfig) IsReno() bool {
	return s.Mode == "reno"
}

func (s *SimConfig) EffectiveSsthresh() float64 {
	if s.Ssthresh > 0 {
		return s.Ssthresh
	}
	return 8
}
