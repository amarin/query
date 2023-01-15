package query

import (
	"strings"
)

// ClauseFromRenderer defines interface to use instead of direct table name in SQL FROM clause.
type ClauseFromRenderer interface {
	// RenderFrom returns string representation of table name or tables join with possible tables aliases.
	RenderFrom() string
}

// TableJoiner defines a table join definition data structure.
// It stores internally join type, right table and join conditions.
type TableJoiner struct {
	rightTable    TableIdent
	joinType      TableJoinType
	joinCondition JoinCondition
}

// JoinCondition returns tables JoinCondition.
func (tableJoiner TableJoiner) JoinCondition() JoinCondition {
	return tableJoiner.joinCondition
}

// JoinType returns TableJoinType.
func (tableJoiner TableJoiner) JoinType() TableJoinType {
	return tableJoiner.joinType
}

// By returns a TableJoiner having (re-)defined which fields are using to join tables.
func (tableJoiner TableJoiner) By(leftTableField FieldDefinition, rightTableField FieldDefinition) TableJoiner {
	tableJoiner.joinCondition = JoinFields(leftTableField, rightTableField)
	return tableJoiner
}

// RenderFrom returns SQL FROM clause filled with required tables.
func (tableJoiner TableJoiner) RenderFrom() string {
	return strings.Join([]string{
		tableJoiner.joinType.String(),
		tableJoiner.rightTable.RenderFrom(),
		kwOn.String(),
		tableJoiner.joinCondition.Render(),
	}, " ")
}

// NewTableJoiner creates new table joiner.
func NewTableJoiner(right TableIdent, joinType TableJoinType, on JoinCondition) TableJoiner {
	return TableJoiner{rightTable: right, joinType: joinType, joinCondition: on}
}
