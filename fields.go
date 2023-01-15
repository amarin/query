package query

import (
	"fmt"
	"strings"
)

// Fields handles list of fields to build SQL SELECT clause.
type Fields struct {
	fieldSpecs []FieldDefinition // fields to select
}

// String returns a string representation of Fields.
func (query Fields) String() string {
	fieldStrings := make([]string, len(query.fieldSpecs))
	for idx, field := range query.fieldSpecs {
		fieldStrings[idx] = field.RenderSpec()
	}
	return "Fields(" + strings.Join(fieldStrings, ",") + ")"
}

// Len returns a length or fields added to Fields.
// Zero length means no fields explicitly defined and '*' will used instead.
func (query Fields) Len() int {
	return len(query.fieldSpecs)
}

// Fields returns a copy of Fields having mustField list to retrieve updated with a list of specified FieldDefinition`s.
func (query Fields) Fields(fieldSpecs ...FieldDefinition) (updated Fields) {
	updated = query
	updated.fieldSpecs = make([]FieldDefinition, len(fieldSpecs)) // reset added fields
	for idx, fieldSpec := range fieldSpecs {
		updated.fieldSpecs[idx] = fieldSpec
	}

	return updated
}

// FieldDefinitions returns attached FieldDefinition list copy.
func (query Fields) FieldDefinitions() (res []FieldDefinition) {
	res = make([]FieldDefinition, len(query.fieldSpecs))
	copy(res, query.fieldSpecs)

	return res
}

// FieldList returns spec list string with their possible aliases to build select query.
func (query Fields) FieldList() string {
	if len(query.fieldSpecs) == 0 {
		return "*"
	}

	fieldSpecs := make([]string, len(query.fieldSpecs))
	for idx, spec := range query.fieldSpecs {
		fieldSpecs[idx] = spec.RenderSpec()
	}

	return strings.Join(fieldSpecs, ", ")
}

// NewFields makes new fields list.
// Takes a slice of strings, FieldName or FieldDefinition.
// Any unexpected type argument will cause an error.
func NewFields(args ...any) Fields {
	fieldSpecs := make([]FieldDefinition, len(args))

	for idx := range args {
		switch value := args[idx].(type) {
		case FieldDefinition:
			fieldSpecs[idx] = value
		case FieldName:
			fieldSpecs[idx] = Field(value)
		case string:
			fieldSpecs[idx] = Field(value)
		default:
			panic(fmt.Sprintf("only string, FieldName or FieldDefinition types allowed"))
		}
	}

	return Fields{fieldSpecs: fieldSpecs}
}
