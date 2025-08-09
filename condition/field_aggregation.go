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
	field    IField
	Function sql.FunctionByDialector
}

func (aggregation Aggregation) toSQL(query *GormQuery, table Table, as string) (string, error) {
	function, isPresent := aggregation.Function.Get(query.Dialector())
	if !isPresent {
		return "", functionError(ErrUnsupportedByDatabase, aggregation.Function)
	}

	columnSQL := ""

	if aggregation.field != nil { // CountAll
		columnSQL = aggregation.field.columnSQL(query, table)
	}

	return fmt.Sprintf(
		"%s AS %s",
		function.ApplyTo(columnSQL, 0),
		as,
	), nil
}
