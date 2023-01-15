package query

import (
	"fmt"
	"strings"
)

// FieldDefinition defines field definition data structure.
// It contains field name, field table name and field name alias.
// Used as field name wrapper in conditions and fields lists of queries.
type FieldDefinition struct {
	fieldName string // original field name, using string here to reduce type conversions in methods.
	tableName string // original field table name, do not required until building joins
	alias     string // retrieve data as this name
}

// FieldOrError creates new FieldDefinition.
// Takes string, FieldName of FieldDefinition instance.
// If FieldDefinition passed as argument it returns immediately with pointer to specified argument and nil error.
// If FieldName or string argument provided it parsed using rules below.
// 1) identification string with no dots and spaces defines field name with empty (unknown) table name and alias;
// 2) "<table>.<name>" form defines both field name and field table name;
// 3) "<name> as|AS <alias>" form defines both field name and its alias;
// 4) "<table>.<name> as|AS <alias>" defines all field name, table name and alias.
// Returns error if no format matched.
// Use Field constructor if argument fully determined and suits one of formats above.
// Note setting alias to same value as field name gives empty alias if table name is empty.
func FieldOrError[T FieldNameParameter](fieldName T) (field *FieldDefinition, err error) {
	var (
		param   any = fieldName
		strName string
	)

	switch value := param.(type) {
	case FieldDefinition:
		return &value, nil
	case FieldName:
		strName = string(value)
	case string:
		strName = value
	default:
		strName = fmt.Sprintf("%v", value)
	}

	tokens := strings.Fields(strName)
	switch {
	case len(tokens) == 1 && !strings.Contains(tokens[0], "."):
		return &FieldDefinition{fieldName: tokens[0], alias: "", tableName: ""}, nil
	case len(tokens) == 1 && strings.Contains(tokens[0], "."):
		nameTokens := strings.Split(tokens[0], ".")
		return &FieldDefinition{fieldName: nameTokens[1], alias: "", tableName: nameTokens[0]}, nil
	case len(tokens) == 3 && strings.ToUpper(tokens[1]) == "AS" && !strings.Contains(tokens[0], "."):
		return &FieldDefinition{fieldName: tokens[0], alias: tokens[2], tableName: ""}, nil
	case len(tokens) == 3 && strings.ToUpper(tokens[1]) == "AS" && strings.Contains(tokens[0], "."):
		nameTokens := strings.Split(tokens[0], ".")
		return &FieldDefinition{fieldName: nameTokens[1], alias: tokens[2], tableName: nameTokens[0]}, nil
	default:
		return nil, fmt.Errorf("%w: unexpected field specification: %v", Error, strName)
	}
}

func mustField[T FieldNameParameter](fieldName T) (field FieldDefinition) {
	var (
		err       error
		param     any = fieldName
		fieldPtr  *FieldDefinition
		takenName FieldName
	)

	switch value := param.(type) {
	case FieldDefinition:
		return value
	case FieldName:
		takenName = value
	case string:
		takenName = FieldName(value)
	default:
		takenName = FieldName(fmt.Sprintf("%v", value))
	}
	fieldPtr, err = FieldOrError(takenName)
	if err != nil {
		panic(err)
	}

	return *fieldPtr
}

// Field creates new mustField identification.
// It can use several forms to define not only mustField name but a table name, an alias or both:
// 1) identification string with no dots and spaces defines mustField name with empty (unknown) table name and alias;
// 2) "<table>.<name>" form defines both mustField name and mustField table name;
// 3) "<name> as <alias>" form defines both mustField name and its alias;
// 4) "<table>.<name> as <alias>" defines all mustField name, table name and alias.
// Panics if no format matched.
// Use FieldOrError constructor if fieldName format is not guaranteed to suit one of formats above.
func Field[T FieldNameParameter](spec T) (field FieldDefinition) {
	var (
		ok    bool
		param any = spec
	)
	if field, ok = param.(FieldDefinition); ok {
		return field
	}

	return mustField(spec)
}

