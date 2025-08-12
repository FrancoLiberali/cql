package condition

import (
	"fmt"

	"github.com/FrancoLiberali/cql/sql"
)

type FieldAggregation struct {
	field IField
}

// Count returns the number of values that are not null
func (fieldAggregation FieldAggregation) Count() NumericResultAggregation {
	return NumericResultAggregation{
		commonAggregation: commonAggregation{
			field:    fieldAggregation.field,
			Function: sql.Count,
		},
	}
}

// // Min returns the minimum value of all values
// func (fieldAggregation FieldAggregation) Min() AnyResultAggregation {
// 	// TODO ver que hacer aca
// 	return AnyResultAggregation{
// 		commonAggregation: commonAggregation{
// 			field:    fieldAggregation.field,
// 			Function: sql.Min,
// 		},
// 	}
// }

// // Max returns the maximum value of all values
// func (fieldAggregation FieldAggregation) Max() AnyResultAggregation {
// 	return AnyResultAggregation{
// 		commonAggregation: commonAggregation{
// 			field:    fieldAggregation.field,
// 			Function: sql.Max,
// 		},
// 	}
// }

type NumericFieldAggregation struct {
	FieldAggregation
}

// Sum calculates the summation of all values
func (fieldAggregation NumericFieldAggregation) Sum() NumericResultAggregation {
	return NumericResultAggregation{
		commonAggregation: commonAggregation{
			field:    fieldAggregation.field,
			Function: sql.Sum,
		},
	}
}

// Average calculates the average (arithmetic mean) of all values
func (fieldAggregation NumericFieldAggregation) Average() NumericResultAggregation {
	return NumericResultAggregation{
		commonAggregation: commonAggregation{
			field:    fieldAggregation.field,
			Function: sql.Average,
		},
	}
}

// And calculates the bitwise AND of all non-null values (null values are ignored)
//
// Not available for: sqlite, sqlserver
func (fieldAggregation NumericFieldAggregation) And() NumericResultAggregation {
	return NumericResultAggregation{
		commonAggregation: commonAggregation{
			field:    fieldAggregation.field,
			Function: sql.BitAndAggregation,
		},
	}
}

// Or calculates the bitwise OR of all non-null values (null values are ignored)
//
// Not available for: sqlite, sqlserver
func (fieldAggregation NumericFieldAggregation) Or() NumericResultAggregation {
	return NumericResultAggregation{
		commonAggregation: commonAggregation{
			field:    fieldAggregation.field,
			Function: sql.BitOrAggregation,
		},
	}
}

type BoolFieldAggregation struct {
	FieldAggregation
}

// All returns true if all the values are true
func (fieldAggregation BoolFieldAggregation) All() BoolResultAggregation {
	return BoolResultAggregation{
		commonAggregation: commonAggregation{
			field:    fieldAggregation.field,
			Function: sql.All,
		},
	}
}

// All returns true if at least one value is true
func (fieldAggregation BoolFieldAggregation) Any() BoolResultAggregation {
	return BoolResultAggregation{
		commonAggregation: commonAggregation{
			field:    fieldAggregation.field,
			Function: sql.Any,
		},
	}
}

// None returns true if all values are false
func (fieldAggregation BoolFieldAggregation) None() BoolResultAggregation {
	return BoolResultAggregation{
		commonAggregation: commonAggregation{
			field:    fieldAggregation.field,
			Function: sql.None,
		},
	}
}

type IAggregation interface {
	toSQL(query *GormQuery) (string, error)
	toSelectSQL(query *GormQuery, as string) (string, error)
	getField() IField
}

type commonAggregation struct {
	field    IField
	Function sql.FunctionByDialector
}

func (aggregation commonAggregation) getSQL() toSQLFunc {
	return aggregation.toSQL
}

func (aggregation commonAggregation) getField() IField {
	return aggregation.field
}

