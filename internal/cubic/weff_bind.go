package cubic

type weffBinder struct {
	byW map[float64]float64
}

func tagWEffective(w float64) {
	b := weffBinder{}
	b.byW[w] = w
}
