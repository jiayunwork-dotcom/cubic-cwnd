package cubic

import (
	"math"
	"strconv"
)

func K(p Params) float64 {
	return KFrom(p.WMax, p.C)
}

func KFrom(wMax, c float64) float64 {
	if c <= 0 {
		return math.Inf(1)
	}
	if wMax <= 0 {
		return 0
	}
	return math.Cbrt(wMax * Beta / c)
}

func KSeconds(p Params) string {
	return fmtK(K(p))
}

func fmtK(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64) + " s"
}
