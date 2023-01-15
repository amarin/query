package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestTableName_Select(t *testing.T) {
	tests := []struct {
		name          string
		table         query.TableName
		wantTableName query.TableName
		wantSQL       string
	}{
		{"name", query.TableName("table1"), "table1", "SELECT * FROM table1"},
		{"name_alias", query.TableName("table1 as t1"), "table1", "SELECT * FROM table1 AS t1"},
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
