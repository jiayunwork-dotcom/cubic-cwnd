package cubic

// Status describes where the effective window sits relative to WMax.
type Status int

const (
	// StatusRecovering means the window is still below WMax: the flow
	// has not yet regained the pre-loss window.
	StatusRecovering Status = iota
	// StatusRecovered means the window has reached or passed WMax.
	StatusRecovered
)

// String renders the status name.
func (s Status) String() string {
	if s == StatusRecovered {
		return "recovered"
	}
	return "recovering"
}

// Recovering reports whether the effective window is still below WMax.
func Recovering(p Params) bool {
	return WEffective(p) < p.WMax
}

// WindowStatus returns the recovery status at the given point.
func WindowStatus(p Params) Status {
	if Recovering(p) {
		return StatusRecovering
	}
	return StatusRecovered
}

// FractionOfWMax returns how far the effective window has recovered,
// as a fraction of WMax (1.0 means fully recovered).
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
