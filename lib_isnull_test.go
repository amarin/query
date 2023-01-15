package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func Test_nullValue_Render(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"field_is_null",
			query.IsNull("f1"),
			0, "f1 IS NULL", []interface{}{},
		},
		{"field_is_not_null",
			query.Not(query.IsNull("f1")),
			0, "f1 IS NOT NULL", []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.want, tt.cond.RenderSQL())
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}
