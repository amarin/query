package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestUpdateBuilder_Operation(t *testing.T) {
	tests := []struct {
		name string
		query.UpdateBuilder
	}{
		{"constructor", query.Update("test")},
		{"via_table_name", query.TableName("t1").Update()},
		{"via_table_ident", query.Table("t1").Update()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, query.DoUpdate, tt.Operation())
		})
	}
}

func TestTableUpdater_Render(t *testing.T) {
	tests := []struct {
		name       string
		query      query.UpdateBuilder
		wantSql    string
		wantParams []interface{}
		wantErr    bool
	}{
		{
			"error_on_no_conditions",
			query.Update("t1").Set(query.FieldName("intField").Value(1)),
			"",
			[]interface{}{},
			true,
		},
		{
			"set_one_field",
			query.Update("test").Where(query.Contains("f1", "a1")).
				Set(query.FieldName("intField").Value(1)),
			"UPDATE test SET intField=$1 WHERE f1 LIKE $2",
			[]interface{}{1, "%a1%"},
			false,
		},
		{
			"set_two_fields_where_field_contains",
			query.SelectManyFrom("test").
				Where(query.Contains("f1", "a1")).
				Update(query.FieldName("intField").Value(1)).
				Set(query.FieldName("strField").Value("aaa")),
			"UPDATE test SET intField=$1, strField=$2 WHERE f1 LIKE $3",
			[]interface{}{1, "aaa", "%a1%"},
			false,
		},
		{
			"set_three_fields_where_field_equals",
			query.SelectManyFrom("test").
				Where(query.FieldName("UUID").EqualTo("00-000-00")).
				Update(query.FieldName("intField").Value(1)).
				Set(query.FieldName("boolField").Value(false)).
				Set(query.FieldName("strField").Value("aaa")),
			"UPDATE test SET intField=$1, boolField=$2, strField=$3 WHERE UUID=$4",
			[]interface{}{1, false, "aaa", "00-000-00"},
			false,
		},
		{
			"set_two_fields_with_and_condition",
			query.Update("test").
				Where(query.FieldName("UUID").EqualTo("00-000-00")).
				Where(query.FieldName("name").EqualTo("testName")).
				Set(query.FieldName("boolField").Value(false)).
				Set(query.FieldName("strField").Value("ccc")),
			"UPDATE test SET boolField=$1, strField=$2 WHERE UUID=$3 AND name=$4",
			[]interface{}{false, "ccc", "00-000-00", "testName"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotParams, err := tt.query.BuildQueryAndParams()
			require.Equalf(t, tt.wantErr, err != nil, "want error %v, got %v", tt.wantErr, err)
			if err != nil {
				return
			}
			require.Equal(t, tt.wantSql, gotSQL)
			require.Equal(t, tt.wantParams, gotParams)
		})
	}
}

func TestUpdateHelper_Set(t *testing.T) {
	updater := query.Update("test")
	require.Empty(t, updater.SetValues())
	require.NotEmpty(t, updater.Set(*query.NewFieldValue("test", nil)).SetValues())
	require.Empty(t, updater.SetValues())
}
