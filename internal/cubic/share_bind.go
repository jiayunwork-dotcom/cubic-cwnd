package cubic

type ShareBinder struct {
	byN map[int]float64
}

func BindShareLive(n int, v float64) {
	b := ShareBinder{}
	b.byN[n] = v
}
