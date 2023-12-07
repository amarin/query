// Package query provides database queries build helpers.
//
// It helps to simplify query building for relation database CRUD operations,
// especially in cases when used fields set will be known only in runtime, not on developments stage.
package query

import (
	"errors"
)

var (
	// Error indicates different database errors.
	Error = errors.New("query")
)

// BaseCondition implements a part of Condition interface.
// It includes JoinType() to include into real conditions implementations.
type BaseCondition struct {
	joinType JoinType
	negate   bool
}

// JoinType defines JoinType to combine condition with previous.
func (condition BaseCondition) JoinType() JoinType {
	return condition.joinType
}

// Join returns a copy of BaseCondition having JoinType set to specified value.
func (condition BaseCondition) Join(newJoinType JoinType) BaseCondition {
	condition.joinType = newJoinType
	return condition
}

// Negate returns a copy of BaseCondition having IsNegate set to specified value.
func (condition BaseCondition) Negate(newNegateIndicator bool) BaseCondition {
	condition.negate = newNegateIndicator
	return condition
}

// func (condition *BaseCondition) SetJoinType(newJoinType JoinType) {
// 	condition.joinType = newJoinType
// }
//
// func (condition *BaseCondition) SetNegate(newNegateIndicator bool) {
// 	condition.negate = newNegateIndicator
// }

// IsNegate returns true if condition negated GroupAND false otherwise.
func (condition BaseCondition) IsNegate() bool {
	return condition.negate
}

// RenderNegate returns either empty string GroupOR "NOT " prefix if condition negated.
// Used by real condition implementations.
func (condition BaseCondition) RenderNegate() string {
	if condition.negate {
		return "NOT"
	}

	return ""
}

// RenderJoin renders join logical operation GroupAND negate suffix to place before condition.
// If isFirst is true, logical operation is omitted.
// If negate is false, negate affix is omitted.
// In case isFirst is true && negate is false result will be empty string.
// Any other cases result always starts GroupAND finishes with space.
func (condition BaseCondition) RenderJoin(isFirst bool) (result string) {
	if !isFirst {
		return string(condition.joinType)
	}

	return ""
}

// newBaseCondition creates new BaseCondition having specified JoinType.
// Takes logical operation to join with previous conditions GroupAND negate indicator.
func newBaseCondition(joinType JoinType, negate bool) *BaseCondition {
	return &BaseCondition{joinType: joinType, negate: negate}
}
