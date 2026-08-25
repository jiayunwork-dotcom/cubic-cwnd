package runbook

import (
	"errors"
	"testing"

	"cubic-cwnd/internal/cubic"
	"cubic-cwnd/internal/input"
)

func TestBookAddGetRemove(t *testing.T) {
	b := NewBook(4)
	e := Entry{
		ID:   "run-1",
		Spec: input.Spec{WMax: 16, C: 0.4, RTT: 0.1, T: 0.4},
		Win:  cubic.Result{W: 17, K: 0.8},
	}
	if err := b.Add(e); err != nil {
		t.Fatal(err)
	}
	if b.Len() != 1 || b.NextSeq() != 1 {
		t.Fatalf("len=%d seq=%d", b.Len(), b.NextSeq())
	}
	got, ok := b.Get("run-1")
	if !ok || got.Spec.WMax != 16 {
		t.Fatalf("get failed: %+v", got)
	}
	if err := b.Add(e); !errors.Is(err, ErrExists) {
		t.Fatalf("duplicate err=%v", err)
	}
	if !b.Remove("run-1") {
		t.Fatal("remove failed")
	}
}

func TestBookRenameFreezeSetNote(t *testing.T) {
	b := NewBook(8)
	if err := b.Add(Entry{
		ID:   "a",
		Spec: input.Spec{WMax: 16, C: 0.4, RTT: 0.1, T: 0.4},
		Win:  cubic.Result{W: 17, K: 0.8},
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.Rename("a", "b"); err != nil {
		t.Fatal(err)
	}
	if err := b.Freeze("b"); err != nil {
		t.Fatal(err)
	}
	if err := b.SetNote("b", "changed"); !errors.Is(err, ErrFrozen) {
		t.Fatalf("frozen set note err=%v", err)
	}
}

func TestValidateRejectsBadEntries(t *testing.T) {
	bad := []Entry{
		{ID: "", Spec: input.Spec{WMax: 1, C: 1, RTT: 1, T: 1}, Win: cubic.Result{W: 1, K: 1}},
		{ID: "x", Spec: input.Spec{WMax: 0, C: 1, RTT: 1, T: 1}, Win: cubic.Result{W: 1, K: 1}},
		{ID: "x", Spec: input.Spec{WMax: 1, C: 1, RTT: 1, T: 1}, Win: cubic.Result{W: -1, K: 1}},
	}
	for i, e := range bad {
		if err := e.Validate(); err == nil {
			t.Fatalf("entry %d should fail", i)
		}
	}
}

func TestDerivedStats(t *testing.T) {
	b := NewBook(16)
	for _, e := range []Entry{
		{ID: "a", Spec: input.Spec{WMax: 16, C: 0.4, RTT: 0.1, T: 0.4}, Win: cubic.Result{W: 17, K: 0.8, Friendly: true}},
		{ID: "b", Spec: input.Spec{WMax: 200, C: 0.4, RTT: 0.2, T: 1}, Win: cubic.Result{W: 210, K: 3, Friendly: false, FastConv: true}},
		{ID: "c", Spec: input.Spec{WMax: 16, C: 0.4, RTT: 0.1, T: 0.2}, Win: cubic.Result{W: 14, K: 0.8, Friendly: true}},
	} {
		if err := b.Add(e); err != nil {
			t.Fatal(err)
		}
	}
	avg, n := b.AverageWindow()
	if n != 3 || avg != 80.33333333333333 {
		t.Fatalf("avg=%v n=%d", avg, n)
	}
	maxW, id := b.MaxWindow()
	if id != "b" || maxW != 210 {
		t.Fatalf("max=%v id=%s", maxW, id)
	}
	if b.FriendlyCount() != 2 || b.FastConvCount() != 1 {
		t.Fatalf("friendly=%d fast=%d", b.FriendlyCount(), b.FastConvCount())
	}
	sim := b.Similar(input.Spec{WMax: 16, C: 0.4, RTT: 0.1, T: 0.4}, 0.05)
	if len(sim) != 2 {
		t.Fatalf("similar=%+v", sim)
	}
	meanK, mk := b.MeanK()
	if mk != 3 || meanK < 1.5333333333333330 || meanK > 1.5333333333333334 {
		t.Fatalf("meanK=%v n=%d", meanK, mk)
	}
}
