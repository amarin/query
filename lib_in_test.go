package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func Test_In_Render(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"in_string",
			query.In("f1", "aaa"),
			0, "f1 IN ($1)", []interface{}{"aaa"},
		},
		{"in_many_strings",
			query.In("f1", []string{"aaa", "132"}),
			0, "f1 IN ($1,$2)", []interface{}{"aaa", "132"},
		},
		{"in_many_strings_with_offset",
			query.In("f1", []string{"aaa", "132"}),
			1, "f1 IN ($2,$3)", []interface{}{"aaa", "132"},
		},
		{"in_int",
			query.In("f1", 1),
			1, "f1 IN ($2)", []interface{}{1},
		},
		{"not_lt",
			query.Not(query.In("f1", 2)),
			2, "f1 NOT IN ($3)", []interface{}{2},
		},
		{"not_lt_many_strings",
			query.Not(query.In("f1", []string{"123", "321"})),
			2, "f1 NOT IN ($3,$4)", []interface{}{"123", "321"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}

func Test_In_And(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"in_and_in",
			query.In("f1", "aaa").And(query.In("f2", "aaa")),
			0, "f1 IN ($1) AND f2 IN ($2)", []interface{}{"aaa", "aaa"},
		},
		{"in_and_con",
			query.In("f1", 1).And(
				query.In("f2", 5).And(
					query.In("f3", 7),
				),
			),
			0, "f1 IN ($1) AND f2 IN ($2) AND f3 IN ($3)", []interface{}{1, 5, 7},
		},
		{"eq_and_dis",
			query.In("f1", 1).And(
				query.In("f2", "2").Or(
					query.In("f3", "3"),
				),
			),
			0, "f1 IN ($1) AND (f2 IN ($2) OR f3 IN ($3))", []interface{}{1, "2", "3"},
		},
		{"ne_and_eq",
			query.Not(query.In("f1", "aaa")).And(query.In("f2", "aaa")),
			0, "f1 NOT IN ($1) AND f2 IN ($2)", []interface{}{"aaa", "aaa"},
		},
		{"ne_and_con",
			query.Not(query.In("f1", 2)).And(
				query.In("f2", 3).And(
					query.In("f3", 5),
				),
			),
			0, "f1 NOT IN ($1) AND f2 IN ($2) AND f3 IN ($3)", []interface{}{2, 3, 5},
		},
		{"ne_and_dis",
			query.Not(query.In("f1", 1)).And(
				query.In("f2", 3).Or(
					query.In("f3", 11),
				),
			),
			7, "f1 NOT IN ($8) AND (f2 IN ($9) OR f3 IN ($10))", []interface{}{1, 3, 11},
		},
		{"ne_and_dis_many",
			query.Not(query.In("f1", []string{"1", "2"})).And(
				query.In("f2", 3).Or(
					query.In("f3", []string{"4", "5"}),
				),
			),
			7,
			"f1 NOT IN ($8,$9) AND (f2 IN ($10) OR f3 IN ($11,$12))",
			[]interface{}{"1", "2", 3, "4", "5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.Render(tt.paramCount))
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}

func Test_In_RenderSQL(t *testing.T) {
	tests := []struct {
		name       string
		cond       query.Condition
		paramCount int
		want       string
		values     []interface{}
	}{
		{"in_string",
			query.In("f1", "aaa"),
			0, "f1 IN (?)", []interface{}{"aaa"},
		},
		{"in_int",
			query.In("f1", 1),
			1, "f1 IN (?)", []interface{}{1},
		},
		{"not_in",
			query.Not(query.In("f1", 2)),
			2, "f1 NOT IN (?)", []interface{}{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cond.RenderSQL())
			require.Equal(t, tt.values, tt.cond.Values())
		})
	}
}
