package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestNewTableOrError(t *testing.T) {
	tests := []struct {
		name          string
		args          query.TableName
		wantErr       bool
		wantTableName query.TableName
		wantAlias     query.TableName
	}{
		{"ok_table_name", "baseTable", false, "baseTable", ""},
		{"ok_name_as_alias", "baseTable as t1", false, "baseTable", "t1"},
		{"ok_name_alias", "baseTable t1", false, "baseTable", "t1"},
		{"nok_unexpected_format", "baseTable t1 t2", true, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := query.TableOrError(tt.args)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.wantTableName, got.TableName())
			require.Equal(t, tt.wantAlias, got.Alias())
		})
	}
}

func TestTable(t *testing.T) {
	tests := []struct {
		name          string
		args          query.TableName
		wantPanic     bool
		wantTableName query.TableName
		wantAlias     query.TableName
	}{
		{"ok_table_name", "baseTable", false, "baseTable", ""},
		{"ok_name_as_alias", "baseTable as t1", false, "baseTable", "t1"},
		{"ok_name_alias", "baseTable t1", false, "baseTable", "t1"},
		{"nok_unexpected_format", "baseTable t1 t2", true, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				require.Panics(t, func() {
					query.Table(tt.args)
				})
				return
			}

			require.NotPanics(t, func() {
				got := query.Table(tt.args)
				require.Equal(t, tt.wantTableName, got.TableName())
				require.Equal(t, tt.wantAlias, got.Alias())
			})
		})
	}
}

func TestTable_As(t *testing.T) {
	tests := []struct {
		name     string
		initial  query.TableIdent
		setAlias query.TableName
		as       query.TableName
	}{
		{"set_alias", query.Table("test"), "one", "one"},
		{"reset_alias", query.Table("test").As("one"), "two", "two"},
		{"set_empty_alias", query.Table("test").As("one"), "", ""},
		{"set_alias_same_as_table_name", query.Table("test"), "test", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.as, tt.initial.As(tt.setAlias).Alias())
		})
	}
}

func TestTable_RenderFrom(t *testing.T) {
	tests := []struct {
		name  string
		table query.TableIdent
		want  string
	}{
		{"name", query.Table("table"), "table"},
		{"name_and_alias", query.Table("table").As("t1"), "table AS t1"},
		{"name_only_if_alias_the_same", query.Table("table").As("table"), "table"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.table.RenderFrom(), tt.want)
		})
	}
}

func TestTable_Select(t *testing.T) {
	tests := []struct {
		name          string
		table         query.TableIdent
		wantTableName query.TableName
		wantSQL       string
	}{
		{"name", query.Table("table1"), "table1", "SELECT * FROM table1"},
		{"name_as", query.Table("table1").As("t1"), "table1", "SELECT * FROM table1 AS t1"},
		{"name_as_internal", query.Table("table1 as t1"), "table1", "SELECT * FROM table1 AS t1"},
		{"name_AS_internal", query.Table("table1 AS t1"), "table1", "SELECT * FROM table1 AS t1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.table.Select()
			require.Equal(t, tt.wantTableName, builder.TableName())
			sql, _, err := builder.BuildQueryAndParams()
			require.NoError(t, err)
			require.Equal(t, tt.wantSQL, sql)
		})
	}
}

func TestTableOrError(t *testing.T) {
	tests := []struct {
		name      string
		args      any
		wantName  query.TableName
		wantAlias query.TableName
		wantErr   bool
	}{
		{"ok_str", "table", "table", "", false},
		{"ok_str_as", "table as t1", "table", "t1", false},
		{"nok_str_err", "table not t1", "", "", true},
		{"ok_name", query.TableName("table"), "table", "", false},
		{"ok_name_as", query.TableName("table as t1"), "table", "t1", false},
		{"nok_name_err", query.TableName("table not t1"), "", "", true},
		{"ok_ident", query.Table("table"), "table", "", false},
		{"ok_ident_as_int", query.Table("table as t1"), "table", "t1", false},
		{"ok_ident_as_ext", query.Table("table").As("t1"), "table", "t1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got *query.TableIdent
				err error
			)

			switch typed := tt.args.(type) {
			case string:
				got, err = query.TableOrError(typed)
			case query.TableName:
				got, err = query.TableOrError(typed)
			case query.TableIdent:
				got, err = query.TableOrError(typed)
			}

			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.NotNil(t, got)
			require.Equal(t, tt.wantName, got.TableName())
			require.Equal(t, tt.wantAlias, got.Alias())
		})
	}
}
