package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValuesMap_FieldValues(t *testing.T) {
	t.Run("empty values map provides empty field values slice", func(t *testing.T) {
		source := make(ValuesMap)
		result, err := source.FieldValues()
		require.NoError(t, err)
		require.Len(t, result, 0)
	})
	t.Run("invalid field name returns error", func(t *testing.T) {
		source := ValuesMap{"field name": 0}
		result, err := source.FieldValues()
		require.Error(t, err)
		require.Nil(t, result)
	})
	t.Run("check with 1 valid key", func(t *testing.T) {
		source := ValuesMap{"some_field": 0}
		result, err := source.FieldValues()
		require.NoError(t, err)
		require.Len(t, result, 1)
	})
	t.Run("check with 2 valid keys", func(t *testing.T) {
		source := ValuesMap{"some_field": 0, "some_another_field": false}
		result, err := source.FieldValues()
		require.NoError(t, err)
		require.Len(t, result, 2)
	})
	t.Run("check with 3 valid keys", func(t *testing.T) {
		source := ValuesMap{"some_field": 0, "some_another_field": false, "once_more": "str"}
		result, err := source.FieldValues()
		require.NoError(t, err)
		require.Len(t, result, 3)
	})
}
