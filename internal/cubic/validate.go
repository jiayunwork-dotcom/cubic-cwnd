package cubic

import (
	"fmt"
	"strings"
)

// ValidationError collects every offending field of a Params value so
// the caller can see exactly which rule failed and why.
type ValidationError struct {
	Fields []string
}

// Error renders the validation failure as a readable message.
func (e *ValidationError) Error() string {
	if e == nil {
		return "cubic: <nil> validation error"
	}
	return "cubic: invalid parameters: " + strings.Join(e.Fields, "; ")
}

// Add records one offending rule with its expected and actual values.
func (e *ValidationError) Add(rule, want, got string) {
	e.Fields = append(e.Fields, fmt.Sprintf("%s (want %s, got %s)", rule, want, got))
}

// Empty reports whether no field was flagged.
func (e *ValidationError) Empty() bool {
	return len(e.Fields) == 0
}

// Validate returns an error when any of the hard rules is violated:
//   - WMax <= 0 (a zero or negative pre-loss window is rejected)
//   - C    <= 0 (a zero or negative scaling factor is rejected)
//   - RTT  <= 0 (a zero or negative RTT is rejected)
//   - T    <  0 (negative elapsed time is rejected)
//
// The checks are independent and all violations are reported together.
func (p Params) Validate() error {
	var ve ValidationError
	if p.WMax <= 0 {
		ve.Add("WMax must be positive", "> 0", fmt.Sprintf("%v", p.WMax))
	}
	if p.C <= 0 {
		ve.Add("C must be positive", "> 0", fmt.Sprintf("%v", p.C))
	}
	if p.RTT <= 0 {
		ve.Add("RTT must be positive", "> 0", fmt.Sprintf("%v", p.RTT))
	}
	if p.T < 0 {
		ve.Add("T must not be negative", ">= 0", fmt.Sprintf("%v", p.T))
	}
	if ve.Empty() {
		return nil
	}
	return &ve
}
