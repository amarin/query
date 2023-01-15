package query

import (
	"fmt"
)

// SortDirection wraps string type to define ordering directions in SelectManyBuilder.
type SortDirection string

// String returns string value of SortDirection.
func (direction SortDirection) String() string {
	return string(direction)
}

const (
	// Ascending defines constant value to indicate ascending ordering required.
	Ascending SortDirection = "ASC"

	// Descending defines constant value to indicate descending ordering required.
	Descending SortDirection = "DESC"
)

// FieldSorting groups together field definition GroupAND SortDirection to build ORDER BY clause in SQL SELECT requests.
type FieldSorting struct {
	FieldDefinition
	direction SortDirection
}

// Direction returns current sorting direction set.
func (o FieldSorting) Direction() SortDirection {
	return o.direction
}

// Render makes an ORDER BY string.
func (o FieldSorting) Render() string {
	return o.FieldDefinition.RenderField() + " " + string(o.direction)
}

// ApplyFieldSpec returns a copy of FieldSorting item with FieldDefinition updated if field name matches.
// If mustField names of original FieldSorting item and argument are differs simply returns a copy of original FieldSorting.
func (o FieldSorting) ApplyFieldSpec(spec FieldDefinition) FieldSorting {
	if o.FieldDefinition.fieldName == spec.fieldName {
		o.FieldDefinition = spec
	}

	return o
}

// applyFieldTable returns a copy of FieldSorting item with FieldDefinition table name updated.
func (o FieldSorting) applyFieldTable(table TableName) FieldSorting {
	o.FieldDefinition.tableName = string(table)

	return o
}

// ASC generates new FieldSorting to build ordering clause with Ascending order by specified field name.
// Note invalid field name leads to panic.
// Use FieldName.Validate before build ordering when taking field name from insecure environment.
// To choose ordering direction programmatically use OrderBy constructor instead.
func ASC(fieldName FieldName) FieldSorting {
	if err := fieldName.Validate(); err != nil {
		panic(fmt.Sprintf("invalid field name: %v", err))
	}
	return FieldSorting{
		FieldDefinition: Field(fieldName),
		direction:       Ascending,
	}
}

// DESC generates new FieldSorting to build ordering clause with Descending order by specified field name.
// Use FieldName.Validate before build ordering when taking field name from insecure environment.
func DESC(fieldName FieldName) FieldSorting {
	if err := fieldName.Validate(); err != nil {
		panic(fmt.Sprintf("invalid field name: %v", err))
	}
	return FieldSorting{
		FieldDefinition: Field(fieldName),
		direction:       Descending,
	}
}

// OrderBy creates FieldSorting instance using field name.
// If optional direction specified it should be either Ascending or Descending, default is Ascending.
// When ordering is known during development use ASC or DESC constructors instead.
// Panics if unexpected order direction requested or field name to order is invalid.
func OrderBy(fieldName FieldName, direction ...SortDirection) FieldSorting {
	switch {
	case len(direction) > 0 && direction[0] == Ascending:
		return ASC(fieldName)
	case len(direction) > 0 && direction[0] == Descending:
		return DESC(fieldName)
	case len(direction) == 0:
		return ASC(fieldName)
	default:
		panic(fmt.Errorf("%w: unexpected order: %v: %v", Error, fieldName, direction))
	}
}
