package query

import (
	"strconv"
)

// Limiter implements query result limiting clause builder and renderer.
// It encapsulates limit parameter to build a queries with expected result limitation,
// an offset parameter to use together with limit when building pagination handlers.
// It also stores received values of any or both of limit and offset internally to return values to parameters substitution
// when fully constructed SQL query is passed to execution.
type Limiter struct {
	limit  uint
	offset uint
}

// Offset returns a copy of Limiter having offset attribute set to specified value.
// Original Limiter instance remains the same.
func (query Limiter) Offset(offset uint) Limiter {
	query.offset = offset
	return query
}

// Limit returns a copy of Limiter having limit attribute set to specified value.
// Original Limiter instance remains the same.
// Set 0 to disable limit in query.
func (query Limiter) Limit(limit uint) Limiter {
	query.limit = limit
	return query
}

// Render renders Limiter provided SQL clause using numbered parameters.
// Takes existed parameters count (0 means no parameters are defined yet).
// Renders parameters substitutions starting from paramNum+1 using "$<number>" notation.
// Implements CountingClauseRenderer.
func (query Limiter) Render(paramNum int) (sql string) {
	sql = "" // start with empty string

	if query.offset > 0 {
		paramNum++
		sql += "OFFSET $" + strconv.Itoa(paramNum)
	}

	if query.limit > 0 {
		if len(sql) > 0 {
			sql += " "
		}
		paramNum++
		sql += "LIMIT $" + strconv.Itoa(paramNum)
	}

	return sql
}

// RenderSQL renders offset and limit clauses using default sql substitution with "?"(question) character.
// Implements RawClauseRenderer.
func (query Limiter) RenderSQL() (sql string) {
	sql = "" // start with empty string

	if query.offset > 0 {
		sql += "OFFSET ?"
	}

	if query.limit > 0 {
		if len(sql) > 0 {
			sql += " "
		}
		sql += "LIMIT ?"
	}

	return sql
}

// Values returns a set of Limiter values to substitute in SQL query generated using RenderSQL or Render methods.
// If either offset or limit attributes is equal to 0 its value is omitted.
// When both offset and limit are zeroes returns empty slice.
func (query Limiter) Values() (params []any) {
	params = make([]any, 0, 2) // max 2 parameters is expected
	switch {
	case query.offset > 0 && query.limit > 0:
		params = append(params, query.offset, query.limit)
	case query.offset > 0 && query.limit == 0:
		params = append(params, query.offset)
	case query.offset == 0 && query.limit > 0:
		params = append(params, query.limit)
	} // default params is empty set

	return params
}
