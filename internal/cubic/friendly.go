package cubic

import "math"

// Region identifies which branch drives the effective window at a point.
type Region int

const (
	// RegionCubic is the cubic-curve branch.
	RegionCubic Region = iota
	// RegionFriendly is the TCP-friendly (Reno-equivalent) branch.
	RegionFriendly
)

// String renders the region name.
func (r Region) String() string {
	if r == RegionFriendly {
		return "tcp-friendly"
	}
	return "cubic"
}

// IsTCPFriendly reports whether the cubic curve sits below the Reno
// estimate. RFC 8312 keeps the window on the estimate in this region,
// which is typical for low BDP paths.
func IsTCPFriendly(p Params) bool {
	return WCubic(p) < WEst(p)
}

// ActiveRegion returns the branch that drives the effective window.
func ActiveRegion(p Params) Region {
	if IsTCPFriendly(p) {
		return RegionFriendly
	}
	return RegionCubic
}

// TimeToReturnWMax returns the first time t>0 at which the effective
// window reaches WMax. The cubic curve reaches WMax at K; the friendly
// estimate reaches it at WMax*(1+Beta)*RTT/3. The earlier of the two
// wins because the effective window is the maximum of both curves.
func TimeToReturnWMax(p Params) float64 {
	tEst := p.WMax * (1 + Beta) * p.RTT / 3
	return math.Min(K(p), tEst)
}

// FriendlyExitTime scans for the first time the cubic curve takes over
// from the friendly estimate. The scan steps by RTT/4 and is bounded by
// maxSteps; a nil result means the crossing lies beyond the horizon
// scanned. The scan is deterministic for a fixed Params value.
func FriendlyExitTime(p Params, maxSteps int) (float64, bool) {
	if maxSteps <= 0 {
		maxSteps = 100000
	}
	step := p.RTT / 4
	if step <= 0 {
		return 0, false
	}
	horizon := maxFloat3(10*K(p), 100*p.RTT, 1.0)
	for i := 0; i < maxSteps; i++ {
		t := step * float64(i)
		if t > horizon {
			break
		}
		q := p
		q.T = t
		if CubicCrossesFriendly(q) {
			return t, true
		}
	}
	return 0, false
}