func (aggregation commonAggregation) toSQL(query *GormQuery) (string, error) {
	columnSQL := ""

	if aggregation.field != nil { // CountAll
		var err error

		table, err := query.GetModelTable(aggregation.field)
		if err != nil {
			return "", err
		}

		columnSQL = aggregation.field.columnSQL(query, table)
	}

	function, isPresent := aggregation.Function.Get(query.Dialector())
	if !isPresent {
		return "", functionError(ErrUnsupportedByDatabase, aggregation.Function)
	}

	return function.ApplyTo(columnSQL, 0), nil
}

func (aggregation commonAggregation) toSelectSQL(query *GormQuery, as string) (string, error) {
	functionSQL, err := aggregation.toSQL(query)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s AS %s",
		functionSQL,
		as,
	), nil
}

func (aggregation commonAggregation) Eq(
	sqlGeneration func(query *GormQuery) (string, error),
	value any,
) AggregationCondition {
	return AggregationCondition{
		aggregation:   aggregation,
		function:      "=",
		sqlGeneration: sqlGeneration,
		values:        []any{value},
	}
}

type NumericAggregationComparable interface {
	getValue() float64
	getSQL() toSQLFunc
	numericAggregationComparable()
}

type NumericResultAggregation struct {
	commonAggregation
}

func (aggregation NumericResultAggregation) getValue() float64 {
	return 0
}

func (aggregation NumericResultAggregation) numericAggregationComparable() {}

func (aggregation NumericResultAggregation) Eq(value NumericAggregationComparable) AggregationCondition {
	return aggregation.commonAggregation.Eq(value.getSQL(), value.getValue())
}

type BoolAggregationComparable interface {
	getValue() bool

	getSQL() toSQLFunc
	boolAggregationComparable()
}

type BoolResultAggregation struct {
	commonAggregation
}

func (aggregation BoolResultAggregation) Eq(value BoolAggregationComparable) AggregationCondition {
	return aggregation.commonAggregation.Eq(value.getSQL(), value.getValue())
}

func (aggregation BoolResultAggregation) getValue() bool {
	return false
}

func (aggregation BoolResultAggregation) boolAggregationComparable() {}

type AggregationCondition struct {
	aggregation        IAggregation
	function           string
	sqlGeneration      func(query *GormQuery) (string, error)
	values             []any
	conditions         []AggregationCondition
	connectionOperator string
	containerOperator  string
}

func (condition AggregationCondition) toSQL(query *GormQuery) (string, []any, error) {
	if condition.connectionOperator != "" { // connector condition
		sqlStrings := []string{}
		values := []any{}

		for _, internalCondition := range condition.conditions {
			internalSQLString, internalValues, err := internalCondition.toSQL(query)
			if err != nil {
				return "", nil, err
			}

			if internalSQLString != "" {
				sqlStrings = append(sqlStrings, internalSQLString)

				values = append(values, internalValues...)
			}
		}

		return connectSQLs(sqlStrings, condition.connectionOperator), values, nil
	}

	if condition.containerOperator != "" { // container condition
		sqlString, values, err := ConnectionAggregationCondition(condition.conditions, sql.And).toSQL(query)
		if err != nil {
			return "", nil, err
		}

		sqlString = condition.containerOperator + " " + sqlString

		return sqlString, values, nil
	}

	functionSQL, err := condition.aggregation.toSQL(query)
	if err != nil {
		return "", nil, err
	}

	if condition.sqlGeneration != nil {
		sql, err := condition.sqlGeneration(query)
		if err != nil {
			return "", nil, err
		}

		return functionSQL + " " + condition.function + " " + sql, []any{}, nil
	}

	return functionSQL + " " + condition.function + " ?", condition.values, nil
}

func ConnectionAggregationCondition(conditions []AggregationCondition, operator sql.Operator) AggregationCondition {
	return AggregationCondition{
		conditions:         conditions,
		connectionOperator: operator.String(),
	}
}

func ContainerAggregationCondition(conditions []AggregationCondition, operator sql.Operator) AggregationCondition {
	return AggregationCondition{
		conditions:        conditions,
		containerOperator: operator.String(),
	}
}
