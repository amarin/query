package query

import (
	"fmt"
	"strconv"
	"strings"
)

// InsertBuilder helps to build SQL INSERT queries.
type InsertBuilder struct {
	BaseBuilder
	tableName TableName
	setValues []FieldValue
}

// RenderFrom returns string representation of table name or tables join with possible tables aliases.
func (inserter InsertBuilder) RenderFrom() string {
	return inserter.tableName.RenderFrom()
}

// SetValues returns a copy of fields values attached to use in create query.
func (inserter InsertBuilder) SetValues() []FieldValue {
	return inserter.setValues
}

// Values generates new InsertBuilder having additional field values to build insert query parameters.
// Can be called multiple times to set as many field values as required.
func (inserter InsertBuilder) Values(values ...FieldValue) InsertBuilder {
	inserter.setValues = append(inserter.setValues, values...)
	return inserter
}

// BuildQueryAndParams generates SQL INSERT query based on the set values.
// Returns SQL INSERT query string, parameters to fill placeholders in driver.
// If any errors occurs returns that error.
func (inserter InsertBuilder) BuildQueryAndParams() (sql string, params []interface{}, err error) {
	params = make([]interface{}, 0)

	// disallow some cases
	switch {
	case len(inserter.setValues) == 0:
		return "", params, fmt.Errorf("%w: no fields to insert", Error)
	case len(inserter.tableName) == 0:
		return "", params, fmt.Errorf("%w: table name empty", Error)
	}

	columns := make([]string, 0)
	valuePlaceholders := make([]string, 0)

	for _, fieldValue := range inserter.setValues {
		params = append(params, fieldValue.Values()...)

		columns = append(columns, fieldValue.fieldName)
		valuePlaceholders = append(valuePlaceholders, "$"+strconv.Itoa(len(params)))
	}

	sql = strings.Join([]string{
		kwInsert.String(), kwInto.String(),
		inserter.tableName.String() + "(" + strings.Join(columns, ", ") + ")",
		kwValues.String(),
		"(" + strings.Join(valuePlaceholders, ", ") + ")",
	}, " ")

	return sql, params, nil
}

// TableName returns table name to insert data into.
func (inserter InsertBuilder) TableName() TableName {
	return inserter.tableName
}

// InsertInto generates InsertBuilder.
// Takes string table name, TableName or TableIdent.
// Use InsertBuilder.Values to finalize InsertBuilder configuration.
func InsertInto[T TableNameParameter](tableNameProvider T) InsertBuilder {
	var (
		p         any = tableNameProvider
		tableName TableName
	)

	switch typedValue := p.(type) {
	case string:
		tableName = TableName(typedValue)
	case TableName:
		tableName = typedValue
	case TableIdent:
		tableName = typedValue.TableName()
	}

	return InsertBuilder{
		BaseBuilder: BaseBuilder{op: DoInsert},
		tableName:   tableName,
		setValues:   make([]FieldValue, 0),
	}
}
