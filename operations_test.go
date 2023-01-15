package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestNot(t *testing.T) {
	tests := []struct {
		name         string
		cond         query.Condition
		expectRender string
	}{
		{"negate_eq",
			query.EqualTo("f1", 0),
			"NOT f1=$1"},
		{"negate_contains",
			query.Contains("f1", "0"),
			"f1 NOT LIKE $1"},
		{"negate_icontains",
			query.IContains("f1", "0"),
			"f1 NOT ILIKE $1"},
		{"negate_is_null",
			query.IsNull("f1"),
			"f1 IS NOT NULL"},
		{"negate_dis",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			"NOT (f1=$1 OR f2=$2)"},
		{"negate_con",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			"NOT (f1=$1 AND f2=$2)"},
		{"negate_dis_not",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.Not(query.EqualTo("f2", 0))),
			"NOT (f1=$1 OR NOT f2=$2)"},
		{"negate_con_not",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 0), query.Not(query.EqualTo("f2", 0))),
			"NOT (f1=$1 AND NOT f2=$2)"},
		{"negate_not_dis",
			query.NewGroup(query.LogicalOR, query.Not(query.EqualTo("f1", 0)), query.EqualTo("f2", 0)),
			"NOT (NOT f1=$1 OR f2=$2)"},
		{"negate_not_con",
			query.NewGroup(query.LogicalAND, query.Not(query.EqualTo("f1", 0)), query.EqualTo("f2", 0)),
			"NOT (NOT f1=$1 AND f2=$2)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			negated := query.Not(tt.cond)
			require.True(t, negated.IsNegate())
			require.Equal(t, tt.expectRender, negated.Render(0))
		})
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		name                 string
		initial              query.Condition
		additionalConditions []query.Condition
		expectRender         string
	}{
		{"eq_and_eq",
			query.EqualTo("f1", 0),
			[]query.Condition{query.EqualTo("f2", 0)},
			"f1=$1 AND f2=$2"},
		{"eq_and_like",
			query.EqualTo("f1", 0),
			[]query.Condition{query.Contains("f2", "0")},
			"f1=$1 AND f2 LIKE $2"},
		{"eq_and_null",
			query.EqualTo("f1", 0),
			[]query.Condition{query.IsNull("f2")},
			"f1=$1 AND f2 IS NULL"},
		{"eq_and_not_eq",
			query.EqualTo("f1", 0),
			[]query.Condition{query.Not(query.EqualTo("f2", 0))},
			"f1=$1 AND NOT f2=$2"},
		{"eq_and_not_like",
			query.EqualTo("f1", 0),
			[]query.Condition{query.Not(query.Contains("f2", "0"))},
			"f1=$1 AND f2 NOT LIKE $2"},
		{"eq_and_not_null",
			query.EqualTo("f1", 0),
			[]query.Condition{query.Not(query.IsNull("f2"))},
			"f1=$1 AND f2 IS NOT NULL"},
		{"not_eq_and_eq",
			query.Not(query.EqualTo("f1", 0)),
			[]query.Condition{query.EqualTo("f2", 0)},
			"NOT f1=$1 AND f2=$2"},
		{"not_eq_and_like",
			query.Not(query.EqualTo("f1", 0)),
			[]query.Condition{query.Contains("f2", "0")},
			"NOT f1=$1 AND f2 LIKE $2"},
		{"not_eq_and_null",
			query.Not(query.EqualTo("f1", 0)),
			[]query.Condition{query.IsNull("f2")},
			"NOT f1=$1 AND f2 IS NULL"},
		{"eq_and_con_eq",
			query.EqualTo("f1", 0),
			[]query.Condition{query.EqualTo("f2", 0).And(query.EqualTo("f3", 0))},
			"f1=$1 AND f2=$2 AND f3=$3"},
		{"eq_and_con_like",
			query.EqualTo("f1", 0),
			[]query.Condition{query.Contains("f2", "0").And(query.Contains("f3", "0"))},
			"f1=$1 AND f2 LIKE $2 AND f3 LIKE $3"},
		{"eq_and_con_null",
			query.EqualTo("f1", 0),
			[]query.Condition{query.IsNull("f2").And(query.IsNull("f3"))},
			"f1=$1 AND f2 IS NULL AND f3 IS NULL"},
		{"con_and_eq",
			query.EqualTo("f1", 0).And(query.EqualTo("f2", 0)),
			[]query.Condition{query.EqualTo("f3", 0)},
			"f1=$1 AND f2=$2 AND f3=$3"},
		{"con_and_like",
			query.EqualTo("f1", 0).And(query.EqualTo("f2", 0)),
			[]query.Condition{query.Contains("f3", "0")},
			"f1=$1 AND f2=$2 AND f3 LIKE $3"},
		{"con_and_null",
			query.EqualTo("f1", 0).And(query.EqualTo("f2", 0)),
			[]query.Condition{query.IsNull("f3")},
			"f1=$1 AND f2=$2 AND f3 IS NULL"},
		{"con_and_dis",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.NewGroup(query.LogicalOR, query.EqualTo("f3", 0), query.EqualTo("f4", 0))},
			"f1=$1 AND f2=$2 AND (f3=$3 OR f4=$4)"},
		{"dis_and_eq",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.EqualTo("f3", 0)},
			"(f1=$1 OR f2=$2) AND f3=$3"},
		{"dis_and_like",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.Contains("f3", "0")},
			"(f1=$1 OR f2=$2) AND f3 LIKE $3"},
		{"dis_and_ilike",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.IContains("f3", "0")},
			"(f1=$1 OR f2=$2) AND f3 ILIKE $3"},
		{"dis_and_null",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.IsNull("f3")},
			"(f1=$1 OR f2=$2) AND f3 IS NULL"},
		{"dis_and_con",
			query.EqualTo("f1", 0).Or(query.EqualTo("f2", 0)),
			[]query.Condition{query.EqualTo("f3", 0).And(query.EqualTo("f4", 0))},
			"(f1=$1 OR f2=$2) AND f3=$3 AND f4=$4"},
		{"dis_and_dis",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.NewGroup(query.LogicalOR, query.EqualTo("f3", 0), query.EqualTo("f4", 0))},
			"(f1=$1 OR f2=$2) AND (f3=$3 OR f4=$4)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.initial.And(tt.additionalConditions...)
			require.Equal(t, tt.expectRender, res.Render(0))
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		name                 string
		initial              query.Condition
		additionalConditions []query.Condition
		expectRender         string
	}{
		{"eq_or_eq",
			query.EqualTo("f1", 0),
			[]query.Condition{query.EqualTo("f2", 0)},
			"f1=$1 OR f2=$2"},
		{"eq_or_like",
			query.EqualTo("f1", 0),
			[]query.Condition{query.Contains("f2", "0")},
			"f1=$1 OR f2 LIKE $2"},
		{"eq_or_null",
			query.EqualTo("f1", 0),
			[]query.Condition{query.IsNull("f2")},
			"f1=$1 OR f2 IS NULL"},
		{"eq_or_not_eq",
			query.EqualTo("f1", 0),
			[]query.Condition{query.Not(query.EqualTo("f2", 0))},
			"f1=$1 OR NOT f2=$2"},
		{"eq_or_not_like",
			query.EqualTo("f1", 0),
			[]query.Condition{query.Not(query.Contains("f2", "0"))},
			"f1=$1 OR f2 NOT LIKE $2"},
		{"eq_or_not_null",
			query.EqualTo("f1", 0),
			[]query.Condition{query.Not(query.IsNull("f2"))},
			"f1=$1 OR f2 IS NOT NULL"},
		{"not_eq_or_eq",
			query.Not(query.EqualTo("f1", 0)),
			[]query.Condition{query.EqualTo("f2", 0)},
			"NOT f1=$1 OR f2=$2"},
		{"not_eq_or_like",
			query.Not(query.EqualTo("f1", 0)),
			[]query.Condition{query.Contains("f2", "0")},
			"NOT f1=$1 OR f2 LIKE $2"},
		{"not_eq_or_null",
			query.Not(query.EqualTo("f1", 0)),
			[]query.Condition{query.IsNull("f2")},
			"NOT f1=$1 OR f2 IS NULL"},
		{"eq_or_con",
			query.EqualTo("f1", 0),
			[]query.Condition{query.EqualTo("f2", 0).And(query.EqualTo("f3", 0))},
			"f1=$1 OR (f2=$2 AND f3=$3)"},
		{"con_or_eq",
			query.EqualTo("f1", 0).And(query.EqualTo("f2", 0)),
			[]query.Condition{query.EqualTo("f3", 0)},
			"(f1=$1 AND f2=$2) OR f3=$3"},
		{"con_or_like",
			query.EqualTo("f1", 0).And(query.EqualTo("f2", 0)),
			[]query.Condition{query.Contains("f3", "0")},
			"(f1=$1 AND f2=$2) OR f3 LIKE $3"},
		{"con_or_null",
			query.EqualTo("f1", 0).And(query.EqualTo("f2", 0)),
			[]query.Condition{query.IsNull("f3")},
			"(f1=$1 AND f2=$2) OR f3 IS NULL"},
		{"con_or_dis",
			query.NewGroup(query.LogicalAND, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.NewGroup(query.LogicalOR, query.EqualTo("f3", 0), query.EqualTo("f4", 0))},
			"(f1=$1 AND f2=$2) OR f3=$3 OR f4=$4"},
		{"dis_or_eq",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.EqualTo("f3", 0)},
			"f1=$1 OR f2=$2 OR f3=$3"},
		{"dis_or_like",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.Contains("f3", "0")},
			"f1=$1 OR f2=$2 OR f3 LIKE $3"},
		{"dis_or_ilike",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.IContains("f3", "0")},
			"f1=$1 OR f2=$2 OR f3 ILIKE $3"},
		{"dis_or_null",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.IsNull("f3")},
			"f1=$1 OR f2=$2 OR f3 IS NULL"},
		{"dis_or_con",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.EqualTo("f3", 0).And(query.EqualTo("f4", 0))},
			"f1=$1 OR f2=$2 OR (f3=$3 AND f4=$4)"},
		{"dis_and_dis",
			query.NewGroup(query.LogicalOR, query.EqualTo("f1", 0), query.EqualTo("f2", 0)),
			[]query.Condition{query.NewGroup(query.LogicalOR, query.EqualTo("f3", 0), query.EqualTo("f4", 0))},
			"f1=$1 OR f2=$2 OR f3=$3 OR f4=$4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.initial.Or(tt.additionalConditions...)
			require.Equal(t, tt.expectRender, res.Render(0))
		})
	}
}
