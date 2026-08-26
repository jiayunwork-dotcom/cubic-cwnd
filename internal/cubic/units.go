package cubic

const DefaultMSSBytes = 1460

func SegmentsToBytes(segments, mssBytes float64) float64 {
	if mssBytes <= 0 {
		mssBytes = DefaultMSSBytes
	}
	return segments * mssBytes
}

func BytesToSegments(bytes, mssBytes float64) float64 {
	if mssBytes <= 0 {
		mssBytes = DefaultMSSBytes
	}
	return bytes / mssBytes
}

func (r Result) WindowBytes(mssBytes float64) float64 {
	return SegmentsToBytes(r.W, mssBytes)
}

func (r Result) WMaxBytes(mssBytes float64) float64 {
	return SegmentsToBytes(r.WMaxRef, mssBytes)
}
