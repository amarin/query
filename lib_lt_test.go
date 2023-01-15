package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func Test_Less_Render(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"lt_string",
			query.Less("f1", "aaa"),
			0, "f1<$1", []interface{}{"aaa"},
		},
		{"lt_int",
			query.Less("f1", 1),
			1, "f1<$2", []interface{}{1},
		},
		{"not_lt",
			query.Not(query.Less("f1", 2)),
			2, "NOT f1<$3", []interface{}{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}

func Test_Less_And(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"eq_and_eq",
			query.Less("f1", "aaa").And(query.Less("f2", "aaa")),
			0, "f1<$1 AND f2<$2", []interface{}{"aaa", "aaa"},
		},
		{"eq_and_con",
			query.Less("f1", 1).And(
				query.Less("f2", 5).And(
					query.Less("f3", 7),
				),
			),
			0, "f1<$1 AND f2<$2 AND f3<$3", []interface{}{1, 5, 7},
		},
		{"eq_and_dis",
			query.Less("f1", 1).And(
				query.Less("f2", "2").Or(
					query.Less("f3", "3"),
				),
			),
			0, "f1<$1 AND (f2<$2 OR f3<$3)", []interface{}{1, "2", "3"},
		},
		{"ne_and_eq",
			query.Not(query.Less("f1", "aaa")).And(query.Less("f2", "aaa")),
			0, "NOT f1<$1 AND f2<$2", []interface{}{"aaa", "aaa"},
		},
		{"ne_and_con",
			query.Not(query.Less("f1", 2)).And(
				query.Less("f2", 3).And(
					query.Less("f3", 5),
				),
			),
			0, "NOT f1<$1 AND f2<$2 AND f3<$3", []interface{}{2, 3, 5},
		},
		{"ne_and_dis",
			query.Not(query.Less("f1", 1)).And(
				query.Less("f2", 3).Or(
					query.Less("f3", 11),
				),
			),
			7, "NOT f1<$8 AND (f2<$9 OR f3<$10)", []interface{}{1, 3, 11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}

func Test_Less_RenderSQL(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"lt_string",
			query.Less("f1", "aaa"),
			0, "f1<?", []interface{}{"aaa"},
		},
		{"lt_int",
			query.Less("f1", 1),
			1, "f1<?", []interface{}{1},
		},
		{"not_lt",
			query.Not(query.Less("f1", 2)),
			2, "NOT f1<?", []interface{}{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.RenderSQL())
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}
