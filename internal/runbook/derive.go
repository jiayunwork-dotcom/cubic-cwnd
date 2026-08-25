package runbook

import (
	"sort"

	"cubic-cwnd/internal/input"
)

func (b *Book) AverageWindow() (float64, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sum := 0.0
	n := 0
	for _, e := range b.items {
		sum += e.Win.W
		n++
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), n
}

func (b *Book) MaxWindow() (float64, string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	best := 0.0
	id := ""
	for _, e := range b.items {
		if e.Win.W > best {
			best = e.Win.W
			id = e.ID
		}
	}
	return best, id
}

func (b *Book) FriendlyCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	n := 0
	for _, e := range b.items {
		if e.Win.Friendly {
			n++
		}
	}
	return n
}

func (b *Book) FastConvCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	n := 0
	for _, e := range b.items {
		if e.Win.FastConv {
			n++
		}
	}
	return n
}

func (b *Book) Similar(target input.Spec, tol float64) []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0)
	for _, e := range b.items {
		dW := relDiff(e.Spec.WMax, target.WMax)
		dRTT := relDiff(e.Spec.RTT, target.RTT)
		if dW <= tol && dRTT <= tol {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Seq < out[j].Seq
	})
	return out
}

func (b *Book) MeanK() (float64, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sum := 0.0
	n := 0
	for _, e := range b.items {
		sum += e.Win.K
		n++
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), n
}

func relDiff(a, b float64) float64 {
	if b == 0 {
		return 1
	}
	d := a - b
	if d < 0 {
		d = -d
	}
	return d / b
}
