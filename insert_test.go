package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestInsertBuilder_Operation(t *testing.T) {
	tests := []struct {
		name string
		query.InsertBuilder
	}{
		{"constructor", query.InsertInto("test")},
		{"via_table_name", query.TableName("t1").InsertInto()},
		{"via_table_ident", query.Table("t1").InsertInto()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, query.DoInsert, tt.Operation())
		})
	}
}

func TestInsertBuilder_BuildQueryAndParams(t *testing.T) {
	tests := []struct {
		name       string
		builder    query.InsertBuilder
		wantSql    string
		wantParams []interface{}
		wantErr    bool
	}{
		{
			"nok_no_fields",
			query.TableName("t1").InsertInto(),
			"",
			[]interface{}{},
			true,
		},
		{
			"nok_empty_table_name",
			query.TableName("").InsertInto().Values(query.FieldName("strField").Value("a1")),
			"",
			[]interface{}{},
			true,
		},
		{
			"ok_set_one_field",
			query.TableName("t1").InsertInto().Values(query.FieldName("strField").Value("a1")),
			"INSERT INTO t1(strField) VALUES ($1)",
			[]interface{}{"a1"},
			false,
		},
		{
			"ok_set_two_field",
			query.TableName("t1").InsertInto().Values(
				query.FieldName("intField").Value(123),
				query.FieldName("strField").Value("a1"),
			),
			"INSERT INTO t1(intField, strField) VALUES ($1, $2)",
			[]interface{}{123, "a1"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotParams, err := tt.builder.BuildQueryAndParams()
			require.Equalf(t, tt.wantErr, err != nil, "want error %v, got %v", tt.wantErr, err)
			if err != nil {
				return
			}

			require.Equal(t, tt.wantSql, gotSQL)
			require.Equal(t, tt.wantParams, gotParams)
		})
	}
}
