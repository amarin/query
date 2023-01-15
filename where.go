package query

// WhereClause wraps conditions Group to handle conditions as whole WHERE clause block.
type WhereClause struct {
	group Group
}

// GroupAND adds a new condition group into current.
// Specified conditions joined into group with logical AND.
// Existed conditions are joined into group too.
// Both groups joined with logical AND operator.
// Returns updated WhereClause containing required condition set, leaving original WhereClause instance intact.
func (query WhereClause) GroupAND(conditions ...Condition) (updated WhereClause) {
	updated = query
	updated.group = updated.group.GroupAND(conditions...)
	return updated
}

// GroupOR creates a new conditions group in where clause.
// Specified conditions joined into group with logical AND.
// Existed conditions are joined into group too.
// Both groups joined with logical OR operator.
// Returns new WhereClause containing required condition set, leaving original WhereClause instance intact.
func (query WhereClause) GroupOR(conditions ...Condition) (updated WhereClause) {
	updated = query
	updated.group = updated.group.GroupOR(conditions...)
	return updated
}

// ApplyFieldSpec makes a copy of WhereClause with updated FieldDefinition on each child when field name matched.
// Implements Condition.
func (query WhereClause) ApplyFieldSpec(spec FieldDefinition) (updated WhereClause) {
	updated = query
	updated.group = updated.group.ApplyFieldSpec(spec).(Group)
	return updated
}

// Conditions returns a copy of attached conditions list.
func (query WhereClause) Conditions() (conditions []Condition) {
	conditions = make([]Condition, len(query.group.conditions))
	copy(conditions, query.group.conditions)

	return conditions
}

// Render renders SQL SELECT clause part for current group.
// Takes existed parameters number (0 means no parameters are defined yet)
// Note Render renders group without brackets itself if called directly.
// Any child conditions groups are enclosed into brackets internally.
func (query WhereClause) Render(parametersCount int) (sql string) {
	return query.group.Render(parametersCount)
}

// RenderSQL renders SQL clause or its part.
// Implementation should render parameters substitutions using standard sql "?"(question) character.
// Values result should contain the same values count as substitutions used in clause.
func (query WhereClause) RenderSQL() (sql string) {
	return query.group.RenderSQL()
}

// Values returns a set of values of grouped conditions.
func (query WhereClause) Values() (values []interface{}) {
	return query.group.Values()
}

// String returns a string representation of WhereClause.
func (query WhereClause) String() string {
	return "Where(" + query.group.RenderSQL() + ")"
}

// NewWhere creates new empty WhereClause instance.
func NewWhere() WhereClause {
	return WhereClause{group: NewGroup(LogicalAND)}
}
