package query

// BaseBuilder defines a base data structure useful for any queries type.
type BaseBuilder struct {
	op Operation // operation aka SELECT, INSERT, UPDATE or DELETE
}

// Operation returns SQL operation of query.
// If not set DoSelect used by default.
func (b BaseBuilder) Operation() Operation {
	return b.op
}
