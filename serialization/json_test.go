package serialization_test

import (
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/stretchr/testify/assert"

	. "github.com/cybersamx/golib/serialization"
)

func TestParseJSON(t *testing.T) {
	t.Parallel()

	type person struct {
		Name string
		Age  int
	}

	t.Run("With non-struct type", func(t *testing.T) {
		json := "mike"
		_, err := ParseJSON[string](strings.NewReader(json))
		assert.Error(t, err)
	})

	t.Run("With struct value type", func(t *testing.T) {
		json := `{"name": "mike", "age": 25}`

		wantPerson := person{
			Name: "mike",
			Age:  25,
		}

		obj, err := ParseJSON[person](strings.NewReader(json))
		assert.NoError(t, err)
		assert.Equal(t, wantPerson, obj)
	})

	t.Run("With struct pointer type", func(t *testing.T) {
		json := `{"name": "mike", "age": 25}`

		wantPersonPtr := &person{
			Name: "mike",
			Age:  25,
		}

		objPtr, err := ParseJSON[*person](strings.NewReader(json))
		assert.NoError(t, err)
		assert.Equal(t, wantPersonPtr, objPtr)
	})

	t.Run("With map[string]any type", func(t *testing.T) {
		json := `{"name": "mike", "age": 25}`

		wantPersonMap := map[string]any{
			"name": "mike",
			"age":  25,
		}

		objMap, err := ParseJSON[map[string]any](strings.NewReader(json))
		assert.NoError(t, err)
		diff := pretty.Compare(wantPersonMap, objMap)
		assert.Empty(t, diff)
	})
}
