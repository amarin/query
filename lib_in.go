package query

import (
	"strconv"
	"strings"
)

const (
	inSymbol = "IN"
)

// in implements any type fields conditions to match records where field value is in than specified value.
type in struct {
	BaseCondition
	FieldValue
}

// ApplyFieldTable makes a copy of contains condition with updated FieldDefinition table name.
// Implements Condition.
func (impl in) ApplyFieldTable(table TableName) Condition {
	impl.FieldValue.FieldDefinition.tableName = string(table)
	return impl
}

// ApplyFieldSpec makes a copy of Condition with updated FieldDefinition if fieldName match.
// If FieldName is not matched it does nothing.
// Implements Condition.
func (impl in) ApplyFieldSpec(spec FieldDefinition) Condition {
	impl.FieldValue = impl.FieldValue.ApplyFieldSpec(spec)
	return impl
}

// Join returns a copy of Group having JoinType set to specified value.
func (impl in) Join(newJoinType JoinType) Condition {
	impl.BaseCondition = impl.BaseCondition.Join(newJoinType)
	return impl
}

// Negate returns a copy of BaseCondition having IsNegate set to specified value.
func (impl in) Negate(newNegateIndicator bool) Condition {
	impl.BaseCondition = impl.BaseCondition.Negate(newNegateIndicator)
	return impl
}

// Render renders SQL SELECT clause part for current field.
func (impl in) Render(paramNum int) string {
	var tokens []string

	placeholders := make([]string, len(impl.Values()))

	for idx := range impl.Values() {
		placeholders[idx] = "$" + strconv.Itoa(paramNum+idx+1)
	}

	placeholdersString := "(" + strings.Join(placeholders, ",") + ")"

	if impl.IsNegate() {
		tokens = []string{impl.RenderSpec(), impl.RenderNegate(), inSymbol, placeholdersString}
	} else {
		tokens = []string{impl.RenderSpec(), inSymbol, placeholdersString}
	}

	return strings.Join(tokens, " ")
}

// RenderSQL renders SQL clause or its part.
// Implementation should render parameters substitutions using standard sql "?"(question) character.
func (impl in) RenderSQL() (sql string) {
	var tokens []string

	placeholders := make([]string, len(impl.Values()))

	for idx := range impl.Values() {
		placeholders[idx] = "?"
	}

	placeholdersString := "(" + strings.Join(placeholders, ",") + ")"

	if impl.IsNegate() {
		tokens = []string{impl.RenderSpec(), impl.RenderNegate(), inSymbol, placeholdersString}
	} else {
		tokens = []string{impl.RenderSpec(), inSymbol, placeholdersString}
	}

	return strings.Join(tokens, " ")
}

// And generates new condition which true on all conditions met.
// Implements Condition.
func (impl in) And(conditions ...Condition) Condition {
	return And(impl, conditions...)
}

// Or generates new condition group which true on either initial condition is true or all of additional are true.
// Implements Condition.
func (impl in) Or(conditions ...Condition) Condition {
	return Or(impl, conditions...)
}

// In generates Condition for any field type to match records having field values in than specified.
func In(fieldName FieldName, value interface{}) Condition {
	return &in{
		BaseCondition: *newBaseCondition(LogicalAND, false),
		FieldValue:    *NewFieldValue(fieldName, value),
	}
}
