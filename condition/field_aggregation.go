package condition

import (
	"fmt"

	"github.com/FrancoLiberali/cql/sql"
)

type FieldAggregation[T any] struct {
	field IField
}

// Count returns the number of values that are not null
func (fieldAggregation FieldAggregation[T]) Count() AggregationResult[float64] {
	return AggregationResult[float64]{
		field:    fieldAggregation.field,
		Function: sql.Count,
	}
}

// Min returns the minimum value of all values
func (fieldAggregation FieldAggregation[T]) Min() AggregationResult[T] {
	return AggregationResult[T]{
		field:    fieldAggregation.field,
		Function: sql.Min,
	}
}

// Max returns the maximum value of all values
func (fieldAggregation FieldAggregation[T]) Max() AggregationResult[T] {
	return AggregationResult[T]{
		field:    fieldAggregation.field,
		Function: sql.Max,
	}
}

type NumericFieldAggregation struct {
	FieldAggregation[float64]
}

// Min returns the minimum value of all values
func (fieldAggregation NumericFieldAggregation) Min() AggregationResult[float64] {
	return fieldAggregation.FieldAggregation.Min()
}

// Max returns the maximum value of all values
func (fieldAggregation NumericFieldAggregation) Max() AggregationResult[float64] {
	return fieldAggregation.FieldAggregation.Max()
}

// Sum calculates the summation of all values
func (fieldAggregation NumericFieldAggregation) Sum() AggregationResult[float64] {
	return AggregationResult[float64]{
		field:    fieldAggregation.field,
		Function: sql.Sum,
	}
}

// Average calculates the average (arithmetic mean) of all values
func (fieldAggregation NumericFieldAggregation) Average() AggregationResult[float64] {
	return AggregationResult[float64]{
		field:    fieldAggregation.field,
		Function: sql.Average,
	}
}

// And calculates the bitwise AND of all non-null values (null values are ignored)
//
// Not available for: sqlite, sqlserver
func (fieldAggregation NumericFieldAggregation) And() AggregationResult[float64] {
	return AggregationResult[float64]{
		field:    fieldAggregation.field,
		Function: sql.BitAndAggregation,
	}
}

// Or calculates the bitwise OR of all non-null values (null values are ignored)
//
// Not available for: sqlite, sqlserver
func (fieldAggregation NumericFieldAggregation) Or() AggregationResult[float64] {
	return AggregationResult[float64]{
		field:    fieldAggregation.field,
		Function: sql.BitOrAggregation,
	}
}

type BoolFieldAggregation struct {
	FieldAggregation[bool]
}

// All returns true if all the values are true
func (fieldAggregation BoolFieldAggregation) All() AggregationResult[bool] {
	return AggregationResult[bool]{
		field:    fieldAggregation.field,
		Function: sql.All,
	}
}

// All returns true if at least one value is true
func (fieldAggregation BoolFieldAggregation) Any() AggregationResult[bool] {
	return AggregationResult[bool]{
		field:    fieldAggregation.field,
		Function: sql.Any,
	}
}

// None returns true if all values are false
func (fieldAggregation BoolFieldAggregation) None() AggregationResult[bool] {
	return AggregationResult[bool]{
		field:    fieldAggregation.field,
		Function: sql.None,
	}
}

type Aggregation interface {
	ToSQL(query *GormQuery) (string, []any, error)
	toSelectSQL(query *GormQuery, as string) (string, error)
	getField() IField
}

type AggregationResult[T any] struct {
	field    IField
	Function sql.FunctionByDialector
}

func (aggregation AggregationResult[T]) getField() IField {
	return aggregation.field
}

func (aggregation AggregationResult[T]) GetValue() T {
	return *new(T)
}

func (aggregation AggregationResult[T]) ToSQL(query *GormQuery) (string, []any, error) {
	columnSQL := ""

	if aggregation.field != nil { // CountAll
		var err error

		table, err := query.GetModelTable(aggregation.field)
		if err != nil {
			return "", nil, err
		}

		columnSQL = aggregation.field.columnSQL(query, table)
	}

	function, isPresent := aggregation.Function.Get(query.Dialector())
	if !isPresent {
		return "", nil, functionError(ErrUnsupportedByDatabase, aggregation.Function)
	}

	return function.ApplyTo(columnSQL, 0), nil, nil
}

