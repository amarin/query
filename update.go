package query

import (
	"fmt"
	"strconv"
	"strings"
)

// UpdateBuilder helps to build SQL UPDATE queries.
type UpdateBuilder struct {
	BaseBuilder
	tableName TableName
	where     Group
	setValues []FieldValue
}

// RenderFrom returns string representation of table name or tables join with possible tables aliases.
func (updater UpdateBuilder) RenderFrom() string {
	return updater.tableName.RenderFrom()
}

// SetValues returns a copy of fields to update.
func (updater UpdateBuilder) SetValues() []FieldValue {
	return updater.setValues
}

// Set generates new UpdateBuilder having extended fields to update array.
// Can be called multiple times to set as many field values as required.
func (updater UpdateBuilder) Set(values ...FieldValue) UpdateBuilder {
	updater.setValues = append(updater.setValues, values...)
	return updater
}

// Where adds fields conditions GroupAND returns modified UpdateBuilder.
func (updater UpdateBuilder) Where(fieldConditions ...Condition) UpdateBuilder {
	updater.where = updater.where.GroupAND(fieldConditions...)
	return updater
}

func (updater UpdateBuilder) fieldsAndValues() (sql string, params []any) {
	kwPairs := make([]string, len(updater.setValues))
	params = make([]any, len(updater.setValues))

	for idx, fieldValue := range updater.setValues {
		params[idx] = fieldValue.Values()[0]
		kwPairs[idx] = fieldValue.fieldName + "=$" + strconv.Itoa(idx+1)
	}

	return strings.Join(kwPairs, ", "), params
}

// BuildQueryAndParams returns query string and params to fill in SQL UPDATE query string.
// If query build failed returns non-nil error.
func (updater UpdateBuilder) BuildQueryAndParams() (sql string, params []interface{}, err error) {
	var fieldsEnum string

	// disallow some cases
	switch {
	case len(updater.where.conditions) == 0:
		return "", params, fmt.Errorf("%w: empty conditions, will not update every record", Error)
	case len(updater.setValues) == 0:
		return "", params, fmt.Errorf("%w: no fields to update set", Error)
	case len(updater.tableName) == 0:
		return "", params, fmt.Errorf("%w: no table name set", Error)
	}

	fieldsEnum, params = updater.fieldsAndValues()
	sql = strings.Join([]string{
		kwUpdate.String(),
		updater.tableName.String(),
		kwSet.String(),
		fieldsEnum,
		kwWhere.String(),
		updater.where.Render(len(params)),
	}, " ")

	params = append(params, updater.where.Values()...)

	return sql, params, nil
}

// TableName returns table name to update.
func (updater UpdateBuilder) TableName() TableName {
	return updater.tableName
}

// Update generates table updater using predefined table helper.
// Takes table name (of string, TableName or TableIdent types).
// Use UpdateBuilder.Where and UpdateBuilder.Set to finalize UpdateBuilder configuration.
func Update[T TableNameParameter](tableNameProvider T) UpdateBuilder {
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

	return UpdateBuilder{
		BaseBuilder: BaseBuilder{op: DoUpdate},
		tableName:   tableName,
		where:       NewGroup(LogicalAND),
		setValues:   make([]FieldValue, 0),
	}
}
