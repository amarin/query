package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func Test_group_FieldName(t *testing.T) {
	group := query.NewGroup(query.LogicalAND)
	require.Empty(t, group.FieldName())
}

func Test_group_add(t *testing.T) {
	tests := []struct {
		name         string
		addFunc      func(conditions ...query.Condition) query.Condition
		fields       []query.Condition
		renderIdx    int
		expectRender string
		expectValues []interface{}
	}{
		{"nil_con_one_gives_one",
			query.NewGroup(query.LogicalAND).And,
			[]query.Condition{query.EqualTo("t1", 1)},
			0,
			"t1=$1",
			[]interface{}{1},
		},
		{"nil_dis_one_gives_one",
			query.NewGroup(query.LogicalAND).And,
			[]query.Condition{query.EqualTo("t1", 1)},
			0,
			"t1=$1",
			[]interface{}{1},
		},
		{"one_con_one_gives_con",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 1)).And,
			[]query.Condition{query.EqualTo("f2", 2)},
			1,
			"f1=$2 AND f2=$3",
			[]interface{}{1, 2},
		},
		{"one_dis_one_gives_dis",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 1)).Or,
			[]query.Condition{query.EqualTo("f2", 2)},
			1,
			"f1=$2 OR f2=$3",
			[]interface{}{1, 2},
		},
		{"con_con_nil_gives_con",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).And,
			nil,
			3,
			"f1=$4 AND f2=$5",
			[]interface{}{1, 2},
		},
		{"con_con_one_gives_con",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).And,
			[]query.Condition{query.Contains("a1", "aaa")},
			3,
			"f1=$4 AND f2=$5 AND a1 LIKE $6",
			[]interface{}{1, 2, "%aaa%"},
		},
		{"dis_dis_one_gives_dis",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).Or,
			[]query.Condition{query.Contains("a1", "aaa")},
			4,
			"f1=$5 OR f2=$6 OR a1 LIKE $7",
			[]interface{}{1, 2, "%aaa%"},
		},
		{"dis_con_one_gives_dis_con_one",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).And,
			[]query.Condition{query.Contains("a1", "aaa")},
			3,
			"(f1=$4 OR f2=$5) AND a1 LIKE $6",
			[]interface{}{1, 2, "%aaa%"},
		},
		{"con_dis_one_gives_con_dis_one",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).Or,
			[]query.Condition{query.Contains("a1", "aaa")},
			3,
			"(f1=$4 AND f2=$5) OR a1 LIKE $6",
			[]interface{}{1, 2, "%aaa%"},
		},
		{"dis_con_dis_gives_dis_con_dis",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).And,
			[]query.Condition{query.NewGroup(
				query.LogicalOR,
				query.EqualTo("a1", 3),
				query.EqualTo("a2", 4),
			)},
			4,
			"(f1=$5 OR f2=$6) AND (a1=$7 OR a2=$8)",
			[]interface{}{1, 2, 3, 4},
		},
		{"con_dis_con_gives_con_dis_con",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).Or,
			[]query.Condition{query.NewGroup(
				query.LogicalAND,
				query.EqualTo("a1", 3),
				query.EqualTo("a2", 4),
			)},
			5,
			"(f1=$6 AND f2=$7) OR (a1=$8 AND a2=$9)",
			[]interface{}{1, 2, 3, 4},
		},
		{"dis_con_con_gives_dis_con_raw",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).And,
			[]query.Condition{query.NewGroup(
				query.LogicalAND,
				query.EqualTo("a1", 3),
				query.EqualTo("a2", 4),
			)},
			10,
			"(f1=$11 OR f2=$12) AND a1=$13 AND a2=$14",
			[]interface{}{1, 2, 3, 4},
		},
		{"con_con_con_gives_con_raw",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).And,
			[]query.Condition{query.NewGroup(
				query.LogicalAND,
				query.EqualTo("a1", 3),
				query.EqualTo("a2", 4),
			)},
			10,
			"f1=$11 AND f2=$12 AND a1=$13 AND a2=$14",
			[]interface{}{1, 2, 3, 4},
		},
		{"con_con_dis_gives_raw_con_dis",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 1), query.EqualTo("f2", 2)).And,
			[]query.Condition{query.NewGroup(
				query.LogicalOR,
				query.EqualTo("a1", 3),
				query.EqualTo("a2", 4),
			)},
			10,
			"f1=$11 AND f2=$12 AND (a1=$13 OR a2=$14)",
			[]interface{}{1, 2, 3, 4},
		},
		{"eq_or_con",
			query.EqualTo("f1", 0).Or,
			[]query.Condition{query.EqualTo("f2", 0).And(query.EqualTo("f3", 0))},
			0,
			"f1=$1 OR (f2=$2 AND f3=$3)",
			[]interface{}{0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := tt.addFunc(tt.fields...)
			require.Equal(t, tt.expectRender, group.Render(tt.renderIdx))
			require.Equal(t, tt.expectValues, group.Values())
		})
	}
}

