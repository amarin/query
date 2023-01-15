package query_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/query"
)

func TestNewFieldSpec(t *testing.T) {
	spec := query.Field("test")
	require.EqualValues(t, spec.FieldName(), "test")
	require.EqualValues(t, spec.Alias(), "")
	require.EqualValues(t, spec.TableName(), "")
}

func TestFieldSpec_As(t *testing.T) {
	spec := query.Field("test")

	t.Run("initial", func(t *testing.T) {
		require.EqualValues(t, spec.FieldName(), "test")
		require.EqualValues(t, spec.Alias(), "")
		require.EqualValues(t, spec.TableName(), "")
	})

	t.Run("use_as", func(t *testing.T) {
		withAlias := spec.As("prod")
		require.EqualValues(t, withAlias.FieldName(), "test")
		require.EqualValues(t, withAlias.Alias(), "prod")
		require.EqualValues(t, withAlias.TableName(), "")
	})

	t.Run("initial_remains_the_same", func(t *testing.T) {
		require.EqualValues(t, spec.FieldName(), "test")
		require.EqualValues(t, spec.Alias(), "")
		require.EqualValues(t, spec.TableName(), "")
	})
}

func TestFieldSpec_Of(t *testing.T) {
	spec := query.Field("test")
	t.Run("initial", func(t *testing.T) {
		require.EqualValues(t, spec.FieldName(), "test")
		require.EqualValues(t, spec.Alias(), "")
		require.EqualValues(t, spec.TableName(), "")
	})

	t.Run("use_of", func(t *testing.T) {
		withAlias := spec.Of("someTable")
		require.EqualValues(t, withAlias.FieldName(), "test")
		require.EqualValues(t, spec.Alias(), "")
		require.EqualValues(t, withAlias.TableName(), "someTable")
	})

	t.Run("initial_remains_the_same", func(t *testing.T) {
		require.EqualValues(t, spec.FieldName(), "test")
		require.EqualValues(t, spec.Alias(), "")
		require.EqualValues(t, spec.TableName(), "")
	})
}

func TestFieldSpec_Render(t *testing.T) {
	tests := []struct {
		name string
		spec query.FieldDefinition
		want string
	}{
		{"only_name", query.Field("test"), "test"},
		{"table_dot_name", query.Field("test").Of("table"), "table.test"},
		{"name_as_alias", query.Field("test").As("Name"), "test AS Name"},
		{"table_name_as_alias", query.Field("test").Of("table").As("Name"), "table.test AS Name"},
		{"empty_table_returns_as without_table", query.Field("test").Of("").As("Name"), "test AS Name"},
		{"empty_alias_returns_as_without_alias", query.Field("test").As(""), "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spec.RenderSpec(); got != tt.want {
				t.Errorf("RenderSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSpec_RenderField(t *testing.T) {
	tests := []struct {
		name string
		spec query.FieldDefinition
		want string
	}{
		{"name", query.Field("test"), "test"},
		{"table_and_name", query.Field("test").Of("table"), "table.test"},
		{"alias_no_table", query.Field("test").As("prod"), "prod"},
		{"table_but_alias", query.Field("test").Of("table").As("prod"), "prod"},
		{"alias_with_table", query.Field("test").As("prod").Of("table"), "prod"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.spec.RenderField())
		})
	}
}

func TestField(t *testing.T) {
	tests := []struct {
		name      string
		args      query.FieldName
		wantErr   bool
		wantName  query.FieldName
		wantTable query.TableName
		wantAlias query.FieldName
	}{
		{"ok_simple_field_name", "f1", false, "f1", "", ""},
		{"ok_table_dot_name", "table.f1", false, "f1", "table", ""},
		{"ok_name_as_alias", "f1 as field1", false, "f1", "", "field1"},
		{"ok_table_dot_name_as_alias", "table.f1 as field1", false, "f1", "table", "field1"},
		{"nok_2_tokens", "table f1", true, "f1", "table", "field1"},
		{"nok_3_tokens_but_not_as", "life is life", true, "f1", "table", "field1"},
		{"nok_3_plus_tokens", "table f1 f2 f3", true, "f1", "table", "field1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := query.FieldOrError(tt.args)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.wantName, got.FieldName())
			require.Equal(t, tt.wantTable, got.TableName())
			require.Equal(t, tt.wantAlias, got.Alias())
		})
	}
}
