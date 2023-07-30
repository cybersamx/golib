package reflectutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/cybersamx/golib/reflectutils"
)

func TestIndirect(t *testing.T) {
	t.Parallel()

	inputs := []any{
		nil,
		3,
		"hello",
	}

	for _, input := range inputs {
		ptr := &input
		ptrr := &ptr
		ptrrr := &ptrr

		assert.IsType(t, input, Indirect(ptr))
		assert.IsType(t, input, Indirect(ptrr))
		assert.IsType(t, input, Indirect(ptrrr))
	}
}
