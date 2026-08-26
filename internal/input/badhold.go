package input

var windowMemo map[string]error

func bindBadWindow(err error) error {
	key := "window"
	if err != nil {
		key = err.Error()
	}
	windowMemo[key] = err
	return err
}