func TestGroup_Render(t *testing.T) {
	tests := []struct {
		name    string
		group   query.Condition
		idx     int
		wantSql string
	}{
		{"neq_dis_eq_no_brackets",
			query.Not(query.EqualTo("f1", 11)).Or(query.EqualTo("f2", 13)),
			0,
			"NOT f1=$1 OR f2=$2",
		},
		{"neq_dis_eq_no_brackets_table",
			query.Not(query.EqualTo("f1", 11)).Or(query.EqualTo("f2", 13)).ApplyFieldTable("test"),
			0,
			"NOT test.f1=$1 OR test.f2=$2",
		},
		{"neq_con_eq_no_brackets",
			query.Not(query.EqualTo("f1", 11)).And(query.EqualTo("f2", 13)),
			0,
			"NOT f1=$1 AND f2=$2",
		},
		{"eq_dis_not_eq_no_brackets",
			query.EqualTo("f1", 11).Or(query.Not(query.EqualTo("f2", 13))),
			0,
			"f1=$1 OR NOT f2=$2",
		},
		{"eq_con_not_eq_no_brackets",
			query.EqualTo("f1", 11).And(query.Not(query.EqualTo("f2", 13))),
			0,
			"f1=$1 AND NOT f2=$2",
		},
		{"eq_con_eq_in_brackets",
			(query.EqualTo("f1", 11).And(query.EqualTo("f2", 13)).(query.Group)).WithBrackets(),
			0,
			"(f1=$1 AND f2=$2)",
		},
		{"eq_dis_eq_in_brackets",
			(query.EqualTo("f1", 11).Or(query.EqualTo("f2", 13)).(query.Group)).WithBrackets(),
			0,
			"(f1=$1 OR f2=$2)",
		},
		{"not_eq_dis_eq_in_brackets",
			query.Not((query.EqualTo("f1", 11).Or(query.EqualTo("f2", 13)).(query.Group)).WithBrackets()),
			0,
			"NOT (f1=$1 OR f2=$2)",
		},
		{"not_eq_dis_eq_in_brackets_table",
			query.Not((query.EqualTo("f1", 11).Or(query.EqualTo("f2", 13)).(query.Group)).WithBrackets()).ApplyFieldTable("test"),
			0,
			"NOT (test.f1=$1 OR test.f2=$2)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantSql, tt.group.Render(tt.idx))
		})
	}
}

func TestGroup_RenderSQL(t *testing.T) {
	tests := []struct {
		name    string
		group   query.Condition
		wantSql string
	}{
		{"neq_dis_eq_no_brackets",
			query.Not(query.EqualTo("f1", 11)).Or(query.EqualTo("f2", 13)),
			"NOT f1=? OR f2=?",
		},
		{"neq_dis_eq_no_brackets_table",
			query.Not(query.EqualTo("f1", 11)).Or(query.EqualTo("f2", 13)).ApplyFieldTable("test"),
			"NOT test.f1=? OR test.f2=?",
		},
		{"neq_con_eq_no_brackets",
			query.Not(query.EqualTo("f1", 11)).And(query.EqualTo("f2", 13)),
			"NOT f1=? AND f2=?",
		},
		{"eq_dis_not_eq_no_brackets",
			query.EqualTo("f1", 11).Or(query.Not(query.EqualTo("f2", 13))),
			"f1=? OR NOT f2=?",
		},
		{"eq_con_not_eq_no_brackets",
			query.EqualTo("f1", 11).And(query.Not(query.EqualTo("f2", 13))),
			"f1=? AND NOT f2=?",
		},
		{"eq_con_eq_in_brackets",
			(query.EqualTo("f1", 11).And(query.EqualTo("f2", 13)).(query.Group)).WithBrackets(),
			"(f1=? AND f2=?)",
		},
		{"eq_dis_eq_in_brackets",
			(query.EqualTo("f1", 11).Or(query.EqualTo("f2", 13)).(query.Group)).WithBrackets(),
			"(f1=? OR f2=?)",
		},
		{"not_eq_dis_eq_in_brackets",
			query.Not((query.EqualTo("f1", 11).Or(query.EqualTo("f2", 13)).(query.Group)).WithBrackets()),
			"NOT (f1=? OR f2=?)",
		},
		{"not_eq_dis_eq_in_brackets_table",
			query.Not((query.EqualTo("f1", 11).Or(query.EqualTo("f2", 13)).(query.Group)).WithBrackets()).ApplyFieldTable("test"),
			"NOT (test.f1=? OR test.f2=?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantSql, tt.group.RenderSQL())
		})
	}
}
