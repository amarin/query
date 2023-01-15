package query

import (
	"fmt"
	"strings"
)

// BaseSelectBuilder helps to build table queries.
// It does not provide pagination or ordering itself
// but contains helper methods to generate SelectListBuilder or SelectOneBuilder.
// It also provides method Count() to generate SelectCountBuilder.
type BaseSelectBuilder struct {
	BaseBuilder
	baseTable TableIdent
	joins     []TableJoiner
	fields    Fields
	where     WhereClause
}

// String returns a string representation of BaseSelectBuilder.
func (query BaseSelectBuilder) String() string {
	return strings.Join([]string{
		fmt.Sprintf("%T(%v)", query, query.RenderFrom()),
		query.fields.String(),
		query.where.String(),
	}, ",")
}

func (query BaseSelectBuilder) RenderSQL() (sql string) {
	tokens := append([]string{},
		DoSelect.String(),
		query.fields.FieldList(),
		kwFrom.String(),
		query.baseTable.RenderFrom(),
	)

	for _, joiner := range query.joins {
		tokens = append(tokens, joiner.RenderFrom())
	}

	if len(query.where.Conditions()) > 0 {
		tokens = append(tokens, kwWhere.String(), query.where.RenderSQL())
	}

	return strings.Join(tokens, " ")
}

func (query BaseSelectBuilder) Render(parametersCount int) (sql string) {
	tokens := append([]string{},
		DoSelect.String(),
		query.fields.FieldList(),
		kwFrom.String(),
		query.baseTable.RenderFrom(),
	)

	for _, joiner := range query.joins {
		tokens = append(tokens, joiner.RenderFrom())
	}

	if len(query.where.Conditions()) > 0 {
		tokens = append(tokens, kwWhere.String(), query.where.Render(parametersCount))
	}

	return strings.Join(tokens, " ")
}

func (query BaseSelectBuilder) Values() (params []any) {
	params = make([]any, 0)

	if len(query.where.Conditions()) > 0 {
		params = append(params, query.where.Values()...)
	}

	return params
}

func (query BaseSelectBuilder) RenderFrom() (fromClause string) {
	fromClauseItems := make([]string, 1+len(query.joins))
	fromClauseItems[0] = query.baseTable.RenderFrom()

	for idx, tableJoiner := range query.joins {
		fromClauseItems[1+idx] = tableJoiner.RenderFrom()
	}

	return strings.Join(fromClauseItems, " ")
}

// BuildQueryAndParams generates sql query string with desired parameters set.
// If query generation failed returns empty query and parameters set or non-nil error.
func (query BaseSelectBuilder) BuildQueryAndParams() (sql string, params []interface{}, err error) {
	params = make([]any, 0)
	tokens := append([]string{},
		DoSelect.String(),
		query.fields.FieldList(),
		kwFrom.String(),
		query.baseTable.RenderFrom(),
	)

	for _, joiner := range query.joins {
		tokens = append(tokens, joiner.RenderFrom())
	}

	if len(query.where.Conditions()) > 0 {
		tokens = append(tokens, kwWhere.String(), query.where.Render(0))
		params = append(params, query.where.Values()...)
	}

	return strings.Join(tokens, " "), params, nil
}

// TableName returns table name to fetch records from.
func (query BaseSelectBuilder) TableName() TableName {
	return query.baseTable.TableName()
}

// Fields returns a copy of SelectManyBuilder having mustField list to retrieve updated with a list of specified FieldDefinition`s.
func (query BaseSelectBuilder) Fields(fieldSpecs ...FieldDefinition) (updated BaseSelectBuilder) {
	updated = query
	updated.fields = updated.fields.Fields(fieldSpecs...)

	return updated
}

// FieldDefinitions returns a copy of attached FieldDefinition list.
func (query BaseSelectBuilder) FieldDefinitions() (res []FieldDefinition) {
	return query.fields.FieldDefinitions()
}

// Where adds fields conditions GroupAND returns modified SelectManyBuilder.
// If any conditions are already added, adds new conditions group joined with logical AND.
func (query BaseSelectBuilder) Where(fieldConditions ...Condition) (updated BaseSelectBuilder) {
	updated = query
	updated.where = updated.where.GroupAND(fieldConditions...)

	return updated
}

// InsertInto generates SQL INSERT query builder using stored table name.
// Insert values are optional and could be set later with InsertBuilder.Values.
// Note generated InsertInto will receive only base TableIdent to generate insert into it.
func (query BaseSelectBuilder) InsertInto(insertValues ...FieldValue) InsertBuilder {
	return InsertInto(query.baseTable).Values(insertValues...)
}

