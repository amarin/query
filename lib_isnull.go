package query

import (
	"strings"
)

type nullValue struct {
	BaseCondition
	FieldDefinition
}

// ApplyFieldSpec makes a copy of Condition with updated FieldDefinition if fieldName match.
// If FieldName is not matched it does nothing.
// Implements Condition.
func (impl nullValue) ApplyFieldSpec(spec FieldDefinition) Condition {
	if impl.FieldDefinition.fieldName == spec.fieldName {
		impl.FieldDefinition = spec
	}

	return impl
}

// ApplyFieldTable makes a copy of contains condition with updated FieldDefinition table name.
// Implements Condition.
func (impl nullValue) ApplyFieldTable(table TableName) Condition {
	impl.FieldDefinition.tableName = string(table)
	return impl
}

// Join returns a copy of Group having JoinType set to specified value.
func (impl nullValue) Join(newJoinType JoinType) Condition {
	impl.BaseCondition = impl.BaseCondition.Join(newJoinType)
	return impl
}

// Negate returns a copy of BaseCondition having IsNegate set to specified value.
func (impl nullValue) Negate(newNegateIndicator bool) Condition {
	impl.BaseCondition = impl.BaseCondition.Negate(newNegateIndicator)
	return impl
}

// Render renders IsNULL condition clause.
func (impl nullValue) Render(_ int) string {
	return impl.RenderSQL()
}

func (impl nullValue) RenderSQL() (sql string) {
	var tokens []string
	if impl.IsNegate() {
		tokens = []string{impl.RenderSpec(), "IS", "NOT", "NULL"}
	} else {
		tokens = []string{impl.RenderSpec(), "IS", "NULL"}
	}
	return strings.Join(tokens, " ")
}

// Values returns empty []interface{} slice as no substitutions required.
func (impl nullValue) Values() []interface{} {
	return []interface{}{}
}

// And generates new condition which true on all conditions met.
// Implements Condition.
func (impl nullValue) And(conditions ...Condition) Condition {
	return NewGroup(LogicalAND, impl).And(conditions...)
}

// Or generates new condition group which true on either initial condition is true or all of additional are true.
// Implements Condition.
func (impl nullValue) Or(conditions ...Condition) Condition {
	return NewGroup(LogicalAND, nil).And(impl).Or(conditions...)
}

// IsNull generates Condition to select rows where specified field name value IS NULL'.
func IsNull(fieldName FieldName) Condition {
	return &nullValue{
		BaseCondition:   *newBaseCondition(LogicalAND, false),
		FieldDefinition: Field(fieldName),
	}
}
