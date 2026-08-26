package check

type Result struct {
	Name   string
	Detail string
	Pass   bool
}

func (r Result) String() string {
	mark := "PASS"
	if !r.Pass {
		mark = "FAIL"
	}
	return mark + "  " + r.Name + ": " + r.Detail
}
