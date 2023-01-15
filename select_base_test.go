package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestSelectFrom(t *testing.T) {
	tests := []struct {
		name          string
		args          any
		wantTableName query.TableName
		wantSQL       string
	}{
		{"str", "table1", "table1", "SELECT * FROM table1"},
		{"str_alias", "table1 as t1", "table1", "SELECT * FROM table1 AS t1"},
		{"name", query.TableName("table1"), "table1", "SELECT * FROM table1"},
		{"name_alias", query.TableName("table1 as t1"), "table1", "SELECT * FROM table1 AS t1"},
		{"table", query.TableName("table1"), "table1", "SELECT * FROM table1"},
		{"table_alias", query.Table("table1 as t1"), "table1", "SELECT * FROM table1 AS t1"},
		{"table_as", query.Table("table1").As("t1"), "table1", "SELECT * FROM table1 AS t1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got query.BaseSelectBuilder

			switch typed := tt.args.(type) {
			case string:
				got = query.SelectFrom(typed)
			case query.TableName:
				got = query.SelectFrom(typed)
			case query.TableIdent:
				got = query.SelectFrom(typed)
			}
			require.Equal(t, tt.wantTableName, got.TableName())
			sql, _, err := got.BuildQueryAndParams()
			require.NoError(t, err)
			require.Equal(t, tt.wantSQL, sql)
		})
	}
}

func TestBaseSelectBuilder_Fields(t *testing.T) {
	expectedTableName := "table13"
	builder := query.SelectFrom(expectedTableName)

	getSQL := func(b query.BaseSelectBuilder) string {
		sql, _, err := b.BuildQueryAndParams()
		require.NoError(t, err)
		return sql
	}

	t.Run("no_fields_by_default", func(t *testing.T) {
		require.Equal(t, "SELECT * FROM table13", getSQL(builder))
	})

	updated := builder.Fields(query.Field("field1"))
	t.Run("add_single_field_by_name", func(t *testing.T) {
		require.Equal(t, "SELECT field1 FROM table13", getSQL(updated))
	})

	updated = updated.Fields(query.Field("field2"), query.Field("field3 as f3"))
	t.Run("fields_resets_fields_list", func(t *testing.T) {
		require.Equal(t, "SELECT field2, field3 AS f3 FROM table13", getSQL(updated))
	})
}

func TestBaseSelectBuilder_Join(t *testing.T) {
	tests := []struct {
		name       string
		left       query.BaseSelectBuilder
		right      query.TableIdent
		joinType   query.TableJoinType
		leftField  query.FieldDefinition
		rightField query.FieldDefinition
		wantSQL    string
	}{
		{
			"base",
			query.SelectFrom("t1"),
			query.Table("t2"),
			query.InnerJoin,
			query.Field("t1.f1"),
			query.Field("t2.f2"),
			"SELECT * FROM t1 INNER JOIN t2 ON t1.f1=t2.f2",
		},
		{
			"left_alias_is_not_infer_join_condition",
			query.SelectFrom("t1"),
			query.Table("t2"),
			query.InnerJoin,
			query.Field("t1.f1 as field1"),
			query.Field("t2.f2"),
			"SELECT * FROM t1 INNER JOIN t2 ON t1.f1=t2.f2",
		},
		{
			"right_alias_is_not_infer_join_condition",
			query.SelectFrom("t1"),
			query.Table("t2"),
			query.InnerJoin,
			query.Field("t1.f1 as field1"),
			query.Field("t2.f2"),
			"SELECT * FROM t1 INNER JOIN t2 ON t1.f1=t2.f2",
		},
		{
			"with_fields",
			query.SelectFrom("t1").Fields(
				query.Field("t1.f1 as field1"),
				query.Field("t2.f13 as field13")),
			query.Table("t2"),
			query.InnerJoin,
			query.Field("t1.f1"),
			query.Field("t2.f2"),
			"SELECT t1.f1 AS field1, t2.f13 AS field13 FROM t1 INNER JOIN t2 ON t1.f1=t2.f2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.left.Join(tt.right, tt.joinType).On(tt.leftField, tt.rightField)
			sql, _, err := got.BuildQueryAndParams()
			require.NoError(t, err)
			require.Equal(t, tt.wantSQL, sql)
		})
	}
}

func TestBaseSelectBuilder_RenderSQL(t *testing.T) {
	tests := []struct {
		name       string
		builder    query.BaseSelectBuilder
		wantSql    string
		wantParams []any
	}{
		{"simple_query", query.SelectFrom("table1"), "SELECT * FROM table1", []any{}},
		{
			"fields",
			query.SelectFrom("table2").Fields(query.Field("f1"), query.Field("f2 as field2"), query.Field("f3")),
			"SELECT f1, f2 AS field2, f3 FROM table2",
			[]any{},
		},
		{
			"join",
			query.
				SelectFrom("table3 as t3").
				LeftJoin(query.Table("table4 as t4")).On(
				query.Field("table3.f1"),
				query.Field("table4.f2")).
				Fields(
					query.Field("t3.f1"),
					query.Field("t3.f2 as field2"),
					query.Field("t4.f3 as field3"),
				),
			"SELECT t3.f1, t3.f2 AS field2, t4.f3 AS field3 FROM table3 AS t3 LEFT JOIN table4 AS t4 ON table3.f1=table4.f2",
			[]any{},
		},
		{
			"join_where",
			query.
				SelectFrom("table3 as t3").
				LeftJoin(query.Table("table4 as t4")).On(
				query.Field("table3.f1"),
				query.Field("table4.f2")).
				Fields(
					query.Field("t3.f1"),
					query.Field("t3.f2 as field2"),
					query.Field("t4.f3 as field3"),
				).Where(
				query.FieldName("t3.f1").EqualTo(13),
			),
			"SELECT t3.f1, t3.f2 AS field2, t4.f3 AS field3 FROM table3 AS t3 LEFT JOIN table4 AS t4 ON table3.f1=table4.f2 WHERE t3.f1=?",
			[]any{13},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantSql, tt.builder.RenderSQL())
			require.Equal(t, tt.wantParams, tt.builder.Values())

		})
	}
}
