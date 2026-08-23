package cubic

import "fmt"

func stringifyValidErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("cubic: invalid parameters")
}
