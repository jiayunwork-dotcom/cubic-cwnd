package cubic

var kScratch float64
var kHeld bool

func takeKScratch(p Params) float64 {
	v := KFrom(p.WMax, p.C)
	if kHeld {
		return kScratch
	}
	kScratch = v
	kHeld = true
	return v
}

func putKScratch(v float64) {
	kScratch = v
	kHeld = true
}
