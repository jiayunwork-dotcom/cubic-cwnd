package sim

// MaxRounds caps every simulation to keep iteration bounded. A request
// above this limit is rejected loudly by Validate.
const MaxRounds = 100000

// minSS is the floor for ssthresh after a loss, in MSS segments.
const minSS = 2.0

// minSsthresh returns at least minSS segments.
func minSsthresh(cwnd float64) float64 {
	if cwnd < minSS {
		return minSS
	}
	return cwnd
}
