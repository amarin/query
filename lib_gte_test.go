package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func Test_GreaterOrEqual_Render(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"gte_string",
			query.GreaterOrEqual("f1", "aaa"),
			0, "f1>=$1", []interface{}{"aaa"},
		},
		{"gte_int",
			query.GreaterOrEqual("f1", 1),
			1, "f1>=$2", []interface{}{1},
		},
		{"not_gte",
			query.Not(query.GreaterOrEqual("f1", 2)),
			2, "NOT f1>=$3", []interface{}{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}

func Test_GreaterOrEqual_And(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"eq_and_eq",
			query.GreaterOrEqual("f1", "aaa").And(query.GreaterOrEqual("f2", "aaa")),
			0, "f1>=$1 AND f2>=$2", []interface{}{"aaa", "aaa"},
		},
		{"eq_and_con",
			query.GreaterOrEqual("f1", 1).And(
				query.GreaterOrEqual("f2", 5).And(
					query.GreaterOrEqual("f3", 7),
				),
			),
			0, "f1>=$1 AND f2>=$2 AND f3>=$3", []interface{}{1, 5, 7},
		},
		{"eq_and_dis",
			query.GreaterOrEqual("f1", 1).And(
				query.GreaterOrEqual("f2", "2").Or(
					query.GreaterOrEqual("f3", "3"),
				),
			),
			0, "f1>=$1 AND (f2>=$2 OR f3>=$3)", []interface{}{1, "2", "3"},
		},
		{"ne_and_eq",
			query.Not(query.GreaterOrEqual("f1", "aaa")).And(query.GreaterOrEqual("f2", "aaa")),
			0, "NOT f1>=$1 AND f2>=$2", []interface{}{"aaa", "aaa"},
		},
		{"ne_and_con",
			query.Not(query.GreaterOrEqual("f1", 2)).And(
				query.GreaterOrEqual("f2", 3).And(
					query.GreaterOrEqual("f3", 5),
				),
			),
			0, "NOT f1>=$1 AND f2>=$2 AND f3>=$3", []interface{}{2, 3, 5},
		},
		{"ne_and_dis",
			query.Not(query.GreaterOrEqual("f1", 1)).And(
				query.GreaterOrEqual("f2", 3).Or(
					query.GreaterOrEqual("f3", 11),
				),
			),
			7, "NOT f1>=$8 AND (f2>=$9 OR f3>=$10)", []interface{}{1, 3, 11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}

func Test_GreaterOrEqual_RenderSQL(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"gte_string",
			query.GreaterOrEqual("f1", "aaa"),
			0, "f1>=?", []interface{}{"aaa"},
		},
		{"gte_int",
			query.GreaterOrEqual("f1", 1),
			1, "f1>=?", []interface{}{1},
		},
		{"not_gte",
			query.Not(query.GreaterOrEqual("f1", 2)),
			2, "NOT f1>=?", []interface{}{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.RenderSQL())
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}