// Update generates table UpdateBuilder.
// Note generated UpdateBuilder will receive only base TableIdent to generate update on it.
func (query BaseSelectBuilder) Update(values ...FieldValue) UpdateBuilder {
	return Update(query.baseTable).Where(query.where.Conditions()...).Set(values...)
}

// Delete generates table DeleteBuilder.
// Note generated DeleteBuilder will receive only base TableIdent to generate update on it.
func (query BaseSelectBuilder) Delete() DeleteBuilder {
	return Delete(query.baseTable).Where(query.where.Conditions()...)
}

// Single makes a SelectSingleBuilder instance from BaseSelectBuilder.
func (query BaseSelectBuilder) Single() SelectSingleBuilder {
	return SelectSingleBuilder{BaseSelectBuilder: query}
}

// Many makes a SelectManyBuilder instance from BaseSelectBuilder setting ordering and pagination to its default values.
// See also Offset, Limit and OrderBy.
func (query BaseSelectBuilder) Many() SelectManyBuilder {
	return SelectManyFromBase(query)
}

// Offset makes SelectManyBuilder instance with offset set to specified value.
func (query BaseSelectBuilder) Offset(offset uint) SelectManyBuilder {
	return query.Many().Offset(offset)
}

// Limit makes a SelectManyBuilder instance from BaseSelectBuilder setting required records limit to specified value.
func (query BaseSelectBuilder) Limit(limit int) SelectManyBuilder {
	return query.Many().Limit(limit)
}

// OrderBy makes a SelectManyBuilder instance from BaseSelectBuilder setting ordering to specified value.
func (query BaseSelectBuilder) OrderBy(orderByFields ...FieldSorting) SelectManyBuilder {
	return query.Many().OrderBy(orderByFields...)
}

// BaseSelectBuilder query builder provides methods to generate SQL SELECT clauses having joined tables source.

// joinIdent generates intermediate IncompleteSelectJoin instance.
// Takes TableIdent to join and TableJoinType constant defining required join type to produce.
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query BaseSelectBuilder) joinIdent(rightTable TableIdent, joinType TableJoinType) IncompleteSelectJoin {
	updated := query
	updated.joins = append(updated.joins, TableJoiner{rightTable: rightTable, joinType: joinType})
	return incompleteBaseJoiner{BaseSelectBuilder: updated}
}

// Join generates intermediate IncompleteSelectJoin instance.
// Takes TableName to join and TableJoinType constant defining required join type to produce.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query BaseSelectBuilder) Join(rightTable TableIdent, joinType TableJoinType) IncompleteSelectJoin {
	return query.joinIdent(rightTable, joinType)
}

// InnerJoin generates intermediate IncompleteSelectJoin instance for InnerJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, InnerJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query BaseSelectBuilder) InnerJoin(rightTable TableIdent) IncompleteSelectJoin {
	return query.Join(rightTable, InnerJoin)
}

// LeftJoin generates intermediate IncompleteSelectJoin instance for LeftJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, LeftJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query BaseSelectBuilder) LeftJoin(rightTable TableIdent) IncompleteSelectJoin {
	return query.Join(rightTable, LeftJoin)
}

// RightJoin generates intermediate IncompleteSelectJoin instance for RightJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, RightJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query BaseSelectBuilder) RightJoin(rightTable TableIdent) IncompleteSelectJoin {
	return query.Join(rightTable, RightJoin)
}

// FullJoin generates intermediate IncompleteSelectJoin instance for FullJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, FullJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (query BaseSelectBuilder) FullJoin(rightTable TableIdent) IncompleteSelectJoin {
	return query.Join(rightTable, FullJoin)
}

// Count creates an CountBuilder using parameters of this SelectManyBuilder.
func (query BaseSelectBuilder) Count() CountBuilder {
	return Count(query)
}

// SelectFrom makes a new BaseSelectBuilder instance.
// Takes string, TableName or TableIdent to build BaseSelectBuilder over it.
func SelectFrom[T TableNameParameter](tableNameProvider T) BaseSelectBuilder {
	var (
		p     any = tableNameProvider
		table TableIdent
	)

	switch typedValue := p.(type) {
	case string:
		table = Table(TableName(typedValue))
	case TableName:
		table = Table(typedValue)
	case TableIdent:
		table = typedValue
	}

	return BaseSelectBuilder{
		baseTable: table,
		joins:     make([]TableJoiner, 0),
		where:     NewWhere(),
		fields:    NewFields(),
	}
}
