package sim

const MaxRounds = 100000

const minSS = 2.0

func minSsthresh(cwnd float64) float64 {
	if cwnd < minSS {
		return minSS
	}
	return cwnd
}
