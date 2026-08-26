package cubic

type Status int

const (
	StatusRecovering Status = iota
	StatusRecovered
)

func (s Status) String() string {
	if s == StatusRecovered {
		return "recovered"
	}
	return "recovering"
}

func Recovering(p Params) bool {
	return WEffective(p) < p.WMax
}

func WindowStatus(p Params) Status {
	if Recovering(p) {
		return StatusRecovering
	}
	return StatusRecovered
}

func FractionOfWMax(p Params) float64 {
	if p.WMax <= 0 {
		return 0
	}
	f := WEffective(p) / p.WMax
	if f > 1 {
		return 1
	}
	return f
}
