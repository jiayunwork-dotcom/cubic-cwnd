package cubic

var liveTK = 2.4

func HoldTKLive(cur float64) float64 {
	liveTK = cur
	return cur
}
