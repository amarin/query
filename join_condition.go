package query

// JoinCondition defines data structure to store fields parameters used in table join.
// Could be created with JoinFields.
// Not intended to use directly but only as a step of TableJoiner definition.
type JoinCondition struct {
	leftField  FieldDefinition
	rightField FieldDefinition
}

// Using transforms JoinCondition into TableJoiner.
// Takes TableJoinType to use.
// Returns TableJoiner having tables set from JoinCondition fields, specified TableJoinType and JoinCondition itself.
func (joinCondition JoinCondition) Using(joinType TableJoinType) TableJoiner {
	return TableJoiner{
		rightTable:    Table(joinCondition.rightField.TableName()),
		joinType:      joinType,
		joinCondition: joinCondition,
	}
}

// InnerJoin generates InnerJoin TableJoiner.
// Shorthand to Using(InnerJoin).
func (joinCondition JoinCondition) InnerJoin() TableJoiner {
	return joinCondition.Using(InnerJoin)
}

// LeftJoin generates LeftJoin TableJoiner.
// Shorthand to Using(LeftJoin).
func (joinCondition JoinCondition) LeftJoin() TableJoiner {
	return joinCondition.Using(LeftJoin)
}

// RightJoin generates RightJoin TableJoiner.
// Shorthand to Using(RightJoin).
func (joinCondition JoinCondition) RightJoin() TableJoiner {
	return joinCondition.Using(RightJoin)
}

// FullJoin generates FullJoin TableJoiner.
// Shorthand to Using(FullJoin).
func (joinCondition JoinCondition) FullJoin() TableJoiner {
	return joinCondition.Using(FullJoin)
}

// Render renders SQL JOIN condition clause.
func (joinCondition JoinCondition) Render() string {
	return joinCondition.leftField.RenderTableSpec() + "=" + joinCondition.rightField.RenderTableSpec()
}

// JoinFields makes JoinCondition parameters to render SQL JOIN condition clause.
// NOTE both FieldDefinition argument MUST already have table name attached.
func JoinFields(fromField FieldDefinition, toField FieldDefinition) JoinCondition {
	return JoinCondition{leftField: fromField, rightField: toField}
}
