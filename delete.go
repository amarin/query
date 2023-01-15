package query

import (
	"fmt"
	"strings"
)

// DeleteBuilder helps to build SQL DELETE queries.
type DeleteBuilder struct {
	BaseBuilder
	tableName TableName
	where     Group
}

// RenderFrom returns string representation of table name or tables join with possible tables aliases.
func (updater DeleteBuilder) RenderFrom() string {
	return updater.tableName.RenderFrom()
}

// Where adds fields conditions GroupAND returns modified DeleteBuilder.
func (updater DeleteBuilder) Where(fieldConditions ...Condition) DeleteBuilder {
	updater.where = updater.where.GroupAND(fieldConditions...)
	return updater
}

// BuildQueryAndParams returns query string and params to fill in SQL DELETE query string.
// If query build failed returns non-nil error.
func (updater DeleteBuilder) BuildQueryAndParams() (sql string, params []interface{}, err error) {
	// disallow some cases
	switch {
	case len(updater.where.conditions) == 0:
		return "", params, fmt.Errorf("%w: empty conditions, will not delete every record", Error)
	case len(updater.tableName) == 0:
		return "", params, fmt.Errorf("%w: no table name set", Error)
	}

	sql = strings.Join([]string{
		kwDelete.String(),
		kwFrom.String(),
		updater.tableName.String(),
		kwWhere.String(),
		updater.where.Render(len(params)),
	}, " ")

	params = append(params, updater.where.Values()...)

	return sql, params, nil
}

// TableName returns table name to delete.
func (updater DeleteBuilder) TableName() TableName {
	return updater.tableName
}

// Delete generates table deleter using predefined table helper.
// Takes table name (of string, TableName or TableIdent types).
// Use DeleteBuilder.Where and DeleteBuilder.Set to finalize DeleteBuilder configuration.
func Delete[T TableNameParameter](tableNameProvider T) DeleteBuilder {
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

	return DeleteBuilder{
		BaseBuilder: BaseBuilder{op: DoDelete},
		tableName:   tableName,
		where:       NewGroup(LogicalAND),
	}
}
