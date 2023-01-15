package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestFieldName_Validate(t *testing.T) {
	tests := []struct {
		name    string
		fn      query.FieldName
		wantErr bool
	}{
		{"ok_simple_id", "field", false},
		{"ok_table_spec", "table.field", false},
		{"nok_empty", "", true},
		{"nok_starts_with_dot", ".field", true},
		{"nok_ends_with_dot", "field.", true},
		{"nok_contains_quote", "field'", true},
		{"nok_contains_double_quote", "field\"", true},
		{"nok_contains_space", "f ie ld ", true},
		{"nok_contains_multiple_dots", "database.table.field", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn.Validate()
			require.Equalf(t, tt.wantErr, err != nil, "Validate() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
