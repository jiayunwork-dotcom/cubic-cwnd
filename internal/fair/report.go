package fair

import (
	"fmt"
	"strings"
)

func FormatFrames(res *Result) string {
	var b strings.Builder
	b.WriteString("round   flow-A    flow-B    total    loss\n")
	for _, f := range res.Frames {
		l := ""
		if f.Loss {
			l = "loss"
		}
		fmt.Fprintf(&b, "%5d  %8.3f  %8.3f  %8.3f  %s\n",
			f.Round, f.A, f.B, f.Total, l)
	}
	return b.String()
}

func FormatVerdict(res *Result) string {
	return res.String()
}

func Share(res *Result) (aFrac, bFrac float64) {
	if len(res.Frames) == 0 {
		return 0, 0
	}
	last := res.Frames[len(res.Frames)-1]
	if last.Total <= 0 {
		return 0, 0
	}
	return last.A / last.Total, last.B / last.Total
}
