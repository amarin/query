package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableJoiner_RenderFrom(t *testing.T) {
	tests := []struct {
		name   string
		joiner TableJoiner
		want   string
	}{
		{"full_definitions_inner",
			NewTableJoiner(
				Table("table2").As("t2"),
				InnerJoin,
				JoinFields(Field("table1.field1"), Field("table2.field2")),
			),
			"INNER JOIN table2 AS t2 ON table1.field1=table2.field2",
		},
		{"full_definitions_left",
			NewTableJoiner(
				Table("table2").As("t2"),
				LeftJoin,
				JoinFields(Field("table1.field1"), Field("table2.field2")),
			),
			"LEFT JOIN table2 AS t2 ON table1.field1=table2.field2",
		},
		{"full_definitions_right",
			NewTableJoiner(
				Table("table2").As("t2"),
				RightJoin,
				JoinFields(Field("table1.field1"), Field("table2.field2")),
			),
			"RIGHT JOIN table2 AS t2 ON table1.field1=table2.field2",
		},
		{"full_definitions_full",
			NewTableJoiner(
				Table("table2").As("t2"),
				FullJoin,
				JoinFields(Field("table1.field1"), Field("table2.field2")),
			),
			"FULL JOIN table2 AS t2 ON table1.field1=table2.field2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.joiner.RenderFrom())
		})
	}
}
