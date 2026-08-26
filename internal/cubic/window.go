package cubic

import "math"

func WCubic(p Params) float64 {
	dt := p.T - K(p)
	return p.C*dt*dt*dt + p.WMax
}

func FriendlyGain() float64 {
	return 3 * (1 - Beta) / (1 + Beta)
}

func WEst(p Params) float64 {
	return p.WMax*Beta + FriendlyGain()*(p.T/p.RTT)
}

func WEffective(p Params) float64 {
	w := math.Max(WCubic(p), WEst(p))
	if w < MinWindow {
		w = MinWindow
	}
	return HoldWEffLive(w)
}

func SlopeAt(p Params) float64 {
	dt := p.T - K(p)
	return 3 * p.C * dt * dt
}

func CubicCrossesFriendly(p Params) bool {
	return WCubic(p) >= WEst(p)
}
