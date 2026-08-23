package cubic

// DefaultMSSBytes is the common 1460-byte MSS.
const DefaultMSSBytes = 1460

// SegmentsToBytes converts a window in MSS segments to bytes.
func SegmentsToBytes(segments, mssBytes float64) float64 {
	if mssBytes <= 0 {
		mssBytes = DefaultMSSBytes
	}
	return segments * mssBytes
}

// BytesToSegments converts a byte window to MSS segments.
func BytesToSegments(bytes, mssBytes float64) float64 {
	if mssBytes <= 0 {
		mssBytes = DefaultMSSBytes
	}
	return bytes / mssBytes
}

// WindowBytes returns the effective window of a result in bytes.
func (r Result) WindowBytes(mssBytes float64) float64 {
	return SegmentsToBytes(r.W, mssBytes)
}

// WMaxBytes returns the reference window in bytes.
func (r Result) WMaxBytes(mssBytes float64) float64 {
	return SegmentsToBytes(r.WMaxRef, mssBytes)
}
