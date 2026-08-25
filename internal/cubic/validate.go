package cubic

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Fields []string
}

func (e *ValidationError) Error() string {
	if e == nil {
		return "cubic: <nil> validation error"
	}
	return "cubic: invalid parameters: " + strings.Join(e.Fields, "; ")
}

func (e *ValidationError) Add(rule, want, got string) {
	e.Fields = append(e.Fields, fmt.Sprintf("%s (want %s, got %s)", rule, want, got))
}

func (e *ValidationError) Empty() bool {
	return len(e.Fields) == 0
}

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
	return bindBadParams(&ve)
}
