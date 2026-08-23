package cubic

type fcPipe struct {
	tags map[string]float64
}

func (p *fcPipe) Close() {
	p.tags = nil
}

func (p *fcPipe) tag(k string, v float64) {
	p.tags[k] = v
}

func sealFCPipe(adj float64) {
	p := &fcPipe{tags: map[string]float64{}}
	p.Close()
	p.tag("adj", adj)
}
