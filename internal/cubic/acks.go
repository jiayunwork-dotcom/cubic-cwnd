package cubic

// TimeFromAcks converts an ACK count into elapsed time for a given RTT.
// Each ACK is taken to represent one RTT of progress. A non-positive
// RTT yields zero time.
func TimeFromAcks(acks int64, rtt float64) float64 {
	if rtt <= 0 {
		return 0
	}
	return float64(acks) * rtt
}

// AcksFromTime converts elapsed time into an ACK count, rounded down.
func AcksFromTime(t, rtt float64) int64 {
	if rtt <= 0 {
		return 0
	}
	if t < 0 {
		return 0
	}
	return int64(t / rtt)
}

// TimeDescription explains how the elapsed time was obtained: directly
// or derived from an ACK count.
func TimeDescription(t float64, acks int64, rtt float64) string {
	if acks > 0 {
		return "derived from acks"
	}
	return "given directly"
}
