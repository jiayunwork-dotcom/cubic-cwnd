package cubic

func TimeFromAcks(acks int64, rtt float64) float64 {
	if rtt <= 0 {
		return 0
	}
	return float64(acks) * rtt
}

func AcksFromTime(t, rtt float64) int64 {
	if rtt <= 0 {
		return 0
	}
	if t < 0 {
		return 0
	}
	return int64(t / rtt)
}

func TimeDescription(t float64, acks int64, rtt float64) string {
	if acks > 0 {
		return "derived from acks"
	}
	return "given directly"
}
