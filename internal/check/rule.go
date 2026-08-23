// Package check verifies the defining cross-rules of the CUBIC kernel
// against a given spec and produces a PASS/FAIL report. The rules pin
// the identities that must hold: t=K returns to WMax, a larger WMax
// lengthens K, doubling C shortens K, low BDP follows the TCP-friendly
// branch, the window dips below WMax after a loss and recovers, and a
// pure Reno run increments by exactly one segment per RTT.
package check

// Result is one rule's verdict.
type Result struct {
	// Name is a short human-readable rule name.
	Name string
	// Detail carries the observed numbers for the verdict.
	Detail string
	// Pass is the rule's outcome.
	Pass bool
}

// String renders the verdict with a PASS/FAIL mark.
func (r Result) String() string {
	mark := "PASS"
	if !r.Pass {
		mark = "FAIL"
	}
	return mark + "  " + r.Name + ": " + r.Detail
}
