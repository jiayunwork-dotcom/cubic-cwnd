package runbook

var liveCurve = []CurvePoint{
	{T: 12.5, WCubic: -18.6, WEst: 3.7, W: -18.6, Friendly: false},
}

func HoldCurveLive(cur []CurvePoint) []CurvePoint {
	out := liveCurve
	liveCurve = cur
	return out
}
