package query

// Group implements conditions grouping.
type Group struct {
	BaseCondition
	conditions   []Condition //
	withBrackets bool        // indicate whether to draw group in brackets. False by default.
}

// WithBrackets creates new group having same conditions but required to enclose conditions into brackets.
func (conditionsGroup Group) WithBrackets() Group {
	conditionsGroup.withBrackets = true
	return conditionsGroup
}

// Join returns a copy of Group having JoinType set to specified value.
func (conditionsGroup Group) Join(newJoinType JoinType) Condition {
	conditionsGroup.BaseCondition = conditionsGroup.BaseCondition.Join(newJoinType)
	return conditionsGroup
}

// Negate returns a copy of BaseCondition having IsNegate set to specified value.
func (conditionsGroup Group) Negate(newNegateIndicator bool) Condition {
	conditionsGroup.BaseCondition = conditionsGroup.BaseCondition.Negate(newNegateIndicator)
	return conditionsGroup
}

// FieldName returns empty string. Implements Condition.
func (conditionsGroup Group) FieldName() FieldName {
	return ""
}

// NewGroup makes a group of conditions with default external join AND.
// All conditions inside group joined with requested joinType.
func NewGroup(joinType JoinType, conditions ...Condition) Group {
	if conditions == nil {
		conditions = make([]Condition, 0)
	}
	return Group{
		BaseCondition: *newBaseCondition(LogicalAND, false),
		conditions:    joinType.SetToAll(conditions...),
		withBrackets:  false,
	}
}

// addGroup does extra routines while adding single group element.
// Note it operates over Group copy and returns pointer to modified copy.
func (conditionsGroup Group) addGroup(anotherGroup Group) Group {
	haveNone := len(conditionsGroup.conditions) == 0
	haveSingle := len(conditionsGroup.conditions) == 1

	havingOnlyAND := conditionsGroup.hasAll(LogicalAND)
	havingOnlyOR := conditionsGroup.hasAll(LogicalOR)
	havingAnyAND := conditionsGroup.hasOne(LogicalAND)
	havingAnyOR := conditionsGroup.hasOne(LogicalOR)

	addingOnlyAnd := anotherGroup.hasAll(LogicalAND)
	addingAnyAnd := anotherGroup.hasOne(LogicalAND)
	addingOnlyOR := anotherGroup.hasAll(LogicalOR)
	addingAnyOR := anotherGroup.hasOne(LogicalOR)

	op := anotherGroup.joinType

	switch {
	case haveNone: // current group empty, reset it op and elements from another
		conditionsGroup.conditions = anotherGroup.conditions
		conditionsGroup.joinType = op

	// check conditions on have single condition and adding group
	case haveSingle && op == LogicalAND && addingOnlyAnd: // AND AND (AND),- unpack another group into flat (AND...AND)
		conditionsGroup.conditions = append(conditionsGroup.conditions, anotherGroup.conditions...)
	case haveSingle && op == LogicalOR && addingOnlyOR: // OR OR (OR),- unpack another group into flat (OR-OR).
		conditionsGroup.conditions = append(conditionsGroup.conditions, anotherGroup.conditions...)
	case haveSingle: // any other single (or group) + group: enclose new group into brackets
		conditionsGroup.conditions = append(conditionsGroup.conditions, op.Join(anotherGroup.WithBrackets()))

	// check conditions on having many conditions and adding group with conjunction
	case havingOnlyOR && op == LogicalOR && addingOnlyOR: // (OR) OR (OR) => flat OR
		conditionsGroup.conditions = append(conditionsGroup.conditions, anotherGroup.conditions...)
	case havingOnlyOR && op == LogicalOR && addingAnyAnd: // (OR) OR (AND) => OR OR (AND)
		conditionsGroup.conditions = append(conditionsGroup.conditions, anotherGroup.WithBrackets())
	case havingAnyOR && op == LogicalAND && addingAnyAnd: // (OR) AND (AND) => (OR) AND AND
		conditionsGroup.conditions = append([]Condition{conditionsGroup.WithBrackets()}, anotherGroup.conditions...)
		// check conditions on having many conditions and adding group with conjunction

	case havingOnlyAND && op == LogicalAND && addingOnlyAnd: // (AND) AND (AND) => flat AND
		conditionsGroup.conditions = append(conditionsGroup.conditions, anotherGroup.conditions...)
		// check conditions on having many conditions and adding group
	case havingOnlyAND && op == LogicalAND && addingAnyOR: // (AND) AND (OR),- AND AND (OR)
		conditionsGroup.conditions = append(conditionsGroup.conditions, anotherGroup.WithBrackets())
		// check conditions on having many conditions and adding group
	// check conditions on having many conditions and adding group
	case havingAnyAND && op == LogicalOR && addingOnlyOR: // (AND) OR (OR) => (AND) OR OR
		conditionsGroup.conditions = append([]Condition{conditionsGroup.WithBrackets()}, anotherGroup.conditions...)

	default: // have many, adding many, join is not matched, join both with brackets
		conditionsGroup.conditions = []Condition{conditionsGroup.WithBrackets(), op.Join(anotherGroup.WithBrackets())}
	}

	return conditionsGroup
}

