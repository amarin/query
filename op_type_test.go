package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOperation_String(t *testing.T) {
	tests := []struct {
		name string
		op   Operation
		want string
	}{
		{"select", DoSelect, "SELECT"},
		{"insert", DoInsert, "INSERT"},
		{"update", DoUpdate, "UPDATE"},
		{"delete", DoDelete, "DELETE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.op.String(), tt.want)
		})
	}
}
