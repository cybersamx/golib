package ioutils

func CloneBytes(buf []byte) []byte {
	if buf == nil {
		return nil
	}

	clone := make([]byte, len(buf))
	copy(clone, buf)

	return clone
}
