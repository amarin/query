package query

// SelectSingleBuilder extends BaseSelectBuilder helps to build SQL SELECT queries to receive single row.
type SelectSingleBuilder struct {
	BaseSelectBuilder
}

// RenderFrom renders a table name or list of tables joins to represent SQL clause FROM contents.
func (query SelectSingleBuilder) RenderFrom() string {
	return query.BaseSelectBuilder.RenderFrom()
}

// BuildQueryAndParams generates sql query string with desired parameters set.
// If query generation failed returns empty query and parameters set or non-nil error.
func (query SelectSingleBuilder) BuildQueryAndParams() (sql string, params []interface{}, err error) {
	sql, params, err = query.BaseSelectBuilder.BuildQueryAndParams()
	return sql + " LIMIT 1", params, nil
}

// FieldList returns spec list string with their possible aliases to build select query.
func (query SelectSingleBuilder) FieldList() string {
	return query.BaseSelectBuilder.fields.FieldList()
}

// Fields returns a copy of SelectManyBuilder having mustField list to retrieve updated with a list of specified FieldDefinition`s.
// Note all fields should be set one step as Fields call resets mustField added before.
func (query SelectSingleBuilder) Fields(fieldSpecs ...FieldDefinition) (updated SelectSingleBuilder) {
	updated.BaseSelectBuilder = updated.BaseSelectBuilder.Fields(fieldSpecs...)
	return updated
}

// FieldDefinitions returns a copy of attached FieldDefinition list.
func (query SelectSingleBuilder) FieldDefinitions() (res []FieldDefinition) {
	return query.BaseSelectBuilder.FieldDefinitions()
}

// Where adds fields conditions GroupAND returns modified SelectManyBuilder.
// If any conditions are already added, adds new conditions group joined with logical AND.
func (query SelectSingleBuilder) Where(fieldConditions ...Condition) SelectSingleBuilder {
	query.where = query.where.GroupAND(fieldConditions...)
	return query
}

// Update generates table UpdateBuilder.
// Note generated UpdateBuilder will use only base table even if join conditions added to SelectManyBuilder instance.
func (query SelectSingleBuilder) Update(values ...FieldValue) UpdateBuilder {
	return Update(query.BaseSelectBuilder.baseTable).Where(query.where.Conditions()...).Set(values...)
}

// BaseSelectBuilder query builder provides methods to generate SQL SELECT clauses having joined tables source.

// joinIdent generates intermediate IncompleteSelectJoin instance.
// Takes TableIdent to join and TableJoinType constant defining required join type to produce.
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectSingleBuilder) joinIdent(rightTable TableIdent, joinType TableJoinType) IncompleteSelectSingleJoiner {
	updated := query
	updated.joins = append(updated.joins, TableJoiner{rightTable: rightTable, joinType: joinType})
	return incompleteSingleJoiner{SelectSingleBuilder: updated}
}

// Join generates intermediate IncompleteSelectJoin instance.
// Takes TableName to join and TableJoinType constant defining required join type to produce.
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectSingleBuilder) Join(rightTable TableName, joinType TableJoinType) IncompleteSelectSingleJoiner {
	return query.joinIdent(Table(rightTable), joinType)
}

// InnerJoin generates intermediate IncompleteSelectJoin instance for InnerJoin type.
// Takes TableName to join.
// Shortcut to Join(rightTable TableName, InnerJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectSingleBuilder) InnerJoin(rightTable TableName) IncompleteSelectSingleJoiner {
	return query.Join(rightTable, InnerJoin)
}

// LeftJoin generates intermediate IncompleteSelectJoin instance for LeftJoin type.
// Takes TableName to join.
// Shortcut to Join(rightTable TableName, LeftJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectSingleBuilder) LeftJoin(rightTable TableName) IncompleteSelectSingleJoiner {
	return query.Join(rightTable, LeftJoin)
}

// RightJoin generates intermediate IncompleteSelectJoin instance for RightJoin type.
// Takes TableName to join.
// Shortcut to Join(rightTable TableName, RightJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectSingleBuilder) RightJoin(rightTable TableName) IncompleteSelectSingleJoiner {
	return query.Join(rightTable, RightJoin)
}

// FullJoin generates intermediate IncompleteSelectJoin instance for FullJoin type.
// Takes TableName to join.
// Shortcut to Join(rightTable TableName, FullJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query SelectSingleBuilder) FullJoin(rightTable TableName) IncompleteSelectSingleJoiner {
	return query.Join(rightTable, FullJoin)
}

// SelectSingleFromBase makes a new SelectSingleBuilder instance using supplied BaseSelectBuilder.
func SelectSingleFromBase(builder BaseSelectBuilder) SelectSingleBuilder {
	return SelectSingleBuilder{
		BaseSelectBuilder: builder,
	}
}

// SelectSingleFrom makes a new SelectSingleBuilder instance.
func SelectSingleFrom[T TableNameParameter](tableNameProvider T) SelectSingleBuilder {
	return SelectSingleFromBase(SelectFrom(tableNameProvider))
}
