package query

import (
	"fmt"
	"strings"
)

// TableIdent defines data structure to represent table in SQL queries generated in query package.
// It encapsulates table name as well as table alias when it required to use.
// Alias is used only in case it is not empty and stored value is different from stored table name.
// Difference requirement leads to empty values set to alias in constructors or As result if alias equals to name.
// Besides TableName provide and As modifier it also provides a convenient method to generate FieldDefinition items.
// Used both as self-contained item in single-table queries generators and as a part o TableJoiner instances
// when SQL JOIN queries are required.
type TableIdent struct {
	name  string // original table name as known in database, use string here to avoid multiple type conversions
	alias string // table name alias to use in conditions
}

// Constructors are TableOrError and Table.
// See also TableName.Ident method to easily translate TableName`s instances into TableIdent`s.

// TableOrError creates new table identification.
// It can use 2 forms to define not table name or table name with alias:
// 1) identification string with no spaces defines table name with no alias;
// 2) "<name> as|AS <alias>" form defines both table name and its alias;
// 3) "<name> <alias>" form defines both table name and its alias;
// Returns error if no format matched.
// Use Table constructor if table name fully determined and suits one of formats above.
// Note setting alias to same value as table name gives empty alias.
func TableOrError[T TableNameParameter](name T) (*TableIdent, error) {
	var (
		strTableName string
		tokens       []string
		p            any = name
	)

	switch typed := p.(type) {
	case string:
		strTableName = typed
	case TableName:
		strTableName = string(typed)
	case TableIdent:
		return &typed, nil // do not parse, already TableIdent
	}

	tokens = strings.Fields(strTableName)
	switch {
	case len(tokens) == 1:
		return &TableIdent{name: tokens[0], alias: ""}, nil
	case len(tokens) == 3 && strings.ToUpper(tokens[1]) == "AS" && len(tokens[2]) > 0 && tokens[2] != tokens[0]:
		return &TableIdent{name: tokens[0], alias: tokens[2]}, nil
	case len(tokens) == 3 && strings.ToUpper(tokens[1]) == "AS" && len(tokens[2]) > 0 && tokens[2] == tokens[0]:
		return &TableIdent{name: tokens[0], alias: ""}, nil
	case len(tokens) == 2 && len(tokens[1]) > 0 && tokens[1] != tokens[0]:
		return &TableIdent{name: tokens[0], alias: tokens[1]}, nil
	default:
		return nil, fmt.Errorf("%w: unexpected table ident: `%v`", Error, name)
	}
}

// Table creates new table identification.
// It can use 2 forms to define not table name or table name with alias:
// 1) identification string with no spaces defines table name with no alias;
// 2) "<name> as|AS <alias>" form defines both table name and its alias;
// 3) "<name> <alias>" form defines both table name and its alias;
// Panics if no format matched.
// Use TableOrError constructor if table name format is not guaranteed to suit one of formats above.
func Table[T TableNameParameter](name T) TableIdent {
	identPtr, err := TableOrError(name)
	if err != nil {
		panic(err)
	}

	return *identPtr
}

// InsertInto generates InsertBuilder instance.
// Optional field values to insert could be set right here.
func (tableIdent TableIdent) InsertInto(insertValues ...FieldValue) InsertBuilder {
	return InsertInto(tableIdent).Values(insertValues...)
}

// Update generates table updater instance.
// Use UpdateBuilder.Where to finalize UpdateBuilder instance configuration before use.
func (tableIdent TableIdent) Update(updateValues ...FieldValue) UpdateBuilder {
	return Update(tableIdent).Set(updateValues...)
}

// Delete generates table updater instance.
// Use DeleteBuilder.Where to finalize DeleteBuilder instance configuration before use.
func (tableIdent TableIdent) Delete(filterConditions ...Condition) DeleteBuilder {
	return Delete(tableIdent).Where(filterConditions...)
}

// Attribute getters are TableName and Alias.

// TableName returns a table name as defined in database.
func (tableIdent TableIdent) TableName() TableName {
	return TableName(tableIdent.name)
}

// Alias returns a table name alias to use in query. Can be empty string.
func (tableIdent TableIdent) Alias() TableName {
	return TableName(tableIdent.alias)
}

