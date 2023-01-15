package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestBaseBuilder_Operation(t *testing.T) {
	tests := []struct {
		name string
		op   query.BaseBuilder
		want query.Operation
	}{
		{"select", query.DoSelect.New(), query.DoSelect},
		{"update", query.DoUpdate.New(), query.DoUpdate},
		{"insert", query.DoInsert.New(), query.DoInsert},
		{"delete", query.DoDelete.New(), query.DoDelete},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.op.Operation())
		})
	}
}
