package system_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/cybersamx/golib/system"
)

func TestIsFileExist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		filePath string
		want     bool
	}{
		{"filesystem.go", true},
		{"../go.mod", true},
		{"not_exist.json", false},
	}

	for _, test := range tests {
		result := IsFileExist(test.filePath)
		assert.Equal(t, test.want, result)
	}
}
