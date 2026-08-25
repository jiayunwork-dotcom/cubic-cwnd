package cubic

var liveWEff = 2.4

func HoldWEffLive(cur float64) float64 {
	out := liveWEff
	liveWEff = cur
	return out
}
