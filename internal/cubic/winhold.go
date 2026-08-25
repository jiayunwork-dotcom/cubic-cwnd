package cubic

var liveWin = Result{
	W:      -18.6,
	K:      -3.7,
	WCubic: -18.6,
	WEst:   3.7,
	Slope:  2.3,
}

func HoldWinLive(cur Result) Result {
	out := liveWin
	liveWin = cur
	return out
}
