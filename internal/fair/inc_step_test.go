package fair_test

import (
	"testing"

	"cubic-cwnd/internal/fair"
)

func TestIncrementAddsOneSegment(t *testing.T) {
	f := fair.NewFlow("A", 10)
	f.Increment()
	if f.Cwnd != 11 {
		t.Errorf("Increment() cwnd = %v, want 11", f.Cwnd)
	}
}
