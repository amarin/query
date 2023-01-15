package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func Test_iContains_Render(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"field_is_like",
			query.IContains("f1", "aaa"),
			0, "f1 ILIKE $1", []interface{}{"%aaa%"},
		},
		{"field_is_like_table",
			query.IContains("f1", "aaa").ApplyFieldTable("test"),
			0, "test.f1 ILIKE $1", []interface{}{"%aaa%"},
		},
		{"field_is_not_like",
			query.Not(query.IContains("f1", "aaa")),
			0, "f1 NOT ILIKE $1", []interface{}{"%aaa%"},
		},
		{"field_is_not_like_table",
			query.Not(query.IContains("f1", "aaa")).ApplyFieldTable("test"),
			0, "test.f1 NOT ILIKE $1", []interface{}{"%aaa%"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}

func Test_iContains_RenderSQL(t *testing.T) {
	tests := []struct {
		name   string
		cond   query.Condition
		want   string
		values []interface{}
	}{
		{"field_is_like",
			query.IContains("f1", "aaa"),
			"f1 ILIKE ?", []interface{}{"%aaa%"},
		},
		{"field_is_like_table",
			query.IContains("f1", "aaa").ApplyFieldTable("test"),
			"test.f1 ILIKE ?", []interface{}{"%aaa%"},
		},
		{"field_is_not_like",
			query.Not(query.IContains("f1", "aaa")),
			"f1 NOT ILIKE ?", []interface{}{"%aaa%"},
		},
		{"field_is_not_like_table",
			query.Not(query.IContains("f1", "aaa")).ApplyFieldTable("test"),
			"test.f1 NOT ILIKE ?", []interface{}{"%aaa%"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.RenderSQL())
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}
