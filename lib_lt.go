package query

import (
	"strconv"
	"strings"
)

const (
	lessSymbol = "<"
)

// less implements any type fields conditions to match records where field value is less than specified value.
type less struct {
	BaseCondition
	FieldValue
}

// ApplyFieldTable makes a copy of contains condition with updated FieldDefinition table name.
// Implements Condition.
func (impl less) ApplyFieldTable(table TableName) Condition {
	impl.FieldValue.FieldDefinition.tableName = string(table)
	return impl
}

// ApplyFieldSpec makes a copy of Condition with updated FieldDefinition if fieldName match.
// If FieldName is not matched it does nothing.
// Implements Condition.
func (impl less) ApplyFieldSpec(spec FieldDefinition) Condition {
	impl.FieldValue = impl.FieldValue.ApplyFieldSpec(spec)
	return impl
}

// Join returns a copy of Group having JoinType set to specified value.
func (impl less) Join(newJoinType JoinType) Condition {
	impl.BaseCondition = impl.BaseCondition.Join(newJoinType)
	return impl
}

// Negate returns a copy of BaseCondition having IsNegate set to specified value.
func (impl less) Negate(newNegateIndicator bool) Condition {
	impl.BaseCondition = impl.BaseCondition.Negate(newNegateIndicator)
	return impl
}

// Render renders SQL SELECT clause part for current field.
func (impl less) Render(paramNum int) string {
	var tokens []string
	if impl.IsNegate() {
		tokens = []string{impl.RenderNegate() + " ", impl.RenderSpec(), lessSymbol, "$" + strconv.Itoa(paramNum+1)}
	} else {
		tokens = []string{impl.RenderSpec(), lessSymbol, "$" + strconv.Itoa(paramNum+1)}
	}
	return strings.Join(tokens, "")
}

// RenderSQL renders SQL clause or its part.
// Implementation should render parameters substitutions using standard sql "?"(question) character.
func (impl less) RenderSQL() (sql string) {
	var tokens []string
	if impl.IsNegate() {
		tokens = []string{impl.RenderNegate() + " ", impl.RenderSpec(), lessSymbol, "?"}
	} else {
		tokens = []string{impl.RenderSpec(), lessSymbol, "?"}
	}
	return strings.Join(tokens, "")
}

// And generates new condition which true on all conditions met.
// Implements Condition.
func (impl less) And(conditions ...Condition) Condition {
	return And(impl, conditions...)
}

// Or generates new condition group which true on either initial condition is true or all of additional are true.
// Implements Condition.
func (impl less) Or(conditions ...Condition) Condition {
	return Or(impl, conditions...)
}

// Less generates Condition for any field type to match records having field values less than specified.
func Less(fieldName FieldName, value interface{}) Condition {
	return &less{
		BaseCondition: *newBaseCondition(LogicalAND, false),
		FieldValue:    *NewFieldValue(fieldName, value),
	}
}
