package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestSelectOneFrom(t *testing.T) {
	tests := []struct {
		name          string
		args          any
		wantTableName query.TableName
		wantSQL       string
	}{
		{"str", "table1", "table1", "SELECT * FROM table1 LIMIT 1"},
		{"str_alias", "table1 as t1", "table1", "SELECT * FROM table1 AS t1 LIMIT 1"},
		{"name", query.TableName("table1"), "table1", "SELECT * FROM table1 LIMIT 1"},
		{"name_alias", query.TableName("table1 as t1"), "table1", "SELECT * FROM table1 AS t1 LIMIT 1"},
		{"table", query.TableName("table1"), "table1", "SELECT * FROM table1 LIMIT 1"},
		{"table_alias", query.Table("table1 as t1"), "table1", "SELECT * FROM table1 AS t1 LIMIT 1"},
		{"table_as", query.Table("table1").As("t1"), "table1", "SELECT * FROM table1 AS t1 LIMIT 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got query.SelectSingleBuilder

			switch typed := tt.args.(type) {
			case string:
				got = query.SelectSingleFrom(typed)
			case query.TableName:
				got = query.SelectSingleFrom(typed)
			case query.TableIdent:
				got = query.SelectSingleFrom(typed)
			}
			require.Equal(t, tt.wantTableName, got.TableName())
			sql, _, err := got.BuildQueryAndParams()
			require.NoError(t, err)
			require.Equal(t, tt.wantSQL, sql)
		})
	}
}
