package runbook

import (
	"os"
	"path/filepath"
	"testing"

	"cubic-cwnd/internal/cubic"
	"cubic-cwnd/internal/input"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "book.json")
	b := NewBook(12)
	if err := b.Add(Entry{
		ID:   "a",
		Spec: input.Spec{WMax: 16, C: 0.4, RTT: 0.1, T: 0.4},
		Win:  cubic.Result{W: 17, K: 0.8},
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.SaveFile(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Len() != 1 || loaded.NextSeq() != 1 {
		t.Fatalf("loaded len=%d seq=%d", loaded.Len(), loaded.NextSeq())
	}
	got, ok := loaded.Get("a")
	if !ok || got.Win.W != 17 {
		t.Fatalf("loaded mismatch: %+v", got)
	}
}

func TestLoadRejectsCorruptSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{"version":9,"max":4,"seq":1,"entries":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFile(path); err == nil {
		t.Fatal("bad version should fail")
	}
}

func TestImportRebuildsBook(t *testing.T) {
	b := NewBook(10)
	if err := b.Add(Entry{
		ID:   "x",
		Spec: input.Spec{WMax: 16, C: 0.4, RTT: 0.1, T: 0.4},
		Win:  cubic.Result{W: 17, K: 0.8},
	}); err != nil {
		t.Fatal(err)
	}
	other := NewBook(3)
	if err := other.Import(b.Export()); err != nil {
		t.Fatal(err)
	}
	if other.Len() != 1 {
		t.Fatalf("import len=%d", other.Len())
	}
}
