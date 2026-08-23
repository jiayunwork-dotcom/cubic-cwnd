package input

import "fmt"

// Error is a typed input error carrying a human-readable message that
// the CLI prints verbatim to stderr.
type Error struct {
	Msg string
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return "input: <nil> error"
	}
	return e.Msg
}

// errf builds a typed input error.
func errf(format string, a ...any) error {
	return &Error{Msg: fmt.Sprintf(format, a...)}
}

// IsInputError reports whether an error originated from input parsing.
func IsInputError(err error) bool {
	_, ok := err.(*Error)
	return ok
}
