package query

// JoinType defines conditions join variants.
type JoinType string

// String returns string type value of JoinType.
func (joinType JoinType) String() string {
	return string(joinType)
}

const (
	// LogicalAND defines logical AND operator to join conditions.
	LogicalAND JoinType = "AND"

	// LogicalOR defines logical OR operator to join conditions.
	LogicalOR JoinType = "OR"
)

// Join wraps specified condition to combine with previous conditions with logical operation specified by JoinType.
// Does the same as And() but in method of JoinType.
func (joinType JoinType) Join(conditions ...Condition) (result Condition) {
	if len(conditions) > 1 {
		return NewGroup(joinType, conditions...)
	}

	result = conditions[0]
	return result.Join(joinType)
}

// SetToAll returns supplied list with join type set to current in every list element.
func (joinType JoinType) SetToAll(conditions ...Condition) (result []Condition) {
	result = make([]Condition, len(conditions))
	for idx, condition := range conditions {
		result[idx] = joinType.Join(condition)
	}

	return result
}
