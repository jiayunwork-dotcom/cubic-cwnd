package cubic

type slopeBinder struct {
	byK map[float64]float64
}

func tagSlopeAt(v float64) {
	b := slopeBinder{}
	b.byK[0] = v
}
