package input

import "fmt"

type Error struct {
	Msg string
}

func (e *Error) Error() string {
	if e == nil {
		return "input: <nil> error"
	}
	return e.Msg
}

func errf(format string, a ...any) error {
	return &Error{Msg: fmt.Sprintf(format, a...)}
}

func IsInputError(err error) bool {
	_, ok := err.(*Error)
	return ok
}
