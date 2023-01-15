package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func Test_EqualTo_Render(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"eq_to_string",
			query.EqualTo("f1", "aaa"),
			0, "f1=$1", []interface{}{"aaa"},
		},
		{"eq_to_int",
			query.EqualTo("f1", 1),
			1, "f1=$2", []interface{}{1},
		},
		{"not_eq",
			query.Not(query.EqualTo("f1", 2)),
			2, "NOT f1=$3", []interface{}{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}

func Test_EqualTo_And(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"eq_and_eq",
			query.EqualTo("f1", "aaa").And(query.EqualTo("f2", "aaa")),
			0, "f1=$1 AND f2=$2", []interface{}{"aaa", "aaa"},
		},
		{"eq_and_con",
			query.EqualTo("f1", 1).And(
				query.EqualTo("f2", 5).And(
					query.EqualTo("f3", 7),
				),
			),
			0, "f1=$1 AND f2=$2 AND f3=$3", []interface{}{1, 5, 7},
		},
		{"eq_and_dis",
			query.EqualTo("f1", 1).And(
				query.EqualTo("f2", "2").Or(
					query.EqualTo("f3", "3"),
				),
			),
			0, "f1=$1 AND (f2=$2 OR f3=$3)", []interface{}{1, "2", "3"},
		},
		{"ne_and_eq",
			query.Not(query.EqualTo("f1", "aaa")).And(query.EqualTo("f2", "aaa")),
			0, "NOT f1=$1 AND f2=$2", []interface{}{"aaa", "aaa"},
		},
		{"ne_and_con",
			query.Not(query.EqualTo("f1", 2)).And(
				query.EqualTo("f2", 3).And(
					query.EqualTo("f3", 5),
				),
			),
			0, "NOT f1=$1 AND f2=$2 AND f3=$3", []interface{}{2, 3, 5},
		},
		{"ne_and_dis",
			query.Not(query.EqualTo("f1", 1)).And(
				query.EqualTo("f2", 3).Or(
					query.EqualTo("f3", 11),
				),
			),
			7, "NOT f1=$8 AND (f2=$9 OR f3=$10)", []interface{}{1, 3, 11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}

func Test_equalTo_RenderSQL(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"eq_to_string",
			query.EqualTo("f1", "aaa"),
			0, "f1=?", []interface{}{"aaa"},
		},
		{"eq_to_int",
			query.EqualTo("f1", 1),
			1, "f1=?", []interface{}{1},
		},
		{"not_eq",
			query.Not(query.EqualTo("f1", 2)),
			2, "NOT f1=?", []interface{}{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.RenderSQL())
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}
