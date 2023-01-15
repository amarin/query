package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestOrder(t *testing.T) {
	tests := []struct {
		name         string
		fieldName    query.FieldName
		direction    []query.SortDirection
		expectRender string
		expectsPanic bool
	}{
		{"ok_field_asc", "field", []query.SortDirection{query.Ascending}, "field ASC", false},
		{"ok_field_desc", "field", []query.SortDirection{query.Descending}, "field DESC", false},
		{"ok_field_spec_asc", "table.field", []query.SortDirection{query.Ascending}, "table.field ASC", false},
		{"ok_field_spec_desc", "table.field", []query.SortDirection{query.Descending}, "table.field DESC", false},
		{"ok_field_default_asc", "field", nil, "field ASC", false},
		{"panics_unknown_direction", "field", []query.SortDirection{"here-and-there"}, "", true},
		{"panics_field_name_starts_with_dot", ".field", []query.SortDirection{query.Descending}, "", true},
		{"panics_field_name_ends_with_dot", "field.", []query.SortDirection{query.Descending}, "", true},
		{"panics_field_name_with_quote", "field'", []query.SortDirection{query.Descending}, "", true},
		{"panics_field_name_with_double_quote", "field\"", []query.SortDirection{query.Descending}, "", true},
		{"panics_field_name_contains_2_dots", "database.table.field", []query.SortDirection{query.Descending}, "", true},
		{"panics_field_name_empty", "", []query.SortDirection{query.Descending}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rendered string

			switch {
			case tt.expectsPanic:
				require.Panics(t, func() {
					rendered = query.OrderBy(tt.fieldName, tt.direction...).Render()
				})
			default:
				require.NotPanics(t, func() {
					rendered = query.OrderBy(tt.fieldName, tt.direction...).Render()
				})
			}
			require.Equal(t, tt.expectRender, rendered)
		})
	}
}

func TestFieldSorting_ApplyFieldSpec(t *testing.T) {
	tests := []struct {
		name   string
		fields query.FieldSorting
		spec   query.FieldDefinition
		want   string
	}{
		{"ok_plain_name",
			query.OrderBy("f1"), query.Field("f1"),
			"f1 ASC"},
		{"ok_table_field_name",
			query.OrderBy("f1"), query.Field("f1").Of("table"),
			"table.f1 ASC"},
		{"ok_alias_name",
			query.OrderBy("f1"), query.Field("f1").As("field1"),
			"field1 ASC"},
		{"ok_table_but_alias",
			query.OrderBy("f1"), query.Field("f1").Of("table").As("field1"),
			"field1 ASC"},
		{"nok_not_matched",
			query.OrderBy("f1"), query.Field("f21").As("field1"),
			"f1 ASC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated := tt.fields.ApplyFieldSpec(tt.spec)
			require.Equal(t, tt.want, updated.Render())
		})
	}
}
