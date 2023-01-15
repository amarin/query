package query

import (
	"strconv"
)

// TableJoinType defines table join type. See InnerJoin, LeftJoin, RightJoin and FullJoin constants.
type TableJoinType int

const (
	// InnerJoin defines constant to indicate (INNER) JOIN returning records that have matching values in both tables.
	InnerJoin TableJoinType = iota

	// LeftJoin defines constant to indicate LEFT (OUTER) JOIN returning all records from the left table,
	// and the matched records from the right table.
	LeftJoin

	// RightJoin defines constant to indicate RIGHT (OUTER) JOIN returning from the right table,
	// and the matched records from the left table.
	RightJoin

	// FullJoin defines constant to indicate FULL (OUTER) JOIN returning all records
	// when there is a match in either left or right table.
	FullJoin
)

// String returns string representation of TableJoinType value.
// Implements fmt.Stringer. Used in queries builder.
func (join TableJoinType) String() string {
	switch join {
	case InnerJoin:
		return "INNER JOIN"
	case LeftJoin:
		return "LEFT JOIN"
	case RightJoin:
		return "RIGHT JOIN"
	case FullJoin:
		return "FULL JOIN"
	default:
		return "unknown JOIN(" + strconv.Itoa(int(join)) + ")"
	}
}

// By generates TableJoiner built as TableJoinType by specified FieldDefinition`s.
func (join TableJoinType) By(leftField FieldDefinition, rightField FieldDefinition) TableJoiner {
	return JoinFields(leftField, rightField).Using(join)
}

// Tables generates IncompleteSelectJoin instance built using leftTable BaseBuilder, right table and empty join condition.
// It should be finalized with IncompleteSelectJoin.On to get BaseSelectBuilder.
func (join TableJoinType) Tables(leftTable TableIdent, rightTable TableIdent) IncompleteSelectJoin {
	return leftTable.Select().joinIdent(rightTable, join)
}
