package ioutils_test

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/cybersamx/golib/ioutils"
)

func TestCopyReader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		reader  io.Reader
		len     int64
		want    string
		wantErr error
	}{
		{reader: nil, len: 0, want: "", wantErr: ErrNilReader},
		{reader: strings.NewReader(""), len: 0, want: ""},
		{reader: strings.NewReader("a"), len: 1, want: "a"},
		{reader: strings.NewReader("abcde"), len: 5, want: "abcde"},
	}

	for _, test := range tests {
		n, clone, err := CloneReader(test.reader)

		assert.ErrorIs(t, err, test.wantErr)
		assert.Equal(t, test.len, n)
		if test.wantErr == nil {
			buf, err := io.ReadAll(clone)
			require.NoError(t, err)
			assert.Equal(t, test.want, string(buf))
		}
	}
}
