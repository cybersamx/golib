package stringsutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/cybersamx/golib/stringsutils"
)

func TestRuneToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input rune
		want  string
	}{
		{rune(0), ""},
		{'a', "a"},
		{'b', "b"},
		{'C', "C"},
		{'D', "D"},
		{'1', "1"},
		{'2', "2"},
		{'0', "0"},
		{'!', "!"},
		{'?', "?"},
		{'∬', "∬"},
		{'∞', "∞"},
		{'✌', "✌"},
		{'❆', "❆"},
		{'ぁ', "ぁ"},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, RuneToString(test.input))
	}
}

func TestEllipsisString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		limit int
		want  string
	}{
		{"", 0, "..."},
		{"abcde", 0, "..."},
		{"", 1, "..."},
		{"abcde", 1, "..."},
		{"", 3, "..."},
		{"abcde", 3, "..."},
		{"", 10, ""},
		{"a", 10, "a"},
		{"abcde", 10, "abcde"},
		{"abcdefghijklmnopq", 10, "abcdefg..."},
	}

	for _, test := range tests {
		actual := EllipsisString(test.input, test.limit)
		assert.Equal(t, test.want, actual)
	}
}

func TestTruncateString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		limit int
		want  string
	}{
		{"", 0, ""},
		{"abcde", 0, ""},
		{"", 10, ""},
		{"a", 10, "a"},
		{"abcde", 10, "abcde"},
		{"abcdefghij", 10, "abcdefghij"},
	}

	for _, test := range tests {
		actual := TruncateString(test.input, test.limit)
		assert.Equal(t, test.want, actual)
	}
}
