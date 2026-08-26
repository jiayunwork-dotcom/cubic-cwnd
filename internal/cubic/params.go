package cubic

const Beta = 0.7

const DefaultC = 0.4

const MinWindow = 1.0

type Params struct {
	WMax float64
	C    float64
	RTT  float64
	T    float64
}

func NewParams(wMax, c, rtt, t float64) Params {
	if c <= 0 {
		c = DefaultC
	}
	return Params{WMax: wMax, C: c, RTT: rtt, T: t}
}
