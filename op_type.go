package query

import (
	"strconv"
)

// Operation defines base query operation type, i.e. SELECT, INSERT, DELETE or UPDATE.
type Operation int

const (
	// DoSelect defines a constant SQL SELECT Operation type.
	DoSelect Operation = iota
	// DoInsert defines a constant SQL INSERT Operation type.
	DoInsert
	// DoUpdate defines a constant SQL UPDATE Operation type.
	DoUpdate
	// DoDelete defines a constant SQL DELETE Operation type.
	DoDelete
)

// String returns a string representation of SQL query operation.
// It equals either to "SELECT", "INSERT", "UPDATE" or "DELETE" when valid Operation used.
// If invalid returns "unknown(<int>)".
func (op Operation) String() string {
	switch op {
	case DoSelect:
		return kwSelect.String()
	case DoInsert:
		return kwInsert.String()
	case DoUpdate:
		return kwUpdate.String()
	case DoDelete:
		return kwDelete.String()
	default:
		return "unknown(" + strconv.Itoa(int(op)) + ")"
	}
}

// New creates new BaseBuilder of desired operation.
func (op Operation) New() BaseBuilder {
	return BaseBuilder{op: op}
}
