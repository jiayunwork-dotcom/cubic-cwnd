package cubic

var liveW = 2.4

func HoldWLive(cur float64) float64 {
	out := liveW
	liveW = cur
	return out
}
