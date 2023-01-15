package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestTableQuery_Render(t *testing.T) {
	tests := []struct {
		name       string
		query      query.SelectManyBuilder
		wantSql    string
		wantParams []interface{}
		wantErr    bool
	}{
		{
			"base",
			query.SelectManyFrom("test"),
			"SELECT * FROM test",
			[]any{},
			false,
		},
		{
			"inner_join",
			query.SelectManyFrom("test").
				InnerJoin("t2").On(query.Field("test.f1"), query.Field("t2.f1")),
			"SELECT * FROM test INNER JOIN t2 ON test.f1=t2.f1",
			[]any{},
			false,
		},
		{
			"left_join",
			query.SelectManyFrom("test").
				LeftJoin("t2").On(query.Field("test.f1"), query.Field("t2.f1")),
			"SELECT * FROM test LEFT JOIN t2 ON test.f1=t2.f1",
			[]any{},
			false,
		},
		{
			"right_join",
			query.SelectManyFrom("test").
				RightJoin("t2").On(query.Field("test.f1"), query.Field("t2.f1")),
			"SELECT * FROM test RIGHT JOIN t2 ON test.f1=t2.f1",
			[]any{},
			false,
		},
		{
			"full_join",
			query.SelectManyFrom("test").
				FullJoin("t2").On(query.Field("test.f1"), query.Field("t2.f1")),
			"SELECT * FROM test FULL JOIN t2 ON test.f1=t2.f1",
			[]any{},
			false,
		},
		{
			"base_with_field_specs",
			query.SelectManyFrom("test").Fields(
				query.Field("f1"),
				query.Field("field2").As("f2"),
			),
			"SELECT f1, field2 AS f2 FROM test",
			[]interface{}{},
			false,
		},
		{
			"join_with_field_specs",
			query.
				SelectManyFrom("t1").
				InnerJoin("t2").
				On(
					query.Field("t1.f1"),
					query.Field("t2.field2")).
				Fields(
					query.Field("t1.f1").As("f1"),
					query.Field("t2.field2").As("f2"),
				),
			"SELECT t1.f1 AS f1, t2.field2 AS f2 FROM t1 INNER JOIN t2 ON t1.f1=t2.field2",
			[]interface{}{},
			false,
		},
		{
			"order_by",
			query.SelectManyFrom("test").OrderBy(query.ASC("field1")),
			"SELECT * FROM test ORDER BY field1 ASC",
			[]interface{}{},
			false,
		},
		{
			"order_by_specs",
			query.SelectManyFrom("test").OrderBy(query.ASC("field1")).Fields(query.Field("field1").As("f1")),
			"SELECT field1 AS f1 FROM test ORDER BY field1 ASC",
			[]interface{}{},
			false,
		},
		{
			"order_by_many",
			query.SelectManyFrom("test").OrderBy(query.ASC("field1"), query.DESC("field2")),
			"SELECT * FROM test ORDER BY field1 ASC, field2 DESC",
			[]interface{}{},
			false,
		},
		{
			"order_by_many_with_table_before",
			query.SelectManyFrom("test").OrderBy(query.ASC("test.field1"), query.DESC("test.field2")),
			"SELECT * FROM test ORDER BY test.field1 ASC, test.field2 DESC",
			[]interface{}{},
			false,
		},
		{
			"limit",
			query.SelectManyFrom("test").Limit(1),
			"SELECT * FROM test LIMIT $1",
			[]interface{}{uint(1)},
			false,
		},
		{
			"offset",
			query.SelectManyFrom("test").Offset(1),
			"SELECT * FROM test OFFSET $1",
			[]interface{}{uint(1)},
			false,
		},
		{
			"offset_limit",
			query.SelectManyFrom("test").Limit(1).Offset(2),
			"SELECT * FROM test OFFSET $1 LIMIT $2",
			[]interface{}{uint(2), uint(1)},
			false,
		},
		{
			"order_offset_limit",
			query.SelectManyFrom("test").
				OrderBy(query.ASC("field1")).
				Limit(4).
				Offset(3).
				OrderBy(query.DESC("field2")),
			"SELECT * FROM test ORDER BY field1 ASC, field2 DESC OFFSET $1 LIMIT $2",
			[]interface{}{uint(3), uint(4)},
			false,
		},
		{
			"where_contains_order_offset_limit",
			query.SelectManyFrom("test").
				Where(query.Contains("f1", "a1")).
				Where(query.Contains("f2", "b2")).
				OrderBy(query.ASC("f3")).
				Limit(6).
				Offset(5).
				OrderBy(query.DESC("f4")),
			"SELECT * FROM test WHERE f1 LIKE $1 AND f2 LIKE $2 ORDER BY f3 ASC, f4 DESC OFFSET $3 LIMIT $4",
			[]interface{}{"%a1%", "%b2%", uint(5), uint(6)},
			false,
		},
		{
			"where_equals_order_offset_limit",
			query.SelectManyFrom("test").
				Where(query.EqualTo("f1", 13)).
				OrderBy(query.ASC("f2")).
				Limit(6).
				Offset(5),
			"SELECT * FROM test WHERE f1=$1 ORDER BY f2 ASC OFFSET $2 LIMIT $3",
			[]interface{}{13, uint(5), uint(6)},
			false,
		},
		{
			"where_single_not",
			query.SelectManyFrom("test").Where(query.Not(query.EqualTo("f1", 13))),
			"SELECT * FROM test WHERE NOT f1=$1",
			[]interface{}{13},
			false,
		},
		{
			"where_not_equal_or_another_equal",
			query.SelectManyFrom("test").Where(
				query.Not(query.EqualTo("f1", 11)),
				query.Or(query.EqualTo("f2", 13)),
			),
			"SELECT * FROM test WHERE NOT f1=$1 OR f2=$2",
			[]interface{}{11, 13},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotParams, err := tt.query.BuildQueryAndParams()
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.wantSql, gotSQL)
			require.Equal(t, tt.wantParams, gotParams)
		})
	}
}

func TestBuilder_Fields(t *testing.T) {
	expectedTableName := "table13"
	builder := query.SelectManyFrom(expectedTableName)
	t.Run("no_field_list_initially", func(t *testing.T) {
		require.Empty(t, builder.FieldDefinitions())
		sql, _, err := builder.BuildQueryAndParams()
		require.NoError(t, err)
		require.Equal(t, "SELECT * FROM table13", sql)
	})

	t.Run("use_Fields", func(t *testing.T) {
		updated := builder.Fields(query.Field("f1").As("a1"))
		require.Len(t, updated.FieldDefinitions(), 1)
		sql, _, err := updated.BuildQueryAndParams()
		require.NoError(t, err)
		require.Equal(t, "SELECT f1 AS a1 FROM table13", sql)
	})
}

func TestBuilder_Operation(t *testing.T) {
	require.Equal(t, query.DoSelect, query.SelectManyFrom("test").Operation())
}

func TestSelectManyFrom(t *testing.T) {
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
			var got query.SelectManyBuilder

			switch typed := tt.args.(type) {
			case string:
				got = query.SelectManyFrom(typed)
			case query.TableName:
				got = query.SelectManyFrom(typed)
			case query.TableIdent:
				got = query.SelectManyFrom(typed)
			}
			require.Equal(t, tt.wantTableName, got.TableName())
			sql, _, err := got.BuildQueryAndParams()
			require.NoError(t, err)
			require.Equal(t, tt.wantSQL, sql)
		})
	}
}
