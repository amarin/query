package query

import (
	"fmt"
	"strings"
)

// ValueProvider implementations provides it's own values to database.
type ValueProvider interface {
	DatabaseValue() interface{}
}

// FieldName wraps string type to provide conditional clauses shortcuts.
type FieldName string

// Validate validates field name.
// Returns error if field name has unexpected characters or empty.
func (fn FieldName) Validate() error {
	fnStr := string(fn)

	switch {
	case strings.Contains(fnStr, "'"):
		return fmt.Errorf("%w: field name contains quote character", Error)
	case strings.Contains(fnStr, "\""):
		return fmt.Errorf("%w: field name contains double quote character", Error)
	case strings.Contains(fnStr, " "):
		return fmt.Errorf("%w: field name contains space character", Error)
	case strings.Count(fnStr, ".") > 1:
		return fmt.Errorf("%w: too many dots inside field name", Error)
	case len(fnStr) == 0:
		return fmt.Errorf("%w: empty field name", Error)
	case strings.HasSuffix(fnStr, "."):
		return fmt.Errorf("%w: field name could not end with dot", Error)
	case strings.HasPrefix(fnStr, "."):
		return fmt.Errorf("%w: field name could not start with dot", Error)
	}

	return nil
}

// String returns string value of FieldName.
func (fn FieldName) String() string {
	return string(fn)
}

// Contains generates field contains text Condition. Wraps Contains(string(*FieldName), value).
func (fn FieldName) Contains(value string) Condition {
	return Contains(fn, value)
}

// IContains generates field lookup confition comparing filed value using case-independent comparison.
// Wraps IContains(string(*FieldName), value).
func (fn FieldName) IContains(value string) Condition {
	return IContains(fn, value)
}

// EqualTo generates field equal to value Condition. Wraps EqualTo(string(*FieldName), value).
func (fn FieldName) EqualTo(value interface{}) Condition {
	return EqualTo(fn, value)
}

// In generates condition to select values where field value in supplied list. Wraps In(string(*FieldName), value).
func (fn FieldName) In(value interface{}) Condition {
	return In(fn, value)
}

// Value generates field value.
// If value is ValueProvider uses ValueProvider.DatabaseValue() instead if direct value.
func (fn FieldName) Value(value interface{}) FieldValue {
	switch typed := value.(type) {
	case ValueProvider:
		return *NewFieldValue(fn, typed.DatabaseValue())
	default:
		return *NewFieldValue(fn, value)
	}
}

// IsNull generates `field is null` Condition. Wraps IsNull(string(*FieldName), value).
func (fn FieldName) IsNull() Condition {
	return IsNull(fn)
}

// Field returns FieldDefinition build from FieldName to use in queries.
// Shorthand to Field((FieldName).fn).
func (fn FieldName) Field() FieldDefinition {
	return Field(fn)
}
