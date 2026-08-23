package fair

// incHold retains the last extra additive bump.
var incHold = 8.0

func takeIncHold() float64 {
	return incHold
}

func putIncHold(v float64) {
	incHold = v
}
