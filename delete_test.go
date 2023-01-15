package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestDeleteBuilder_Operation(t *testing.T) {
	tests := []struct {
		name string
		query.DeleteBuilder
	}{
		{"constructor", query.Delete("test")},
		{"via_table_name", query.TableName("t1").Delete()},
		{"via_table_ident", query.Table("t1").Delete()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, query.DoDelete, tt.Operation())
		})
	}
}

func TestTableDeleter_Render(t *testing.T) {
	tests := []struct {
		name       string
		query      query.DeleteBuilder
		wantSql    string
		wantParams []interface{}
		wantErr    bool
	}{
		{
			"error_on_no_conditions",
			query.Delete("t1"),
			"",
			[]interface{}{},
			true,
		},
		{
			"delete_by_one_field",
			query.Delete("test").Where(query.Contains("f1", "a1")),
			"DELETE FROM test WHERE f1 LIKE $1",
			[]interface{}{"%a1%"},
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

func TestDeleteHelper_Set(t *testing.T) {
	updater := query.Update("test")
	require.Empty(t, updater.SetValues())
	require.NotEmpty(t, updater.Set(*query.NewFieldValue("test", nil)).SetValues())
	require.Empty(t, updater.SetValues())
}
