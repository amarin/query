package query

import (
	"context"
)

// TableNameParameter allows use a string, TableName or TableIdent in cases where TableName or TableIdent required.
type TableNameParameter interface {
	string | TableName | TableIdent
}

// FieldNameParameter allows use a string, FieldName or FieldDefinition in cases where FieldName or FieldDefinition required.
type FieldNameParameter interface {
	string | FieldName | FieldDefinition
}

// IncompleteSelectJoin items is generated in BaseSelectBuilder join generation functions
// such as SelectFromTableBuilder.Join, BaseSelectBuilder.InnerJoin, SelectFromTableBuilder.LeftJoin,
// BaseSelectBuilder.RightJoin or BaseSelectBuilder.FullJoin.
// It provides single method On required to finalize JOIN clause builder and return to BaseSelectBuilder instance.
type IncompleteSelectJoin interface {
	// On finalises SQL JOIN definition adding left and right tables fields onto JOIN ... ON clause.
	// Note both left and right FieldDefinition's required to contain table name.
	// Returns updated BaseSelectBuilder instance.
	On(leftField FieldDefinition, rightField FieldDefinition) (updated BaseSelectBuilder)
}

// IncompleteSelectManyJoiner items is generated in SelectManyBuilder join generation functions
// such as SelectFromTableBuilder.Join, SelectManyBuilder.InnerJoin, SelectManyBuilder.LeftJoin,
// SelectManyBuilder.RightJoin or SelectManyBuilder.FullJoin.
// It provides single method On required to finalize JOIN clause builder and return to SelectManyBuilder instance.
// See also IncompleteSelectJoin interface for serve joins on BaseSelectBuilder.
type IncompleteSelectManyJoiner interface {
	// On finalises SQL JOIN definition adding left and right tables fields onto JOIN ... ON clause.
	// Note both left and right FieldDefinition's required to contain table name.
	// Returns updated SelectManyBuilder instance.
	On(leftField FieldDefinition, rightField FieldDefinition) (updated SelectManyBuilder)
}

// IncompleteSelectSingleJoiner items is generated in SelectSingleBuilder join generation functions
// such as SelectSingleBuilder.Join, SelectSingleBuilder.InnerJoin, SelectSingleBuilder.LeftJoin,
// SelectSingleBuilder.RightJoin or SelectSingleBuilder.FullJoin.
// It provides single method By required to finalize JOIN clause builder and return to SelectSingleBuilder instance.
// See also IncompleteSelectJoin interface for serve joins on BaseSelectBuilder.
type IncompleteSelectSingleJoiner interface {
	// On finalises SQL JOIN definition adding left and right tables fields onto JOIN ... ON clause.
	// Note both left and right FieldDefinition's required to contain table name.
	// Returns updated SelectSingleBuilder instance.
	On(leftField FieldDefinition, rightField FieldDefinition) (updated SelectSingleBuilder)
}

// Fetcher requires implementation could fetch records from underline database using prepared query SelectManyBuilder.
type Fetcher interface {
	FetchCtx(ctx context.Context, queryParams SelectManyBuilder, target interface{}) (err error)
}

// Getter requires implementation could get single record from underline database using prepared query SelectManyBuilder.
type Getter interface {
	GetCtx(ctx context.Context, queryParams SelectSingleBuilder, target interface{}) (err error)
}

// Updater requires implementation provides single record update method using prepared UpdateBuilder.
type Updater interface {
	UpdateOneCtx(ctx context.Context, updateParams UpdateBuilder) (err error)
}

// Deleter requires implementation provides records delete method using prepared DeleteBuilder.
type Deleter interface {
	DeleteManyCtx(ctx context.Context, updateParams DeleteBuilder) (affectedRows int, err error)
}

// FetchCounter requires implementation could fetch records with total records count using prepared query SelectManyBuilder.
type FetchCounter interface {
	FetchCountCtx(ctx context.Context, queryParams SelectManyBuilder, target interface{}) (totalRows int, err error)
}

// Counter requires implementation could fetch total records count using prepared query CountBuilder.
type Counter interface {
	CountCtx(ctx context.Context, helper CountBuilder) (rowsCount int, err error)
}

// FetcherCounter provides both FetchCtx and CountCtx methods.
type FetcherCounter interface {
	Fetcher
	Counter
}

// CountingClauseRenderer requires implementations could render its SQL clause part using numbered parameters.
type CountingClauseRenderer interface {
	// Render renders SQL clause or its part with respect of parameters count added previously.
	// Takes existed parameters count (0 means no parameters are defined yet).
	// Implementation should render parameters substitutions starting from number = paramNum+1.
	// It expected substitutions are rendered using "$<number>" notation, as expected by postgresql driver.
	// Values result should contain the same values count as substitutions used in clause.
	// If implementation is condition and condition is negated, implementation SHOULD write negate prefix itself.
	Render(parametersCount int) (sql string)
}

// ValuesProvider requires clause renderer implementations should provide substitution values slice.
type ValuesProvider interface {
	// Values returns a set of parameters to substitute when SQL query is fully constructed and passed to execution.
	// Parameters count should be the same as substitutions used in Render result.
	// If implementation requires nothing to substitute it should return empty slice.
	Values() []any
}

// RawClauseRenderer requires implementations could render its SQL clause part using default parameters substitution.
type RawClauseRenderer interface {
	// RenderSQL renders SQL clause or its part.
	// Implementation should render parameters substitutions using standard sql "?"(question) character.
	// Values result should contain the same values count as substitutions used in clause.
	// If implementation is condition and condition is negated, implementation SHOULD write negate prefix itself.
	RenderSQL() (sql string)
}

// Inserter interface defines field condition methods.
// Inserter requires implementation provides single record insert method using prepared InsertBuilder.
type Inserter interface {
	InsertOneCtx(ctx context.Context, updateParams InsertBuilder) (err error)
}

type Condition interface {
	ValuesProvider
	CountingClauseRenderer
	RawClauseRenderer

	// JoinType defines JoinType to combine condition with previous.
	JoinType() JoinType

	// RenderJoin renders join logical operation GroupAND negate suffix to place before condition.
	RenderJoin(isFirst bool) (result string)

	// IsNegate returns true if condition negated GroupAND false otherwise.
	IsNegate() bool

	// Join returns a copy of Condition having JoinType set to specified value.
	Join(newJoinType JoinType) Condition

	// Negate returns a copy of Condition having IsNegate set to specified value.
	Negate(newNegate bool) Condition

	// RenderNegate renders negate prefix if condition is market as negated.
	// It should return single "NOT" with no spaces around if negated GroupAND empty string otherwise.
	RenderNegate() string

	// And generates new condition which true on all conditions met.
	And(conditions ...Condition) Condition

	// Or generates new condition group which true on either initial condition is true GroupOR all of additional are true.
	Or(conditions ...Condition) Condition

	// FieldName returns field name of condition applied to.
	FieldName() FieldName

	// ApplyFieldSpec makes a copy of Condition with updated FieldDefinition if fieldName match.
	// If FieldName is not matched it does nothing.
	ApplyFieldSpec(spec FieldDefinition) Condition

	// ApplyFieldTable makes a copy of Condition with updated FieldDefinition table name.
	ApplyFieldTable(table TableName) Condition
}

// TableNameProvider requires instances could provide used table name.
type TableNameProvider interface {
	TableName() TableName
}
