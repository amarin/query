package query_test

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

type UUID struct{ res string }

func (U UUID) Value() (driver.Value, error) {
	return U.res, nil
}

func TestFieldValue_Values(t *testing.T) {
	uuid1 := UUID{"dbc72284-f696-4a5d-98af-79568b0d5141"}
	uuid2 := UUID{"dbc72284-f696-4a5d-98af-79568b0d5142"}
	tests := []struct {
		name       string
		value      interface{}
		wantPanics bool
		wantResult []interface{}
	}{
		{"int", 42, false, []any{42}},
		{"int64", int64(42), false, []any{int64(42)}},
		{"int32", int32(42), false, []any{int32(42)}},
		{"int16", int16(42), false, []any{int16(42)}},
		{"int8", int8(42), false, []any{int8(42)}},
		{"uint", 42, false, []any{42}},
		{"uint64", uint64(42), false, []any{uint64(42)}},
		{"uint32", uint32(42), false, []any{uint32(42)}},
		{"uint16", uint16(42), false, []any{uint16(42)}},
		{"uint8", uint8(42), false, []any{uint8(42)}},
		{"float64", float64(3.14), false, []any{float64(3.14)}},
		{"float32", float32(2.71), false, []any{float32(2.71)}},
		{"false", false, false, []any{false}},
		{"true", true, false, []any{true}},
		{"str", "i_am", false, []any{"i_am"}},
		{"[]str", []string{"one", "two"}, false, []any{"one", "two"}},
		{"uuid", uuid1, false, []any{uuid1.res}},
		{"[]uuid", []UUID{uuid1, uuid2}, false, []any{uuid1.res, uuid2.res}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := query.Field("field").Of("table")
			fieldValue := field.Value(tt.value)
			var values []any
			takeValues := func() {
				values = fieldValue.Values()
			}
			if tt.wantPanics {
				require.Panics(t, takeValues)
				return
			} else {
				require.NotPanics(t, takeValues)
			}
			require.Equal(t, tt.wantResult, values)
		})
	}
}
