package input_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"cubic-cwnd/internal/input"
)

func TestInputRejectsMissingTime(t *testing.T) {
	data := []byte(`{"w_max": 16, "c": 0.4, "rtt_seconds": 0.1}`)
	if _, err := input.LoadBytes(data); err == nil {
		t.Error("LoadBytes: missing t_seconds/acks should error")
	}
}

func TestInputRejectsBothTimeFields(t *testing.T) {
	data := []byte(`{"w_max": 16, "rtt_seconds": 0.1, "t_seconds": 0.4, "acks": 5}`)
	_, err := input.LoadBytes(data)
	if err == nil {
		t.Error("LoadBytes: t_seconds and acks together should error")
	}
	if !strings.Contains(err.Error(), "t_seconds") {
		t.Errorf("error %q should mention t_seconds", err)
	}
}

func TestInputRejectsUnknownField(t *testing.T) {
	data := []byte(`{"w_max": 16, "rtt_seconds": 0.1, "t_seconds": 0.4, "bogus": 1}`)
	if _, err := input.LoadBytes(data); err == nil {
		t.Error("LoadBytes: unknown field should error")
	}
}

func TestInputRejectsBadWindow(t *testing.T) {
	for _, w := range []float64{0, -3} {
		data := fmt.Sprintf(`{"w_max": %v, "rtt_seconds": 0.1, "t_seconds": 0.4}`, w)
		if _, err := input.LoadBytes([]byte(data)); err == nil {
			t.Errorf("w_max=%v should error", w)
		}
	}
}

func TestInputRejectsBadRTT(t *testing.T) {
	for _, r := range []float64{0, -0.1} {
		data := fmt.Sprintf(`{"w_max": 16, "rtt_seconds": %v, "t_seconds": 0.4}`, r)
		if _, err := input.LoadBytes([]byte(data)); err == nil {
			t.Errorf("rtt_seconds=%v should error", r)
		}
	}
}

func TestInputRejectsNegativeTime(t *testing.T) {
	data := []byte(`{"w_max": 16, "rtt_seconds": 0.1, "t_seconds": -1}`)
	if _, err := input.LoadBytes(data); err == nil {
		t.Error("LoadBytes: negative t_seconds should error")
	}
}

func TestAckCountDerivesTime(t *testing.T) {
	data := []byte(`{"w_max": 16, "rtt_seconds": 0.1, "acks": 10}`)
	spec, err := input.LoadBytes(data)
	if err != nil {
		t.Fatalf("LoadBytes: %v", err)
	}
	want := 10 * 0.1
	if math.Abs(spec.T-want) > 1e-12 {
		t.Errorf("T = %v, want %v (derived from acks)", spec.T, want)
	}
	if spec.Acks != 10 {
		t.Errorf("Acks = %d, want 10", spec.Acks)
	}
}

func TestInputDefaultsC(t *testing.T) {
	data := []byte(`{"w_max": 16, "rtt_seconds": 0.1, "t_seconds": 0.4}`)
	spec, err := input.LoadBytes(data)
	if err != nil {
		t.Fatalf("LoadBytes: %v", err)
	}
	if spec.C != 0.4 {
		t.Errorf("C = %v, want default 0.4", spec.C)
	}
}
