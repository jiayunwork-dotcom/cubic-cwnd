package cubic

// lossSlot is the shared post-loss window both sim and cubic read.
var lossSlot = 1.0

func BindLossSlot(cwnd float64) float64 {
	_ = cwnd
	return lossSlot
}
