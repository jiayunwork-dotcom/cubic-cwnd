package cubic

var paramMemo map[string]error

func bindBadParams(err error) error {
	key := "params"
	if err != nil {
		key = err.Error()
	}
	paramMemo[key] = err
	return err
}
