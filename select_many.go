package query

import (
	"strconv"
)

const (
	noLimit int = -1
)

// SelectManyBuilder extends BaseSelectBuilder helps to build SQL SELECT queries with ordering, offset and limit.
type SelectManyBuilder struct {
	BaseSelectBuilder
	offset    uint
	limit     int
	order     []FieldSorting
	tableSpec bool // render query with field spec including table names
}

// RenderFrom renders a table name or list of tables joins to represent SQL clause FROM contents.
func (query SelectManyBuilder) RenderFrom() string {
	return query.BaseSelectBuilder.RenderFrom()
}

// BuildQueryAndParams generates sql query string with desired parameters set.
// If query generation failed returns empty query and parameters set or non-nil error.
func (query SelectManyBuilder) BuildQueryAndParams() (sql string, params []interface{}, err error) {
	sql, params, err = query.BaseSelectBuilder.BuildQueryAndParams()

	if len(query.order) > 0 {
		sql += " ORDER BY "
		for idx, order := range query.order {
			if idx > 0 {
				sql += ", "
			}
			sql += order.Render()
		}
	}

	if query.offset > 0 {
		params = append(params, query.offset)
		sql += " OFFSET $" + strconv.Itoa(len(params))
	}

	if query.limit > 0 {
		params = append(params, uint(query.limit))
		sql += " LIMIT $" + strconv.Itoa(len(params))

	}

	return sql, params, nil
}

// FieldList returns spec list string with their possible aliases to build select query.
func (query SelectManyBuilder) FieldList() string {
	return query.BaseSelectBuilder.fields.FieldList()
}

// Fields returns a copy of SelectManyBuilder having mustField list to retrieve updated with a list of specified FieldDefinition`s.
// Note all fields should be set one step as Fields call resets mustField added before.
func (query SelectManyBuilder) Fields(fieldSpecs ...FieldDefinition) (updated SelectManyBuilder) {
	updated = query
	updated.BaseSelectBuilder = updated.BaseSelectBuilder.Fields(fieldSpecs...)

	return updated
}

// FieldDefinitions returns a copy of attached FieldDefinition list.
func (query SelectManyBuilder) FieldDefinitions() (res []FieldDefinition) {
	return query.BaseSelectBuilder.FieldDefinitions()
}

// applySpecToOrdering checks if field spec suits for any ordering fields and replaces its spec if required
func (query SelectManyBuilder) applySpecToOrdering(spec FieldDefinition) SelectManyBuilder {
	for idx, ordering := range query.order {
		if spec.fieldName == ordering.fieldName {
			query.order[idx] = query.order[idx].ApplyFieldSpec(spec)
		}
	}

	return query
}

// applySpecToOrdering checks if field spec suits for any ordering fields and replaces its spec if required
func (query SelectManyBuilder) applyFieldTableToOrdering(tableName TableName) SelectManyBuilder {
	for idx := range query.order {
		query.order[idx] = query.order[idx].applyFieldTable(tableName)
	}

	return query
}

// Offset sets offset GroupAND returns modified SelectManyBuilder.
func (query SelectManyBuilder) Offset(offset uint) SelectManyBuilder {
	query.offset = offset
	return query
}

// Limit sets limit GroupAND returns modified SelectManyBuilder. Values -1 to disable limit in query.
func (query SelectManyBuilder) Limit(limit int) SelectManyBuilder {
	query.limit = limit
	return query
}

// OrderBy adds ordering by fields GroupAND returns modified SelectManyBuilder.
func (query SelectManyBuilder) OrderBy(orderByFields ...FieldSorting) SelectManyBuilder {
	query.order = append(query.order, orderByFields...)
	return query
}

// Where adds fields conditions GroupAND returns modified SelectManyBuilder.
// If any conditions are already added, adds new conditions group joined with logical AND.
func (query SelectManyBuilder) Where(fieldConditions ...Condition) SelectManyBuilder {
	query.where = query.where.GroupAND(fieldConditions...)
	return query
}

