package condition

import (
	"fmt"

	"github.com/FrancoLiberali/cql/sql"
)

type FieldAggregation struct {
	field IField
}

// Count returns the number of values that are not null
func (fieldAggregation FieldAggregation) Count() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.Count,
	}
}

// Min returns the minimum value of all values
func (fieldAggregation FieldAggregation) Min() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.Min,
	}
}

// Max returns the maximum value of all values
func (fieldAggregation FieldAggregation) Max() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.Max,
	}
}

type NumericFieldAggregation struct {
	FieldAggregation
}

// Sum calculates the summation of all values
func (fieldAggregation NumericFieldAggregation) Sum() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.Sum,
	}
}

// Average calculates the average (arithmetic mean) of all values
func (fieldAggregation NumericFieldAggregation) Average() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.Average,
	}
}

// And calculates the bitwise AND of all non-null values (null values are ignored)
//
// Not available for: sqlite, sqlserver
func (fieldAggregation NumericFieldAggregation) And() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.BitAndAggregation,
	}
}

// Or calculates the bitwise OR of all non-null values (null values are ignored)
//
// Not available for: sqlite, sqlserver
func (fieldAggregation NumericFieldAggregation) Or() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.BitOrAggregation,
	}
}

type BoolFieldAggregation struct {
	FieldAggregation
}

// All returns true if all the values are true
func (fieldAggregation BoolFieldAggregation) All() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.All,
	}
}

// All returns true if at least one value is true
func (fieldAggregation BoolFieldAggregation) Any() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.Any,
	}
}

// None returns true if all values are false
func (fieldAggregation BoolFieldAggregation) None() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		Function: sql.None,
	}
}

type Aggregation struct {
	field        IField
	Function     sql.FunctionByDialector
	havingValues []any
}

func (aggregation Aggregation) Eq(value int) Aggregation {
	// TODO ver si siempre es int, depende del resultado que devuelva la agregacion
	aggregation.havingValues = []any{value}
	return aggregation
}

func (aggregation Aggregation) toSQL(query *GormQuery, table Table) (string, error) {
	function, isPresent := aggregation.Function.Get(query.Dialector())
	if !isPresent {
		return "", functionError(ErrUnsupportedByDatabase, aggregation.Function)
	}

	columnSQL := ""

	if aggregation.field != nil { // CountAll
		columnSQL = aggregation.field.columnSQL(query, table)
	}

	return function.ApplyTo(columnSQL, 0), nil
}

func (aggregation Aggregation) toSelectSQL(query *GormQuery, table Table, as string) (string, error) {
	functionSQL, err := aggregation.toSQL(query, table)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s AS %s",
		functionSQL,
		as,
	), nil
}

func (aggregation Aggregation) toHavingSQL(query *GormQuery, table Table) (string, []any, error) {
	functionSQL, err := aggregation.toSQL(query, table)
	if err != nil {
		return "", nil, err
	}

	// TODO hardcodeado el =
	return functionSQL + " = ?", aggregation.havingValues, nil
}