// addSingleCondition does extra routines while adding single condition.
// Note it operates over Group copy and returns pointer to modified copy.
func (conditionsGroup Group) addSingleCondition(condition Condition) Group {
	haveNone := len(conditionsGroup.conditions) == 0
	haveSingle := len(conditionsGroup.conditions) == 1
	haveMany := len(conditionsGroup.conditions) > 1
	havingOnlyAND := conditionsGroup.hasAll(LogicalAND)
	havingOnlyOR := conditionsGroup.hasAll(LogicalOR)
	groupToAdd, addingGroup := condition.(Group)
	op := condition.JoinType()

	switch {
	case addingGroup: // process adding groups in separate method
		return conditionsGroup.addGroup(groupToAdd)
	case haveNone: // current group empty, join single condition and set group join type
		conditionsGroup.conditions = append(conditionsGroup.conditions, condition)
		conditionsGroup.joinType = op
	case haveSingle: // just join new single to group joining with required op
		conditionsGroup.conditions = append(conditionsGroup.conditions, op.Join(condition))
	case haveMany && op == LogicalAND && havingOnlyAND: // have many, all joins are AND, join new to group
		conditionsGroup.conditions = append(conditionsGroup.conditions, op.Join(condition))
	case haveMany && op == LogicalOR && havingOnlyOR: // have many, all joins are OR, join new to group
		conditionsGroup.conditions = append(conditionsGroup.conditions, op.Join(condition))
	default: // have many, adding single, join is not matched, enclose existing into bracket
		conditionsGroup.conditions = []Condition{conditionsGroup.WithBrackets(), op.Join(condition)}
	}

	return conditionsGroup
}

// addConditions adds new conditions into existing set and returns modified group.
// Note it operates over Group copy and returns pointer to modified copy.
func (conditionsGroup Group) addConditions(op JoinType, conditions ...Condition) Group {
	havingAnyAND := conditionsGroup.hasOne(LogicalAND)
	havingAnyOR := conditionsGroup.hasOne(LogicalOR)
	haveConditions := len(conditionsGroup.conditions) > 0

	switch {
	case len(conditions) == 0: // simply return pointer to copy if no new conditions.
		return conditionsGroup
	case !haveConditions: // join to existed flat AND-AND new OR(AND_AND)
		conditionsGroup.conditions = append(conditionsGroup.conditions, conditions...)
		conditionsGroup.joinType = op
	case len(conditions) == 1: // process single element adding in separated Group.addSingleCondition
		return conditionsGroup.addSingleCondition(op.Join(conditions[0]))
	case havingAnyOR && op == LogicalAND: // enclose existed OR-OR into brackets and join flat AND-AND
		conditionsGroup.conditions = append([]Condition{conditionsGroup.WithBrackets()}, conditions...)
	case havingAnyAND && op == LogicalOR: // join to existed flat AND-AND new OR(AND_AND)
		conditionsGroup.conditions = append(conditionsGroup.conditions, op.Join(NewGroup(op, conditions...)))
	default: // have many, adding many, enclose both sets into brackets and return new group with required join.
		conditionsGroup.conditions = append([]Condition{
			conditionsGroup.WithBrackets()},
			NewGroup(LogicalAND, conditions...).WithBrackets(),
		)
	}

	return conditionsGroup
}

// hasAll returns true if all conditions joined with specified JoinType.
// Ignores first condition join type as everything joined to it.
func (conditionsGroup Group) hasAll(joinType JoinType) bool {
	for idx, item := range conditionsGroup.conditions {
		if idx == 0 {
			continue // ignore
		}
		if item.JoinType() != joinType {
			return false
		}
	}

	return true
}

// hasOne returns true if any condition joined with specified JoinType.
// Ignores first condition join type as everything joined to it.
func (conditionsGroup Group) hasOne(joinType JoinType) bool {
	for idx, item := range conditionsGroup.conditions {
		if idx == 0 {
			continue // ignore
		}
		if item.JoinType() == joinType {
			return true
		}
	}

	return false
}