// Update generates table UpdateBuilder.
// Note generated UpdateBuilder will use only base table even if join conditions added to SelectManyBuilder instance.
func (query SelectManyBuilder) Update(values ...FieldValue) UpdateBuilder {
	return Update(query.BaseSelectBuilder.TableName()).Where(query.where.Conditions()...).Set(values...)
}

// Delete generates table DeleteBuilder.
// Note generated DeleteBuilder will use only base table even if join conditions added to SelectManyBuilder instance.
func (query SelectManyBuilder) Delete() DeleteBuilder {
	return Delete(query.BaseSelectBuilder.TableName()).Where(query.where.Conditions()...)
}

// BaseSelectBuilder query builder provides methods to generate SQL SELECT clauses having joined tables source.

// joinIdent generates intermediate IncompleteSelectJoin instance.
// Takes TableIdent to join and TableJoinType constant defining required join type to produce.
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectManyBuilder) joinIdent(rightTable TableIdent, joinType TableJoinType) IncompleteSelectManyJoiner {
	updated := query
	updated.joins = append(updated.joins, TableJoiner{rightTable: rightTable, joinType: joinType})
	return incompleteManyJoiner{SelectManyBuilder: updated}
}

// Join generates intermediate IncompleteSelectJoin instance.
// Takes TableName to join and TableJoinType constant defining required join type to produce.
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectManyBuilder) Join(rightTable TableName, joinType TableJoinType) IncompleteSelectManyJoiner {
	return query.joinIdent(Table(rightTable), joinType)
}

// InnerJoin generates intermediate IncompleteSelectJoin instance for InnerJoin type.
// Takes TableName to join.
// Shortcut to Join(rightTable TableName, InnerJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectManyBuilder) InnerJoin(rightTable TableName) IncompleteSelectManyJoiner {
	return query.Join(rightTable, InnerJoin)
}

// LeftJoin generates intermediate IncompleteSelectJoin instance for LeftJoin type.
// Takes TableName to join.
// Shortcut to Join(rightTable TableName, LeftJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectManyBuilder) LeftJoin(rightTable TableName) IncompleteSelectManyJoiner {
	return query.Join(rightTable, LeftJoin)
}

// RightJoin generates intermediate IncompleteSelectJoin instance for RightJoin type.
// Takes TableName to join.
// Shortcut to Join(rightTable TableName, RightJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectManyBuilder) RightJoin(rightTable TableName) IncompleteSelectManyJoiner {
	return query.Join(rightTable, RightJoin)
}

// FullJoin generates intermediate IncompleteSelectJoin instance for FullJoin type.
// Takes TableName to join.
// Shortcut to Join(rightTable TableName, FullJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectManyBuilder) FullJoin(rightTable TableName) IncompleteSelectManyJoiner {
	return query.Join(rightTable, FullJoin)
}

// Count creates an CountBuilder using parameters of this SelectManyBuilder.
func (query SelectManyBuilder) Count() CountBuilder {
	return Count(query.BaseSelectBuilder)
}

// SelectManyFromBase makes a new query SelectManyBuilder instance using supplied BaseSelectBuilder.
func SelectManyFromBase(builder BaseSelectBuilder) SelectManyBuilder {
	return SelectManyBuilder{
		BaseSelectBuilder: builder,
		offset:            0,
		limit:             noLimit,
		order:             make([]FieldSorting, 0),
		tableSpec:         false, // table prefix in mustField names is not required unless building JOIN
	}
}

// SelectManyFrom makes a new query SelectManyBuilder instance.
func SelectManyFrom[T TableNameParameter](tableNameProvider T) SelectManyBuilder {
	return SelectManyBuilder{
		BaseSelectBuilder: SelectFrom(tableNameProvider),
		offset:            0,
		limit:             noLimit,
		order:             make([]FieldSorting, 0),
		tableSpec:         false, // table prefix in mustField names is not required unless building JOIN
	}
}
