package query

import (
	"strconv"
	"strings"
)

// contains implements string fields compare using LIKE operator.
// It adds <field_name> LIKE '%<value>%' to SQL SELECT clause.
type contains struct {
	BaseCondition
	FieldValue
}

// ApplyFieldTable makes a copy of contains condition with updated FieldDefinition table name.
// Implements Condition.
func (impl contains) ApplyFieldTable(table TableName) Condition {
	impl.FieldValue = impl.FieldValue.ApplyFieldTable(table)
	return impl
}

// Join returns a copy of Group having JoinType set to specified value.
func (impl contains) Join(newJoinType JoinType) Condition {
	impl.BaseCondition = impl.BaseCondition.Join(newJoinType)
	return impl
}

// Negate returns a copy of BaseCondition having IsNegate set to specified value.
func (impl contains) Negate(newNegateIndicator bool) Condition {
	impl.BaseCondition = impl.BaseCondition.Negate(newNegateIndicator)
	return impl
}

// ApplyFieldSpec makes a copy of Condition with updated FieldDefinition if fieldName match.
// If FieldName is not matched it does nothing.
// Implements Condition.
func (impl contains) ApplyFieldSpec(spec FieldDefinition) Condition {
	impl.FieldValue = impl.FieldValue.ApplyFieldSpec(spec)
	return impl
}

// And generates new condition which true on all conditions met.
// Implements Condition.
func (impl contains) And(conditions ...Condition) Condition {
	return NewGroup(LogicalAND, impl).And(conditions...)
}

// Or generates new condition group which true on either initial condition is true or all of additional are true.
// Implements Condition.
func (impl contains) Or(conditions ...Condition) Condition {
	return NewGroup(LogicalOR, impl).Or(conditions...)
}

// Render renders SQL SELECT clause part for current mustField.
// Takes existed parameters count (0 means no parameters are defined yet).
func (impl contains) Render(paramNum int) string {
	var tokens []string
	if impl.IsNegate() {
		tokens = []string{impl.RenderSpec(), impl.RenderNegate(), "LIKE", "$" + strconv.Itoa(paramNum+1)}
	} else {
		tokens = []string{impl.RenderSpec(), "LIKE", "$" + strconv.Itoa(paramNum+1)}
	}

	return strings.Join(tokens, " ")
}

// RenderSQL renders SQL clause or its part.
// Implementation should render parameters substitutions using standard sql "?"(question) character.
func (impl contains) RenderSQL() (sql string) {
	var tokens []string
	if impl.IsNegate() {
		tokens = []string{impl.RenderSpec(), impl.RenderNegate(), "LIKE", "?"}
	} else {
		tokens = []string{impl.RenderSpec(), "LIKE", "?"}
	}

	return strings.Join(tokens, " ")
}

// Values provides single string value wrapped to percent sign to fill SQL LIKE clause.
func (impl contains) Values() []interface{} {
	return []interface{}{"%" + (impl.value).(string) + "%"}
}

// Contains generates Condition to compare string like fields using LIKE '%value%'.
// See IContains condition generator to make case-independent `contains` comparison.
func Contains(fieldName FieldName, value string) Condition {
	return &contains{
		BaseCondition: *newBaseCondition(LogicalAND, false),
		FieldValue:    *NewFieldValue(fieldName, value),
	}
}
