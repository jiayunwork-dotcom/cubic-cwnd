package sim

import "fmt"

type Mode int

const (
	ModeCubic Mode = iota
	ModeReno
)

func (m Mode) String() string {
	if m == ModeReno {
		return "reno"
	}
	return "cubic"
}

func ParseMode(s string) (Mode, error) {
	switch s {
	case "", "cubic":
		return ModeCubic, nil
	case "reno":
		return ModeReno, nil
	default:
		return ModeCubic, fmt.Errorf("sim: unknown mode %q (want cubic or reno)", s)
	}
}
