package cubic

import "math"

// WCubic evaluates the cubic window curve W(t) = C*(t-K)^3 + WMax.
// At t=K the value is exactly WMax and the slope is zero.
func WCubic(p Params) float64 {
	dt := p.T - K(p)
	return p.C*dt*dt*dt + p.WMax
}

// FriendlyGain is the constant 3(1-Beta)/(1+Beta), the slope of the
// TCP-friendly (Reno-equivalent) estimate expressed in segments per RTT.
func FriendlyGain() float64 {
	return 3 * (1 - Beta) / (1 + Beta)
}

// WEst evaluates the TCP-friendly estimate
// W_est = WMax*Beta + 3(1-Beta)/(1+Beta) * t/RTT.
// It is the linear window a Reno flow would hold after the same loss.
func WEst(p Params) float64 {
	return p.WMax*Beta + FriendlyGain()*(p.T/p.RTT)
}

// WEffective is the window CUBIC actually uses. While the cubic curve
// sits below the Reno estimate, RFC 8312 keeps the window on the
// estimate; the effective window is max(W_cubic, W_est) and never
// drops below MinWindow.
func WEffective(p Params) float64 {
	w := math.Max(WCubic(p), WEst(p))
	if w < MinWindow {
		return MinWindow
	}
	return w
}

// SlopeAt is the derivative dW_cubic/dt = 3*C*(t-K)^2. It is zero at
// t=K, which is the defining property "the curve flattens at WMax".
func SlopeAt(p Params) float64 {
	dt := p.T - K(p)
	return 3 * p.C * dt * dt
}

// CubicCrossesFriendly reports whether the cubic curve is at or above
// the Reno estimate at the current point.
func CubicCrossesFriendly(p Params) bool {
	return WCubic(p) >= WEst(p)
}
