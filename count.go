package query

import (
	"strings"
)

// CountBuilder helps to build count queries.
// It utilizes BaseSelectBuilder methods
type CountBuilder struct {
	baseBuilder BaseSelectBuilder
	countField  FieldDefinition
}

// RenderFrom returns string representation of table name or tables join with possible tables aliases.
// Implements ClauseFromRenderer.
func (query CountBuilder) RenderFrom() string {
	return query.baseBuilder.RenderFrom()
}

// BuildQueryAndParams generates sql query string with desired parameters set.
// If query generation failed returns empty query and parameters set or non-nil error.
func (query CountBuilder) BuildQueryAndParams() (sql string, params []interface{}, err error) {
	tokens := append([]string{},
		DoSelect.String(),
		kwCount.String()+"("+query.countField.RenderField()+")",
		kwFrom.String(),
		query.baseBuilder.RenderFrom(),
	)

	if len(query.baseBuilder.where.Conditions()) > 0 {
		tokens = append(tokens, kwWhere.String(), query.baseBuilder.where.Render(0))
		params = query.baseBuilder.where.Values()
	}

	return strings.Join(tokens, " "), params, nil
}

// Count prepares SQL SELECT COUNT query builder.
// Takes BaseSelectBuilder instance and optional mustField name to count over it.
// If no mustField name specified default '*' will used.
// Returns CountBuilder instance.
func Count(query BaseSelectBuilder, fieldToCount ...FieldName) (countBuilder CountBuilder) {
	countBuilder = CountBuilder{baseBuilder: query}

	if len(fieldToCount) > 0 {
		countBuilder.countField = Field(fieldToCount[0])
	} else {
		countBuilder.countField = Field("*")
	}

	return countBuilder
}
