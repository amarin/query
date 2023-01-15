package query

// Not sets negate flag to single condition or group.
func Not(condition Condition) Condition {
	return condition.Negate(true)
}

// And makes condition or condition group with logical join AND.
// If single condition is specified, its join type is changed to LogicalAND.
func And(condition Condition, additionalConditions ...Condition) Condition {
	if len(additionalConditions) == 0 {
		return LogicalAND.Join(condition)
	}

	var (
		group        Group
		alreadyGroup bool
	)

	if group, alreadyGroup = condition.(Group); !alreadyGroup {
		group = NewGroup(LogicalAND, condition)
	}

	return group.And(additionalConditions...)
}

// Or makes condition or condition group with logical join OR.
// If single condition is specified, its join type is changed to LogicalOR.
func Or(condition Condition, additionalConditions ...Condition) Condition {
	if len(additionalConditions) == 0 {
		return LogicalOR.Join(condition)
	}

	var (
		group        Group
		alreadyGroup bool
	)

	if group, alreadyGroup = condition.(Group); !alreadyGroup {
		group = NewGroup(LogicalOR, condition)
	}

	return LogicalOR.Join(group.Or(additionalConditions...))
}