// GroupAND adds a new condition group into current.
// Specified conditions joined into group with logical AND.
// Existed conditions are joined into group too.
// Both groups joined with logical AND operator.
// Returns new group containing required condition set.
func (conditionsGroup Group) GroupAND(conditions ...Condition) Group {
	return conditionsGroup.addConditions(LogicalAND, conditions...)
}

// And adds a new condition group into current.
// Specified conditions joined into group with logical AND.
// Existed conditions are joined into group too.
// Both groups joined with logical AND operator.
func (conditionsGroup Group) And(conditions ...Condition) Condition {
	return conditionsGroup.GroupAND(conditions...)
}

// GroupOR creates a new condition group.
// Specified conditions joined into group with logical AND.
// Existed conditions are joined into group too.
// Both groups joined with logical OR operator.
// Returns new group containing required condition set.
// Used internally to manage groups on externally called Or.
func (conditionsGroup Group) GroupOR(conditions ...Condition) Group {
	return conditionsGroup.addConditions(LogicalOR, conditions...)
}

// Or creates a new condition group.
// Specified conditions joined into group with logical AND.
// Existed conditions are joined into group too.
// Both groups joined with logical OR operator.
func (conditionsGroup Group) Or(conditions ...Condition) Condition {
	return conditionsGroup.GroupOR(conditions...)
}

// Render renders SQL SELECT clause part for current group.
// Takes existed parameters number (0 means no parameters are defined yet)
// Note Render renders group without brackets itself if called directly.
// Any child conditions groups are enclosed into brackets internally.
func (conditionsGroup Group) Render(parametersCount int) (sql string) {
	sql = ""

	if len(conditionsGroup.conditions) == 0 {
		return sql
	}

	for idx, condition := range conditionsGroup.conditions {
		joinToken := condition.RenderJoin(idx == 0)
		if len(joinToken) > 0 {
			sql += " " + joinToken + " "
		}

		sql += condition.Render(parametersCount)
		parametersCount += len(condition.Values())
	}

	if conditionsGroup.IsNegate() && len(conditionsGroup.conditions) > 1 { // force set brackets when negate
		conditionsGroup.withBrackets = true
	}

	if conditionsGroup.withBrackets { // if render with brackets
		sql = "(" + sql + ")"
	}

	if conditionsGroup.IsNegate() {
		return conditionsGroup.RenderNegate() + " " + sql
	}

	return sql
}

// RenderSQL renders SQL clause or its part.
// Implementation should render parameters substitutions using standard sql "?"(question) character.
func (conditionsGroup Group) RenderSQL() (sql string) {
	sql = ""

	if len(conditionsGroup.conditions) == 0 {
		return sql
	}

	for idx, condition := range conditionsGroup.conditions {
		joinToken := condition.RenderJoin(idx == 0)
		if len(joinToken) > 0 {
			sql += " " + joinToken + " "
		}

		sql += condition.RenderSQL()
	}

	if conditionsGroup.IsNegate() && len(conditionsGroup.conditions) > 1 { // force set brackets when negate
		conditionsGroup.withBrackets = true
	}

	if conditionsGroup.withBrackets { // if render with brackets
		sql = "(" + sql + ")"
	}

	if conditionsGroup.IsNegate() {
		return conditionsGroup.RenderNegate() + " " + sql
	}

	return sql
}

// Values returns a set of values of grouped conditions.
func (conditionsGroup Group) Values() (values []interface{}) {
	values = make([]interface{}, 0)

	for _, condition := range conditionsGroup.conditions {
		values = append(values, condition.Values()...)
	}

	return values
}

// ApplyFieldSpec makes a copy of Group with updated FieldDefinition on each child when field name matched.
// Implements Condition.
func (conditionsGroup Group) ApplyFieldSpec(spec FieldDefinition) Condition {
	for idx := range conditionsGroup.conditions {
		conditionsGroup.conditions[idx] = conditionsGroup.conditions[idx].ApplyFieldSpec(spec)
	}

	return conditionsGroup
}

// ApplyFieldTable makes a copy of Condition with updated FieldDefinition table name.
// Implements Condition.
func (conditionsGroup Group) ApplyFieldTable(table TableName) Condition {
	for idx := range conditionsGroup.conditions {
		conditionsGroup.conditions[idx] = conditionsGroup.conditions[idx].ApplyFieldTable(table)
	}

	return conditionsGroup
}
