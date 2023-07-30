package ioutils

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrNilReader  = errors.New("reader must not be nil")
	ErrBufferCopy = errors.New("failed to copy to buffer")
)

// CloneReader copy the content from an io.Reader and return a new io.Reader.
func CloneReader(reader io.Reader) (int64, io.Reader, error) {
	if reader == nil {
		return 0, nil, ErrNilReader
	}

	clone := new(bytes.Buffer)
	n, err := io.Copy(clone, reader)

	if err != nil {
		return 0, nil, ErrBufferCopy
	}

	return n, clone, nil
}