func (aggregation AggregationResult[T]) toSelectSQL(query *GormQuery, as string) (string, error) {
	functionSQL, _, err := aggregation.ToSQL(query)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s AS %s",
		functionSQL,
		as,
	), nil
}

func (aggregation AggregationResult[T]) applyOperator(value ValueOfType[T], operator sql.Operator) AggregationCondition {
	return AggregationCondition{
		aggregation:        aggregation,
		operator:           operator,
		valueSQLGeneration: value.ToSQL,
	}
}

// EqualTo
func (aggregation AggregationResult[T]) Eq(value ValueOfType[T]) AggregationCondition {
	return aggregation.applyOperator(value, sql.Eq)
}

// NotEqualTo
func (aggregation AggregationResult[T]) NotEq(value ValueOfType[T]) AggregationCondition {
	return aggregation.applyOperator(value, sql.NotEq)
}

// LessThan
func (aggregation AggregationResult[T]) Lt(value ValueOfType[T]) AggregationCondition {
	return aggregation.applyOperator(value, sql.Lt)
}

// LessThanOrEqualTo
func (aggregation AggregationResult[T]) LtOrEq(value ValueOfType[T]) AggregationCondition {
	return aggregation.applyOperator(value, sql.LtOrEq)
}

// GreaterThan
func (aggregation AggregationResult[T]) Gt(value ValueOfType[T]) AggregationCondition {
	return aggregation.applyOperator(value, sql.Gt)
}

// GreaterThanOrEqualTo
func (aggregation AggregationResult[T]) GtOrEq(value ValueOfType[T]) AggregationCondition {
	return aggregation.applyOperator(value, sql.GtOrEq)
}

func (aggregation AggregationResult[T]) applyListOperator(values []T, operator sql.Operator) AggregationCondition {
	return AggregationCondition{
		aggregation: aggregation,
		operator:    operator,
		valueSQLGeneration: func(*GormQuery) (string, []any, error) {
			return "", []any{values}, nil
		},
	}
}

func (aggregation AggregationResult[T]) In(values []T) AggregationCondition {
	return aggregation.applyListOperator(values, sql.ArrayIn)
}

func (aggregation AggregationResult[T]) NotIn(values []T) AggregationCondition {
	return aggregation.applyListOperator(values, sql.ArrayNotIn)
}

// Pattern in all databases:
//   - An underscore (_) in pattern stands for (matches) any single character.
//   - A percent sign (%) matches any sequence of zero or more characters.
//
// Additionally in SQLServer:
//   - Square brackets ([ ]) matches any single character within the specified range ([a-f]) or set ([abcdef]).
//   - [^] matches any single character not within the specified range ([^a-f]) or set ([^abcdef]).
//
// WARNINGS:
//   - SQLite: LIKE is case-insensitive unless case_sensitive_like pragma (https://www.sqlite.org/pragma.html#pragma_case_sensitive_like) is true.
//   - SQLServer, MySQL: the case-sensitivity depends on the collation used in compared column.
//   - PostgreSQL: LIKE is always case-sensitive
//
// refs:
//   - mysql: https://dev.mysql.com/doc/refman/8.0/en/string-comparison-functions.html#operator_like
//   - postgresql: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-LIKE
//   - sqlserver: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/like-transact-sql?view=sql-server-ver16
//   - sqlite: https://www.sqlite.org/lang_expr.html#like
func (aggregation AggregationResult[T]) Like(value string) AggregationCondition {
	return AggregationCondition{
		aggregation: aggregation,
		operator:    sql.Like,
		valueSQLGeneration: func(*GormQuery) (string, []any, error) {
			return "", []any{value}, nil
		},
	}
}

type AggregationCondition struct {
	aggregation        Aggregation
	operator           sql.Operator
	valueSQLGeneration func(*GormQuery) (string, []any, error)
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

	functionSQL, _, err := condition.aggregation.ToSQL(query)
	if err != nil {
		return "", nil, err
	}

	sql, values, err := condition.valueSQLGeneration(query)
	if err != nil {
		return "", nil, err
	}

	if sql != "" {
		return functionSQL + " " + condition.operator.String() + " " + sql, []any{}, nil
	}

	return functionSQL + " " + condition.operator.String() + " ?", values, nil
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
