package query

import (
	"strconv"
	"strings"
)

// equalTo implements any type fields conditions to select those records where field value is equal to specified value.
type equalTo struct {
	BaseCondition
	FieldValue
}

// ApplyFieldTable makes a copy of contains condition with updated FieldDefinition table name.
// Implements Condition.
func (impl equalTo) ApplyFieldTable(table TableName) Condition {
	impl.FieldValue.FieldDefinition.tableName = string(table)
	return impl
}

// ApplyFieldSpec makes a copy of Condition with updated FieldDefinition if fieldName match.
// If FieldName is not matched it does nothing.
// Implements Condition.
func (impl equalTo) ApplyFieldSpec(spec FieldDefinition) Condition {
	impl.FieldValue = impl.FieldValue.ApplyFieldSpec(spec)
	return impl
}

// Join returns a copy of Group having JoinType set to specified value.
func (impl equalTo) Join(newJoinType JoinType) Condition {
	impl.BaseCondition = impl.BaseCondition.Join(newJoinType)
	return impl
}

// Negate returns a copy of BaseCondition having IsNegate set to specified value.
func (impl equalTo) Negate(newNegateIndicator bool) Condition {
	impl.BaseCondition = impl.BaseCondition.Negate(newNegateIndicator)
	return impl
}

// Render renders SQL SELECT clause part for current field.
func (impl equalTo) Render(paramNum int) string {
	var tokens []string
	if impl.IsNegate() {
		tokens = []string{impl.RenderNegate() + " ", impl.RenderSpec(), "=$" + strconv.Itoa(paramNum+1)}
	} else {
		tokens = []string{impl.RenderSpec(), "=$" + strconv.Itoa(paramNum+1)}
	}
	return strings.Join(tokens, "")
}

// RenderSQL renders SQL clause or its part.
// Implementation should render parameters substitutions using standard sql "?"(question) character.
func (impl equalTo) RenderSQL() (sql string) {
	var tokens []string
	if impl.IsNegate() {
		tokens = []string{impl.RenderNegate() + " ", impl.RenderSpec(), "=?"}
	} else {
		tokens = []string{impl.RenderSpec(), "=?"}
	}
	return strings.Join(tokens, "")
}

// And generates new condition which true on all conditions met.
// Implements Condition.
func (impl equalTo) And(conditions ...Condition) Condition {
	return And(impl, conditions...)
}

// Or generates new condition group which true on either initial condition is true or all of additional are true.
// Implements Condition.
func (impl equalTo) Or(conditions ...Condition) Condition {
	return Or(impl, conditions...)
}

// EqualTo generates Condition for any field type equalTo to value.
func EqualTo(fieldName FieldName, value interface{}) Condition {
	return &equalTo{
		BaseCondition: *newBaseCondition(LogicalAND, false),
		FieldValue:    *NewFieldValue(fieldName, value),
	}
}
