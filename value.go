package query

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// FieldValue joins field name with its value. Used both for conditions build and for update tables clause generation.
type FieldValue struct {
	FieldDefinition
	value interface{} // untyped field value. Used to provide value into database.
}

// ApplyFieldTable makes a copy of FieldValue with updated FieldDefinition table name.
// Implements Condition.
func (fieldValue FieldValue) ApplyFieldTable(table TableName) FieldValue {
	fieldValue.FieldDefinition.tableName = string(table)
	return fieldValue
}

// ApplyFieldSpec makes a copy of FieldValue with updated FieldDefinition on when mustField name matched.
// If field name is not matched it simply returns a copy of original FieldValue.
func (fieldValue FieldValue) ApplyFieldSpec(spec FieldDefinition) FieldValue {
	if fieldValue.FieldDefinition.fieldName == spec.fieldName {
		fieldValue.FieldDefinition = spec
	}

	return fieldValue
}

// Values returns a slice of interface{} containing single field value.
func (fieldValue FieldValue) Values() (result []interface{}) {
	var err error

	switch typedValue := fieldValue.value.(type) {
	case driver.Valuer:
		singleValue, _ := typedValue.Value()
		return []any{singleValue}
	case []driver.Valuer:
		result = make([]any, len(typedValue))

		for idx, value := range typedValue {
			singleValue, _ := value.Value()
			result[idx] = singleValue
		}

		return result
	case ValueProvider:
		return []any{typedValue.DatabaseValue()}
	case []ValueProvider:
		result = make([]any, len(typedValue))

		for idx, value := range typedValue {
			result[idx] = value.DatabaseValue()
		}

		return result
	case []string:
		result = make([]any, len(typedValue))

		for idx, value := range typedValue {
			result[idx] = value
		}

		return result
	}
	// no known types is found, check if array of elements implemented driver.Valuer received,
	// and translate using driver.Valuer when possible.
	valueReflection := reflect.ValueOf(fieldValue.value)

	if valueReflection.Kind() == reflect.Slice {
		valuerType := reflect.TypeOf((*driver.Valuer)(nil)).Elem()
		sliceType := valueReflection.Type()
		elementsReflection := sliceType.Elem()
		if elementsReflection.Implements(valuerType) {
			res := make([]any, valueReflection.Len())
			for idx := 0; idx < len(res); idx++ {
				elem := valueReflection.Index(idx)
				rawElem := elem.Interface()
				valuer := rawElem.(driver.Valuer)
				if valuer == nil {
					res[idx] = nil
				} else {
					if res[idx], err = valuer.Value(); err != nil {
						panic(fmt.Errorf("%w: translate value %v(%T): %v", Error, rawElem, rawElem, err))
					}
				}
			}

			return res
		}
	}
	return []any{fieldValue.value}
}

// NewFieldValue creates new FieldValue having specified field name and value.
func NewFieldValue(fieldName FieldName, fieldValue interface{}) *FieldValue {
	return &FieldValue{
		FieldDefinition: Field(fieldName),
		value:           fieldValue,
	}
}
