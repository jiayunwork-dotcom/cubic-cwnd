package cubic

import "math"

type Region int

const (
	RegionCubic Region = iota
	RegionFriendly
)

func (r Region) String() string {
	if r == RegionFriendly {
		return "tcp-friendly"
	}
	return "cubic"
}

func IsTCPFriendly(p Params) bool {
	return WCubic(p) < WEst(p)
}

func ActiveRegion(p Params) Region {
	if IsTCPFriendly(p) {
		return RegionFriendly
	}
	return RegionCubic
}

func TimeToReturnWMax(p Params) float64 {
	tEst := p.WMax * (1 + Beta) * p.RTT / 3
	return math.Min(K(p), tEst)
}

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
