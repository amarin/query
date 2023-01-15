package query

// TableName wraps string to identify table names.
type TableName string

// String returns table name as string.
func (tableName TableName) String() string {
	return string(tableName)
}

// InsertInto generates table inserter instance.
func (tableName TableName) InsertInto(insertValues ...FieldValue) InsertBuilder {
	return InsertInto(tableName).Values(insertValues...)
}

// Update generates UpdateBuilder instance.
// Use UpdateBuilder.Where to finalize UpdateBuilder instance configuration before use.
func (tableName TableName) Update(updateValues ...FieldValue) UpdateBuilder {
	return Update(tableName).Set(updateValues...)
}

// Delete generates DeleteBuilder instance.
// Use DeleteBuilder.Where to finalize DeleteBuilder instance configuration before use.
func (tableName TableName) Delete(filterConditions ...Condition) DeleteBuilder {
	return Delete(tableName).Where(filterConditions...)
}

// Ident generates a TableIdent instance from a TableName.
// See also TableOrError and Table constructors.
func (tableName TableName) Ident() TableIdent {
	return Table(tableName)
}

// Field generates a FieldDefinition instance having TableName set.
// It's a shorthand for Field(fieldName).Of(baseTable).
func (tableName TableName) Field(fieldName FieldName) FieldDefinition {
	return Field(fieldName).Of(tableName)
}

// RenderFrom returns TableName string value. Implements ClauseFromRenderer.
// Returns a string like "<tableName>".
func (tableName TableName) RenderFrom() string {
	return string(tableName)
}

// Select generates BaseSelectBuilder helper.
// Condition's will be used to filter values if specified.
// To get SelectSingleBuilder or SelectManyBuilder directly use SelectSingleFrom or SelectManyFrom instead.
func (tableName TableName) Select(filterConditions ...Condition) BaseSelectBuilder {
	return SelectFrom(tableName).Where(filterConditions...)
}

// SelectOne generates table query helper to fetch single value.
// Condition's will be used to filter values if specified.
// Shorthand to SelectSingleFrom((TableName)).Where(filterConditions...).
func (tableName TableName) SelectOne(filterConditions ...Condition) SelectSingleBuilder {
	return SelectSingleFrom(tableName).Where(filterConditions...)
}

// SelectMany generates SelectManyBuilder query helper to fetch many values at once.
// Condition's will be used to filter values if specified.
// Shorthand to SelectManyFrom((TableName)).Where(filterConditions...)
func (tableName TableName) SelectMany(filterConditions ...Condition) SelectManyBuilder {
	return SelectManyFrom(tableName).Where(filterConditions...)
}

// Join generates intermediate IncompleteSelectJoin instance.
// Takes TableName to join and TableJoinType constant defining required join type to produce.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Call IncompleteSelectJoin.On will return BaseSelectBuilder with join builder data finished.
func (tableName TableName) Join(rightTable TableIdent, joinType TableJoinType) IncompleteSelectJoin {
	return tableName.Select().Join(rightTable, joinType)
}

// InnerJoin generates intermediate IncompleteSelectJoin instance for InnerJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, InnerJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (tableName TableName) InnerJoin(rightTable TableIdent) IncompleteSelectJoin {
	return tableName.Join(rightTable, InnerJoin)
}

// LeftJoin generates intermediate IncompleteSelectJoin instance for LeftJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, LeftJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (tableName TableName) LeftJoin(rightTable TableIdent) IncompleteSelectJoin {
	return tableName.Join(rightTable, LeftJoin)
}

// RightJoin generates intermediate IncompleteSelectJoin instance for RightJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, RightJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (tableName TableName) RightJoin(rightTable TableIdent) IncompleteSelectJoin {
	return tableName.Join(rightTable, RightJoin)
}

// FullJoin generates intermediate IncompleteSelectJoin instance for FullJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, FullJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (tableName TableName) FullJoin(rightTable TableIdent) IncompleteSelectJoin {
	return tableName.Join(rightTable, FullJoin)
}

// // Count generates table rows count query builder.
// // Specified Condition's will be used to filter values.
// func (baseTable TableName) Count(filterConditions ...Condition) CountBuilder {
// 	return SelectManyFrom(baseTable, nil).Where(filterConditions...).Count()
// }
