package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestLimiter_Render(t *testing.T) {
	tests := []struct {
		name     string
		fields   query.Limiter
		paramNum int
		wantSql  string
	}{
		{"empty", query.Limiter{}, 0, ""},
		{"empty_shifted", query.Limiter{}, 1, ""},
		{"limit", new(query.Limiter).Limit(1), 0, "LIMIT $1"},
		{"limit_shifted", new(query.Limiter).Limit(1), 1, "LIMIT $2"},
		{"offset", new(query.Limiter).Offset(1), 0, "OFFSET $1"},
		{"offset_shifted", new(query.Limiter).Offset(1), 1, "OFFSET $2"},
		{"limit_and_offset", new(query.Limiter).Offset(1).Limit(1), 0,
			"OFFSET $1 LIMIT $2"},
		{"limit_and_offset_shifted", new(query.Limiter).Offset(1).Limit(1), 1,
			"OFFSET $2 LIMIT $3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantSql, tt.fields.Render(tt.paramNum))
		})
	}
}

func TestLimiter_RenderSQL(t *testing.T) {
	tests := []struct {
		name    string
		fields  query.Limiter
		wantSql string
	}{
		{"empty", query.Limiter{}, ""},
		{"limit", new(query.Limiter).Limit(1), "LIMIT ?"},
		{"offset", new(query.Limiter).Offset(1), "OFFSET ?"},
		{"limit_and_offset", new(query.Limiter).Offset(1).Limit(1), "OFFSET ? LIMIT ?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantSql, tt.fields.RenderSQL())
		})
	}
}

func TestLimiter_Values(t *testing.T) {
	tests := []struct {
		name       string
		fields     query.Limiter
		wantValues []any
	}{
		{"empty", query.Limiter{}, []any{}},
		{"limit", new(query.Limiter).Limit(1), []any{uint(1)}},
		{"offset", new(query.Limiter).Offset(2), []any{uint(2)}},
		{"limit_and_offset", new(query.Limiter).Offset(3).Limit(4), []any{uint(3), uint(4)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantValues, tt.fields.Values())
		})
	}
}