// RenderFrom render table identification string to fill SQL FROM clause.
// The result is "<name>" or "<name> AS <alias>".
// Implements ClauseFromRenderer.
func (tableIdent TableIdent) RenderFrom() string {
	switch {
	case len(tableIdent.alias) > 0 && tableIdent.alias != tableIdent.name:
		return tableIdent.name + " AS " + tableIdent.alias
	default:
		return tableIdent.name
	}
}

// As returns a copy of TableIdent having Alias set to specified value, leaving original TableIdent settings intact.
// Note setting alias same as table name will set empty alias, see why in TableIdent structure definition description.
func (tableIdent TableIdent) As(alias TableName) (updated TableIdent) {
	updated = tableIdent
	if string(alias) == updated.name {
		alias = ""
	}
	updated.alias = string(alias)

	return updated
}

// Field generates new FieldDefinition having table name filled&
// If TableIdent alias is not empty it used to fill FieldDefinition table name, otherwise TableIdent.TableName will be used.
// Note TableIdent alias change with As() will not propagate to fields generated before.
func (tableIdent TableIdent) Field(name FieldName) FieldDefinition {
	return Field(name).Of(TableName(tableIdent.aliasOrName()))
}

// aliasOrName returns tableIdent alias if not empty or name if alias is not set.
func (tableIdent TableIdent) aliasOrName() string {
	if len(tableIdent.alias) > 0 {
		return tableIdent.alias
	}

	return tableIdent.name
}

// Select generates BaseSelectBuilder.
// Specified Condition's will be used to filter values.
// Shorthand to SelectFrom((TableIdent)).Where(filterConditions...)
func (tableIdent TableIdent) Select(filterConditions ...Condition) BaseSelectBuilder {
	return SelectFrom(tableIdent).Where(filterConditions...)
}

// SelectOne generates table query helper to fetch single row.
// Specified Condition's will be used to filter values.
// Shorthand to SelectSingleFrom((TableIdent)).Where(filterConditions...)
func (tableIdent TableIdent) SelectOne(filterConditions ...Condition) SelectSingleBuilder {
	return SelectSingleFrom(tableIdent).Where(filterConditions...)
}

// SelectMany generates table query helper to fetch many values.
// Specified Condition's will be used to filter values.
// Shorthand to SelectManyFrom((TableIdent)).Where(filterConditions...)
func (tableIdent TableIdent) SelectMany(filterConditions ...Condition) SelectManyBuilder {
	return SelectManyFrom(tableIdent).Where(filterConditions...)
}

// Count generates CountBuilder builder.
// Specified Condition's will be used to filter values.
func (tableIdent TableIdent) Count(filterConditions ...Condition) CountBuilder {
	return SelectFrom(tableIdent).Where(filterConditions...).Count()
}

// Join generates intermediate IncompleteSelectJoin instance.
// Takes TableName to join and TableJoinType constant defining required join type to produce.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Call IncompleteSelectJoin.On will return BaseSelectBuilder with join builder data finished.
func (tableIdent TableIdent) Join(rightTable TableIdent, joinType TableJoinType) IncompleteSelectJoin {
	return tableIdent.Select().Join(rightTable, joinType)
}

// InnerJoin generates intermediate IncompleteSelectJoin instance for InnerJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, InnerJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (tableIdent TableIdent) InnerJoin(rightTable TableIdent) IncompleteSelectJoin {
	return tableIdent.Join(rightTable, InnerJoin)
}

// LeftJoin generates intermediate IncompleteSelectJoin instance for LeftJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, LeftJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (tableIdent TableIdent) LeftJoin(rightTable TableIdent) IncompleteSelectJoin {
	return tableIdent.Join(rightTable, LeftJoin)
}

// RightJoin generates intermediate IncompleteSelectJoin instance for RightJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, RightJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (tableIdent TableIdent) RightJoin(rightTable TableIdent) IncompleteSelectJoin {
	return tableIdent.Join(rightTable, RightJoin)
}

// FullJoin generates intermediate IncompleteSelectJoin instance for FullJoin type.
// Takes right table name or ident to join.
// Panics if rightTable is not instance of string, TableName or TableIdent types.
// Shortcut to Join(rightTable TableName, FullJoin).
// Call IncompleteSelectJoin.On will return updated BaseSelectBuilder with join builder data finished.
func (tableIdent TableIdent) FullJoin(rightTable TableIdent) IncompleteSelectJoin {
	return tableIdent.Join(rightTable, FullJoin)
}
