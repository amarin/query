package query

// incompleteBaseJoiner is not exported implementation of IncompleteSelectJoin just for prevent manual instance building.
// The goal is to provide chain style method to generate every clause types.
type incompleteBaseJoiner struct {
	BaseSelectBuilder
}

// On finalises SQL JOIN definition adding left and right tables fields onto JOIN ... ON clause.
// Note both left and right FieldDefinition's required to contain table name.
// Returns updated BaseSelectBuilder instance.
func (joiner incompleteBaseJoiner) On(left FieldDefinition, right FieldDefinition) (updated BaseSelectBuilder) {
	updated = joiner.BaseSelectBuilder
	unfinishedJoinIdx := len(updated.joins) - 1
	updated.joins[unfinishedJoinIdx] = updated.joins[unfinishedJoinIdx].By(left, right)

	return updated
}

// incompleteManyJoiner implements IncompleteSelectJoin preventing direct instance usage.
// // The goal is to provide chain style method to generate every clause types.
type incompleteManyJoiner struct {
	SelectManyBuilder
}

// On finalises SQL JOIN definition adding left and right tables fields onto JOIN ... ON clause.
// Note both left and right FieldDefinition's required to contain table name.
// Returns updated BaseSelectBuilder instance.
func (joiner incompleteManyJoiner) On(left FieldDefinition, right FieldDefinition) (updated SelectManyBuilder) {
	updated = joiner.SelectManyBuilder
	unfinishedJoinIdx := len(updated.joins) - 1
	updated.joins[unfinishedJoinIdx] = updated.joins[unfinishedJoinIdx].By(left, right)

	return updated
}

// incompleteManyJoiner implements IncompleteSelectJoin preventing direct instance usage.
// // The goal is to provide chain style method to generate every clause types.
type incompleteSingleJoiner struct {
	SelectSingleBuilder
}

// On finalises SQL JOIN definition adding left and right tables fields onto JOIN ... ON clause.
// Note both left and right FieldDefinition's required to contain table name.
// Returns updated BaseSelectBuilder instance.
func (joiner incompleteSingleJoiner) On(left FieldDefinition, right FieldDefinition) (updated SelectSingleBuilder) {
	updated = joiner.SelectSingleBuilder
	unfinishedJoinIdx := len(updated.joins) - 1
	updated.joins[unfinishedJoinIdx] = updated.joins[unfinishedJoinIdx].By(left, right)

	return updated
}
