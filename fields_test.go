package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestFields_FieldList(t *testing.T) {
	tests := []struct {
		name   string
		fields query.Fields
		want   string
	}{
		{"empty", query.NewFields(), "*"},
		{"str", query.NewFields("test"), "test"},
		{"str_as", query.NewFields("test as t1"), "test AS t1"},
		{"str_table_as", query.NewFields("table.test as t1"), "table.test AS t1"},
		{"name", query.NewFields(query.FieldName("test")), "test"},
		{"name_as", query.NewFields(query.FieldName("test as t1")), "test AS t1"},
		{"name_table_as", query.NewFields(query.FieldName("table.test as t1")), "table.test AS t1"},
		{"name_table_as_x2",
			query.NewFields(query.FieldName("t1.f1 as s1"), query.FieldName("t2.f2 as s2")),
			"t1.f1 AS s1, t2.f2 AS s2"},
		{"field", query.NewFields(query.Field("test")), "test"},
		{"field_as", query.NewFields(query.Field("test").As("t1")), "test AS t1"},
		{"field_table_as", query.NewFields(query.Field("test").As("t1").Of("table")), "table.test AS t1"},
		{"mix_str_and_name", query.NewFields("test", query.FieldName("test1")), "test, test1"},
		{"mix_str_and_field", query.NewFields("test", query.Field("test1")), "test, test1"},
		{"mix_name_and_field", query.NewFields(query.FieldName("test"), query.Field("test1")), "test, test1"},
		{"mix_str_name_and_field",
			query.NewFields(
				"f1",
				query.FieldName("table1.field2"),
				query.Field("field3").Of("table2").As("f3"),
			),
			"f1, table1.field2, table2.field3 AS f3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.fields.FieldList())
		})
	}
}
