package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func Test_FieldCondition_ApplyFieldTable(t *testing.T) {
	tests := []struct {
		name       string
		fields     query.Condition
		table      query.TableName
		wantString string
	}{
		{"contains", query.Contains("f1", "v1"), "test", "test.f1 LIKE $1"},
		{"equal", query.EqualTo("f1", "v1"), "test", "test.f1=$1"},
		{"is_null", query.IsNull("f1"), "test", "test.f1 IS NULL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated := tt.fields.ApplyFieldTable(tt.table)
			require.Equal(t, tt.wantString, updated.Render(0))
		})
	}
}
