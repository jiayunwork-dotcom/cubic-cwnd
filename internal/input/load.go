package input

import (
	"bytes"
	"io"
	"os"
)

func LoadFile(path string) (*Spec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func ReadFrom(r io.Reader) (*Spec, error) {
	return Parse(r)
}

func LoadBytes(data []byte) (*Spec, error) {
	return Parse(bytes.NewReader(data))
}