// FieldName returns mustField name as wrapped FieldName string type.
func (fieldIdent FieldDefinition) FieldName() FieldName {
	return FieldName(fieldIdent.fieldName)
}

// Alias returns mustField name alias. Same as FieldName until As not called.
func (fieldIdent FieldDefinition) Alias() FieldName {
	return FieldName(fieldIdent.alias)
}

// TableName returns table name of mustField pointed by FieldDefinition.
func (fieldIdent FieldDefinition) TableName() TableName {
	return TableName(fieldIdent.tableName)
}

// As returns a copy of FieldDefinition with mustField name alias set to specified fieldName.
// Note setting alias to same name as field name renders FieldDefinition without alias unless table name is not empty.
func (fieldIdent FieldDefinition) As(alias FieldName) FieldDefinition {
	fieldIdent.alias = string(alias)
	return fieldIdent
}

// Of returns a copy of FieldDefinition with table name set to specified baseTable.
func (fieldIdent FieldDefinition) Of(tableName TableName) FieldDefinition {
	fieldIdent.tableName = string(tableName)
	return fieldIdent
}

// RenderSpec returns a field specification to use in SQL queries as a fetch items enumeration.
func (fieldIdent FieldDefinition) RenderSpec() string {
	items := make([]string, 0, 4)
	if len(fieldIdent.tableName) > 0 {
		items = append(items, fieldIdent.tableName+".")
	}
	items = append(items, fieldIdent.fieldName)
	switch {
	case len(fieldIdent.alias) > 0 && fieldIdent.alias != fieldIdent.fieldName:
		// alias differs field name, render alias
		items = append(items, " "+kwAs.String()+" ", fieldIdent.alias)
	case len(fieldIdent.alias) > 0 && fieldIdent.alias == fieldIdent.fieldName && len(fieldIdent.tableName) > 0:
		// alias matches field name but table is not empty, render alias too
		items = append(items, " "+kwAs.String()+" ", fieldIdent.alias)
	}

	return strings.Join(items, "")
}

// RenderField returns a mustField identification to use in SQL queries in conditional or sorting clauses.
func (fieldIdent FieldDefinition) RenderField() string {
	items := make([]string, 0, 4)
	if len(fieldIdent.tableName) > 0 {
		items = append(items, fieldIdent.tableName+".")
	}
	items = append(items, fieldIdent.fieldName)
	switch {
	case len(fieldIdent.alias) > 0 && fieldIdent.alias != fieldIdent.fieldName:
		// use alias if defined and differs from original mustField name
		return fieldIdent.alias
	case len(fieldIdent.tableName) > 0:
		// use table.mustField form as mustField name is not empty
		return fieldIdent.tableName + "." + fieldIdent.fieldName
	default:
		return fieldIdent.fieldName
	}
}

// RenderTableSpec returns a field name in form <table_name>.<field_name> as string value.
// Used to render SQL JOIN conditional clauses.
// If value of FieldName type required use TableFieldName instead.
func (fieldIdent FieldDefinition) RenderTableSpec() string {
	return fieldIdent.tableName + "." + fieldIdent.fieldName
}

// TableFieldName returns a field name in form <table_name>.<field_name> as FieldName value.
// If value of string type required use RenderTableSpec instead.
func (fieldIdent FieldDefinition) TableFieldName() FieldName {
	return FieldName(fieldIdent.RenderTableSpec())
}

// Value generates FieldValue item using current field definition.
// Note FieldValue keeps its own copy of field definition taken in constructor,
// so any further changes to FieldDefinition will not affect generated FieldValue.
// To sync changes made in FieldDefinition after FieldValue generated, use FieldValue.ApplyFieldSpec.
func (fieldIdent FieldDefinition) Value(value any) FieldValue {
	return NewFieldValue(fieldIdent.FieldName(), value).ApplyFieldTable(fieldIdent.TableName())
}
