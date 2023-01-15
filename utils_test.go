package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestUniqFieldNames(t *testing.T) {
	tests := []struct {
		name    string
		initial []query.FieldName
		uniq    []query.FieldName
	}{
		{"empty_list", []query.FieldName{}, []query.FieldName{}},
		{"single_foo", []query.FieldName{"foo"}, []query.FieldName{"foo"}},
		{"single_empty", []query.FieldName{""}, []query.FieldName{""}},
		{"double_a", []query.FieldName{"a", "a"}, []query.FieldName{"a"}},
		{"many_ordered", []query.FieldName{"a", "a", "a", "b", "b", "b", "c", "c", "c"}, []query.FieldName{"a", "b", "c"}},
		{"many_mixed", []query.FieldName{"a", "a", "b", "c", "a", "b", "c", "c"}, []query.FieldName{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.uniq, query.UniqFieldNames(tt.initial))
		})
	}
}
