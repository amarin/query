package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestTableJoinType_String(t *testing.T) {
	tests := []struct {
		name string
		join query.TableJoinType
		want string
	}{
		{"inner", query.InnerJoin, "INNER JOIN"},
		{"left", query.LeftJoin, "LEFT JOIN"},
		{"right", query.RightJoin, "RIGHT JOIN"},
		{"full", query.FullJoin, "FULL JOIN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.join.String())
		})
	}
}

func TestTableJoinType_By(t *testing.T) {
	tests := []struct {
		name       string
		join       query.TableJoinType
		leftField  query.FieldDefinition
		rightField query.FieldDefinition
		wantFrom   string
	}{
		{"inner_join",
			query.InnerJoin,
			query.Table("t1").Field("f1"),
			query.Table("t2").Field("f2"),
			"INNER JOIN t2 ON t1.f1=t2.f2",
		},
		{"left_join",
			query.LeftJoin,
			query.Table("t1").Field("f1"),
			query.Table("t2").Field("f2"),
			"LEFT JOIN t2 ON t1.f1=t2.f2",
		},
		{"right_join",
			query.RightJoin,
			query.Table("t1").Field("f1"),
			query.Table("t2").Field("f2"),
			"RIGHT JOIN t2 ON t1.f1=t2.f2",
		},
		{"full_join",
			query.FullJoin,
			query.Table("t1").Field("f1"),
			query.Table("t2").Field("f2"),
			"FULL JOIN t2 ON t1.f1=t2.f2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantFrom, tt.join.By(tt.leftField, tt.rightField).RenderFrom())
		})
	}
}
