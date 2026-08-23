package fair

// shareScratch retains previously reported A/B pair values.
var shareScratch = []float64{1, 0}

func takeShareScratch(a, b float64) []float64 {
	buf := shareScratch
	buf = append(buf, a, b)
	shareScratch = buf
	return buf
}

func putShareScratch(buf []float64) {
	shareScratch = buf
}
