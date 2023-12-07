package query

import (
	"errors"
	"fmt"
)

// ValuesMap allows to define FieldValue set in mapping style.
type ValuesMap map[string]any

// FieldValues represents ValuesMap as FieldValues.
// If table name points to valid field name, it will be used as FieldValue field definition table.
// Returns non-nil error if any key is not valid field name.
// Note table spec is not applied to elements, its up to caller to adopt values when required.
func (valuesMap ValuesMap) FieldValues() ([]FieldValue, error) {
	res := make([]FieldValue, len(valuesMap))
	idx := 0
	for k, v := range valuesMap {
		fieldName := FieldName(k)
		if err := fieldName.Validate(); err != nil {
			return nil, errors.Join(fmt.Errorf("%w: invalid field name `%v`", Error, k), err)
		}
		res[idx] = *NewFieldValue(fieldName, v)

		idx++
	}

	return res, nil
}
