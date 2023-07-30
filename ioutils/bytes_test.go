package ioutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/cybersamx/golib/ioutils"
)

func TestCloneBytes(t *testing.T) {
	t.Parallel()

	bufs := [][]byte{
		nil,
		[]byte(""),
		[]byte("a"),
		[]byte("abc"),
	}

	for _, buf := range bufs {
		clone := CloneBytes(buf)

		assert.Equal(t, buf, clone)
	}
}
